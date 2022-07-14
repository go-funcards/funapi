package v1

import (
	"context"
	"github.com/go-funcards/funapi/proto/authz_service/v1"
	"github.com/go-funcards/grpc-pool"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ v1.SubjectClient = (*SubjectService)(nil)

type SubjectService struct {
	Pool *grpcpool.Pool
}

func (s *SubjectService) SaveSub(ctx context.Context, in *v1.SaveSubRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewSubjectClient(conn).SaveSub(ctx, in, opts...)
}

func (s *SubjectService) DeleteSub(ctx context.Context, in *v1.DeleteSubRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewSubjectClient(conn).DeleteSub(ctx, in, opts...)
}

func (s *SubjectService) DeleteRef(ctx context.Context, in *v1.DeleteRefRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewSubjectClient(conn).DeleteRef(ctx, in, opts...)
}

func (s *SubjectService) GetSub(ctx context.Context, in *v1.SubRequest, opts ...grpc.CallOption) (*v1.SubResponse, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewSubjectClient(conn).GetSub(ctx, in, opts...)
}
