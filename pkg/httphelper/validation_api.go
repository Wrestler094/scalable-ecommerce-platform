package httphelper

// FieldError описывает ошибку конкретного поля.
type FieldError struct {
	Field   string
	Message string
}

// Validator — интерфейс валидации на уровне HTTP.
// Возвращает слайс FieldError, если есть ошибки.
type Validator interface {
	Validate(any) []FieldError
}
