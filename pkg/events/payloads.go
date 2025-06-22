package events

type PaymentSuccessfulPayload struct {
	OrderID int64   `json:"order_id"`
	UserID  int64   `json:"user_id"`
	Amount  float64 `json:"amount"`
}
