package main

import (
	"context"
	"log"
	"time"

	"github.com/3n0ugh/allotropes/framework/application"
	"github.com/3n0ugh/allotropes/internal/config"
	"github.com/3n0ugh/allotropes/internal/database"
	"github.com/3n0ugh/allotropes/pkg/movie"
)

func main() {
	cfg := config.ReadConfig()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	cb, err := database.OpenConnectionCB(cfg)
	if err != nil {
		log.Fatal(err)

	}

	movieRefController := movie.InitController(ctx, cfg, cb)

	a := application.App{
		Name:           "Movpic",
		Port:           8080,
		Controllers:    []application.Controller{movieRefController},
		SwaggerEnabled: true,
	}
	a.Setup()
	a.Run()
}
