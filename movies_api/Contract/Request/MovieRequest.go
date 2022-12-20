package request

type Movie struct {
	Id          int32   `json:"id" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Rating      float64 `json:"rating" validate:"required"`
	Image       string  `json:"image" default:""`
	Created_at  string  `json:"created_at" validate:"required"`
	Updated_at  string  `json:"updated_at" validate:"required"`
}
