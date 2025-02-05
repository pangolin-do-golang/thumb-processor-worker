package domain

type Event struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	Metadata interface{}
}
