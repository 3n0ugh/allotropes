package internal

type MoviesRequest struct {
	Page int `query:"page"`
	Size int `query:"size"`
}

type MovieByIDRequest struct {
	ID int `path:"id"`
}

type UpdateMovieRequest struct {
	ID    int   `path:"id"`
	Movie Movie `json:"movie"`
}

type DeleteMovieRequest struct {
	ID int `path:"id"`
}

type AddMovieRequest struct {
	Movie Movie `json:"movie"`
}
