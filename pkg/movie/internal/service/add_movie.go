package service

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/3n0ugh/allotropes/framework/application"
	"github.com/3n0ugh/allotropes/internal/errors"
	"github.com/3n0ugh/allotropes/internal/middleware"
	"github.com/3n0ugh/allotropes/pkg/movie/internal/domain"
	"github.com/couchbase/gocb/v2"
)

type AddMovie struct {
	Repo *gocb.Bucket
}

type AddMovieRequest struct {
	Movie domain.Movie `json:"movie"`
}

type AddMovieResponse struct{}

func NewAddMovie(repo *gocb.Bucket) *AddMovie {
	return &AddMovie{Repo: repo}
}

func (m *AddMovie) Route(ctx context.Context) application.Route {
	return application.Route{
		Name:        "Add Movie",
		Description: "Add movie",
		Method:      http.MethodPost,
		Path:        "/v1/movies",
		Headers:     map[string]string{},
		Middlewares: []func(http.Handler) http.Handler{middleware.Auth},
		Handler:     m.endpoint(ctx),
		Request:     AddMovieRequest{},
		Response:    AddMovieResponse{},
	}
}

func (m *AddMovie) endpoint(ctx context.Context) application.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (any, error) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, errors.NewBadRequestError("unaccepted body", "read body")
		}

		var movie domain.Movie

		err = json.Unmarshal(body, &movie)
		if err != nil {
			return nil, errors.NewBadRequestError("unaccepted body", errors.Wrap(err, "movie body unmarshal").Error())
		}

		return m.handle(ctx, AddMovieRequest{Movie: movie})
	}
}

func (m *AddMovie) handle(ctx context.Context, r AddMovieRequest) (*AddMovieResponse, error) {
	err := r.Movie.Validate()
	if err != nil {
		return nil, errors.NewBadRequestError(err.Error(), errors.Wrap(err, "validation").Error())
	}

	err = m.repo(ctx, r.Movie)
	if err != nil {
		return nil, errors.NewInternalServerError(errors.Wrap(err, "couchbase query").Error())
	}

	return &AddMovieResponse{}, nil
}

func (m *AddMovie) repo(_ context.Context, movie domain.Movie) error {
	_, err := m.Repo.Scope("movie").Collection("movie").Insert(strconv.Itoa(movie.ID), movie, &gocb.InsertOptions{Timeout: 5 * time.Second})
	if err != nil {
		return errors.Wrap(err, "couchbase query")
	}
	return nil
}
