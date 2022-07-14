package binding

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-funcards/validate"
)

func BindAllAndValidate(c *gin.Context, obj any) bool {
	return BindCtx(c, obj) && BindUri(c, obj) && BindQuery(c, obj) &&
		BindHeader(c, obj) && BindBody(c, obj) && Validate(c, obj)
}

func Bind(c *gin.Context, obj any) bool {
	return withError(c, c.ShouldBind(obj))
}

func BindUri(c *gin.Context, obj any) bool {
	return withError(c, c.ShouldBindUri(obj))
}

func BindUriAndValidate(c *gin.Context, obj any) bool {
	return BindUri(c, obj) && Validate(c, obj)
}

func BindQuery(c *gin.Context, obj any) bool {
	return withError(c, c.ShouldBindQuery(obj))
}

func BindQueryAndValidate(c *gin.Context, obj any) bool {
	return BindQuery(c, obj) && Validate(c, obj)
}

func BindHeader(c *gin.Context, obj any) bool {
	return withError(c, c.ShouldBindHeader(obj))
}

func BindBody(c *gin.Context, obj any) bool {
	return withError(c, c.ShouldBindBodyWith(obj, BindingBody(c)))
}

func BindBodyAndValidate(c *gin.Context, obj any) bool {
	return BindBody(c, obj) && Validate(c, obj)
}

func BindingBody(c *gin.Context) binding.BindingBody {
	bind, ok := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)
	if !ok {
		bind = binding.JSON
	}
	return bind
}

// BindCtx works only for strings
func BindCtx(c *gin.Context, obj any) bool {
	m := make(map[string][]string)
	for k, v := range c.Keys {
		if str, ok := v.(string); ok {
			m[k] = []string{str}
		}
	}
	return withError(c, bindCtx(m, obj))
}

func Validate(c *gin.Context, obj any) bool {
	if err := validate.Default.ValidateStruct(obj); err != nil {
		_ = c.Error(err)
		return false
	}
	return true
}

func bindCtx(m map[string][]string, obj any) error {
	if err := binding.MapFormWithTag(obj, m, "ctx"); err != nil {
		return err
	}
	if binding.Validator == nil {
		return nil
	}
	return binding.Validator.ValidateStruct(obj)
}

func withError(c *gin.Context, err error) bool {
	if err != nil {
		_ = c.Error(err)
		return false
	}
	return true
}
