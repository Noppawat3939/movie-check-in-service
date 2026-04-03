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
	ID         uuid.UUID         `json:"id"`
	ShowTimeID uuid.UUID         `json:"showtime_id"`
	SeatID     uuid.UUID         `json:"seat_id"`
	Status     ReservationStatus `json:"status"`

	ReservedAt time.Time `json:"reserved_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
