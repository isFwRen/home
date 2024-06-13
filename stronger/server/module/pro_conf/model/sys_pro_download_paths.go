/**
 * @Author: 阿先
 * @Description:
 * @Date: 2020/11/30 16:20
 */

package model

import (
	"server/module/sys_base/model"
)

type SysProDownloadPaths struct {
	model.Model
	Scan           string `json:"scan" form:"scan" gorm:"comment:'扫描cmd'"`                                    //扫描cmd
	TimeZone       string `json:"timeZone" form:"timeZone" gorm:"comment:'远程主机时区'"`                           //远程主机时区
	FetchBill      string `json:"fetchBill" form:"fetchBill" gorm:"comment:'下载单据cmd'"`                        //下载单据cmd
	MaxDownload    int    `json:"maxDownload" form:"maxDownload" gorm:"comment:'一次最多下载数'"`                    //一次最多下载数
	Backup         string `json:"backup" form:"backup" gorm:"comment:'备份cmd'"`                                //备份cmd
	DownloadClean  string `json:"downloadClean" form:"downloadClean" gorm:"comment:'清理下载文件cmd'"`              //清理下载文件cmd
	Upload         string `json:"upload" form:"upload" gorm:"comment:'回传cmd'"`                                //回传cmd
	UploadRename   string `json:"uploadRename" form:"uploadRename" gorm:"comment:'回传后的改名操作cmd'"`              //回传后的改名操作cmd
	MaxConnections int    `json:"maxConnections" form:"maxConnections" gorm:"comment:'最大连接数'"`                //最大连接数
	FetchFile      string `json:"fetchFile" form:"fetchFile" gorm:"comment:'客户下载文件cmd'"`                      //客户下载文件cmd
	DownloadType   int    `json:"downloadType" form:"downloadType" gorm:"comment:'下载标记类型0：测试内部，1：测试客户，2：正式'"` //下载标记类型0：测试内部，1：测试客户，2：正式
	ProId          string `json:"proId" form:"proId" gorm:"comment:'项目所属id'"`                                 //项目所属id
	ProCode        string `json:"proCode" form:"proCode" gorm:"comment:'项目所属编码'"`                             //项目所属编码
	ProName        string `json:"proName" form:"proName" gorm:"comment:'项目名称'"`                               //项目名称
	IsDownload     bool   `json:"isDownload" form:"isDownload" gorm:"comment:'下载启用'"`                         //下载启用
	IsUpload       bool   `json:"isUpload" form:"isUpload" gorm:"comment:'回传启用'"`                             //回传启用
}
