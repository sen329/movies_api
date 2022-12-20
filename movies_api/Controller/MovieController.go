package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	request "movies/Contract/Request"
	response "movies/Contract/Response"
	model "movies/Model"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/cache/v8"
)

type MovieController struct {
	Db    *sql.DB
	Cache *cache.Cache
}

func (m MovieController) MapEndpoints(route *gin.Engine) {
	route.GET("/Movies", m.getMovies)
	route.GET("/Movies/:ID", m.getMovie)
	route.POST("/Movies", m.addMovie)
	route.PATCH("/Movies/:ID", m.updateMovie)
	route.DELETE("/Movies/:ID", m.deleteMovie)
}

func (m MovieController) getMovies(c *gin.Context) {
	ctx := context.TODO()
	key := "movies"
	var movies []response.Movie

	if err := m.Cache.Get(ctx, key, &movies); err == nil {
		c.JSON(200, movies)
		return
	}

	query, err := m.Db.Query("SELECT id, title, description, rating, image, created_at, updated_at FROM movies;")
	if err != nil {
		message := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something wrong while fetching data: " + err.Error(),
		}
		c.JSON(500, message)
		return
	}

	for query.Next() {
		var movie model.Movie

		err := query.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.Rating, &movie.Image, &movie.Created_at, &movie.Updated_at)
		if err != nil {
			message := response.ErrorResponse{
				ErrorCode: 500,
				Message:   "Something wrong while scan data: " + err.Error(),
			}
			c.JSON(500, message)
			return
		}

		responseMovie := response.Movie{
			Id:          movie.Id,
			Title:       movie.Title,
			Description: movie.Description,
			Rating:      movie.Rating,
			Image:       movie.Image,
			Created_at:  movie.Created_at.Format("2006-01-02 15:04:05"),
			Updated_at:  movie.Updated_at.Format("2006-01-02 15:04:05"),
		}

		movies = append(movies, responseMovie)
	}

	if movies == nil {
		message := response.ErrorResponse{
			ErrorCode: 404,
			Message:   "No movies registered",
		}
		c.JSON(404, message)
		return
	}

	if err := m.Cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: movies,
		TTL:   time.Second * 5,
	}); err != nil {
		message := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something wrong while caching data: " + err.Error(),
		}
		c.JSON(500, message)
		return
	}

	c.JSON(200, movies)

}

func (m MovieController) getMovie(c *gin.Context) {
	param := c.Param("ID")
	ctx := context.TODO()
	key := "movie_" + param
	var movie model.Movie

	if err := m.Cache.Get(ctx, key, &movie); err == nil {
		c.JSON(200, movie)
		return
	}

	query := m.Db.QueryRow("SELECT id, title, description, rating, image, created_at, updated_at FROM movies WHERE id = ?;", param)
	if query == nil {
		message := response.ErrorResponse{
			ErrorCode: 404,
			Message:   "No movies registered",
		}
		c.JSON(404, message)
		return
	}

	err := query.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.Rating, &movie.Image, &movie.Created_at, &movie.Updated_at)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			message := response.ErrorResponse{
				ErrorCode: 404,
				Message:   "No movies registered",
			}
			c.JSON(404, message)
			return
		}
		message := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something wrong while scan data: " + err.Error(),
		}
		c.JSON(500, message)
		return
	}

	responseMovie := response.Movie{
		Id:          movie.Id,
		Title:       movie.Title,
		Description: movie.Description,
		Rating:      movie.Rating,
		Image:       movie.Image,
		Created_at:  movie.Created_at.Format("2006-01-02 15:04:05"),
		Updated_at:  movie.Updated_at.Format("2006-01-02 15:04:05"),
	}

	if err := m.Cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: responseMovie,
		TTL:   time.Second * 60,
	}); err != nil {
		message := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something wrong while caching data: " + err.Error(),
		}
		c.JSON(500, message)
		return
	}

	c.JSON(200, responseMovie)

}

func (m MovieController) addMovie(c *gin.Context) {
	var requestMovieBody request.Movie

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		message := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Request body error: " + err.Error(),
		}
		c.JSON(400, message)
		return
	}

	json.Unmarshal(body, &requestMovieBody)

	validate := validator.New()
	err = validate.Struct(requestMovieBody)
	if err != nil {
		message := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Request body error: " + err.Error(),
		}
		c.JSON(400, message)
		return
	}

	stmt, err := m.Db.Prepare("INSERT INTO movies(id, title, description, rating, image, created_at, updated_at) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		message := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something wrong with preparing stmt: " + err.Error(),
		}
		c.JSON(500, message)
		return
	}

	_, err = stmt.Exec(requestMovieBody.Id, requestMovieBody.Title, requestMovieBody.Description, requestMovieBody.Rating,
		requestMovieBody.Image, requestMovieBody.Created_at, requestMovieBody.Updated_at)
	if err != nil {
		message := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something wrong with entering data to database: " + err.Error(),
		}
		c.JSON(500, message)
		return
	}

	c.JSON(200, "OK")

}

func (m MovieController) updateMovie(c *gin.Context) {
	param := c.Param("ID")
	var requestMovieBody request.Movie

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		message := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Request body error: " + err.Error(),
		}
		c.JSON(500, message)
		return
	}

	json.Unmarshal(body, &requestMovieBody)

	validate := validator.New()
	err = validate.Struct(requestMovieBody)
	if err != nil {
		message := response.ErrorResponse{
			ErrorCode: 400,
			Message:   "Request body error: " + err.Error(),
		}
		c.JSON(400, message)
		return
	}

	stmt, err := m.Db.Prepare("UPDATE movies SET title = ?, description = ?, rating = ?, image = ?, created_at = ?,  updated_at = ? WHERE id = ?")
	if err != nil {
		message := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something wrong with preparing stmt: " + err.Error(),
		}
		c.JSON(500, message)
		return
	}

	_, err = stmt.Exec(requestMovieBody.Title, requestMovieBody.Description, requestMovieBody.Rating,
		requestMovieBody.Image, requestMovieBody.Created_at, requestMovieBody.Updated_at, param)
	if err != nil {
		message := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something wrong with updating data to database: " + err.Error(),
		}
		c.JSON(500, message)
		return
	}

	c.JSON(200, "OK")
}

func (m MovieController) deleteMovie(c *gin.Context) {
	param := c.Param("ID")

	stmt, err := m.Db.Prepare("DELETE FROM movies WHERE id = ?")
	if err != nil {
		message := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something wrong with preparing stmt: " + err.Error(),
		}
		c.JSON(500, message)
		return
	}

	_, err = stmt.Exec(param)
	if err != nil {
		message := response.ErrorResponse{
			ErrorCode: 500,
			Message:   "Something wrong with deleting data from database: " + err.Error(),
		}
		c.JSON(500, message)
		return
	}

	c.JSON(200, "OK")
}
