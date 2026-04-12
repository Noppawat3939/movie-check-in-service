package domain

import "errors"

var (
	ErrBodyInvalid            = errors.New("invalid request body")
	ErrSeatAlreadyReserved    = errors.New("seat already reserved")
	ErrLockNotAquired         = errors.New("seat is being processed, please try again")
	ErrReservationNotEligible = errors.New("reservation cannot be changed")
	ErrReservationNotFound    = errors.New("reservatoin not found")
	ErrMovieNotFound          = errors.New("movie not found")
	ErrSeatUnchanged          = errors.New("new seat is the same current seat")
)
