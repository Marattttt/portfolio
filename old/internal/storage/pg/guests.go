package pg

import (
	"context"
	"fmt"

	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/Marattttt/portfolio/portfolio_back/internal/models"
	"github.com/Marattttt/portfolio/portfolio_back/internal/storage"
	"github.com/Marattttt/portfolio/portfolio_back/internal/storage/storageutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository struct {
	pool   *pgxpool.Pool
	logger applog.Logger
}

func New(pool *pgxpool.Pool, logger applog.Logger) UsersRepository {
	return UsersRepository{
		pool:   pool,
		logger: logger,
	}
}

func (u UsersRepository) Get(ctx context.Context, id int) (*models.User, error) {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("beginning a transaction: %w", err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	user, err := u.GetTx(ctx, &tx, id)
	if err != nil {
		return nil, err
	}

	_ = tx.Commit(ctx)

	return user, nil
}

// Returns nil if the entity was not found or has a non-null deletion time
func (u UsersRepository) GetTx(ctx context.Context, tx *pgx.Tx, id int) (*models.User, error) {
	var err error

	const query = `
		SELECT 
			* FROM users
		WHERE
			user_id = $1
		LIMIT 1
	`

	rows, err := (*tx).Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("executing query %s: %w", query, err)
	}
	defer rows.Close()

	users := storageutils.PGScanMany[models.User](ctx, u.logger, &rows)

	if len(users) == 0 {
		return nil, storage.EntityNotExistsError{ID: id}
	}
	if len(users) > 1 {
		return nil, fmt.Errorf("executing get by key %d returned %d entities, when expected to return 1", id, len(users))
	}

	return &users[0], nil
}

func (u UsersRepository) GetName(ctx context.Context, paramgs storage.QueryParams) ([]models.User, error) {

	return nil, nil
}

func (u UsersRepository) GetNameTx(ctx context.Context, tx *pgx.Tx, params storage.QueryParams) ([]models.User, error) {

}

func generateQuery(params storage.QueryParams) (string, pgx.NamedArgs) {
	base := "FROM users SELECT * WHERE "
	args := pgx.NamedArgs{}
	isFirstParam := true
	if params.Name != nil {
		if !isFirstParam {
			base += " AND "
		}
		base += "name = @name"
		args["name"] = *params.Name
	}

	return base

}

// The newG parameter will have its id updated if the operation is successful
//
// If the id of the entity is non-zero, a check for existing entity is performed,
// and if one exists ErrEntityExists is returned
func (u UsersRepository) Create(ctx context.Context, newU *models.User) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning a transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	err = u.CreateTx(ctx, &tx, newU)
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
func (u UsersRepository) CreateTx(ctx context.Context, tx *pgx.Tx, newU *models.User) error {
	if newU.ID != 0 {
		oldU, err := u.GetTx(ctx, tx, newU.ID)
		if err != nil {
			return err
		}

		if oldU != nil {
			return storage.EntityExistsError{ID: oldU.ID}
		}
	}

	const query = `
		INSERT INTO 
			users
			(name, salt, password, created_at, deleted_at)
		VALUES 
			(@name, @salt, @password, @created_at, @deleted_at)
		RETURNIG id
	`
	args := pgx.NamedArgs{
		"name":       newU.Name,
		"salt":       newU.Salt,
		"password":   newU.Password,
		"created_at": newU.CreatedAt,
		"deleted_at": newU.DeletedAt,
	}

	rows := u.pool.QueryRow(ctx, query, args)

	var newID int

	if err := rows.Scan(&newID); err != nil {
		return fmt.Errorf("inserting %v into users: %w", newU, err)
	}

	newU.ID = newID
	return nil
}

func (u UsersRepository) Update(ctx context.Context, user *models.User) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	err = u.Updatetx(ctx, &tx, user)
	if err != nil {
		return err
	}

	_ = tx.Commit(ctx)
	return nil
}

func (u UsersRepository) Updatetx(ctx context.Context, tx *pgx.Tx, user *models.User) error {
	old, err := u.GetTx(ctx, tx, user.ID)
	if err != nil {
		return err
	}

	if old == nil {
		return storage.EntityNotExistsError{ID: user.ID}
	}

	const query = `UPDATE users
		SET 
			name = @name,
			salt = @salt,
			password = @password,
			created_at = @created_at,
			deleted_at = @deleted_at
		WHERE 
			id = @id
		RETURNING *`

	args := pgx.NamedArgs{
		"id":         user.ID,
		"name":       user.Name,
		"salt":       user.Salt,
		"password":   user.Password,
		"created_at": user.CreatedAt,
		"deleted_at": user.DeletedAt,
	}

	rows, err := u.pool.Query(ctx, query, args)
	if err != nil {
		return fmt.Errorf("Executing update query: %w", err)
	}
	defer rows.Close()

	if rows.CommandTag().RowsAffected() == 0 {
		return storage.ErrNoEffect
	}

	results := storageutils.PGScanMany[models.User](ctx, u.logger, &rows)

	if len(results) == 0 {
		return storage.ErrNoEffect
	}

	if len(results) > 1 {
		return fmt.Errorf(
			"executing update by key %d caused %d updates; affected entitites; %v",
			user.ID, len(results), results)
	}

	user.UpdateWith(results[0])

	return nil
}

func (u UsersRepository) Delete(ctx context.Context, id int) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning a transaction: %w", err)
	}

	err = u.DeleteTx(ctx, &tx, id)
	if err != nil {
		return err
	}

	return nil
}

func (g UsersRepository) DeleteTx(ctx context.Context, tx *pgx.Tx, id int) error {
	const query = `
		UPDATE 
			guests 
		SET 
			deleted_at = CURRENT_TIMESTAMP 
		WHERE 
			id = @id RETURNING id
	`

	args := pgx.NamedArgs{"id": id}

	rows, err := (*tx).Query(ctx, query, args)
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
