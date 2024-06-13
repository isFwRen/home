/**
 * @Author: 阿先
 * @Description:
 * @Date: 2020/11/2 13:44
 */

package service

import (
	"fmt"
	dingtalk_robot "github.com/JetBlink/dingtalk-notify-go-sdk"
	"reflect"
	"server/global"
	"server/module/dingding/model"
	dingding2 "server/module/dingding/model/request"
	"strconv"
	"time"
)


func SendMsgToDingdingGroups(groupMsg model.DingdingGroupMsgs) (err error, groupInter model.DingdingGroupMsgs)  {
	fmt.Println(groupMsg)
	db := global.GDb.Model(&model.DingdingGroups{})
	dingdingGroupIds := groupMsg.DingdingGroupId
	var failGroups = ""
	fmt.Println("dingdingGroupIds", dingdingGroupIds)
	for _, groupId := range dingdingGroupIds {
		//根据id查询单个钉钉群记录，根据access_token发送消息
		var dingdingGroup model.DingdingGroups
		db.Where("id = ?", groupId).First(&dingdingGroup)
		fmt.Println("dingdingGroup：" + strconv.Itoa(int(groupId)), dingdingGroup)
		accessToken := dingdingGroup.AccessToken
		groupName := dingdingGroup.GroupName
		secret := dingdingGroup.Secret
		if &accessToken == nil || &secret == nil || accessToken == "" || secret == "" {
			if failGroups != "" {
				failGroups += "," + groupName
			} else {
				failGroups = groupName
			}
			continue
		}
		robot := dingtalk_robot.NewRobot(accessToken, secret)
		err := robot.SendTextMessage(groupMsg.SendMsg, []string{}, false); if err != nil{
			if failGroups != "" {
				failGroups += "," + groupName
			} else {
				failGroups = groupName
			}
		}
	}
	//SendStatus:1,success;2,exception
	if failGroups == "" {
		groupMsg.SendStatus = 1
	} else {
		groupMsg.SendStatus = 2
		groupMsg.FailReason = failGroups + "群发送失败，请联系开发人员！"
	}
	groupMsg.SendDate = time.Now()
	fmt.Println("插入", groupMsg)
	err = global.GDb.Model(groupMsg).Create(&groupMsg).Error
	//err = nil
	fmt.Println("err", err)
	return err, groupMsg
}

func SelectDingdingGroupMsgListByPage(info dingding2.OddsDingdingGroupMsgsStruct) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageIndex - 1)
	db := global.GDb.Model(&model.DingdingGroupMsgs{})

	t := reflect.TypeOf(info)
	if _, ok := t.FieldByName("StartTime"); ok {
		db = db.Where("sendDate >= ?", info.StartTime)
	}
	if _, ok := t.FieldByName("EndTime"); ok {
		db = db.Where("sendDate <= ?", info.EndTime)
	}
	var groupList []model.DingdingGroupMsgs
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Find(&groupList).Error
	return err, groupList, total
}
