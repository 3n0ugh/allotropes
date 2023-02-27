package movie

import (
	"github.com/3n0ugh/allotropes/framework/application"
	"github.com/3n0ugh/allotropes/internal/config"
	"github.com/3n0ugh/allotropes/pkg/movie/internal"
	"github.com/couchbase/gocb/v2"
)

func InitController(c config.Config, db *gocb.Bucket) application.Controller {
	repository := internal.NewRespository(db)
	service := internal.NewService(repository)
	routes := internal.SetRoute(service)

	return application.Controller{
		Name:        "Movie",
		Description: "Movie related services",
		Routes:      routes,
	}
}
