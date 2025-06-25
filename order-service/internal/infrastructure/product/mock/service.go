package mock

import (
	"context"
	"math/rand"
	"time"
)

// MockProductService возвращает случайные цены от 1000.00 до 5000.00
type MockProductService struct {
	seededRand *rand.Rand
	minPrice   float64
	maxPrice   float64
}

func NewMockProductService() *MockProductService {
	return &MockProductService{
		seededRand: rand.New(rand.NewSource(time.Now().UnixNano())),
		minPrice:   1000.0,
		maxPrice:   5000.0,
	}
}

func (s *MockProductService) GetPrice(_ context.Context, _ int64) (float64, error) {
	price := s.minPrice + s.seededRand.Float64()*(s.maxPrice-s.minPrice)
	return round(price, 2), nil
}

func round(f float64, places int) float64 {
	shift := float64(1)
	for i := 0; i < places; i++ {
		shift *= 10
	}
	return float64(int(f*shift+0.5)) / shift
}
