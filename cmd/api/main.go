package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	v1 "github.com/olzzhas/grpc-sneakershop/internal/controller/grpc/v1"
	proto_sneaker_service "github.com/olzzhas/grpc-sneakershop/service/sneaker_service/service/v1"
	proto_user_service "github.com/olzzhas/grpc-sneakershop/service/user_service/service/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"
)

const (
	grpcHostPort = "0.0.0.0:8082"
)

func main() {
	grpcServer := grpc.NewServer()
	listen, err := net.Listen("tcp", grpcHostPort)
	if err != nil {
		panic(err)
	}

	proto_sneaker_service.RegisterSneakerServiceServer(
		grpcServer,
		v1.NewSneakerServer(proto_sneaker_service.UnimplementedSneakerServiceServer{}),
	)

	proto_user_service.RegisterUserServiceServer(
		grpcServer,
		v1.NewUserServer(proto_user_service.UnimplementedUserServiceServer{}),
	)

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = proto_sneaker_service.RegisterSneakerServiceHandlerFromEndpoint(context.Background(), mux, grpcHostPort, opts)
	if err != nil {
		panic(err)
	}

	err = proto_user_service.RegisterUserServiceHandlerFromEndpoint(context.Background(), mux, grpcHostPort, opts)
	if err != nil {
		panic(err)
	}

	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() (err error) {
		return grpcServer.Serve(listen)
	})
	g.Go(func() (err error) {
		return http.ListenAndServe(":8081", mux)
	})

	err = g.Wait()
	if err != nil {
		panic(err)
	}

}
