package internal

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/3n0ugh/allotropes/framework/application"
	"github.com/3n0ugh/allotropes/internal/errors"
	"github.com/3n0ugh/allotropes/internal/middleware"
	"github.com/go-chi/chi"
)

func SetRoute(svc Service) []application.Route {
	return []application.Route{
		{
			Name:        "Movies",
			Description: "Get movies by page and size",
			Method:      http.MethodGet,
			Path:        "/v1/movies",
			Handler:     application.HandlerFunc(MoviesEndpoint(svc)),
			Request:     MoviesRequest{},
			Response:    MoviesResponse{},
		},
		{
			Name:        "MoviesByID",
			Description: "Get movie by id",
			Method:      http.MethodGet,
			Path:        "/v1/movies/{id}",
			Handler:     application.HandlerFunc(MovieByIDEndpoint(svc)),
			Request:     MovieByIDRequest{},
			Response:    MovieByIDResponse{},
		},
		{
			Name:        "Update Movie",
			Description: "Update movie by id",
			Method:      http.MethodPut,
			Path:        "/v1/movie/{id}",
			Middlewares: []func(http.Handler) http.Handler{middleware.Auth},
			Handler:     application.HandlerFunc(UpdateMovieEndpoint(svc)),
			Request:     UpdateMovieRequest{},
			Response:    UpdateMovieResponse{},
		},
		{
			Name:        "Delete Movie",
			Description: "Delete movie by id",
			Method:      http.MethodDelete,
			Path:        "/v1/movie/{id}",
			Middlewares: []func(http.Handler) http.Handler{middleware.Auth},
			Handler:     application.HandlerFunc(DeleteMovieEndpoint(svc)),
			Request:     DeleteMovieRequest{},
			Response:    DeleteMovieResponse{},
		},
		{
			Name:        "Add Movie",
			Description: "Add movie",
			Method:      http.MethodPost,
			Path:        "/v1/movies",
			Middlewares: []func(http.Handler) http.Handler{middleware.Auth},
			Handler:     application.HandlerFunc(AddMovieEndpoint(svc)),
			Request:     AddMovieRequest{},
			Response:    AddMovieResponse{},
		},
	}
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) (response any, err error)

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	res, err := h(w, r)
	if err != nil {
		log.Printf("err: %s", err.Error())
		res = err
		if e, ok := err.(*errors.Error); ok {
			w.WriteHeader(e.StatusCode)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("err: %s", err.Error())
	}
}

func MoviesEndpoint(svc Service) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (any, error) {
		var req MoviesRequest

		if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
			req.Page = p
		}

		if p, err := strconv.Atoi(r.URL.Query().Get("pageSize")); err == nil {
			req.Size = p
		}

		return svc.Movies(req)
	}
}

func MovieByIDEndpoint(svc Service) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (any, error) {
		idStr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, errors.NewBadRequestError("id must be integer", errors.Wrap(err, "id conversion").Error())
		}

		return svc.MovieByID(MovieByIDRequest{ID: id})
	}
}

func UpdateMovieEndpoint(svc Service) HandlerFunc {
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

		var movie Movie

		err = json.Unmarshal(body, &movie)
		if err != nil {
			return nil, errors.NewBadRequestError("unaccepted body", errors.Wrap(err, "movie body unmarshal").Error())
		}

		err = movie.Validate()
		if err != nil {
			return nil, errors.NewBadRequestError(err.Error(), errors.Wrap(err, "validation").Error())
		}

		return svc.UpdateMovie(UpdateMovieRequest{
			ID:    id,
			Movie: movie,
		})
	}
}

func DeleteMovieEndpoint(svc Service) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (any, error) {
		idStr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, errors.NewBadRequestError("id must be integer", errors.Wrap(err, "id conversion").Error())
		}

		return svc.DeleteMovie(DeleteMovieRequest{ID: id})
	}
}

func AddMovieEndpoint(svc Service) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (any, error) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, errors.NewBadRequestError("unaccepted body", "read body")
		}

		var movie Movie

		err = json.Unmarshal(body, &movie)
		if err != nil {
			return nil, errors.NewBadRequestError("unaccepted body", errors.Wrap(err, "movie body unmarshal").Error())
		}

		err = movie.Validate()
		if err != nil {
			return nil, errors.NewBadRequestError(err.Error(), errors.Wrap(err, "validation").Error())
		}

		return svc.AddMovie(AddMovieRequest{Movie: movie})
	}
}
