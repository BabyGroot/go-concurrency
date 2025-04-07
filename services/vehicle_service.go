package services

import (
	"myproject/models"
	"myproject/repositories"
)

type VehicleService struct {
	vehicleRepo  *repositories.VehicleRepository
	locationRepo *repositories.LocationRepository
}

func NewVehicleService(vehicleRepo *repositories.VehicleRepository, locationRepo *repositories.LocationRepository) *VehicleService {
	return &VehicleService{
		vehicleRepo:  vehicleRepo,
		locationRepo: locationRepo,
	}
}

func (s *VehicleService) GetVehiclesByLocationName(locationName string) ([]models.Vehicle, error) {
	location, err := s.locationRepo.FindByPermalink(locationName)
	if err != nil {
		return nil, err
	}

	// Then get vehicles at this location
	return s.vehicleRepo.GetVehiclesByLocation(location.ID)
}
