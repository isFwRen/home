/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/6 15:19
 */

package service

import (
	"server/global"
	"server/module/msg_manager/model"
	model2 "server/module/msg_manager/model/request"
)

//GetDingtalkNoticeSendMsgPage 群通知信息记录
func GetDingtalkNoticeSendMsgPage(dingtalkNoticeMsgReqSearch model2.DingtalkNoticeMsgReq) (err error, total int64, list []model.DingtalkNotice) {
	limit := dingtalkNoticeMsgReqSearch.PageSize
	offset := dingtalkNoticeMsgReqSearch.PageSize * (dingtalkNoticeMsgReqSearch.PageIndex - 1)
	db := global.GDb.Model(&model.DingtalkNotice{})
	err = db.Count(&total).Error
	if err != nil {
		return err, 0, list
	}
	err = db.Order("created_at").Limit(limit).Offset(offset).Find(&list).Error
	return err, total, list
}

//GetDingtalkGroupByIds 根据ids获取钉钉群
func GetDingtalkGroupByIds(ids []string) (err error, list []model.DingtalkGroup) {
	err = global.GDb.Order("created_at").Where("id in (?)", ids).Find(&list).Error
	return err, list
}

//AddDingtalkNotice 新增钉钉群所发的信息记录
func AddDingtalkNotice(notice model.DingtalkNotice) (err error, id string) {
	err = global.GDb.Model(&model.DingtalkNotice{}).Create(&notice).Error
	return err, notice.ID
}

//UpdateDingtalkNoticeStatus 更新发送状态
func UpdateDingtalkNoticeStatus(notice model.DingtalkNotice) error {
	return global.GDb.Where("id = ?", notice.ID).Updates(&notice).Updates(map[string]interface{}{
		"send_status": notice.SendStatus,
		"fail_reason": notice.FailReason,
	}).Error
}

//GetDingtalkGroupById 根据id获取钉钉群
func GetDingtalkGroupById(id string) (err error, list model.DingtalkGroup) {
	err = global.GDb.Where("id = ?", id).First(&list).Error
	return err, list
}
