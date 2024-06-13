package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/global"
	"server/global/response"
	"server/module/pro_manager/model"
	request2 "server/module/pro_manager/model/request"
	"server/module/pro_manager/service"
	responseSysBase "server/module/sys_base/model/response"
	"server/utils"
	api2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"
)

// GetAnnouncementByPage
// @Tags pro-manager(项目管理--公告管理)
// @Summary 项目管理--获取公告管理分页
// @Auth xingqiyi
// @Date 2022/7/4 16:21
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param startTime      		query   string   false   "开始时间格式2022-07-06T16:00:00.000Z"
// @Param endTime        		query   string   false   "结束时间格式2022-07-06T16:00:00.000Z"
// @Param proCode	        	query   string   false    "项目编码"
// @Param releaseType       	query   int      false    "发布类型"
// @Param status        		query   string   false    "发布状态"
// @Param pageIndex      		query   int      true     "页码"
// @Param pageSize       		query   int      true     "数量"
// @Param orderBy        		query   string   false    "排序JSON.stringify([["CreatedAt","desc"]])"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/announcement-manager/page [get]
func GetAnnouncementByPage(c *gin.Context) {
	var announcementPageReq request2.AnnouncementPageReq
	err := c.ShouldBindQuery(&announcementPageReq)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, total, list := service.GetAnnouncementByPage(announcementPageReq)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List:      list,
			Total:     total,
			PageIndex: announcementPageReq.PageIndex,
			PageSize:  announcementPageReq.PageSize,
		}, c)
	}
}

// AddAnnouncement
// @Tags pro-manager(项目管理--公告管理)
// @Summary 项目管理--新增一个公告
// @Auth xingqiyi
// @Date 2022/7/4 16:34
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Announcement true "公告实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/announcement-manager/add [post]
func AddAnnouncement(c *gin.Context) {
	var announcement model.Announcement
	err := c.ShouldBindJSON(&announcement)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	verify := utils.Rules{
		"Title":       {utils.NotEmpty()},
		"ReleaseType": {utils.NotEmpty()},
		"ProCode":     {utils.NotEmpty()},
		"Content":     {utils.NotEmpty()},
	}
	verifyErr := utils.Verify(announcement, verify)
	if verifyErr != nil {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}
	//customClaims, err := api.GetUserByToken(c)
	//if err != nil {
	//	response.FailWithParamErr(err, c)
	//	return
	//}
	//announcement.ReleaseDate = time.Now()
	//announcement.ReleaseUserName = customClaims.NickName
	//announcement.ReleaseUserCode = customClaims.Code
	announcement.Status = 1
	err = service.AddAnnouncement(announcement)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithData("创建成功", c)
	}
}

// ChangeStatusAnnouncementById
// @Tags pro-manager(项目管理--公告管理)
// @Summary 项目管理--删除、恢复(待发布)、发布、取消发布公告
// @Auth xingqiyi
// @Date 2022/7/6 14:14
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.AnnouncementChangeStatusReq true "id数组"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/announcement-manager/change-status [post]
func ChangeStatusAnnouncementById(c *gin.Context) {
	var announcementChangeStatusReq request2.AnnouncementChangeStatusReq
	err := c.ShouldBindJSON(&announcementChangeStatusReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	customClaims, err := api2.GetUserByToken(c)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	rows := service.ChangeStatusAnnouncementById(announcementChangeStatusReq, customClaims)
	if rows < 1 {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", rows), c)
	} else {
		response.OkDetailed(rows, "更新成功", c)
	}
}

// UpdateAnnouncementById
// @Tags pro-manager(项目管理--公告管理)
// @Summary 项目管理--编辑公告
// @Auth xingqiyi
// @Date 2022/7/6 14:21
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Announcement true "公告实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/announcement-manager/edit [post]
func UpdateAnnouncementById(c *gin.Context) {
	var announcement model.Announcement
	err := c.ShouldBindJSON(&announcement)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	verify := utils.Rules{
		"ID":          {utils.NotEmpty()},
		"Title":       {utils.NotEmpty()},
		"ReleaseType": {utils.NotEmpty()},
		"ProCode":     {utils.NotEmpty()},
		"Content":     {utils.NotEmpty()},
	}
	verifyErr := utils.Verify(announcement, verify)
	if verifyErr != nil {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}
	announcement.Status = 1
	err = service.UpdateAnnouncementById(announcement)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkWithData("更新成功", c)
	}
}

// GetHomePageAnnouncement
// @Tags homepage(首页--我的主页)
// @Summary 首页--主页获取公告
// @Auth xingqiyi
// @Date 2022/7/11 10:15
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param startTime      		query   string   false   "开始时间格式2022-07-06T16:00:00.000Z"
// @Param endTime        		query   string   false   "结束时间格式2022-07-06T16:00:00.000Z"
// @Param releaseType       	query   int      true    "发布类型1: 公告通知, 2: 规则动态"
// @Param title        			query   string   false    "标题"
// @Param proCode        		query   string   false    "项目"
// @Param pageIndex      		query   int      true     "页码"
// @Param pageSize       		query   int      true     "数量"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /homepage/home/announcement [get]
func GetHomePageAnnouncement(c *gin.Context) {
	var announcementPageHomeReq request2.AnnouncementPageHomeReq
	err := c.ShouldBindQuery(&announcementPageHomeReq)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	verify := utils.Rules{
		"ReleaseType": {utils.Gt("0")},
	}
	verifyErr := utils.Verify(announcementPageHomeReq, verify)
	if verifyErr != nil {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}

	err, total, list := service.GetHomePageAnnouncement(announcementPageHomeReq)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List:      list,
			Total:     total,
			PageIndex: announcementPageHomeReq.PageIndex,
			PageSize:  announcementPageHomeReq.PageSize,
		}, c)
	}
}

// AnnouncementView
// @Tags homepage(首页--我的主页)
// @Summary 首页--新增访问量
// @Auth xingqiyi
// @Date 2022/8/16 09:54
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id      query   string   true    "公告id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /homepage/home/announcement-view [get]
func AnnouncementView(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.FailWithMessage("传参有误哦", c)
		return
	}
	//global.GRedis.Incr(global.AnnouncementViewPath + id)
	//response.Ok(c)

	err := service.AddView(id)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithoutCode(err.Error(), fmt.Sprintf("新增失败"), c)
	} else {
		response.OkWithMessage("新增成功", c)
	}
}
