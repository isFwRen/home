package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/lib/pq"
	"os"
	"path"
	"server/global"
	re "server/global/request"
	"server/global/response"
	"server/module/pro_conf/model"
	model2 "server/module/pro_manager/model"
	"server/module/pro_manager/model/request"
	resp "server/module/pro_manager/model/response"
	"server/module/pro_manager/service"
	"server/utils"
	"strings"
)

// FieldsRuleSync
// @Tags Fields Rule (字段规则)
// @Summary 字段规则--同步数据
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/fieldsRule/sync [get]
func FieldsRuleSync(c *gin.Context) {
	proCode := c.Query("proCode")
	err := service.FieldsRuleSync(proCode)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("同步数据失败, err : %v", err), c)
	} else {
		response.Ok(c)
	}
}

// GetFieldsRuleList
// @Tags Fields Rule (字段规则)
// @Summary 字段规则--查询
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param fieldsName query string true "字段名称"
// @Param rule query string true "录入规则"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/fieldsRule/list [get]
func GetFieldsRuleList(c *gin.Context) {
	var search request.GFR
	if err := c.ShouldBindWith(&search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetFieldsRuleListVerify := utils.Rules{
		"proCode": {utils.NotEmpty()},
	}
	GetFieldsRuleListVerifyErr := utils.Verify(search, GetFieldsRuleListVerify)
	if GetFieldsRuleListVerifyErr != nil {
		response.FailWithMessage(GetFieldsRuleListVerifyErr.Error(), c)
		return
	}
	err, list, total := service.GetFieldsRuleList(search)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("查询字段规则失败, err : %v", err), c)
	} else {
		response.GetOkWithData(resp.QualityRes{List: list, Total: total}, c)
	}
}

// EditFieldsRule
// @Tags Fields Rule (字段规则)
// @Summary 字段规则--编辑
// @accept application/json
// @Produce application/json
// @Param id body string true "id"
// @Param file formData file true "文件"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/fieldsRule/edit [post]
func EditFieldsRule(c *gin.Context) {
	form, _ := c.MultipartForm()
	id := form.Value["id"]
	if len(id) == 0 {
		response.GetFailWithData(fmt.Sprintf("id不能为空"), c)
		return
	}

	// 获取上传文件
	Files := form.File["file"]
	for _, p := range Files {
		pName := p.Filename
		types := strings.Split(pName, ".")[len(strings.Split(pName, "."))-1]
		if types != "png" && types != "PNG" && types != "jpg" && types != "JPG" && types != "tif" && types != "TIF" {
			response.FailWithMessage("每张影像格式必须为png、PNG、jpg、JPG、tif和TIF", c)
			return
		}
		//if types != "PNG" && types != "JPG" && types != "TIF" {
		//	response.FailWithMessage("每张影像格式必须为PNG、JPG和TIF", c)
		//	return
		//}
	}

	var fieldRule model2.SysFieldRule
	err := global.GDb.Model(&model2.SysFieldRule{}).Where("id in ? ", id).Find(&fieldRule).Error
	if err != nil {
		fmt.Println(err.Error())
		response.GetFailWithData(fmt.Sprintf("编辑失败"), c)
		return
	}

	var proInformation model.SysProject
	err = global.GDb.Model(&model.SysProject{}).Where("id = ? ", fieldRule.ProId).Find(&proInformation).Error
	if err != nil {
		fmt.Println(err.Error())
		response.GetFailWithData(fmt.Sprintf("编辑失败"), c)
		return
	}

	for _, file := range Files {
		filename := file.Filename

		// 线上 uploads/file
		basicPath := global.GConfig.LocalUpload.FilePath + "/字段规则/" + proInformation.Code + "/"
		// 设置文件需要保存的指定位置并设置保存的文件名字
		dst := path.Join(basicPath, filename)
		fmt.Println(filename + " has been saved in this path " + basicPath)
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
		fieldRule.RulePicture = pq.StringArray{}
		fieldRule.RulePicture = append(fieldRule.RulePicture, dst)
	}
	fieldRule.InputRule = "有"
	err = service.EditFieldsRule(fieldRule)
	if err != nil {
		fmt.Println(err)
		response.GetFailWithData(fmt.Sprintf("编辑失败"), c)
	} else {
		response.Ok(c)
	}
}

// DeleteFieldsRule
// @Tags Fields Rule (字段规则)
// @Summary 字段规则--删除
// @accept application/json
// @Produce application/json
// @Param data body re.RmById true "id"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/fieldsRule/delete [post]
func DeleteFieldsRule(c *gin.Context) {
	var reqId re.RmById
	_ = c.ShouldBindJSON(&reqId)
	IdVerifyErr := utils.Verify(reqId, utils.CustomizeMap["IdVerify"])
	if IdVerifyErr != nil {
		response.FailWithMessage(IdVerifyErr.Error(), c)
		return
	}
	err := service.DeleteFieldsRule(reqId.Ids)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", err), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// ExportFieldsRule
// @Tags Fields Rule (字段规则)
// @Summary 字段规则--导出
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param rule query string true "录入规则"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/fieldsRule/export [get]
func ExportFieldsRule(c *gin.Context) {
	var req request.ExportFR
	if err := c.ShouldBindWith(&req, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	ExportFieldsRuleVerify := utils.Rules{
		"proCode": {utils.NotEmpty()},
	}
	ExportFieldsRuleVerifyErr := utils.Verify(req, ExportFieldsRuleVerify)
	if ExportFieldsRuleVerifyErr != nil {
		response.FailWithMessage(ExportFieldsRuleVerifyErr.Error(), c)
		return
	}
	err, p := service.ExportFieldsRule(req)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("导出失败，%v", err), c)
	} else {
		response.GetOkWithData(resp.FRReq{List: p}, c)
	}
}

// UploadFieldsRule
// @Tags Fields Rule (字段规则)
// @Summary 字段规则--导入
// @accept application/json
// @Produce application/json
// @Param proCode body string true "项目编码"
// @Param file formData file true "文件"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/fieldsRule/upload [post]
func UploadFieldsRule(c *gin.Context) {
	form, _ := c.MultipartForm()
	proCode := form.Value["proCode"]
	if len(proCode) == 0 {
		response.GetFailWithData(fmt.Sprintf("proCode不能为空"), c)
		return
	}

	// 获取上传文件
	Files := form.File["file"]
	for _, p := range Files {
		pName := p.Filename
		types := strings.Split(pName, ".")[len(strings.Split(pName, "."))-1]
		if types != "png" && types != "PNG" && types != "jpg" && types != "JPG" && types != "tif" && types != "TIF" {
			response.FailWithMessage("每张影像格式必须为png、PNG、jpg、JPG、tif和TIF", c)
			return
		}
		//if types != "PNG" && types != "JPG" && types != "TIF" {
		//	response.FailWithMessage("每张影像格式必须为PNG、JPG和TIF", c)
		//	return
		//}
	}

	err := service.UploadFieldsRule(Files, proCode, c)
	if err != nil {
		fmt.Println(err)
		response.FailWithMessage(fmt.Sprintf("导入失败, %s", err.Error()), c)
	} else {
		response.Ok(c)
	}
}
