package events

type PaymentSuccessfulPayload struct {
	OrderUUID string  `json:"order_uuid"`
	UserID    int64   `json:"user_id"`
	Amount    float64 `json:"amount"`
}
