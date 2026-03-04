package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/devloeprsabbir/go-elasticsearch/models"
)

const vehicleIndex = "vehicles"

// BulkIndexVehicles indexes vehicles in batches using Elasticsearch's bulk API.
// batchSize controls how many documents are sent per HTTP request.
func BulkIndexVehicles(vehicles []models.Vehicle, batchSize int) error {
	if EsClient == nil {
		return fmt.Errorf("elasticsearch client is not initialised")
	}
	if batchSize <= 0 {
		batchSize = 500
	}

	total := len(vehicles)
	indexed := 0

	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		batch := vehicles[i:end]

		var buf bytes.Buffer
		for _, v := range batch {
			// Action line
			meta := map[string]interface{}{
				"index": map[string]interface{}{
					"_index": vehicleIndex,
					"_id":    v.UniqueID,
				},
			}
			metaBytes, err := json.Marshal(meta)
			if err != nil {
				log.Printf("marshal meta error for %s: %v", v.UniqueID, err)
				continue
			}
			buf.Write(metaBytes)
			buf.WriteByte('\n')

			// Document line
			docBytes, err := json.Marshal(v)
			if err != nil {
				log.Printf("marshal doc error for %s: %v", v.UniqueID, err)
				continue
			}
			buf.Write(docBytes)
			buf.WriteByte('\n')
		}

		res, err := EsClient.Bulk(
			bytes.NewReader(buf.Bytes()),
			EsClient.Bulk.WithIndex(vehicleIndex),
			EsClient.Bulk.WithRefresh("false"), // don't refresh per batch — refresh at end
		)
		if err != nil {
			return fmt.Errorf("bulk request failed at batch %d-%d: %w", i, end, err)
		}

		var bulkResp map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&bulkResp); err != nil {
			res.Body.Close()
			return fmt.Errorf("bulk response decode error: %w", err)
		}
		res.Body.Close()

		// Count errors in bulk response
		errCount := 0
		if items, ok := bulkResp["items"].([]interface{}); ok {
			for _, item := range items {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if idx, ok := itemMap["index"].(map[string]interface{}); ok {
						if status, ok := idx["status"].(float64); ok && status >= 400 {
							errCount++
						}
					}
				}
			}
		}

		indexed += len(batch) - errCount
		log.Printf("bulk indexed batch %d-%d: %d ok, %d errors", i, end, len(batch)-errCount, errCount)
	}

	// Force a refresh so documents are immediately searchable
	refreshRes, err := EsClient.Indices.Refresh(EsClient.Indices.Refresh.WithIndex(vehicleIndex))
	if err != nil {
		log.Printf("post-bulk refresh error: %v", err)
	} else {
		refreshRes.Body.Close()
	}

	log.Printf("BulkIndexVehicles complete: %d/%d documents indexed", indexed, total)
	return nil
}

