package service

import (
	"database/sql"
	repository "movies/Repository"
)

type MovieService struct {
	Db *sql.DB
}

func (m MovieService) GetMovies() (*sql.Rows, error) {

	r := repository.MovieRepository{
		Db: m.Db,
	}

	queryResult, err := r.GetMovies()
	if err != nil {
		return nil, err
	}

	return queryResult, nil

}

func (m MovieService) GetMovie(param string) *sql.Row {
	r := repository.MovieRepository{
		Db: m.Db,
	}

	queryResult := r.GetMovie(param)

	return queryResult
}
