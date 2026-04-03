package domain

import (
	"time"

	"github.com/google/uuid"
)

type Seat struct {
	ID         uuid.UUID `json:"id"`
	ShowTimeID uuid.UUID `json:"showtime_id"`
	SeatNumber string    `json:"seat_number"`
	CreatedAt  time.Time `json:"created_at"`
}
