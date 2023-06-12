package entity

import (
	proto_user_model "github.com/olzzhas/grpc-sneakershop/service/user_service/model/v1"
)

type User struct {
	ID       uint32
	Name     string
	Email    string
	Password string
	Age      uint32
}

func (u *User) ToProto() *proto_user_model.User {
	return &proto_user_model.User{
		Id:       u.ID,
		Name:     u.Name,
		Age:      u.Age,
		Email:    u.Email,
		Password: u.Password,
	}
}
