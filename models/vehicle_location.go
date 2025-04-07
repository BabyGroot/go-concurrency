package models

import "time"

type VehicleLocation struct {
	ID         uint `gorm:"primarykey"`
	VehicleID  uint
	LocationID uint
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Add relationships (optional but helpful)
	Vehicle  Vehicle  `gorm:"foreignKey:VehicleID"`
	Location Location `gorm:"foreignKey:LocationID"`
}
