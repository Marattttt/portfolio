package main

import (
	"fmt"

	"github.com/Marattttt/portfolio/portfolio_back/internal/db"
	"github.com/spf13/viper"
)

func main() {
	vpr := viper.New()
	dbConfig, err := db.CreateConfig(*vpr)

	if err != nil {
		panic(err)
	}

	fmt.Println(dbConfig)
}
