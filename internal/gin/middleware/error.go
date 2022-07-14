package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-funcards/funapi/internal/gin/httputil"
	"github.com/go-funcards/validate"
	"github.com/go-playground/validator/v10"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

var gRPCCodes = map[codes.Code]int{
	codes.Internal:         http.StatusInternalServerError,
	codes.NotFound:         http.StatusNotFound,
	codes.AlreadyExists:    http.StatusConflict,
	codes.Unauthenticated:  http.StatusUnauthorized,
	codes.PermissionDenied: http.StatusForbidden,
}

func APIError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() || len(c.Errors) == 0 || c.IsAborted() {
			return
		}

		for i := len(c.Errors) - 1; i >= 0; i-- {
			err := c.Errors[i].Err
			for {
				if abort(c, err) {
					return
				}
				if err = errors.Unwrap(err); err == nil {
					break
				}
			}
		}

		sc := http.StatusBadRequest
		c.AbortWithStatusJSON(sc, httputil.Errors[sc])
	}
}

func abort(c *gin.Context, err error) bool {
	switch tmp := err.(type) {
	case validate.SliceValidateError:
		if len(tmp) > 0 {
			if _, ok := tmp[0].(validator.ValidationErrors); ok {
				return abort(c, tmp[0])
			}
		}
	case validator.ValidationErrors:
		sc := http.StatusUnprocessableEntity
		c.AbortWithStatusJSON(sc, httputil.Errors[sc].SetErrors(validationErrors(tmp)))
		return true
	case *httputil.APIError:
		c.AbortWithStatusJSON(tmp.StatusCode, err)
		return true
	}

	if s, ok := status.FromError(err); ok {
		sc := http.StatusInternalServerError
		if found, ok := gRPCCodes[s.Code()]; ok {
			sc = found
		} else if s.Code() == codes.InvalidArgument {
			for _, e := range s.Details() {
				if br, ok := e.(*errdetails.BadRequest); ok {
					data := gin.H{}
					for _, fv := range br.FieldViolations {
						if f, ok := data[fv.GetField()]; ok {
							data[fv.GetField()] = append(f.([]string), fv.GetDescription())
						} else {
							data[fv.GetField()] = []string{fv.GetDescription()}
						}
					}
					sc = http.StatusBadRequest
					c.AbortWithStatusJSON(sc, httputil.Errors[sc].SetErrors(data))
					return true
				}
			}
		}
		c.AbortWithStatusJSON(sc, httputil.Errors[sc])
		return true
	}
	return false
}

func validationErrors(vErrors validator.ValidationErrors) gin.H {
	data := gin.H{}
	for _, ve := range vErrors {
		if i, ok := data[ve.Field()]; ok {
			field := i.(gin.H)
			field[ve.Tag()] = ve.Error()
		} else {
			data[ve.Field()] = gin.H{
				ve.Tag(): ve.Error(),
			}
		}
	}
	return data
}

//var r = regexp.MustCompile(`^.+'(.+)' tag$`)
//var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
//var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
//
//func toSnakeCase(str string) string {
//	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
//	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
//	return strings.ToLower(snake)
//}
