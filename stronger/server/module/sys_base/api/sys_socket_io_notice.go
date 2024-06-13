/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/10/27 10:00
 */

package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"server/global"
	"server/global/response"
	"server/module/msg_manager/model"
	"server/module/sys_base/service"
	"strings"
	"time"
)

func CustomerNotice(c *gin.Context) {
	var reqParam struct {
		FileName string    `json:"fileName" form:"fileName"`
		SendTime time.Time `json:"sendTime" form:"sendTime"`
	}
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	isSend := global.GSocketIo.BroadcastToNamespace("/global-notice", "customerNotice", reqParam)
	if !isSend {
		response.FailWithMessage("发送失败", c)
		return
	}
	response.OkWithMessage("发送成功", c)
}

// BusinessPush 业务socket io通知推送
func BusinessPush(c *gin.Context) {
	var reqParam model.BusinessPush
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	reqParam.ID = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	global.GLog.Info("reqParam", zap.Any("s", reqParam))
	err, ids := service.GetUserIdsByMenu(reqParam.ProCode)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithParamErr(err, c)
		return
	}
	reqParam.CreatedAt = time.Now()
	reqParam.UpdatedAt = time.Now()
	sends := make([]model.BusinessPushSend, 0)
	for _, uid := range ids {
		send := model.BusinessPushSend{
			PushId: reqParam.ID, // 推送记录表
			UserId: uid,         // 推送到的人，用户id
			//ReadFlag: false,       // '阅读状态（f未读，t已读）',
		}
		send.CreatedAt = time.Now()
		send.UpdatedAt = time.Now()
		//fmt.Println(send)
		//fmt.Println(global.GSocketConnMsgMap[uid])
		if global.GSocketConnMsgMap[uid] != nil {
			fmt.Println(global.GSocketConnMsgMap[uid])
			eventName := []string{"download", "upload"}
			global.GSocketConnMsgMap[uid].Emit(eventName[reqParam.Type-1], map[string]interface{}{
				"businessPush":     reqParam,
				"businessPushSend": send,
			})
			send.IsPush = true
		}
		sends = append(sends, send)
	}
	err = service.SaveBusinessPush(reqParam, sends)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithParamErr(err, c)
		return
	}
}
