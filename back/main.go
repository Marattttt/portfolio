package main

import (
	"log"

	"github.com/Marattttt/portfolio/portfolio_back/internal/api/handlers"
	"github.com/Marattttt/portfolio/portfolio_back/internal/api/middleware"
	"github.com/Marattttt/portfolio/portfolio_back/internal/appconfig"
	"github.com/gin-gonic/gin"
)

func main() {
	var globalConf appconfig.AppConfig
	if conf, _, err := appconfig.CreateAppConfig(); err != nil {
		log.Fatal(err)
	} else {
		globalConf = *conf
	}

	// Initialize the dbconnection pool
	if _, err := globalConf.DB.Connect(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	middleware.AddMiddleware(r, &globalConf)

	if err := handlers.SetupHandlers(r); err != nil {
		log.Fatal(err)
	}

	r.Run(globalConf.Server.ListenOn)
	log.Fatal("Server stopped working!")
}
