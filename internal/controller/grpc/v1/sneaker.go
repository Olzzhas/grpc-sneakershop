package v1

import (
	"context"
	"github.com/olzzhas/grpc-sneakershop/internal/domain/entity"
	proto_sneaker_model "github.com/olzzhas/grpc-sneakershop/service/sneaker_service/model/v1"
	proto_sneaker_service "github.com/olzzhas/grpc-sneakershop/service/sneaker_service/service/v1"
)

type sneakerServer struct {
	proto_sneaker_service.UnimplementedSneakerServiceServer
}

func NewUserServer(unimplementedSneakerServiceServer proto_sneaker_service.UnimplementedSneakerServiceServer) *sneakerServer {
	return &sneakerServer{UnimplementedSneakerServiceServer: unimplementedSneakerServiceServer}
}

func (s *sneakerServer) GetSneakers(ctx context.Context, req *proto_sneaker_service.GetSneakersRequest) (*proto_sneaker_service.GetSneakersResponse, error) {
	example := entity.Sneaker{
		Model: "Under Armour Curry 9",
		Price: 89990,
	}

	return &proto_sneaker_service.GetSneakersResponse{
		Sneaker: []*proto_sneaker_model.Sneaker{
			example.ToProto(),
		},
	}, nil
}

func (s *sneakerServer) UpdateSneaker(ctx context.Context, req *proto_sneaker_service.UpdateSneakerRequest) (*proto_sneaker_service.UpdateSneakerResponse, error) {
	return nil, nil
}
