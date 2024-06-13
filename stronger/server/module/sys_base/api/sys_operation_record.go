package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global/response"
	sys_base2 "server/module/sys_base/model"
	request2 "server/module/sys_base/model/request"
	resp "server/module/sys_base/model/response"
	"server/module/sys_base/service"
)

// CreateSysOperationRecord
// @Tags SysOperationRecord
// @Summary 操作记录--创建SysOperationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body sys_base2.SysOperationRecord true "创建SysOperationRecord"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysOperationRecord/createSysOperationRecord [post]
func CreateSysOperationRecord(c *gin.Context) {
	var sysOperationRecord sys_base2.SysOperationRecord
	_ = c.ShouldBindJSON(&sysOperationRecord)
	err := service.CreateSysOperationRecord(sysOperationRecord)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteSysOperationRecord
// @Tags SysOperationRecord
// @Summary 操作记录--删除SysOperationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body sys_base2.SysOperationRecord true "删除SysOperationRecord"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysOperationRecord/deleteSysOperationRecord [delete]
func DeleteSysOperationRecord(c *gin.Context) {
	var sysOperationRecord sys_base2.SysOperationRecord
	_ = c.ShouldBindJSON(&sysOperationRecord)
	err := service.DeleteSysOperationRecord(sysOperationRecord)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", err), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteSysOperationRecordByIds
// @Tags SysOperationRecord
// @Summary 操作记录--批量删除SysOperationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.ReqIds true "批量删除SysOperationRecord"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysOperationRecord/deleteSysOperationRecordByIds [delete]
func DeleteSysOperationRecordByIds(c *gin.Context) {
	var IDS request2.ReqIds
	_ = c.ShouldBindJSON(&IDS)
	err := service.DeleteSysOperationRecordByIds(IDS)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", err), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// UpdateSysOperationRecord
// @Tags SysOperationRecord
// @Summary 操作记录--更新SysOperationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body sys_base2.SysOperationRecord true "更新SysOperationRecord"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sysOperationRecord/updateSysOperationRecord [put]
func UpdateSysOperationRecord(c *gin.Context) {
	var sysOperationRecord sys_base2.SysOperationRecord
	_ = c.ShouldBindJSON(&sysOperationRecord)
	err := service.UpdateSysOperationRecord(&sysOperationRecord)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindSysOperationRecord
// @Tags SysOperationRecord
// @Summary 操作记录--用id查询SysOperationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body sys_base2.SysOperationRecord true "用id查询SysOperationRecord"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sysOperationRecord/findSysOperationRecord [get]
func FindSysOperationRecord(c *gin.Context) {
	var sysOperationRecord sys_base2.SysOperationRecord
	_ = c.ShouldBindQuery(&sysOperationRecord)
	err, resysOperationRecord := service.GetSysOperationRecord(sysOperationRecord.ID)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkWithData(gin.H{"resysOperationRecord": resysOperationRecord}, c)
	}
}

// GetSysOperationRecordList
// @Tags SysOperationRecord
// @Summary 操作记录--分页获取SysOperationRecord列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.SysOperationRecordSearch true "分页获取SysOperationRecord列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysOperationRecord/getSysOperationRecordList [get]
func GetSysOperationRecordList(c *gin.Context) {
	var pageInfo request2.SysOperationRecordSearch
	_ = c.ShouldBindQuery(&pageInfo)
	err, list, total := service.GetSysOperationRecordInfoList(pageInfo)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取数据失败，%v", err), c)
	} else {
		response.OkWithData(resp.PageResult{
			List:      list,
			Total:     total,
			PageIndex: pageInfo.PageIndex,
			PageSize:  pageInfo.PageSize,
		}, c)
	}
}
