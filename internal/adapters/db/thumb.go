package db

import (
	"github.com/google/uuid"
	"github.com/pangolin-do-golang/thumb-processor-worker/internal/core/thumb"
)

type PostgresThumbRepository struct {
	db IDB
}

type ThumbPostgres struct {
	BaseModel
}

func (op ThumbPostgres) TableName() string {
	return "thumb"
}

func NewPostgresThumbRepository(db IDB) *PostgresThumbRepository {
	return &PostgresThumbRepository{db: db}
}

func (r *PostgresThumbRepository) Update(_ thumb.Thumb) error {

	return nil
}

func (r *PostgresThumbRepository) Create(_ *thumb.Thumb) (*thumb.Thumb, error) {

	return nil, nil
}

func (r *PostgresThumbRepository) Get(_ uuid.UUID) (*thumb.Thumb, error) {
	return nil, nil
}

func (r *PostgresThumbRepository) GetAll() ([]thumb.Thumb, error) {
	return nil, nil
}
