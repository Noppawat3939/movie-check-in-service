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
	ListReservation(ctx context.Context, showtimeID uuid.UUID) ([]domain.Reservation, error)
	CreateReservation(ctx context.Context, req domain.CreateReservationRequest) (*domain.CreateReservationResponse, error)
	ChangeReservation(ctx context.Context, req domain.ChangeReservationRequest) (*domain.CreateReservationResponse, error)
}

type reservationUsecase struct {
	reservationRepo        postgresl.ReservationRepository
	lockRepo               redis.LockRepository
	reservationLockLogRepo postgresl.ReservationLockLogRepository
}

func NewReservationUsecase(reservationRepo postgresl.ReservationRepository, lockRepo redis.LockRepository, reservationLockLogRepo postgresl.ReservationLockLogRepository) ReservationUsecase {
	return &reservationUsecase{reservationRepo, lockRepo, reservationLockLogRepo}
}

func (u *reservationUsecase) ListReservation(ctx context.Context, showtimeID uuid.UUID) ([]domain.Reservation, error) {
	return u.reservationRepo.ListReservationByShowtimeID(ctx, showtimeID)
}

func (u *reservationUsecase) CreateReservation(ctx context.Context, req domain.CreateReservationRequest) (*domain.CreateReservationResponse, error) {
	var reservation *domain.Reservation
	// prevent concurrency requests to reserve same showtime and seat
	err := u.withSeatLock(ctx, req.ShowTimeID, req.SeatID, func() error {
		// ensure has 1 request below processes
		count, err := u.reservationRepo.CountByShowTimeAndSeat(ctx, req.ShowTimeID, req.SeatID)
		if err != nil {
			return err
		}
		if count > 0 {
			return domain.ErrSeatAlreadyReserved
		}

		reservation = &domain.Reservation{
			ID:         uuid.New(),
			ShowTimeID: req.ShowTimeID,
			SeatID:     req.SeatID,
			Status:     domain.ReservationConfirmed,
			ReservedAt: time.Now(),
		}

		return u.reservationRepo.Create(ctx, reservation)
	})

	if err != nil {
		return nil, err
	}

	return &domain.CreateReservationResponse{
		ID:         reservation.ID,
		ShowtimeID: reservation.ShowTimeID,
		SeatID:     reservation.SeatID,
		Status:     reservation.Status,
		ReservedAt: reservation.ReservedAt,
	}, nil
}

func (u *reservationUsecase) ChangeReservation(ctx context.Context, req domain.ChangeReservationRequest) (*domain.CreateReservationResponse, error) {
	// check current reseved
	existing, err := u.reservationRepo.FindByID(ctx, req.ReservationID)
	if err != nil {
		return nil, domain.ErrReservationNotFound
	}
	if existing.Status != domain.ReservationConfirmed {
		return nil, domain.ErrReservationNotEligible
	}
	if req.NewSeatID == existing.SeatID {
		return nil, domain.ErrSeatUnchanged
	}

	var newReservation *domain.Reservation

	err = u.withSeatLock(ctx, existing.ShowTimeID, req.NewSeatID, func() error {
		// check new seat available
		count, err := u.reservationRepo.CountByShowTimeAndSeat(ctx, existing.ShowTimeID, req.NewSeatID)
		if err != nil {
			return err
		}

		if count > 0 {
			return domain.ErrSeatAlreadyReserved
		}

		// cancel old and create new seat
		newReservation = &domain.Reservation{
			ID:         uuid.New(),
			ShowTimeID: existing.ShowTimeID,
			SeatID:     req.NewSeatID,
			Status:     domain.ReservationConfirmed,
			ReservedAt: time.Now(),
		}

		return u.reservationRepo.CancelAndCreate(ctx, req.ReservationID, newReservation)

	})
	if err != nil {
		return nil, err
	}

	return &domain.CreateReservationResponse{
		ID:         newReservation.ID,
		ShowtimeID: newReservation.ShowTimeID,
		SeatID:     newReservation.SeatID,
		Status:     newReservation.Status,
		ReservedAt: newReservation.ReservedAt,
	}, nil
}

func buildLockArgs(showtimeID uuid.UUID, seatID uuid.UUID) (string, string) {
	lockKey := fmt.Sprintf("lock:showtime:%s:seat:%s", showtimeID, seatID)
	lockValue := uuid.NewString()
	return lockKey, lockValue
}

func (u *reservationUsecase) withSeatLock(ctx context.Context, showtimeID, seatID uuid.UUID, fn func() error) error {
	key, value := buildLockArgs(showtimeID, seatID)
	// lock data in 10sec
	acuired, err := u.lockRepo.AcquireLock(ctx, key, value, 10*time.Second)
	if err != nil {
		// failed acquire lock
		go u.lockFailedLog(context.Background(), key, showtimeID, seatID) // use background prevent request cancelled before create log
		return err
	}
	if !acuired {
		go u.lockFailedLog(context.Background(), key, showtimeID, seatID)
		return domain.ErrLockNotAquired
	}

	// acquired then create log
	go u.lockAquiredLog(context.Background(), key, showtimeID, seatID)

	defer func() {
		// release lock complete then create log
		if err := u.lockRepo.ReleaseLock(ctx, key, value); err == nil {
			go u.lockReleasedLog(context.Background(), key, showtimeID, seatID)
		}
	}()

	return fn()
}

func (u *reservationUsecase) lockAquiredLog(ctx context.Context, lockKey string, showtimeID, seatID uuid.UUID) error {
	now := time.Now()
	return u.reservationLockLogRepo.Create(ctx, &domain.ReservationLockLog{
		ID:         uuid.New(),
		ShowTimeID: showtimeID,
		SeatID:     seatID,
		LockKey:    lockKey,
		Status:     domain.LockStatusAquired,
		AquiredAt:  now,
		CreatedAt:  now,
	})
}

func (u *reservationUsecase) lockFailedLog(ctx context.Context, lockKey string, showtimeID, seatID uuid.UUID) error {
	return u.reservationLockLogRepo.Create(ctx, &domain.ReservationLockLog{
		ID:         uuid.New(),
		ShowTimeID: showtimeID,
		SeatID:     seatID,
		LockKey:    lockKey,
		Status:     domain.LockStatusFailed,
		CreatedAt:  time.Now(),
	})
}

func (u *reservationUsecase) lockReleasedLog(ctx context.Context, lockKey string, showtimeID, seatID uuid.UUID) error {
	now := time.Now()
	return u.reservationLockLogRepo.Create(ctx, &domain.ReservationLockLog{
		ID:         uuid.New(),
		ShowTimeID: showtimeID,
		SeatID:     seatID,
		LockKey:    lockKey,
		Status:     domain.LockStatusReleased,
		ReleasedAt: now,
		CreatedAt:  now,
	})
}
