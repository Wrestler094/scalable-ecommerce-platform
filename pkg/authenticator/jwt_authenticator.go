package authenticator

import (
	"errors"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidClaims        = errors.New("invalid claims")
	ErrInvalidSubClaim      = errors.New("invalid sub claim")
	ErrUnexpectedSignMethod = errors.New("unexpected signing method")
)

type accessTokenClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// jwtAuthenticator implements Authenticator
type jwtAuthenticator struct {
	secret        string
	signingMethod jwt.SigningMethod
}

func NewJWTAuthenticator(secret string) Authenticator {
	return &jwtAuthenticator{
		secret:        secret,
		signingMethod: jwt.SigningMethodHS256,
	}
}

func (a *jwtAuthenticator) Validate(tokenStr string) (int64, string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &accessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != a.signingMethod.Alg() {
			return nil, ErrUnexpectedSignMethod
		}
		return []byte(a.secret), nil
	})
	if err != nil || !token.Valid {
		return 0, "", ErrInvalidToken
	}

	claims, ok := token.Claims.(*accessTokenClaims)
	if !ok {
		return 0, "", ErrInvalidClaims
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, "", ErrInvalidSubClaim
	}

	return userID, claims.Role, nil
}
