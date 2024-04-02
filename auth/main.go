package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	config, err := createConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	if err := configureApp(*config); err != nil {
		log.Fatal(err)
	}

	configFormatted, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(configFormatted))
}
