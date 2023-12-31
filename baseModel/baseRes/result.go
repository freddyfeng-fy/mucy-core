package baseRes

import (
	"github.com/freddyfeng-fy/mucy-core/baseModel/baseError"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 响应结构体
type Response struct {
	Success   bool        `json:"success"`
	ErrorCode int         `json:"errorCode"` // 自定义错误码
	Data      interface{} `json:"data"`      // 数据
	Message   string      `json:"message"`   // 信息
}

// Success 响应成功 ErrorCode 为 0 表示成功
func SuccessFul(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		true,
		0,
		nil,
		"ok",
	})
}

// Success 响应成功 ErrorCode 为 0 表示成功
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		true,
		0,
		data,
		"ok",
	})
}

// Fail 响应失败 ErrorCode 不为 0 表示失败
func Fail(c *gin.Context, errorCode int, msg string) {
	c.JSON(http.StatusOK, Response{
		false,
		errorCode,
		nil,
		msg,
	})
}

// FailByError 失败响应 返回自定义错误的错误码、错误信息
func FailByError(c *gin.Context, error baseError.CustomError) {
	Fail(c, error.ErrorCode, error.ErrorMsg)
}

// ValidateFail 请求参数验证失败
func ValidateFail(c *gin.Context, msg string) {
	Fail(c, baseError.Errors.ValidateError.ErrorCode, msg)
}

// BusinessFail 业务逻辑失败
func BusinessFail(c *gin.Context, msg string) {
	Fail(c, baseError.Errors.BusinessError.ErrorCode, msg)
}

func TokenFail(c *gin.Context) {
	Fail(c, baseError.Errors.TokenError.ErrorCode, "The login authorization is invalid")
}

func AuthFail(c *gin.Context) {
	Fail(c, baseError.Errors.TokenError.ErrorCode, "No permissions")
}

func ServerError(c *gin.Context, err interface{}) {
	msg := "Internal Server Error"
	c.JSON(http.StatusInternalServerError, Response{
		false,
		http.StatusInternalServerError,
		nil,
		msg,
	})
	c.Abort()
}
