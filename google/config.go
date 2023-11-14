package google

type Conf struct {
	AccessToken string    `mapstructure:"accessToken" json:"accessToken" yaml:"accessToken"`
	ReCaptCha   ReCaptCha `mapstructure:"reCaptCha" json:"reCaptCha" yaml:"reCaptCha"`
}

type ReCaptCha struct {
	ProjectID    string `mapstructure:"projectID" json:"projectID" yaml:"projectID"`
	RecaptchaKey string `mapstructure:"recaptchaKey" json:"recaptchaKey" yaml:"recaptchaKey"`
}
