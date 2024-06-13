/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/4/12 15:24
 */

package model

import "server/module/sys_base/model"

var IsDeletedMap = map[string]bool{
	"N": false,
	"Y": true,
}

//ImageFile 图片文件
type ImageFile struct {
	BpoSendRemark    string `json:"bpoSendRemark" form:"bpoSendRemark"`
	BpoSpeDesc       string `json:"bpoSpeDesc" form:"bpoSpeDesc"`
	BranchNo         string `json:"branchNo" form:"branchNo"`
	ClaimNo          string `json:"claimNo" form:"claimNo"`
	ClaimTpaId       string `json:"claimTpaId" form:"claimTpaId"`
	CmsImageInfoList []struct {
		FileName  string `json:"fileName" form:"fileName"`
		FileType  string `json:"fileType" form:"fileType"`
		Id        string `json:"id" form:"id"`
		ImageType string `json:"imageType" form:"imageType"`
		ImageUrl  string `json:"imageUrl" form:"imageUrl"`
	} `json:"cmsImageInfoList" form:"cmsImageInfoList"`
	Url string `json:"url" form:"url"`
}

//Hospital 医疗机构
type Hospital struct {
	BeginNo          string         `json:"beginNo" form:"beginNo"`
	HospitalInfoList []HospitalInfo `json:"hospitalInfoList" form:"hospitalInfoList"`
	TotalNum         string         `json:"totalNum" form:"totalNum"`
}

type HospitalInfo struct {
	BranchNo      string `json:"branchNo" form:"branchNo"`
	DefTime       string `json:"defTime" form:"defTime"`
	HospitalCode  string `json:"hospitalCode" form:"hospitalCode"`
	HospitalGrade string `json:"hospitalGrade" form:"hospitalGrade"`
	HospitalName  string `json:"hospitalName" form:"hospitalName"`
	IsDeleted     string `json:"isDeleted" form:"isDeleted"`
}

//Catalog 医疗目录
type Catalog struct {
	AreaCode    string `json:"areaCode" form:"areaCode"`
	BranchNo    string `json:"branchNo" form:"branchNo"`
	CatalogCode string `json:"catalogCode" form:"catalogCode"`
	CatalogName string `json:"catalogName" form:"catalogName"`
	DefTime     string `json:"defTime" form:"defTime"`
	IsDeleted   string `json:"isDeleted" form:"isDeleted"`
	ResourceUrl string `json:"resourceUrl" form:"resourceUrl"`
	Url         string `json:"url" form:"url"`
}

//UpdateConstLog 更新常量日志
type UpdateConstLog struct {
	model.Model
	Type      int    `json:"type" form:"type"`
	Name      string `json:"name" form:"name"`
	ProCode   string `json:"proCode" form:"proCode"`
	IsDeleted bool   `json:"isDeleted" form:"isDeleted"`
	RemoteUrl string `json:"remoteUrl" form:"remoteUrl"`
	LocalUrl  string `json:"localUrl" form:"localUrl"`
	OtherInfo string `json:"otherInfo" form:"otherInfo"`
	IsUpdated bool   `json:"isUpdated" form:"isUpdated"`
}
