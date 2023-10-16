package baseReq

type PageReq struct {
	Page  int `form:"page" json:"page"`
	Limit int `form:"limit" json:"limit"`
}
