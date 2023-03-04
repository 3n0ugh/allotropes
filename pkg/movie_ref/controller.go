package movie

import (
	"github.com/3n0ugh/allotropes/framework/application"
	"github.com/3n0ugh/allotropes/internal/config"
	"github.com/3n0ugh/allotropes/pkg/movie_ref/internal/service"
	"github.com/couchbase/gocb/v2"
)

func InitController(c config.Config, db *gocb.Bucket) application.Controller {
	addMovieSvc := service.NewAddMovie(db)

	return application.Controller{
		Name:        "Movie",
		Description: "Movie related services",
		Routes: []application.Route{
			addMovieSvc.Route(),
		},
	}
}
