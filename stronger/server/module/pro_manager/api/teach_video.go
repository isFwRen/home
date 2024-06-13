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
	scheduleService "server/module/training_guide/service"
	"server/utils"
	"strconv"
	"strings"
	api2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"
)

// BlockSync
// @Tags Teach Video (教学视频)
// @Summary 教学视频--同步数据
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/teachVideo/sync [get]
func BlockSync(c *gin.Context) {
	proCode := c.Query("proCode")
	err := service.BlockSync(proCode)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("同步数据失败, err : %v", err), c)
	} else {
		response.Ok(c)
	}
}

// GetTaskTeachVideoList
// @Tags Teach Video (教学视频(录入系统))
// @Summary 教学视频(录入系统)--查询
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param blockName query string false "分块名称"
// @Param rule query string false "教学视频"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/teachVideo/task/list [get]
func GetTaskTeachVideoList(c *gin.Context) {
	var search request.TV
	if err := c.ShouldBindWith(&search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetTeachVideoListVerify := utils.Rules{
		"proCode": {utils.NotEmpty()},
	}
	GetTeachVideoListVerifyErr := utils.Verify(search, GetTeachVideoListVerify)
	if GetTeachVideoListVerifyErr != nil {
		response.FailWithMessage(GetTeachVideoListVerifyErr.Error(), c)
		return
	}
	err, list, total := service.GetTeachVideoList(search)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("查询教学视频失败, err : %v", err), c)
	} else {

		user, err1 := api2.GetUserByToken(c)
		if err1 != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		needReadList, err2 := scheduleService.GetNeedReadList2(3, user.ID)
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
			response.GetOkWithData(resp.TVReq{List: list, Total: total}, c)
			return
		}
		response.GetOkWithData(resp.TVReq{List: convMap, Total: total}, c)
	}
}

// GetTeachVideoList
// @Tags Teach Video (教学视频)
// @Summary 教学视频--查询
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param blockName query string false "分块名称"
// @Param rule query string false "教学视频"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/teachVideo/list [get]
func GetTeachVideoList(c *gin.Context) {
	var search request.TV
	if err := c.ShouldBindWith(&search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetTeachVideoListVerify := utils.Rules{
		"proCode": {utils.NotEmpty()},
	}
	GetTeachVideoListVerifyErr := utils.Verify(search, GetTeachVideoListVerify)
	if GetTeachVideoListVerifyErr != nil {
		response.FailWithMessage(GetTeachVideoListVerifyErr.Error(), c)
		return
	}
	err, list, total := service.GetTeachVideoListCM(search)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("查询教学视频失败, err : %v", err), c)
	} else {
		response.GetOkWithData(resp.TVReq{List: list, Total: total}, c)
	}
}

// EditTeachVideo
// @Tags Teach Video (教学视频)
// @Summary 教学视频--编辑
// @accept application/json
// @Produce application/json
// @Param id formData string true "id"
// @Param file formData file true "文件"
// @Param isRequired formData int true "是否必学"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/teachVideo/edit [post]
func EditTeachVideo(c *gin.Context) {
	form, _ := c.MultipartForm()
	id := form.Value["id"]
	isRequired := form.Value["isRequired"]
	if len(id) == 0 {
		response.GetFailWithData(fmt.Sprintf("id不能为空"), c)
		return
	}

	// 获取上传文件
	Files := form.File["file"]

	var tv model2.TeachVideo
	err := global.GDb.Model(&model2.TeachVideo{}).Where("id in ? ", id).Find(&tv).Error
	if err != nil {
		fmt.Println(err.Error())
		response.GetFailWithData(fmt.Sprintf("编辑失败"), c)
		return
	}

	var proInformation model.SysProject
	err = global.GDb.Model(&model.SysProject{}).Where("id = ? ", tv.ProId).Find(&proInformation).Error
	if err != nil {
		fmt.Println(err.Error())
		response.GetFailWithData(fmt.Sprintf("编辑失败"), c)
		return
	}
	if len(Files) == 0 {
		response.GetFailWithData(fmt.Sprintf("编辑失败"), c)
		return
	}
	for _, file := range Files {
		filename := file.Filename
		if tv.SysBlockName != strings.Split(filename, ".")[0] {
			response.GetFailWithData(fmt.Sprintf("编辑失败"), c)
			fmt.Println("***编辑失败***")
			fmt.Println("***编辑失败***", tv.SysBlockName, strings.Split(filename, ".")[0])
			return
		}
		// 线上 uploads/file
		basicPath := global.GConfig.LocalUpload.FilePath + "/教学视频/" + proInformation.Code + "/"
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
		tv.Video = pq.StringArray{}
		tv.Video = append(tv.Video, dst)
	}
	tv.InputRule = "有"
	tv.IsRequired, _ = strconv.Atoi(isRequired[0])
	err = service.EditTeachVideo(tv)
	if err != nil {
		fmt.Println(err)
		response.GetFailWithData(fmt.Sprintf("编辑失败"), c)
	} else {
		response.Ok(c)
	}
}

// DeleteTeachVideo
// @Tags Teach Video (教学视频)
// @Summary 教学视频--删除
// @accept application/json
// @Produce application/json
// @Param data body re.RmById true "id"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/teachVideo/delete [post]
func DeleteTeachVideo(c *gin.Context) {
	var reqId re.RmById
	_ = c.ShouldBindJSON(&reqId)
	IdVerifyErr := utils.Verify(reqId, utils.CustomizeMap["IdVerify"])
	if IdVerifyErr != nil {
		response.FailWithMessage(IdVerifyErr.Error(), c)
		return
	}
	err := service.DeleteTeachVideo(reqId.Ids)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", err), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// ExportTeachVideo
// @Tags Teach Video (教学视频)
// @Summary 教学视频--导出
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/teachVideo/export [get]
func ExportTeachVideo(c *gin.Context) {
	var req request.ExportTV
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
	err, p := service.ExportTeachVideo(req)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("导出失败，%v", err), c)
	} else {
		response.GetOkWithData(resp.FRReq{List: p}, c)
	}
}

// UploadTeachVideo
// @Tags Teach Video (教学视频)
// @Summary 教学视频--导入
// @accept application/json
// @Produce application/json
// @Param proCode body string true "项目编码"
// @Param file formData file true "文件"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/teachVideo/upload [post]
func UploadTeachVideo(c *gin.Context) {
	form, _ := c.MultipartForm()
	proCode := form.Value["proCode"]
	if len(proCode) == 0 {
		response.GetFailWithData(fmt.Sprintf("proCode不能为空"), c)
		return
	}

	// 获取上传文件
	Files := form.File["file"]

	err := service.UploadTeachVideo(Files, proCode, c)
	if err != nil {
		fmt.Println(err)
		response.FailWithMessage(fmt.Sprintf("导入失败, %s", err.Error()), c)
	} else {
		response.Ok(c)
	}
}
