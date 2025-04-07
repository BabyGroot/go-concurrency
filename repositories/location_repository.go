package repositories

import (
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
