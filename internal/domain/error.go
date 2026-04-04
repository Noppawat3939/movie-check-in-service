package domain

import "errors"

const (
	ErrBodyInvalidMsg         = "invalid request body"
	ErrSeatAlreadyReservedMsg = "seat already reserved"
	ErrLockNotAquiredMsg      = "seat is being processed, please try again"
)

var (
	ErrSeatAlreadyReserved = errors.New(ErrSeatAlreadyReservedMsg)
	ErrLockNotAquired      = errors.New(ErrLockNotAquiredMsg)
)
