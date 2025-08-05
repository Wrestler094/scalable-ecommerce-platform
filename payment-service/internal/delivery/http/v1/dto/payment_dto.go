package dto

import (
	"github.com/google/uuid"
)

// ====== Pay ======

type PayRequest struct {
	OrderUUID      uuid.UUID `json:"order_uuid" validate:"required,uuid4"`
	Amount         float64   `json:"amount" validate:"required,gt=0"`
	IdempotencyKey string    `json:"idempotency_key" validate:"required"`
}

type PayResponse struct {
	Message string `json:"message"`
}
