package clientutil

import (
	"context"
	"github.com/go-funcards/funapi/internal/gin/httputil"
	v1Board "github.com/go-funcards/funapi/proto/board_service/v1"
)

func BoardsRequest(id string) *v1Board.BoardsRequest {
	return &v1Board.BoardsRequest{
		PageIndex: 0,
		PageSize:  1,
		BoardIds:  []string{id},
	}
}

func GetBoard(ctx context.Context, client v1Board.BoardClient, id string) (*v1Board.BoardsResponse_Board, error) {
	response, err := client.GetBoards(ctx, BoardsRequest(id))
	if err != nil {
		return nil, err
	}
	if len(response.GetBoards()) != 1 {
		return nil, httputil.ErrNotFound
	}
	return response.GetBoards()[0], nil
}
