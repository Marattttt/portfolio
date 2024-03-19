package models

import (
	"log/slog"
	"time"
)

type User struct {
	ID     int    `db:"id"`
	Name   string `db:"name"`
	Salt   []byte `db:"salt"`
	Secret []byte `db:"secret"`

	CreatedAt time.Time `db:"created_at"`

	// Nil for not deletd, any value for deleted
	DeletedAt *time.Time `db:"deleted_at"`
}

func (u User) LogValue() slog.Value {
	gvals := []slog.Attr{
		slog.Int("id", u.ID),
		slog.String("name", u.Name),
		slog.Time("createdAt", u.CreatedAt),
	}

	if u.DeletedAt != nil {
		gvals = append(gvals, slog.Time("deletedAt", *u.DeletedAt))
	}

	return slog.GroupValue(gvals...)
}

func (u *User) UpdateWith(g1 User) {
	u.ID = g1.ID
	u.Name = g1.Name
	u.Salt = g1.Salt
	u.Secret = g1.Secret
	u.CreatedAt = g1.CreatedAt
	u.DeletedAt = g1.DeletedAt
}
