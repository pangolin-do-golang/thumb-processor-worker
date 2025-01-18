package thumb

import (
	"github.com/google/uuid"
)

type Thumb struct {
	ID uuid.UUID `json:"id"`
}

type IThumbService interface {
	GetThumb(id uuid.UUID) (*Thumb, error)
}
