package guests

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
type Guests struct {
	logger applog.Logger
	db     storage.GuestsRepository
}

func New(l applog.Logger, repo storage.GuestsRepository) Guests {
	return Guests{
		logger: l,
		db:     repo,
	}
}

func NewFromConfig(l applog.Logger, conf *config.AppConfig) (*Guests, error) {
	var g Guests

	if conf.Storage.DB != nil {
		err := g.addDb(l, *conf.Storage.DB)
		if err != nil {
			return nil, err
		}
	}

	return &g, nil
}

func (g *Guests) addDb(l applog.Logger, conf dbconfig.DbConfig) error {
	pg, err := pg.NewGuestsPGRepository(conf, l)
	if err != nil {
		return fmt.Errorf("creating pg guests repository: %w", err)
	}
	g.db = pg
	return nil
}

// Returns nil if any error is encountered
func (g Guests) GetGuest(ctx context.Context, id int) *models.Guest {
	guest, err := g.db.Get(ctx, id)

	if err != nil {
		if errors.Is(err, storage.EntityNotExistsError{ID: id}) {
			return nil
		}
		g.logger.Error(ctx, applog.DB, "failed to retrieve guest by id from db", err, slog.Int("entityId", id))
		return nil
	}

	return guest
}

func (g Guests) NewGuest(ctx context.Context, guest models.Guest) (*models.Guest, error) {
	err := g.db.Create(ctx, &guest)
	if err != nil {
		return nil, fmt.Errorf("creating in pg: %w", err)
	}

	return &guest, nil
}
