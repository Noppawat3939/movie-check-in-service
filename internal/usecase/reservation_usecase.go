package usecase

import (
	"check-in/internal/domain"
	"check-in/internal/infra/postgresl"
	"check-in/internal/infra/redis"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ReservationUsecase interface {
	CreateReservation(ctx context.Context, req domain.CreateReservationRequest) (*domain.CreateReservationResponse, error)
	ListReservation(ctx context.Context, showtimeID uuid.UUID) ([]domain.Reservation, error)
}

type reservationUsecase struct {
	reservationRepo postgresl.ReservationRepository
	lockRepo        redis.LockRepository
}

func NewReservationUsecase(reservationRepo postgresl.ReservationRepository, lockRepo redis.LockRepository) ReservationUsecase {
	return &reservationUsecase{reservationRepo, lockRepo}
}

func (u *reservationUsecase) CreateReservation(ctx context.Context, req domain.CreateReservationRequest) (*domain.CreateReservationResponse, error) {
	// prevent concurrency requests to reserve same showtime and seat
	lockKey := fmt.Sprintf("lock:showtime:%s:seat:%s", req.ShowTimeID, req.SeatID)
	lockValue := uuid.NewString()
	acquired, err := u.lockRepo.AcquireLock(ctx, lockKey, lockValue, 10*time.Second)
	if err != nil {
		return nil, err
	}
	if !acquired {
		return nil, domain.ErrLockNotAquired
	}
	// clear lock
	defer u.lockRepo.ReleaseLock(ctx, lockKey, lockValue)

	// ensure has 1 request below processes
	count, err := u.reservationRepo.CountByShowTimeAndSeat(ctx, req.ShowTimeID, req.SeatID)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, domain.ErrSeatAlreadyReserved
	}

	reservation := &domain.Reservation{
		ID:         uuid.New(),
		ShowTimeID: req.ShowTimeID,
		SeatID:     req.SeatID,
		Status:     domain.ReservationConfirmed,
		ReservedAt: time.Now(),
	}

	if err := u.reservationRepo.Create(ctx, reservation); err != nil {
		return nil, err
	}

	resp := &domain.CreateReservationResponse{
		ID:         reservation.ID,
		ShowtimeID: reservation.ShowTimeID,
		SeatID:     reservation.SeatID,
		Status:     reservation.Status,
		ReservedAt: reservation.ReservedAt,
	}

	return resp, nil
}

func (u *reservationUsecase) ListReservation(ctx context.Context, showtimeID uuid.UUID) ([]domain.Reservation, error) {
	return u.reservationRepo.ListReservationByShowtimeID(ctx, showtimeID)
}
