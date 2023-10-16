package i18n

import (
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type (
	// Bundle i18n.Bundle
	Bundle = i18n.Bundle

	// Message i18n.Message
	Message = i18n.Message

	// LocalizeConfig i18n.LocalizeConfig
	LocalizeConfig = i18n.LocalizeConfig

	// Data TemplateData
	Data = map[string]interface{}
)

const (
	GinI18nKey       = "i18n"
	GinI18nAcceptKey = "i18n-accept"
)

// NewBundle new bundle
func NewBundle(tag language.Tag, tomls ...string) *Bundle {
	bundle := i18n.NewBundle(tag)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	for _, file := range tomls {
		bundle.LoadMessageFile(file)
	}
	return bundle
}

// MustFormat must
func MustFormat(c *gin.Context, lc *i18n.LocalizeConfig) string {
	return MustLocalizer(c).MustLocalize(lc)
}

func Text(c *gin.Context, messageId string) string {
	return TextArgs(c, messageId, nil)
}

func TextArgs(c *gin.Context, messageId string, args map[string]interface{}) string {
	message := &i18n.Message{
		ID: messageId,
	}
	if localizer, ok := GetLocalizer(c); ok {
		return localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: message,
			TemplateData:   args,
		})
	}
	return formatInternalMessage(message, args)
}

// MustLocalizer i18n
func MustLocalizer(c *gin.Context) *i18n.Localizer {
	localizer, ok := GetLocalizer(c)
	if !ok {
		panic(errors.New("context no has i18n localizer"))
	}
	return localizer
}

// GetLocalizer i18n
func GetLocalizer(c *gin.Context) (*i18n.Localizer, bool) {
	if v, ok := c.Get(GinI18nKey); ok {
		if l, b := v.(*i18n.Localizer); b {
			return l, true
		}
	}
	return nil, false
}

func formatInternalMessage(message *i18n.Message, args map[string]interface{}) string {
	if args == nil {
		return message.Other
	}
	tpl := i18n.NewMessageTemplate(message)
	msg, err := tpl.Execute("other", args, nil)
	if err != nil {
		panic(err)
	}
	return msg
}

func GetLocalizerAccept(c *gin.Context) string {
	return c.Keys[GinI18nAcceptKey].(string)
}
