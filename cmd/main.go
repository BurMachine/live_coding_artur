package main

import (
	"awesomeProject/internal/app"
	"awesomeProject/internal/config"
	"context"
	"github.com/caarlos0/env/v6"
	"log"
)

func main() {
	cfg := &config.Config{}

	if err := env.Parse(cfg); err != nil {
		log.Fatalf("failed to retrieve env variables, %v", err)
	}

	err := app.Run(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
}
