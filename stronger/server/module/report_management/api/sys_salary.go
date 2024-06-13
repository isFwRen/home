package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"os"
	"path"
	"regexp"
	"server/global"
	"server/global/response"
	"server/module/report_management/model/request"
	R "server/module/report_management/model/response"
	"server/module/report_management/service"
	"server/utils"
)

// GetPtSalaryTask
// @Tags Salary (工资数据(录入系统))
// @Summary	工资数据(录入系统)--查询工资数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param start query string true "开始年月"
// @Param end query string true "结束年月"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /report-management/pt/salary-task/list [get]
func GetPtSalaryTask(c *gin.Context) {
	var Search request.SysSalarySearch
	if err := c.ShouldBindWith(&Search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetPtSalaryVerify := utils.Rules{
		"start":     {utils.NotEmpty()},
		"end":       {utils.NotEmpty()},
		"PageIndex": {utils.Gt("0")},
		"PageSize":  {utils.Gt("0")},
	}
	GetPtSalaryVerifyErr := utils.Verify(Search, GetPtSalaryVerify)
	if GetPtSalaryVerifyErr != nil {
		response.FailWithMessage(GetPtSalaryVerifyErr.Error(), c)
		return
	}

	uid := c.Request.Header.Get("x-user-id")
	if uid == "" {
		response.FailWithMessage("x-user-id is empty", c)
		return
	}

	err, PtSalary, total := service.GetPtSalaryTask(Search, uid)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkWithData(R.SalaryRes{
			List:  PtSalary,
			Total: total,
		}, c)
	}
}

// GetPtSalary
// @Tags Salary (工资数据)
// @Summary	工资数据--查询工资数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param ym query string true "年月"
// @Param code query string false "工号"
// @Param name query string false "姓名"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /report-management/pt/salary/list [get]
func GetPtSalary(c *gin.Context) {
	var Search request.SysSalarySearch
	if err := c.ShouldBindWith(&Search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetPtSalaryVerify := utils.Rules{
		"Ym":        {utils.NotEmpty()},
		"PageIndex": {utils.Gt("0")},
		"PageSize":  {utils.Gt("0")},
	}
	GetPtSalaryVerifyErr := utils.Verify(Search, GetPtSalaryVerify)
	if GetPtSalaryVerifyErr != nil {
		response.FailWithMessage(GetPtSalaryVerifyErr.Error(), c)
		return
	}
	err, PtSalary, total := service.GetPtSalary(Search)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkWithData(R.SalaryRes{
			List:  PtSalary,
			Total: total,
		}, c)
	}
}

// PtSalaryFileDownload
// @Tags Salary (工资数据)
// @Summary	工资数据--导出PT工资数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ym query string true "年月"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /report-management/pt/salary/download [get]
func PtSalaryFileDownload(c *gin.Context) {
	FileName := c.Query("ym") + "PT工资表.xlsx"
	fmt.Println("FileName", FileName)
	// 线上 uploads/file
	basicPath := global.GConfig.LocalUpload.FilePath + "Salary/" + "pt/" + c.Query("ym")[:4] + "/"
	// 本地
	//basicPath := "./"
	fmt.Println("PtSalaryFileDownload path", basicPath+FileName)
	isExist, err := pathExists(basicPath + FileName)
	if err != nil {
		response.FailWithDetailed(response.ERROR, "", "查看是否有该月份的工资表失败, "+fmt.Sprintf("err : %s", err.Error()), c)
	}
	if isExist {
		c.FileAttachment(basicPath+FileName, FileName)
	} else {
		response.FailWithMessage("该月份的工资表不存在", c)
	}
}

