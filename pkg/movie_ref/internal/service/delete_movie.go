package service

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/3n0ugh/allotropes/framework/application"
	"github.com/3n0ugh/allotropes/internal/errors"
	"github.com/3n0ugh/allotropes/internal/middleware"
	"github.com/couchbase/gocb/v2"
	"github.com/go-chi/chi"
)

type DeleteMovie struct {
	Repo *gocb.Bucket
}

type DeleteMovieRequest struct {
	ID int `path:"id"`
}

type DeleteMovieResponse struct{}

func NewDeleteMovie(repo *gocb.Bucket) *DeleteMovie {
	return &DeleteMovie{Repo: repo}
}

func (m *DeleteMovie) Route(ctx context.Context) application.Route {
	return application.Route{
		Name:        "Delete Movie",
		Description: "Delete movie",
		Method:      http.MethodDelete,
		Path:        "/v1/movies/{id}",
		Headers:     map[string]string{},
		Middlewares: []func(http.Handler) http.Handler{middleware.Auth},
		Handler:     m.endpoint(ctx),
		Request:     DeleteMovieRequest{},
		Response:    DeleteMovieResponse{},
	}
}

func (m *DeleteMovie) endpoint(ctx context.Context) application.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (any, error) {
		idStr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, errors.NewBadRequestError("id must be integer", errors.Wrap(err, "id conversion").Error())
		}

		return m.handle(ctx, DeleteMovieRequest{ID: id})
	}
}

func (m *DeleteMovie) handle(ctx context.Context, r DeleteMovieRequest) (*DeleteMovieResponse, error) {
	err := m.repo(ctx, r.ID)
	if err != nil {
		return nil, errors.NewInternalServerError(errors.Wrap(err, "couchbase query").Error())
	}

	return &DeleteMovieResponse{}, nil
}

func (m *DeleteMovie) repo(_ context.Context, id int) error {
	_, err := m.Repo.Scope("movie").Collection("movie").Remove(strconv.Itoa(id), &gocb.RemoveOptions{Timeout: 5 * time.Second})
	if err != nil {
		return errors.Wrap(err, "couchbase query")
	}
	return nil
}
