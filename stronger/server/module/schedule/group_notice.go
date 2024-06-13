/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/26 14:58
 */

package schedule

import (
	"fmt"
	"server/global"
	"server/module/load/service"
	"server/module/msg_manager/model"
	service2 "server/module/msg_manager/service"
	"server/utils"
)

//SendMsg 发送通知
func SendMsg(notice model.GroupNoticeTwo) {
	fmt.Println(notice)
	global.GLog.Info("发送通知【" + notice.ProCode + "】")
	//[单量通知]目前B0101项目一码单量150、报销单50；二码单量120、报销单34，请安排时间上线处理。
	db := global.ProDbMap[notice.ProCode+"_task"]
	if db == nil {
		global.GLog.Error("没有该项目【" + notice.ProCode + "】的连接")
		return
	}
	_, op1Count := service.CountBlockOpNum(notice.ProCode+"_task", "op1")
	_, op2Count := service.CountBlockOpNum(notice.ProCode+"_task", "op2")
	_, op1CountBXD := service.CountBlockOpNumAndName(notice.ProCode+"_task", "op1", "报销单")
	_, op2CountBXD := service.CountBlockOpNumAndName(notice.ProCode+"_task", "op1", "报销单")

	msg := fmt.Sprintf("[单量通知]目前%s项目一码单量%d、报销单%d；二码单量%d、报销单%d，请安排时间上线处理。",
		notice.ProCode, op1Count, op1CountBXD, op2Count, op2CountBXD)

	err, group := service2.GetDingtalkGroupById(notice.GroupId)
	if err != nil {
		global.GLog.Error("没有该群【" + notice.ProCode + "-" + notice.GroupId + "】的群机器人")
		return
	}
	robot := utils.NewRobot(group.AccessToken, group.Secret)
	err = robot.SendTextMessage(msg, []string{}, true)
	if err != nil {
		global.GLog.Error("【" + notice.ProCode + "-" + notice.GroupId + "】定时发送失败")
		return
	}
}
