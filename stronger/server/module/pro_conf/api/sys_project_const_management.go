package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mozillazg/go-pinyin"
	"os"
	"path"
	"regexp"
	"server/global"
	"server/global/response"
	model2 "server/module/pro_conf/model"
	"server/module/pro_conf/model/request"
	response2 "server/module/pro_conf/model/response"
	"server/module/pro_conf/service"
	responseSysBase "server/module/sys_base/model/response"
	"server/utils"
	u "server/utils"
	"strings"
)

// GetConstTablesByProjectName
// @Tags Const (常量管理)
// @Summary	常量管理--获取常量表的基本信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/project-cost-management/list-const-name [get]
func GetConstTablesByProjectName(c *gin.Context) {
	//根据项目名称获取常量表的基本信息
	var R request.ConstTableWithProject
	R.Conf = "const"
	R.ProCode = c.Query("proCode")
	ConstTableVerify := utils.Rules{
		"ProCode": {utils.NotEmpty()},
		//"Conf":    {utils.NotEmpty()},
	}
	ConstTableVerifyErr := utils.Verify(R, ConstTableVerify)
	if ConstTableVerifyErr != nil {
		response.FailWithMessage(ConstTableVerifyErr.Error(), c)
		return
	}
	GetTablesErr, constTables, count := service.GetConstTablesByProjectName(c, R)
	if GetTablesErr != nil {
		response.FailWithDetailed(response.ERROR, response2.ConstTablesWithName{ConstTables: constTables}, fmt.Sprintf("%v", GetTablesErr), c)
	} else {
		response.OkDetailed(response2.ConstTablesWithName{ConstTables: constTables, Count: count}, "获取成功", c)
	}
}

// PutConstTableByArr
// @Tags Const (常量管理)
// @Summary 常量管理--直接保存传进来的数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.UpdateConstTableStructWithItems true "保存常量"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /pro-config/project-cost-management/putConstTableByArr [post]
func PutConstTableByArr(c *gin.Context) {
	var R request.UpdateConstTableStructWithItems
	R.Conf = "const"
	_ = c.ShouldBindJSON(&R)
	ConstTableVerify := utils.Rules{
		"Id":       {utils.NotEmpty()},
		"Rev":      {utils.NotEmpty()},
		"Project":  {utils.NotEmpty()},
		"FileName": {utils.NotEmpty()},
	}
	ConstTableVerifyErr := utils.Verify(R, ConstTableVerify)
	if ConstTableVerifyErr != nil {
		response.FailWithMessage(ConstTableVerifyErr.Error(), c)
		return
	}
	PutTablesErr, constTable := service.PutConstTableByArr(c, R)
	if PutTablesErr != nil {
		response.FailWithDetailed(response.ERROR, response2.ConstTableResp{ConstTable: constTable}, fmt.Sprintf("%v", PutTablesErr), c)
	} else {
		response.OkDetailed(response2.ConstTableResp{ConstTable: constTable}, "保存成功", c)
	}
}

// DelConstTableLineById
// @Tags Const (常量管理)
// @Summary	常量管理--根据id和rev删除常量表的行
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.DelConstTableLineArr true "删除常量"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /pro-config/project-cost-management/delete-line [post]
func DelConstTableLineById(c *gin.Context) {
	var R request.DelConstTableLineArr
	_ = c.ShouldBindJSON(&R)
	DelConstTableVerify := utils.Rules{
		"DbName": {utils.NotEmpty()},
	}
	DelConstTableVerifyErr := utils.Verify(R, DelConstTableVerify)
	if DelConstTableVerifyErr != nil {
		response.FailWithMessage(DelConstTableVerifyErr.Error(), c)
		return
	}
	DelErr, _ := service.DelConstTableLineById(c, R)
	if DelErr != nil {
		response.FailWithDetailed(response.ERROR, nil, fmt.Sprintf("%v", DelErr), c)
	} else {
		response.OkDetailed(nil, "删除成功", c)
	}
}

