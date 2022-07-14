package categories

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-funcards/funapi/internal/gin/binding"
	"github.com/go-funcards/funapi/internal/gin/httputil"
	"github.com/go-funcards/funapi/internal/handlers"
	"github.com/go-funcards/funapi/internal/handlers/v1/clientutil"
	v1Category "github.com/go-funcards/funapi/proto/category_service/v1"
	"github.com/go-funcards/slice"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

var _ handlers.Handler = (*Handler)(nil)

type Handler struct {
	*handlers.BaseBoard
	CategoryService v1Category.CategoryClient
	Log             *zap.Logger
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	g := rg.Group("/categories")
	{
		g.GET("", h.list)
		g.POST("", h.create)
		g.PATCH("", h.updateMany)

		b := g.Group("/:category_id")
		{
			b.GET("", h.read)
			b.PATCH("", h.update)
			b.DELETE("", h.delete)
		}
	}
}

// @Summary Category List
// @Tags Categories
// @ModuleID listCategory
// @Accept json
// @Produce json
// @Param board_id query string true "Board ID" format(uuid)
// @Param page_index query int false "Page Index" minimum(0)
// @Param page_size query int false "Page Size" minimum(1) maximum(1000)
// @Success 200 {object} categories.PageResponse
// @Failure 400,401,403,500 {object} httputil.APIError
// @Router /categories [get]
// @Security BearerAuth
func (h *Handler) list(c *gin.Context) {
	h.Log.Debug("category handler::list bind")
	req := PageReq()
	if !binding.BindQueryAndValidate(c, &req) {
		return
	}

	ctx := context.TODO()

	if !h.IsGranted(ctx, c, req.BoardID, "READ") {
		return
	}

	h.Log.Debug("category handler::list call gRPC /CategoryClient/GetCategories")
	response, err := h.CategoryService.GetCategories(ctx, req.toRead())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, PageResp(response, req))
}

// @Summary Create Category
// @Tags Categories
// @ModuleID createCategory
// @Accept json
// @Param payload body categories.CreateCategoryDTO true "Category data"
// @Success 201
// @Failure 400,401,403,422,500 {object} httputil.APIError
// @Header 201 {string} Location "/categories/{category_id}"
// @Router /categories [post]
// @Security BearerAuth
func (h *Handler) create(c *gin.Context) {
	h.Log.Debug("category handler::create bind")
	var dto CreateCategoryDTO
	if !binding.BindCtx(c, &dto) || !binding.BindBodyAndValidate(c, &dto) {
		return
	}

	ctx := context.TODO()

	if !h.IsGranted(ctx, c, dto.BoardID, "CREATE") {
		return
	}

	id := uuid.NewString()

	h.Log.Debug("category handler::create call gRPC /CategoryClient/CreateCategory")
	if _, err := h.CategoryService.CreateCategory(ctx, dto.toCreate(id)); err != nil {
		_ = c.Error(err)
		return
	}

	httputil.Created(c, id)
}

