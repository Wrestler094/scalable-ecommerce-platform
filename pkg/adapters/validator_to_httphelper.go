package adapters

import (
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/httphelper"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/validator"
)

// HttpValidatorAdapter адаптирует интерфейс validator.Validator к интерфейсу httphelper.Validator.
type HttpValidatorAdapter struct {
	inner validator.Validator
}

// NewHttpValidatorAdapter возвращает адаптированную реализацию httphelper.Validator
func NewHttpValidatorAdapter(inner validator.Validator) httphelper.Validator {
	return &HttpValidatorAdapter{inner}
}

// Validate адаптирует ошибки из validator.FieldError в httphelper.FieldError.
func (a *HttpValidatorAdapter) Validate(i any) []httphelper.FieldError {
	src := a.inner.Validate(i)
	if src == nil {
		return nil
	}

	out := make([]httphelper.FieldError, len(src))
	for i, err := range src {
		out[i] = httphelper.FieldError{
			Field:   err.Field(),
			Message: err.Message(),
		}
	}

	return out
}
