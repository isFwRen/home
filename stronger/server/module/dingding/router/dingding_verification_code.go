package router

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	api2 "server/module/dingding/api"
)

func InitDingingBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	// Create a limiter struct.
	// This is a generic middleware to rate-limit HTTP requests.
	limiter := tollbooth.NewLimiter(2, nil)

	DingingBase := Router.Group("dinging/")
	{
		DingingBase.POST("captcha",tollbooth_gin.LimitHandler(limiter), api2.Captcha) //发送验证码（数字)
	}
	return DingingBase
}