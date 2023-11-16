package signin

import (
	"context"
	googleConf "github.com/freddyfeng-fy/mucy-core/google"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
	"net/http"
	"net/url"
)

var (
	oauth2Config *oauth2.Config
	conf         *googleConf.Conf
)

func InitSigninConfig(config *googleConf.Conf) {
	conf = config
	oauth2Config = &oauth2.Config{
		ClientID:     conf.OAuth.ClientID,
		ClientSecret: conf.OAuth.ClientSecret,
		RedirectURL:  conf.OAuth.RedirectURL,
		Scopes:       []string{people.UserinfoEmailScope, people.UserinfoProfileScope},
		Endpoint:     google.Endpoint,
	}
}

func GoogleSignin() string {
	return oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func GoogleCallback(code string) (err error, userInfo *people.Person) {
	rootCtx := context.Background()
	var httpClient *http.Client
	// 如果提供了代理地址，则使用代理
	if conf.OAuth.Proxy != "" {
		proxyURL, err := url.Parse(conf.OAuth.Proxy)
		if err != nil {
			return err, nil
		}
		transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
		httpClient = &http.Client{Transport: transport}
	} else {
		httpClient = http.DefaultClient
	}
	ctx := context.WithValue(rootCtx, oauth2.HTTPClient, httpClient)
	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return
	}
	// 使用token创建一个新的服务
	peopleService, err := people.NewService(ctx, option.WithTokenSource(oauth2Config.TokenSource(ctx, token)))
	if err != nil {
		return
	}
	// 获取用户的信息
	userInfo, err = peopleService.People.Get("people/me").
		PersonFields("names,photos,gender,emailAddresses,metadata").
		Do()
	if err != nil {
		return
	}
	return
}
