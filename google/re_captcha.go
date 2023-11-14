package google

import (
	"github.com/freddyfeng-fy/mucy-core/jsons"
	"io"
	"net/http"
	"net/url"
)

var (
	conf *Conf
)

func InitReCaptchaConfig(config *Conf) {
	conf = config
}

func VerifyRecaptcha(responseToken string) (*RecaptchaResponse, error) {
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{
			"secret":   {conf.ReCaptCha.RecaptchaKey},
			"response": {responseToken},
		})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var recaptchaResp RecaptchaResponse
	err = jsons.Cjson.Unmarshal(body, &recaptchaResp)
	if err != nil {
		return nil, err
	}

	return &recaptchaResp, nil
}