// @Summary Update Many Categories
// @Tags Categories
// @ModuleID updateManyCategories
// @Accept json
// @Param payload body categories.UpdateManyCategoriesDTO true "Categories data"
// @Success 204
// @Failure 400,401,403,404,422,500 {object} httputil.APIError
// @Router /categories [patch]
// @Security BearerAuth
func (h *Handler) updateMany(c *gin.Context) {
	h.Log.Debug("category handler::updateMany bind")
	var dto UpdateManyCategoriesDTO
	if !binding.BindBodyAndValidate(c, &dto) {
		return
	}

	ctx := context.TODO()

	h.Log.Debug("category handler::updateMany call gRPC /CategoryClient/GetCategories")
	data, err := h.CategoryService.GetCategories(ctx, &v1Category.CategoriesRequest{
		PageIndex: 0,
		PageSize:  uint32(len(dto.Data)),
		CategoryIds: slice.Map(dto.Data, func(item UpdateCategoryDTO) string {
			return item.CategoryID
		}),
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	boards := make(map[string]bool)

	for _, category := range data.GetCategories() {
		if _, ok := boards[category.GetBoardId()]; !ok {
			boards[category.GetBoardId()] = true
			if !h.IsGranted(ctx, c, category.GetBoardId(), "UPDATE") {
				return
			}
		}
	}

	for _, item := range dto.Data {
		if _, ok := boards[item.BoardID]; !ok && len(item.BoardID) > 0 {
			boards[item.BoardID] = true
			if !h.IsGranted(ctx, c, item.BoardID, "CREATE") {
				return
			}
		}
	}

	h.Log.Debug("category handler::updateMany call gRPC /CategoryClient/UpdateManyCategories")
	if _, err = h.CategoryService.UpdateManyCategories(ctx, dto.toUpdateMany()); err != nil {
		_ = c.Error(err)
		return
	}

	httputil.NoContent(c)
}

// @Summary Read Category
// @Tags Categories
// @ModuleID readCategory
// @Accept json
// @Produce json
// @Param category_id path string true "Category ID" format(uuid)
// @Success 200 {object} categories.Category
// @Failure 400,401,403,404,500 {object} httputil.APIError
// @Router /categories/{category_id} [get]
// @Security BearerAuth
func (h *Handler) read(c *gin.Context) {
	h.Log.Debug("category handler::read bind")
	var dto ReadCategoryDTO
	if !binding.BindUriAndValidate(c, &dto) {
		return
	}

	ctx := context.TODO()

	h.Log.Debug("category handler::read call gRPC /CategoryClient/GetCategories")
	category, err := clientutil.GetCategory(ctx, h.CategoryService, dto.CategoryID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if !h.IsGranted(ctx, c, category.GetBoardId(), "READ") {
		return
	}

	c.JSON(http.StatusOK, CreateCategory(category))
}

// @Summary Update Category
// @Tags Categories
// @ModuleID updateCategory
// @Accept json
// @Param category_id path string true "Category ID" format(uuid)
// @Param payload body categories.UpdateCategoryDTO true "Category data"
// @Success 204
// @Failure 400,401,403,404,422,500 {object} httputil.APIError
// @Router /categories/{category_id} [patch]
// @Security BearerAuth
func (h *Handler) update(c *gin.Context) {
	h.Log.Debug("category handler::update bind")
	var dto UpdateCategoryDTO
	if !binding.BindBody(c, &dto) || !binding.BindUriAndValidate(c, &dto) {
		return
	}

	ctx := context.TODO()

	h.Log.Debug("category handler::update call gRPC /CategoryClient/GetCategories")
	category, err := clientutil.GetCategory(ctx, h.CategoryService, dto.CategoryID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if !h.IsGranted(ctx, c, category.GetBoardId(), "UPDATE") {
		return
	}

	if len(dto.BoardID) > 0 && dto.BoardID != category.GetBoardId() && !h.IsGranted(ctx, c, dto.BoardID, "CREATE") {
		return
	}

	h.Log.Debug("category handler::update call gRPC /CategoryClient/UpdateCategory")
	if _, err = h.CategoryService.UpdateCategory(ctx, dto.toUpdate()); err != nil {
		_ = c.Error(err)
		return
	}

	httputil.NoContent(c)
}

// @Summary Delete Category
// @Tags Categories
// @ModuleID deleteCategory
// @Param category_id path string true "Category ID" format(uuid)
// @Success 204
// @Failure 400,401,403,404,500 {object} httputil.APIError
// @Router /categories/{category_id} [delete]
// @Security BearerAuth
func (h *Handler) delete(c *gin.Context) {
	h.Log.Debug("category handler::delete bind")
	var dto DeleteCategoryDTO
	if !binding.BindUriAndValidate(c, &dto) {
		return
	}

	ctx := context.TODO()

	h.Log.Debug("category handler::delete call gRPC /CategoryClient/GetCategories")
	category, err := clientutil.GetCategory(ctx, h.CategoryService, dto.CategoryID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if !h.IsGranted(ctx, c, category.GetBoardId(), "DELETE") {
		return
	}

	h.Log.Debug("category handler::delete call gRPC /CategoryClient/DeleteCategory")
	if _, err = h.CategoryService.DeleteCategory(ctx, dto.toDelete()); err != nil {
		_ = c.Error(err)
		return
	}

	httputil.NoContent(c)
}
