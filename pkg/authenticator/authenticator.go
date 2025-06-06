package authenticator

type Role string

const (
	Admin Role = "admin"
	User  Role = "user"
)

type Authenticator interface {
	Validate(token string) (userID int64, role string, err error)
}
