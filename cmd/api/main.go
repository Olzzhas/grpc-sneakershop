package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	proto_sneaker_service "github.com/olzzhas/grpc-sneakershop/service/sneaker_service/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func main() {
	grpcServer := grpc.NewServer()

	proto_sneaker_service.RegisterSneakerServiceServer(grpcServer, proto_sneaker_service.UnimplementedSneakerServiceServer{})

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
