package v1

import (
	"context"
	"github.com/go-funcards/funapi/proto/tag_service/v1"
	"github.com/go-funcards/grpc-pool"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ v1.TagClient = (*TagService)(nil)

type TagService struct {
	Pool *grpcpool.Pool
}

func (s *TagService) CreateTag(ctx context.Context, in *v1.CreateTagRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewTagClient(conn).CreateTag(ctx, in, opts...)
}

func (s *TagService) UpdateTag(ctx context.Context, in *v1.UpdateTagRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewTagClient(conn).UpdateTag(ctx, in, opts...)
}

func (s *TagService) DeleteTag(ctx context.Context, in *v1.DeleteTagRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewTagClient(conn).DeleteTag(ctx, in, opts...)
}

func (s *TagService) GetTags(ctx context.Context, in *v1.TagsRequest, opts ...grpc.CallOption) (*v1.TagsResponse, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewTagClient(conn).GetTags(ctx, in, opts...)
}
