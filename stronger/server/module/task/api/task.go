package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"server/global"
	_func "server/global/func"
	"server/global/response"
	"server/module/load/model"
	"server/module/load/service"
	"server/module/task"
	"server/module/task/project"
	"server/module/task/request"
	responseOp "server/module/task/response"
	taskService "server/module/task/service"
	"server/utils"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RefreshProConf
// @Tags task
// @Summary 重启管理--刷新配置
// @Auth xingqiyi
// @Date 2022/2/16 3:42 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/conf/refresh-pro-conf [post]
func RefreshProConf(c *gin.Context) {
	task.Init()
	response.OkDetailed(make(map[string]interface{}), "刷新成功", c)
}

// GetProConf
// @Tags task
// @Summary 录入管理--获取项目配置
// @Auth xingqiyi
// @Date 2022/2/16 3:41 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/conf [get]
func GetProConf(c *gin.Context) {
	response.OkDetailed(task.CacheProConf, "领取成功", c)
}

// GetOpNum
// @Tags task
// @Summary 录入管理--获取任务数量
// @Auth xingqiyi
// @Date 2022/2/16 3:39 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/opNum [get]
func GetOpNum(c *gin.Context) {
	var taskGet request.TaskGet
	_ = c.ShouldBindQuery(&taskGet)
	code := taskGet.Code
	fmt.Println("11111111111111111111:", code)
	data := task.TaskListOpNum
	data.Accuracy = 1
	for _, userDayOutput := range task.UserDayOutputs {
		if userDayOutput.Code == code {
			data.Character = userDayOutput.SummaryFieldCharacter
			data.Accuracy = userDayOutput.SummaryAccuracyRate
			break
		}
	}
	response.OkDetailed(data, "查询成功", c)
}

func ReleaseBill(c *gin.Context) {
	var billRelease request.BillRelease
	_ = c.ShouldBindJSON(&billRelease)
	id := billRelease.ID
	if id != "" {
		task.BroadcastRelease(id, "", "", "")
	}
	response.OkDetailed(map[string]string{}, "操作成功", c)
}

// ReleaseExitBlock
// @Tags task/(释放领取分块)
// @Summary 退出释放领取分块
// @Auth sf
// @Date 2023/10/7  下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param code        query   string   true    "工号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"操作成功"}"
// @Router /task/releaseExitBlock [post]
func ReleaseExitBlock(c *gin.Context) {
	var taskRelease request.TaskRelease
	_ = c.ShouldBindJSON(&taskRelease)
	code := taskRelease.Code
	ops := []string{"op0", "op1", "op2", "opq"}
	for _, op := range ops {
		err, block := service.GetUserTaskBlock(task.TaskProCode, op, code)
		fmt.Println("---------已领取手上的分块------------:", err)
		if err == nil {
			err = ReleaseFunc(block, "", op, false)
			if err != nil {
				response.FailWithMessage(fmt.Sprintf("操作失败"), c)
			}
			// fmt.Println("11111111111111111111:", code)
		}
	}
	response.OkDetailed(map[string]string{}, "操作成功", c)
}

func ReleaseFunc(block model.ProjectBlock, code, op string, isBroad bool) error {
	data := make(map[string]interface{})
	mblock := reflect.ValueOf(block)
	mOp := strings.Replace(op, "o", "O", -1)
	data[op+"_code"] = code
	clearop := "op1"
	if op == "op1" {
		clearop = "op2"
	}
	if code != "" {
		data[op+"_stage"] = op + "Cache"
		data[op+"_apply_at"] = time.Now()
		if block.IsCompetitive && (op == "op1" || op == "op2") {
			data[clearop+"_code"] = "0"
			data[clearop+"_apply_at"] = time.Now()
			data[clearop+"_submit_at"] = time.Now()
			data[clearop+"_stage"] = "done"
		}
	} else {
		data[op+"_stage"] = op
		data[op+"_apply_at"] = time.Time{}
		if block.IsCompetitive && (op == "op1" || op == "op2") {
			data[clearop+"_code"] = ""
			data[clearop+"_apply_at"] = time.Time{}
			data[clearop+"_submit_at"] = time.Time{}
			data[clearop+"_stage"] = clearop
		}
	}
	if isBroad {
		opCode := mblock.FieldByName(mOp + "Code").Interface().(string)
		task.BroadcastRelease("", block.ID, op, opCode)
	}
	return service.UpdateBlock(task.TaskProCode, data, block.ID, op)
}

