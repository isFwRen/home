package service

import (
	"server/global"
	"server/module/msg_manager/model"
	model2 "server/module/msg_manager/model/request"
	requestSysBase "server/module/sys_base/model/request"
)

//AddDingtalkGroup 新增钉钉群机器人
func AddDingtalkGroup(dingtalkGroup model.DingtalkGroup) error {
	return global.GDb.Create(&dingtalkGroup).Error
}

//GetDingtalkGroupByPage 分页获取钉钉群机器人
func GetDingtalkGroupByPage(dingtalkGroupSearch model2.DingtalkGroupReq) (err error, total int64, list []model.DingtalkGroup) {
	limit := dingtalkGroupSearch.PageSize
	offset := dingtalkGroupSearch.PageSize * (dingtalkGroupSearch.PageIndex - 1)
	db := global.GDb.Model(&model.DingtalkGroup{}).Where("name like ? ", "%"+dingtalkGroupSearch.Name+"%")
	if dingtalkGroupSearch.ProCode != "" {
		db = db.Where("pro_code = ? ", dingtalkGroupSearch.ProCode)
	}
	if dingtalkGroupSearch.Env > 0 {
		db = db.Where("env = ? ", dingtalkGroupSearch.Env)
	}
	err = db.Count(&total).Error
	if err != nil {
		return err, 0, list
	}
	err = db.Order("created_at").Limit(limit).Offset(offset).Find(&list).Error
	return err, total, list
}

//DelDingtalkGroupByIds 根据id数组删除PT群配置
func DelDingtalkGroupByIds(idsIntReq requestSysBase.ReqIds) (rows int64) {
	return global.GDb.Delete(&[]model.DingtalkGroup{}, idsIntReq.Ids).RowsAffected
}

//UpdateDingtalkGroupById 根据id更新PT群配置
func UpdateDingtalkGroupById(dingtalkGroup model.DingtalkGroup) (err error) {
	return global.GDb.Where("id = ?", dingtalkGroup.ID).Updates(&dingtalkGroup).Updates(map[string]interface{}{
		"name":         dingtalkGroup.Name,
		"env":          dingtalkGroup.Env,
		"pro_code":     dingtalkGroup.ProCode,
		"secret":       dingtalkGroup.Secret,
		"access_token": dingtalkGroup.AccessToken,
	}).Error
}
