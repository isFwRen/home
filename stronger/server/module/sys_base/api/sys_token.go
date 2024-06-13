/**
 * @Author: 星期一
 * @Description:
 * @Date: 2021/1/5 上午10:49
 */

package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"server/global"
	"server/middleware"
	"server/module/sys_base/model/request"
)

// GetUserByToken 在token获取user信息
func GetUserByToken(c *gin.Context) (customClaims *request.CustomClaims, err error) {
	r := c.Request
	token := r.Header.Get(global.XToken)
	if token == "" {
		return customClaims, errors.New("token为空")
	}
	customClaims, err = middleware.NewJWT().ParseToken(token)
	if err != nil {
		return customClaims, errors.New("token不正确")
	}
	return customClaims, nil
}
