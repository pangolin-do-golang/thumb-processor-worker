package domain

type Event struct {
	ID         string `json:"id"`
	Path       string `json:"path"`
	OwnerEmail string `json:"owner_email"`
	Metadata   interface{}
}
