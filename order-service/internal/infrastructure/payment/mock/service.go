package mock

import (
	"context"
	"fmt"

	"order-service/internal/domain"
)

// MockPaymentService возвращает фиктивную ссылку на оплату
type MockPaymentService struct{}

func NewMockPaymentService() *MockPaymentService {
	return &MockPaymentService{}
}

func (s *MockPaymentService) CreatePayment(_ context.Context, order domain.Order) (string, error) {
	return fmt.Sprintf("https://pay.mock/%s", order.UUID), nil
}
