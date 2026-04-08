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

type SeatUsecase interface {
	GetSeatAvailable(ctx context.Context, showtimeID uuid.UUID) ([]domain.SeatsAvailableResponse, error)
}

type seatUsecase struct {
	seatRepo        postgresl.SeatRepository
	reservationRepo postgresl.ReservationRepository
	cache           *redis.Cache
}

func NewSeatUsecase(seatRepo postgresl.SeatRepository, reservationRepo postgresl.ReservationRepository, cache *redis.Cache) SeatUsecase {
	return &seatUsecase{seatRepo, reservationRepo, cache}
}

func (u *seatUsecase) GetSeatAvailable(ctx context.Context, showtimeID uuid.UUID) ([]domain.SeatsAvailableResponse, error) {
	seats, err := u.getSeatsByShowTimeID(ctx, showtimeID)
	if err != nil {
		return nil, err
	}

	reservations, err := u.reservationRepo.ListReservationByShowtimeID(ctx, showtimeID)
	if err != nil {
		return nil, err
	}

	// map reservations status confirmed to set
	reservedSet := make(map[uuid.UUID]struct{})
	for _, reserve := range reservations {
		if reserve.Status == domain.ReservationConfirmed && reserve.SeatID != uuid.Nil {
			reservedSet[reserve.SeatID] = struct{}{}
		}
	}

	var resp []domain.SeatsAvailableResponse
	// map reserved to seats before response
	for _, seat := range seats {
		_, isReserved := reservedSet[seat.ID]

		resp = append(resp, domain.SeatsAvailableResponse{
			ID:         seat.ID,
			ShowTimeID: seat.ShowTimeID,
			SeatNumber: seat.SeatNumber,
			IsReserved: isReserved,
			CreatedAt:  seat.CreatedAt,
		})
	}

	return resp, nil
}

// helpers
func (u *seatUsecase) getSeatsByShowTimeID(ctx context.Context, showtimeID uuid.UUID) ([]domain.Seat, error) {
	// get cached
	cacheKey := fmt.Sprintf("seat:showtime:%s", showtimeID)
	var seats []domain.Seat
	if err := u.cache.Get(ctx, cacheKey, seats); err == nil {
		return seats, nil
	}

	seats, err := u.seatRepo.FindByShowtimeID(ctx, showtimeID)
	if err != nil {
		return nil, err
	}

	// set cache seats by showtime
	_ = u.cache.Set(ctx, cacheKey, seats, 3*time.Hour)

	return seats, nil
}
