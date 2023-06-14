package entity

import (
	proto_authentication_model "github.com/olzzhas/grpc-sneakershop/service/authentication_service/model"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    uint32
	Expiry    time.Time
	Scope     string
}

func (t *Token) ToProto() *proto_authentication_model.Token {
	return &proto_authentication_model.Token{
		Token:  t.Hash,
		Expiry: timestamppb.New(t.Expiry),
	}
}
