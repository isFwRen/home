package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global/response"
	sys_base2 "server/module/sys_base/model"
	resp "server/module/sys_base/model/response"
	"server/module/sys_base/service"

	"server/utils"
)

// @Tags authority
// @Summary 创建角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.dao.SysAuthority true "创建角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /authority/createAuthority [post]

func CreateAuthority(c *gin.Context) {
	var auth sys_base2.SysAuthority
	_ = c.ShouldBindJSON(&auth)
	AuthorityVerify := utils.Rules{
		"AuthorityId":   {utils.NotEmpty()},
		"AuthorityName": {utils.NotEmpty()},
		"ParentId":      {utils.NotEmpty()},
	}
	AuthorityVerifyErr := utils.Verify(auth, AuthorityVerify)
	if AuthorityVerifyErr != nil {
		response.FailWithMessage(AuthorityVerifyErr.Error(), c)
		return
	}
	err, authBack := service.CreateAuthority(auth)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		//response.OkWithData(sys_base3.SysAuthorityResponse{Authority: authBack}, c)
		response.OkWithData(authBack, c)
	}
}

// @Tags authority
// @Summary 设置角色资源权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysAuthority true "设置角色资源权限"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /authority/updateAuthority [post]

//func UpdateAuthority(c *gin.Context) {
//	var auth model.SysAuthority
//	_ = c.ShouldBindJSON(&auth)
//	AuthorityVerify := utils.Rules{
//		"AuthorityId":   {utils.NotEmpty()},
//		"AuthorityName": {utils.NotEmpty()},
//		"ParentId":      {utils.NotEmpty()},
//	}
//	AuthorityVerifyErr := utils.Verify(auth, AuthorityVerify)
//	if AuthorityVerifyErr != nil {
//		response.FailWithMessage(AuthorityVerifyErr.Error(), c)
//		return
//	}
//	err, authority := service.UpdateAuthority(auth)
//	if err != nil {
//		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
//	} else {
//		response.OkWithData(resp.SysAuthorityResponse{Authority: authority}, c)
//	}
//}

// @Tags authority
// @Summary 分页获取角色列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "分页获取用户列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /authority/getAuthorityList [post]

//func GetAuthorityList(c *gin.Context) {
//	var pageInfo request.PageInfo
//	_ = c.ShouldBindJSON(&pageInfo)
//	PageVerifyErr := utils.Verify(pageInfo, utils.CustomizeMap["PageVerify"])
//	if PageVerifyErr != nil {
//		response.FailWithMessage(PageVerifyErr.Error(), c)
//		return
//	}
//	err, list, total := service.GetAuthorityInfoList(pageInfo)
//	if err != nil {
//		response.FailWithMessage(fmt.Sprintf("获取数据失败，%v", err), c)
//	} else {
//		response.OkWithData(resp.PageResult{
//			List:     list,
//			Total:    total,
//			Page:     pageInfo.Page,
//			PageSize: pageInfo.PageSize,
//		}, c)
//	}
//}

// @Tags authority
// @Summary 设置角色资源权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysAuthority true "设置角色资源权限"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /authority/setDataAuthority [post]

//func SetDataAuthority(c *gin.Context) {
//	var auth model.SysAuthority
//	_ = c.ShouldBindJSON(&auth)
//	AuthorityIdVerifyErr := utils.Verify(auth, utils.CustomizeMap["AuthorityIdVerify"])
//	if AuthorityIdVerifyErr != nil {
//		response.FailWithMessage(AuthorityIdVerifyErr.Error(), c)
//		return
//	}
//	err := service.SetDataAuthority(auth)
//	if err != nil {
//		response.FailWithMessage(fmt.Sprintf("设置关联失败，%v", err), c)
//	} else {
//		response.Ok(c)
//	}
//}

// @Tags authority
// @Summary 默认获取全部角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /authority/getAuthority [post]

func GetAllAuthority(c *gin.Context) {
	err, u := service.GetAllAuthority()
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取数据失败，%v", err), c)
	} else {
		response.OkWithData(resp.Authority{
			List:  u,
		}, c)
	}
}


// @Tags authority
// @Summary 创建新角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AuthorityManager true "创建新角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /authority/createAuthority [post]

func CreateNewAuthority(c *gin.Context){
	var auth sys_base2.AuthorityManager
	_ = c.ShouldBindJSON(&auth)
	err := service.CreateNewAuthority(auth)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithMessage("创建成功",c)
	}
}


// @Tags authority
// @Summary 修改角色信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AuthorityManager true "修改角色信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /authority/updateAuthority [post]

func UpdateAuthority(c *gin.Context){
	var auth sys_base2.AuthorityManager
	_ = c.ShouldBindJSON(&auth)
	err := service.UpdateAuthority(auth)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkWithMessage("更新成功",c)
	}
}

func UpdateAuthorityPower(c *gin.Context){
	var auth sys_base2.AuthorityManager
	_ = c.ShouldBindJSON(&auth)
	err := service.UpdateAuthorityPower(auth)
	//fmt.Println(auth)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkWithMessage("更新成功",c)
	}
}

func DeleteAuthority(c *gin.Context){
	var auth sys_base2.AuthorityManager
	_ = c.ShouldBindJSON(&auth)
	err := service.DeleteAuthority(auth)
	fmt.Println(err)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", err), c)
	} else {
		response.OkWithMessage("删除成功",c)
	}
}