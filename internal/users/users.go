package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/Marattttt/portfolio/portfolio_back/internal/config"
	"github.com/Marattttt/portfolio/portfolio_back/internal/config/dbconfig"
	"github.com/Marattttt/portfolio/portfolio_back/internal/models"
	"github.com/Marattttt/portfolio/portfolio_back/internal/storage"
	"github.com/Marattttt/portfolio/portfolio_back/internal/storage/pg"
)

// Handles both guests and their visits
// Transactions, if needed should be deined in the passed dbconn
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

func NewFromConfig(l applog.Logger, conf *config.AppConfig) (*Users, error) {
	var u Users

	if conf.Storage.DB != nil {
		err := u.addDb(l, *conf.Storage.DB)
		if err != nil {
			return nil, err
		}
	}

	return &u, nil
}

func (u *Users) addDb(l applog.Logger, conf dbconfig.DbConfig) error {
	pg, err := pg.NewUsersPGRepository(conf, l)
	if err != nil {
		return fmt.Errorf("creating pg users repository: %w", err)
	}
	u.db = pg
	return nil
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

func (u Users) Create(ctx context.Context, user models.User) (*models.User, error) {
	err := u.db.Create(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("creating user in pg: %w", err)
	}

	return &user, nil
}
