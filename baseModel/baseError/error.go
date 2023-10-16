package baseError

import (
	"github.com/gin-gonic/gin"
	"mucy-core/i18n"
)

type CustomError struct {
	ErrorCode int
	ErrorMsg  string
}

type CustomErrors struct {
	BusinessError CustomError
	ValidateError CustomError
	TokenError    CustomError
}

var c = &gin.Context{}

var Errors = CustomErrors{
	BusinessError: CustomError{40000, i18n.Text(c, "UnknownException")},
	ValidateError: CustomError{42200, "请求参数错误"},
	TokenError:    CustomError{40100, "登录授权失效"},
}
