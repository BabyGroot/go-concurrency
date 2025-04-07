package repositories

import (
	"fmt"
	"myproject/database"
	"myproject/models"
	"strings"
	"sync"
)

type VehicleRepository struct {
	db *database.Database
}

func NewVehicleRepository(db *database.Database) *VehicleRepository {
	return &VehicleRepository{db: db}
}

func (r *VehicleRepository) GetAll() ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	result := r.db.Find(&vehicles)
	if result.Error != nil {
		return nil, result.Error
	}
	return vehicles, nil
}

func (r *VehicleRepository) GetAllPaginated(page, pageSize int) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	offset := (page - 1) * pageSize

	result := r.db.Offset(offset).Limit(pageSize).Find(&vehicles)
	if result.Error != nil {
		return nil, result.Error
	}
	return vehicles, nil
}

func (r *VehicleRepository) Update(user *models.Vehicle) error {
	return r.db.Save(user).Error
}

func (r *VehicleRepository) BatchUpdatePermalinks(vehicles []models.Vehicle) error {
	// Channels for collecting results
	errorChan := make(chan error, len(vehicles))
	semaphore := make(chan struct{}, 5) // Limit concurrency to 5 goroutines
	var wg sync.WaitGroup

	for _, vehicle := range vehicles {
		wg.Add(1)
		go func(v models.Vehicle) {
			defer wg.Done()

			// Acquire semaphore slot
			semaphore <- struct{}{}
			defer func() { <-semaphore }() // Release slot when done

			// Generate permalink if empty
			if v.Permalink == "" {
				v.Permalink = strings.ToLower(strings.ReplaceAll(v.Name, " ", "-"))

				// Update in database
				if err := r.Update(&v); err != nil {
					errorChan <- err
				}
			}
		}(vehicle)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errorChan)

	// Collect any errors
	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("batch update had %d errors: %v", len(errors), errors)
	}

	return nil
}

func (r *VehicleRepository) GetVehiclesByLocation(locationID uint) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle

	result := r.db.Joins("JOIN vehicle_locations ON vehicles.id = vehicle_locations.vehicle_id").
		Where("vehicle_locations.location_id = ?", locationID).
		Find(&vehicles)

	if result.Error != nil {
		return nil, result.Error
	}

	return vehicles, nil
}

func (r *VehicleRepository) GetVehiclesByLocationName(locationName string, locationRepo *LocationRepository) ([]models.Vehicle, error) {
	// First find the location by name
	location, err := locationRepo.FindByPermalink(locationName)
	if err != nil {
		return nil, err
	}

	// Then get vehicles at this location
	return r.GetVehiclesByLocation(location.ID)
}

func (r *VehicleRepository) GetVehicleStatsByLocation(locationID uint) (map[string]int, error) {
	// Example implementation that counts vehicles by type/category
	// First, make sure you have a "type" or "category" field in your Vehicle model
	// If not, you can use other fields for statistics

	stats := make(map[string]int)

	// Option 1: If you have a vehicle_types table or field
	// This query counts vehicles by type
	rows, err := r.db.Raw(`
        SELECT v.class as vehicle_class, COUNT(*) as count
        FROM vehicles v
        JOIN vehicle_locations vl ON v.id = vl.vehicle_id
        WHERE vl.location_id = ?
        GROUP BY v.class
    `, locationID).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var vehicleClass string
		var count int
		if err := rows.Scan(&vehicleClass, &count); err != nil {
			return nil, err
		}
		stats[vehicleClass] = count
	}

	// If we don't have vehicle types, we can provide other useful statistics
	// Add total vehicle count
	var totalCount int64
	if err := r.db.Model(&models.VehicleLocation{}).
		Where("location_id = ?", locationID).
		Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats["total"] = int(totalCount)

	return stats, nil
}
