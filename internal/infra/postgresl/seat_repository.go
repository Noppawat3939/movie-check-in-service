package postgresl

import (
	"check-in/internal/domain"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeatRepository interface {
	FindByShowtimeID(ctx context.Context, showtimeID uuid.UUID) ([]domain.Seat, error)
}

type seatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) SeatRepository {
	return &seatRepository{db}
}

func (r *seatRepository) FindByShowtimeID(ctx context.Context, showtimeID uuid.UUID) ([]domain.Seat, error) {
	var data []domain.Seat

	if err := r.db.WithContext(ctx).Where("showtime_id = ?", showtimeID).Order("seat_number ASC").Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}
