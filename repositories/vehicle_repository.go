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

func (r *VehicleRepository) Create(vehicle *models.Vehicle) error {
	return r.db.Create(vehicle).Error
}

func (r *VehicleRepository) FindByID(id uint) (*models.Vehicle, error) {
	var user models.Vehicle
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *VehicleRepository) Update(user *models.Vehicle) error {
	return r.db.Save(user).Error
}

func (r *VehicleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Vehicle{}, id).Error
}
