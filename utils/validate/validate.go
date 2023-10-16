package validate

import (
	"errors"
	"github.com/freddyfeng-fy/mucy-core/i18n"
	"github.com/freddyfeng-fy/mucy-core/utils/strs"
	"github.com/gin-gonic/gin"
	"regexp"
)

// IsUsername 验证用户名合法性，用户名必须由5-12位(数字、字母、_、-)组成，且必须以字母开头。
func IsUsername(username string) error {
	if strs.IsBlank(username) {
		return errors.New("请输入用户名")
	}
	matched, err := regexp.MatchString("^[0-9a-zA-Z_-]{5,12}$", username)
	if err != nil || !matched {
		return errors.New("用户名必须由5-12位(数字、字母、_、-)组成，且必须以字母开头")
	}
	matched, err = regexp.MatchString("^[a-zA-Z]", username)
	if err != nil || !matched {
		return errors.New("用户名必须由5-12位(数字、字母、_、-)组成，且必须以字母开头")
	}
	return nil
}

// IsEmail 验证是否是合法的邮箱
func IsEmail(email string) (err error) {
	if strs.IsBlank(email) {
		err = errors.New("邮箱格式不符合规范")
		return
	}
	pattern := `^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$`
	matched, _ := regexp.MatchString(pattern, email)
	if !matched {
		err = errors.New("邮箱格式不符合规范")
	}
	return
}

// IsPassword 是否是合法的密码
func IsPassword(c *gin.Context, password, rePassword string) error {
	if strs.IsBlank(password) {
		return errors.New(i18n.Text(c, "passwordNotNull"))
	}
	if strs.RuneLen(password) < 6 {
		return errors.New(i18n.Text(c, "passwordTooSimple"))
	}
	if password != rePassword {
		return errors.New(i18n.Text(c, "twoPasswordsEnteredNotMatch"))
	}
	return nil
}

// IsURL 是否是合法的URL
func IsURL(url string) bool {
	if strs.IsBlank(url) {
		return false
	}
	ok, _ := regexp.MatchString(`(?:(?:https?://|[a-z0-9.]?)+(?:(?:[.]))+(?:[a-z]{2,3})+(?=\s|$))`, url)
	if !ok {
		return true
	}
	return false
}
