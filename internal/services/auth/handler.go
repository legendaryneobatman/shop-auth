package auth

import (
	"context"
	"errors"
	"strconv"

	v1 "github.com/legendaryneobatman/shop-proto-repo/gen/go/api/auth/v1"
	"golang.org/x/crypto/bcrypt"

	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	v1.UnimplementedAuthServiceServer
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SignIn(ctx context.Context, in *v1.SignInRequest) (*v1.SignInResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}

	tokenPair, err := h.service.Authenticate(ctx, SignInRequest{
		username: in.Username,
		password: in.Password,
	})
	if err != nil {
		var mismatchErr error = bcrypt.ErrMismatchedHashAndPassword
		if errors.Is(err, mismatchErr) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.SignInResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (h *Handler) SignUp(ctx context.Context, in *v1.SignUpRequest) (*v1.SignUpResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}

	id, err := h.service.RegisterUser(ctx, SignUpInput{
		name:     in.Name,
		username: in.Username,
		password: in.Password,
	})
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.SignUpResponse{Id: strconv.Itoa(id)}, nil
}

func (h *Handler) Refresh(ctx context.Context, in *v1.RefreshRequest) (*v1.RefreshResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}

	tokenPair, err := h.service.RefreshTokens(in.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &v1.RefreshResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (h *Handler) Logout(ctx context.Context, in *v1.LogoutRequest) (*v1.LogoutResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}
	userID, err := h.service.ParseTokenForUserID(in.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}

	if err := h.service.RevokeAllTokens(userID); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.LogoutResponse{}, nil
}

func (h *Handler) LogoutAll(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "LogoutAll is not supported yet")
}
