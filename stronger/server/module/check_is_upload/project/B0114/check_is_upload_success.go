/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/10/20 11:40
 */

package B0114

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"server/global"
	"server/module/download/project"
	"server/module/pro_manager/model"
	"strings"
)

// CheckIsUploadSuccess 检查客户路径是否还有
func CheckIsUploadSuccess(bill model.ProjectBill) (err error, isUploadSuccess bool) {
	cache, ok := global.GProConf[global.GConfig.System.ProCode]
	if !ok {
		err := errors.New("没有该项目[" + global.GConfig.System.ProCode + "]配置")
		global.GLog.Error("", zap.Error(err))
		return err, false
	}
	if cache.DownloadPaths.Scan == "" {
		err := errors.New("该项目[" + global.GConfig.System.ProCode + "]扫描配置为空配置")
		global.GLog.Error("", zap.Error(err))
		return err, false
	}
	global.GLog.Info("scanCmd", zap.Any("", cache.DownloadPaths.Scan))
	err, stdout, _ := project.ShellOut(cache.DownloadPaths.Scan)
	fmt.Println("stdout::" + stdout)
	if err != nil {
		global.GLog.Error("scan run cmd failed", zap.Error(err))
		return err, false
	}
	if strings.Index(stdout, bill.BillName+".zip") == -1 {
		return nil, true
	}
	return err, false
}
