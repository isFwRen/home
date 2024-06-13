package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"os"
	"path"
	"server/global"
	"server/global/response"
	"server/module/pro_manager/model"
	"server/module/pro_manager/model/request"
	resp "server/module/pro_manager/model/response"
	server "server/module/pro_manager/service"
	scheduleService "server/module/training_guide/service"
	"server/utils"
	"strconv"
	"strings"
	"time"
	api2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"
	model2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/model"
)

// GetTaskReimbursementFormTemplate
// @Tags Reimbursement Form Template (报销单模板(录入系统))
// @Summary 报销单模板(录入系统)--查询
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param name query string true "影像名称"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/reimbursementFormTemplate/task/list [get]
func GetTaskReimbursementFormTemplate(c *gin.Context) {
	var search request.GetRFT
	if err := c.ShouldBindWith(&search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}

	GetReimbursementFormTemplateVerify := utils.Rules{
		"ProCode": {utils.NotEmpty()},
	}
	GetReimbursementFormTemplateVerifyErr := utils.Verify(search, GetReimbursementFormTemplateVerify)
	if GetReimbursementFormTemplateVerifyErr != nil {
		response.FailWithMessage(GetReimbursementFormTemplateVerifyErr.Error(), c)
		return
	}
	err, list, total := server.GetReimbursementFormTemplate(search)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("查询报销单模板失败, err : %v", err), c)
	} else {

		user, err1 := api2.GetUserByToken(c)
		if err1 != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		needReadList, err2 := scheduleService.GetNeedReadList2(2, user.ID)
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
			response.GetOkWithData(resp.RFT{List: list, Total: total}, c)
			return
		}

		response.GetOkWithData(resp.RFT{List: convMap, Total: total}, c)
	}
}

// GetReimbursementFormTemplate
// @Tags Reimbursement Form Template (报销单模板)
// @Summary 报销单模板--查询
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param name query string true "影像名称"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/reimbursementFormTemplate/list [get]
func GetReimbursementFormTemplate(c *gin.Context) {
	var search request.GetRFT
	if err := c.ShouldBindWith(&search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}

	GetReimbursementFormTemplateVerify := utils.Rules{
		"ProCode": {utils.NotEmpty()},
	}
	GetReimbursementFormTemplateVerifyErr := utils.Verify(search, GetReimbursementFormTemplateVerify)
	if GetReimbursementFormTemplateVerifyErr != nil {
		response.FailWithMessage(GetReimbursementFormTemplateVerifyErr.Error(), c)
		return
	}
	err, list, total := server.GetReimbursementFormTemplate(search)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("查询报销单模板失败, err : %v", err), c)
	} else {
		response.GetOkWithData(resp.RFT{List: list, Total: total}, c)
	}
}

