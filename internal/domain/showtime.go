package domain

import (
	"time"

	"github.com/google/uuid"
)

type ShowTime struct {
	ID          uuid.UUID `json:"id"`
	MovieID     uuid.UUID `json:"movie_id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	TotalSeats  int       `json:"total_seats"`
	BookedSeats int       `json:"booked_seats"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
