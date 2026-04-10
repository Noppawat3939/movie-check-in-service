package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReservationLockStatus string

const (
	LockStatusAquired  ReservationLockStatus = "acquired"
	LockStatusFailed   ReservationLockStatus = "failed"
	LockStatusReleased ReservationLockStatus = "released"
)

type ReservationLockLog struct {
	ID         uuid.UUID             `json:"id" gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	ShowTimeID uuid.UUID             `json:"showtime_id" gorm:"column:showtime_id"`
	SeatID     uuid.UUID             `json:"seat_id" gorm:"column:seat_id"`
	LockKey    string                `json:"lock_key" gorm:"column:lock_key"`
	Status     ReservationLockStatus `json:"status" gorm:"column:status; type:lock_status"`
	AquiredAt  time.Time             `json:"acquired_at" gorm:"column:acquired_at"`
	ReleasedAt time.Time             `json:"released_at" gorm:"column:released_at"`
	CreatedAt  time.Time             `json:"created_at" gorm:"column:created_at"`
}
