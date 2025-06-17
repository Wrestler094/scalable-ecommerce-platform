package domain

import (
	"errors"
)

var (
	ErrDuplicatePayment              = errors.New("idempotency conflict")
	ErrIdempotencyRegistrationFailed = errors.New("idempotency registration failed")
)
