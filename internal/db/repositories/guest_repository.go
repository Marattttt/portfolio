package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/Marattttt/portfolio/portfolio_back/internal/config"
	"github.com/Marattttt/portfolio/portfolio_back/internal/db/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/gorm"
)

type Guests struct {
	pool   *pgxpool.Pool
	logger *applog.AppLogger
}

// Contains id of the entity
type ErrEntityExists struct {
	Id int
}

// Error interface for ErrEntityExists
func (e ErrEntityExists) Error() string {
	return fmt.Sprintf("Entity id %d already exists", e.Id)
}

// Contains id of the entity
type ErrEntityNotExists struct {
	Id int
}

// Error interface for ErrEntityNotExists
func (e ErrEntityNotExists) Error() string {
	return fmt.Sprintf("Entity id %d does not exist", e.Id)
}

func NewGuestsRepository(conf *config.DbConfig, l *applog.AppLogger) (*Guests, error) {
	if conf.Pool == nil {
		return nil, fmt.Errorf("pgxpool.Pool is nil in dbconfig")
	}
	return &Guests{
		pool:   conf.Pool,
		logger: l,
	}, nil
}

// Nil if not found or an error was encountered
func (g Guests) Get(ctx context.Context, id int) *models.Guest {
	q := "SELECT * FROM guests WHERE guest_id = $1 LIMIT 1"

	rows, err := g.pool.Query(ctx, q, id)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var guest *models.Guest

	for rows.Next() {
		if guest != nil {
			g.logger.Error(applog.DB, "Expected 1 entity, received %d", rows.CommandTag().RowsAffected())
			return nil
		}

		if err = rows.Scan(&guest); err != nil {
			g.logger.Error(applog.DB, err)
		}
	}

	if guest.DeletedAt != nil {
		return nil
	}

	return guest
}

func (g Guests) Create(ctx context.Context, newG *models.Guest) error {
	if newG.ID != 0 {
		oldG := g.Get(ctx, newG.ID)
		if oldG != nil {
			return ErrEntityExists{Id: oldG.ID}
		}

	}

	q := "INSERT INTO guests(name, salt, secret, created_at, deleted_at) VALUES(@name, @salt, @secret, @created_at, @deleted_at) RETURNIG id"
	args := pgx.NamedArgs{
		"name":       newG.Name,
		"salt":       newG.Salt,
		"secret":     newG.Secret,
		"created_at": newG.CreatedAt,
		"deleted_at": newG.DeletedAt,
	}

	rows := g.pool.QueryRow(ctx, q, args)

	var newId int

	if err := rows.Scan(&newId); err != nil {
		return err
	}

	newG.ID = newId
	return nil
}

func (g Guests) Update(ctx context.Context, guest *models.Guest) error {
	old := g.Get(ctx, guest.ID)
	if old == nil {
		return ErrEntityNotExists{Id: guest.ID}
	}

	return nil
}

func (g Guests) Delete(id int) (*gorm.DB, *models.Guest) {
	var guest models.Guest

	resul := g.dbconn.Delete(&guest, id)
	return resul, &guest
}

type GuestsSTD struct {
	dbconn *gorm.DB
}

func NewGuestsSTDRepository(dbConn *gorm.DB) GuestsSTD {
	// _ := sql.
	return GuestsSTD{
		dbconn: dbConn,
	}
}

func (g GuestsSTD) Get(id int) (*gorm.DB, *models.Guest) {
	var guest models.Guest

	res := g.dbconn.First(&guest, id)
	return res, &guest
}

func (g GuestsSTD) Create(newguest *models.Guest) *gorm.DB {
	res := g.dbconn.Create(newguest)
	return res
}

func (g GuestsSTD) Update(guest *models.Guest) *gorm.DB {
	resul := g.dbconn.Save(&guest)
	return resul
}

func (g GuestsSTD) Delete(id int) (*gorm.DB, *models.Guest) {
	var guest models.Guest

	resul := g.dbconn.Delete(&guest, id)
	return resul, &guest
}
