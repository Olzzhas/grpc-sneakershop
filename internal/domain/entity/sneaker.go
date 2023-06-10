package entity

import proto_sneaker_model "github.com/olzzhas/grpc-sneakershop/service/sneaker_service/model/v1"

type Sneaker struct {
	ID    string
	Model string
	Price uint32
}

func (s *Sneaker) ToProto() *proto_sneaker_model.Sneaker {
	return &proto_sneaker_model.Sneaker{
		ID:    s.ID,
		Model: s.Model,
		Price: s.Price,
	}
}
