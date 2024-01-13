package models

import "time"

type Guest struct {
	ID        uint   `gorm:"column:guest_id"`
	Name      string `gorm:"type:bytea"`
	Salt      string `gorm:"type:bytea"`
	Secret    string
	CreatedAt time.Time

	Visits []Visit
}
