package authenticator

import (
	"net/http"
	"strings"
)

func RequireRoles(v Authenticator, allowedRoles ...Role) func(http.Handler) http.Handler {
	roleSet := make(map[Role]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			userID, roleStr, err := v.Validate(token)

			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			role := Role(roleStr)
			if _, ok := roleSet[role]; !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			ctx := WithUserID(r.Context(), userID)
			ctx = WithUserRole(ctx, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
