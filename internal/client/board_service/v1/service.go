package v1

import (
	"context"
	"github.com/go-funcards/funapi/proto/board_service/v1"
	"github.com/go-funcards/grpc-pool"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ v1.BoardClient = (*BoardService)(nil)

type BoardService struct {
	Pool *grpcpool.Pool
}

func (s *BoardService) CreateBoard(ctx context.Context, in *v1.CreateBoardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewBoardClient(conn).CreateBoard(ctx, in, opts...)
}

func (s *BoardService) UpdateBoard(ctx context.Context, in *v1.UpdateBoardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewBoardClient(conn).UpdateBoard(ctx, in, opts...)
}

func (s *BoardService) DeleteBoard(ctx context.Context, in *v1.DeleteBoardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewBoardClient(conn).DeleteBoard(ctx, in, opts...)
}

func (s *BoardService) GetBoards(ctx context.Context, in *v1.BoardsRequest, opts ...grpc.CallOption) (*v1.BoardsResponse, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewBoardClient(conn).GetBoards(ctx, in, opts...)
}
