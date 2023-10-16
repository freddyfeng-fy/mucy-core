package utils

import (
	"github.com/asaskevich/govalidator"
	"github.com/freddyfeng-fy/mucy-core/utils/strs"
	"github.com/go-playground/validator/v10"
	"regexp"
)

// ValidateMobile 校验手机号
func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	ok, _ := regexp.MatchString(`^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$`, mobile)
	if !ok {
		return false
	}
	return true
}

func ValidateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	return govalidator.IsEmail(email)
}

func ValidateURL(fl validator.FieldLevel) bool {
	url := fl.Field().String()
	return govalidator.IsURL(url)
}

func ValidateTitleLength(fl validator.FieldLevel) bool {
	title := fl.Field().String()
	if strs.RuneLen(title) > 128 {
		return false
	}
	return true
}
