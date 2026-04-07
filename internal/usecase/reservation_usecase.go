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
	ChangeReservation(ctx context.Context, req domain.ChangeReservationRequest) (*domain.CreateReservationResponse, error)
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

func (u *reservationUsecase) ChangeReservation(ctx context.Context, req domain.ChangeReservationRequest) (*domain.CreateReservationResponse, error) {
	// check current reseved
	existing, err := u.reservationRepo.FindByID(ctx, req.ReservationID)
	if err != nil {
		return nil, err
	}
	if existing.Status != domain.ReservationConfirmed {
		return nil, domain.ErrReservationNotEligible
	}

	// lock new seat
	lockKey := fmt.Sprintf("lock:showtime:%s:seat:%s", existing.ShowTimeID, req.NewSeatID)
	lockValue := uuid.NewString()
	acquired, err := u.lockRepo.AcquireLock(ctx, lockKey, lockValue, 10*time.Second)
	if err != nil {
		return nil, err
	}
	if !acquired {
		return nil, domain.ErrLockNotAquired
	}
	defer u.lockRepo.ReleaseLock(ctx, lockKey, lockValue)

	// check new seat available
	count, err := u.reservationRepo.CountByShowTimeAndSeat(ctx, existing.ShowTimeID, req.NewSeatID)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, domain.ErrSeatAlreadyReserved
	}

	// cancel old and create new seat
	newReservation := &domain.Reservation{
		ID:         uuid.New(),
		ShowTimeID: existing.ShowTimeID,
		SeatID:     req.NewSeatID,
		Status:     domain.ReservationConfirmed,
		ReservedAt: time.Now(),
	}

	if err := u.reservationRepo.CancelAndCreate(ctx, req.ReservationID, newReservation); err != nil {
		return nil, err
	}

	resp := &domain.CreateReservationResponse{
		ID:         newReservation.ID,
		ShowtimeID: newReservation.ShowTimeID,
		SeatID:     newReservation.SeatID,
		Status:     newReservation.Status,
		ReservedAt: newReservation.ReservedAt,
	}

	return resp, nil
}
