package postgresl

import (
	"check-in/internal/domain"
	"context"

	"gorm.io/gorm"
)

type ReservationLockLogRepository interface {
	Create(ctx context.Context, data *domain.ReservationLockLog) error
}

type reservationLockLogRepository struct {
	db *gorm.DB
}

func NewReservationLockLogRepository(db *gorm.DB) ReservationLockLogRepository {
	return &reservationLockLogRepository{db}
}

func (r *reservationLockLogRepository) Create(ctx context.Context, data *domain.ReservationLockLog) error {
	return r.db.WithContext(ctx).Create(data).Error
}
