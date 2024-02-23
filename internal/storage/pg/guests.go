package pg

import (
	"context"
	"fmt"

	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/Marattttt/portfolio/portfolio_back/internal/config/dbconfig"
	"github.com/Marattttt/portfolio/portfolio_back/internal/models"
	"github.com/Marattttt/portfolio/portfolio_back/internal/storage"
	"github.com/Marattttt/portfolio/portfolio_back/internal/storage/storageutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GuestsPGRepository struct {
	pool   *pgxpool.Pool
	logger applog.Logger
}

func NewGuestsPGRepository(conf dbconfig.DbConfig, l applog.Logger) (*GuestsPGRepository, error) {
	if conf.Pool == nil {
		return nil, fmt.Errorf("pgxpool.Pool is nil in dbconfig")
	}

	gpg := GuestsPGRepository{
		pool:   conf.Pool,
		logger: l,
	}

	return &gpg, nil
}

func (g GuestsPGRepository) Get(ctx context.Context, id int) (*models.Guest, error) {
	tx, err := g.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("beginning a transaction: %w", err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	guest, err := g.GetTx(ctx, &tx, id)
	if err != nil {
		return nil, err
	}

	_ = tx.Commit(ctx)

	return guest, nil
}

// Returns nil if the entity was not found or has a non-null deletion time
func (g GuestsPGRepository) GetTx(ctx context.Context, tx *pgx.Tx, id int) (*models.Guest, error) {
	var err error

	const q = `SELECT * FROM guests
		WHERE
			guest_id = $1
		LIMIT 1`

	rows, err := (*tx).Query(ctx, q, id)
	if err != nil {
		return nil, fmt.Errorf("executing query %s: %w", q, err)
	}
	defer rows.Close()

	guests := storageutils.PGScanMany[models.Guest](ctx, g.logger, &rows)

	if len(guests) == 0 {
		return nil, storage.EntityNotExistsError{ID: id}
	}
	if len(guests) > 1 {
		return nil, fmt.Errorf("executing get by key %d returned %d entities, when expected to return 1", id, len(guests))
	}

	return &guests[0], nil
}

// The newG parameter will have its id updated if the operation is successful
//
// If the id of the entity is non-zero, a check for existing entity is performed,
// and if one exists ErrEntityExists is returned
func (g GuestsPGRepository) Create(ctx context.Context, newG *models.Guest) error {
	tx, err := g.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	err = g.CreateTx(ctx, &tx, newG)
	if err == nil {
		_ = tx.Commit(ctx)
		return nil
	}
	return err
}

// The tx provided should be committed or rolled back by the caller;
// The newG parameter will have its id updated if the operation is successful
//
// If the id of the entity is non-zero, a check for existing entity is performed,
// and if one exists ErrEntityExists is returned
func (g GuestsPGRepository) CreateTx(ctx context.Context, tx *pgx.Tx, newG *models.Guest) error {
	if newG.ID != 0 {
		oldG, err := g.GetTx(ctx, tx, newG.ID)
		if err != nil {
			return err
		}

		if oldG != nil {
			return storage.EntityExistsError{ID: oldG.ID}
		}
	}

	const q = `INSERT INTO guests
		(name, salt, secret, created_at, deleted_at)
		VALUES (@name, @salt, @secret, @created_at, @deleted_at)
		RETURNIG id`
	args := pgx.NamedArgs{
		"name":       newG.Name,
		"salt":       newG.Salt,
		"secret":     newG.Secret,
		"created_at": newG.CreatedAt,
		"deleted_at": newG.DeletedAt,
	}

	rows := g.pool.QueryRow(ctx, q, args)

	var newID int

	if err := rows.Scan(&newID); err != nil {
		return fmt.Errorf("inserting %v into guests: %w", newG, err)
	}

	newG.ID = newID
	return nil
}

func (g GuestsPGRepository) Update(ctx context.Context, guest *models.Guest) error {
	tx, err := g.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	err = g.Updatetx(ctx, &tx, guest)
	if err != nil {
		return err
	}

	_ = tx.Commit(ctx)
	return nil
}

func (g GuestsPGRepository) Updatetx(ctx context.Context, tx *pgx.Tx, guest *models.Guest) error {
	old, err := g.GetTx(ctx, tx, guest.ID)
	if err != nil {
		return err
	}

	if old == nil {
		return storage.EntityNotExistsError{ID: guest.ID}
	}

	const q = `UPDATE guests 
		SET 
			name = @name,
			salt = @salt,
			secret = @secret,
			created_at = @created_at,
			deleted_at = @deleted_at
		WHERE 
			id = @id
		RETURNING *`

	args := pgx.NamedArgs{
		"id":         guest.ID,
		"name":       guest.Name,
		"salt":       guest.Salt,
		"secret":     guest.Secret,
		"created_at": guest.CreatedAt,
		"deleted_at": guest.DeletedAt,
	}

	rows, err := g.pool.Query(ctx, q, args)
	if err != nil {
		return fmt.Errorf("executing update query: %w", err)
	}
	defer rows.Close()

	if rows.CommandTag().RowsAffected() == 0 {
		return storage.ErrNoEffect
	}

	results := storageutils.PGScanMany[models.Guest](ctx, g.logger, &rows)

	if len(results) == 0 {
		return storage.ErrNoEffect
	}

	if len(results) > 1 {
		return fmt.Errorf(
			"executing update by key %d caused %d updates; affected entitites; %v",
			guest.ID, len(results), results)
	}

	guest.UpdateWith(results[0])

	return nil
}

func (g GuestsPGRepository) Delete(ctx context.Context, id int) error {
	tx, err := g.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning a transaction: %w", err)
	}

	err = g.DeleteTx(ctx, &tx, id)
	if err != nil {
		return err
	}

	return nil
}

func (g GuestsPGRepository) DeleteTx(ctx context.Context, tx *pgx.Tx, id int) error {
	const q = "UPDATE guests SET deleted_at = CURRENT_TIMESTAMP WHERE id = @id RETURNING id"
	args := pgx.NamedArgs{"id": id}

	rows, err := (*tx).Query(ctx, q, args)
	if err != nil {
		return fmt.Errorf("executing delete by id %d: %w", id, err)
	}
	defer rows.Close()

	deletedIds := storageutils.PGScanMany[int](ctx, g.logger, &rows)

	if len(deletedIds) == 0 {
		return storage.ErrNoEffect
	}
	if len(deletedIds) > 1 {
		_ = (*tx).Rollback(ctx)
		return fmt.Errorf("Delete by id %d caused %d deletions", id, len(deletedIds))
	}

	return nil
}
