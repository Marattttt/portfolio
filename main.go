package main

import (
	"log"

	"github.com/Marattttt/portfolio/portfolio_back/internal/api/handlers"
	"github.com/Marattttt/portfolio/portfolio_back/internal/api/middleware"
	"github.com/Marattttt/portfolio/portfolio_back/internal/appconfig"
	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/gin-gonic/gin"
)

func main() {
	var globalConf appconfig.AppConfig
	if conf, _, err := appconfig.CreateAppConfig(); err != nil {
		log.Fatalf(applog.Config, err)
	} else {
		globalConf = *conf
	}

	applog.Setup(globalConf.Log)

	// Initialize the dbconnection pool
	if _, err := globalConf.DB.Connect(); err != nil {
		applog.Fatal(applog.Db, err)
	}

	r := gin.Default()

	middleware.AddMiddleware(r, &globalConf)

	if err := handlers.SetupHandlers(r); err != nil {
		applog.Fatal(applog.Http, err)
	}

	r.Run(globalConf.Server.ListenOn)
	applog.Fatal("Server stopped working")
}
