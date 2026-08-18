package main

import (
	"log"

	"go-api/internal/config"
	"go-api/internal/matrix"
	"go-api/internal/statistics"
)

func main() {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	client := statistics.NewClient(cfg)
	app := matrix.NewApp(cfg, client)
	log.Fatal(app.Listen(":" + cfg.Port))
}
