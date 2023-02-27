package internal

import "time"

type Movie struct {
	Cast []struct {
		CastOrder int    `json:"CastOrder"`
		Gender    string `json:"Gender"`
		Name      string `json:"Name"`
		PersonID  int    `json:"PersonID"`
	} `json:"Cast"`
	Crew []struct {
		Department string `json:"Department"`
		Job        string `json:"Job"`
		PersonID   int    `json:"PersonID"`
	} `json:"Crew"`
	Genres      []string  `json:"Genres"`
	ID          int       `json:"ID"`
	Keywords    []string  `json:"Keywords"`
	Language    string    `json:"Language"`
	ReleaseDate time.Time `json:"ReleaseDate"`
	Runtime     int       `json:"Runtime"`
	ShortStory  string    `json:"ShortStory"`
	Story       string    `json:"Story"`
	Title       string    `json:"Title"`
	Company     []string  `json:"Company"`
	Country     []string  `json:"Country"`
}

func (m Movie) Validate() error { return nil }

type Pagination struct {
	TotalCount int `json:"-"`

	Rel  string `json:"rel"`
	Next string `json:"next"`
	Prev string `json:"prev"`
}

func (p *Pagination) Set() {
}
