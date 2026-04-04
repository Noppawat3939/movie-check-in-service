package usecase

import (
	"check-in/internal/domain"
	"check-in/internal/infra/postgresl"
	"context"
	"errors"
	"fmt"
	"time"
)

type ReverationUsecase interface {
	CreateReveration(ctx context.Context, req domain.CreateReservationRequest) (*domain.CreateReservationResponse, error)
}

type reverationUsecase struct {
	reservationRepo postgresl.ReservationRepository
}

func NewReverationUsecase(reservationRepo postgresl.ReservationRepository) ReverationUsecase {
	return &reverationUsecase{reservationRepo}
}

func (u *reverationUsecase) CreateReveration(ctx context.Context, req domain.CreateReservationRequest) (*domain.CreateReservationResponse, error) {
	count, err := u.reservationRepo.CountByShowTimeAndSeat(ctx, req.ShowTimeID, req.SeatID)
	if err != nil {
		return nil, err
	}
	fmt.Printf("count %d", count)
	if count > 0 {
		return nil, errors.New("seat already reserved")
	}

	reservation := &domain.Reservation{
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
