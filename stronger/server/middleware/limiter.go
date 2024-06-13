/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/5 3:45 下午
 */

package middleware

import (
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"server/global/response"
)

//Limiter 自定义限流中间件
func Limiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			//c.Data(httpError.StatusCode, lmt.GetMessageContentType(), []byte(httpError.Message))
			response.FailWithMessage(fmt.Sprintf("请求失败，点慢点，我裂开啦"), c)
			//response.FailWithMessage(fmt.Sprintf("请求失败，%v", httpError.Message), c)
			c.Abort()
		} else {
			c.Next()
		}
	}
}
