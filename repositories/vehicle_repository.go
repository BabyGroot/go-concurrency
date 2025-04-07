package repositories

import (
	"myproject/database"
	"myproject/models"
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
