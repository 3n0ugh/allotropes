package main

import (
	"log"

	"github.com/3n0ugh/allotropes/framework/application"
	"github.com/3n0ugh/allotropes/internal/config"
	"github.com/3n0ugh/allotropes/internal/database"
	"github.com/3n0ugh/allotropes/pkg/movie"
)

func main() {
	cfg := config.ReadConfig()

	cb, err := database.OpenConnectionCB(cfg)
	if err != nil {
		log.Fatal(err)

	}

	controller := movie.InitController(cfg, cb)

	a := application.App{
		Name:           "Movpic",
		Port:           8080,
		Controllers:    []application.Controller{controller},
		SwaggerEnabled: true,
	}
	a.Setup()
	a.Run()
}
