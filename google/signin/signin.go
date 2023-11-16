package signin

import (
	"context"
	googleConf "github.com/freddyfeng-fy/mucy-core/google"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
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
		Scopes: []string{people.UserinfoEmailScope, people.UserinfoProfileScope, people.UserEmailsReadScope,
			people.ContactsOtherReadonlyScope, people.ContactsReadonlyScope, people.DirectoryReadonlyScope,
			people.UserAddressesReadScope, people.UserBirthdayReadScope, people.UserGenderReadScope,
			people.UserOrganizationReadScope, people.UserPhonenumbersReadScope},
		Endpoint: google.Endpoint,
	}
}

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
	userInfo, err = peopleService.People.Get("people/me").
		PersonFields("addresses,ageRanges,biographies,birthdays,calendarUrls,clientData,coverPhotos," +
			"emailAddresses,events,externalIds,genders,imClients,interests,locales,locations,memberships,metadata," +
			"miscKeywords,names,nicknames,occupations,organizations,phoneNumbers,photos,relations,sipAddresses," +
			"skills,urls,userDefined").
		Do()
	if err != nil {
		return
	}
	return
}
