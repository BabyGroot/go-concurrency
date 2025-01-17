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
