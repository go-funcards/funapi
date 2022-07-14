package v1

import (
	"context"
	"github.com/go-funcards/funapi/proto/category_service/v1"
	"github.com/go-funcards/grpc-pool"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ v1.CategoryClient = (*CategoryService)(nil)

type CategoryService struct {
	Pool *grpcpool.Pool
}

func (s *CategoryService) CreateCategory(ctx context.Context, in *v1.CreateCategoryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewCategoryClient(conn).CreateCategory(ctx, in, opts...)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, in *v1.UpdateCategoryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewCategoryClient(conn).UpdateCategory(ctx, in, opts...)
}

func (s *CategoryService) UpdateManyCategories(ctx context.Context, in *v1.UpdateManyCategoriesRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewCategoryClient(conn).UpdateManyCategories(ctx, in, opts...)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, in *v1.DeleteCategoryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewCategoryClient(conn).DeleteCategory(ctx, in, opts...)
}

func (s *CategoryService) GetCategories(ctx context.Context, in *v1.CategoriesRequest, opts ...grpc.CallOption) (*v1.CategoriesResponse, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewCategoryClient(conn).GetCategories(ctx, in, opts...)
}
