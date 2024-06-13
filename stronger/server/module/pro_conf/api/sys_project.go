package api

import (
	"fmt"
	"go.uber.org/zap"
	"server/global"
	"server/global/response"
	api2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"

	"github.com/flosch/pongo2/v4"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	//"server/module/other/service"
	"server/module/pro_conf/const_data"
	"server/module/pro_conf/model"
	"server/module/pro_conf/model/request"
	responseProConf "server/module/pro_conf/model/response"
	"server/module/pro_conf/service"
	serviceProConf "server/module/pro_conf/service"
	model2 "server/module/sys_base/model"
	requestSysBase "server/module/sys_base/model/request"
	responseSysBase "server/module/sys_base/model/response"
	"server/utils"
	"strconv"
	"strings"
	"time"
)

// AddSysProject
// @Tags ProConfig
// @Summary 配置管理--新增项目配置
// @Security ApiKeyAuth
// @Accept json
// @Produce  json
// @Param data body model.SysProject true "新增SysProject"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /pro-config/sys-project/add [post]
func AddSysProject(c *gin.Context) {
	customClaims, err := api2.GetUserByToken(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var sysProject model.SysProject
	_ = c.ShouldBindJSON(&sysProject)
	ProVerify := utils.Rules{
		"Name": {utils.NotEmpty()},
		"Code": {utils.NotEmpty()},
		"Type": {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(sysProject, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	//新增项目数据库连接
	sysProject.InnerIp = "192.168.202.18"
	sysProject.OutIp = "www.i-confluence.com"
	sysProject.DbTask = "user=postgres password=Change.Postgres host=127.0.0.1 dbname=" + sysProject.Code + "_task port=5433 sslmode=disable TimeZone=Asia/Shanghai"
	sysProject.DbHistory = "user=postgres password=Change.Postgres host=" + sysProject.InnerIp + " dbname=winter_" + sysProject.Code + " port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	err = serviceProConf.AddSysProject(sysProject, customClaims)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithData("创建成功", c)
	}
}

// GetSysProjectByPage
// @Tags ProConfig
// @Summary 配置管理--获取项目配置分页
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /pro-config/sys-project/page [get]
func GetSysProjectByPage(c *gin.Context) {
	var pageInfo requestSysBase.SysProjectRecordSearch
	//err := c.ShouldBind(&pageInfo)
	// 数据模型绑定查询字符串验证
	if err := c.ShouldBindWith(&pageInfo, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}

	err, list, total := serviceProConf.GetSysProjectByPage(pageInfo)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
	} else {
		response.OkWithData(responseSysBase.PageResult{
			List:      list,
			Total:     total,
			PageIndex: pageInfo.PageIndex,
			PageSize:  pageInfo.PageSize,
		}, c)
	}
}

// GetSysProjectList
// @Tags ProConfig
// @Summary 配置管理--获取项目配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /pro-config/sys-project/list [get]
func GetSysProjectList(c *gin.Context) {
	err, list, total := serviceProConf.GetSysProjectList()
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
	} else {
		response.OkWithData(responseSysBase.PageResult{
			List:      list,
			Total:     total,
			PageIndex: 1,
			PageSize:  int(total),
		}, c)
	}
}

// UpdateSysProjectById
// @Tags ProConfig
// @Summary 配置管理--更新项目配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysProjectUpdateRecord true "项目配置实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /pro-config/sys-project/edit [post]
func UpdateSysProjectById(c *gin.Context) {
	var sysProject request.SysProjectUpdateRecord
	_ = c.ShouldBindJSON(&sysProject)
	err := serviceProConf.UpdateSysProjectById(sysProject)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败"), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags ProConfig
// @Summary 配置管理--批量删除SysProject
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body requestSysBase.ReqIds true "批量删除SysProject"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /pro-config/sys-project/rmSysProjectByIds [post]

func RmSysProjectByIds(c *gin.Context) {
	var IDS requestSysBase.ReqIds
	_ = c.ShouldBindJSON(&IDS)
	rows := serviceProConf.RmSysProjectByIds(IDS)
	if rows == 0 {
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkDetailed(responseSysBase.RowResult{
			Row: rows,
		}, "删除成功", c)
	}
}

// AddSysProTemplate
// @Tags ProConfig
// @Summary 配置管理--新增项目模板配置
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body request.SysProTemplate true "模板实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /pro-config/sys-template/add [post]
func AddSysProTemplate(c *gin.Context) {
	var sysProTemplate request.SysProTemplate
	_ = c.ShouldBindJSON(&sysProTemplate)

	ProVerify := utils.Rules{
		"Name":  {utils.NotEmpty()},
		"ProId": {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(sysProTemplate, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	var insertSysProTemplate model.SysProTemplate
	insertSysProTemplate.Name = sysProTemplate.Name
	insertSysProTemplate.ProId = sysProTemplate.ProId
	err := serviceProConf.AddSysProTemplate(insertSysProTemplate)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithData("创建成功", c)
	}
}

// GetSysProTemplateById
// @Tags ProConfig
// @Summary 配置管理--根据id获取项目模板
// @Security ApiKeyAuth
// @Produce  application/json
// @Param id path string true "模版id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /pro-config/sys-template/get-info/{id} [get]
func GetSysProTemplateById(c *gin.Context) {
	sysProTemplateID, err1 := strconv.Atoi(c.Param("id"))

	if sysProTemplateID < 1 || err1 != nil {
		response.FailWithMessage("传参id有误哦", c)
		return
	}

	err, reSysProTemplate := serviceProConf.GetSysProTemplateById(sysProTemplateID)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkDetailed(reSysProTemplate, "查询成功", c)
	}
}

// GetSysProTempListByProId
// @Tags ProConfig
// @Summary 配置管理--根据项目名称获取所有模板
// @Security ApiKeyAuth
// @Produce  application/json
// @Param proId query string false "项目id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /pro-config/sys-template/list [get]
func GetSysProTempListByProId(c *gin.Context) {
	proId := c.Query("proId")
	if proId == "" {
		response.FailWithMessage("传参有误哦", c)
		return
	}
	//intProId, _ := strconv.ParseInt(proId, 10, 64)
	err, reSysProTemplateList := serviceProConf.GetSysProTempListByProId(proId)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkDetailed(reSysProTemplateList, "查询成功", c)
	}
}

// GetSysProTempBlockByTempId
// @Tags ProConfig
// @Summary 配置管理--根据模板id获取项目模板分块
// @Security ApiKeyAuth
// @Produce  application/json
// @Param tempId query string true "模板id"
// @Param blockName query string false "分块名字"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /pro-config/sys-block/get-info [get]
func GetSysProTempBlockByTempId(c *gin.Context) {
	var infoSearch model2.InfoSearch
	if err := c.ShouldBindWith(&infoSearch, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	//tempId, _ := strconv.Atoi(c.Query("tempId"))
	//tempName := c.Query("tempName")
	tempId := infoSearch.TempId
	blockName := infoSearch.BlockName
	err, list, total, max := serviceProConf.GetSysProTempBlockByTempId(tempId, blockName)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkWithData(responseSysBase.HasMaxPageResult{
			List:     list,
			Total:    total,
			MaxCode:  strconv.Itoa(max),
			Page:     0,
			PageSize: 0,
		}, c)
	}
}

// UpdateImages
// @Tags ProConfig
// @Summary 配置管理--更新模板图片
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce  json
// @Param tempImages formData file true "图片"
// @Param sysProTempId formData string true "模板id"
// @Param proCode formData string true "项目编码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /pro-config/sys-template/update-images [post]
func UpdateImages(c *gin.Context) {
	form, err := c.MultipartForm()
	files := form.File["tempImages"]
	sysProTempId := c.PostForm("sysProTempId")
	proCode := c.PostForm("proCode")

	pathArr := make([]string, 0)
	if files == nil {
		response.FailWithMessage("没有图片", c)
		return
	}
	for _, file := range files {
		err, sysProTempBIdPath := utils.SaveImgFile(c, file, file.Filename, global.GConfig.LocalUpload.FilePath+"/"+proCode+"/sys_pro_temp/"+sysProTempId+"/")
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		pathArr = append(pathArr, sysProTempBIdPath)
	}
	//intId, _ := strconv.Atoi(sysProTempId)
	if sysProTempId == "" {
		response.FailWithMessage("sysProTempId为空", c)
		return
	}
	reRow := serviceProConf.UpdateImages(pathArr, sysProTempId)
	if reRow != 1 {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags ProConfig
// @Summary 配置管理--删除所有然后插入新的
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body request.SysProjectBlocks true "删除所有然后插入新的"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /pro-config/sys-block/delete-and-add-all [post]

func RmAllAndAddSysProTempBlockByTempId(c *gin.Context) {
	var sysProjectBlocks request.SysProjectBlocks
	_ = c.ShouldBindJSON(&sysProjectBlocks)
	if sysProjectBlocks.TempId == "" {
		response.FailWithMessage("传参tempId有误哦", c)
		return
	}
	//arr := make([]model.SysProTempB, 0)
	for _, item := range sysProjectBlocks.SysProTempBArr {
		if item.ProTempId == "" || item.ProTempId != sysProjectBlocks.TempId {
			response.FailWithMessage("ProTempId有误哦", c)
			return
		}
	}
	//err := service.RmAllAndAddSysProTempBlockByTempId(arr,sysProjectBlocks.TempId)
	err, reSysProTempBList := serviceProConf.RmAllAndAddSysProTempBlockByTempId(sysProjectBlocks)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkDetailed(reSysProTempBList, "更新成功", c)
	}
}

// AddBlock
// @Tags ProConfig
// @Summary 配置管理-- 新增分块
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body model.SysProTempB true "分块"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /pro-config/sys-block/add [post]
func AddBlock(c *gin.Context) {
	var block model.SysProTempB
	_ = c.ShouldBindJSON(&block)
	ProVerify := utils.Rules{
		"Name": {utils.NotEmpty()},
		//"Code":      {utils.NotEmpty()},
		"ProTempId": {utils.NotEmpty()},
	}

	ProVerifyErr := utils.Verify(block, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	err := serviceProConf.AddBlock(block)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("新增失败，%v", err), c)
	} else {
		response.OkDetailed(err, "新增成功", c)
	}
}

// EditBlock
// @Tags ProConfig
// @Summary 配置管理-- 修改分块
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body model.SysProTempB true "分块"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /pro-config/sys-block/edit [post]
func EditBlock(c *gin.Context) {
	var block model.SysProTempB
	_ = c.ShouldBindJSON(&block)

	ProVerify := utils.Rules{
		"Id":        {utils.NotEmpty()},
		"Name":      {utils.NotEmpty()},
		"Code":      {utils.NotEmpty()},
		"ProTempId": {utils.NotEmpty()},
	}

	ProVerifyErr := utils.Verify(block, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	err := serviceProConf.EditBlock(block)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("编辑失败，%v", err), c)
	} else {
		response.OkDetailed(err, "编辑成功", c)
	}
}

// UpdateBlockCoordinate
// @Tags ProConfig
// @Summary 配置管理--分块截图配置
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body request.SysReqBlockCoordinate true "分块截图配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /pro-config/sys-block/edit-crop-coordinate [post]
func UpdateBlockCoordinate(c *gin.Context) {
	var sysReqBlockCoordinate request.SysReqBlockCoordinate
	_ = c.ShouldBindJSON(&sysReqBlockCoordinate)

	ProVerify := utils.Rules{
		"BlockId": {utils.NotEmpty()},
		//"CoorType":   {utils.NotEmpty()},
		//"PicPage":    {utils.NotEmpty()},
		"Coordinate": {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(sysReqBlockCoordinate, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	reRow := serviceProConf.UpdateBlockCoordinate(sysReqBlockCoordinate)
	if reRow < 1 {
		response.FailWithMessage(fmt.Sprintf("更新失败"), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// RmAllAndAddBlockRelations
// @Tags ProConfig
// @Summary 配置管理--删除所有然后插入新的分块关系
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body request.SysBlocksRelations true "删除所有然后插入新的分块关系"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /pro-config/sys-block-relation/delete-and-add-all [post]
func RmAllAndAddBlockRelations(c *gin.Context) {
	var sysBlocksRelationsObj request.SysBlocksRelations
	_ = c.ShouldBindJSON(&sysBlocksRelationsObj)

	if sysBlocksRelationsObj.BlockId == "" {
		response.FailWithMessage("传参BlockId有误哦", c)
		return
	}
	for _, item := range sysBlocksRelationsObj.TempBlockRelationArr {
		if item.TempBId == "" || item.TempBId != sysBlocksRelationsObj.BlockId {
			response.FailWithMessage("TempBId有误哦", c)
			return
		}
		if item.MyType < 0 || item.MyType != sysBlocksRelationsObj.MyType {
			response.FailWithMessage("MyType有误哦", c)
			return
		}
		ProVerify := utils.Rules{
			//"MyType":   {utils.NotEmpty()},
			"PreBId":   {utils.NotEmpty()},
			"PreBName": {utils.NotEmpty()},
			"PreBCode": {utils.NotEmpty()},
		}
		ProVerifyErr := utils.Verify(item, ProVerify)
		if ProVerifyErr != nil {
			response.FailWithMessage(ProVerifyErr.Error(), c)
			return
		}
	}

	err, reSysProTempBList := serviceProConf.RmAllAndAddBlockRelations(sysBlocksRelationsObj)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkDetailed(reSysProTempBList, "更新成功", c)
	}
}

// GetBlockRelationsByBId
// @Tags ProConfig
// @Summary 配置管理--根据分块id获取分块关系
// @Security ApiKeyAuth
// @Produce  application/json
// @Param id path string true "分块id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /pro-config/sys-block-relation/get-info/{id} [get]
func GetBlockRelationsByBId(c *gin.Context) {
	//bId, err1 := strconv.ParseInt(c.Param("id"), 10, 64)
	bId := c.Param("id")

	if bId == "" {
		response.FailWithMessage("传参id有误哦", c)
		return
	}

	var sysBlockRelationType responseProConf.SysBlockRelationType
	var err interface{}
	err, sysBlockRelationType.OneBlockRelation = serviceProConf.GetBlockRelationsByBIdWithType(bId, "0")
	err, sysBlockRelationType.TwoBlockRelation = serviceProConf.GetBlockRelationsByBIdWithType(bId, "1")
	err, sysBlockRelationType.ThreeBlockRelation = serviceProConf.GetBlockRelationsByBIdWithType(bId, "2")
	err, sysBlockRelationType.TempBFRelation = serviceProConf.GetTempBFRelationByBId(bId)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
	} else {
		response.OkDetailed(sysBlockRelationType, "获取成功", c)
	}
}

// CopyTemp
// @Tags ProConfig
// @Summary 配置管理--复制整个模板
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body model.SysProTemplate true "复制整个模板"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /pro-config/sys-template/copy-temp [post]
func CopyTemp(c *gin.Context) {
	template := model.SysProTemplate{}
	_ = c.ShouldBindJSON(&template)
	ProVerify := utils.Rules{
		"Name":  {utils.NotEmpty()},
		"ProId": {utils.NotEmpty()},
		"ID":    {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(template, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	//sysProTempId := template.ID
	//tempName := c.PostForm("tempName")

	//sysProTempIdInt, err1 := strconv.Atoi(sysProTempId)
	err, tempId := serviceProConf.CopyTemp(template.ID, template.Name)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
	} else {
		response.GetOkWithData(tempId, c)
	}
}

// AddFields
// @Tags ProConfig
// @Summary 配置管理--新增字段
// @Auth xingqiyi
// @Date 2020/11/3 4:11 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysProField true "字段实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-field/add [post]
func AddFields(c *gin.Context) {

	var sysProField model.SysProField
	_ = c.ShouldBindJSON(&sysProField)

	ProVerify := utils.Rules{
		"Name":         {utils.NotEmpty()},
		"Code":         {utils.NotEmpty()},
		"MyOrder":      {utils.NotEmpty()},
		"ProId":        {utils.NotEmpty()},
		"InputProcess": {utils.Gt("0")},
	}
	ProVerifyErr := utils.Verify(sysProField, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	err := serviceProConf.AddFields(sysProField)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("操作失败，%v", err), c)
	} else {
		response.OkWithData("操作成功", c)
	}
}

// GetSysFieldsByPage
// @Tags ProConfig
// @Summary 配置管理--获取项目字段分页
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query int true "页码"
// @Param pageSize query int true "页数量"
// @Param proId query string true "项目id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-field/page [get]
func GetSysFieldsByPage(c *gin.Context) {
	var pageInfo requestSysBase.SysProFieldsSearch
	_ = c.ShouldBindQuery(&pageInfo)
	ProVerify := utils.Rules{
		"ProId":     {utils.NotEmpty()},
		"PageIndex": {utils.Gt("0")},
		"PageSize":  {utils.Gt("0")},
	}
	ProVerifyErr := utils.Verify(pageInfo.BasePageInfo, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	err, list, total := serviceProConf.GetSysFieldsByPage(pageInfo)
	err1, maxCode := serviceProConf.GetMaxCodeSysFieldsByProName(pageInfo.ProId)
	maxCodeInt, _ := strconv.Atoi(strings.Replace(maxCode, "fc", "", -1))
	maxCodeIntAdd := maxCodeInt + 1
	str := "fc"
	if maxCodeIntAdd < 10 {
		str = str + "00"
	} else if 10 <= maxCodeIntAdd && maxCodeIntAdd < 100 {
		str = str + "0"
	}
	maxCodeStr := str + strconv.Itoa(maxCodeIntAdd)
	if err != nil || err1 != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v,%v", err, err1), c)
	} else {
		response.OkWithData(responseSysBase.HasMaxPageResult{
			List:     list,
			Total:    total,
			MaxCode:  maxCodeStr,
			Page:     pageInfo.PageIndex,
			PageSize: pageInfo.PageSize,
		}, c)
	}
}

// UpdateSysFieldsById
// @Tags ProConfig
// @Summary 配置管理--更新字段
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysProField true "字段实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-field/edit [post]
func UpdateSysFieldsById(c *gin.Context) {
	var sysProField model.SysProField
	_ = c.ShouldBindJSON(&sysProField)
	ProVerify := utils.Rules{
		"id": {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(sysProField, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	err := serviceProConf.UpdateSysFieldsById(sysProField)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkWithData("更新成功", c)
	}
}

// RmSysFieldsById
// @Tags ProConfig
// @Summary 配置管理--删除字段配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body requestSysBase.ReqIds true "id数组"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /pro-config/sys-field/delete [delete]
func RmSysFieldsById(c *gin.Context) {
	var idsIntReq requestSysBase.ReqIds
	_ = c.ShouldBindJSON(&idsIntReq)
	rows := serviceProConf.RmSysProjectById(idsIntReq)
	if rows < 1 {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", rows), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DelBlockByIds
// @Tags ProConfig
// @Summary 配置管理--删除项目模板的分块
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body requestSysBase.ReqIds true "id数组"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /pro-config/sys-block/delete [delete]
func DelBlockByIds(c *gin.Context) {
	var idsIntReq requestSysBase.ReqIds
	_ = c.ShouldBindJSON(&idsIntReq)
	rows := serviceProConf.DelBlockByIds(idsIntReq)
	if rows < 1 {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", rows), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// GetSysFieldsList
// @Tags ProConfig
// @Summary 配置管理--获取所有field
// @Auth xingqiyi
// @Date 2020/11/4 3:29 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proId query string true "项目id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-field/list [get]
func GetSysFieldsList(c *gin.Context) {
	//sysProId := c.Query("proId")

	//sysProId, _ := strconv.ParseInt(c.Query("proId"), 10, 64)
	sysProId := c.Query("proId")
	if sysProId == "" {
		response.FailWithMessage("传参sysProId有误哦", c)
		return
	}

	err, reSysProTemplateList := serviceProConf.GetSysFieldsList(sysProId)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkDetailed(reSysProTemplateList, "查询成功", c)
	}
}

// RmAllAndAddBFRelations
// @Tags ProConfig
// @Summary 配置管理--删除所有然后插入新的分块字段关系
// @Auth xingqiyi
// @Date 2020/11/5 3:13 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysBFRelations true "删除所有然后插入新的分块字段关系"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-block-field-relation/delete-and-add-all [post]
func RmAllAndAddBFRelations(c *gin.Context) {
	var sysBFRelations request.SysBFRelations
	_ = c.ShouldBindJSON(&sysBFRelations)

	if sysBFRelations.BlockId == "" {
		response.FailWithMessage("传参BlockId有误哦", c)
		return
	}
	if sysBFRelations.TempBFRelationArr == nil {
		response.FailWithMessage("没有字段哦", c)
		return
	}
	for _, item := range sysBFRelations.TempBFRelationArr {
		if item.TempBId == "" || item.TempBId != sysBFRelations.BlockId {
			response.FailWithMessage("TempBId有误哦", c)
			return
		}
		ProVerify := utils.Rules{
			"FId":   {utils.NotEmpty()},
			"FName": {utils.NotEmpty()},
			"FCode": {utils.NotEmpty()},
		}
		ProVerifyErr := utils.Verify(item, ProVerify)
		if ProVerifyErr != nil {
			response.FailWithMessage(ProVerifyErr.Error(), c)
			return
		}
	}

	err, reSysProTempBList := serviceProConf.RmSysFieldsById(sysBFRelations)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkDetailed(reSysProTempBList, "更新成功", c)
	}
}

// AddExportNode
// @Tags ProConfig
// @Summary 配置管理--增加新的节点
// @Auth xingqiyi
// @Date 2020/11/12 10:08 上午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysExportNode true "增加新的节点"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"新增成功"}"
// @Router /pro-config/sys-export/add-node [post]
func AddExportNode(c *gin.Context) {
	var sysExportNode model.SysExportNode
	_ = c.ShouldBindJSON(&sysExportNode)

	ProVerify := utils.Rules{
		"ExportId": {utils.NotEmpty()},
		"Name":     {utils.NotEmpty()},
		"MyType":   {utils.NotEmpty()},
		"MyOrder":  {utils.Gt("0")},
	}
	ProVerifyErr := utils.Verify(sysExportNode, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	err := serviceProConf.AddExportNode(sysExportNode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//更新模板
	err = UpdateSysExportXml(sysExportNode.ExportId)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithData("创建成功", c)
	}
}

// UpdateExportNode
// @Tags ProConfig
// @Summary 配置管理--根据id更新节点
// @Auth xingqiyi
// @Date 2020/11/12 10:20 上午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysExportNode true "根据id更新节点"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /pro-config/sys-export/edit-node [post]
func UpdateExportNode(c *gin.Context) {
	var sysExportNode model.SysExportNode
	_ = c.ShouldBindJSON(&sysExportNode)

	ProVerify := utils.Rules{
		"ExportId": {utils.NotEmpty()},
		"Name":     {utils.NotEmpty()},
		"MyType":   {utils.NotEmpty()},
		"MyOrder":  {utils.Gt("0")},
		"Id":       {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(sysExportNode, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	sysExportNode.UpdatedAt = time.Now()
	err := serviceProConf.UpdateExportNodes(sysExportNode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	global.GLog.Info(sysExportNode.ExportId)
	//更新模板
	err = UpdateSysExportXml(sysExportNode.ExportId)
	if err != nil {
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithData("更新成功", c)
	}
}

/**
 * @Tags ProConfig
 * @Summary 配置管理--根据项目名字进行重启
 * @Auth xingqiyi
 * @Date 2020/11/19 10:40 上午
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body model.TempConfig true "根据项目名字进行重启"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"重启成功"}"
 * @Router /base/ssh [post]
 */

//func SendSshApp(c *gin.Context) {
//	var R model.TempConfig
//	_ = c.ShouldBindJSON(&R)
//	fmt.Println(R)
//	if R.TempName != "" {
//		ip, pwd, shell, _ := serviceProConf.SearchTempApp(R)
//		cli := service.New(ip, "root", pwd, 22)
//		output, err := cli.Run(shell)
//		fmt.Printf("%v\n%v", output, err)
//		response.OkWithMessage("重启成功"+output, c)
//	} else {
//		response.FailWithMessage("请选择需要重置的项目", c)
//	}
//
//}
//
//func SendSshMain(c *gin.Context) {
//	var R model.TempConfig
//	_ = c.ShouldBindJSON(&R)
//	fmt.Println(R)
//	if R.TempName != "" {
//		ip, pwd, shell, _ := serviceProConf.SearchTempMain(R)
//		cli := service.New(ip, "root", pwd, 22)
//		output, err := cli.Run(shell)
//		fmt.Printf("%v\n%v", output, err)
//		response.OkWithMessage("重启成功"+output, c)
//	} else {
//		response.FailWithMessage("请选择需要重置的项目", c)
//	}
//}

// AddInspection
// @Tags ProConfig
// @Summary 配置管理--增加审核配置
// @Auth xingqiyi
// @Date 2020/11/23 下午5:29
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysInspection true "审核配置实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-inspection/add [post]
func AddInspection(c *gin.Context) {
	var sysInspection model.SysInspection
	_ = c.ShouldBindJSON(&sysInspection)

	ProVerify := utils.Rules{
		"ProId":       {utils.NotEmpty()},
		"ProName":     {utils.NotEmpty()},
		"XmlNodeName": {utils.NotEmpty()},
		"XmlNodeCode": {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(sysInspection, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	err := serviceProConf.AddInspection(sysInspection)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithData("创建成功", c)
	}
}

// UpdateInspection
// @Tags ProConfig
// @Summary 配置管理--根据id更新审核配置
// @Auth xingqiyi
// @Date 2020/11/23 下午5:49
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysInspection true "根据id更新审核配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-inspection/edit [post]
func UpdateInspection(c *gin.Context) {
	var sysInspection model.SysInspection
	_ = c.ShouldBindJSON(&sysInspection)

	ProVerify := utils.Rules{
		"ProId":       {utils.NotEmpty()},
		"ProName":     {utils.NotEmpty()},
		"XmlNodeName": {utils.NotEmpty()},
		"XmlNodeCode": {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(sysInspection, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	sysInspection.UpdatedAt = time.Now()
	reRow := serviceProConf.UpdateInspection(sysInspection)
	if reRow != 1 {
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithData("更新成功", c)
	}
}

// RmSysInspectionByIds
// @Tags ProConfig
// @Summary 配置管理--批量根据id删除审核配置
// @Auth xingqiyi
// @Date 2020/11/23 下午8:23
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body requestSysBase.ReqIds true "批量根据id删除"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-inspection/delete [delete]
func RmSysInspectionByIds(c *gin.Context) {
	var IDS requestSysBase.ReqIds
	_ = c.ShouldBindJSON(&IDS)
	rows := serviceProConf.RmSysInspectionByIds(IDS)
	if rows == 0 {
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkDetailed(responseSysBase.RowResult{
			Row: rows,
		}, "删除成功", c)
	}
}

// AddExport
// @Tags ProConfig
// @Summary 配置管理--新增导出配置
// @Auth xingqiyi
// @Date 2020/12/22 下午3:05
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysExport true "导出实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-export/add [post]
func AddExport(c *gin.Context) {
	var sysExport model.SysExport
	_ = c.ShouldBindJSON(&sysExport)

	ProVerify := utils.Rules{
		"ProId":   {utils.NotEmpty()},
		"ProName": {utils.NotEmpty()},
		"TempVal": {utils.NotEmpty()},
		"XmlType": {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(sysExport, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	err := serviceProConf.AddExport(sysExport)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithData("创建成功", c)
	}
}

// UpdateExport
// @Tags ProConfig
// @Summary 配置管理--根据id更新导出模板
// @Auth xingqiyi
// @Date 2020/12/22 下午3:05
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysExport true "导出配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /pro-config/sys-export/edit [post]
func UpdateExport(c *gin.Context) {
	var sysExport model.SysExport
	_ = c.ShouldBindJSON(&sysExport)

	ProVerify := utils.Rules{
		"ProId":   {utils.NotEmpty()},
		"ProName": {utils.NotEmpty()},
		"TempVal": {utils.NotEmpty()},
		"XmlType": {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(sysExport, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}
	sysExport.UpdatedAt = time.Now()
	reRow := serviceProConf.UpdateExport(sysExport)
	if reRow != 1 {
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithData("更新成功", c)
	}
}

// GetExportAndNodesByProId
// @Tags ProConfig
// @Summary 配置管理--根据项目id获取导出节点
// @Auth xingqiyi
// @Date 2020/12/22 下午3:25
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "页数量"
// @Param proId query string true "项目id"
// @Param name query string false "项目名称"
// @Param fieldLike query string false "其他字段模糊查询"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-export/get-info-page [get]
func GetExportAndNodesByProId(c *gin.Context) {
	var search requestSysBase.SysExportNodesSearch
	//_ = c.ShouldBind(&search)

	if err := c.ShouldBindWith(&search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}

	ProVerify := utils.Rules{
		"PageIndex": {utils.NotEmpty()},
		"PageSize":  {utils.NotEmpty()},
		"ProId":     {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(search, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	err, sysExport, total := serviceProConf.GetExportAndNodesByProId(search)
	if err != nil {
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkDetailed(responseSysBase.PageResult{
			List:      sysExport,
			Total:     total,
			PageIndex: search.PageIndex,
			PageSize:  search.PageSize,
		}, "查询成功", c)
	}
}

// DelByIds
// @Tags ProConfig
// @Summary 配置管理--根据id删除导出配置
// @Auth xingqiyi
// @Date 2020/12/23 上午11:54
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body requestSysBase.ReqIds true "根据id删除"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /pro-config/sys-export/delete [delete]
func DelByIds(c *gin.Context) {
	var idsIntReq requestSysBase.ReqIds
	err := c.ShouldBindJSON(&idsIntReq)
	if err != nil {
		response.FailWithMessage("BindJSON错误", c)
	}

	if len(idsIntReq.Ids) <= 0 {
		response.FailWithMessage("请传入正确的id数组", c)
	}
	//更新模板
	err, sysExportNode := serviceProConf.GetExportNodeById(idsIntReq.Ids[0])
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err, rows := serviceProConf.RmSysExportNodeByIds(idsIntReq)
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	err = UpdateSysExportXml(sysExportNode.ExportId)

	if rows == 0 || err != nil {
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkDetailed(responseSysBase.RowResult{
			Row: rows,
		}, "删除成功", c)
	}
}

// ChangeOrder
// @Tags ProConfig
// @Summary 配置管理--将序号插到序号前
// @Auth xingqiyi
// @Date 2020/12/23 下午2:00
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.ChangeOrder true "将序号插到序号前"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"插入成功"}"
// @Router /pro-config/sys-export/change-order [post]
func ChangeOrder(c *gin.Context) {
	var insertTo request.ChangeOrder
	_ = c.ShouldBindJSON(&insertTo)

	ProVerify := utils.Rules{
		"ExportId":   {utils.NotEmpty()},
		"StartOrder": {utils.NotEmpty()},
		"EndOrder":   {utils.NotEmpty()},
		"StartId":    {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(insertTo, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	if insertTo.StartOrder <= 0 ||
		insertTo.EndOrder <= 0 ||
		insertTo.StartId == "" ||
		insertTo.ExportId == "" {
		response.OkWithMessage("请传正确的数据，小于等于0了", c)
		return
	}

	err := serviceProConf.ChangeOrder(insertTo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//更新模板
	err = UpdateSysExportXml(insertTo.ExportId)
	if err != nil {
		response.FailWithMessage("插入失败", c)
	} else {
		response.OkWithMessage("插入成功", c)
	}
}

// BlockChangeOrder
// @Tags ProConfig
// @Summary 配置管理--将序号插到序号前
// @Auth xingqiyi
// @Date 2020/12/23 下午2:00
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SwitchOrder true "交换顺序"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"插入成功"}"
// @Router /pro-config/sys-block/change-order [post]
func BlockChangeOrder(c *gin.Context) {
	var insertTo request.SwitchOrder
	_ = c.ShouldBindJSON(&insertTo)

	ProVerify := utils.Rules{
		"EndId":      {utils.NotEmpty()},
		"StartOrder": {utils.Gt("-1")},
		"EndOrder":   {utils.Gt("-1")},
		"StartId":    {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(insertTo, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	if insertTo.StartOrder < 0 ||
		insertTo.EndOrder < 0 ||
		insertTo.StartId == "" ||
		insertTo.EndId == "" {
		response.OkWithMessage("请传正确的数据，小于等于0了", c)
		return
	}

	err := serviceProConf.BlockChangeOrder(insertTo)
	if err != nil {
		response.FailWithMessage("交换失败", c)
	} else {
		response.OkWithMessage("交换成功", c)
	}
}

// ExportNodeByExportId
// @Tags ProConfig
// @Summary 配置管理--根据exportId导出到Excel
// @Auth xingqiyi
// @Date 2020/12/28 上午11:20
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param exportId query string true "导出的id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-export/export [get]
func ExportNodeByExportId(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			response.FailWithMessage("导出失败", c)
		}
	}()

	var sysExportNode model.SysExportNode

	//method:post
	//_ = c.ShouldBindJSON(&sysExportNode)

	//method : get
	//exportId, _ := strconv.Atoi(c.Query("exportId"))
	exportId := c.Query("exportId")
	sysExportNode.ExportId = exportId

	ProVerify := utils.Rules{
		"ExportId": {utils.Gt("0")},
	}
	ProVerifyErr := utils.Verify(sysExportNode, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	//获取当前的export的node list
	sysExportNodeList, err1 := serviceProConf.ExportNodeByExportId(sysExportNode.ExportId)
	if err1 != nil {
		panic(err1)
	}
	if len(sysExportNodeList) == 0 {
		response.FailWithMessage("没有数据", c)
		return
	}
	path := global.GConfig.LocalUpload.FilePath + global.PathExportNode + exportId + "/"
	name := "SysExportNode.xlsx"
	err := utils.ExportBigExcel(path, name, "ExportNodeList", sysExportNodeList)

	//err, file := utils.ExportExcel(model.SysExportNodeExport{}, sysExportNodeList)
	//err = file.Save(global.GConfig.LocalUpload.FilePath + "SysExportNode.xlsx")
	if err != nil {
		panic(err)
	}

	//测试模板
	//TestExportXml()
	c.FileAttachment(path+name, name)
}

/////////////////////////////////////////////下面是xml模板操作//////////////////////////////////////////
/////////////////////////////////////////////下面是xml模板操作//////////////////////////////////////////
/////////////////////////////////////////////下面是xml模板操作//////////////////////////////////////////

// UpdateSysExportXml 更新xmlVal
func UpdateSysExportXml(exportId string) error {
	list, err := serviceProConf.ExportNodeByExportId(exportId)
	if err != nil {
		return err
	}
	sysExport, err := serviceProConf.GetExportById(exportId)
	if err != nil {
		return err
	}
	tempVal := SaveExportTemplate(list, sysExport)
	err = serviceProConf.UpdateExportXmlValById(sysExport.ID, tempVal)
	return err
}

// SaveExportTemplate 保存为xml模板
func SaveExportTemplate(list []model.SysExportNode, sysExport model.SysExport) (data string) {
	var encoding = sysExport.XmlType
	global.GLog.Info("encoding：：：" + encoding)
	data = "<?xml version=\"1.0\" encoding=\"" + encoding + "\"?>\n"
	var t, deepFor = 0, 0
	for _, item := range list {
		name := strings.Trim(item.Name, "")
		switch item.MyType {
		case "1":
			for j := 0; j < t; j++ {
				data = data + "	"
			}
			data = data + "<" + name + ">\n"
			t++
		case "2":
			t--
			for j := 0; j < t; j++ {
				data = data + "	"
			}
			data = data + "</" + name + ">\n"
		case "3":
			for j := 0; j < t; j++ {
				data = data + "	"
			}
			data = data + "<" + name + ">" + item.FixedValue + "</" + name + ">\n"
		case "4":
			for j := 0; j < t; j++ {
				data = data + "	"
			}

			data = data + "<" + name + ">{{ '" + item.OneFields + "," + item.TwoFields + "," + item.ThreeFields + "' | " + global.ProCodeId[sysExport.ProId] + "FilterFunc:items" + strconv.Itoa(deepFor) + " }}</" + name + ">\n"
		case "5":
			data = data + "{% " + item.Name + " %}\n"
			if strings.Index(item.Name, "for") != -1 && strings.Index(item.Name, "in") != -1 {
				deepFor++
			}
			if strings.Index(item.Name, "endfor") != -1 {
				deepFor--
			}
		}
	}
	//fmt.Println(data)
	return data
}

//模板的函数

func MyTempFilterFunc(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	//fmt.Println("qw31231231")
	//in.Interface()
	fCode := strings.Split(in.String(), ",")
	for _, item := range fCode {
		fmt.Println(item)
	}
	//global.GLog.Info(in.Interface())
	fmt.Println(param.Interface())
	return pongo2.AsValue(in.Interface()), nil
}

// TestExportXml
// 测试更新和render
func TestExportXml() {
	//UpdateSysExportXml(3)
	//sysExport, _ := serviceProConf.GetExportById("3")
	str := `<?xml version="1.0" encoding="gb2312"?>
			<start>
				<one>
					<two>{{ "in|pp" | MyFunc:param }}</two>
				</one>
				{% for admin in users %}
				<three>{{ admin |  MyFunc:"qweqwe" }}</three>
				{% endfor %}
				{% if is_admin %}
				<p>This user is an admin!</p>
				{% endif %}
			</start>`
	err := pongo2.RegisterFilter("MyFunc", MyTempFilterFunc)
	if err != nil {
		panic(err)
	}
	tpl, err := pongo2.FromString(str)
	if err != nil {
		panic(err)
	}
	//var arr = [5]float32{1000.0, 2.0, 3.4, 7.0, 50.0}
	var arr = []interface{}{map[string]interface{}{"11": "122", "111": "1111"}, map[string]interface{}{"121": "1232", "1141": "11511"}}
	out, err := tpl.Execute(pongo2.Context{"name": "1111", "users": arr})
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}

//模板样式
//<?xml version="1.0" encoding="gb2312"?>
//<start>
//	<one>
//		<two>{{ in | MyFunc:param }}</two>
//	</one>
//	{% for admin in users %}
//	<three>{{ admin |  MyFunc:"qweqwe" }}</three>
//	{% endfor %}
//	{% if is_admin %}
//	<p>This user is an admin!</p>
//	{% endif %}
//</start>

/////////////////////////////////////////////上面是xml模板操作//////////////////////////////////////////
/////////////////////////////////////////////上面是xml模板操作//////////////////////////////////////////
/////////////////////////////////////////////上面是xml模板操作//////////////////////////////////////////

// AddSysQuality
// @Tags ProConfig
// @Summary 配置管理--增加质检配置
// @Auth xingqiyi
// @Date 2021/1/4 下午5:19
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysQuality true "增加质检配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"增加成功"}"
// @Router /pro-config/sys-quality/add [post]
func AddSysQuality(c *gin.Context) {
	var sysQuality model.SysQuality
	e := c.ShouldBindJSON(&sysQuality)
	if e != nil {
		response.FailWithMessage(fmt.Sprintf("%v", e), c)
		return
	}

	ProVerify := utils.Rules{
		//"Name":              {utils.NotEmpty()},
		//"ParentXmlNodeId":   {utils.NotEmpty()},
		"ParentXmlNodeName": {utils.NotEmpty()},
		//"XmlNodeId":         {utils.NotEmpty()},
		"XmlNodeName":  {utils.NotEmpty()},
		"FieldName":    {utils.NotEmpty()},
		"FieldCode":    {utils.NotEmpty()},
		"InputType":    {utils.Gt("0")},
		"BelongType":   {utils.Gt("0")},
		"WidthPrecent": {utils.Gt("0")},
		"MyOrder":      {utils.Gt("0")},
		"ProId":        {utils.NotEmpty()},
		"ProName":      {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(sysQuality, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	customClaims, err1 := api2.GetUserByToken(c)
	if err1 != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err1), c)
	}

	sysQuality.CreatedBy = customClaims.NickName
	sysQuality.UpdatedBy = customClaims.NickName

	err := serviceProConf.AddSysQuality(sysQuality)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithData("创建成功", c)
	}
}

// RmAddSysQualityByIds
// @Tags ProConfig
// @Summary 配置管理--批量根据id删除质检配置
// @Auth xingqiyi
// @Date 2021/1/5 上午9:35
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body requestSysBase.ReqIds true "批量根据id删除"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /pro-config/sys-quality/delete [delete]
func RmAddSysQualityByIds(c *gin.Context) {
	var IDS requestSysBase.ReqIds
	_ = c.ShouldBindJSON(&IDS)
	rows := serviceProConf.RmAddSysQualityByIds(IDS)
	if rows == 0 {
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkDetailed(responseSysBase.RowResult{
			Row: rows,
		}, "删除成功", c)
	}
}

// UpdateSysQuality
// @Tags ProConfig
// @Summary 配置管理--更新质检配置
// @Auth xingqiyi
// @Date 2021/1/5 上午11:00
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysQuality true "更新质检配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /pro-config/sys-quality/edit [post]
func UpdateSysQuality(c *gin.Context) {
	var sysQuality model.SysQuality
	_ = c.ShouldBindJSON(&sysQuality)

	ProVerify := utils.Rules{
		"Id":   {utils.NotEmpty()},
		"Name": {utils.NotEmpty()},
		//"ParentXmlNodeId":   {utils.NotEmpty()},
		"ParentXmlNodeName": {utils.NotEmpty()},
		//"XmlNodeId":         {utils.NotEmpty()},
		"XmlNodeName":  {utils.NotEmpty()},
		"FieldName":    {utils.NotEmpty()},
		"FieldCode":    {utils.NotEmpty()},
		"InputType":    {utils.Gt("0")},
		"BelongType":   {utils.Gt("0")},
		"WidthPrecent": {utils.Gt("0")},
		"MyOrder":      {utils.Gt("0")},
		"ProId":        {utils.NotEmpty()},
		"ProName":      {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(sysQuality, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	customClaims, err1 := api2.GetUserByToken(c)
	if err1 != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err1), c)
	}
	sysQuality.UpdatedBy = customClaims.NickName
	reRow := serviceProConf.UpdateSysQuality(sysQuality)
	if reRow < 1 {
		response.FailWithMessage(fmt.Sprintf("更新失败"), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// GetSysQualityByPage
// @Tags ProConfig
// @Summary 配置管理--分页获取质检配置
// @Auth xingqiyi
// @Date 2021/1/5 上午11:17
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body requestSysBase.SearchSysQuality true "分页获取质检配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-quality/page [get]
func GetSysQualityByPage(c *gin.Context) {
	var search requestSysBase.SearchSysQuality
	_ = c.ShouldBindQuery(&search)

	ProVerify := utils.Rules{
		"Page":     {utils.Gt("0")},
		"PageSize": {utils.Gt("0")},
		"ProId":    {utils.NotEmpty()},
	}
	ProVerifyErr := utils.Verify(search, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	err, sysExport, total, maxOrder := serviceProConf.GetSysQualityByPage(search)
	if err != nil {
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkDetailed(responseSysBase.BasePageResult{
			List:      sysExport,
			Total:     total,
			PageIndex: search.PageIndex,
			PageSize:  search.PageSize,
			MaxOrder:  maxOrder,
		}, "查询成功", c)
	}
}

// GetSysQualityByType
// @Tags ProConfig
// @Summary 配置管理--根据类型获取质检配置
// @Auth xingqiyi
// @Date 2021/1/7 上午10:28
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model2.ProId true "根据类型获取质检配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-quality/list [get]
func GetSysQualityByType(c *gin.Context) {
	var proId model2.ProId
	err := c.ShouldBindJSON(&proId)
	if err != nil {
		response.FailWithMessage("参数有误", c)
		return
	}

	ProVerify := utils.Rules{
		"ProId": {utils.Gt("0")},
	}
	ProVerifyErr := utils.Verify(proId, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	//err, myType := service.GetGroupByProId(proId.ProId)
	//if err != nil {
	//	response.FailWithMessage(fmt.Sprintf("查询分类失败，%v", err), c)
	//	return
	//}
	myType := []int{1, 2, 3, 4, 5, 6, 7, 8}
	type arrSq []model.SysQuality
	var tempMapType map[string]arrSq
	tempMapType = make(map[string]arrSq)
	//var aList = [4][][]model.SysQuality{}
	//arrType := []string{"申请人信息", "被保人信息", "受托人信息", "其他信息", "受益人信息", "领款人信息", "账单信息", "出险信息"}
	//1:申请人信息,2:被保人信息,3:受托人信息,4:其他信息,5:受益人信息,6:领款人信息,7:账单信息,8:出险信息
	for _, item := range myType {
		err1, list := serviceProConf.GetByTypeAndProId(item, proId.ProId)
		if err1 != nil {
			response.FailWithMessage(fmt.Sprintf("根据分类查询失败，%v", err1), c)
			return
		}
		if len(list) > 0 {
			//tempMapType[arrType[item]] = list
			tempMapType[strconv.Itoa(item)] = list
		}
		//switch item {
		//case 1, 2, 3, 4:
		//	if len(list) > 0 {
		//tempMapType[arrType[item]] = list
		//aList[0] = append(aList[0], tempMapType[strconv.Itoa(item)])
		//}
		//case 5, 6:
		//	if len(list) > 0 {
		//tempMapType[strconv.Itoa(item)+"aaaa2"] = list
		//aList[1] = append(aList[1], tempMapType[strconv.Itoa(item)])
		//}
		//case 7:
		//	if len(list) > 0 {
		//tempMapType[strconv.Itoa(item)+"aaaa3"] = list
		//aList[2] = append(aList[2], tempMapType[strconv.Itoa(item)])
		//}
		//case 8:
		//	if len(list) > 0 {
		//tempMapType[strconv.Itoa(item)+"aaaa4"] = list
		//aList[3] = append(aList[3], tempMapType[strconv.Itoa(item)])
		//}

		//}
	}

	//var aList = [4][]model.SysQuality{}
	//var err1 error
	//err1, aList[0] = service.GetByTypeAndProIdAndIds([]int{1, 2, 3, 4}, proId.ProId)
	//if err1 != nil {
	//	response.FailWithMessage(fmt.Sprintf("根据分类查询失败，%v", err1), c)
	//	return
	//}
	//err1, aList[1] = service.GetByTypeAndProIdAndIds([]int{5, 6}, proId.ProId)
	//if err1 != nil {
	//	response.FailWithMessage(fmt.Sprintf("根据分类查询失败，%v", err1), c)
	//	return
	//}
	//err1, aList[2] = service.GetByTypeAndProId(7, proId.ProId)
	//err1, aList[3] = service.GetByTypeAndProId(8, proId.ProId)
	response.OkDetailed(tempMapType, "查询成功", c)

}

// GetSysInspectionByPage
// @Tags ProConfig
// @Summary 配置管理--分页获取审核配置
// @Auth xingqiyi
// @Date 2021/1/5 上午11:17
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "页数量"
// @Param proId    query string true "项目id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-inspection/page [get]
func GetSysInspectionByPage(c *gin.Context) {
	var search requestSysBase.SearchSysInspection
	_ = c.BindQuery(&search)

	ProVerify := utils.Rules{
		"PageIndex": {utils.Gt("0")},
		"PageSize":  {utils.Gt("0")},
		"ProId":     {utils.Gt("0")},
	}
	ProVerifyErr := utils.Verify(search, ProVerify)
	if ProVerifyErr != nil {
		response.FailWithMessage(ProVerifyErr.Error(), c)
		return
	}

	err, sysInspection, total := serviceProConf.GetSysInspectionByPage(search)
	if err != nil {
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkDetailed(responseSysBase.BasePageResult{
			List:      sysInspection,
			Total:     total,
			PageIndex: search.PageIndex,
			PageSize:  search.PageSize,
		}, "查询成功", c)
	}
}

// GetValidationType
// @Tags ProConfig
// @Summary 配置管理--获取审核配置校验类型
// @Auth xingqiyi
// @Date 2021年10月27日13:55:48
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-inspection/validation-type [get]
func GetValidationType(c *gin.Context) {
	response.OkDetailed(const_data.InspectionValidation, "查询成功", c)
}

// GetConst
// @Tags ProConfig
// @Summary 配置管理--获取字段配置校验和下拉流程
// @Auth xingqiyi
// @Date 2021年10月27日13:55:48
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-field/get-const [get]
func GetConst(c *gin.Context) {
	obj := map[string]interface{}{
		"validation":   const_data.FieldValidation,
		"inputProcess": const_data.FieldProcess,
		"dateLimit":    const_data.DateLimit,
	}
	response.OkDetailed(obj, "查询成功", c)
}

// GetQualityConst
// @Tags ProConfig
// @Summary 配置管理--获取质检常量
// @Auth xingqiyi
// @Date 2021年10月27日13:55:48
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-quality/get-const [get]
func GetQualityConst(c *gin.Context) {
	obj := map[string]interface{}{
		"qualityInputType":   const_data.QualityInputType,
		"qualityBelongType":  const_data.QualityBelongType,
		"qualityBeneficiary": const_data.QualityBeneficiary,
		"qualityBillInfo":    const_data.QualityBillInfo,
	}
	response.OkDetailed(obj, "查询成功", c)
}

// GetSysQualityFormat
// @Tags ProConfig
// @Summary 配置管理--获取格式化后的质检配置
// @Auth xingqiyi
// @Date 2021年10月27日13:55:48
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proId query string true "-"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/sys-quality/format [get]
func GetSysQualityFormat(c *gin.Context) {
	proId := c.Query("proId")
	if proId == "" {
		response.FailWithMessage("传参有误哦", c)
		return
	}
	err, qualities := serviceProConf.GetSysQualityFormat(proId)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
		return
	}
	formatQualityList := make(map[int][]model.SysQuality, 0)
	for _, quality := range qualities {
		formatQualityList[quality.BelongType] = append(formatQualityList[quality.BelongType], quality)
	}
	response.OkDetailed(formatQualityList, "查询成功", c)
}

func RefreshIndex(c *gin.Context) {
	serviceProConf.RefreshIndex()
}

// RefreshProConf 刷新管理配置
// @Tags ProConfig
// @Summary 重启管理--刷新管理配置
// @Auth xingqiyi
// @Date 2022/3/17 09:15
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysProCode true "项目编码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/refresh-pro-conf [post]
func RefreshProConf(c *gin.Context) {
	var proCode request.SysProCode
	err := c.ShouldBindJSON(&proCode)
	if err != nil {
		response.FailWithMessage("参数有误", c)
		return
	}
	proVerify := utils.Rules{
		"ProCode": {utils.NotEmpty()},
	}
	proVerifyErr := utils.Verify(proCode, proVerify)
	if proVerifyErr != nil {
		response.FailWithMessage(proVerifyErr.Error(), c)
		return
	}

	err, constMap := utils.LoadConst(proCode.ProCode)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("刷新失败,%v", err.Error()), c)
		return
	}
	conf := global.GProConf[proCode.ProCode]
	conf.ConstTable = constMap
	global.GProConf[proCode.ProCode] = conf

	if err != nil {
		response.FailWithMessage("刷新失败", c)
	} else {
		response.OkWithMessage("刷新成功", c)
	}
}

func MoveProConf(c *gin.Context) {
	var info requestSysBase.ReqMoveConf
	//err := c.ShouldBind(&pageInfo)
	// 数据模型绑定查询字符串验证
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage("参数有误", c)
		return
	}
	proVerify := utils.Rules{
		"ProCode": {utils.NotEmpty()},
		"Mtype":   {utils.NotEmpty()},
	}
	proVerifyErr := utils.Verify(info, proVerify)
	if proVerifyErr != nil {
		response.FailWithMessage(proVerifyErr.Error(), c)
		return
	}

	err = service.MoveData(info.ProCode, info.Mtype, info.TemplateId)
	fmt.Println("-----------info---------------", info, err)
	if err != nil {
		response.FailWithMessage("操作失败", c)
	} else {
		response.OkWithMessage("配置更新成功", c)
	}

}

// FetchFtpMonitorConf
// @Tags ProConfig
// @Summary 配置管理--获取ftp监控配置
// @Description 返回参考实体类
// @Date 2024/4/2
// @Security XToken
// @Security XUserID
// @Security XCode
// @Security XIsIntranet
// @Accept json
// @Produce json
// @Param proCode        query   string   true    "项目代码"
// @Success 200 {object} response.Response
// @Router /pro-config/sys-ftp-monitor/info [get]
func FetchFtpMonitorConf(c *gin.Context) {
	var proCode model.SysFtpMonitorInfoReq
	err := c.ShouldBindQuery(&proCode)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, info := service.FetchFtpMonitorConf(proCode.ProCode)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(info, c)
	}
}

// EditFtpMonitorConf
// @Tags ProConfig
// @Summary 配置管理--编辑ftp监控配置
// @Description
// @Date 2024/4/2
// @Security XToken
// @Security XUserID
// @Security XCode
// @Security XIsIntranet
// @Accept json
// @Produce json
// @Param data body model.SysFtpMonitor true "请求实体类"
// @Success 200 {object} response.Response
// @Router /pro-config/sys-ftp-monitor/edit [post]
func EditFtpMonitorConf(c *gin.Context) {
	var editReq model.SysFtpMonitor
	err := c.ShouldBindJSON(&editReq)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	if editReq.ID == "" {
		u, err := api2.GetUserByToken(c)
		if err != nil {
			response.FailWithParamErr(err, c)
			return
		}
		editReq.CreatedName = u.NickName
		editReq.CreatedCode = u.Code
	}
	row := service.EditFtpMonitorConf(editReq)
	if row == 0 {
		global.GLog.Error("", zap.Any("操作失败 RowsAffected:", row))
		response.FailWithMessage(fmt.Sprintf("失败"), c)
	} else {
		response.OkWithData(row, c)
	}
}
