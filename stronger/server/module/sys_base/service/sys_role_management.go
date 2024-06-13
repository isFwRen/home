package service

import (
	"errors"
	"server/global"
	"server/module/sys_base/model"
	"server/module/sys_base/model/request"
	"strconv"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

//var RoleStatus = map[int]string{
//	1: "正常",
//	2: "停用",
//}

func GetRoleInformation(status string) (err error, roles []model.SysRoles, total int64) {
	var role []model.SysRoles
	Status, _ := strconv.Atoi(status)
	db := global.GDb.Model(&model.SysRoles{})
	err = db.Count(&total).Error
	if status != "" {
		db = db.Where("status = ? ", Status)
	}
	err = db.Find(&role).Error
	if err != nil {
		return err, roles, total
	}
	return err, role, total
}

func AddRoleInformation(role request.AddRoleInformation) error {
	var total int64
	err := global.GDb.Model(&model.SysRoles{}).Where("name = ? ", role.RoleName).Count(&total).Error
	if err != nil {
		return err
	}
	if total > 0 {
		return errors.New("该角色已经存在")
	}
	var newRole model.SysRoles
	newRole.Name = role.RoleName
	newRole.Status = role.RoleStatus
	newRole.Remark = role.RoleDescribe
	newRole.CreatedBy = role.RoleCreatedBy
	err = global.GDb.Model(&model.SysRoles{}).Create(&newRole).Error
	return err
}

func EditRoleInformation(role request.EditRoleInformation) error {
	var R model.SysRoles
	var total int64
	err := global.GDb.Model(&model.SysRoles{}).Where("name = ? ", role.RoleBeforeName).Find(&R).Count(&total).Error
	if err != nil {
		return err
	}
	if total == 1 {
		if role.RoleNewName != "" {
			err = global.GDb.Model(&model.SysRoles{}).Where("name = ? ", role.RoleBeforeName).Updates(map[string]interface{}{
				"name":       role.RoleNewName,
				"status":     role.RoleStatus,
				"remark":     role.RoleDescribe,
				"updated_by": role.RoleUpdatedBy,
			}).Error
		} else {
			err = global.GDb.Model(&model.SysRoles{}).Where("name = ? ", role.RoleBeforeName).Updates(map[string]interface{}{
				"status":     role.RoleStatus,
				"remark":     role.RoleDescribe,
				"updated_by": role.RoleUpdatedBy,
			}).Error
		}
		return err
	}
	return err
}

func DeleteRoleInformation(R request.DeleteRoleInformation) error {
	var role model.SysRoles
	for _, id := range R.Ids {
		err := global.GDb.Model(&model.SysRoles{}).Where("id = ? ", id).Delete(&role).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// GetRoleMenuRelationsByRoleId 根据roleId获取所有menuId
func GetRoleMenuRelationsByRoleId(roleId string) (err error, sysRoleMenuRelations []model.SysRoleMenuRelations) {
	err = global.GDb.Model(&model.SysRoleMenuRelations{}).
		Where("role_id = ? ", roleId).Find(&sysRoleMenuRelations).Error
	return err, sysRoleMenuRelations
}

// GetPmPermissionByUId 根据UId获取PM所有perm
func GetPmPermissionByUId(uid string) (err error, s []model.SysProPermission) {
	err = global.GDb.Model(&model.SysProPermission{}).
		Where("user_id = ? and has_pm = true", uid).Find(&s).Error
	return err, s
}

// GetAllPermissionByUId 根据UId获取所有perm
func GetAllPermissionByUId(uid string) (err error, s []model.SysProPermission) {
	err = global.GDb.Model(&model.SysProPermission{}).
		Where("user_id = ? and pro_code in (?)", uid, global.GConfig.System.ProArr).Order("pro_code").Find(&s).Error
	return err, s
}

// FetchPermissionByUId 根据UId获取所有perm
func FetchPermissionByUId(uid string) (err error, s []model.SysProPermission) {
	err = global.GDb.Model(&model.SysProPermission{}).
		Where("user_id = ?", uid).Order("pro_code").Find(&s).Error
	return err, s
}

// GetRoleByRoleId 根据roleId获取role
func GetRoleByRoleId(roleId string) (err error, sysRole model.SysRoles) {
	err = global.GDb.Model(&model.SysRoles{}).
		Where("id = ? ", roleId).First(&sysRole).Error
	return err, sysRole
}

// SetRoleMenuRelation 设置角色菜单关系
func SetRoleMenuRelation(roleAndMenuRelationForms []request.RoleAndMenuRelationForm) (row int64, err error) {
	for _, roleAndMenuRelationForm := range roleAndMenuRelationForms {
		if roleAndMenuRelationForm.IsSelect {
			tx := global.GDb.Begin()
			row = tx.Model(&model.SysRoleMenuRelations{}).Where("id = ?", roleAndMenuRelationForm.ID).Debug().
				FirstOrCreate(&roleAndMenuRelationForm).RowsAffected
			err, menu := GetMenuById(roleAndMenuRelationForm.MenuId)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
			if menu.MenuType == 1 {
				e, err := GoCasbin()
				if err != nil {
					tx.Rollback()
					return 0, err
				}
				if menu.ApiId == "" || menu.Api == "" || menu.Action == "" {
					tx.Rollback()
					return 0, errors.New("菜单有误")
				}
				_, err = e.AddPolicy(roleAndMenuRelationForm.RoleId, menu.Api, menu.Action)
				if err != nil {
					tx.Rollback()
					return 0, err
				}
			}
			if tx.Commit().Error != nil {
				tx.Rollback()
			}
		} else {
			tx := global.GDb.Begin()
			row = tx.Model(&model.SysRoleMenuRelations{}).
				Where("id = ?", roleAndMenuRelationForm.ID).
				Delete(&roleAndMenuRelationForm).RowsAffected
			err, menu := GetMenuById(roleAndMenuRelationForm.MenuId)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
			if menu.MenuType == 1 {
				e, err := GoCasbin()
				if err != nil {
					tx.Rollback()
					return 0, err
				}
				_, err = e.RemoveFilteredPolicy(0, roleAndMenuRelationForm.RoleId, menu.Api, menu.Action)
				if err != nil {
					tx.Rollback()
					return 0, err
				}
			}
			if tx.Commit().Error != nil {
				tx.Rollback()
			}
		}
	}
	return row, nil
}

func GoCasbin() (*casbin.Enforcer, error) {
	p := global.GConfig.Postgresql
	dsn := "user=" + p.Username + " password=" + p.Password + " host=" + p.Host + " dbname=" + p.Dbname + " port=" + p.Port + " " + p.Config
	global.GLog.Info(dsn)
	a, err := gormadapter.NewAdapter(global.GConfig.System.DbType, dsn, true)
	if err != nil {
		global.GLog.Error("2222:::" + err.Error())
		return nil, err
	}
	e, err := casbin.NewEnforcer(global.GConfig.Casbin.ModelPath, a)
	if err != nil {
		global.GLog.Error("1111:::" + err.Error())
		return nil, err
	}
	e.AddFunction("ParamsMatch", ParamsMatchFunc)
	err = e.LoadPolicy()
	return e, err
}
