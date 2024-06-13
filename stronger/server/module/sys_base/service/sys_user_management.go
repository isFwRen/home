package service

import (
	"gorm.io/gorm/clause"
	"server/global"
	"server/module/sys_base/model"
	"server/module/sys_base/model/request"
	"time"
	model3 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/model"
)

func GetUserInformation(info request.UserManagement) (err error, list interface{}, total int64) {
	var user []model.SysUser
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	db := global.GDb.Model(&model.SysUser{})
	if info.StartTime != "" && info.EndTime != "" {
		StartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.StartTime, time.Local)
		EndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.EndTime, time.Local)
		db = db.Where("mount_guard_date >= ? AND mount_guard_date <= ?", StartTime, EndTime)
	}
	if info.Code != "" {
		db = db.Where("code like ? ", "%"+info.Code+"%")
	}
	if info.Name != "" {
		db = db.Where("nick_name like ? ", "%"+info.Name+"%")
	}
	if info.Phone != "" {
		db = db.Where("phone like ? ", "%"+info.Phone+"%")
	}
	if info.Status != "" {
		if info.Status != "false" {
			db = db.Where("status = ? ", true)
		} else {
			db = db.Where("status = ? ", false)
		}
	}

	if info.Role != "" {
		var roleA []model.SysRoles
		err = global.GDb.Model(&model.SysRoles{}).Where("name like ? ", "%"+info.Role+"%").Find(&roleA).Error
		if err != nil {
			return err, nil, 0
		}
		role := make([]string, 0)
		for _, v := range roleA {
			role = append(role, v.ID)
		}
		db = db.Where("(role_id) IN ? ", role)
	}
	//var proTotal int64
	//var permissionTotal int64
	var Um []model.UserManagement
	var UmItem model.UserManagement
	err = db.Count(&total).Error
	if total > 0 {
		err = db.Order("id desc").Limit(limit).Offset(offset).Find(&user).Error
		if err != nil {
			return err, nil, 0
		}
		for _, u := range user {
			//var pro []proConf.SysProject
			//err = global.GDb.Model(&proConf.SysProject{}).Find(&pro).Count(&proTotal).Error
			//if err != nil {
			//	return err, nil, 0
			//}
			//var proPermission []model.SysProPermission
			//err = global.GDb.Model(&model.SysProPermission{}).Where("user_code = ? ", u.Code).Find(&proPermission).Count(&permissionTotal).Error
			//if err != nil {
			//	return err, nil, 0
			//}
			//if proTotal != permissionTotal {
			//	for _, i := range pro {
			//		exist := false
			//		for _, j := range proPermission {
			//			if i.Code == j.ProCode {
			//				exist = true
			//				break
			//			}
			//		}
			//		if !exist {
			//			item := model.SysProPermission{
			//				ProCode:   i.Code,
			//				ProName:   i.Name,
			//				HasOp0:    false,
			//				HasOp1:    false,
			//				HasOp2:    false,
			//				HasOpq:    false,
			//				HasInNet:  false,
			//				HasOutNet: false,
			//				ObjectId:  u.ObjectId,
			//				UserId:    u.ID,
			//				ProId:     i.ID,
			//				HasPm:     false,
			//				UserCode:  u.Code,
			//			}
			//			proPermission = append(proPermission, item)
			//		}
			//	}
			//}
			var role model.SysRoles
			err = global.GDb.Model(&model.SysRoles{}).Where("id = ? ", u.RoleId).Find(&role).Error
			if err != nil {
				return err, nil, 0
			}
			UmItem.Role = role.Name
			UmItem.SysUser = u
			//UmItem.SysProPermission = proPermission
			Um = append(Um, UmItem)
		}
		return err, Um, total
	}
	return err, nil, total
}

