package users

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-funcards/funapi/internal/gin/binding"
	"github.com/go-funcards/funapi/internal/gin/httputil"
	"github.com/go-funcards/funapi/internal/handlers"
	"github.com/go-funcards/funapi/internal/handlers/v1/clientutil"
	"github.com/go-funcards/funapi/proto/user_service/v1"
	"go.uber.org/zap"
	"net/http"
)

var _ handlers.Handler = (*Handler)(nil)

type Handler struct {
	UserService v1.UserClient
	IsGranted   httputil.IsGrantedFn
	Log         *zap.Logger
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	g := rg.Group("/users")
	{
		g.GET("/me", h.me)
		g.PATCH("/:user_id", h.update)
	}
}

// @Summary Update User
// @Tags Users
// @Description Update User Info
// @ModuleID updateUser
// @Accept json
// @Produce json
// @Param user_id path string true "user id" format(uuid)
// @Param payload body users.UpdateUserDTO true "User data"
// @Success 204
// @Failure 400,401,403,404,409,422,500 {object} httputil.APIError
// @Router /users/{user_id} [patch]
// @Security BearerAuth
func (h *Handler) update(c *gin.Context) {
	ctx := context.TODO()

	h.Log.Debug("user handler::create bind")
	var dto UpdateUserDTO
	if !binding.BindUri(c, &dto) || !binding.BindBodyAndValidate(c, &dto) {
		return
	}

	if err := h.IsGranted(ctx, c, dto.UserID, "", "UPDATE"); err != nil {
		_ = c.Error(err)
		return
	}

	if _, err := h.UserService.UpdateUser(ctx, dto.toUpdate()); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Get Authenticated User
// @Tags Users
// @Description Return authenticated user info
// @ModuleID me
// @Produce json
// @Success 200 {object} users.User
// @Failure 400,401,404,500 {object} httputil.APIError
// @Router /users/me [get]
// @Security BearerAuth
func (h *Handler) me(c *gin.Context) {
	ctx := context.TODO()

	h.Log.Debug("user handler::me")
	user, err := clientutil.GetUser(ctx, h.UserService, httputil.GetUserID(c))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, CreateUser(user))
}
