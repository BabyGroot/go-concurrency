package models

import "time"

type Vehicle struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"size:255;not null"`
	//StartDate  date `gorm:"size:255;not null;unique"`
	Permalink string `gorm:"size:255;null"`
}
