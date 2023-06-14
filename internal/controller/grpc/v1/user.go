package v1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/olzzhas/grpc-sneakershop/internal/domain/entity"
	proto_authentication_service "github.com/olzzhas/grpc-sneakershop/service/authentication_service/service"
	proto_user_model "github.com/olzzhas/grpc-sneakershop/service/user_service/model/v1"
	proto_user_service "github.com/olzzhas/grpc-sneakershop/service/user_service/service/v1"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"log"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type userServer struct {
	proto_user_service.UnimplementedUserServiceServer
}

func NewUserServer(unimplementedUserServiceServer proto_user_service.UnimplementedUserServiceServer) *userServer {
	return &userServer{UnimplementedUserServiceServer: unimplementedUserServiceServer}
}

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://postgres:olzhas4@localhost:5432/sneakershop?sslmode=disable")
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
	db, err := connectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT * FROM users ORDER BY id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	users := []*entity.User{}

	for rows.Next() {
		var user entity.User

		err := rows.Scan(
			&user.ID,
			&user.CreatedAt,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Activated,
			&user.Version,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return &proto_user_service.GetUsersResponse{
		User: []*proto_user_model.User{
			// don't forget to fix
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

	query := `
		INSERT INTO users (name, email, password_hash, activated)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version
	`

	var password entity.Password
	passHash, err := password.Set(req.Password)
	if err != nil {
		return nil, err
	}

	reqUser := entity.User{
		Name:      req.Name,
		Password:  password,
		Email:     req.Email,
		Activated: false,
	}

	args := []any{reqUser.Name, reqUser.Email, passHash, false}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = db.QueryRowContext(ctx, query, args...).Scan(&reqUser.ID, &reqUser.CreatedAt, &reqUser.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, ErrDuplicateEmail
		default:
			return nil, err
		}
	}

	conn, err := grpc.Dial("localhost:8082", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto_authentication_service.NewAuthenticationServiceClient(conn)

	response, err := client.CreateAuthenticationToken(context.Background(), &proto_authentication_service.CreateAuthenticationTokenRequest{UserId: reqUser.ID, Password: req.Password})
	fmt.Println(response)

	rabbitconn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
	}
	defer conn.Close()

	ch, err := rabbitconn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        response.Token.Token,
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	log.Printf(" [x] Sent %s\n", string(response.Token.Token))

	return &proto_user_service.CreateUserResponse{
		Id: reqUser.ID,
	}, nil

}
