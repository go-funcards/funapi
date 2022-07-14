package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-funcards/funapi/internal/gin/httputil"
	"github.com/go-funcards/funapi/internal/handlers/v1/clientutil"
	"github.com/go-funcards/funapi/proto/board_service/v1"
)

type Handler interface {
	Register(rg *gin.RouterGroup)
}

type BaseBoard struct {
	BoardService v1.BoardClient
	IsGrantedFn  httputil.IsGrantedFn
}

func (h *BaseBoard) GetBoard(ctx context.Context, boardID string) (*v1.BoardsResponse_Board, error) {
	return clientutil.GetBoard(ctx, h.BoardService, boardID)
}

func (h *BaseBoard) IsGranted(ctx context.Context, c *gin.Context, boardID, act string) bool {
	board, err := h.GetBoard(ctx, boardID)
	if err != nil {
		_ = c.Error(err)
		return false
	}

	if err = h.IsGrantedFn(ctx, c, board.GetOwnerId(), board.GetBoardId(), act); err != nil {
		_ = c.Error(err)
		return false
	}

	return true
}
