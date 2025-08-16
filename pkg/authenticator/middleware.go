package authenticator

import (
	"net/http"
	"strconv"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/httphelper"
)

// RequireAuth проверяет заголовки от Gateway и обогащает контекст
func RequireAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Authenticated") != "true" {
				httphelper.RespondError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			userIDStr := r.Header.Get("X-User-ID")
			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil || userID == 0 {
				httphelper.RespondError(w, http.StatusUnauthorized, "invalid user ID")
				return
			}

			roleStr := r.Header.Get("X-User-Role")
			role := Role(roleStr)

			ctx := WithUserID(r.Context(), userID)
			ctx = WithUserRole(ctx, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAdmin - только для админов
func RequireAdmin() func(http.Handler) http.Handler {
	return RequireRoles(Admin)
}

// RequireRoles проверяет роли (без JWT, используя заголовки)
func RequireRoles(allowedRoles ...Role) func(http.Handler) http.Handler {
	roleSet := make(map[Role]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Сначала базовая аутентификация
			RequireAuth()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Затем проверка роли
				role, ok := UserRole(r.Context())
				if !ok {
					httphelper.RespondError(w, http.StatusUnauthorized, "role not found")
					return
				}

				if _, allowed := roleSet[role]; !allowed {
					httphelper.RespondError(w, http.StatusForbidden, "forbidden")
					return
				}

				next.ServeHTTP(w, r)
			})).ServeHTTP(w, r)
		})
	}
}

// TODO: For future use. Check before using.
// RequireOwnerOrAdmin - для ресурсов, где доступ только у владельца или админа
func RequireOwnerOrAdmin(getOwnerID func(r *http.Request) (int64, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Базовая аутентификация
			if r.Header.Get("X-Authenticated") != "true" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			userIDStr := r.Header.Get("X-User-ID")
			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil || userID == 0 {
				http.Error(w, "invalid user ID", http.StatusUnauthorized)
				return
			}

			roleStr := r.Header.Get("X-User-Role")
			role := Role(roleStr)

			// Админы могут всё
			if role != Admin {
				// Проверяем владельца
				ownerID, err := getOwnerID(r)
				if err != nil {
					http.Error(w, "resource not found", http.StatusNotFound)
					return
				}

				if ownerID != userID {
					http.Error(w, "forbidden", http.StatusForbidden)
					return
				}
			}

			ctx := WithUserID(r.Context(), userID)
			ctx = WithUserRole(ctx, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
