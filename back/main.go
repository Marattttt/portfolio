package main

import (
	"log"

	"github.com/Marattttt/portfolio/portfolio_back/internal/api/handlers"
	"github.com/Marattttt/portfolio/portfolio_back/internal/db/dbconfig"
	"github.com/Marattttt/portfolio/portfolio_back/internal/db/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	vpr := viper.New()
	vpr.SetEnvPrefix("PORTFOLIO")
	dbConf, err := dbconfig.CreateConfig(*vpr)
	if err != nil {
		panic(err)
	}

	dsn := dbConf.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	guest := models.Guest{}

	db.Last(&guest)

	r := gin.Default()
	handlers.SetupHandlers(r)

	listenOn := ":" + vpr.GetString("PORT")
	r.Run(listenOn)
	log.Fatal("Server stopped working!")
}
