package models

import "time"

type Visit struct {
	ID        uint `gorm:"column:visit_id"`
	VisitedAt time.Time
}
