package movie

import (
	"context"

	"github.com/3n0ugh/allotropes/framework/application"
	"github.com/3n0ugh/allotropes/internal/config"
	"github.com/3n0ugh/allotropes/pkg/movie/internal/service"
	"github.com/couchbase/gocb/v2"
)

func InitController(ctx context.Context, c config.Config, db *gocb.Bucket) application.Controller {
	addMovieSvc := service.NewAddMovie(db)
	getMoviesSvc := service.NewGetMovies(db)
	getMovieByIDSvc := service.NewGetMovieByID(db)
	updateMovieSvc := service.NewUpdateMovie(db)
	deleteMovieSvc := service.NewDeleteMovie(db)

	return application.Controller{
		Name:        "Movie Ref",
		Description: "Movie Ref related services",
		Routes: []application.Route{
			addMovieSvc.Route(ctx),
			getMoviesSvc.Route(ctx),
			getMovieByIDSvc.Route(ctx),
			updateMovieSvc.Route(ctx),
			deleteMovieSvc.Route(ctx),
		},
	}
}
