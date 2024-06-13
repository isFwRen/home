package request

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

// Logins User login structure
type Logins struct {
	Username   string `json:"username" example:"工号(内网传, 外网不用传)"`
	Phone      string `json:"phone" example:"手机号(内网不用传, 外网传)"`
	Password   string `json:"password" example:"密码"`
	IsIntranet bool   `bson:"isIntranet"` //判断内外网登录, 内网传true, 外网传false
	Captcha    string `json:"captcha" example:"验证码-对应获取验证码接口返回来的json数据里面的Captcha字段"`
	CaptchaId  string `json:"captchaId" example:"验证码ID-对应获取验证码接口返回来的json数据里面的CaptchaId字段"`
}

// SetUserAuth Modify  user's auth structure
type SetUserAuth struct {
	UUID        uuid.UUID `json:"uuid"`
	AuthorityId string    `json:"authorityId"`
}

type QueryUserId struct {
	IDCard string `json:"IDCard"`
}

type ResignationStruct struct {
	Username string `json:"username" example:"工号"`
	Nickname string `json:"nickname" example:"姓名"`
	Reason   string `json:"reason" example:"离职理由"`
}

type ChangePasswordStruct struct {
	Username    string `json:"username" example:"工号(内网传, 外网不用传)"`
	Phone       string `json:"phone" example:"手机号(内网不用传, 外网传)"`
	IsIntranet  bool   `json:"isIntranet"` //判断内外网登录, 内网传true, 外网传false
	Captcha     string `json:"captcha" example:"验证码-对应获取验证码接口返回来的json数据里面的Captcha字段"`
	CaptchaId   string `json:"captchaId" example:"验证码ID-对应获取验证码接口返回来的json数据里面的CaptchaId字段"`
	Password    string `json:"oldpass" example:"旧密码"`
	NewPassword string `json:"newPass" example:"新密码"`
}

type ResetPassword struct {
	Username   string `json:"username" example:"工号"`
	Phone      string `json:"phone" example:"手机号(内网不用传, 外网传)"`
	IsIntranet bool   `json:"isIntranet"` //判断内外网登录, 内网传true, 外网传false
	Captcha    string `json:"captcha" example:"验证码-对应获取验证码接口返回来的json数据里面的Captcha字段"`
	CaptchaId  string `json:"captchaId" example:"验证码ID-对应获取验证码接口返回来的json数据里面的CaptchaId字段"`
}

type ChangeRoleForm struct {
	Id        string    `json:"id" form:"id" binding:"required"`
	RoleId    string    `json:"roleId" form:"roleId" binding:"required"`
	UpdatedAt time.Time `json:"updatedAt" form:"updatedAt" binding:"required"`
}
