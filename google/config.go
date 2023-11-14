package google

type Conf struct {
	OAuth     OAuth     `mapstructure:"auth" json:"auth" yaml:"auth"`
	ReCaptCha ReCaptCha `mapstructure:"reCaptCha" json:"reCaptCha" yaml:"reCaptCha"`
}

type OAuth struct {
	ClientID     string `mapstructure:"clientID" json:"clientID" yaml:"clientID"`
	ClientSecret string `mapstructure:"clientSecret" json:"clientSecret" yaml:"clientSecret"`
	RedirectURL  string `mapstructure:"redirectURL" json:"redirectURL" yaml:"redirectURL"`
}

type ReCaptCha struct {
	ProjectID    string `mapstructure:"projectID" json:"projectID" yaml:"projectID"`
	RecaptchaKey string `mapstructure:"recaptchaKey" json:"recaptchaKey" yaml:"recaptchaKey"`
}

type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	ChallengeTs string   `json:"challenge_ts"` // challenge timestamp
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"` // error codes if any
}
