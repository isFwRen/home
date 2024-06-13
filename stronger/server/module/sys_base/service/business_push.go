package service

import (
	"server/global"
	"server/module/msg_manager/model"
)

// GetUserIdsByMenu 根据权限menu 获取角色 再获取用户id
func GetUserIdsByMenu(proCode string) (err error, ids []string) {
	//50955ef387cf40bbbbfee963c137234e 业务通知  /main/notices/business
	//db, ok := global.ProDbMap[proCode]
	//if !ok {
	//	return global.ProDbErr, ids
	//}
	//err = global.GDb.Raw("SELECT u.id FROM sys_role_menu_relations a,sys_users u WHERE a.menu_id = '50955ef387cf40bbbbfee963c137234e' and a.role_id = u.role_id").Scan(&ids).Error
	err = global.GUserDb.Raw("SELECT b.sys_user_id FROM sys_role_menus a,sys_user_roles b WHERE a.sys_menu_id = '50955ef387cf40bbbbfee963c137234e' and a.sys_role_id = b.sys_role_id").Scan(&ids).Error
	return err, ids
}

// SaveBusinessPush 保存推送的消息
func SaveBusinessPush(msg model.BusinessPush, sends []model.BusinessPushSend) error {
	db, ok := global.ProDbMap[msg.ProCode]
	if !ok {
		return global.ProDbErr
	}
	tx := db.Begin()
	//fmt.Println(msg)
	err := tx.Model(&model.BusinessPush{}).Create(&msg).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	//fmt.Println(sends)
	err = tx.Model(&model.BusinessPushSend{}).Create(&sends).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
