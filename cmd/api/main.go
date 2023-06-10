package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	v1 "github.com/olzzhas/grpc-sneakershop/internal/controller/grpc/v1"
	proto_sneaker_service "github.com/olzzhas/grpc-sneakershop/service/sneaker_service/service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func main() {
	grpcServer := grpc.NewServer()

	proto_sneaker_service.RegisterSneakerServiceServer(
		grpcServer,
		v1.NewUserServer(proto_sneaker_service.UnimplementedSneakerServiceServer{}),
	)

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := proto_sneaker_service.RegisterSneakerServiceHandlerFromEndpoint(context.Background(), mux, "0.0.0.0:8082", opts)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(":8081", mux)
	if err != nil {
		panic(err)
	}

}