// IndexVehicle indexes (or re-indexes) a single vehicle document.
func IndexVehicle(v *models.Vehicle) error {
	if EsClient == nil {
		return fmt.Errorf("elasticsearch client is not initialised")
	}

	docBytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	res, err := EsClient.Index(
		vehicleIndex,
		bytes.NewReader(docBytes),
		EsClient.Index.WithDocumentID(v.UniqueID),
		EsClient.Index.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("index request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("indexing error: %s", res.String())
	}
	return nil
}

// GetVehicleByListingURL fetches an exact vehicle document from Elasticsearch by listing_url.
func GetVehicleByListingURL(ctx context.Context, listingURL string) (*models.Vehicle, error) {
	if EsClient == nil {
		return nil, fmt.Errorf("elasticsearch client is not initialised")
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"listing_url": listingURL,
			},
		},
		"size": 1,
	}
	queryBytes, _ := json.Marshal(query)

	res, err := EsClient.Search(
		EsClient.Search.WithContext(ctx),
		EsClient.Search.WithIndex(vehicleIndex),
		EsClient.Search.WithBody(bytes.NewReader(queryBytes)),
	)
	if err != nil {
		return nil, fmt.Errorf("search by listing_url failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	hitsMap, _ := r["hits"].(map[string]interface{})
	hits, _ := hitsMap["hits"].([]interface{})
	if len(hits) == 0 {
		return nil, fmt.Errorf("vehicle not found for listing_url: %s", listingURL)
	}

	hitMap, _ := hits[0].(map[string]interface{})
	source, _ := hitMap["_source"].(map[string]interface{})
	sourceBytes, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	var v models.Vehicle
	if err := json.Unmarshal(sourceBytes, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// SearchSimilarVehicles returns the best matching vehicles for a given reference vehicle
// using a weighted bool query across make, model, mileage range, fuel type, transmission,
// color, body type, year, and optional feature flags.
func SearchSimilarVehicles(ctx context.Context, ref *models.Vehicle, page, pageSize int) ([]models.Vehicle, int64, error) {
	if EsClient == nil {
		return nil, 0, fmt.Errorf("elasticsearch client is not initialised")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 15
	}

	// --- Build must clauses (hard filters) ---
	mustClauses := []interface{}{}

	// Always exclude the reference vehicle itself
	mustClauses = append(mustClauses, map[string]interface{}{
		"bool": map[string]interface{}{
			"must_not": map[string]interface{}{
				"term": map[string]interface{}{"listing_url": ref.ListingURL},
			},
		},
	})

	// Only search available vehicles
	mustClauses = append(mustClauses, map[string]interface{}{
		"term": map[string]interface{}{"is_vehicle_available": true},
	})

	// --- Build should clauses (weighted similarity) ---
	shouldClauses := []interface{}{}

	// Make — highest weight
	if ref.Make != "" {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"match": map[string]interface{}{
				"make": map[string]interface{}{
					"query": ref.Make,
					"boost": 10,
				},
			},
		})
	}

	// Model — very important
	if ref.Model != "" {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"match": map[string]interface{}{
				"model": map[string]interface{}{
					"query": ref.Model,
					"boost": 8,
				},
			},
		})
	}

	// Model version / trim
	if ref.ModelVersion != "" {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"match": map[string]interface{}{
				"model_version": map[string]interface{}{
					"query": ref.ModelVersion,
					"boost": 4,
				},
			},
		})
	}

	// Body type
	if ref.BodyType != "" {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"body_type": map[string]interface{}{
					"value": strings.ToLower(ref.BodyType),
					"boost": 5,
				},
			},
		})
	}

	// Fuel type — important for running costs
	if ref.FuelType != "" {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"fuel_type": map[string]interface{}{
					"value": strings.ToLower(ref.FuelType),
					"boost": 6,
				},
			},
		})
	}

	// Transmission
	if ref.Transmission != "" {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"transmission": map[string]interface{}{
					"value": strings.ToLower(ref.Transmission),
					"boost": 5,
				},
			},
		})
	}

	// Drive train
	if ref.DriveTrain != "" {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"drive_train": map[string]interface{}{
					"value": strings.ToLower(ref.DriveTrain),
					"boost": 3,
				},
			},
		})
	}

	// Color
	if ref.Color != "" {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"color": map[string]interface{}{
					"value": strings.ToLower(ref.Color),
					"boost": 3,
				},
			},
		})
	}

	// Engine type
	if ref.EngineType != "" {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"match": map[string]interface{}{
				"engine_type": map[string]interface{}{
					"query": ref.EngineType,
					"boost": 4,
				},
			},
		})
	}

	// Power HP — fuzzy range: ±30% of the reference power
	if ref.PowerHP != "" {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"match": map[string]interface{}{
				"power_hp": map[string]interface{}{
					"query": ref.PowerHP,
					"boost": 3,
				},
			},
		})
	}

	// Condition
	if ref.Condition != "" {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"condition": map[string]interface{}{
					"value": strings.ToLower(ref.Condition),
					"boost": 4,
				},
			},
		})
	}

	// Production / construction year
	if ref.ProductionYear != "" {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"production_year": map[string]interface{}{
					"value": ref.ProductionYear,
					"boost": 4,
				},
			},
		})
	} else if ref.ConstructionYear != "" {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"construction_year": map[string]interface{}{
					"value": ref.ConstructionYear,
					"boost": 4,
				},
			},
		})
	}

	// Mileage similarity — match same 10k range bucket using mileage_km text match
	if ref.MileageKM != "" {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"match": map[string]interface{}{
				"mileage_km": map[string]interface{}{
					"query": ref.MileageKM,
					"boost": 5,
				},
			},
		})
	}

	// Boolean feature flags — each adds a small boost when present in both
	flags := map[string]**bool{
		"navigation_system":       &ref.NavigationSystem,
		"bluetooth":               &ref.Bluetooth,
		"sunroof":                 &ref.Sunroof,
		"panoramic_roof":          &ref.PanoramicRoof,
		"heated_seats":            &ref.HeatedSeats,
		"leather_interior":        &ref.LeatherInterior,
		"all_wheel_drive":         &ref.AllWheelDrive,
		"cruise_control":          &ref.CruiseControl,
		"adaptive_cruise_control": &ref.AdaptiveCruiseControl,
		"parking_camera":          &ref.ParkingCamera,
		"led_headlights":          &ref.LedHeadlights,
		"apple_carplay":           &ref.AppleCarPlay,
		"android_auto":            &ref.AndroidAuto,
		"warranty":                &ref.Warranty,
	}
	for field, ptr := range flags {
		if *ptr != nil && **ptr {
			shouldClauses = append(shouldClauses, map[string]interface{}{
				"term": map[string]interface{}{
					field: map[string]interface{}{
						"value": true,
						"boost": 1,
					},
				},
			})
		}
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":                 mustClauses,
				"should":               shouldClauses,
				"minimum_should_match": 1,
			},
		},
		"from": (page - 1) * pageSize,
		"size": pageSize,
	}

	queryBytes, _ := json.Marshal(query)
	res, err := EsClient.Search(
		EsClient.Search.WithContext(ctx),
		EsClient.Search.WithIndex(vehicleIndex),
		EsClient.Search.WithBody(bytes.NewReader(queryBytes)),
		EsClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("elasticsearch search failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("search error: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, 0, err
	}

	hitsMap, _ := r["hits"].(map[string]interface{})
	hits, _ := hitsMap["hits"].([]interface{})

	var total int64
	if totalMap, ok := hitsMap["total"].(map[string]interface{}); ok {
		if val, ok := totalMap["value"].(float64); ok {
			total = int64(val)
		}
	}

	vehicles := make([]models.Vehicle, 0, len(hits))
	for _, hit := range hits {
		hitMap, _ := hit.(map[string]interface{})
		source, _ := hitMap["_source"].(map[string]interface{})
		srcBytes, err := json.Marshal(source)
		if err != nil {
			continue
		}
		var v models.Vehicle
		if err := json.Unmarshal(srcBytes, &v); err != nil {
			continue
		}
		vehicles = append(vehicles, v)
	}

	return vehicles, total, nil
}
