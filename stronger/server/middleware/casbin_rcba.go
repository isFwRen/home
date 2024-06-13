package middleware

import (
	"fmt"
	"regexp"
	"server/global"
	"server/global/response"
	sys_base2 "server/module/sys_base/model/request"
	"server/module/sys_base/service"

	"github.com/gin-gonic/gin"
)

// CasbinHandler 权限拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if global.GConfig.System.Env == "develop" {
			c.Next()
			return
		}
		claims, _ := c.Get("claims")
		waitUse := claims.(*sys_base2.CustomClaims)
		// 获取请求的URI (请求的api)
		obj := c.Request.URL.RequestURI()
		// 获取请求方法 （可以换成项目用请求方法也行每个项目的api应该都差不多的）pro+role
		act := c.Request.Method
		// 获取用户的角色 （roleId）
		sub := waitUse.RoleId
		fmt.Println(waitUse)
		fmt.Println("--------------------------", obj, sub, act)
		if waitUse.Code == "6482" || waitUse.RoleName == "admin" {
			c.Next()
			return
		}
		e, err := service.GoCasbin()
		fmt.Println("------------e--------------", e)
		if err != nil {
			response.Result(response.NotLogin, gin.H{}, err.Error(), c)
			c.Abort()
			return
		}
		// 判断策略中是否存在
		success, err := e.Enforce(sub, obj, act)
		fmt.Println("-----------success-------------", success)
		if err != nil {
			response.Result(response.NotLogin, gin.H{}, err.Error(), c)
			c.Abort()
			return
		}
		if global.GConfig.System.Env == "develop" || success {
			c.Next()
		} else {
			fmt.Println("-----------!!!!!!!!!!!!!----------")
			response.Result(response.NotLogin, gin.H{}, global.RoleErr.Error(), c)
			c.Abort()
			return
		}
	}
}

// PermHandler 项目权限拦截器
func PermHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.Request.Header.Get("x-user-id")    //登录用户id
		proCode := c.Request.Header.Get("pro-code") //当前录入项目编码
		process := c.Request.Header.Get("process")  //领取的流程op0、op1、op2、opq
		uri := c.Request.RequestURI                 //uri
		if noNeedProcess(uri) {
			c.Next()
			return
		}
		if uid == "" || proCode == "" || process == "" {
			response.Result(response.ERROR, gin.H{}, "没有传x-user-id、pro-code或process", c)
			c.Abort()
			return
		}
		err, has := service.GetRedisPerm(uid, proCode, process)
		if err != nil || len(has) != 1 {
			global.GLog.Error(fmt.Sprintf("权限长度：%d", len(has)))
			response.Result(response.ERROR, gin.H{}, err.Error(), c)
			c.Abort()
			return
		}

		if global.GConfig.System.Env == "develop" || has[0] == "1" {
			c.Next()
		} else {
			response.Result(response.ERROR, gin.H{}, global.RoleErr.Error(), c)
			c.Abort()
			return
		}
	}
}

// noNeedProcess 不需要传流程在header
func noNeedProcess(str string) bool {
	List := []string{
		"/task/opNum",
		"/task/conf",
		"/task/releaseBlock",
		"/task/releaseBill",
		// "/task/op",
		// "/task/submit",
	}
	for _, s := range List {
		if match, _ := regexp.MatchString(s, str); match {
			return true
		}
	}
	return false
}
