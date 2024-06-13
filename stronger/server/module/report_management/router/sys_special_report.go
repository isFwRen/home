/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/5/16 14:18
 */

package router

import (
	"github.com/gin-gonic/gin"
	"server/module/report_management/api"
)

func InitSpecialReport(Router *gin.RouterGroup) {
	//limiter := tollbooth.NewLimiter(1, nil)
	specialReportRouter := Router.Group("special-report/")
	//Use(middleware.Limiter(limiter)).
	{
		specialReportRouter.GET("new-hospital-catalogue/page", api.PageNewHospitalAndCatalogue)     //目录外数据列表
		specialReportRouter.GET("new-hospital-catalogue/export", api.ExportNewHospitalAndCatalogue) //导出目录外数据

		specialReportRouter.GET("extract-agency/page", api.PageExtractAgency)                       //机构抽取列表
		specialReportRouter.GET("extract-agency/export", api.ExportExtractAgency)                   //导出机构抽取
		specialReportRouter.GET("extract-agency/export-real-time", api.ExportExtractAgencyRealTime) //实时机构抽取列表导出
	}
}
