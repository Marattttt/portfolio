package models

import "time"

type Guest struct {
	ID        uint `gorm:"column:guest_id"`
	Name      string
	Salt      *string
	Secret    string
	CreatedAt time.Time

	Visits []Visit
}
