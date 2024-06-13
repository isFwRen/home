package api

//
//import (
//	"errors"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"go.uber.org/zap"
//	"os"
//	"server/global"
//	"server/global/response"
//	"server/module/training_management/model/request"
//	"server/module/training_management/service"
//	"server/utils"
//	"time"
//)
//
//// PageTrainingManagement
//// @Tags training-management(培训管理)
//// @Summary 培训设置-分页查询培训设置
//// @Auth lhc
//// @Date 2023/11/24 21:02 下午
//// @Security ApiKeyAuth
//// @Security XCode
//// @Security XToken
//// @Security XUserId
//// @Security ProCode
//// @accept application/json
//// @Produce application/json
//// @Param pageIndex   query   int   true    "Index"
//// @Param pageSize   query   int   true    "Size"
//// @Param projectCode query string	false    "项目编码"
//// @Param isAt query string	false    "日期"
//// @Param userName query string	false    "姓名"
//// @Param userCode query string	false    "工号"
//// @Param auditStatus query int	false    "审核状态 1.待审核 2.审核通过  3.审核未通过"
//// @Success 200 {string} string "{"success":true,"data":{},"msg":"操作成功"}"
//// @Router /training-management/page [get]
//func PageTrainingManagement(c *gin.Context) {
//	var req request.ReqTraining
//	err := c.ShouldBindQuery(&req)
//	if err != nil {
//		response.FailWithMessage(err.Error(), c)
//		return
//	}
//	fmt.Println("==================req===", req)
//	err, list, total := service.PageTrainingManagement(req)
//	if err != nil {
//		response.FailWithParamErr(err, c)
//		return
//	} else {
//		response.OkWithData(request.PageResult{
//			List:      list,
//			Total:     total,
//			PageIndex: req.PageIndex,
//			PageSize:  req.PageSize,
//		}, c)
//	}
//}
//
//// InfoTrainingManagement
//// @Tags training-management(培训管理)
//// @Summary 培训设置-查询详情
//// @Auth lhc
//// @Date 2023/11/24 21:02 下午
//// @Security ApiKeyAuth
//// @Security XCode
//// @Security XToken
//// @Security XUserId
//// @Security ProCode
//// @accept application/json
//// @Produce application/json
//// @Param id   query   string   true    "id"
//// @Success 200 {string} string "{"success":true,"data":{},"msg":"操作成功"}"
//// @Router /training-management/info [get]
//func InfoTrainingManagement(c *gin.Context) {
//	id := c.Query("id")
//	if id == "" {
//		response.FailWithParamErr(errors.New("ID不能为空"), c)
//		return
//	}
//	err, info := service.InfoTrainingManagement(id)
//	if err != nil {
//		response.FailWithParamErr(err, c)
//		return
//	} else {
//		response.OkWithData(info, c)
//		return
//	}
//
//}
//
//// EditTrainingManagement
//// @Tags training-management(培训管理)
//// @Summary 培训设置-审核
//// @Auth lhc
//// @Date 2023/11/24 21:02 下午
//// @Security ApiKeyAuth
//// @Security XCode
//// @Security XToken
//// @Security XUserId
//// @Security ProCode
//// @accept application/json
//// @Produce application/json
//// @Param data   body   request.ReqAuditStatus   true    "id和审核状态1.待审核 2.审核通过  3.审核未通过"
//// @Success 200 {string} string "{"success":true,"data":{},"msg":"操作成功"}"
//// @Router /training-management/edit [post]
//func EditTrainingManagement(c *gin.Context) {
//	var reqId request.ReqAuditStatus
//	if err := c.ShouldBindJSON(&reqId); err != nil {
//		response.FailWithParamErr(err, c)
//		return
//	}
//	err := service.EditTrainingManagement(reqId)
//	if err != nil {
//		response.FailWithParamErr(err, c)
//		return
//	} else {
//		response.OkWithData("操作成功！", c)
//	}
//}
//
//// ExportTrainingManagementInfo
//// @Tags training-management(培训管理)
//// @Summary 培训设置-导出
//// @Auth lhc
//// @Date 2023/11/24 21:02 下午
//// @Security ApiKeyAuth
//// @Security XCode
//// @Security XToken
//// @Security XUserId
//// @Security ProCode
//// @accept application/json
//// @Produce application/json
//// @Param ids   body   request.ReqAuditStatus   true    "id和审核状态1.待审核 2.审核通过  3.审核未通过"
//// @Success 200 {string} string "{"success":true,"data":{},"msg":"操作成功"}"
//// @Router /training-management/export [post]
//func ExportTrainingManagementInfo(c *gin.Context) {
//	var req request.ReqAuditStatus
//	if err := c.ShouldBindJSON(&req); err != nil {
//		response.FailWithParamErr(err, c)
//		return
//	}
//	list, err := service.ExportTrainingManagementInfo(req.Ids)
//	if err != nil {
//		response.FailWithParamErr(err, c)
//		return
//	}
//	path := global.GConfig.LocalUpload.FilePath + global.PathTrainingManagementExport
//	err = os.MkdirAll(path, os.ModePerm)
//	if err != nil {
//		global.GLog.Error("upload file fail:", zap.Any("err", err))
//		return
//	}
//	name := fmt.Sprintf("培训管理-培训管理信息表-%v.xlsx", time.Now().Format("20060102"))
//	err = utils.ExportBigExcel(path, name, "sheet", list)
//	if err != nil {
//		response.FailWithMessage(err.Error(), c)
//		return
//	}
//	c.FileAttachment(path+name, name)
//}
