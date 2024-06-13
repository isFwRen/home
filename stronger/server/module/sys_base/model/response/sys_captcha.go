package response

type SysCaptchaResponse struct {
	CaptchaId string `json:"captchaId"`
	Captcha string `json:"captcha"`
	PicPath   string `json:"picPath"`
}
