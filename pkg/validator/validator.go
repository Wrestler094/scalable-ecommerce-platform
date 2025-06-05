package validator

// FieldError — контракт ошибки валидации
type FieldError interface {
	Field() string
	Message() string
}

// Validator — контракт валидации любых структур
type Validator interface {
	Validate(any) []FieldError
}
