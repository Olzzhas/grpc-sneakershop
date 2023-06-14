package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"testing"
	"time"

	"github.com/olzzhas/grpc-sneakershop/internal/domain/entity"
	proto_authentication_model "github.com/olzzhas/grpc-sneakershop/service/authentication_service/model"
	proto_authentication_service "github.com/olzzhas/grpc-sneakershop/service/authentication_service/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type authServer struct {
	proto_authentication_service.UnimplementedAuthenticationServiceServer
}

func NewAuthenticationServer(unimplementedAuthenticationServiceServer proto_authentication_service.UnimplementedAuthenticationServiceServer) *authServer {
	return &authServer{UnimplementedAuthenticationServiceServer: unimplementedAuthenticationServiceServer}
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

func TestCreateAuthenticationToken(t *testing.T) {
	server := NewAuthenticationServer(proto_authentication_service.UnimplementedAuthenticationServiceServer{})
	ctx := context.Background()
	req := &proto_authentication_service.CreateAuthenticationTokenRequest{
		UserId: 123,
	}

	resp, err := server.CreateAuthenticationToken(ctx, req)

	if err != nil {
		t.Errorf("CreateAuthenticationToken returned an unexpected error: %v", err)
	}

	if resp.Token == nil {
		t.Error("CreateAuthenticationToken response does not contain a token")
	} else {

	}
}

func TestGenerateToken(t *testing.T) {
	userID := uint32(123)
	ttl := 30 * time.Minute
	scope := "authentication"

	token, err := generateToken(userID, ttl, scope)

	if err != nil {
		t.Errorf("generateToken returned an unexpected error: %v", err)
	}

	if token == nil {
		t.Error("generateToken returned a nil token")
	} else {

	}
}
