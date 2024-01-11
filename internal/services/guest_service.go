package services

import (
	"log"
	"time"

	"github.com/Marattttt/portfolio/portfolio_back/internal/db/models"
	"github.com/Marattttt/portfolio/portfolio_back/internal/db/repositories"
	"gorm.io/gorm"
)

type Guests struct {
	dbconn *gorm.DB
}

func NewGuestsService(dbConn *gorm.DB) Guests {
	return Guests{
		dbconn: dbConn,
	}
}

func (g Guests) GetGuest(id int) *models.Guest {
	repo := repositories.NewGuestsRepository(g.dbconn)
	res, guest := repo.Get(id)

	if res.Error != nil {
		log.Default()
		return nil
	}

	return guest
}

func (g Guests) NewGuest(guest models.Guest) (*models.Guest, error) {
	if len(guest.Visits) == 0 {
		guest.Visits = []models.Visit{
			{
				VisitedAt: time.Now(),
			},
		}
	}

	repo := repositories.NewGuestsRepository(g.dbconn)

	res, newGuest := repo.Create(guest)

	if res.Error != nil {
		return nil, res.Error
	}

	return newGuest, nil
}

func (g Guests) AddVisit(guest models.Guest) (*models.Guest, error) {
	if guest.ID == 0 {

	}

	repo := repositories.NewVisitsService(g.dbconn)
}
