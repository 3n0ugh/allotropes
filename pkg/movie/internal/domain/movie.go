package domain

import "time"

type Movie struct {
	Cast        []Cast    `json:"Cast"`
	Crew        []Crew    `json:"Crew"`
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