func ReleaseBlock(c *gin.Context) {
	var taskRelease request.TaskRelease
	_ = c.ShouldBindJSON(&taskRelease)
	code := taskRelease.Code
	op := taskRelease.Op
	id := taskRelease.ID
	fmt.Println("------------------------ReleaseBlock-----------------------:", op, code, id)
	err, block := service.SelectBlockByID(task.TaskProCode, id)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("操作失败"), c)
		return
	}
	// data := make(map[string]interface{})
	// data[op+"_code"] = code
	// clearop := "op1"
	// if op == "op1" {
	// 	clearop = "op2"
	// }
	// if code != "" {
	// 	data[op+"_stage"] = op + "Cache"
	// 	data[op+"_apply_at"] = time.Now()
	// 	if block.IsCompetitive && (op == "op1" || op == "op2") {
	// 		data[clearop+"_code"] = "0"
	// 		data[clearop+"_apply_at"] = time.Now()
	// 		data[clearop+"_submit_at"] = time.Now()
	// 		data[clearop+"_stage"] = "done"
	// 	}
	// } else {
	// 	data[op+"_stage"] = op
	// 	data[op+"_apply_at"] = time.Time{}
	// 	if block.IsCompetitive && (op == "op1" || op == "op2") {
	// 		data[clearop+"_code"] = ""
	// 		data[clearop+"_apply_at"] = time.Time{}
	// 		data[clearop+"_submit_at"] = time.Time{}
	// 		data[clearop+"_stage"] = clearop
	// 	}
	// }
	// task.BroadcastRelease("", block.ID, op)
	// err = service.UpdateBlock(task.TaskProCode, data, block.ID, op)
	err = ReleaseFunc(block, code, op, true)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("操作失败"), c)
		return
	}
	// fmt.Println("11111111111111111111:", code)
	response.OkDetailed(map[string]string{}, "操作成功", c)
}

func GetOpModifyBlock(c *gin.Context) {
	var taskGet request.TaskGet
	_ = c.ShouldBindQuery(&taskGet)
	op := taskGet.Op
	code := taskGet.Code
	num := taskGet.Num
	fmt.Println("11111111111111111111:", op, code)
	err, block := service.GetUserModifyTaskBlock(task.TaskProCode, op, code, num)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.FailWithMessage(fmt.Sprintf("没有可修改分块"), c)
		return
	}
	err, data := getBillData(block, code, op)
	if err != nil {
		fmt.Println("---------领取修改分块失败------------:", err)
		response.FailWithMessage(fmt.Sprintf("领取失败"), c)
		return
	} else {
		err, sysProject := service.GetSysProject(global.GConfig.System.ProCode)
		// NextTime = float64(sysProject.CacheTime)
		// err = formatImages(&data)
		data.CacheTime = sysProject.CacheTime
		global.GLog.Error("formatImages", zap.Error(err))
		response.OkDetailed(data, "领取成功", c)
		return
	}
}

// GetOpTask
// @Tags task
// @Summary 录入管理--领取任务
// @Auth xingqiyi
// @Date 2022/2/16 3:37 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/op [get]
func GetOpTask(c *gin.Context) {
	var taskGet request.TaskGet
	_ = c.ShouldBindQuery(&taskGet)
	op := taskGet.Op
	code := taskGet.Code
	fmt.Println("11111111111111111111:", op, code)
	err, block := service.GetUserTaskBlock(task.TaskProCode, op, code)
	fmt.Println("---------已领取手上的分块------------:", err)
	if err == nil {
		err, data := getBillData(block, code, op)
		fmt.Println("---------getBillData------------:", err)
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("领取失败"), c)
			return
		}
		//err = formatImages(&data)
		global.GLog.Error("formatImages", zap.Error(err))
		//response.OkDetailedEncrypt(data, "领取成功", c)
		response.OkDetailed(data, "领取成功", c)
		return
	}
	// err, data := GetTaskBlock(task.TaskProCode, op, code, []string{}, false)
	// if err != nil {
	// 	response.FailWithMessage(fmt.Sprintf("没有分块"), c)
	// } else {
	// 	response.OkDetailed(data, "领取成功", c)
	// }
	// return
	for true {
		err, block = taskService.GetTaskBlock(task.TaskProCode, op, code)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("---------查询可领取新分块失败------------:", err)
			response.FailWithMessage(fmt.Sprintf("没有分块"), c)
			return
		} else if err != nil {
			fmt.Println("---------查询可领取异常------------:", err)
			response.FailWithMessage(fmt.Sprintf("领取异常"), c)
			return
		}
		err, data := getBillData(block, code, op)
		if err != nil {
			fmt.Println("---------领取新分块失败------------:", err)
			// response.FailWithMessage(fmt.Sprintf("领取失败"), c)
		} else {
			//err = formatImages(&data)
			//global.GLog.Error("formatImages", zap.Error(err))
			response.OkDetailed(data, "领取成功", c)
			//response.OkDetailedEncrypt(data, "领取成功", c)
			return
		}

	}

}

