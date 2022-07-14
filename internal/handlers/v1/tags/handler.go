package tags

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-funcards/funapi/internal/gin/binding"
	"github.com/go-funcards/funapi/internal/gin/httputil"
	"github.com/go-funcards/funapi/internal/handlers"
	"github.com/go-funcards/funapi/internal/handlers/v1/clientutil"
	v1Tag "github.com/go-funcards/funapi/proto/tag_service/v1"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

var _ handlers.Handler = (*Handler)(nil)

type Handler struct {
	*handlers.BaseBoard
	TagService v1Tag.TagClient
	Log        *zap.Logger
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	g := rg.Group("/tags")
	{
		g.GET("", h.list)
		g.POST("", h.create)

		b := g.Group("/:tag_id")
		{
			b.GET("", h.read)
			b.PATCH("", h.update)
			b.DELETE("", h.delete)
		}
	}
}

// @Summary Tag List
// @Tags Tags
// @ModuleID listTag
// @Accept json
// @Produce json
// @Param board_id query string true "Board ID" format(uuid)
// @Param page_index query int false "Page Index" minimum(0)
// @Param page_size query int false "Page Size" minimum(1) maximum(1000)
// @Success 200 {object} tags.PageResponse
// @Failure 400,401,403,500 {object} httputil.APIError
// @Router /tags [get]
// @Security BearerAuth
func (h *Handler) list(c *gin.Context) {
	h.Log.Debug("tag handler::list bind")
	req := PageReq()
	if !binding.BindQueryAndValidate(c, &req) {
		return
	}

	ctx := context.TODO()

	if !h.IsGranted(ctx, c, req.BoardID, "READ") {
		return
	}

	h.Log.Debug("tag handler::list call gRPC /TagClient/GetTags")
	response, err := h.TagService.GetTags(ctx, req.toRead())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, PageResp(response, req))
}

// @Summary Create Tag
// @Tags Tags
// @ModuleID createTag
// @Accept json
// @Param payload body tags.CreateTagDTO true "Tag data"
// @Success 201
// @Failure 400,401,403,422,500 {object} httputil.APIError
// @Header 201 {string} Location "/tags/{tag_id}"
// @Router /tags [post]
// @Security BearerAuth
func (h *Handler) create(c *gin.Context) {
	h.Log.Debug("tag handler::create bind")
	var dto CreateTagDTO
	if !binding.BindCtx(c, &dto) || !binding.BindBodyAndValidate(c, &dto) {
		return
	}

	ctx := context.TODO()

	if !h.IsGranted(ctx, c, dto.BoardID, "CREATE") {
		return
	}

	id := uuid.NewString()

	h.Log.Debug("tag handler::create call gRPC /TagClient/CreateTag")
	if _, err := h.TagService.CreateTag(ctx, dto.toCreate(id)); err != nil {
		_ = c.Error(err)
		return
	}

	httputil.Created(c, id)
}

// @Summary Read Tag
// @Tags Tags
// @ModuleID readTag
// @Accept json
// @Produce json
// @Param tag_id path string true "Tag ID" format(uuid)
// @Success 200 {object} tags.Tag
// @Failure 400,401,403,404,500 {object} httputil.APIError
// @Router /tags/{tag_id} [get]
// @Security BearerAuth
func (h *Handler) read(c *gin.Context) {
	h.Log.Debug("tag handler::read bind")
	var dto ReadTagDTO
	if !binding.BindUriAndValidate(c, &dto) {
		return
	}

	ctx := context.TODO()

	h.Log.Debug("tag handler::read call gRPC /TagClient/GetTags")
	tag, err := clientutil.GetTag(ctx, h.TagService, dto.TagID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if !h.IsGranted(ctx, c, tag.GetBoardId(), "READ") {
		return
	}

	c.JSON(http.StatusOK, CreateTag(tag))
}

// @Summary Update Tag
// @Tags Tags
// @ModuleID updateTag
// @Accept json
// @Param tag_id path string true "Tag ID" format(uuid)
// @Param payload body tags.UpdateTagDTO true "Tag data"
// @Success 204
// @Failure 400,401,403,404,422,500 {object} httputil.APIError
// @Router /tags/{tag_id} [patch]
// @Security BearerAuth
func (h *Handler) update(c *gin.Context) {
	h.Log.Debug("tag handler::update bind")
	var dto UpdateTagDTO
	if !binding.BindUri(c, &dto) || !binding.BindBodyAndValidate(c, &dto) {
		return
	}

	ctx := context.TODO()

	h.Log.Debug("tag handler::update call gRPC /TagClient/GetTags")
	tag, err := clientutil.GetTag(ctx, h.TagService, dto.TagID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if !h.IsGranted(ctx, c, tag.GetBoardId(), "UPDATE") {
		return
	}

	if len(dto.BoardID) > 0 && dto.BoardID != tag.GetBoardId() && !h.IsGranted(ctx, c, dto.BoardID, "CREATE") {
		return
	}

	h.Log.Debug("tag handler::update call gRPC /TagClient/UpdateTag")
	if _, err = h.TagService.UpdateTag(ctx, dto.toUpdate()); err != nil {
		_ = c.Error(err)
		return
	}

	httputil.NoContent(c)
}

// @Summary Delete Tag
// @Tags Tags
// @ModuleID deleteTag
// @Param tag_id path string true "Tag ID" format(uuid)
// @Success 204
// @Failure 400,401,403,404,500 {object} httputil.APIError
// @Router /tags/{tag_id} [delete]
// @Security BearerAuth
func (h *Handler) delete(c *gin.Context) {
	h.Log.Debug("tag handler::delete bind")
	var dto DeleteTagDTO
	if !binding.BindUriAndValidate(c, &dto) {
		return
	}

	ctx := context.TODO()

	h.Log.Debug("tag handler::delete call gRPC /TagClient/GetTags")
	tag, err := clientutil.GetTag(ctx, h.TagService, dto.TagID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if !h.IsGranted(ctx, c, tag.GetBoardId(), "DELETE") {
		return
	}

	h.Log.Debug("tag handler::delete call gRPC /TagClient/DeleteTag")
	if _, err = h.TagService.DeleteTag(ctx, dto.toDelete()); err != nil {
		_ = c.Error(err)
		return
	}

	httputil.NoContent(c)
}
