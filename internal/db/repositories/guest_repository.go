package repositories

import (
	"github.com/Marattttt/portfolio/portfolio_back/internal/db/models"
	"gorm.io/gorm"
)

type Guests struct {
	dbconn *gorm.DB
}

func NewGuestsRepository(dbConn *gorm.DB) Guests {
	return Guests{
		dbconn: dbConn,
	}
}

func (g Guests) Get(id int) (*gorm.DB, *models.Guest) {
	var guest models.Guest

	res := g.dbconn.First(&guest, id)
	return res, &guest
}

func (g Guests) Create(newguest *models.Guest) *gorm.DB {
	if g.dbconn == nil {
		panic("ehee")
	}
	res := g.dbconn.Create(newguest)
	return res
}

func (g Guests) Update(guest *models.Guest) *gorm.DB {
	resul := g.dbconn.Save(&guest)
	return resul
}

func (g Guests) Delete(id int) (*gorm.DB, *models.Guest) {
	var guest models.Guest

	resul := g.dbconn.Delete(&guest, id)
	return resul, &guest
}
