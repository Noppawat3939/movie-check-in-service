package domain

import (
	"time"

	"github.com/google/uuid"
)

type Seat struct {
	ID         uuid.UUID `json:"id" gorm:"column:id"`
	ShowTimeID uuid.UUID `json:"showtime_id" gorm:"column:showtime_id"`
	SeatNumber string    `json:"seat_number" gorm:"column:seat_number"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
}
