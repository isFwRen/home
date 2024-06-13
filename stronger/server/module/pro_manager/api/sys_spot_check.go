package api

import (
	"fmt"
	"server/global"
	"server/global/response"
	"server/module/pro_manager/model"
	request2 "server/module/pro_manager/model/request"
	"server/module/pro_manager/service"
	responseSysBase "server/module/sys_base/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetSysSpotCheckByPage(c *gin.Context) {
	// fmt.Println("-------------------------------------------")
	var sysSpotCheckQuery request2.SysSpotCheckQuery
	err := c.ShouldBindQuery(&sysSpotCheckQuery)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	// fmt.Println("--------------------sysSpotCheckQuery-----------------------", sysSpotCheckQuery)
	err, total, list := service.GetSysSpotCheckByPage(sysSpotCheckQuery)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List:      list,
			Total:     total,
			PageIndex: sysSpotCheckQuery.PageIndex,
			PageSize:  sysSpotCheckQuery.PageSize,
		}, c)
	}
}

func AddSysSpotCheck(c *gin.Context) {
	var sysSpotCheck model.SysSpotCheck
	err := c.ShouldBindJSON(&sysSpotCheck)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// sysSpotCheck.ReleaseDate = time.Now()

	// sysSpotCheck.Status = 1
	err = service.InsertSysSpotCheck(sysSpotCheck)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithData("创建成功", c)
	}
}

func UpdateSysSpotCheck(c *gin.Context) {
	var sysSpotCheck model.SysSpotCheck
	err := c.ShouldBindJSON(&sysSpotCheck)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = service.UpdateSysSpotCheckById(sysSpotCheck, sysSpotCheck.ID)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkWithData("更新成功", c)
	}
}
