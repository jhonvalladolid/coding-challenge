package main

import (
	"log"

	"go-api/internal/config"
	"go-api/internal/matrix"
)

func main() {
	cfg := config.Load()
	app := matrix.NewApp(cfg)
	log.Fatal(app.Listen(":" + cfg.Port))
}
