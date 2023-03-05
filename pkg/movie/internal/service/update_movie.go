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
	"github.com/go-chi/chi"
)

type UpdateMovie struct {
	Repo *gocb.Bucket
}

type UpdateMovieRequest struct {
	ID    int          `path:"id"`
	Movie domain.Movie `json:"movie"`
}

type UpdateMovieResponse struct{}

func NewUpdateMovie(repo *gocb.Bucket) *UpdateMovie {
	return &UpdateMovie{Repo: repo}
}

func (m *UpdateMovie) Route(ctx context.Context) application.Route {
	return application.Route{
		Name:        "Update Movie",
		Description: "Update movie by id",
		Method:      http.MethodGet,
		Path:        "/v1/movies/{id}",
		Headers:     map[string]string{},
		Middlewares: []func(http.Handler) http.Handler{middleware.Auth},
		Handler:     m.endpoint(ctx),
		Request:     UpdateMovieRequest{},
		Response:    UpdateMovieResponse{},
	}
}

func (m *UpdateMovie) endpoint(ctx context.Context) application.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (any, error) {
		idStr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, errors.NewBadRequestError("id must be integer", errors.Wrap(err, "id conversion").Error())
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, errors.NewBadRequestError("unaccepted body", "read body")
		}

		var movie domain.Movie

		err = json.Unmarshal(body, &movie)
		if err != nil {
			return nil, errors.NewBadRequestError("unaccepted body", errors.Wrap(err, "movie body unmarshal").Error())
		}

		return m.handle(ctx, UpdateMovieRequest{ID: id, Movie: movie})
	}
}

func (m *UpdateMovie) handle(ctx context.Context, r UpdateMovieRequest) (*UpdateMovieResponse, error) {
	err := r.Movie.Validate()
	if err != nil {
		return nil, errors.NewBadRequestError(err.Error(), errors.Wrap(err, "validation").Error())
	}

	err = m.repo(ctx, r.ID, r.Movie)
	if err != nil {
		return nil, errors.NewInternalServerError(errors.Wrap(err, "couchbase query").Error())
	}

	return &UpdateMovieResponse{}, nil
}

func (m *UpdateMovie) repo(ctx context.Context, id int, movie domain.Movie) error {
	_, err := m.Repo.Scope("movie").Collection("movie").Upsert(strconv.Itoa(movie.ID), movie, &gocb.UpsertOptions{Timeout: 5 * time.Second})
	if err != nil {
		return errors.Wrap(err, "couchbase query")
	}
	return nil
}