// func GetTaskBlock(proCode string, op string, code string, blackList []string, isCompetitive bool) (error, responseOp.OpData) {
// 	for true {
// 		err, block := service.GetTaskBlock(task.TaskProCode, op, []string{}, isCompetitive)
// 		if err != nil {
// 			fmt.Println("---------查询可领取新分块失败------------:", err)
// 			// response.FailWithMessage(fmt.Sprintf("没有分块"), c)
// 			if op == "op2" {
// 				return GetTaskBlock(task.TaskProCode, "op1", code, []string{}, true)
// 			}
// 			return err, responseOp.OpData{}
// 		}
// 		err, data := getBillData(block, code, op)
// 		if err != nil {
// 			fmt.Println("---------领取新分块失败------------:", err)
// 			// response.FailWithMessage(fmt.Sprintf("领取失败"), c)
// 		} else {
// 			// response.OkDetailed(data, "领取成功", c)
// 			return nil, data
// 		}

// 	}
// }

// SubmitOpTask
// @Tags casbin
// @Summary 录入管理--提交任务
// @Auth xingqiyi
// @Date 2022/2/16 3:40 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/submit [post]
func SubmitOpTask(c *gin.Context) {
	var taskSubmit request.TaskSubmit
	_ = c.ShouldBindJSON(&taskSubmit)
	fmt.Println("Bill:", taskSubmit.Bill)
	fmt.Println("Fields:", taskSubmit.Fields)
	err := taskService.SaveSubmitData(task.TaskProCode, taskSubmit)
	if err != nil {
		fmt.Println("SubmitOpTask  err:", err)
		response.FailWithMessage(fmt.Sprintf("更新失败"), c)
	} else {
		response.OkDetailed(responseOp.OpData{}, "更新成功", c)
	}
}

// KeyLog
// @Tags casbin
// @Summary 录入管理--按键日志
// @Auth xingqiyi
// @Date 2022/2/16 3:40 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param billNum        body   string   true    "单号"
// @Param log       	 body   string   true    "按键操作日志"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/keyLog [post]
func KeyLog(c *gin.Context) {
	var keyLog request.KeyLog
	_ = c.ShouldBindJSON(&keyLog)
	fmt.Println("-------keyLog------------:", keyLog)
	fileDir := global.GConfig.LocalUpload.FilePath + global.GConfig.System.ProCode + "/key_log/" +
		time.Now().Format("2006/01/02")
	exists, err := utils.PathExists(fileDir)
	if err != nil {
		response.OkDetailed(responseOp.OpData{}, "更新失败", c)
		return
	}
	if !exists {
		err = utils.CreateDir(fileDir)
		if err != nil {
			response.OkDetailed(responseOp.OpData{}, "更新失败", c)
			return
		}
	}
	logs := ""
	for _, log := range keyLog.Log {
		logs += "\r\n" + log
	}
	err = ioutil.WriteFile(fileDir+"/"+keyLog.BillNum+".txt", []byte(logs), 0666)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败"), c)
	} else {
		response.OkDetailed(responseOp.OpData{}, "更新成功", c)
	}
}

