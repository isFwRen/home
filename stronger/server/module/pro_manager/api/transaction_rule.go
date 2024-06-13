package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	scheduleService "server/module/training_guide/service"
	api2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"

	//"gorm.io/gorm"
	"os"
	"path"
	//"reflect"
	"server/global"
	"server/global/response"
	"server/module/pro_manager/model/request"
	response2 "server/module/pro_manager/model/response"
	server "server/module/pro_manager/service"
	"server/module/sys_base/model"
	"server/utils"
	"strings"
)

// GetTaskTransactionRule
// @Tags Transaction Rule (业务规则(录入系统))
// @Summary	业务规则(录入系统)--查询
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param ruleType query string true "规则类型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/transactionRule/task/list [get]
func GetTaskTransactionRule(c *gin.Context) {
	var search request.GetTransactionRule
	if err := c.ShouldBindWith(&search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetTransactionRuleVerify := utils.Rules{
		"proCode": {utils.NotEmpty()},
	}
	GetTransactionRuleVerifyErr := utils.Verify(search, GetTransactionRuleVerify)
	if GetTransactionRuleVerifyErr != nil {
		response.FailWithMessage(GetTransactionRuleVerifyErr.Error(), c)
		return
	}
	err, list, total := server.GetTransactionRule(search)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("查询业务规则失败, err : %v", err), c)
	} else {
		user, err1 := api2.GetUserByToken(c)
		if err1 != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		needReadList, err2 := scheduleService.GetNeedReadList2(1, user.ID)
		if err2 != nil {
			response.FailWithMessage(err2.Error(), c)
			return
		}

		convMap, err3 := scheduleService.InsertIsReadStatus(needReadList, list)
		if err3 != nil {
			response.FailWithMessage(err2.Error(), c)
			return
		}
		if convMap == nil {
			response.OkWithData(response2.TransactionRuleRes{
				List:  list,
				Total: total,
			}, c)
			return
		}
		response.OkWithData(response2.TransactionRuleRes{
			List:  convMap,
			Total: total,
		}, c)
	}
}

// GetTransactionRule
// @Tags Transaction Rule (业务规则)
// @Summary	业务规则--查询
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param ruleType query string true "规则类型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/transactionRule/list [get]
func GetTransactionRule(c *gin.Context) {
	var search request.GetTransactionRule
	if err := c.ShouldBindWith(&search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetTransactionRuleVerify := utils.Rules{
		"proCode": {utils.NotEmpty()},
	}
	GetTransactionRuleVerifyErr := utils.Verify(search, GetTransactionRuleVerify)
	if GetTransactionRuleVerifyErr != nil {
		response.FailWithMessage(GetTransactionRuleVerifyErr.Error(), c)
		return
	}
	err, list, total := server.GetTransactionRule(search)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("查询业务规则失败, err : %v", err), c)
	} else {
		response.OkWithData(response2.TransactionRuleRes{
			List:  list,
			Total: total,
		}, c)
	}
}

// AddTransactionRule
// @Tags Transaction Rule (业务规则)
// @Summary	业务规则--新增
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data formData request.AddTransactionRule true "新增"
// @Param file formData file true "文件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/transactionRule/add [post]
func AddTransactionRule(c *gin.Context) {
	var add request.AddTransactionRule
	err := c.ShouldBind(&add)
	fmt.Println(err)
	GetTransactionRuleVerify := utils.Rules{
		"proCode": {utils.NotEmpty()},
	}
	GetTransactionRuleVerifyErr := utils.Verify(add, GetTransactionRuleVerify)
	if GetTransactionRuleVerifyErr != nil {
		response.FailWithMessage(GetTransactionRuleVerifyErr.Error(), c)
		return
	}
	form, _ := c.MultipartForm()
	excelFile := form.File["file"]
	if len(excelFile) == 0 {
		response.FailWithMessage("没有对应的规则文件", c)
		return
	}
	for _, file := range excelFile {
		fileName := file.Filename

		if file.Size > 5*1048576 {
			response.FailWithMessage("文件大小要限制在5MB之内", c)
			return
		}
		if strings.Index(fileName, ".pdf") == -1 {
			response.FailWithMessage("文件格式不对", c)
			return
		}
		basicPath := global.GConfig.LocalUpload.FilePath + "/业务规则/" + add.ProCode + "/"
		dst := path.Join(basicPath, fileName)
		fmt.Println("dst", dst)
		isExist, err := exists(basicPath)
		if err != nil {
			return
		}
		if !isExist {
			err = os.MkdirAll(basicPath, 0777)
			if err != nil {
				return
			}
			err = os.Chmod(basicPath, 0777)
			if err != nil {
				return
			}
		}
		// 上传文件到指定的路径
		saveErr := c.SaveUploadedFile(file, dst)
		if saveErr != nil {
			fmt.Println("saveErr", saveErr)
			continue
		}

		add.DocsPath = dst
	}
	err = server.AddTransactionRule(add)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("新增业务规则失败, err : %v", err), c)
	} else {
		response.Ok(c)
	}
}

