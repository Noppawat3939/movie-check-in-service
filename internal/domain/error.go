package domain

import "errors"

const (
	ErrBodyInvalidMsg            = "invalid request body"
	ErrSeatAlreadyReservedMsg    = "seat already reserved"
	ErrLockNotAquiredMsg         = "seat is being processed, please try again"
	ErrReservationNotFoundMsg    = "reservatoin not found"
	ErrReservationNotEligibleMsg = "reservation cannot be changed"
	ErrMovieNotFoundMsg          = "movie not found"
	ErrSeatUnchangedMsg          = "new seat is the same current seat"
)

var (
	ErrSeatAlreadyReserved    = errors.New(ErrSeatAlreadyReservedMsg)
	ErrLockNotAquired         = errors.New(ErrLockNotAquiredMsg)
	ErrReservationNotEligible = errors.New(ErrReservationNotEligibleMsg)
	ErrReservationNotFound    = errors.New(ErrReservationNotFoundMsg)
	ErrMovieNotFound          = errors.New(ErrMovieNotFoundMsg)
	ErrSeatUnchanged          = errors.New(ErrSeatUnchangedMsg)
)
