package response

type Movie struct {
	Id          int32   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	Image       string  `json:"image" default:""`
	Created_at  string  `json:"created_at"`
	Updated_at  string  `json:"updated_at"`
}
