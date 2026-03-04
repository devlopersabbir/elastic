package elastic

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	config "github.com/devloeprsabbir/go-elasticsearch/config"
	"github.com/elastic/go-elasticsearch/v9"
)

var EsClient *elasticsearch.Client

// ElasticClient initialises the Elasticsearch client and creates the vehicles index.
func ElasticClient(cfg *config.Config) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.ElasticAddresses(),
	})
	if err != nil {
		log.Fatalf("failed to connect to elasticsearch: %v", err)
	}

	info, err := es.Info()
	if err != nil {
		log.Fatalf("elasticsearch info call failed: %v", err)
	}
	defer info.Body.Close()
	if info.IsError() {
		log.Fatalf("elasticsearch info error: %s", info.String())
	}
	log.Println("Elasticsearch connection OK:", strings.TrimRight(info.String(), "\n"))

	EsClient = es
	createVehicleIndex(es)
}

// vehicleIndexMapping defines the Elasticsearch index mappings for the vehicle_data table.
// Key fields that drive similarity are mapped with keyword sub-fields for exact-match
// aggregations and a text field for full-text search.
var vehicleIndexMapping = map[string]interface{}{
	"settings": map[string]interface{}{
		"number_of_shards":   2,
		"number_of_replicas": 1,
		"analysis": map[string]interface{}{
			"analyzer": map[string]interface{}{
				"lowercase_analyzer": map[string]interface{}{
					"type":      "custom",
					"tokenizer": "keyword",
					"filter":    []string{"lowercase", "trim"},
				},
			},
		},
	},
	"mappings": map[string]interface{}{
		"properties": buildVehicleMappingProperties(),
	},
}

func buildVehicleMappingProperties() map[string]interface{} {
	// Helper to create a text field that also has a keyword sub-field
	textKeyword := func() map[string]interface{} {
		return map[string]interface{}{
			"type": "text",
			"fields": map[string]interface{}{
				"keyword": map[string]interface{}{
					"type":         "keyword",
					"ignore_above": 256,
				},
			},
		}
	}
	keyword := func() map[string]interface{} {
		return map[string]interface{}{"type": "keyword"}
	}
	boolField := func() map[string]interface{} {
		return map[string]interface{}{"type": "boolean"}
	}

	return map[string]interface{}{
		// --- Core identifiers ---
		"unique_id":            keyword(),
		"vehicle_id":           keyword(),
		"data_source":          keyword(),
		"listing_url":          keyword(),
		"is_vehicle_available": boolField(),

		// --- Key similarity fields: make, model, fuel, transmission, color, mileage ---
		"make":          textKeyword(),
		"model":         textKeyword(),
		"model_version": textKeyword(),
		"model_range":   textKeyword(),
		"trim_line":     textKeyword(),
		"title":         textKeyword(),
		"vehicle_type":  keyword(),
		"category":      keyword(),
		"body_type":     keyword(),

		// Engine / powertrain
		"fuel_type":         keyword(),
		"fuel_category":     keyword(),
		"primary_fuel":      keyword(),
		"engine_type":       keyword(),
		"transmission":      keyword(),
		"transmission_type": keyword(),
		"drive_train":       keyword(),
		"power_kw":          textKeyword(),
		"power_hp":          textKeyword(),
		"displacement_ccm":  textKeyword(),
		"cylinders":         keyword(),
		"gears":             keyword(),

		// Mileage / age
		"mileage_km":              textKeyword(),
		"mileage_in_km":           textKeyword(),
		"first_registration":      keyword(),
		"first_registration_date": keyword(),
		"production_year":         keyword(),
		"construction_year":       keyword(),

		// Price
		"price": textKeyword(),

		// Physical
		"color":              keyword(),
		"color_original":     keyword(),
		"manufacturer_color": keyword(),
		"paint_type":         keyword(),
		"seats":              keyword(),
		"doors":              keyword(),
		"upholstery":         keyword(),
		"upholstery_color":   keyword(),
		"interior":           keyword(),
		"interior_color":     keyword(),
		"interior_type":      keyword(),

		// Condition / history
		"condition":        keyword(),
		"damage_condition": keyword(),
		"had_accident":     keyword(),
		"previous_owners":  keyword(),
		"offer_type":       keyword(),

		// Location
		"country_code": keyword(),
		"city":         textKeyword(),
		"seller_name":  textKeyword(),

		// Description (full-text only)
		"description": map[string]interface{}{"type": "text"},

		// Boolean feature flags for similarity scoring
		"abs":                     boolField(),
		"esp":                     boolField(),
		"navigation_system":       boolField(),
		"bluetooth":               boolField(),
		"sunroof":                 boolField(),
		"panoramic_roof":          boolField(),
		"heated_seats":            boolField(),
		"leather_interior":        boolField(),
		"all_wheel_drive":         boolField(),
		"cruise_control":          boolField(),
		"adaptive_cruise_control": boolField(),
		"parking_camera":          boolField(),
		"led_headlights":          boolField(),
		"touchscreen":             boolField(),
		"apple_carplay":           boolField(),
		"android_auto":            boolField(),
		"start_stop_system":       boolField(),
		"warranty":                boolField(),
	}
}

func createVehicleIndex(es *elasticsearch.Client) {
	// Check if the index already exists
	existRes, err := es.Indices.Exists([]string{"vehicles"})
	if err != nil {
		log.Printf("could not check vehicles index existence: %v", err)
		return
	}
	defer existRes.Body.Close()
	if existRes.StatusCode == 200 {
		log.Println("Elasticsearch index 'vehicles' already exists — skipping creation")
		return
	}

	mappingBytes, _ := json.Marshal(vehicleIndexMapping)
	res, err := es.Indices.Create(
		"vehicles",
		es.Indices.Create.WithBody(bytes.NewReader(mappingBytes)),
	)
	if err != nil {
		log.Fatalf("failed creating vehicles index: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("create vehicles index error: %s", res.String())
	} else {
		log.Println("Elasticsearch index 'vehicles' created successfully")
	}
}
