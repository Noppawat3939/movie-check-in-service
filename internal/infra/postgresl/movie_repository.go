package postgresl

import (
	"check-in/internal/domain"
	"context"

	"gorm.io/gorm"
)

type MovieRepository interface {
	FindAll(ctx context.Context) ([]domain.Movie, error)
}

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{db}
}

func (r *movieRepository) FindAll(ctx context.Context) ([]domain.Movie, error) {
	var data []domain.Movie

	if err := r.db.WithContext(ctx).Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}
