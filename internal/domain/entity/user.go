package entity

import (
	"errors"
	proto_user_model "github.com/olzzhas/grpc-sneakershop/service/user_service/model/v1"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        uint32    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  Password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

func (u *User) ToProto() *proto_user_model.User {
	return &proto_user_model.User{
		Id:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

type Password struct {
	plaintext *string
	hash      []byte
}

func (p *Password) Set(plaintextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return nil, err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return hash, nil
}

func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
