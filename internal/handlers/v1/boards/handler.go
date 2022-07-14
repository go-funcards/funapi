package boards

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-funcards/funapi/internal/gin/binding"
	"github.com/go-funcards/funapi/internal/gin/httputil"
	"github.com/go-funcards/funapi/internal/handlers"
	v1Authz "github.com/go-funcards/funapi/proto/authz_service/v1"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

var _ handlers.Handler = (*Handler)(nil)

type Handler struct {
	*handlers.BaseBoard
	MemberHandler  handlers.Handler
	SubjectService v1Authz.SubjectClient
	Log            *zap.Logger
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	g := rg.Group("/boards")
	{
		g.GET("", h.list)
		g.POST("", h.create)

		b := g.Group("/:board_id")
		{
			b.GET("", h.read)
			b.PATCH("", h.update)
			b.DELETE("", h.delete)

			h.MemberHandler.Register(b)
		}
	}
}

// @Summary Board List
// @Tags Boards
// @ModuleID listBoard
// @Accept json
// @Produce json
// @Param page_index query int false "Page Index" minimum(0)
// @Param page_size query int false "Page Size" minimum(1) maximum(1000)
// @Success 200 {object} boards.PageResponse
// @Failure 400,401,403,500 {object} httputil.APIError
// @Router /boards [get]
// @Security BearerAuth
func (h *Handler) list(c *gin.Context) {
	ctx := context.TODO()

	h.Log.Debug("board handler::list bind")
	req := PageReq()
	if !binding.BindCtx(c, &req) || !binding.BindQueryAndValidate(c, &req) {
		return
	}

	h.Log.Debug("board handler::list call gRPC /BoardClient/GetBoards")
	response, err := h.BoardService.GetBoards(ctx, req.toRead())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, PageResp(response, req))
}

// @Summary Create Board
// @Tags Boards
// @ModuleID createBoard
// @Accept json
// @Param payload body boards.CreateBoardDTO true "Board data"
// @Success 201
// @Failure 400,401,403,422,500 {object} httputil.APIError
// @Header 201 {string} Location "/boards/{board_id}"
// @Router /boards [post]
// @Security BearerAuth
func (h *Handler) create(c *gin.Context) {
	ctx := context.TODO()

	if err := h.IsGrantedFn(ctx, c, "", "", "CREATE"); err != nil {
		_ = c.Error(err)
		return
	}

	h.Log.Debug("board handler::create bind")
	var dto CreateBoardDTO
	if !binding.BindCtx(c, &dto) || !binding.BindBodyAndValidate(c, &dto) {
		return
	}

	id := uuid.NewString()

	h.Log.Debug("board handler::create call gRPC /BoardClient/CreateBoard")
	if _, err := h.BoardService.CreateBoard(ctx, dto.toCreate(id)); err != nil {
		_ = c.Error(err)
		return
	}

	httputil.Created(c, id)
}

// @Summary Read Board
// @Tags Boards
// @ModuleID readBoard
// @Accept json
// @Produce json
// @Param board_id path string true "Board ID" format(uuid)
// @Success 200 {object} boards.Board
// @Failure 400,401,403,404,500 {object} httputil.APIError
// @Router /boards/{board_id} [get]
// @Security BearerAuth
func (h *Handler) read(c *gin.Context) {
	ctx := context.TODO()

	h.Log.Debug("board handler::read bind")
	var dto ReadBoardDTO
	if !binding.BindUriAndValidate(c, &dto) {
		return
	}

	h.Log.Debug("board handler::read call gRPC /BoardClient/GetBoards")
	board, err := h.GetBoard(ctx, dto.BoardID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if err = h.IsGrantedFn(ctx, c, board.GetOwnerId(), board.GetBoardId(), "READ"); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, CreateBoard(board))
}

// @Summary Update Board
// @Tags Boards
// @ModuleID updateBoard
// @Accept json
// @Param board_id path string true "Board ID" format(uuid)
// @Param payload body boards.UpdateBoardDTO true "Board data"
// @Success 204
// @Failure 400,401,403,404,422,500 {object} httputil.APIError
// @Router /boards/{board_id} [patch]
// @Security BearerAuth
func (h *Handler) update(c *gin.Context) {
	ctx := context.TODO()

	h.Log.Debug("board handler::update bind")
	var dto UpdateBoardDTO
	if !binding.BindUri(c, &dto) || !binding.BindBodyAndValidate(c, &dto) {
		return
	}

	if !h.IsGranted(ctx, c, dto.BoardID, "UPDATE") {
		return
	}

	h.Log.Debug("board handler::update call gRPC /BoardClient/UpdateBoard")
	if _, err := h.BoardService.UpdateBoard(ctx, dto.toUpdate()); err != nil {
		_ = c.Error(err)
		return
	}

	httputil.NoContent(c)
}

// @Summary Delete Board
// @Tags Boards
// @ModuleID deleteBoard
// @Param board_id path string true "Board ID" format(uuid)
// @Success 204
// @Failure 400,401,403,404,500 {object} httputil.APIError
// @Router /boards/{board_id} [delete]
// @Security BearerAuth
func (h *Handler) delete(c *gin.Context) {
	ctx := context.TODO()

	h.Log.Debug("board handler::delete bind")
	var dto DeleteBoardDTO
	if !binding.BindUriAndValidate(c, &dto) {
		return
	}

	if !h.IsGranted(ctx, c, dto.BoardID, "DELETE") {
		return
	}

	h.Log.Debug("board handler::delete call gRPC /BoardClient/DeleteBoard")
	if _, err := h.BoardService.DeleteBoard(ctx, dto.toDelete()); err != nil {
		_ = c.Error(err)
		return
	}

	h.Log.Debug("board handler::delete call gRPC /SubjectClient/DeleteRef")
	if _, err := h.SubjectService.DeleteRef(ctx, dto.toDeleteRef()); err != nil {
		_ = c.Error(err)
		return
	}

	httputil.NoContent(c)
}
