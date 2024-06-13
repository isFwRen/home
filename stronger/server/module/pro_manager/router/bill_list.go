/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/3 10:31 上午
 */

package router

import (
	"server/middleware"
	"server/module/pro_manager/api"

	"github.com/gin-gonic/gin"
)

func InitBillList(Router *gin.RouterGroup) {
	//limiter := tollbooth.NewLimiter(10, nil)
	billListRouter := Router.Group("pro-manager/").
		//Use(middleware.GinRecovery(false)).
		//Use(middleware.Limiter(limiter)).
		//Use(middleware.CasbinHandler()).
		Use(middleware.SysLogger(2))
	{
		billListRouter.GET("bill-list/page", api.GetBillByPage)             //案件列表
		billListRouter.DELETE("bill-list/delete", api.DelByIdsAndProCode)   //批量删除软删除
		billListRouter.GET("bill-list/dict-const", api.GetDictConst)        //获取列表的常量
		billListRouter.POST("bill-list/recover", api.RecoverBill)           //恢复
		billListRouter.POST("bill-list/export-err-bill", api.ExportErrBill) //导出异常
		billListRouter.POST("bill-list/force-export", api.ForceExport)      //强制导出
		billListRouter.POST("bill-list/set-upload-type", api.SetUploadType) //手动自动回传
		billListRouter.POST("bill-list/reload", api.Reload)                 //重加载
		billListRouter.POST("bill-list/remark", api.Remark)                 //备注

		billListRouter.POST("bill-list/save-xml", api.SaveXml) //保存xml
		billListRouter.POST("bill-list/upload", api.Upload)    //回传xml

		billListRouter.POST("bill-list/edit-bill-result-data", api.EditBillResultData) //编辑结果数据
		billListRouter.POST("bill-list/zbj-bill-result-data", api.ZbjBillResultData)   //编辑结果数据
		billListRouter.POST("bill-list/is-quality", api.IsQuality)                     //是否有人在质检
		billListRouter.POST("bill-list/save-bill-result-data", api.SaveBillResultData) //保存单据
		billListRouter.GET("bill-list/see-bill-result-data", api.SeeBillResultData)    //查看单据

		billListRouter.GET("bill-list/get-field-info", api.GetFieldInfo)        //获取字段、分块和字段配置信息
		billListRouter.GET("bill-list/get-block-img", api.GetBlockImg)          //获取分块相关图片
		billListRouter.POST("bill-list/edit-feedback-val", api.EditFeedbackVal) //修改字段客户反馈值

		billListRouter.POST("bill-list/set-practice", api.SetPractice)                      //设置练习
		billListRouter.GET("bill-list/get-log", api.GetLog)                                 //获取单据的修改日志
		billListRouter.GET("bill-list/qing-dan/get", api.GetQingDan)                        //获取清单
		billListRouter.GET("bill-list/qing-dan/export", api.ExportQingDan)                  //导出清单
		billListRouter.GET("bill-list/qqq", api.QQQ)                                        //test id int64
		billListRouter.GET("bill-list/get-time-liness-briefing", api.GetTimeLinessBriefing) //test id int64

		billListRouter.GET("bill-del-list/page", api.GetBillDelByPage)
		billListRouter.GET("bill-del-list/sum", api.GetBillDelSum)
		billListRouter.GET("bill-del-list/export", api.ExportBillDel)
		billListRouter.GET("bill-del-list/history", api.HistoryDelList)

		billListRouter.GET("bill-deduction_details/page", api.GetDeductionDetailsByPage)
		billListRouter.GET("bill-deduction_details/export", api.ExportDeductionDetails)

		billListRouter.POST("bill-list/save-qualities-demo", api.SetDataInRedis)
		billListRouter.POST("bill-list/get-qualities-demo", api.GetDataInRedis)
	}
}
