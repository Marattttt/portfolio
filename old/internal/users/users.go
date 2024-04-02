package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/Marattttt/portfolio/portfolio_back/internal/models"
	"github.com/Marattttt/portfolio/portfolio_back/internal/storage"
)

type Users struct {
	logger applog.Logger
	db     storage.UsersRepository
}

func New(l applog.Logger, repo storage.UsersRepository) Users {
	return Users{
		logger: l,
		db:     repo,
	}
}

// Returns nil if any error is encountered
func (u Users) Get(ctx context.Context, id int) *models.User {
	user, err := u.db.Get(ctx, id)

	if err != nil {
		if errors.Is(err, storage.EntityNotExistsError{ID: id}) {
			return nil
		}
		u.logger.Error(ctx, applog.DB, "failed to retrieve user by id from db", err, slog.Int("entityId", id))
		return nil
	}

	return user
}

func (u Users) GetName(ctx context.Context, name string) {

}

func (u Users) Create(ctx context.Context, user models.User) (*models.User, error) {
	err := u.db.Create(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("creating user in pg: %w", err)
	}

	return &user, nil
}
