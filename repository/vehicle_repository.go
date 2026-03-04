package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/devloeprsabbir/go-elasticsearch/models"
	"gorm.io/gorm"
)

// VehicleRepository handles all PostgreSQL operations for vehicle_data.
type VehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) *VehicleRepository {
	return &VehicleRepository{db: db}
}

// GetByListingURL retrieves a single vehicle record by its listing URL.
func (r *VehicleRepository) GetByListingURL(ctx context.Context, listingURL string) (*models.Vehicle, error) {
	var v models.Vehicle
	err := r.db.WithContext(ctx).
		Where("listing_url = ?", listingURL).
		First(&v).Error
	if err != nil {
		return nil, fmt.Errorf("vehicle not found for listing_url %q: %w", listingURL, err)
	}
	return &v, nil
}

// GetByUniqueIDs fetches vehicles by their unique_id, preserving Elasticsearch score order.
func (r *VehicleRepository) GetByUniqueIDs(ctx context.Context, ids []string) ([]models.Vehicle, error) {
	if len(ids) == 0 {
		return []models.Vehicle{}, nil
	}
	var vehicles []models.Vehicle
	err := r.db.WithContext(ctx).
		Where("unique_id IN ?", ids).
		Find(&vehicles).Error
	if err != nil {
		return nil, err
	}

	vehicleMap := make(map[string]models.Vehicle, len(vehicles))
	for _, v := range vehicles {
		vehicleMap[v.UniqueID] = v
	}

	ordered := make([]models.Vehicle, 0, len(ids))
	for _, id := range ids {
		if v, ok := vehicleMap[id]; ok {
			ordered = append(ordered, v)
		}
	}
	return ordered, nil
}

// FetchAllForIndexing streams all vehicle records from PostgreSQL in batches
// and calls the provided callback for each batch. This avoids loading 24k+ records
// into memory at once.
func (r *VehicleRepository) FetchAllForIndexing(ctx context.Context, batchSize int, callback func([]models.Vehicle) error) error {
	if batchSize <= 0 {
		batchSize = 1000
	}

	var offset int
	for {
		var batch []models.Vehicle
		err := r.db.WithContext(ctx).
			Offset(offset).
			Limit(batchSize).
			Order("unique_id").
			Find(&batch).Error
		if err != nil {
			return fmt.Errorf("fetch batch at offset %d: %w", offset, err)
		}
		if len(batch) == 0 {
			break
		}

		log.Printf("Fetched PG batch offset=%d count=%d", offset, len(batch))
		if err := callback(batch); err != nil {
			return fmt.Errorf("indexing callback error at offset %d: %w", offset, err)
		}

		if len(batch) < batchSize {
			break // last page
		}
		offset += batchSize
	}
	return nil
}

// CountVehicles returns total number of vehicles in the table.
func (r *VehicleRepository) CountVehicles(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Vehicle{}).Count(&count).Error
	return count, err
}
