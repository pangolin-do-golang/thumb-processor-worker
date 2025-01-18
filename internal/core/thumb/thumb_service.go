package thumb

import (
	"github.com/google/uuid"
)

type Service struct {
}

func NewThumbService() *Service {
	return &Service{}
}

func (s *Service) GetThumb(_ uuid.UUID) (*Thumb, error) {
	var thumb *Thumb

	return thumb, nil
}