// PtSalaryFileUpload
// @Tags Salary (工资数据)
// @Summary	工资数据--导入PT工资数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param file formData file true "文件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /report-management/pt/salary/upload [post]
func PtSalaryFileUpload(c *gin.Context) {
	form, _ := c.MultipartForm()
	// 获取上传文件
	excelFiles := form.File["file"]
	fmt.Println("excelFiles", len(excelFiles))
	failArray := make([]string, 0)
	fp := make([]string, 0)
	for _, excelFile := range excelFiles {
		filename := excelFile.Filename
		err := CheckSalaryFileName(filename, "pt")
		if err != nil {
			failArray = append(failArray, err.Error())
			break
		}
		fp = append(fp, filename)
		// 线上 uploads/file
		basicPath := global.GConfig.LocalUpload.FilePath + "/Salary/" + "pt/" + filename[:4] + "/"
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
			failArray = append(failArray, filename+"保存失败")
			continue
		}
	}
	fmt.Println("failArray", len(failArray))
	if len(failArray) != 0 {
		response.FailWithDetailed(response.ERROR, failArray, "保存工资数据表失败,请重新上传", c)
	} else {
		//msg := service.PtSalaryFileUpload(fp)
		//fmt.Println(msg)
		//保存数据
		uid := c.Request.Header.Get("x-user-id")
		if uid == "" {
			response.FailWithMessage("x-user-id is empty", c)
			return
		}
		err, a := utils.GetRedisExport("PtSalaryFileUpload", uid)
		if err == nil && a == "true" {
			response.FailWithMessage("正在导入!", c)
			return
		}
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(fmt.Sprintf("工资数据--导入PT工资数据导出失败, err: %s", r))
					global.GLog.Error(fmt.Sprintf("工资数据--导入PT工资数据导出失败, err: %s", r))
					err = utils.DelRedisExport("PtSalaryFileUpload", uid)
					if err != nil {
						fmt.Println(fmt.Sprintf("工资数据--导入PT工资数据删除导出缓存失败, err: %s", err.Error()))
						global.GLog.Error(fmt.Sprintf("工资数据--导入PT工资数据删除导出缓存失败, err: %s", err.Error()))
					}
				}
			}()
			err := utils.SetRedisExport("PtSalaryFileUpload", uid)
			if err != nil {
				//global.GSocketConnSendMsgMap[uid].Emit("global-send-message", response.ExportResponse{
				//	Code: 400,
				//	Data: "",
				//	Msg:  fmt.Sprintf("PtSalaryFileUpload SetRedisExport err: %s", err.Error()),
				//})
				err = utils.DelRedisExport("PtSalaryFileUpload", uid)
				return
			}
			msg := service.PtSalaryFileUpload(fp)
			if msg != "" {
				//global.GSocketConnSendMsgMap[uid].Emit("global-send-message", response.ExportResponse{
				//	Code: 400,
				//	Data: "",
				//	Msg:  fmt.Sprintf("PtSalaryFileUpload SetRedisExport err: %s", msg),
				//})
				err = utils.DelRedisExport("PtSalaryFileUpload", uid)
				return
			}
			//global.GSocketConnSendMsgMap[uid].Emit("global-send-message", response.ExportResponse{
			//	Code: 200,
			//	Data: "",
			//	Msg:  "导入完成!",
			//})
			err = utils.DelRedisExport("PtSalaryFileUpload", uid)
			return
		}()
		response.Ok(c)
	}
}

// GetInternalSalary
// @Tags Salary (工资数据)
// @Summary	工资数据--查询内部工资数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param ym formData string true "年月"
// @Param code formData string false "工号"
// @Param name formData string false "姓名"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /report-management/internal/salary/list [get]
func GetInternalSalary(c *gin.Context) {
	var Search request.SysSalarySearch
	if err := c.ShouldBindWith(&Search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetPtSalaryVerify := utils.Rules{
		"Ym":        {utils.NotEmpty()},
		"PageIndex": {utils.Gt("0")},
		"PageSize":  {utils.Gt("0")},
	}
	GetPtSalaryVerifyErr := utils.Verify(Search, GetPtSalaryVerify)
	if GetPtSalaryVerifyErr != nil {
		response.FailWithMessage(GetPtSalaryVerifyErr.Error(), c)
		return
	}

	uid := c.Request.Header.Get("x-user-id")
	if uid == "" {
		response.FailWithMessage("x-user-id is empty", c)
		return
	}

	err, InternalSalary, total := service.GetInternalSalary(Search, uid)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkWithData(R.SalaryRes{
			List:  InternalSalary,
			Total: total,
		}, c)
	}
}

// InternalSalaryFileDownload
// @Tags Salary (工资数据)
// @Summary	工资数据--导出内部工资数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ym formData string true "年月"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /report-management/internal/salary/download [get]
func InternalSalaryFileDownload(c *gin.Context) {
	FileName := c.Query("ym") + "内部工资表.xlsx"
	fmt.Println("FileName", FileName)
	// 线上 uploads/file
	basicPath := global.GConfig.LocalUpload.FilePath + "/Salary/" + "internal/" + c.Query("ym")[:4] + "/"
	// 本地
	//basicPath := "./"
	fmt.Println("InternalSalaryFileDownload path", basicPath+FileName)
	isExist, err := pathExists(basicPath + FileName)
	if err != nil {
		response.FailWithDetailed(response.ERROR, "", "查看是否有该月份的工资表失败, "+fmt.Sprintf("err : %s", err.Error()), c)
	}
	if isExist {
		c.FileAttachment(basicPath+FileName, FileName)
	} else {
		response.FailWithMessage("该月份的工资表不存在", c)
	}
}

