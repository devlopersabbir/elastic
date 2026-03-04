package main

import (
	"log"

	"github.com/devloeprsabbir/go-elasticsearch/config"
	"github.com/devloeprsabbir/go-elasticsearch/elastic"
	handler "github.com/devloeprsabbir/go-elasticsearch/handler"
	"github.com/devloeprsabbir/go-elasticsearch/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.LoadEnv()

	// Connect to PostgreSQL (Azure)
	config.ConnectDatabase(cfg)

	// Connect to Elasticsearch & create index
	elastic.ElasticClient(cfg)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(corsMiddleware())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// ---------- Vehicle Routes ----------
	vehicleRepo := repository.NewVehicleRepository(config.DB)
	vehicleHandler := handler.NewVehicleHandler(vehicleRepo, elastic.EsClient)

	v1 := r.Group("/api/v1")
	{
		vehicles := v1.Group("/vehicles")
		{
			// POST /api/v1/vehicles/search
			// Body: { "listing_url": "...", "page": 1, "page_size": 15 }
			// Returns: top 10–15 most similar vehicles from Elasticsearch
			vehicles.POST("/search", vehicleHandler.SearchVehicles)

			// POST /api/v1/vehicles/index
			// Triggers async bulk-indexing of all vehicles from PG → Elasticsearch
			vehicles.POST("/index", vehicleHandler.TriggerIndexing)

			// GET /api/v1/vehicles/status
			// Returns PG count vs ES indexed count
			vehicles.GET("/status", vehicleHandler.GetVehicleStatus)
		}
	}

	log.Printf("Vehicle Search Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