// EditTransactionRule
// @Tags Transaction Rule (业务规则)
// @Summary	业务规则--修改
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.AddTransactionRule true "修改"
// @Param file formData file true "文件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/transactionRule/edit [post]
func EditTransactionRule(c *gin.Context) {
	var add request.AddTransactionRule
	_ = c.ShouldBind(&add)

	EditTransactionRuleVerify := utils.Rules{
		"proCode": {utils.NotEmpty()},
	}
	EditTransactionRuleVerifyErr := utils.Verify(add, EditTransactionRuleVerify)
	if EditTransactionRuleVerifyErr != nil {
		response.FailWithMessage(EditTransactionRuleVerifyErr.Error(), c)
		return
	}
	uid := c.Request.Header.Get("x-user-id")
	if uid == "" {
		response.FailWithMessage("x-user-id is empty", c)
		return
	}
	var user model.SysUser
	err := global.GDb.Model(&model.SysUser{}).Where("id = ? ", uid).Find(&user).Error
	if err != nil {
		response.FailWithMessage("没有找到修改人的基本信息", c)
		return
	}
	form, _ := c.MultipartForm()
	excelFile := form.File["file"]
	for _, file := range excelFile {
		fileName := file.Filename
		if file.Size > 5*1048576 {
			response.FailWithMessage("文件大小要限制在5MB之内", c)
			return
		}
		if strings.Index(fileName, ".pdf") == -1 {
			response.FailWithMessage("文件格式不对", c)
			return
		}
		basicPath := global.GConfig.LocalUpload.FilePath + "/业务规则/" + add.ProCode + "/"
		dst := path.Join(basicPath, fileName)
		fmt.Println("dst", dst)
		isExist, err := exists(basicPath)
		if err != nil {
			return
		}
		if !isExist {
			err = os.MkdirAll(basicPath, 0777)
			if err != nil {
				return
			}
			err = os.Chmod(basicPath, 0777)
			if err != nil {
				return
			}
		}
		// 上传文件到指定的路径
		saveErr := c.SaveUploadedFile(file, dst)
		if saveErr != nil {
			fmt.Println("saveErr", saveErr)
			continue
		}

		add.DocsPath = dst
	}
	add.UpdatedName = user.NickName
	err = server.EditTransactionRule(add)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("修改业务规则失败, err : %v", err), c)
	} else {
		response.Ok(c)
	}
}

// DeleteTransactionRule
// @Tags Transaction Rule (业务规则)
// @Summary	业务规则--删除
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.Rm true "id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/transactionRule/delete [post]
func DeleteTransactionRule(c *gin.Context) {
	var reqId request.Rm
	_ = c.ShouldBindJSON(&reqId)
	IdVerifyErr := utils.Verify(reqId, utils.CustomizeMap["IdVerify"])
	if IdVerifyErr != nil {
		response.FailWithMessage(IdVerifyErr.Error(), c)
		return
	}
	err := server.DeleteTransactionRule(reqId)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", err), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}
