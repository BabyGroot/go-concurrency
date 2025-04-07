package repositories

import (
	"fmt"
	"myproject/database"
	"myproject/models"
)

type LocationRepository struct {
	db *database.Database
}

func NewLocationRepository(db *database.Database) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) FindByPermalink(permalink string) (*models.Location, error) {
	var location models.Location
	result := r.db.Where("permalink = ?", permalink).First(&location)

	if result.Error != nil {
		return nil, result.Error
	}

	return &location, nil
}

func (r *LocationRepository) GetAll() ([]models.Location, error) {
	var vehicles []models.Location
	result := r.db.Find(&vehicles)
	if result.Error != nil {
		return nil, result.Error
	}
	return vehicles, nil
}

func (r *LocationRepository) GetAllWithVehicleCounts() ([]map[string]interface{}, error) {
	// Get all locations first
	locations, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	// Fan out to get vehicle counts for each location
	type locationResult struct {
		Location models.Location
		Count    int
		Error    error
	}

	resultChan := make(chan locationResult, len(locations))

	// Limit concurrency to avoid database connection issues
	maxConcurrency := 5
	sem := make(chan struct{}, maxConcurrency)

	for _, loc := range locations {
		go func(l models.Location) {
			sem <- struct{}{}        // Acquire token
			defer func() { <-sem }() // Release token

			// Count vehicles at this location
			var count int64
			err := r.db.Model(&models.VehicleLocation{}).
				Where("location_id = ?", l.ID).
				Count(&count).Error

			resultChan <- locationResult{
				Location: l,
				Count:    int(count),
				Error:    err,
			}
		}(loc)
	}

	// Fan in results
	results := make([]map[string]interface{}, 0, len(locations))
	var resultErrors []error

	for i := 0; i < len(locations); i++ {
		result := <-resultChan

		if result.Error != nil {
			resultErrors = append(resultErrors, result.Error)
			continue
		}

		results = append(results, map[string]interface{}{
			"permalink": result.Location.Permalink,
			"name":      result.Location.Name,
			"count":     result.Count,
		})
	}

	if len(resultErrors) > 0 {
		return results, fmt.Errorf("errors occurred while getting vehicle counts: %v", resultErrors)
	}

	return results, nil
}

func (r *LocationRepository) GetAllPaginated(page, pageSize int) ([]models.Location, error) {
	var vehicles []models.Location
	offset := (page - 1) * pageSize

	result := r.db.Offset(offset).Limit(pageSize).Find(&vehicles)
	if result.Error != nil {
		return nil, result.Error
	}
	return vehicles, nil
}

func (r *LocationRepository) Update(user *models.Location) error {
	return r.db.Save(user).Error
}
