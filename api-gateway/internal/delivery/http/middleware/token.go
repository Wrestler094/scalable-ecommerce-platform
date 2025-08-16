package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/authenticator"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
)

type TokenMiddleware struct {
	authenticator authenticator.Authenticator
	logger        logger.Logger
}

func NewTokenMiddleware(authenticator authenticator.Authenticator, logger logger.Logger) *TokenMiddleware {
	return &TokenMiddleware{
		authenticator: authenticator,
		logger:        logger,
	}
}

// ProcessToken обрабатывает токен аутентификации:
// - Если токен есть и валиден - устанавливает заголовки с информацией о пользователе
// - Если токена нет или он невалиден - просто пропускает запрос дальше
func (am *TokenMiddleware) ProcessToken() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "TokenMiddleware.ProcessToken"

			requestID := middleware.GetReqID(r.Context())
			log := am.logger.WithOp(op).WithRequestID(requestID)

			// Очищаем пользовательские заголовки
			userHeaders := []string{
				"X-User-ID",
				"X-User-Role",
				"X-Authenticated",
			}

			for _, header := range userHeaders {
				r.Header.Del(header)
			}

			authHeader := r.Header.Get("Authorization")
			token := am.extractBearerToken(authHeader)
			if token == "" {
				log.Debug("token not found")
				r.Header.Del("Authorization")
				next.ServeHTTP(w, r)
				return
			}

			userID, role, err := am.authenticator.Validate(token)
			if err != nil {
				log.WithError(err).Debug("failed to validate token")
				r.Header.Del("Authorization")
				next.ServeHTTP(w, r)
				return
			}

			// Устанавливаем заголовки
			r.Header.Set("X-User-ID", strconv.FormatInt(userID, 10))
			r.Header.Set("X-User-Role", role)
			r.Header.Set("X-Authenticated", "true")

			next.ServeHTTP(w, r)
		})
	}
}

// extractBearerToken извлекает токен из Authorization заголовка
// Ожидает формат "Bearer <token>" и возвращает только токен
func (am *TokenMiddleware) extractBearerToken(authHeader string) string {
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return ""
	}

	token := strings.TrimPrefix(authHeader, bearerPrefix)
	token = strings.TrimSpace(token)

	if token == "" {
		return ""
	}

	return token
}
