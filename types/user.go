package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost     = 12
	minNameLen     = 2
	minPasswordLen = 7
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}

type UserParamsForCreate struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserParamsForUpdate struct {
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
	Email     string `bson:"email" json:"email"`
	Password  string `json:"password"`
}

func NewUserFromParams(params UserParamsForCreate) (*User, error) {
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

func (params UserParamsForCreate) Validate() map[string]string {
	errors := map[string]string{}

	if len(params.FirstName) < minNameLen {
		errors["firstName"] = fmt.Sprintf("First Name length should be at least %d characters", minNameLen)
	}
	if len(params.LastName) < minNameLen {
		errors["lastName"] = fmt.Sprintf("Last Name length should be at least %d characters", minNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("Password length should be at least %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("Email is invalid")
	}

	return errors
}

func UpdatedUserFromParams(params UserParamsForUpdate) (*User, error) {
	return &User{
		FirstName: params.FirstName,
		LastName:  params.LastName,
	}, nil
}

func (params UserParamsForUpdate) Validate() map[string]string {
	errors := map[string]string{}

	if len(params.FirstName) < minNameLen {
		errors["firstName"] = fmt.Sprintf("First Name length should be at least %d characters", minNameLen)
	}
	if len(params.LastName) < minNameLen {
		errors["lastName"] = fmt.Sprintf("Last Name length should be at least %d characters", minNameLen)
	}
	if len(params.Password) > 0 {
		errors["password"] = fmt.Sprintf("Password cannot be updated at this point")
	}
	if len(params.Email) > 0 {
		errors["email"] = fmt.Sprintf("Email cannot be updated")
	}

	return errors
}

// //
func isEmailValid(e string) bool {
	emailReg := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	return emailReg.MatchString(e)
}
