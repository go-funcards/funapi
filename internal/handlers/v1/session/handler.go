package session

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-funcards/funapi/internal/gin/binding"
	"github.com/go-funcards/funapi/internal/gin/httputil"
	"github.com/go-funcards/funapi/internal/handlers"
	v1Authz "github.com/go-funcards/funapi/proto/authz_service/v1"
	v1User "github.com/go-funcards/funapi/proto/user_service/v1"
	"github.com/go-funcards/jwt"
	"github.com/go-funcards/token"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

var _ handlers.Handler = (*Handler)(nil)

type Handler struct {
	UserService    v1User.UserClient
	SubjectService v1Authz.SubjectClient
	TokenService   token.Service
	Log            *zap.Logger
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	g := rg.Group("/session")
	{
		g.POST("/create", h.create)
		g.POST("/credentials", h.credentials)
		g.POST("/refresh", h.refreshToken)
	}
}

// @Summary Session By Creating User
// @Tags Session
// @Description Return session of created user
// @ModuleID createSession
// @Accept json
// @Produce json
// @Param payload body session.CreateUserDTO true "User data"
// @Success 200 {object} token.Session
// @Failure 400,409,422,500 {object} httputil.APIError
// @Router /session/create [post]
func (h *Handler) create(c *gin.Context) {
	h.Log.Debug("session handler::create bind")
	var dto CreateUserDTO
	if !binding.BindBodyAndValidate(c, &dto) {
		return
	}

	ctx := context.TODO()
	req := dto.toCreate(uuid.NewString(), "ROLE_USER")

	if _, err := h.SubjectService.SaveSub(ctx, &v1Authz.SaveSubRequest{
		SubId: req.GetUserId(),
		Roles: req.GetRoles(),
	}); err != nil {
		_ = c.Error(err)
		return
	}

	if _, err := h.UserService.CreateUser(ctx, req); err != nil {
		_ = c.Error(err)
		return
	}

	session, err := h.TokenService.SessByUser(ctx, jwt.User{
		UserID: req.GetUserId(),
		Name:   req.GetName(),
		Email:  req.GetEmail(),
	})
	if err == nil {
		h.session(c, session)
	} else {
		_ = c.Error(err)
	}
}

// @Summary Session By Credentials
// @Tags Session
// @Description Return session by credentials
// @ModuleID credentials
// @Accept json
// @Produce json
// @Param payload body session.CredentialsDTO true "Credentials"
// @Success 200 {object} token.Session
// @Failure 400,404,422,500 {object} httputil.APIError
// @Router /session/credentials [post]
func (h *Handler) credentials(c *gin.Context) {
	h.Log.Debug("session handler::credentials bind CredentialsDTO")
	var dto CredentialsDTO
	if !binding.BindBodyAndValidate(c, &dto) {
		return
	}

	ctx := context.TODO()

	user, err := h.UserService.GetUserByEmailAndPassword(ctx, dto.toRequest())
	if err != nil {
		_ = c.Error(err)
		return
	}

	session, err := h.TokenService.SessByUser(ctx, jwt.User{
		UserID: user.GetUserId(),
		Name:   user.GetName(),
		Email:  user.GetEmail(),
	})
	if err == nil {
		h.session(c, session)
	} else {
		_ = c.Error(err)
	}
}

// @Summary Session By Refresh Token
// @Tags Session
// @Description Return session by refresh token
// @ModuleID refreshToken
// @Accept json
// @Produce json
// @Param payload body session.RefreshTokenDTO true "Refresh Token"
// @Success 200 {object} token.Session
// @Failure 400,404,422,500 {object} httputil.APIError
// @Router /session/refresh [post]
func (h *Handler) refreshToken(c *gin.Context) {
	h.Log.Debug("session handler::refreshToken bind RefreshTokenDTO")
	var dto RefreshTokenDTO
	if !binding.BindBodyAndValidate(c, &dto) {
		return
	}

	session, err := h.TokenService.SessByRefreshToken(context.TODO(), dto.RefreshToken)
	if err == nil {
		h.session(c, session)
	} else {
		_ = c.Error(err)
	}
}

func (h *Handler) session(c *gin.Context, session token.Session) {
	httputil.NoCache(c)
	c.JSON(http.StatusOK, session)
}
