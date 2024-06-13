package response

import (
	uuid "github.com/satori/go.uuid"
	"server/module/sys_base/model"
	sys_base2 "server/module/sys_base/model/request"
)

type SysUserResponse struct {
	User model.SysUser `json:"user"`
}

type LoginResponse struct {
	User      model.SysUser `json:"user"`
	Token     string        `json:"token"`
	Secret    string        `json:"secret"`
	ExpiresAt int64         `json:"expiresAt"`
}

type RegisterAndLoginStruct struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Captcha   string `json:"captcha"`
	CaptchaId string `json:"captchaId"`
}

// Modify  user's auth structure
type SetUserAuth struct {
	UUID        uuid.UUID `json:"uuid"`
	AuthorityId string    `json:"authorityId"`
}

type QueryUserId struct {
	IDCard string `json:"IDCard"`
}

type RegisterResponse struct {
	RegisterUser sys_base2.RegisterStruct `json:"register_user"`
}

type QueryUserResponse struct {
	Username string `json:"Code"`
}

type ResignationResponse struct {
	Resignation sys_base2.ResignationStruct `json:"resignation"`
}
