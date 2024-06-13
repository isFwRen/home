package request

import (
	"github.com/dgrijalva/jwt-go"
)

//CustomClaims Custom claims structure
type CustomClaims struct {
	//UUID        uuid.UUID
	ID       string
	NickName string
	Code     string
	Phone    string
	DingId   string
	RoleId   string
	RoleName string
	//SysMenu  []model.SysMenu
	jwt.StandardClaims
}