func AddUser(adduser request.UserAdd) error {
	sysUser := model.SysUser{}
	RelationShip := model.UserRoleRelationShip{}
	Role := model.SysRoles{}
	sysUser.Code = adduser.Code
	sysUser.NickName = adduser.NickName
	sysUser.Staff = adduser.Staff
	sysUser.Phone = adduser.Phone
	sysUser.IsMobile = adduser.IsMobileTerminal
	sysUser.Referees = adduser.Referees
	EntryDate, _ := time.ParseInLocation("2006-01-02 15:04:05", adduser.EntryDate, time.Local)
	MountGuardDate, _ := time.ParseInLocation("2006-01-02 15:04:05", adduser.MountGuardDate, time.Local)
	LeaveDate, _ := time.ParseInLocation("2006-01-02 15:04:05", adduser.LeaveDate, time.Local)
	sysUser.EntryDate = EntryDate
	sysUser.MountGuardDate = MountGuardDate
	sysUser.LeaveDate = LeaveDate
	sysUser.Status = true
	err := global.GDb.Model(&model.SysUser{}).Create(&sysUser).Error
	if err != nil {
		return err
	}
	err = global.GDb.Model(&model.SysUser{}).Where("nick_name = ? ", adduser.NickName).Find(&sysUser).Error
	if err != nil {
		return err
	}
	err = global.GDb.Model(&model.SysRoles{}).Where("name = ? ", adduser.Staff).Find(&Role).Error
	if err != nil {
		return err
	}
	RelationShip.UserId = sysUser.ID
	RelationShip.RoleId = Role.ID
	err = global.GDb.Model(&model.UserRoleRelationShip{}).Create(&RelationShip).Error
	return err
}

func UpdateProPermissionInUserInformation(info request.ProPermissionInUserInformation) (err error) {
	//var total int64
	for _, per := range info.SysProPermission {
		//err = global.GDb.Model(&model.SysProPermission{}).Where("pro_code = ? AND pro_name = ? AND user_code = ? ", per.ProCode, per.ProName, per.UserCode).Count(&total).Error
		//if err != nil {
		//	return err
		//}\
		if per.ID == "" {
			permissions := model.SysProPermission{
				ProCode:   per.ProCode,
				ProName:   per.ProName,
				HasOp0:    per.HasOp0,
				HasOp1:    per.HasOp1,
				HasOp2:    per.HasOp2,
				HasOpq:    per.HasOpq,
				HasInNet:  per.HasInNet,
				HasOutNet: per.HasOutNet,
				UserCode:  per.UserCode,
				UserId:    per.UserId,
				ObjectId:  per.ObjectId,
				HasPm:     per.HasPm,
				ProId:     global.ProCodeId[per.ProCode],
			}
			err = global.GDb.Model(&model.SysProPermission{}).Create(&permissions).Error
			if err != nil {
				return err
			}
		} else {
			err = global.GDb.Debug().Model(&model.SysProPermission{}).Where("id = ? ", per.ID).Updates(map[string]interface{}{
				"has_op0":     per.HasOp0,
				"has_op1":     per.HasOp1,
				"has_op2":     per.HasOp2,
				"has_opq":     per.HasOpq,
				"has_in_net":  per.HasInNet,
				"has_out_net": per.HasOutNet,
				"has_pm":      per.HasPm,
				//"user_id":     per.UserId,
				//"user_code":   per.UserCode,
				//"object_id": per.ObjectId,
				//"pro_id":    per.ProId,
			}).Error
			if err != nil {
				return err
			}
		}
	}
	return err
}

// ChangeRole 修改用户角色
func ChangeRole(changeRoleForm request.ChangeRoleForm) (row int64) {
	return global.GDb.Model(model.SysUser{}).
		Where("id = ? and updated_at = ?", changeRoleForm.Id, changeRoleForm.UpdatedAt).
		Updates(map[string]interface{}{
			"role_id":    changeRoleForm.RoleId,
			"updated_at": time.Now(),
		}).RowsAffected
}