// DelConstTableById
// @Tags Const (常量管理)
// @Summary	常量管理--删除数据库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.DelConstTable true "删除常量表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /pro-config/project-cost-management/delete-table [post]
func DelConstTableById(c *gin.Context) {
	var R request.DelConstTable
	_ = c.ShouldBind(&R)

	DelConstTableVerify := utils.Rules{
		"DbName": {utils.NotEmpty()},
	}

	DelConstTableVerifyErr := utils.Verify(R, DelConstTableVerify)
	if DelConstTableVerifyErr != nil {
		response.FailWithMessage(DelConstTableVerifyErr.Error(), c)
		return
	}
	DelErr := service.DeleteDatabase(R.DbName)
	if DelErr != nil {
		response.FailWithDetailed(response.ERROR, "", fmt.Sprintf("%v", DelErr), c)
	} else {
		response.OkDetailed("", "删除成功", c)
	}
}

// PutConstTableWithExcel
// @Tags Const (常量管理)
// @Summary	常量管理--上传常量表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proId formData string true "项目编码"
// @Param file formData file true "文件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /pro-config/project-cost-management/putConstTableByExcel [post]
func PutConstTableWithExcel(c *gin.Context) {
	var C request.ConstRequest
	//读取多表文件，然后其中一个表没保存成功的话，就将该表名称添加到数组中，最后如果该数组不为空，则提示有文件没有保存成功
	projectCode := c.Request.PostFormValue("proId")
	if projectCode == "" {
		response.FailWithMessage("The argument of \"project\" is not allowed to be an empty string", c)
		return
	}
	C.ProCode = projectCode
	form, _ := c.MultipartForm()
	// 获取上传文件
	excelFiles := form.File["file"]
	fmt.Println(len(excelFiles))
	// 线上 uploads/file
	basicPath := global.GConfig.LocalUpload.FilePath + "const/" + projectCode + "/"
	// 本地
	//basicPath := "D:/Projects/tempfiles/"
	fmt.Println("The path of saving excels is " + basicPath)
	relationShips := make(map[string]string, 0)
	failArray := make([]string, 0)
	for _, excelFile := range excelFiles {
		filename := excelFile.Filename
		// 判断文件名字是否符合要求
		name := strings.Replace(filename, ".xlsx", "", -1)
		err, dbname := ChineseNameToPinyinName(projectCode, name)
		if err != nil {
			failArray = append(failArray, err.Error())
			continue
		}
		relationShips[filename] = dbname
		// 设置文件需要保存的指定位置并设置保存的文件名字
		dst := path.Join(basicPath, filename)
		fmt.Println(filename + "has been saved in this path " + basicPath)
		isExist, err := u.Exists(basicPath)
		if err != nil {
			failArray = append(failArray, err.Error())
			continue
		}
		if !isExist {
			err = os.MkdirAll(basicPath, 0777)
			if err != nil {
				failArray = append(failArray, err.Error())
				continue
			}
			err = os.Chmod(basicPath, 0777)
			if err != nil {
				failArray = append(failArray, err.Error())
				continue
			}
		}
		// 上传文件到指定的路径
		saveErr := c.SaveUploadedFile(excelFile, dst)
		if saveErr != nil {
			failArray = append(failArray, filename+", 保存失败!")
			continue
		}
	}
	C.Relationship = relationShips
	//_, _ = service.PutConstTableWithExcel(C.ProCode, C.Relationship)
	if len(failArray) != 0 {
		response.FailWithDetailed(response.ERROR, failArray, "上传部分常量表失败, 请重新上传保存失败的常量表!", c)
	} else {
		//保存数据
		uid := c.Request.Header.Get("x-user-id")
		if uid == "" {
			response.FailWithMessage("x-user-id is empty", c)
			return
		}
		err, a := utils.GetRedisExport("PutConstTableWithExcel", uid)
		if err == nil && a == "true" {
			response.FailWithMessage("正在导入!", c)
			return
		}
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(fmt.Sprintf("常量管理--上传常量表导入失败, err: %s", r))
					global.GLog.Error(fmt.Sprintf("常量管理--上传常量表导入失败, err: %s", r))
					err = utils.DelRedisExport("PutConstTableWithExcel", uid)
					if err != nil {
						fmt.Println(fmt.Sprintf("常量管理--上传常量表删除导出缓存失败, err: %s", err.Error()))
						global.GLog.Error(fmt.Sprintf("常量管理--上传常量表删除导出缓存失败, err: %s", err.Error()))
					}
				}
			}()
			err := utils.SetRedisExport("PutConstTableWithExcel", uid)
			if err != nil {
				//global.GSocketConnSendMsgMap[uid].Emit("constReply", response.ExportResponse{
				//	Code: 400,
				//	Data: "",
				//	Msg:  fmt.Sprintf("PtSalaryFileUpload SetRedisExport err: %s", err.Error()),
				//})
				err = utils.DelRedisExport("PutConstTableWithExcel", uid)
				return
			}
			_, _ = service.PutConstTableWithExcel(C.ProCode, C.Relationship)
			//if msg != "" {
			//	global.GSocketConnSendMsgMap[uid].Emit("constReply", response.ExportResponse{
			//		Code: 400,
			//		Data: "",
			//		Msg:  fmt.Sprintf("PtSalaryFileUpload SetRedisExport err: %s", msg),
			//	})
			//	err = utils.DelRedisExport("PutConstTableWithExcel", uid)
			//	return
			//}
			//global.GSocketConnSendMsgMap[uid].Emit("constReply", response.ExportResponse{
			//	Code: 200,
			//	Data: "",
			//	Msg:  "上传常量表成功!",
			//})
			err = utils.DelRedisExport("PutConstTableWithExcel", uid)
			return
		}()
		response.Ok(c)
	}
}

