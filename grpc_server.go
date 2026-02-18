package shop

import (
	"context"
	"net"

	v1 "github.com/legendaryneobatman/shop-proto-repo/gen/go/api/auth/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	v1.UnimplementedAuthServer
}

func (s *server) SignIn(ctx context.Context, in *v1.SignInRequest) (*v1.SignInResponse, error) {
	return &v1.SignInResponse{
		AccessToken:  "",
		RefreshToken: "",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logrus.Fatalf("failed listen to tcp %s", err.Error())
	}

	s := grpc.NewServer()
	v1.RegisterAuthServer(s, &server{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
}
