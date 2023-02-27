package internal

import (
	"strconv"
	"time"

	"github.com/3n0ugh/allotropes/internal/errors"
	"github.com/couchbase/gocb/v2"
)

type Repository interface {
	Movies(page, pageSize int) ([]Movie, error)
	MovieByID(id int) (*Movie, error)
	AddMovie(movie Movie) error
	UpdateMovie(id int, movie Movie) error
	DeleteMovie(id int) error
}

type repository struct {
	db *gocb.Bucket
}

func NewRespository(db *gocb.Bucket) Repository {
	return &repository{db: db}
}

func (r *repository) Movies(page, pageSize int) ([]Movie, error) {
	query := "SELECT movie FROM `movie`.movie.movie OFFSET $1 LIMIT $2"

	rows, err := r.db.Scope("movie").Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{page * pageSize, pageSize},
	})

	var movies []Movie

	for rows.Next() {
		var movie struct {
			M Movie `json:"movie"`
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

func (r *repository) MovieByID(ID int) (*Movie, error) {
	doc, err := r.db.Scope("movie").Collection("movie").Get(strconv.Itoa(ID), &gocb.GetOptions{
		Timeout: time.Second * 3,
	})
	if err != nil {
		return nil, errors.Wrap(err, "couchbase query")
	}

	var movie Movie
	if err := doc.Content(&movie); err != nil {
		return nil, errors.Wrap(err, "row parse")
	}

	return &movie, nil
}

func (r *repository) AddMovie(movie Movie) error {
	_, err := r.db.Scope("movie").Collection("movie").Insert(strconv.Itoa(movie.ID), movie, &gocb.InsertOptions{Timeout: 5 * time.Second})
	if err != nil {
		return errors.Wrap(err, "couchbase query")
	}
	return nil
}

func (r *repository) UpdateMovie(id int, movie Movie) error {
	_, err := r.db.Scope("movie").Collection("movie").Upsert(strconv.Itoa(movie.ID), movie, &gocb.UpsertOptions{Timeout: 5 * time.Second})
	if err != nil {
		return errors.Wrap(err, "couchbase query")
	}
	return nil
}

func (r *repository) DeleteMovie(id int) error {
	_, err := r.db.Scope("movie").Collection("movie").Remove(strconv.Itoa(id), &gocb.RemoveOptions{Timeout: 5 * time.Second})
	if err != nil {
		return errors.Wrap(err, "couchbase query")
	}
	return nil
}
