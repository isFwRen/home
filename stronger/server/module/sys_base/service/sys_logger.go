/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/2/9 9:33 上午
 */

package service

import (
	"server/global"
	"server/module/sys_base/model"
	"server/module/sys_base/model/request"
)

//AddSysLogger 新增用户查看的日志
func AddSysLogger(logger model.SysLogger) (err error) {
	err = global.GDb.Model(&model.SysLogger{}).Create(&logger).Error
	return err
}

//GetPageByType 分页获取系统日志
func GetPageByType(param request.SysLogger) (err error, total int64, loggers []model.SysLogger) {

	limit := param.PageSize
	offset := param.PageSize * (param.PageIndex - 1)

	db := global.GDb.Where("log_type = ?", param.LogType)
	if param.ProCode != "" {
		db = db.Where("pro_code = ? ", param.ProCode)
	}
	if param.FunctionModule != "" {
		db = db.Where("function_module = ? ", param.FunctionModule)
	}
	if param.ModuleOperation != "" {
		db = db.Where("module_operation LIKE ?", "%"+param.ModuleOperation+"%")
	}
	if param.OperationPeople != "" {
		db = db.Where("operation_name LIKE ? or operation_code LIKE ? ", "%"+param.OperationPeople+"%", "%"+param.OperationPeople+"%")
	}

	db = db.Model(&model.SysLogger{}).Where("created_at BETWEEN ? AND ? ", param.StartTime, param.EndTime)

	err = db.Count(&total).Error
	err = db.Order("created_at desc").Limit(limit).Offset(offset).Find(&loggers).Error
	return err, total, loggers
}
