package service

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/3n0ugh/allotropes/framework/application"
	"github.com/3n0ugh/allotropes/internal/errors"
	"github.com/3n0ugh/allotropes/pkg/movie/internal/domain"
	"github.com/couchbase/gocb/v2"
	"github.com/go-chi/chi"
)

type GetMovieByID struct {
	Repo *gocb.Bucket
}

type GetMovieByIDRequest struct {
	ID int `path:"id"`
}

type GetMovieByIDResponse struct {
	Movie domain.Movie `json:"movie"`
}

func NewGetMovieByID(repo *gocb.Bucket) *GetMovieByID {
	return &GetMovieByID{Repo: repo}
}

func (m *GetMovieByID) Route(ctx context.Context) application.Route {
	return application.Route{
		Name:        "Get Movie",
		Description: "Get movie by id",
		Method:      http.MethodPut,
		Path:        "/v1/movies/{id}",
		Headers:     map[string]string{},
		Handler:     m.endpoint(ctx),
		Request:     GetMovieByIDRequest{},
		Response:    GetMovieByIDResponse{},
	}
}

func (m *GetMovieByID) endpoint(ctx context.Context) application.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (any, error) {
		idStr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, errors.NewBadRequestError("id must be integer", errors.Wrap(err, "id conversion").Error())
		}

		return m.handle(ctx, GetMovieByIDRequest{ID: id})
	}
}

func (m *GetMovieByID) handle(ctx context.Context, r GetMovieByIDRequest) (*GetMovieByIDResponse, error) {
	movie, err := m.repo(r.ID)
	if err != nil {
		return nil, errors.NewInternalServerError(errors.Wrap(err, "couchbase query").Error())
	}

	return &GetMovieByIDResponse{Movie: *movie}, nil
}

func (m *GetMovieByID) repo(ID int) (*domain.Movie, error) {
	doc, err := m.Repo.Scope("movie").Collection("movie").Get(strconv.Itoa(ID), &gocb.GetOptions{
		Timeout: time.Second * 3,
	})
	if err != nil {
		return nil, errors.Wrap(err, "couchbase query")
	}

	var movie domain.Movie
	if err := doc.Content(&movie); err != nil {
		return nil, errors.Wrap(err, "row parse")
	}

	return &movie, nil
}
