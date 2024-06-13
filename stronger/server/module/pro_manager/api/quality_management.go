package api

import "C"
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"os"
	"path"
	"server/global"
	re "server/global/request"
	"server/global/response"
	"server/module/pro_manager/model/request"
	resp "server/module/pro_manager/model/response"
	server "server/module/pro_manager/service"
	"server/utils"
	"strings"
)

func Creatable(c *gin.Context) {
	server.CreateQualityTable()
}

// GetQualityManagement
// @Tags Quality Management (质量管理)
// @Summary 质量管理--查询
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param billName query string true "案件号"
// @Param wrongFieldName query string true "错误字段"
// @Param responsibleName query string true "责任人"
// @Param month query string true "日期YYYY-MM"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/quality/management/list [get]
func GetQualityManagement(c *gin.Context) {
	var search request.QualityRequest
	if err := c.ShouldBindWith(&search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetQualityListVerify := utils.Rules{
		"month": {utils.NotEmpty()},
	}
	GetQualityListVerifyErr := utils.Verify(search, GetQualityListVerify)
	if GetQualityListVerifyErr != nil {
		response.FailWithMessage(GetQualityListVerifyErr.Error(), c)
		return
	}
	err, list, total := server.GetQualityManagement(search)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("查询质量管理失败, err : %v", err), c)
	} else {
		response.GetOkWithData(resp.QualityRes{List: list, Total: total}, c)
	}
}

// AddQualityManagement
// @Tags Quality Management (质量管理)
// @Summary 质量管理--新增
// @accept application/json
// @Produce application/json
// @Param data body request.QualitiesAeRequest true "新增"
// @Param file formData file true "文件"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/quality/management/add [post]
func AddQualityManagement(c *gin.Context) {
	var add request.QualitiesAeRequest
	_ = c.ShouldBind(&add)

	QualitiesAddVerify := utils.Rules{
		"ProCode": {utils.NotEmpty()},
	}
	QualitiesAddVerifyErr := utils.Verify(add, QualitiesAddVerify)
	if QualitiesAddVerifyErr != nil {
		response.FailWithMessage(QualitiesAddVerifyErr.Error(), c)
		return
	}

	form, _ := c.MultipartForm()
	// 获取上传文件
	excelFiles := form.File["file"]
	for _, excelFile := range excelFiles {
		if excelFile.Size > 5*1048576 {
			response.FailWithMessage("文件大小要限制在5MB之内", c)
			return
		}
		filename := add.BillName + "-" + add.WrongFieldName + ".png"
		// 线上 uploads/file
		basicPath := global.GConfig.LocalUpload.FilePath + "/quality-management/add/" + add.FeedbackDate + "/" + add.BillName + "/"
		// 本地
		//basicPath := "./Salary/" + filename[:4] + "/"
		fmt.Println("The path of saving excels is " + basicPath)
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
		saveErr := c.SaveUploadedFile(excelFile, dst)
		if saveErr != nil {
			fmt.Println("saveErr", saveErr)
			continue
		}
		add.ImagePath = append(add.ImagePath, basicPath+filename)
	}

	err := server.AddQualityManagement(add)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("新增质量管理失败, err : %v", err), c)
	} else {
		response.Ok(c)
	}
}

// EditQualityManagement
// @Tags Quality Management (质量管理)
// @Summary 质量管理--修改
// @accept application/json
// @Produce application/json
// @Param data body request.QualitiesAeRequest true "修改"
// @Param file formData file true "文件"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/quality/management/edit [post]
func EditQualityManagement(c *gin.Context) {
	var edit request.QualitiesAeRequest
	_ = c.ShouldBind(&edit)
	QualitiesAddVerify := utils.Rules{
		"ProCode": {utils.NotEmpty()},
	}
	QualitiesAddVerifyErr := utils.Verify(edit, QualitiesAddVerify)
	if QualitiesAddVerifyErr != nil {
		response.FailWithMessage(QualitiesAddVerifyErr.Error(), c)
		return
	}

	form, _ := c.MultipartForm()
	// 获取上传文件
	excelFiles := form.File["file"]
	for _, excelFile := range excelFiles {
		if excelFile.Size > 5*1048576 {
			response.FailWithMessage("文件大小要限制在5MB之内", c)
			return
		}
		filename := edit.BillName + "-" + edit.WrongFieldName + ".png"
		// 线上 uploads/file
		basicPath := global.GConfig.LocalUpload.FilePath + "/quality-management/add/" + edit.FeedbackDate + "/" + edit.BillName + "/"
		// 本地
		//basicPath := "./Salary/" + filename[:4] + "/"
		fmt.Println("The path of saving excels is " + basicPath)
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
		saveErr := c.SaveUploadedFile(excelFile, dst)
		if saveErr != nil {
			fmt.Println("saveErr", saveErr)
			continue
		}
		edit.ImagePath = append(edit.ImagePath, dst)
	}

	err := server.EditQualityManagement(edit)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("修改质量管理失败, err : %v", err), c)
	} else {
		response.Ok(c)
	}
}

