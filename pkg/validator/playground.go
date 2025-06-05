package validator

import (
	"errors"
	"strings"

	v10 "github.com/go-playground/validator/v10"
)

// PlaygroundValidator — адаптер над go-playground/validator.v10.
type PlaygroundValidator struct {
	validator *v10.Validate
}

// NewPlaygroundValidator создаёт новый экземпляр PlaygroundValidator.
func NewPlaygroundValidator() Validator {
	return &PlaygroundValidator{validator: v10.New()}
}

// playgroundFieldError реализует интерфейс FieldError.
type playgroundFieldError struct {
	field   string
	message string
}

func (e playgroundFieldError) Field() string   { return e.field }
func (e playgroundFieldError) Message() string { return e.message }

// Validate валидирует переданную структуру i и возвращает срез ошибок FieldError.
// Если ошибок нет, возвращается nil.
// Если ошибка не связана с валидацией, возвращается пустой срез.
func (v *PlaygroundValidator) Validate(i any) []FieldError {
	err := v.validator.Struct(i)
	if err == nil {
		return nil
	}

	var ve v10.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]FieldError, 0, len(ve))
		for _, fe := range ve {
			out = append(out, playgroundFieldError{
				field:   normalizeFieldName(fe.Field()),
				message: messageForTag(fe.Tag(), fe.Param()),
			})
		}
		return out
	}

	return []FieldError{}
}

// normalizeFieldName преобразует имя поля в lowerCamelCase.
// Например, "Email" в "email".
func normalizeFieldName(field string) string {
	return strings.ToLower(field[:1]) + field[1:]
}
