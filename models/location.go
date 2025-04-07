package models

import "time"

type Location struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"size:255;not null"`
	Permalink string `gorm:"size:255;null"`
}