// GetTableTop
// @Tags Const (常量管理)
// @Summary	常量管理--获取常量表表头
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param dbName query string true "数据库名字"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /pro-config/project-cost-management/list-table-top [get]
func GetTableTop(c *gin.Context) {
	dbname := c.Query("dbName")
	if dbname == "" {
		response.FailWithMessage(errors.New("dbname is empty").Error(), c)
		return
	}
	GetTablesErr, tabletop := service.GetTableTop(c, dbname)
	if GetTablesErr != nil {
		response.FailWithDetailed(response.ERROR, tabletop, fmt.Sprintf("%v", GetTablesErr), c)
	} else {
		response.OkDetailed2(tabletop, "获取成功", c)
	}
}

// ExportConstExcel
// @Tags Const (常量管理)
// @Summary	常量管理--导出常量表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param constName query string true "常量表对应的数据库名称"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /pro-config/project-cost-management/export-const [get]
func ExportConstExcel(c *gin.Context) {
	proCode := c.Query("proCode")
	if proCode == "" {
		response.FailWithMessage("proCode is empty", c)
		return
	}
	dbName := c.Query("constName")
	if dbName == "" {
		response.FailWithMessage("constName is empty", c)
		return
	}

	uid := c.Request.Header.Get("x-user-id")
	if uid == "" {
		response.FailWithMessage("x-user-id is empty", c)
		return
	}
	err, a := utils.GetRedisExport("exportConst", uid)
	if err == nil && a == "true" {
		response.FailWithMessage("正在导出!", c)
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(fmt.Sprintf("常量管理--导出常量表导出失败, err: %s", r))
				global.GLog.Error(fmt.Sprintf("常量管理--导出常量表导出失败, err: %s", r))
				err = utils.DelRedisExport("exportConst", uid)
				if err != nil {
					fmt.Println(fmt.Sprintf("常量管理--导出常量表删除导出缓存失败, err: %s", err.Error()))
					global.GLog.Error(fmt.Sprintf("常量管理--导出常量表删除导出缓存失败, err: %s", err.Error()))
				}
			}
		}()
		//可以广播同一个登录人的客户端的写法
		//global.GSocketIo.BroadcastToRoom("/global-export", uid, "outputStatistics", "所要发送的消息")
		err := utils.SetRedisExport("exportConst", uid)
		if err != nil {
			//global.GSocketConnSendMsgMap[uid].Emit("constReply", response.ExportResponse{
			//	Code: 400,
			//	Data: "",
			//	Msg:  fmt.Sprintf("exportConst SetRedisExport err: %s", err.Error()),
			//})
			err = utils.DelRedisExport("exportConst", uid)
			return
		}
		err, p := service.ExportConstExcel(proCode, dbName)
		fmt.Println("path", p)
		fmt.Println("uid", uid)
		if err != nil {
			//global.GSocketConnSendMsgMap[uid].Emit("constReply", response.ExportResponse{
			//	Code: 400,
			//	Data: "",
			//	Msg:  fmt.Sprintf("exportConst err: %s", err.Error()),
			//})
			err = utils.DelRedisExport("exportConst", uid)
			return
		}
		//fmt.Println("123", global.GSocketConnSendMsgMap)
		//global.GSocketConnSendMsgMap[uid].Emit("constReply", response.ExportResponse{
		//	Code: 200,
		//	Data: p,
		//	Msg:  "导出完成!",
		//})
		err = utils.DelRedisExport("exportConst", uid)
		fmt.Println("导出完成!")
		return
	}()

	//_, _ = service.ExportConstExcel(proCode, dbName)
	//if err != nil {
	//	global.GLog.Error(err.Error())
	//	response.FailWithMessage(fmt.Sprintf("导出失败，%v", err), c)
	//} else {
	response.Ok(c)
	//}
}

