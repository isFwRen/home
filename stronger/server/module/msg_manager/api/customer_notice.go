/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/10/18 09:52
 */

package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"server/global"
	"server/global/response"
	"server/module/msg_manager/model"
	"server/module/msg_manager/service"
	responseSysBase "server/module/sys_base/model/response"
	"server/utils"
	"strconv"
	"strings"
	"time"
	api2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"
)

// GetCustomerNoticeByPage
// @Tags msg-manager(消息管理--客户通知)
// @Summary 消息管理--获取客户通知分页
// @Auth xingqiyi
// @Date 2022/7/4 16:21
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode        query   string   false    "项目代码"
// @Param pageIndex      query   int      true    "页码"
// @Param pageSize       query   int      true    "数量"
// @Param startTime      query   string   true   "开始时间今天开始时间格式'2022-07-06T16:00:00.000Z'"
// @Param endTime        query   string   true   "结束时间今天结束时间格式'2022-07-06T16:00:00.000Z'"
// @Param msgType        query   int      false    "消息类型"
// @Param status       	 query   int      false    "状态"
// @Param orderBy        query   string   false   "排序JSON.stringify([['CreatedAt','desc']])"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /msg-manager/customer-notice/page [get]
func GetCustomerNoticeByPage(c *gin.Context) {
	var customerNoticeSearch model.CustomerNoticeSearchReq
	err := c.ShouldBindQuery(&customerNoticeSearch)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, total, groups := service.GetCustomerNoticeByPage(customerNoticeSearch)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		err, count := service.GetCustomerNoticeByStatus(customerNoticeSearch.ProCode)
		if err != nil {
			global.GLog.Error("", zap.Error(err))
			response.FailWithMessage(fmt.Sprintf("获取失败"), c)
			return
		}
		response.OkWithData(responseSysBase.BasePageResult{
			List:      groups,
			Total:     total,
			MaxOrder:  count,
			PageIndex: customerNoticeSearch.PageIndex,
			PageSize:  customerNoticeSearch.PageSize,
		}, c)
	}
}

// Reply
// @Tags msg-manager(消息管理--客户通知)
// @Summary 消息管理--回复客户
// @Auth xingqiyi
// @Date 2022/7/6 14:21
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.CustomerNotice true "客户通知实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"回复成功"}"
// @Router /msg-manager/customer-notice/reply [post]
func Reply(c *gin.Context) {
	var customerNotice model.CustomerNotice
	err := c.ShouldBindJSON(&customerNotice)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	customClaims, err := api2.GetUserByToken(c)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}

	verify := utils.Rules{
		"ID":      {utils.NotEmpty()},
		"ProCode": {utils.NotEmpty()},
	}
	verifyErr := utils.Verify(customerNotice, verify)
	if verifyErr != nil {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}

	err, customerNotice = service.GetCustomerNoticeById(customerNotice)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	customerNotice.DealUserCode = customClaims.Code
	customerNotice.DealUserName = customClaims.NickName

	//生成xml
	err, xmlPath, xmlDir := creatXML(customerNotice)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("生成xml失败失败，%v", err), c)
		return
	}

	//传到客户路径receive_response
	project, ok := global.GProConf[customerNotice.ProCode]
	if !ok {
		global.GLog.Error(fmt.Sprintf("回传错误id:::%v,错误:::没有找到项目配置:::%v,%v", customerNotice.ID, err, customerNotice.ProCode))
		return
	}
	sh := fmt.Sprintf(project.UploadPaths.Upload, xmlDir, customerNotice.FileName, customerNotice.FileName)
	sh = strings.Replace(sh, "Receive_claim", "receive_response", -1)
	global.GLog.Info("上传cmd:::", zap.Any("sh", sh))
	err, s, s2 := utils.ShellOut(sh)
	global.GLog.Info("上传回显", zap.Any("std out", s))
	global.GLog.Error("上传回显", zap.Any("std err", s2))
	if err != nil {
		global.GLog.Error("上传错误err", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	sh = fmt.Sprintf(project.UploadPaths.UploadRename, customerNotice.FileName, customerNotice.FileName)
	sh = strings.Replace(sh, "Receive_claim", "receive_response", -1)
	global.GLog.Info("上传重命名cmd:::", zap.Any("sh", sh))
	err, s, s2 = utils.ShellOut(sh)
	global.GLog.Info("上传重命名回显", zap.Any("std out", s))
	global.GLog.Error("上传重命名回显", zap.Any("std err", s2))
	if err != nil {
		global.GLog.Error("上传重命名错误err", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	//更新数据库
	customerNotice.UploadPath = xmlPath
	err = service.Reply(customerNotice)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkWithData("更新成功", c)
	}
}

// creatXML 生成xml
func creatXML(customerNotice model.CustomerNotice) (err error, xmlFile, xmlFileDir string) {
	timeStr := time.Now().Format("2006/01/02 15:04:05")
	isNo := "N"
	if customerNotice.IsReply {
		isNo = "Y"
	}
	//理赔：30， 秒赔：33
	typeId := "30"
	if customerNotice.MsgType == 2 {
		typeId = "33"
	}
	xmlStr := `
		<?xml version="1.0" encoding="utf-8"?>
		<NEWCASE>
		<RESPONSE_ID>W` + customerNotice.FileName + `</RESPONSE_ID> 
		<REQUEST_ID>` + customerNotice.FileName + `</REQUEST_ID>
		<RESPONSE_TIME>` + timeStr + `</RESPONSE_TIME>
		<TYPE_ID>` + typeId + `</TYPE_ID>
		<COMPANY_ID>09</COMPANY_ID>
		<IS_NO>` + isNo + `</IS_NO>
		<COUNT>` + strconv.Itoa(customerNotice.ExpectNum) + `</COUNT> 
		</NEWCASE>
	`
	err, xmlFile, xmlFileDir = saveXml(customerNotice, xmlStr)
	return err, xmlFile, xmlFileDir
}

// saveXml 保存成xml文件
func saveXml(customerNotice model.CustomerNotice, xml string) (err error, xmlFile, xmlFileDir string) {
	xmlFileDir = global.GConfig.LocalUpload.FilePath + customerNotice.ProCode + "/upload_xml/" +
		fmt.Sprintf("%v/%v/%v/%v/%v/",
			customerNotice.SendTime.Year(), int(customerNotice.SendTime.Month()),
			customerNotice.SendTime.Day(), global.PathCustomerNotice, customerNotice.MsgType)
	exists, err := utils.PathExists(xmlFileDir)
	if err != nil {
		return err, "", ""
	}
	if !exists {
		err = utils.CreateDir(xmlFileDir)
		if err != nil {
			return err, "", ""
		}
	}
	xmlFile = xmlFileDir + "/" + customerNotice.FileName + ".xml"
	err = ioutil.WriteFile(xmlFile, []byte(xml), 0666)
	return err, xmlFile, xmlFileDir
}
