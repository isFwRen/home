package schedule

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"server/global"
	pf "server/module/pro_conf/model"
	billServer "server/module/pro_manager/service"
	"server/module/sys_base/service"
	"time"
)

// 获取时效简报
func GetTimeLinessBriefingCron() {
	//存入当前proCode
	//取出proCode  并且进行查询  固定
	code, err := global.GRedis.Get("ageing_briefing:code").Result()
	if err != nil {
		global.GLog.Error("项目编码有误请检查" + err.Error())
		return
	}
	if code == "" {
		global.GLog.Error("项目编码为空" + err.Error())
		return
	}
	//项目编码必须在数据库中
	var projectMap []pf.SysProject
	err = global.GDb.Model(&pf.SysProject{}).Find(&projectMap).Error
	if err != nil {
		global.GLog.Error("获取项目失败" + err.Error())
		return
	}
	billMap := make(map[string]string)
	for _, project := range projectMap {
		billMap[project.Code] = project.Code
	}
	if billMap[code] == "" {
		global.GLog.Error("项目编码有误，请重新输入" + err.Error())
		return
	}
	//获取项目的时效简报
	err, contractBillMap := billServer.GetTimeLinessBriefing(code)
	//存入redis中
	//转换JSON
	marshal, err := json.Marshal(contractBillMap)
	if err != nil {
		global.GLog.Error("转换失败，请检查数据是否有误" + err.Error())
		return
	}
	//存入把简报redis  并发送简报
	resultMap := global.GRedis.Set("ageing_briefing:"+code, marshal, 3600*time.Second)
	fmt.Println(resultMap)
	//取出redis  用scoket推送
	contractBills, err := global.GRedis.Get("ageing_briefing:" + code).Result()
	if err != nil {
		global.GLog.Error("项目编码有误请检查" + err.Error())
		return
	}

	//-*-----------------------------推送 socket
	err, ids := service.GetUserIdsByMenu(code)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		return
	}
	for _, uid := range ids {
		if global.GSocketConnMsgMap[uid] != nil {
			reqParam := "ageingBriefing:" + code
			global.GSocketConnMsgMap[uid].Emit("briefing", map[string]interface{}{
				"codePush":       reqParam,
				"ageingBriefing": contractBills,
			})
		}
	}
}
