package internal

type MoviesResponse struct {
	TotalCount int        `json:"totalCount"`
	Movies     []Movie    `json:"movies"`
	Pagination Pagination `json:"pagination"`
}

type MovieByIDResponse struct {
	Movie Movie `json:"movie"`
}

type UpdateMovieResponse struct{}

type DeleteMovieResponse struct{}

type AddMovieResponse struct{}
