package auth

import (
	"context"
	"net/http"
	"shop-auth/internal/models"

	v1 "github.com/legendaryneobatman/shop-proto-repo/gen/go/api/auth/v1"
)

type Handler struct {
	v1.UnimplementedAuthServiceServer
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SignIn(ctx context.Context, in *v1.SignInRequest) (*v1.SignInResponse, error) {
	tokenPair, err := h.service.Authenticate(ctx, SignInRequest{
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
func (h *Handler) SignUp(ctx context.Context, in *v1.SignUpRequest) (*v1.SignUpResponse, error) {
	user, err := h.CreateUser(&models.User{
		Name:     input.Name,
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		return ErrUserAlreadyExists(err)
	}

	c.JSON(http.StatusCreated, SignUpOutput{
		ID: user.ID,
	})
	return nil
}
func Refresh(ctx context.Context) (, error) {}
func Logout(ctx context.Context) (, error) {}
func LogoutAll(ctx context.Context) (, error) {}