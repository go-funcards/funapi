package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-funcards/funapi/internal/gin/httputil"
	"github.com/go-funcards/jwt"
	"net/http"
	"strings"
)

// ErrInvalidAuthHeader indicates authorization header is invalid, could for example have the wrong Realm name
var ErrInvalidAuthHeader = errors.New("authorization header is invalid")

var (
	messages = map[int]*httputil.APIError{
		http.StatusBadRequest: {
			StatusCode: http.StatusBadRequest,
			ErrorCode:  "invalid_authorization",
			Message:    "Authorization header is invalid",
		},
		http.StatusUnauthorized: {
			StatusCode: http.StatusUnauthorized,
			ErrorCode:  "invalid_token",
			Message:    "JWT token is invalid or expired",
		},
	}
)

func Authorize(verifier jwt.Verifier, authScheme string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.Split(c.Request.Header.Get("Authorization"), " ")
		if len(authHeader) != 2 || strings.ToLower(authHeader[0]) != strings.ToLower(authScheme) {
			_ = c.Error(ErrInvalidAuthHeader)
			c.AbortWithStatusJSON(http.StatusBadRequest, c.Error(messages[http.StatusBadRequest]).Err)
			return
		}

		user, err := verifier.ExtractUser(authHeader[1])
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, c.Error(messages[http.StatusUnauthorized]).Err)
			return
		}

		httputil.SetAccessToken(c, authHeader[1])
		httputil.SetUser(c, user)

		c.Next()
	}
}