// UploadQualityManagement
// @Tags Quality Management (质量管理)
// @Summary	质量管理--上传质量管理
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param file formData file true "文件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /pro-manager/quality/management/push [post]
func UploadQualityManagement(c *gin.Context) {
	form, _ := c.MultipartForm()
	// 获取上传文件
	excelFiles := form.File["file"]
	fmt.Println(len(excelFiles))
	// 线上 uploads/file
	basicPath := global.GConfig.LocalUpload.FilePath + "质量管理/" + "upload/"
	// 本地
	//basicPath := "D:/Projects/tempfiles/"
	fmt.Println("The path of saving excels is " + basicPath)
	excelArr := make(map[string]string, 0)
	failArray := make([]string, 0)
	for _, excelFile := range excelFiles {
		filename := excelFile.Filename
		f := strings.Replace(filename, ".xlsx", "", -1)
		// 判断文件是否符合要求
		dst1 := path.Join(basicPath, f)
		dst := path.Join(dst1, filename)
		fmt.Println(filename + "has been saved in this path " + basicPath)
		isExist, err := exists(dst1)
		if err != nil {
			failArray = append(failArray, err.Error())
			continue
		}
		if !isExist {
			err = os.MkdirAll(dst1, 0777)
			if err != nil {
				failArray = append(failArray, err.Error())
				continue
			}
			err = os.Chmod(dst1, 0777)
			if err != nil {
				failArray = append(failArray, err.Error())
				continue
			}
		}
		// 上传文件到指定的路径
		fmt.Println("dst", dst)
		saveErr := c.SaveUploadedFile(excelFile, dst)
		if saveErr != nil {
			fmt.Println("saveErr", saveErr)
			failArray = append(failArray, filename+", 保存失败!")
			continue
		}
		excelArr[filename] = dst1
	}

	if len(failArray) != 0 {
		response.FailWithDetailed(response.ERROR, failArray, "上传质量管理表失败, 请重新上传!", c)
	} else {
		//保存数据
		//uid := c.Request.Header.Get("x-user-id")
		//if uid == "" {
		//	response.FailWithMessage("x-user-id is empty", c)
		//	return
		//}
		//err, a := utils.GetRedisExport("PutQualityWithExcel", uid)
		//if err == nil && a == "true" {
		//	response.FailWithMessage("正在导入!", c)
		//	return
		//}
		//go func() {
		//	defer func() {
		//		if r := recover(); r != nil {
		//			fmt.Println(fmt.Sprintf("质量管理--上传质量管理表导入失败, err: %s", r))
		//			global.GLog.Error(fmt.Sprintf("质量管理--上传质量管理表导入失败, err: %s", r))
		//			err = utils.DelRedisExport("PutQualityWithExcel", uid)
		//			if err != nil {
		//				fmt.Println(fmt.Sprintf("质量管理--上传质量管理表删除导出缓存失败, err: %s", err.Error()))
		//				global.GLog.Error(fmt.Sprintf("质量管理--上传质量管理表删除导出缓存失败, err: %s", err.Error()))
		//			}
		//		}
		//	}()
		//	err := utils.SetRedisExport("PutQualityWithExcel", uid)
		//	if err != nil {
		//		global.GSocketConnSendMsgMap[uid].Emit("qualityReply", response.ExportResponse{
		//			Code: 400,
		//			Data: "",
		//			Msg:  fmt.Sprintf("qualityFileUpload SetRedisExport err: %s", err.Error()),
		//		})
		//		err = utils.DelRedisExport("PutQualityWithExcel", uid)
		//		return
		//	}
		//	_ = server.UploadQualityManagement(excelArr)
		//	//if msg != "" {
		//	//	global.GSocketConnSendMsgMap[uid].Emit("constReply", response.ExportResponse{
		//	//		Code: 400,
		//	//		Data: "",
		//	//		Msg:  fmt.Sprintf("PtSalaryFileUpload SetRedisExport err: %s", msg),
		//	//	})
		//	//	err = utils.DelRedisExport("PutConstTableWithExcel", uid)
		//	//	return
		//	//}
		//	global.GSocketConnSendMsgMap[uid].Emit("qualityReply", response.ExportResponse{
		//		Code: 200,
		//		Data: "",
		//		Msg:  "上传质量管理表成功!",
		//	})
		//	err = utils.DelRedisExport("PutQualityWithExcel", uid)
		//	return
		//}()
		err := server.UploadQualityManagement(excelArr)
		if err != nil {
			response.GetFailWithData(fmt.Sprintf("上传质量管理失败, err : %v", err), c)
		} else {
			response.Ok(c)
		}
	}
}

