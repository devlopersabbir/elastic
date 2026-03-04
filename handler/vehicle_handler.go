package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/devloeprsabbir/go-elasticsearch/elastic"
	"github.com/devloeprsabbir/go-elasticsearch/models"
	"github.com/devloeprsabbir/go-elasticsearch/repository"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/gin-gonic/gin"
)

// VehicleHandler handles all vehicle-related HTTP routes.
type VehicleHandler struct {
	repo *repository.VehicleRepository
	es   *elasticsearch.Client
}

func NewVehicleHandler(repo *repository.VehicleRepository, es *elasticsearch.Client) *VehicleHandler {
	return &VehicleHandler{repo: repo, es: es}
}

// SearchVehicles godoc
//
// POST /api/v1/vehicles/search
//
// Request body:
//
//	{
//	  "listing_url": "https://...",
//	  "page": 1,
//	  "page_size": 15
//	}
//
// The handler:
//  1. Looks up the reference vehicle by listing_url (tries Elasticsearch first, PG fallback)
//  2. Runs a weighted similarity search across 15+ fields in Elasticsearch
//  3. Returns the top 10-15 most similar vehicles with pagination
func (h *VehicleHandler) SearchVehicles(c *gin.Context) {
	var req models.VehicleSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 50 {
		req.PageSize = 15
	}

	ctx := c.Request.Context()

	// Step 1: Resolve the reference vehicle from the listing URL
	refVehicle, err := resolveReferenceVehicle(ctx, h, req.ListingURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "reference vehicle not found",
			"details": err.Error(),
		})
		return
	}

	// Step 2: Search for similar vehicles in Elasticsearch
	results, total, err := elastic.SearchSimilarVehicles(ctx, refVehicle, req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.VehicleSearchResponse{
		Total:        total,
		Page:         req.Page,
		PageSize:     req.PageSize,
		QueryVehicle: refVehicle,
		Results:      results,
	})
}

// resolveReferenceVehicle tries to get the reference vehicle.
// It first tries Elasticsearch (fast), then falls back to PostgreSQL (complete).
func resolveReferenceVehicle(ctx context.Context, h *VehicleHandler, listingURL string) (*models.Vehicle, error) {
	// Try Elasticsearch first
	v, err := elastic.GetVehicleByListingURL(ctx, listingURL)
	if err == nil && v != nil {
		return v, nil
	}
	log.Printf("listing_url not found in Elasticsearch (%v), falling back to PostgreSQL", err)

	// Fallback to PostgreSQL
	return h.repo.GetByListingURL(ctx, listingURL)
}

// TriggerIndexing godoc
// POST /api/v1/vehicles/index
//
// Triggers a background re-indexing of all vehicle records from PostgreSQL → Elasticsearch.
// Returns immediately with 202 Accepted; indexing runs asynchronously.
// Progress is logged to the server log.
func (h *VehicleHandler) TriggerIndexing(c *gin.Context) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		count, err := h.repo.CountVehicles(ctx)
		if err != nil {
			log.Printf("TriggerIndexing: count error: %v", err)
			return
		}
		log.Printf("TriggerIndexing: starting indexing of %d vehicles", count)

		err = h.repo.FetchAllForIndexing(ctx, 500, func(batch []models.Vehicle) error {
			return elastic.BulkIndexVehicles(batch, 500)
		})
		if err != nil {
			log.Printf("TriggerIndexing: error during indexing: %v", err)
		} else {
			log.Printf("TriggerIndexing: completed successfully (%d vehicles)", count)
		}
	}()

	c.JSON(http.StatusAccepted, models.IndexTriggerResponse{
		Message:   "Indexing started in background. Check server logs for progress.",
		Scheduled: true,
	})
}

// GetVehicleStatus godoc
// GET /api/v1/vehicles/status
//
// Returns current count of indexed vehicles in Elasticsearch vs PostgreSQL.
func (h *VehicleHandler) GetVehicleStatus(c *gin.Context) {
	ctx := c.Request.Context()

	pgCount, err := h.repo.CountVehicles(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "pg count failed: " + err.Error()})
		return
	}

	var esCount int64
	esCountRes, esErr := elastic.EsClient.Count(
		elastic.EsClient.Count.WithIndex("vehicles"),
	)
	if esErr == nil && !esCountRes.IsError() {
		body, readErr := io.ReadAll(esCountRes.Body)
		esCountRes.Body.Close()
		if readErr == nil {
			var cr map[string]interface{}
			if json.Unmarshal(body, &cr) == nil {
				if v, ok := cr["count"].(float64); ok {
					esCount = int64(v)
				}
			}
		}
	} else if esCountRes != nil {
		esCountRes.Body.Close()
	}

	var indexedPct float64
	if pgCount > 0 {
		indexedPct = float64(esCount) / float64(pgCount) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"postgres_count":  pgCount,
		"elastic_count":   esCount,
		"indexed_percent": indexedPct,
	})
}