// AddReimbursementFormTemplate
// @Tags Reimbursement Form Template (报销单模板)
// @Summary 报销单模板--新增
// @accept application/json
// @Produce application/json
// @Param proCode formData string true "项目编码"
// @Param isRequired formData string true "是否必学"
// @Param file formData file true "文件"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/reimbursementFormTemplate/add [post]
func AddReimbursementFormTemplate(c *gin.Context) {
	var Pics []model.PicInformation
	form, _ := c.MultipartForm()
	proCode := form.Value["proCode"]
	isRequired := form.Value["isRequired"]
	if len(proCode) == 0 {
		response.FailWithMessage("项目编码不能为空", c)
		return
	}
	uid := c.Request.Header.Get("x-user-id")
	if uid == "" {
		response.FailWithMessage("x-user-id is empty", c)
		return
	}
	var user model2.SysUser
	err := global.GUserDb.Model(&model2.SysUser{}).Where("id = ? ", uid).Find(&user).Error
	if err != nil {
		response.FailWithMessage("没有找到修改人的基本信息", c)
		return
	}
	pic := form.File["file"]
	if len(pic) > 50 {
		response.FailWithMessage("影像上传数量不能超过50张", c)
		return
	}
	for _, p := range pic {
		pName := p.Filename
		types := strings.Split(pName, ".")[len(strings.Split(pName, "."))-1]
		if types != "png" && types != "PNG" && types != "jpg" && types != "JPG" && types != "tif" && types != "TIF" {
			response.FailWithMessage("每张影像格式必须为png、PNG、jpg、JPG、tif和TIF", c)
			return
		}
		if p.Size > 2*1024*1024 {
			response.FailWithMessage("每张影像的大小限制为2M", c)
			return
		}
	}
	for _, p := range pic {
		var PicsItem model.PicInformation
		pName := p.Filename
		// 线上 uploads/file
		basicPath := global.GConfig.LocalUpload.FilePath + "/报销单模板/" + proCode[0] + "/"
		// 本地
		//basicPath := "./Salary/" + filename[:4] + "/"
		// 设置文件需要保存的指定位置并设置保存的文件名字
		dst := path.Join(basicPath, pName)
		fmt.Println(pName + " has been saved in this path " + basicPath)
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
		saveErr := c.SaveUploadedFile(p, dst)
		if saveErr != nil {
			fmt.Println("saveErr", saveErr)
			continue
		}
		PicsItem.ProCode = proCode[0]
		PicsItem.IsRequired, _ = strconv.Atoi(isRequired[0])
		for i, v := range strings.Split(pName, ".") {
			if i == len(strings.Split(pName, "."))-1 {
				break
			}
			if i != 0 && i != len(strings.Split(pName, "."))-1 {
				PicsItem.Name += "."
			}
			PicsItem.Name += v
		}
		PicsItem.EditTime = time.Now().Format("2006-01-02 15:04:05")
		PicsItem.Path = dst
		PicsItem.Size = float64(p.Size)
		PicsItem.EditName = user.Name
		PicsItem.Types = strings.Split(pName, ".")[len(strings.Split(pName, "."))-1]
		Pics = append(Pics, PicsItem)
	}
	err = server.AddReimbursementFormTemplate(proCode[0], Pics)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("新增报销单模板失败, err : %v", err), c)
	} else {
		response.Ok(c)
	}
}

// DeleteReimbursementFormTemplate
// @Tags Reimbursement Form Template (报销单模板)
// @Summary	报销单模板--删除
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.Rm true "id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/reimbursementFormTemplate/delete [post]
func DeleteReimbursementFormTemplate(c *gin.Context) {
	var reqId request.Rm
	_ = c.ShouldBindJSON(&reqId)
	IdVerifyErr := utils.Verify(reqId, utils.CustomizeMap["IdVerify"])
	if IdVerifyErr != nil {
		response.FailWithMessage(IdVerifyErr.Error(), c)
		return
	}
	err := server.DeleteReimbursementFormTemplate(reqId)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", err), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// ReNameReimbursementFormTemplate
// @Tags Reimbursement Form Template (报销单模板)
// @Summary	报销单模板--重命名
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode body string true "项目编码"
// @Param id body string true "id"
// @Param name body string true "新文件名称"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/reimbursementFormTemplate/rename [post]
func ReNameReimbursementFormTemplate(c *gin.Context) {
	var rn request.RFTRename
	_ = c.ShouldBindJSON(&rn)
	ReNameReimbursementFormTemplateVerify := utils.Rules{
		"Id":   {utils.NotEmpty()},
		"Name": {utils.NotEmpty()},
	}
	ReNameReimbursementFormTemplateVerifyErr := utils.Verify(rn, ReNameReimbursementFormTemplateVerify)
	if ReNameReimbursementFormTemplateVerifyErr != nil {
		response.FailWithMessage(ReNameReimbursementFormTemplateVerifyErr.Error(), c)
		return
	}
	err := server.ReNameReimbursementFormTemplate(rn)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("重命名失败，%v", err), c)
	} else {
		response.OkWithMessage("重命名成功", c)
	}
}
