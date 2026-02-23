package auth

import (
	"context"

	v1 "github.com/legendaryneobatman/shop-proto-repo/gen/go/api/auth/v1"
)

type Handler struct {
	v1.UnimplementedAuthServiceServer
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (s *Handler) SignIn(ctx context.Context, in *v1.SignInRequest) (*v1.SignInResponse, error) {
	tokenPair, err := s.service.Authenticate(ctx, SignInRequest{
		username: in.Username,
		password: in.Password,
	})
	if err != nil {
		return nil, err
	}

	return &v1.SignInResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
