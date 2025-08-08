package valueObjects

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	PlainText string
	Hash      []byte
}

func NewPassword(plainText string) *Password {
	return &Password{PlainText: plainText}
}

func (p *Password) GetHash() []byte {
	return p.Hash
}

func (p *Password) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.PlainText), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	p.Hash = hash

	return nil
}

func (p *Password) VerifyPassword(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err //internal server error
		}
	}

	return true, nil
}
