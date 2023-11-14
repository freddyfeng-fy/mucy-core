package google

type Conf struct {
	AccessToken string    `mapstructure:"accessToken" json:"accessToken" yaml:"accessToken"`
	ReCaptCha   ReCaptCha `mapstructure:"reCaptCha" json:"reCaptCha" yaml:"reCaptCha"`
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