// ChineseNameToPinyinName 常量表excel中文转英文
func ChineseNameToPinyinName(proCode string, name string) (error, string) {
	//判断name是否符合项目编码_项目名称_常量表名称
	NameArr := strings.Split(name, "_")
	if NameArr[0] != proCode || len(NameArr) != 3 {
		return errors.New("文件名称请按照一下格式: 项目编码_项目名称_常量表名称 "), ""
	}
	//CouchDB, Creates a new database. The database name {db} must be composed by following next rules:
	//1. Name must begin with a lowercase letter (a-z)
	//2. Lowercase characters (a-z)
	//3. Digits (0-9)
	//4. Any of the characters _, $, (, ), +, -, and /
	//reference: https://docs.couchdb.org/en/stable/api/database/common.html#get--db

	//匹配中文字符： [u4e00-u9fa5]
	reg1 := regexp.MustCompile("^([\u4e00-\u9fa5]+)$")
	//匹配大写英文： [A-Z]
	reg2 := regexp.MustCompile("^[A-Z]$")
	////匹配小写英文： [a-z]
	//reg3 := regexp.MustCompile("^[a-z]$")
	////匹配数字： [0-9]
	//reg4 := regexp.MustCompile("^[0-9]$")
	////匹配其他字符
	//reg5 := regexp.MustCompile("^[_$()+-/]$")
	reg6 := regexp.MustCompile("[a-z0-9_$()+/-]")
	arg := pinyin.NewArgs()
	arg.Style = pinyin.TONE2
	dbName := ""
	isError := false
	for _, v := range name {
		//if unicode.Is(unicode.Han, r) {
		//	fmt.Print(string(r))
		//}
		isWrongName := false
		str := string(v)
		if reg1.MatchString(str) {
			//fmt.Println("ChineseNameToPinyinName: ", string(v), pinyin.Pinyin(str, arg))
			pinYin := pinyin.Pinyin(str, arg)
			dbName += pinYin[0][0]
			isWrongName = true
		}
		if reg2.MatchString(str) {
			dbName += strings.ToLower(str)
			isWrongName = true
		}
		if reg6.MatchString(str) {
			dbName += str
			isWrongName = true
		}
		//if str == "_" || str == "$" || str == "(" || str == ")" || str == "+" || str == "-" || str == "/" {
		//	dbName += str
		//	isWrongName = true
		//}
		if !isWrongName {
			isError = true
			break
		}
	}
	if isError {
		//请看函数头部规则信息
		return errors.New("文件名称含有不允许的符号, 详细信息请咨询软件！"), ""
	}
	return nil, dbName
}

// PageOperationLog
// @Tags Const (常量管理)
// @Summary	常量管理--获取常量操作日志
// @Description
// @Date 2024/3/14
// @Security XToken
// @Security XUserID
// @Security XCode
// @Security XIsIntranet
// @Accept json
// @Produce json
// @Param data body model2.ConstLogReq true "请求实体类"
// @Success 200 {object} response.Response
// @Router /pro-config/const/operation-log/page [post]
func PageOperationLog(c *gin.Context) {
	var constLogReq model2.ConstLogReq
	err := c.ShouldBindJSON(&constLogReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err, total, list := service.FetchConstLogPage(constLogReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(responseSysBase.PageResult{
		List:      list,
		Total:     total,
		PageIndex: constLogReq.PageIndex,
		PageSize:  constLogReq.PageSize,
	}, c)
}
