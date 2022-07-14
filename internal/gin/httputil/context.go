package httputil

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-funcards/funapi/internal/client/authz_service/v1"
	proto "github.com/go-funcards/funapi/proto/authz_service/v1"
	"github.com/go-funcards/jwt"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strings"
)

const (
	User        = "user"
	UserID      = "user_id"
	AccessToken = "access_token"
)

type IsGrantedFn func(ctx context.Context, c *gin.Context, owner, ref, act string) error

func IsGranted(client proto.AuthorizationCheckerClient, name string) func(ctx context.Context, c *gin.Context, owner, ref, act string) error {
	return func(ctx context.Context, c *gin.Context, owner, ref, act string) error {
		if v1.IsGranted(ctx, client, GetUserID(c), v1.NewObject(name, owner, ref), act) {
			return nil
		}
		return ErrForbidden
	}
}

func SetUser(c *gin.Context, user jwt.User) {
	c.Set(User, user)
	c.Set(UserID, user.UserID)
}

func GetUser(c *gin.Context) jwt.User {
	if u, ok := c.Get(User); ok {
		return u.(jwt.User)
	}
	return jwt.User{}
}

func GetUserID(c *gin.Context) string {
	return c.GetString(UserID)
}

func SetAccessToken(c *gin.Context, accessToken string) {
	c.Set(AccessToken, accessToken)
}

func GetAccessToken(c *gin.Context) string {
	return c.GetString(AccessToken)
}

func GRPCAuthContext(ctx context.Context, c *gin.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", GetAccessToken(c))
}

func GRPCAuth(c *gin.Context) context.Context {
	return GRPCAuthContext(context.TODO(), c)
}

func CtxWithUserID(c *gin.Context) (context.Context, string) {
	return GRPCAuth(c), GetUserID(c)
}

func Created(c *gin.Context, id string) {
	c.Header("Location", fmt.Sprintf("%s/%s", strings.TrimRight(c.Request.URL.Path, "/"), id))
	c.Status(http.StatusCreated)
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func NoCache(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-store")
	c.Writer.Header().Set("Pragma", "no-cache")
}
