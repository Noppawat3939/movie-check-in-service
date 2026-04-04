package domain

import "errors"

const (
	ErrBodyInvalidMsg         = "invalid request body"
	ErrSeatAlreadyReservedMsg = "seat already reserved"
)

var (
	ErrSeatAlreadyReserved = errors.New(ErrSeatAlreadyReservedMsg)
)