func getBillData(block model.ProjectBlock, code string, op string) (error, responseOp.OpData) {
	// if isUpdate {
	mblock := reflect.ValueOf(&block).Elem()
	mOp := strings.Replace(op, "o", "O", -1)
	mblock.FieldByName(mOp + "Code").Set(reflect.ValueOf(code))
	mblock.FieldByName(mOp + "ApplyAt").Set(reflect.ValueOf(time.Now()))
	mblock.FieldByName(mOp + "Stage").Set(reflect.ValueOf(op + "Cache"))
	// block.Stage = op + "Cache"
	data := make(map[string]interface{})
	data[op+"_code"] = code
	data[op+"_apply_at"] = time.Now()
	data[op+"_stage"] = op + "Cache"
	if block.IsCompetitive && (op == "op1" || op == "op2") {
		clearop := "op1"
		if op == "op1" {
			clearop = "op2"
		}
		data[clearop+"_code"] = "0"
		data[clearop+"_apply_at"] = time.Now()
		data[clearop+"_submit_at"] = time.Now()
		data[clearop+"_stage"] = "done"
	}
	err := service.UpdateBlock(task.TaskProCode, data, block.ID, op)
	// return err, responseOp.OpData{}
	// }
	if err != nil {
		return err, responseOp.OpData{}
	}
	err, bill := service.SelectBillByID(task.TaskProCode, block.BillID)
	if err != nil {
		return err, responseOp.OpData{}
	}
	err, fields := service.SelectOpFieldsByBlockID(task.TaskProCode, block.ID)
	opFields := _func.FieldsFormat(fields)
	codeValues := project.BillValue(strings.Replace(task.TaskProCode, "_task", "", -1), block, fields, op)
	return err, responseOp.OpData{
		Bill:       bill,
		Block:      block,
		Fields:     opFields,
		CodeValues: codeValues,
	}
}

// UploadImage
// @Tags task
// @Summary 录入管理--上传编辑的图片
// @Auth xingqiyi
// @Date 2022/2/16 3:40 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/uploadImage [post]
func UploadImage(c *gin.Context) {
	form, _ := c.MultipartForm()
	file := form.File["file"][0]
	path := form.Value["path"][0]
	name := form.Value["name"][0]
	// sysProTempId := c.PostForm("sysProTempId")
	// pathArr := make([]string, 0)
	// sysProTempBIdpath := ""
	// for _, file := range files {
	err, sysProTempBIdpath := utils.SaveImgFile(c, file, name, global.GConfig.LocalUpload.FilePath+path)
	fmt.Println("err:", err)
	fmt.Println("sysProTempBIdpath:", sysProTempBIdpath)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	}
	// pathArr = append(pathArr, sysProTempBIdpath)
	// }
	// reRow := service2.UpdateImages(pathArr, sysProTempId)
	// if reRow != 1 {
	// response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	// } else {
	response.OkDetailed(name, "上传成功", c)
	// }
}

// GetImageSzie
// @Tags task
// @Summary 录入管理--查看图片
// @Auth sf
// @Date 2022/2/16 3:40 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param path    query   string   true    "图片路径"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /task/getImageSize [get]
func GetImageSzie(c *gin.Context) {
	var imageData request.ImageData
	_ = c.ShouldBindQuery(&imageData)
	url := global.GConfig.LocalUpload.FilePath + imageData.Path
	fmt.Println("url:", url)
	cmd := fmt.Sprintf(`ls -sh %s`, url)
	fmt.Println("cmd:", cmd)
	err, stdout, _ := utils.ShellOut(cmd)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	}
	size := strings.Split(stdout, " ")[0]
	sizeFloat := 0.0
	if strings.Index(size, "K") != -1 {
		size = RegReplace(size, `K`, "")
		sizeFloat, _ = strconv.ParseFloat(size, 64)
	} else {
		size = RegReplace(size, `M`, "")
		sizeFloat, _ = strconv.ParseFloat(size, 64)
		sizeFloat = sizeFloat * 1024
	}

	response.OkDetailed(sizeFloat, "查询成功", c)

}

func RegReplace(data string, query string, value string) string {
	reg := regexp.MustCompile(query)
	return reg.ReplaceAllString(data, value)
}

// 读取图片信息
func formatImages(bill *responseOp.OpData) error {
	bill.Bins = responseOp.BillAndBlockBin{}
	//byteArr := make([][]byte, 0)
	//for _, picture := range bill.Bill.Pictures {
	//	url := global.GConfig.LocalUpload.FilePath + bill.Bill.DownloadPath + picture
	//	global.GLog.Info("bill:::" + url)
	//	contentByte, err := ioutil.ReadFile(url)
	//	if err != nil {
	//		return err
	//	}
	//	byteArr = append(byteArr, contentByte)
	//}
	//bill.Bins.BillBin = byteArr
	url := global.GConfig.LocalUpload.FilePath + bill.Bill.DownloadPath + bill.Block.Picture
	global.GLog.Info("block:::" + url)
	contentByte, err := ioutil.ReadFile(url)
	if err != nil {
		return err
	}
	bill.Bins.BlockBin = contentByte
	return nil
}

