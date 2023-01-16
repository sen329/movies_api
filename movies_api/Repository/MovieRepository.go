package repository

import (
	"database/sql"
)

type MovieRepository struct {
	Db *sql.DB
}

func (m MovieRepository) GetMovies() (*sql.Rows, error) {
	query, err := m.Db.Query("SELECT id, title, description, rating, image, created_at, updated_at FROM movies;")
	if err != nil {
		return nil, err
	}
	return query, nil
}

func (m MovieRepository) GetMovie(param string) *sql.Row {
	query := m.Db.QueryRow("SELECT id, title, description, rating, image, created_at, updated_at FROM movies WHERE id = ?;", param)
	return query
}
