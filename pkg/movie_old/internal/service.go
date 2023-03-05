package internal

import "github.com/3n0ugh/allotropes/internal/errors"

type Service interface {
	Movies(r MoviesRequest) (*MoviesResponse, error)
	MovieByID(r MovieByIDRequest) (*MovieByIDResponse, error)
	UpdateMovie(r UpdateMovieRequest) (*UpdateMovieResponse, error)
	DeleteMovie(r DeleteMovieRequest) (*DeleteMovieResponse, error)
	AddMovie(r AddMovieRequest) (*AddMovieResponse, error)
}

type service struct {
	Repo Repository
}

func NewService(r Repository) Service {
	return &service{Repo: r}
}

func (s *service) Movies(r MoviesRequest) (*MoviesResponse, error) {
	if r.Size > 20 {
		r.Size = 20
	}

	movies, err := s.Repo.Movies(r.Page, r.Size)
	if err != nil {
		return nil, errors.NewInternalServerError(errors.Wrap(err, "couchbase query").Error())
	}
	return &MoviesResponse{Movies: movies}, nil
}

func (s *service) MovieByID(r MovieByIDRequest) (*MovieByIDResponse, error) {
	movie, err := s.Repo.MovieByID(r.ID)
	if err != nil {
		return nil, errors.NewInternalServerError(errors.Wrap(err, "couchbase query").Error())
	}

	return &MovieByIDResponse{Movie: *movie}, nil
}

func (s *service) UpdateMovie(r UpdateMovieRequest) (*UpdateMovieResponse, error) {
	err := r.Movie.Validate()
	if err != nil {
		return nil, errors.NewBadRequestError(err.Error(), errors.Wrap(err, "validation").Error())
	}

	err = s.Repo.UpdateMovie(r.ID, r.Movie)
	if err != nil {
		return nil, errors.NewInternalServerError(errors.Wrap(err, "couchbase query").Error())
	}

	return &UpdateMovieResponse{}, nil
}

func (s *service) DeleteMovie(r DeleteMovieRequest) (*DeleteMovieResponse, error) {
	err := s.Repo.DeleteMovie(r.ID)
	if err != nil {
		return nil, errors.NewInternalServerError(errors.Wrap(err, "couchbase query").Error())
	}

	return &DeleteMovieResponse{}, nil
}

func (s *service) AddMovie(r AddMovieRequest) (*AddMovieResponse, error) {
	err := r.Movie.Validate()
	if err != nil {
		return nil, errors.NewBadRequestError(err.Error(), errors.Wrap(err, "validation").Error())
	}

	err = s.Repo.AddMovie(r.Movie)
	if err != nil {
		return nil, errors.NewInternalServerError(errors.Wrap(err, "couchbase query").Error())
	}

	return &AddMovieResponse{}, nil
}
