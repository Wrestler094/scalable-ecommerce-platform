package httphelper

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse — стандартная JSON-обёртка для ошибок.
type ErrorResponse struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

// jsonFieldError — сериализуемая ошибка поля (только для HTTP-ответа).
type jsonFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// RespondJSON — универсальный JSON-ответ.
func RespondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// RespondError — отправка базовой ошибки с сообщением.
func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, ErrorResponse{
		Error: message,
	})
}

// RespondValidationErrors — отправка ошибок валидации в читаемом формате.
func RespondValidationErrors(w http.ResponseWriter, fields []FieldError) {
	if len(fields) == 0 {
		RespondError(w, http.StatusBadRequest, "invalid data")
		return
	}

	RespondJSON(w, http.StatusBadRequest, ErrorResponse{
		Error:   "invalid data",
		Details: marshalFieldErrors(fields),
	})
}

// marshalFieldErrors маппит FieldError в сериализуемые jsonFieldError.
func marshalFieldErrors(fields []FieldError) []jsonFieldError {
	out := make([]jsonFieldError, 0, len(fields))

	for _, f := range fields {
		out = append(out, jsonFieldError{
			Field:   f.Field,
			Message: f.Message,
		})
	}

	return out
}
