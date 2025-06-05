package authenticator

type Authenticator interface {
	Validate(token string) (userID int64, role string, err error)
}
