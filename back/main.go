package main

import (
	"log"

	"github.com/Marattttt/portfolio/portfolio_back/internal/api/handlers"
	"github.com/Marattttt/portfolio/portfolio_back/internal/appconfig"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	globalConf, _, err := appconfig.CreateAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	dsn := globalConf.DB.GetDSN()

	_, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	if err := handlers.SetupHandlers(r); err != nil {
		log.Fatal(err)
	}

	r.Run(globalConf.Server.ListenOn)
	log.Fatal("Server stopped working!")
}
