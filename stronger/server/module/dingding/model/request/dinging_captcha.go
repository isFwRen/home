package request

type GetCaptcha struct {
	Phone      string `json:"phone" example:"手机号"`
}
