package hasher

import (
	"errors"
	"fmt"

	"user-service/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

var _ domain.PasswordHasher = (*BcryptHasher)(nil)

type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	const op = "BcryptHasher.Hash"

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: failed to hash password: %w", op, err)
	}

	return string(bytes), nil
}

func (h *BcryptHasher) Compare(hashed, plain string) (bool, error) {
	const op = "BcryptHasher.Compare"

	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, fmt.Errorf("%s: compare error: %w", op, err)
	}

	return true, nil
}
