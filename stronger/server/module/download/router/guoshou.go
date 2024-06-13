/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/4/12 14:49
 */

package router

import (
	"github.com/gin-gonic/gin"
	"server/module/download/api"
)

func InitDownload(Router *gin.RouterGroup) {
	router := Router.Group("gs/")
	{
		router.POST("upload-file", api.UploadFile) //推送案件
		router.POST("upload-hospital", api.UploadHospital) //推送医疗机构
		router.POST("upload-catalogue", api.UploadCatalogue) //推送医疗目录
	}
}
