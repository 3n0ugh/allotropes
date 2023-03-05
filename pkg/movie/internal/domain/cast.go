package domain

type Cast struct {
	CastOrder int    `json:"CastOrder"`
	Gender    string `json:"Gender"`
	Name      string `json:"Name"`
	PersonID  int    `json:"PersonID"`
}
