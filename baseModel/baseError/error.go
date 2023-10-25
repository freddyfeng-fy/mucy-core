package baseError

import (
	"github.com/gin-gonic/gin"
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
	BusinessError: CustomError{40000, "Unknown exception"},
	ValidateError: CustomError{42200, "The request parameter is incorrect"},
	TokenError:    CustomError{40100, "The login authorization is invalid"},
}
