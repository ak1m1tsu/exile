package main

import (
	"log"
	"os"

	"github.com/insan1a/exile/internal/app"
)

func main() {
	cfg, err := app.LoadConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	if err = app.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
