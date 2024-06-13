package router

import (
	"server/middleware"
	"server/module/pro_conf/api"

	"github.com/gin-gonic/gin"
)

func InitSysProject(Router *gin.RouterGroup) {
	//var tbOptions limiter.ExpirableOptions
	//tbOptions.DefaultExpirationTTL = time.Second
	//tbOptions.ExpireJobInterval = 0
	//limiter := tollbooth.NewLimiter(2, nil)
	sysProjectRouter := Router.Group("pro-config/").
		//Use(middleware.GinRecovery(false)).
		//Use(middleware.JWTAuth()).
		Use(middleware.SysLogger(1))
	//Use(middleware.Limiter(limiter)).
	//Use(middleware.CasbinHandler())
	//Use(middleware.OperationRecord())
	{

		//项目配置
		sysProjectRouter.GET("sys-project/page", api.GetSysProjectByPage)   // 获取项目配置分页
		sysProjectRouter.GET("sys-project/list", api.GetSysProjectList)     // 获取项目配置
		sysProjectRouter.POST("sys-project/add", api.AddSysProject)         // 新增一个项目配置
		sysProjectRouter.POST("sys-project/edit", api.UpdateSysProjectById) // 根据id更新
		//sysProjectRouter.POST("rmSysProjectByIds", api.RmSysProjectByIds)       // 根据id删除

		//项目模板
		sysProjectRouter.POST("sys-template/update-images", api.UpdateImages)        // 更新模板图片
		sysProjectRouter.POST("sys-template/add", api.AddSysProTemplate)             // 新增项目模板
		sysProjectRouter.GET("sys-template/list", api.GetSysProTempListByProId)      // 根据项目名称获取所有模板
		sysProjectRouter.GET("sys-template/get-info/:id", api.GetSysProTemplateById) // 根据id获取项目模板
		sysProjectRouter.POST("sys-template/copy-temp", api.CopyTemp)                // 复制整个模板

		//项目模板的分块
		sysProjectRouter.GET("sys-block/get-info", api.GetSysProTempBlockByTempId) // 根据模板id获取项目模板分块
		//sysProjectRouter.POST("sys-block/delete-and-add-all", api.RmAllAndAddSysProTempBlockByTempId) // 删除所有然后插入新的
		sysProjectRouter.POST("sys-block/edit", api.EditBlock)                             // 更新
		sysProjectRouter.POST("sys-block/add", api.AddBlock)                               // 插入新的
		sysProjectRouter.DELETE("sys-block/delete", api.DelBlockByIds)                     // 删除项目模板的分块
		sysProjectRouter.POST("sys-block/edit-crop-coordinate", api.UpdateBlockCoordinate) // 分块截图配置
		sysProjectRouter.POST("sys-block/change-order", api.BlockChangeOrder)              // 交换顺序

		//分块关系
		sysProjectRouter.POST("sys-block-relation/delete-and-add-all", api.RmAllAndAddBlockRelations) // 删除所有然后插入新的分块关系
		sysProjectRouter.GET("sys-block-relation/get-info/:id", api.GetBlockRelationsByBId)           // 根据分块id获取分块关系

		//字段配置
		sysProjectRouter.POST("sys-field/add", api.AddFields)            //添加字段
		sysProjectRouter.GET("sys-field/page", api.GetSysFieldsByPage)   //获取项目字段分页
		sysProjectRouter.GET("sys-field/list", api.GetSysFieldsList)     //获取当前项目所有字段
		sysProjectRouter.POST("sys-field/edit", api.UpdateSysFieldsById) //更新字段
		sysProjectRouter.DELETE("sys-field/delete", api.RmSysFieldsById) //删除字段配置
		sysProjectRouter.GET("sys-field/get-const", api.GetConst)        //获取字段配置校验和下拉流程

		//问题件配置
		sysProjectRouter.GET("sys-issue/list", api.GetIssuesByFieldId)   //获取字段问题件配置
		sysProjectRouter.POST("sys-issue/edit", api.EditIssuesByFieldId) //编辑字段问题件配置

		//字段审核配置
		sysProjectRouter.GET("sys-field-check/list", api.GetFieldCheckByFieldId)   //获取字段问题件配置
		sysProjectRouter.POST("sys-field-check/edit", api.EditFieldCheckByFieldId) //编辑字段问题件配置

		//分块字段关系
		sysProjectRouter.POST("sys-block-field-relation/delete-and-add-all", api.RmAllAndAddBFRelations) // 删除所有然后插入新的分块字段关系

		//导出配置
		sysProjectRouter.POST("sys-export/add", api.AddExport)                         // 增加新的导出模板
		sysProjectRouter.POST("sys-export/add-node", api.AddExportNode)                // 增加新的节点
		sysProjectRouter.POST("sys-export/edit", api.UpdateExport)                     // 根据id更新导出模板
		sysProjectRouter.POST("sys-export/edit-node", api.UpdateExportNode)            // 根据id更新导出节点
		sysProjectRouter.GET("sys-export/get-info-page", api.GetExportAndNodesByProId) // 根据项目Id获取节点分页
		sysProjectRouter.DELETE("sys-export/delete", api.DelByIds)                     // 根据id删除导出配置
		sysProjectRouter.POST("sys-export/change-order", api.ChangeOrder)              // 将序号插到序号前
		sysProjectRouter.GET("sys-export/export", api.ExportNodeByExportId)            // 根据exportId导出到Excel
		sysProjectRouter.GET("sys-export/refresh-index", api.RefreshIndex)             // 更新序号
		//sysProjectRouter.POST("updateExports", api.UpdateExports) // 根据id更新节点

		//审核配置
		sysProjectRouter.POST("sys-inspection/add", api.AddInspection)                // 增加审核配置
		sysProjectRouter.POST("sys-inspection/edit", api.UpdateInspection)            // 根据id更新审核配置
		sysProjectRouter.DELETE("sys-inspection/delete", api.RmSysInspectionByIds)    // 批量根据id删除审核配置
		sysProjectRouter.GET("sys-inspection/page", api.GetSysInspectionByPage)       // 分页获取审核配置
		sysProjectRouter.GET("sys-inspection/validation-type", api.GetValidationType) // 分页获取审核配置
		//sysProjectRouter.POST("getSysInspectionByProId", api.GetSysInspectionByProId) // 根据项目id获取审核配置

		//质检配置
		sysProjectRouter.POST("sys-quality/add", api.AddSysQuality)             // 增加质检配置
		sysProjectRouter.POST("sys-quality/edit", api.UpdateSysQuality)         // 根据id更新质检配置
		sysProjectRouter.DELETE("sys-quality/delete", api.RmAddSysQualityByIds) // 批量根据id删除质检配置
		sysProjectRouter.GET("sys-quality/page", api.GetSysQualityByPage)       // 分页获取质检配置
		sysProjectRouter.GET("sys-quality/list", api.GetSysQualityByType)       // 根据类型获取质检配置
		sysProjectRouter.GET("sys-quality/get-const", api.GetQualityConst)
		sysProjectRouter.GET("sys-quality/format", api.GetSysQualityFormat) // 获取格式化后的质检配置
		//sysProjectRouter.POST("getSysQualityByPage",tollbooth_gin.LimitHandler(limiter), api.GetSysQualityByPage)   // 分页获取质检配置

		//ftp监控配置
		sysProjectRouter.GET("sys-ftp-monitor/info", api.FetchFtpMonitorConf) //获取ftp监控配置
		sysProjectRouter.POST("sys-ftp-monitor/edit", api.EditFtpMonitorConf) //编辑ftp监控配置

		//配置复制更新到正式
		sysProjectRouter.POST("sys-project/move-conf", api.MoveProConf)
	}
	//刷新管理配置
	sysProRouter := Router.Group("pro-config/").
		//Use(middleware.GinRecovery(false)).
		Use(middleware.SysLogger(3))
	//Use(middleware.Limiter(limiter))
	{
		sysProRouter.POST("refresh-pro-conf", api.RefreshProConf) // 刷新管理项目配置
	}

}
