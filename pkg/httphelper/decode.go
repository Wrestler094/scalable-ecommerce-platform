package httphelper

import (
	"encoding/json"
	"net/http"
)

// DecodeJSON декодирует JSON из тела запроса в структуру типа T.
// Возвращает ошибку, если данные некорректны.
func DecodeJSON[T any](r *http.Request, w http.ResponseWriter) (T, error) {
	var req T
	var zero T

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		return zero, err
	}

	return req, nil
}
