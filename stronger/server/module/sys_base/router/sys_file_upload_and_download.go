package router

import (
	"github.com/gin-gonic/gin"
	"server/module/sys_base/api"
)

func InitFileUploadAndDownloadRouter(Router *gin.RouterGroup) {
	FileUploadAndDownloadGroup := Router.Group("file-upload-and-download/")
	// .Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		FileUploadAndDownloadGroup.POST("/upload", api.UploadFile) // 上传文件
		//FileUploadAndDownloadGroup.POST("/list", api.GetFileList)  // 获取上传文件列表
		//FileUploadAndDownloadGroup.POST("/delete", api.DeleteFile)                                   // 删除指定文件
		FileUploadAndDownloadGroup.POST("/breakpoint", api.BreakpointContinue)                       // 断点续传
		//FileUploadAndDownloadGroup.GET("/find-file", api.FindFile)                                   // 查询当前文件成功的切片
		//FileUploadAndDownloadGroup.POST("/breakpoint-continue-finish", api.BreakpointContinueFinish) // 查询当前文件成功的切片
		//FileUploadAndDownloadGroup.POST("/remove-chunk", api.RemoveChunk)                            // 查询当前文件成功的切片
	}
}
