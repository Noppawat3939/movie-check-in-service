package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReservationStatus string

const (
	ReservationConfirmed ReservationStatus = "confirmed"
	ReservationCancelled ReservationStatus = "cancelled"
	ReservationExpired   ReservationStatus = "expired"
)

type Reservation struct {
	ID         uuid.UUID         `json:"id" gorm:"column:id"`
	ShowTimeID uuid.UUID         `json:"showtime_id" gorm:"column:showtime_id"`
	SeatID     uuid.UUID         `json:"seat_id" gorm:"column:seat_id"`
	Status     ReservationStatus `json:"status" gorm:"column:status"`
	ReservedAt time.Time         `json:"reserved_at" gorm:"column:reserved_at"`
	CreatedAt  time.Time         `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time         `json:"updated_at" gorm:"column:updated_at"`
}

type CreateReservationRequest struct {
	ShowTimeID uuid.UUID `json:"showtime_id" binding:"required"`
	SeatID     uuid.UUID `json:"seat_id" binding:"required"`
}

type CreateReservationResponse struct {
	ID         uuid.UUID         `json:"id"`
	ShowtimeID uuid.UUID         `json:"showtime_id"`
	SeatID     uuid.UUID         `json:"seat_id"`
	Status     ReservationStatus `json:"status"`
	ReservedAt time.Time         `json:"reserved_at"`
}
