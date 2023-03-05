package service

import (
	"context"
	"net/http"
	"strconv"

	"github.com/3n0ugh/allotropes/framework/application"
	"github.com/3n0ugh/allotropes/internal/errors"
	"github.com/3n0ugh/allotropes/internal/pagination"
	"github.com/3n0ugh/allotropes/pkg/movie/internal/domain"
	"github.com/couchbase/gocb/v2"
)

type GetMovies struct {
	Repo *gocb.Bucket
}

type GetMoviesRequest struct {
	Page int `query:"page"`
	Size int `query:"size"`
}

type GetMoviesResponse struct {
	TotalCount int              `json:"totalCount"`
	Movies     []domain.Movie   `json:"movies"`
	Pagination pagination.Model `json:"pagination"`
}

func NewGetMovies(repo *gocb.Bucket) *GetMovies {
	return &GetMovies{Repo: repo}
}

func (m *GetMovies) Route(ctx context.Context) application.Route {
	return application.Route{
		Name:        "Get Movies",
		Description: "Get movies by page and page size",
		Method:      http.MethodGet,
		Path:        "/v1/movies",
		Headers:     map[string]string{},
		Handler:     m.endpoint(ctx),
		Request:     GetMoviesRequest{},
		Response:    GetMoviesResponse{},
	}
}

func (m *GetMovies) endpoint(ctx context.Context) application.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (any, error) {
		var req GetMoviesRequest

		if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
			req.Page = p
		}

		if p, err := strconv.Atoi(r.URL.Query().Get("pageSize")); err == nil {
			req.Size = p
		}

		return m.handle(ctx, req)
	}
}

func (m *GetMovies) handle(ctx context.Context, r GetMoviesRequest) (*GetMoviesResponse, error) {
	if r.Size > 20 {
		r.Size = 20
	}

	movies, err := m.repo(ctx, r.Page, r.Size)
	if err != nil {
		return nil, errors.NewInternalServerError(errors.Wrap(err, "couchbase query").Error())
	}
	return &GetMoviesResponse{Movies: movies}, nil
}

func (m *GetMovies) repo(ctx context.Context, page, pageSize int) ([]domain.Movie, error) {
	query := "SELECT movie FROM `movie`.movie.movie OFFSET $1 LIMIT $2"

	rows, err := m.Repo.Scope("movie").Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{page * pageSize, pageSize},
	})

	var movies []domain.Movie

	for rows.Next() {
		var movie struct {
			M domain.Movie `json:"movie"`
		}

		err := rows.Row(&movie)
		if err != nil {
			return nil, errors.Wrap(err, "row parse")
		}

		movies = append(movies, movie.M)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows")
	}

	return movies, nil
}
