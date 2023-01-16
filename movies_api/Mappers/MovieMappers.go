package mapper

import (
	response "movies/Contract/Response"
	model "movies/Model"
)

func MapMovies(m model.Movie) response.Movie {
	response := response.Movie{
		Id:          m.Id,
		Title:       m.Title,
		Description: m.Description,
		Rating:      m.Rating,
		Image:       m.Image,
		Created_at:  m.Created_at.Format("2006-01-02 15:04:05"),
		Updated_at:  m.Updated_at.Format("2006-01-02 15:04:05"),
	}

	return response
}
