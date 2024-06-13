package middleware

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"regexp"
	"server/global"
	"server/global/response"
	sys_base2 "server/module/sys_base/model/request"
	"server/module/sys_base/service"
	"server/utils"
	"strings"
	"time"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localSstorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.GetHeader(global.XToken)
		uid := c.GetHeader(global.XUserId)
		code := c.GetHeader(global.XCode)
		//modelToken := model.JwtBlacklist{
		//	Jwt: token,
		//}
		if whitelistsCheck(c.Request.RequestURI) || global.GConfig.System.Env == "develop" {
			c.Next()
			return
		}
		global.GLog.Info("1111", zap.Any("s", "i"))
		if token == "" {
			//response.Result(response.NotLogin, gin.H{
			//	"reload": true,
			//}, global.NoLogin.Error(), c)
			//c.Abort()

			c.AbortWithStatusJSON(response.NotLogin, response.Response{
				Code: response.NotLogin,
				Data: "", Msg: global.NoLogin.Error(),
			})
			return
		}
		global.GLog.Info("2222", zap.Any("s", "i"))
		//还要做一步是否在数据库|Redis是否有
		err1, rToken := service.GetRedisJWT(uid)
		global.GLog.Info("uid", zap.Any(uid, rToken))
		if rToken != token || err1 != nil {
			//response.Result(response.NotLogin, gin.H{
			//	"reload": true,
			//}, TokenLoginOtherExpired.Error(), c)
			global.GLog.Error("2222222222222222", zap.Error(err1))
			c.AbortWithStatusJSON(response.NotLogin, response.Response{
				Code: response.NotLogin,
				Data: err1, Msg: TokenNoLoginExpired.Error(),
			})
			return
		}

		//验证登录otp
		secret, err := service.GetRedisSecret(uid)
		codeNow := utils.GetNowCode(secret)
		global.GLog.Info("code codeCache", zap.Any(code, codeNow))
		if err != nil {
			global.GLog.Error("get cache secret err", zap.Error(err))
			c.AbortWithStatusJSON(response.NotLogin, response.Response{
				Code: response.NotLogin,
				Data: err.Error(), Msg: global.NoLogin.Error(),
			})
			return
		}
		pass, err := utils.ValidateCode(code, secret)
		if err != nil {
			global.GLog.Error("validate code err", zap.Any(code, err))
			c.AbortWithStatusJSON(response.NotLogin, response.Response{
				Code: response.NotLogin,
				Data: err.Error(), Msg: global.NoLogin.Error(),
			})
			return
		}
		if !pass {
			global.GLog.Error("pass", zap.Error(err))
			c.AbortWithStatusJSON(response.NotLogin, response.Response{
				Code: response.NotLogin,
				Data: "请校准系统时间后重新登录", Msg: global.NoLogin.Error(),
			})
			return
		}
		//if service.IsBlacklist(token, modelToken) {
		//	response.Result(response.ERROR, gin.H{
		//		"reload": true,
		//	}, "您的帐户异地登陆或令牌失效", c)
		//	c.Abort()
		//	return
		//}
		global.GLog.Info("444444", zap.Any("s", "i"))
		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				//response.Result(response.NotLogin, gin.H{
				//	"reload": true,
				//}, TokenExpired.Error(), c)
				//c.Abort()
				c.AbortWithStatusJSON(response.NotLogin, response.Response{
					Code: response.NotLogin,
					Data: "", Msg: TokenExpired.Error(),
				})
				return
			}
			//response.Result(response.NotLogin, gin.H{
			//	"reload": true,
			//}, err.Error(), c)
			c.AbortWithStatusJSON(response.NotLogin, response.Response{
				Code: response.NotLogin,
				Data: "", Msg: err.Error(),
			})
			c.Abort()
			return
		}
		global.GLog.Info("55555", zap.Any("s", "i"))
		c.Set("claims", claims)
		c.Next()
	}
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired        = errors.New("您的帐户令牌失效")
	TokenNoLoginExpired = errors.New("您的帐户登录状态不存在或失效")
	TokenNotValidYet    = errors.New("帐户令牌不可用")
	TokenMalformed      = errors.New("帐户令牌格式错误")
	TokenInvalid        = errors.New("不能解析该令牌")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.GConfig.JWT.SigningKey),
	}
}

// CreateToken 创建一个token
func (j *JWT) CreateToken(claims sys_base2.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析 token
func (j *JWT) ParseToken(tokenString string) (*sys_base2.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &sys_base2.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	fmt.Println("token1", time.Now().Unix())
	fmt.Println("token2", token.Claims.(*sys_base2.CustomClaims).ExpiresAt)
	if err != nil {
		global.GLog.Error("ParseToken", zap.Error(err))
		if ve, ok := err.(*jwt.ValidationError); ok {
			global.GLog.Error("ve", zap.Any("err", ve))
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				fmt.Println("Token is expired")
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*sys_base2.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid

	}

}

// RefreshToken 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &sys_base2.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*sys_base2.CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}

// whitelistsCheck 白名单检查
func whitelistsCheck(str string) bool {
	for _, s := range WhiteList {
		if match, _ := regexp.MatchString(s, str); match {
			if strings.Index(str, "socket.io") == -1 {
				global.GLog.Info("str", zap.Any("白名单检查", str))
			}
			return true
		}
	}
	return false
}
