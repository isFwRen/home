/**
 * @Author: 阿先
 * @Description:
 * @Date: 2020/11/30 17:20
 */

package service

import (
	"fmt"
	"server/global"
	"server/module/pro_conf/model"
	"server/module/pro_conf/model/request"
)

func GetDownloadPathByProjectId(proId string) (err error, projectConfigDownloadPaths []model.SysProDownloadPaths) {
	err = global.GDb.Where("pro_id = ?", proId).Find(&projectConfigDownloadPaths).Error
	return err, projectConfigDownloadPaths
}

func SetDownloadPathAvailable(sysProPathsList request.SysProPathsList) (err error) {
	tx := global.GDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			global.GLog.Error(fmt.Sprintf("%v", r))
			tx.Rollback()
		}
	}()

	for _, v := range sysProPathsList.SysProPathsList {
		tx.Model(&model.SysProDownloadPaths{}).Where("id = ?", v.ID).
			Updates(map[string]interface{}{"is_download": v.IsDownload, "is_upload": v.IsUpload})
	}

	return tx.Commit().Error
}
