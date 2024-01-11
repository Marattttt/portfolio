package models

import "time"

type Visit struct {
	ID        uint `gorm:"column:visit_id"`
	GuestID   uint
	VisitedAt time.Time
}
