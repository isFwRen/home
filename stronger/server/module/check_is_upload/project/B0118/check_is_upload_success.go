/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/10/20 11:40
 */

package B0118

import (
	"errors"
	"fmt"
	"server/global"
	"server/module/download/project"
	"server/module/pro_manager/model"
	"strings"

	"go.uber.org/zap"
)

// CheckIsUploadSuccess 检查客户路径是否还有
func CheckIsUploadSuccess(bill model.ProjectBill) (err error, isUploadSuccess bool) {
	cache, ok := global.GProConf[global.GConfig.System.ProCode]
	if !ok {
		err := errors.New("没有该项目[" + global.GConfig.System.ProCode + "]配置")
		global.GLog.Error("", zap.Error(err))
		return err, false
	}
	cache.DownloadPaths.Scan = "docker exec -it download_system /usr/bin/curl -u ICSBPOzhuhaihuiliu:ICSBPOzhuhaihuiliu#123 -f -s -k ftps://202.108.196.190:991/input/"
	//cache.DownloadPaths.Scan = "curl -u ICSBPOzhuhaihuiliu:ICSBPOzhuhaihuiliu#123 -f -s -k ftps://202.108.196.190:991/input/"
	if cache.DownloadPaths.Scan == "" {
		err := errors.New("该项目[" + global.GConfig.System.ProCode + "]扫描配置为空配置")
		global.GLog.Error("", zap.Error(err))
		return err, false
	}
	global.GLog.Info("scanCmd", zap.Any("", cache.DownloadPaths.Scan))
	err, stdout, _ := project.SpecialShell(cache.DownloadPaths.Scan)
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
