/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/12/6 16:54
 */

package upload

import (
	"errors"
	"go.uber.org/zap"
	"server/global"
	model2 "server/module/pro_conf/model"
	"server/module/pro_manager/model"
	"server/module/upload/project"
	"server/module/upload/project/B0108"
	"server/module/upload/project/B0113"
	"server/module/upload/project/B0114"
	"server/module/upload/project/B0121"
	"server/module/upload/project/B0122"
	"server/module/upload/project/guoshou"
	"server/utils"
	"time"
)

func BillUploadAdapter(reqParam model.ProCodeAndId, uploadPath model2.SysProDownloadPaths) error {
	lockVal := "locked:" + reqParam.ProCode + ":" + reqParam.ID
	global.GLog.Info("", zap.Any("lockVal", "lockVal"))
	acquired, err := utils.Lock(lockVal, 5*time.Minute)
	if err != nil {
		return err
	}
	if acquired {
		defer utils.Unlock(lockVal)
		switch reqParam.ProCode {
		case "B0108":
			return B0108.BillUpload(reqParam, uploadPath)
		case "B0113":
			return B0113.BillUpload(reqParam, uploadPath)
		case "B0121":
			return B0121.BillUpload(reqParam, uploadPath)
		case "B0122":
			return B0122.BillUpload(reqParam, uploadPath)
		case "B0114":
			return B0114.BillUpload(reqParam, uploadPath)
		case "B0103", "B0106", "B0110":
			return guoshou.BillUpload(reqParam, uploadPath)
		default:
			return project.BillUpload(reqParam, uploadPath)
		}
	}
	return errors.New("获取锁失败")
}
