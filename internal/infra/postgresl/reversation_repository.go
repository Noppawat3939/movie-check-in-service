package postgresl

import (
	"check-in/internal/domain"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReservationRepository interface {
	Create(ctx context.Context, data *domain.Reservation) error
	CountByShowTimeAndSeat(ctx context.Context, showtimeID uuid.UUID, seatID uuid.UUID) (int64, error)
	ListReservationByShowtimeID(ctx context.Context, showtimeID uuid.UUID) ([]domain.Reservation, error)
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReversationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db}
}

func (r *reservationRepository) Create(ctx context.Context, data *domain.Reservation) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *reservationRepository) CountByShowTimeAndSeat(ctx context.Context,
	showtimeID uuid.UUID,
	seatID uuid.UUID) (int64, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&domain.Reservation{}).
		Where(map[string]interface{}{
			"showtime_id": showtimeID,
			"seat_id":     seatID,
			"status":      domain.ReservationConfirmed}).
		Count(&count).
		Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *reservationRepository) ListReservationByShowtimeID(ctx context.Context, showtimeID uuid.UUID) ([]domain.Reservation, error) {
	var data []domain.Reservation

	if err := r.db.WithContext(ctx).Where("showtime_id = ?", showtimeID).Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}
