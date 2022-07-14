package v1

import (
	"context"
	"github.com/go-funcards/funapi/proto/user_service/v1"
	"github.com/go-funcards/grpc-pool"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ v1.UserClient = (*UserService)(nil)

type UserService struct {
	Pool *grpcpool.Pool
}

func (s *UserService) CreateUser(ctx context.Context, in *v1.CreateUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewUserClient(conn).CreateUser(ctx, in, opts...)
}

func (s *UserService) UpdateUser(ctx context.Context, in *v1.UpdateUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewUserClient(conn).UpdateUser(ctx, in, opts...)
}

func (s *UserService) DeleteUser(ctx context.Context, in *v1.DeleteUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewUserClient(conn).DeleteUser(ctx, in, opts...)
}

func (s *UserService) GetUsers(ctx context.Context, in *v1.UsersRequest, opts ...grpc.CallOption) (*v1.UsersResponse, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewUserClient(conn).GetUsers(ctx, in, opts...)
}

func (s *UserService) GetUserByEmailAndPassword(ctx context.Context, in *v1.UserByEmailAndPasswordRequest, opts ...grpc.CallOption) (*v1.UserResponse, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewUserClient(conn).GetUserByEmailAndPassword(ctx, in, opts...)
}
