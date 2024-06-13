/**
 * @Author: 阿先
 * @Description:
 * @Date: 2020/11/2 11:30
 */

package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global/response"
	dingding2 "server/module/dingding/model"
	dingding3 "server/module/dingding/model/request"
	dingding4 "server/module/dingding/model/response"
	"server/module/dingding/service"
	"server/module/sys_base/model"
	response2 "server/module/sys_base/model/response"
	"server/utils"
)

/**
 * @Author 阿先
 * @Date 11:45 2020/11/2
 * @Tags
 * @Summary	添加钉钉群
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
 * @Router /dingdingGroups/addGroup [post]
 */
func AddGroup(c *gin.Context){
	var R dingding3.AddGroupStruct
	fmt.Println("get in here")
	_ = c.ShouldBindJSON(&R)
	DingdingVerify := utils.Rules{
		"Index":              {utils.NotEmpty()},
		"GroupName":          {utils.NotEmpty()},
		"ProjectCode":        {utils.NotEmpty()},
		"ProjectEnvironment": {utils.NotEmpty()},
	}
	DingdingVerifyErr := utils.Verify(R, DingdingVerify)
	if DingdingVerifyErr != nil {
		response.FailWithMessage(DingdingVerifyErr.Error(), c)
		return
	}
	dingdingGroup := &dingding2.DingdingGroups{Index: R.Index, GroupName: R.GroupName, ProjectCode: R.ProjectCode, ProjectEnvironment: R.ProjectEnvironment, AccessToken:R.AccessToken, Secret:R.Secret}
	err, groupReturn := service.AddGroup(*dingdingGroup)
	if err != nil {
		response.FailWithDetailed(response.ERROR, dingding4.DingdingGroupResponse{DingdingGroup: groupReturn}, fmt.Sprintf("%v", err), c)
	} else {
		response.OkDetailed(dingding4.DingdingGroupResponse{DingdingGroup: groupReturn}, "添加成功", c)
	}
}

/**
 * @Author 阿先
 * @Date 15:46 2020/11/2
 * @Tags
 * @Summary 分页获取钉钉群列表
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
 * @Router /dingdingGroups/GetGroupList [post]
 */
func GetGroupList(c *gin.Context)  {
	var pageInfo dingding3.OddsDingdingGroupStruct
	_ = c.ShouldBindJSON(&pageInfo)
	fmt.Println(pageInfo)
	PageVerifyErr := utils.Verify(pageInfo, utils.CustomizeMap["PageVerify"])
	fmt.Println(PageVerifyErr)
	if PageVerifyErr != nil {
		response.FailWithMessage(PageVerifyErr.Error(), c)
		return
	}
	err, list, total := service.SelectGroupListByPage(pageInfo)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取数据失败，%v", err), c)
	} else {
		response.OkWithData(response2.PageResult{
			List:     list,
			Total:    total,
			PageIndex:     pageInfo.PageIndex,
			PageSize: pageInfo.PageSize,
		}, c)
	}
}

/**
 * @Author 阿先
 * @Date 11:54 2020/11/3
 * @Tags
 * @Summary 根据id删除钉钉群
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
 * @Router /dingdingGroups/DelGroup [post]
 */
func DelGroup(c *gin.Context)  {
	var reqId model.GetById
	_ = c.ShouldBindJSON(&reqId)
	IdVerifyErr := utils.Verify(reqId, utils.CustomizeMap["IdVerify"])
	if IdVerifyErr != nil {
		response.FailWithMessage(IdVerifyErr.Error(), c)
		return
	}
	err := service.DelGroup(reqId.Id)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", err), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}

}

/**
 * @Author 阿先
 * @Date 13:50 2020/11/3
 * @Tags
 * @Summary	 根据id更新钉钉群
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
 * @Router [post]
 */
func UpdateGroup(c *gin.Context)  {
	var params dingding2.DingdingGroups
	_ = c.ShouldBindJSON(&params)
	if err := service.UpdateGroup(params, params.ID); err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkWithMessage("修改成功", c)
	}
}