func GetUserPermissionInformation(info request.GetPermissionInformation, uid string) (err error, list interface{}, total int64) {
	var UidProPermission []model.SysProPermission
	err = global.GDb.Model(&model.SysProPermission{}).Where("user_id = ? ", uid).Order("pro_code").Find(&UidProPermission).Error
	if err != nil {
		return nil, nil, 0
	}
	var ProPermission []model.SysProPermission
	err = global.GDb.Model(&model.SysProPermission{}).Where("user_id = ? ", info.UserId).Order("pro_code").Find(&ProPermission).Error
	if err != nil {
		return nil, nil, 0
	}
	var user model3.SysUser
	err = global.GUserDb.Model(&model3.SysUser{}).Where("id = ? ", info.UserId).Find(&user).Error
	if err != nil {
		return nil, nil, 0
	}
	r := make([]model.SysProPermission, 0)
	for _, v1 := range UidProPermission {
		r = append(r, v1)
	}
	for i, v1 := range r {
		exist := false
		for _, v2 := range ProPermission {
			if v1.ProCode == v2.ProCode {
				r[i].HasOp0 = v2.HasOp0
				r[i].HasOp1 = v2.HasOp1
				r[i].HasOp2 = v2.HasOp2
				r[i].HasOpq = v2.HasOpq
				r[i].HasInNet = v2.HasInNet
				r[i].HasOutNet = v2.HasOutNet
				r[i].HasPm = v2.HasPm
				r[i].ID = v2.ID
				exist = true
			}
		}
		if !exist {
			r[i].ID = ""
			r[i].HasOp0 = false
			r[i].HasOp1 = false
			r[i].HasOp2 = false
			r[i].HasOpq = false
			r[i].HasInNet = false
			r[i].HasOutNet = false
			r[i].HasPm = false
		}
		//r[i].ObjectId = user.ObjectId
		r[i].UserId = user.ID
		r[i].UserCode = user.Code
	}
	return err, r, int64(len(r))
}

// GetPermissionByProCode 根据项目编码批量获取权限
func GetPermissionByProCode(proCode string) (sysProPermissions []model.SysProPermission, err error) {
	db := global.GDb.Model(&model.SysProPermission{})
	if proCode != "all" && proCode != "" {
		db = db.Where("pro_code = ?", proCode)
	}
	err = db.Joins("inner join sys_users on sys_pro_permissions.user_id = sys_users.id and sys_users.status = true").Find(&sysProPermissions).Error
	return sysProPermissions, err
}

func FetchUserIDAndObjectIDByCode(code string) (err error, u model.SysUser) {
	err = global.GDb.Model(&model.SysUser{}).Where("code = ?", code).Select("id,code,object_id").Find(&u).Error
	return err, u
}

// UpdateProPermission 根据项目编码个工号批量更新权限
func UpdateProPermission(permissions []model.SysProPermission) (rows int64) {
	//tx := global.GDb.Model(&model.SysProPermission{}).Begin()
	//defer func() {
	//	if r := recover(); r != nil {
	//		tx.Rollback()
	//	}
	//}()
	//for _, permission := range permissions {
	//err = tx.Where("pro_code = ? and user_code = ?", permission.ProCode, permission.UserCode).
	//	Updates(map[string]interface{}{
	//		"has_op0":     permission.HasOp0,
	//		"has_op1":     permission.HasOp1,
	//		"has_op2":     permission.HasOp2,
	//		"has_opq":     permission.HasOpq,
	//		"has_pm":      permission.HasPm,
	//		"has_in_net":  permission.HasInNet,
	//		"has_out_net": permission.HasOutNet,
	//	}).Error
	rows = global.GDb.Model(&model.SysProPermission{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "pro_code"}, {Name: "user_code"}},
		DoUpdates: clause.AssignmentColumns([]string{"has_op0", "has_op1", "has_op2", "has_opq", "has_pm", "has_in_net", "has_out_net"}),
	}).Create(&permissions).RowsAffected
	//if rows == 0 {
	//	global.GLog.Error("没有更新呢")
	//	tx.Rollback()
	//}
	//}
	return rows
}

// UserPage 用户管理用户列表
func UserPage(info request.UserManagement) (err error, list []model3.SysUser, total int64) {
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)

	// db := global.GDb.Model(&model.SysUser{}).Select("id,code,name,phone,email,created_at,updated_at,comments")
	db := global.GUserDb.Model(&model3.SysUser{})
	if info.Code != "" {
		db = db.Where("code like ?", info.Code+"%")
	}
	if info.Name != "" {
		db = db.Where("name like ?", info.Name+"%")
	}
	if info.Phone != "" {
		db = db.Where("phone like ?", info.Phone+"%")
	}
	//if info.Status {
	//	db = db.Where("status = ?", info.Status)
	//}

	//db = db.Preload("SysRoles")

	err = db.Count(&total).Error
	if err != nil {
		return err, list, 0
	}
	err = db.Order("created_at").Limit(limit).Offset(offset).Find(&list).Error
	return err, list, total
}
