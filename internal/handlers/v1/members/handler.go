package members

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-funcards/funapi/internal/gin/binding"
	"github.com/go-funcards/funapi/internal/gin/httputil"
	"github.com/go-funcards/funapi/internal/handlers"
	v1Authz "github.com/go-funcards/funapi/proto/authz_service/v1"
	"go.uber.org/zap"
)

var _ handlers.Handler = (*Handler)(nil)

type Handler struct {
	*handlers.BaseBoard
	SubjectService v1Authz.SubjectClient
	Log            *zap.Logger
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	g := rg.Group("/members")
	{
		m := g.Group("/:member_id")
		{
			m.PUT("", h.save)
			m.DELETE("", h.delete)
		}
	}
}

// @Summary Save Board Member
// @Tags Boards
// @ModuleID saveBoardMember
// @Accept json
// @Param board_id path string true "Board ID" format(uuid)
// @Param member_id path string true "Member ID" format(uuid)
// @Param payload body members.SaveMemberDTO true "Member data"
// @Success 204
// @Failure 400,401,403,404,422,500 {object} httputil.APIError
// @Router /boards/{board_id}/members/{member_id} [put]
// @Security BearerAuth
func (h *Handler) save(c *gin.Context) {
	ctx := context.TODO()

	h.Log.Debug("member handler::add bind")
	var dto SaveMemberDTO
	if !binding.BindUri(c, &dto) || !binding.BindBodyAndValidate(c, &dto) {
		return
	}

	if !h.IsGranted(ctx, c, dto.BoardID, "SAVE_MEMBER") {
		return
	}

	h.Log.Debug("member handler::add call gRPC /SubjectClient/SaveSub")
	if _, err := h.SubjectService.SaveSub(ctx, dto.toSaveSub()); err != nil {
		_ = c.Error(err)
		return
	}

	h.Log.Debug("member handler::add call gRPC /BoardClient/UpdateBoard")
	if _, err := h.BoardService.UpdateBoard(ctx, dto.toUpdate()); err != nil {
		_ = c.Error(err)
		return
	}

	httputil.NoContent(c)
}

// @Summary Delete Board Member
// @Tags Boards
// @ModuleID deleteBoardMember
// @Param board_id path string true "Board ID" format(uuid)
// @Param member_id path string true "Member ID" format(uuid)
// @Success 204
// @Failure 400,401,403,404,500 {object} httputil.APIError
// @Router /boards/{board_id}/members/{member_id} [delete]
// @Security BearerAuth
func (h *Handler) delete(c *gin.Context) {
	ctx := context.TODO()

	h.Log.Debug("member handler::delete bind")
	var dto DeleteMemberDTO
	if !binding.BindUriAndValidate(c, &dto) {
		return
	}

	if !h.IsGranted(ctx, c, dto.BoardID, "DELETE_MEMBER") {
		return
	}

	h.Log.Debug("member handler::delete call gRPC /BoardClient/UpdateBoard")
	if _, err := h.BoardService.UpdateBoard(ctx, dto.toUpdate()); err != nil {
		_ = c.Error(err)
		return
	}

	h.Log.Debug("member handler::add call gRPC /SubjectClient/SaveSub")
	if _, err := h.SubjectService.SaveSub(ctx, dto.toSaveSub()); err != nil {
		_ = c.Error(err)
		return
	}

	httputil.NoContent(c)
}
