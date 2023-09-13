package types

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost     = 12
	minNameLen     = 2
	minPasswordLen = 7
)

type User struct {
	ID                string `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string `bson:"firstName" json:"firstName"`
	LastName          string `bson:"lastName" json:"lastName"`
	Email             string `bson:"email" json:"email"`
	EncryptedPassword string `bson:"encryptedPassword" json:"-"`
}

type UserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params UserParams) Validate() error {
	if len(params.FirstName) < minNameLen {
		return fmt.Errorf("First Name length should be at leash %d characters", minNameLen)
	}
	if len(params.LastName) < minNameLen {
		return fmt.Errorf("Last Name length should be at leash %d characters", minNameLen)
	}
	if len(params.Password) < minPasswordLen {
		return fmt.Errorf("Last Name length should be at leash %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email) {
		return fmt.Errorf("Email is invalid")
	}

	return nil
}

func isEmailValid(e string) bool {
	emailReg := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	return emailReg.MatchString(e)
}

func NewUserFromParams(params UserParams) (*User, error) {
	encryptedPw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encryptedPw),
	}, nil
}
