package v1

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"github.com/olzzhas/grpc-sneakershop/internal/domain/entity"
	proto_authentication_model "github.com/olzzhas/grpc-sneakershop/service/authentication_service/model"
	proto_authentication_service "github.com/olzzhas/grpc-sneakershop/service/authentication_service/service"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type authServer struct {
	proto_authentication_service.UnimplementedUserServiceServer
}

func (s *authServer) CreateAuthenticationToken(ctx context.Context, req *proto_authentication_service.CreateAuthenticationTokenRequest) (*proto_authentication_service.CreateAuthenticationTokenResponse, error) {
	db, err := connectDB()
	if err != nil {
		return nil, fmt.Errorf("error while connecting to database: %s", err)
	}

	defer db.Close()

	token, err := generateToken(req.UserId, 30*time.Minute, "authentication")
	if err != nil {
		return nil, err
	}

	return &proto_authentication_service.CreateAuthenticationTokenResponse{
		Token: &proto_authentication_model.Token{
			Token:  token.Hash,
			Expiry: timestamppb.New(token.Expiry),
		},
	}, nil

}

func generateToken(user_id uint32, ttl time.Duration, scope string) (*entity.Token, error) {
	token := &entity.Token{
		UserID: user_id,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}
