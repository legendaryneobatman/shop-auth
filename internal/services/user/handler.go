package user

import (
	"context"

	v1 "github.com/legendaryneobatman/shop-proto-repo/gen/go/api/auth/v1"
)

type Handler struct {
	service Service
	v1.UnimplementedUserServiceServer
}

func (h *Handler) Get(ctx context.Context, in v1.GetUserRequest) (v1.GetUserResponse, error) {
	h.service
}
func (h *Handler) GetById(ctx context.Context, in v1.GetUserRequest) (v1.GetUserResponse, error) {
}
func (h *Handler) Create(ctx context.Context, in v1.CreateUserRequest) (v1.CreateUserResponse, error) {
}
func (h *Handler) Edit(ctx context.Context, in v1.UpdateUserRequest) (v1.UpdateUserResponse, error) {}
func (h *Handler) Delete(ctx context.Context, in v1.DeleteUserRequest) (v1.DeleteUserResponse, error) {
}
