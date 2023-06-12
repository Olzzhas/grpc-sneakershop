package v1

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/olzzhas/grpc-sneakershop/internal/domain/entity"
	proto_user_model "github.com/olzzhas/grpc-sneakershop/service/user_service/model/v1"
	proto_user_service "github.com/olzzhas/grpc-sneakershop/service/user_service/service/v1"
	"math/rand"
)

type userServer struct {
	proto_user_service.UnimplementedUserServiceServer
}

func NewUserServer(unimplementedUserServiceServer proto_user_service.UnimplementedUserServiceServer) *userServer {
	return &userServer{UnimplementedUserServiceServer: unimplementedUserServiceServer}
}

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://postgres:olzhas4@localhost/sneakershop?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("error while sql open: %s", err)
	}

	// Проверка подключения к базе данных
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func (s *userServer) GetUsers(ctx context.Context, req *proto_user_service.GetUsersRequest) (*proto_user_service.GetUsersResponse, error) {
	example := entity.User{
		ID:       1,
		Name:     "Olzhas",
		Age:      19,
		Email:    "211337@astanait.edu.kz",
		Password: "olzzhas",
	}

	return &proto_user_service.GetUsersResponse{
		User: []*proto_user_model.User{
			example.ToProto(),
		},
	}, nil
}

func (s *userServer) UpdateUser(ctx context.Context, req *proto_user_service.UpdateUserRequest) (*proto_user_service.UpdateUserResponse, error) {
	return nil, nil
}

func (s *userServer) CreateUser(ctx context.Context, req *proto_user_service.CreateUserRequest) (*proto_user_service.CreateUserResponse, error) {
	db, err := connectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	testUser := entity.User{
		ID:       uint32(rand.Intn(1000000)),
		Name:     "Olzhas",
		Age:      19,
		Email:    "211337@astanait.edu.kz",
		Password: "olzzhas",
	}

	args := []any{testUser.ID, testUser.Name, testUser.Email, testUser.Password, testUser.Age}

	query := `
		INSERT INTO users (id ,name, email, password, age)
		VALUES ($1, $2, $3, $4, $5)
	`

	db.QueryRow(query, args...)

	return &proto_user_service.CreateUserResponse{
		Id: testUser.ID,
	}, nil

}
