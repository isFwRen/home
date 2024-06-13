package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global/response"
	sys_base2 "server/module/sys_base/model"
	response2 "server/module/sys_base/model/response"
	"server/module/sys_base/service"
)

func GetTabsList(c *gin.Context) {
	var tabs sys_base2.TabsModel
	_ = c.ShouldBindJSON(&tabs)

	err, reSysTabsList := service.GetTabsList()
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkDetailed(reSysTabsList, "查询成功", c)
	}
}

func UpdateTabs(c *gin.Context) {
	var tabs sys_base2.TabsModel
	_ = c.ShouldBindJSON(&tabs)

	row := service.UpdateTabs(tabs)
	if row == 0 {
		response.FailWithMessage(fmt.Sprintf("更新失败"), c)
	} else {
		response.OkDetailed(response2.RowResult{
			Row: row,
		}, "更新成功", c)
	}
}

func GetTabsLast(c *gin.Context) {
	var tabs sys_base2.TabsModel
	_ = c.ShouldBindJSON(&tabs)

	err, reSysTabs := service.GetTabsLast()
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkDetailed(reSysTabs, "查询成功", c)
	}
}


func AddTabs(c *gin.Context){
	var tabs sys_base2.TabsModel
	_ = c.ShouldBindJSON(&tabs)

	err := service.AddTabs(tabs)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithData("创建成功", c)
	}
}

func RemoveTabs(c *gin.Context){
	var tabs sys_base2.TabsModel
	_ = c.ShouldBindJSON(&tabs)
	if tabs.Name == "0" {
		response.FailWithMessage("主页不能移除!", c)
		return
	}
	rows := service.RemoveTabs(tabs)
	if rows == 0 {
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkDetailed(response2.RowResult{
			Row: rows,
		}, "删除成功", c)
	}

}

