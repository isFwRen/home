package utils

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"server/global"
	"server/module/pro_conf/model"
	"server/utils"
	"strconv"
	"strings"
)

// SendMsg 监控发送消息
func SendMsg(billNameArr []string, conf model.SysFtpMonitor) error {
	if len(billNameArr) > 0 {
		msg := fmt.Sprintf(conf.WrongMsg, strings.Join(billNameArr, ","))
		//发消息
		robot := utils.NewRobot("b72bc04ad782bdeea40828d064df3869fd49c4121ed0bbcbb1ec8295508c1b01", "SEC467698ef65ac33aa5e2f04e2c12eb9b5a55325446df8be33594927b91412ab30")
		err := robot.SendTextMessage(msg, []string{}, true)
		//更新消息
		if err != nil {
			global.GLog.Error("B0118 download Monitor的消息发送失败", zap.Error(err))
		}
		global.GLog.Info("B0118 download Monitor的消息发送成功", zap.Any("", msg))

		obj := map[string]interface{}{
			"title":   "下载通知",                        // 标题
			"msg":     msg,                           // 内容
			"type":    1,                             // 消息类型1:下载2:上传
			"stage":   0,                             // '状态（0，正常，1已删除）',
			"proCode": global.GConfig.System.ProCode, // 项目编码
		}
		global.GLog.Info("s", zap.Any("s", obj))
		//广播通知
		err1, res := utils.HttpRequest("http://localhost:"+strconv.Itoa(global.GConfig.System.CommonPort)+"/sys-socket-io-notice/business-push", obj)
		if err1 != nil {
			global.GLog.Error("通知失败::" + ":::" + err.Error())
		}
		global.GLog.Error("通知", zap.Any("res", res))
		global.GLog.Info("通知到/sys-socket-io-notice/business-push成功")
		return errors.New(err.Error() + "\n" + err1.Error())
	}
	return nil
}
