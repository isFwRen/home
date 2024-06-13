/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/8/30 15:07
 */

package monitor

import (
	"errors"
	"server/module/monitor/project/B0108"
	"server/module/monitor/project/B0113"
	"server/module/monitor/project/B0114"
	"server/module/monitor/project/B0118"
	"server/module/monitor/project/B0121"
	"server/module/monitor/project/B0122"
	proConf "server/module/pro_conf/model"
)

// AdapterMonitor 监控
func AdapterMonitor(conf proConf.SysFtpMonitor) error {
	switch conf.ProCode {
	case "B0108":
		err := B0108.ScanDownload()
		if err != nil {
			return err
		}
		err = B0108.ScanUpload()
		//time.Sleep(5 * time.Minute)
		return err
	case "B0113":
		return B0113.ScanDownload(conf)
	case "B0114":
		return B0114.ScanDownload(conf)
	case "B0118":
		return B0118.ScanDownload(conf)
	case "B0121":
		return B0121.ScanDownload(conf)
	case "B0122":
		return B0122.ScanDownload(conf)
	default:
		return errors.New("该项目没有自定义导出异常单处理")
	}
}
