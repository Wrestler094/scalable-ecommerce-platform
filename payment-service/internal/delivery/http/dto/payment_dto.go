package dto

// ====== Pay ======

type PayRequest struct {
	OrderID        int64   `json:"order_id" validate:"required,gt=0"`
	Amount         float64 `json:"amount" validate:"required,gt=0"`
	IdempotencyKey string  `json:"idempotency_key" validate:"required"`
}

type PayResponse struct {
	Message string `json:"message"`
}
