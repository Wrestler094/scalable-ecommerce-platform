package authenticator

type Role string

const (
	Admin Role = "admin"
	User  Role = "user"
)

// Headers используемые для передачи контекста между Gateway и сервисами
const (
	// HeaderAuthenticated - флаг что пользователь аутентифицирован
	HeaderAuthenticated = "X-Authenticated"
	// HeaderUserID - ID пользователя
	HeaderUserID = "X-User-ID"
	// HeaderUserRole - роль пользователя
	HeaderUserRole = "X-User-Role"
)

type Authenticator interface {
	Validate(token string) (userID int64, role string, err error)
}
