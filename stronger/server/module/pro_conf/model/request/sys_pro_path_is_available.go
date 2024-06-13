/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/2 10:15 上午
 */

package request

type SysProPathsList struct {
	SysProPathsList []SysProPaths
}

type SysProPaths struct {
	ID         string `json:"id" form:"id" gorm:"primary_key"`
	IsDownload bool   `json:"isDownload" form:"isDownload" gorm:"comment:'下载启用'"`
	IsUpload   bool   `json:"isUpload" form:"isUpload" gorm:"comment:'回传启用'"`
}
