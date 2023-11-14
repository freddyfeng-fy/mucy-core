package signin

import (
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

var (
	oauth2Config = &oauth2.Config{
		ClientID:     "YOUR_CLIENT_ID",
		ClientSecret: "YOUR_CLIENT_SECRET",
		RedirectURL:  "YOUR_REDIRECT_URL",
		Scopes:       []string{people.UserinfoEmailScope},
		Endpoint:     google.Endpoint,
	}
)

func GoogleSignin() string {
	return oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func GoogleCallback(code string) (err error, userInfo *people.Person) {
	ctx := context.Background()
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
	userInfo, err = peopleService.People.Get("people/me").PersonFields("emailAddresses").Do()
	if err != nil {
		return
	}
	return
}
