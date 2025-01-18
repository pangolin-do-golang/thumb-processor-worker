package db

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
