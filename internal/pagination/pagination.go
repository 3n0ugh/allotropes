package pagination

type Model struct {
	TotalCount int `json:"-"`

	Rel  string `json:"rel"`
	Next string `json:"next"`
	Prev string `json:"prev"`
}

func (p *Model) Set() {
}
