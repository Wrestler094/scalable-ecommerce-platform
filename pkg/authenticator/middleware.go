package authenticator

import (
	"net/http"
	"strconv"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/httphelper"
)

// extractAuthInfo - это внутренний middleware для извлечения информации о пользователе из заголовков,
// ее валидации и добавления в контекст.
func extractAuthInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(HeaderAuthenticated) != "true" {
			httphelper.RespondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		userIDStr := r.Header.Get(HeaderUserID)
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil || userID == 0 {
			httphelper.RespondError(w, http.StatusUnauthorized, "invalid user ID")
			return
		}

		roleStr := r.Header.Get(HeaderUserRole)
		if roleStr == "" {
			httphelper.RespondError(w, http.StatusUnauthorized, "user role is missing")
			return
		}
		role := Role(roleStr)

		ctx := WithUserID(r.Context(), userID)
		ctx = WithUserRole(ctx, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth проверяет, что пользователь аутентифицирован (через заголовки от Gateway)
// и обогащает контекст запроса.
func RequireAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// extractAuthInfo выполняет всю работу по проверке и извлечению данных.
		return extractAuthInfo(next)
	}
}

// RequireAdmin - это middleware, которое разрешает доступ только администраторам.
func RequireAdmin() func(http.Handler) http.Handler {
	return RequireRoles(Admin)
}

// RequireRoles проверяет, имеет ли аутентифицированный пользователь одну из разрешенных ролей.
func RequireRoles(allowedRoles ...Role) func(http.Handler) http.Handler {
	roleSet := make(map[Role]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		// Сначала выполняем базовую аутентификацию и извлечение данных.
		return extractAuthInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Затем проверяем роль из контекста.
			role, ok := UserRole(r.Context())
			if !ok {
				// Эта ошибка не должна возникать, если extractAuthInfo отработал корректно.
				httphelper.RespondError(w, http.StatusUnauthorized, "role not found in context")
				return
			}

			if _, allowed := roleSet[role]; !allowed {
				httphelper.RespondError(w, http.StatusForbidden, "forbidden: insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		}))
	}
}

// RequireOwnerOrAdmin разрешает доступ, если пользователь является владельцем ресурса или администратором.
// getOwnerID - функция для получения ID владельца ресурса из запроса.
// TODO: For future use. Check before using.
func RequireOwnerOrAdmin(getOwnerID func(r *http.Request) (int64, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Сначала выполняем базовую аутентификацию.
		return extractAuthInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем userID и role из контекста, добавленные extractAuthInfo.
			userID, _ := UserID(r.Context()) // Ошибку можно игнорировать, т.к. extractAuthInfo уже все проверил
			role, _ := UserRole(r.Context())

			// Администраторы имеют доступ всегда.
			if role == Admin {
				next.ServeHTTP(w, r)
				return
			}

			// Если не админ, проверяем, является ли пользователь владельцем.
			ownerID, err := getOwnerID(r)
			if err != nil {
				// Например, если ресурс не найден по параметрам из URL.
				httphelper.RespondError(w, http.StatusNotFound, "resource not found")
				return
			}

			if ownerID != userID {
				httphelper.RespondError(w, http.StatusForbidden, "forbidden: you are not the owner of this resource")
				return
			}

			// Если проверки пройдены, передаем управление дальше.
			next.ServeHTTP(w, r)
		}))
	}
}
