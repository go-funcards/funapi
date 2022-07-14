package v1

import (
	"context"
	"github.com/go-funcards/funapi/proto/card_service/v1"
	"github.com/go-funcards/grpc-pool"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ v1.CardClient = (*CardService)(nil)

type CardService struct {
	Pool *grpcpool.Pool
}

func (s *CardService) CreateCard(ctx context.Context, in *v1.CreateCardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewCardClient(conn).CreateCard(ctx, in, opts...)
}

func (s *CardService) UpdateCard(ctx context.Context, in *v1.UpdateCardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewCardClient(conn).UpdateCard(ctx, in, opts...)
}

func (s *CardService) UpdateManyCards(ctx context.Context, in *v1.UpdateManyCardsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewCardClient(conn).UpdateManyCards(ctx, in, opts...)
}

func (s *CardService) DeleteCard(ctx context.Context, in *v1.DeleteCardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewCardClient(conn).DeleteCard(ctx, in, opts...)
}

func (s *CardService) GetCards(ctx context.Context, in *v1.CardsRequest, opts ...grpc.CallOption) (*v1.CardsResponse, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewCardClient(conn).GetCards(ctx, in, opts...)
}
