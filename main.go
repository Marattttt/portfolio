package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/Marattttt/portfolio/portfolio_back/internal/config"
)

func main() {
	conf, err := config.New()
	if err != nil {
		log.Fatalf("While creating appconfig: %v\n", err)
	}

	err = json.NewEncoder(os.Stderr).Encode(conf)
	if err != nil {
		log.Fatalf("While logging created config to stderr: %v", err)
	}

	logger, err := applog.New(conf.Log)
	if err != nil {
		log.Fatalf("While creating logger: %v", err)
	}

	marshalledConf, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		log.Fatalf("Marshalling created config: %v", err)
	}

	logger.Info(applog.Config, "Using config", "conf", conf)
	log.Println("Starting up using config: \n" + string(marshalledConf))

	// if err := applog.Setup(conf.Log); err != nil {
	// 	stdlog.Fatal(err)
	// }

	// if err := applog.SetupOtherLogs(conf.Log); err != nil {
	// 	stdlog.Fatal(err)
	// }

	// r := gin.Default()

	// middleware.AddMiddleware(r, &conf)

	// if err := handlers.SetupHandlers(r); err != nil {
	// 	applog.Fatal(applog.HTTP, err)
	// }

	// r.Run(conf.Server.ListenOn)
	// applog.Fatal("Server stopped working")
}
