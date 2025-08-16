package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/authenticator"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
)

// userHeadersToClean - это список заголовков, которые устанавливаются этим middleware.
// Их необходимо очищать в начале каждого запроса, чтобы предотвратить подделку со стороны клиента.
var userHeadersToClean = []string{
	authenticator.HeaderUserID,
	authenticator.HeaderUserRole,
	authenticator.HeaderAuthenticated,
}

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
// - Если токен есть и валиден - устанавливает заголовки с информацией о пользователе.
// - Если токена нет или он невалиден - просто пропускает запрос дальше как анонимный.
func (am *TokenMiddleware) ProcessToken() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "TokenMiddleware.ProcessToken"

			requestID := middleware.GetReqID(r.Context())
			log := am.logger.WithOp(op).WithRequestID(requestID)

			// Очищаем пользовательские заголовки, чтобы клиент не мог их подделать.
			for _, header := range userHeadersToClean {
				r.Header.Del(header)
			}

			authHeader := r.Header.Get("Authorization")
			token := am.extractBearerToken(authHeader)
			if token == "" {
				// Токен не предоставлен, обрабатываем запрос как анонимный.
				// Удаляем заголовок Authorization, чтобы он не просочился дальше.
				r.Header.Del("Authorization")
				next.ServeHTTP(w, r)
				return
			}

			userID, role, err := am.authenticator.Validate(token)
			if err != nil {
				log.WithError(err).Debug("failed to validate token")
				// Токен невалиден, обрабатываем запрос как анонимный.
				r.Header.Del("Authorization")
				next.ServeHTTP(w, r)
				return
			}

			// Токен валиден. Устанавливаем заголовки для внутренних сервисов.
			r.Header.Set(authenticator.HeaderUserID, strconv.FormatInt(userID, 10))
			r.Header.Set(authenticator.HeaderUserRole, role)
			r.Header.Set(authenticator.HeaderAuthenticated, "true")

			next.ServeHTTP(w, r)
		})
	}
}

// extractBearerToken извлекает токен из заголовка Authorization.
// Ожидает формат "Bearer <token>" и возвращает только сам токен.
func (am *TokenMiddleware) extractBearerToken(authHeader string) string {
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return ""
	}

	return strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
}