// GetThumbnail
// @Tags task
// @Summary 录入管理--获取初审缩列图
// @Auth xingqiyi
// @Date 2022/4/6 16:53
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param billId query string true "单据id"
// @Param proCode query string true "项目编码"
// @Param pageIndex query int true "页码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/get-thumbnail [get]
func GetThumbnail(c *gin.Context) {
	var taskImage request.TaskImageThumbnail
	_ = c.ShouldBindQuery(&taskImage)
	//fmt.Println(taskImage)
	err, bill := service.SelectBillByID(taskImage.ProCode+"_task", taskImage.BillId)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取单据失败，%v", err), c)
	}
	pictures := bill.Pictures
	maxIndex := taskImage.PageIndex * task.ImagePageSize
	minIndex := (taskImage.PageIndex - 1) * task.ImagePageSize
	if minIndex > len(pictures) {
		response.FailWithMessage(fmt.Sprintf("没有那么图片啦，%v", err), c)
		return
	}
	if len(pictures) > maxIndex {
		pictures = pictures[minIndex:maxIndex]
	} else {
		pictures = pictures[minIndex:len(pictures)]
	}
	byteArr := make([][]byte, 0)
	for _, picture := range pictures {
		url := global.GConfig.LocalUpload.FilePath + bill.DownloadPath + picture
		contentByte, err := ioutil.ReadFile(url)
		if err != nil {
			global.GLog.Info("err:::" + err.Error())
		}
		byteArr = append(byteArr, contentByte)
	}
	response.OkDetailed(byteArr, "获取成功", c)
}

// GetImageByIndex
// @Tags task
// @Summary 录入管理--获取初审原图
// @Auth xingqiyi
// @Date 2022/4/6 20:33
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param billId query string true "单据id"
// @Param proCode query string true "项目编码"
// @Param index query int true "第几张图片"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/get-image-by-index [get]
func GetImageByIndex(c *gin.Context) {
	var taskImage request.TaskImage
	_ = c.ShouldBindQuery(&taskImage)
	//fmt.Println(taskImage)
	err, bill := service.SelectBillByID(taskImage.ProCode+"_task", taskImage.BillId)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取单据失败，%v", err), c)
		return
	}
	if len(bill.Pictures) < (taskImage.Index + 1) {
		response.FailWithMessage(fmt.Sprintf("没有那么多图片，%v", err), c)
		return
	}
	url := global.GConfig.LocalUpload.FilePath + bill.DownloadPath + bill.Pictures[taskImage.Index]
	contentByte, err := ioutil.ReadFile(url)
	if err != nil {
		global.GLog.Info("err:::" + err.Error())
		response.FailWithMessage(fmt.Sprintf("读取图片失败，%v", err), c)
		return
	}
	response.OkDetailed(contentByte, "获取成功", c)
}

// GetImageByBlockId
// @Tags task
// @Summary 录入管理--获取分块图片
// @Auth xingqiyi
// @Date 2022/5/5 09:49
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param blockId query string true "分块id"
// @Param proCode query string true "项目编码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/get-image-by-block-id [get]
func GetImageByBlockId(c *gin.Context) {
	var taskImage request.TaskImage
	_ = c.ShouldBindQuery(&taskImage)
	err, block := service.SelectBlockByID(taskImage.ProCode+"_task", taskImage.BillId)
	err, bill := service.SelectBillByID(taskImage.ProCode+"_task", taskImage.BillId)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取分块失败，%v", err), c)
		return
	}
	url := global.GConfig.LocalUpload.FilePath + bill.DownloadPath + block.Picture
	contentByte, err := ioutil.ReadFile(url)
	if err != nil {
		global.GLog.Info("err:::" + err.Error())
		response.FailWithMessage(fmt.Sprintf("读取图片失败，%v", err), c)
		return
	}
	response.OkDetailed(contentByte, "获取成功", c)
}
