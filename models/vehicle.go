package models

import "gorm.io/gorm"

type User struct {
	gorm.Model        // Adds ID, CreatedAt, UpdatedAt, DeletedAt fields
	Name       string `gorm:"size:255;not null"`
	// Email    string    `gorm:"size:255;not null;unique"`
	// Age      int       `gorm:"not null"`
	// Profile  Profile   `gorm:"foreignKey:UserID"`
}

type Profile struct {
	gorm.Model
	UserID  uint   `gorm:"not null"`
	Address string `gorm:"size:255"`
	Phone   string `gorm:"size:20"`
}
