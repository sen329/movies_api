package query

import (
	"context"
	"database/sql"
	response "movies/Contract/Response"
	mapper "movies/Mappers"
	model "movies/Model"
	service "movies/Services"
	"strings"
	"time"

	"github.com/go-redis/cache/v8"
)

type MovieQuery struct {
	Db    *sql.DB
	Cache *cache.Cache
}

var emptyMovieResponse response.Movie

func (m MovieQuery) GetAllMovieQuery() ([]response.Movie, error) {
	ctx := context.TODO()
	key := "movies"
	var movies []response.Movie

	if err := m.Cache.Get(ctx, key, &movies); err == nil {
		return movies, nil
	}

	service := service.MovieService{
		Db: m.Db,
	}

	rows, err := service.GetMovies()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var movie model.Movie

		err := rows.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.Rating, &movie.Image, &movie.Created_at, &movie.Updated_at)
		if err != nil {
			return nil, err
		}

		responseMovie := mapper.MapMovies(movie)

		movies = append(movies, responseMovie)
	}

	if err := m.Cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: movies,
		TTL:   time.Second * 5,
	}); err != nil {
		return nil, err
	}

	return movies, nil
}

func (m MovieQuery) GetMovieById(param string) (response.Movie, error) {
	var movie model.Movie
	var response response.Movie
	ctx := context.TODO()
	key := "movie_" + param

	if err := m.Cache.Get(ctx, key, &response); err == nil {
		return response, nil
	}

	service := service.MovieService{
		Db: m.Db,
	}

	row := service.GetMovie(param)

	err := row.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.Rating, &movie.Image, &movie.Created_at, &movie.Updated_at)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return emptyMovieResponse, err
		}
		return emptyMovieResponse, err
	}

	response = mapper.MapMovies(movie)

	if err := m.Cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: response,
		TTL:   time.Second * 60,
	}); err != nil {
		return emptyMovieResponse, err
	}

	return response, nil
}
