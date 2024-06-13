package api

import (
	"fmt"
	"server/global"
	"server/global/response"
	model2 "server/module/pro_conf/model"
	"server/module/pro_conf/service"
	responseSysBase "server/module/sys_base/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetFieldCheckByFieldId(c *gin.Context) {
	fmt.Println("--------------GetFieldCheckByFieldId----------------------")
	var SysProFieldCheckEditReq model2.SysProFieldCheck
	err := c.ShouldBindQuery(&SysProFieldCheckEditReq)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, list := service.GetFieldCheckListByFieldId(SysProFieldCheckEditReq.FId)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List:      list,
			Total:     0,
			PageIndex: 1,
			PageSize:  0,
		}, c)
	}
}

func EditFieldCheckByFieldId(c *gin.Context) {
	fmt.Println("--------------EditFieldCheckByFieldId----------------------")
	var reqParam model2.SysFieldCheckEditReq
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err = service.EditFieldCheckByFieldId(reqParam.FId, reqParam.List)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("编辑失败"), c)
	} else {
		response.OkWithData("编辑成功", c)
	}
}
