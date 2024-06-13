/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/18 15:11
 */

package service

import (
	"server/global"
	"server/module/msg_manager/model"
	model2 "server/module/msg_manager/model/request"
	model3 "server/module/pro_conf/model"
)

//GetTaskNoticeSendMsgPage 分页获取录入通知
func GetTaskNoticeSendMsgPage(dingtalkNoticeMsgReqSearch model2.DingtalkNoticeMsgReq) (err error, total int64, list []model.TaskNotice) {
	limit := dingtalkNoticeMsgReqSearch.PageSize
	offset := dingtalkNoticeMsgReqSearch.PageSize * (dingtalkNoticeMsgReqSearch.PageIndex - 1)
	db := global.GDb.Model(&model.TaskNotice{})
	err = db.Count(&total).Error
	if err != nil {
		return err, 0, list
	}
	err = db.Order("created_at").Limit(limit).Offset(offset).Find(&list).Error
	return err, total, list
}

//GetProPortByCodes 获取项目端口
func GetProPortByCodes(codeArr []string) (err error, list []model3.SysProject) {
	err = global.GDb.Model(&model3.SysProject{}).Select("name,code,in_app_port,out_app_port").Where("id in (?)", codeArr).Find(&list).Error
	return err, list
}

//AddTaskNotice 新增录入发送通知记录
func AddTaskNotice(notice model.TaskNotice) (err error, id string) {
	err = global.GDb.Model(&model.TaskNotice{}).Create(&notice).Error
	return err, notice.ID
}

//AddTaskNotices 新增多条录入发送通知记录
func AddTaskNotices(notices []model.TaskNotice) (err error) {
	return global.GDb.Model(&model.TaskNotice{}).Create(&notices).Error
}

//UpdateTaskNoticeStatus 更新发送状态
func UpdateTaskNoticeStatus(notice model.TaskNotice) error {
	return global.GDb.Where("id = ?", notice.ID).Updates(&notice).Updates(map[string]interface{}{
		"send_status": notice.SendStatus,
		"fail_reason": notice.FailReason,
	}).Error
}
