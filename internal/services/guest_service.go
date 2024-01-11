package services

import (
	"fmt"
	"log"
	"time"

	"github.com/Marattttt/portfolio/portfolio_back/internal/db/models"
	"github.com/Marattttt/portfolio/portfolio_back/internal/db/repositories"
	"gorm.io/gorm"
)

// Handles both guests and their visits
// Transactions, if needed should be deined in the passed dbconn
type Guests struct {
	dbconn *gorm.DB
}

// Guest does not exist in the database
type GuestDoesNotExist struct {
	Guest *models.Guest
}

func (err GuestDoesNotExist) Error() string {
	if err.Guest != nil {
		return fmt.Sprintf("Guest id: %d does not exist", err.Guest.ID)
	}
	return "Requested guest does not exist"
}

func NewGuestsService(dbConn *gorm.DB) Guests {
	return Guests{
		dbconn: dbConn,
	}
}

// Returns nil if any error is encountered
func (g Guests) GetGuest(id int) *models.Guest {
	repo := repositories.NewGuestsRepository(g.dbconn)
	res, guest := repo.Get(id)

	if res.Error != nil {
		log.Default()
		return nil
	}

	return guest
}

// Creates a new guest with one visit at time.Now()
func (g Guests) NewGuest(guest models.Guest) (*models.Guest, error) {
	if len(guest.Visits) == 0 {
		guest.Visits = []models.Visit{
			{
				VisitedAt: time.Now(),
			},
		}
	}

	repo := repositories.NewGuestsRepository(g.dbconn)

	res := repo.Create(&guest)

	if res.Error != nil {
		return nil, res.Error
	}

	return &guest, nil
}

// Adds a visit at time.Now() to guest with the specified id
func (g Guests) AddVisit(guestId int) (*models.Guest, error) {
	guestsRepo := repositories.NewGuestsRepository(g.dbconn)

	guest := g.GetGuest(guestId)
	if guest == nil {
		return nil, GuestDoesNotExist{Guest: guest}
	}

	guest.Visits = append(guest.Visits, models.Visit{
		VisitedAt: time.Now(),
	})

	res := guestsRepo.Update(guest)

	if res.Error != nil {
		return nil, res.Error
	}

	return guest, nil
}
