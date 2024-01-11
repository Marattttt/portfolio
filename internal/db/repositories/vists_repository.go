package repositories

import (
	"github.com/Marattttt/portfolio/portfolio_back/internal/db/models"
	"gorm.io/gorm"
)

type Visits struct {
	dbconn *gorm.DB
}

func NewVisitsService(conn *gorm.DB) Visits {
	return Visits{
		dbconn: conn,
	}
}

func (v Visits) Get(id int) (*gorm.DB, *models.Visit) {
	var visit models.Visit
	result := v.dbconn.First(&visit, id)
	return result, &visit
}

func (v Visits) GetByGuestId(id int) (*gorm.DB, []models.Visit) {
	var visits []models.Visit
	res := v.dbconn.Model(&models.Visit{}).Where("guest_id = ?", id).Find(&visits)
	return res, visits
}

func (v Visits) Create(newVisit *models.Visit) *gorm.DB {
	res := v.dbconn.Create(&newVisit)
	return res
}

func (v Visits) Update(visit *models.Visit) *gorm.DB {
	result := v.dbconn.Save(visit)
	return result
}

func (v Visits) Delete(id int) (*gorm.DB, *models.Visit) {
	var visit models.Visit
	result := v.dbconn.Delete(&visit, id)
	return result, &visit
}
