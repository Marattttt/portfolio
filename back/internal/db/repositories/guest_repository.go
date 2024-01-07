package repositories

import (
	"github.com/Marattttt/portfolio/portfolio_back/internal/db/models"
	"gorm.io/gorm"
)

type Guests struct {
	dbconn *gorm.DB
}

func NewGuest(dbConn *gorm.DB) Guests {
	return Guests{
		dbconn: dbConn,
	}
}

func (g Guests) Get(id int) (*gorm.DB, *models.Guest) {
	var guest models.Guest

	result := g.dbconn.First(&guest, id)
	return result, &guest
}

func (g Guests) Create(newguest *models.Guest) (*gorm.DB, *models.Guest) {
	res := g.dbconn.Create(&newguest)
	return res, newguest
}

func (g Guests) Update(id int) (*gorm.DB, *models.Guest) {
	var guest models.Guest

	result := g.dbconn.Save(&guest)
	return result, &guest
}

func (g Guests) Delete(id int) (*gorm.DB, *models.Guest) {
	var guest models.Guest

	result := g.dbconn.Delete(&guest, id)
	return result, &guest
}