// InternalSalaryFileUpload
// @Tags Salary (工资数据)
// @Summary	工资数据--导入内部工资数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param file formData file true "文件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /report-management/internal/salary/upload [post]
func InternalSalaryFileUpload(c *gin.Context) {
	form, _ := c.MultipartForm()
	// 获取上传文件
	excelFiles := form.File["file"]
	fmt.Println("excelFiles", len(excelFiles))
	failArray := make([]string, 0)
	fp := make([]string, 0)
	for _, excelFile := range excelFiles {
		filename := excelFile.Filename
		err := CheckSalaryFileName(filename, "")
		if err != nil {
			failArray = append(failArray, err.Error())
			break
		}
		fp = append(fp, filename)
		// 线上 uploads/file
		basicPath := global.GConfig.LocalUpload.FilePath + "/Salary/" + "internal/" + filename[:4] + "/"
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
			failArray = append(failArray, filename+"上传失败")
			continue
		}
	}
	fmt.Println("failArray", len(failArray))
	//_ = service.InternalSalaryFileUpload([]string{"202112交付内部工资模板.xlsx"})
	if len(failArray) != 0 {
		response.FailWithDetailed(response.ERROR, failArray, "上传工资数据表失败,请重新上传", c)
	} else {
		//msg := service.InternalSalaryFileUpload(fp)
		//fmt.Println(msg)
		//保存数据
		uid := c.Request.Header.Get("x-user-id")
		if uid == "" {
			response.FailWithMessage("x-user-id is empty", c)
			return
		}
		err, a := utils.GetRedisExport("InternalSalaryFileUpload", uid)
		if err == nil && a == "true" {
			response.FailWithMessage("正在导入!", c)
			return
		}
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(fmt.Sprintf("工资数据--导入内部工资数据导出失败, err: %s", r))
					global.GLog.Error(fmt.Sprintf("工资数据--导入内部工资数据导出失败, err: %s", r))
					err = utils.DelRedisExport("InternalSalaryFileUpload", uid)
					if err != nil {
						fmt.Println(fmt.Sprintf("工资数据--导入内部工资数据删除导出缓存失败, err: %s", err.Error()))
						global.GLog.Error(fmt.Sprintf("工资数据--导入内部工资数据删除导出缓存失败, err: %s", err.Error()))
					}
				}
			}()
			err := utils.SetRedisExport("InternalSalaryFileUpload", uid)
			if err != nil {
				//global.GSocketConnSendMsgMap[uid].Emit("global-send-message", response.ExportResponse{
				//	Code: 400,
				//	Data: "",
				//	Msg:  fmt.Sprintf("InternalSalaryFileUpload SetRedisExport err: %s", err.Error()),
				//})
				err = utils.DelRedisExport("InternalSalaryFileUpload", uid)
				return
			}
			msg := service.InternalSalaryFileUpload(fp)
			if msg != "" {
				//global.GSocketConnSendMsgMap[uid].Emit("global-send-message", response.ExportResponse{
				//	Code: 400,
				//	Data: "",
				//	Msg:  fmt.Sprintf("InternalSalaryFileUpload SetRedisExport err: %s", msg),
				//})
				err = utils.DelRedisExport("InternalSalaryFileUpload", uid)
				return
			}
			//global.GSocketConnSendMsgMap[uid].Emit("global-send-message", response.ExportResponse{
			//	Code: 200,
			//	Data: "",
			//	Msg:  "导入完成!",
			//})
			err = utils.DelRedisExport("InternalSalaryFileUpload", uid)
			return
		}()
		response.Ok(c)
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

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	fmt.Println(err)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CheckSalaryFileName
// PT工资表：年+月+PT工资表，比如202101PT工资表
// 内部工资表：年+月+内部工资表，比如202101内部工资表
func CheckSalaryFileName(filename string, types string) error {
	msg := "请按以下格式来命名: " + "\n" + "  1. PT工资表: 年+月+PT工资表, 比如202101PT工资表.xlsx" + "\n" + "  2. 内部工资表：年+月+内部工资表, 比如202101内部工资表.xlsx"
	fmt.Println(filename[6:])
	if types == "pt" {
		//length = len("202112PT工资表")
		if filename[6:] != "PT工资表.xlsx" {
			return errors.New(msg)
		}
	} else {
		//length = len("202112内部工资表")
		if filename[6:] != "内部工资表.xlsx" {
			return errors.New(msg)
		}
	}
	reg := regexp.MustCompile("[0-9]")
	fmt.Println(filename[:6])
	if !reg.MatchString(filename[:6]) && len(filename[:6]) == 6 {
		return errors.New(msg)
	}
	return nil
}
