package models

import (
	"log/slog"
	"time"
)

type Guest struct {
	ID     int    `db:"id"`
	Name   string `db:"name"`
	Salt   []byte `db:"salt"`
	Secret []byte `db:"secret"`

	CreatedAt time.Time `db:"created_at"`

	// Nil for not deletd, any value for deleted
	DeletedAt *time.Time `db:"deleted_at"`

	// Related visits
	Visits []Visit `db:"-"`
}

func (g Guest) LogValue() slog.Value {
	gvals := []slog.Attr{
		slog.Int("id", g.ID),
		slog.String("name", g.Name),
		slog.Time("createdAt", g.CreatedAt),
	}

	if g.DeletedAt != nil {
		gvals = append(gvals, slog.Time("deletedAt", *g.DeletedAt))
	}

	gvals = append(gvals, slog.Int("numberOfVisits", len(g.Visits)))

	return slog.GroupValue(gvals...)
}
