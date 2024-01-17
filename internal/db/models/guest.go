package models

import "time"

type Guest struct {
	ID        uint `gorm:"column:guest_id"`
	Name      string
	Salt      []byte `gorm:"type:bytea"`
	Secret    []byte `gorm:"type:bytea"`
	CreatedAt time.Time

	Visits []Visit
}