// ExportQualityData
// @Tags Quality Management (质量管理)
// @Summary	质量管理--导出质量管理
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param month query string true "日期YYYY-MM"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /pro-manager/quality/management/export [get]
func ExportQualityData(c *gin.Context) {
	//proCode := c.Query("proCode")
	//if proCode == "" {
	//	response.FailWithMessage("proCode is empty", c)
	//	return
	//}
	var search request.QualitiesExport
	if err := c.ShouldBindWith(&search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	ExportQualityVerify := utils.Rules{
		"month": {utils.NotEmpty()},
	}
	ExportQualityVerifyErr := utils.Verify(search, ExportQualityVerify)
	if ExportQualityVerifyErr != nil {
		response.FailWithMessage(ExportQualityVerifyErr.Error(), c)
		return
	}

	//uid := c.Request.Header.Get("x-user-id")
	//if uid == "" {
	//	response.FailWithMessage("x-user-id is empty", c)
	//	return
	//}
	//err, a := utils.GetRedisExport("exportQuality", uid)
	//if err == nil && a == "true" {
	//	response.FailWithMessage("正在导出!", c)
	//	return
	//}
	//go func() {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			fmt.Println(fmt.Sprintf("质量管理--导出质量管理表导出失败, err: %s", r))
	//			global.GLog.Error(fmt.Sprintf("质量管理--导出质量管理表导出失败, err: %s", r))
	//			err = utils.DelRedisExport("exportQuality", uid)
	//			if err != nil {
	//				fmt.Println(fmt.Sprintf("质量管理--导出质量管理表删除导出缓存失败, err: %s", err.Error()))
	//				global.GLog.Error(fmt.Sprintf("质量管理--导出质量管理表删除导出缓存失败, err: %s", err.Error()))
	//			}
	//		}
	//	}()
	//	//可以广播同一个登录人的客户端的写法
	//	//global.GSocketIo.BroadcastToRoom("/global-export", uid, "outputStatistics", "所要发送的消息")
	//	err := utils.SetRedisExport("exportQuality", uid)
	//	if err != nil {
	//		global.GSocketConnSendMsgMap[uid].Emit("qualityReply", response.ExportResponse{
	//			Code: 400,
	//			Data: "",
	//			Msg:  fmt.Sprintf("exportQuality SetRedisExport err: %s", err.Error()),
	//		})
	//		err = utils.DelRedisExport("exportQuality", uid)
	//		return
	//	}
	//	err, p := server.ExportQualityData(strings.Split(search.ProCode, ","), search.StartTime, search.EndTime)
	//	fmt.Println("path", p)
	//	fmt.Println("uid", uid)
	//	if err != nil {
	//		global.GSocketConnSendMsgMap[uid].Emit("qualityReply", response.ExportResponse{
	//			Code: 400,
	//			Data: "",
	//			Msg:  fmt.Sprintf("exportQuality err: %s", err.Error()),
	//		})
	//		err = utils.DelRedisExport("exportQuality", uid)
	//		return
	//	}
	//	global.GSocketConnSendMsgMap[uid].Emit("qualityReply", response.ExportResponse{
	//		Code: 200,
	//		Data: p,
	//		Msg:  "导出完成!",
	//	})
	//	err = utils.DelRedisExport("exportQuality", uid)
	//	fmt.Println("导出完成!")
	//	return
	//}()

	err, p := server.ExportQualityData(strings.Split(search.ProCode, ","), search.Month)
	fmt.Println("path", p)
	fmt.Println("err", err)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("导出失败，%v", err), c)
	} else {
		response.GetOkWithData(resp.QualityRes{List: p}, c)
	}
}

// DeleteQualityData
// @Tags  Quality Management (质量管理)
// @Summary	质量管理--删除
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body re.RmById true "id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/quality/management/delete [post]
func DeleteQualityData(c *gin.Context) {
	var reqId re.RmById
	_ = c.ShouldBindJSON(&reqId)
	IdVerifyErr := utils.Verify(reqId, utils.CustomizeMap["IdVerify"])
	if IdVerifyErr != nil {
		response.FailWithMessage(IdVerifyErr.Error(), c)
		return
	}
	err := server.DeleteQualityData(reqId.Ids)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", err), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
