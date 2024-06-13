package service

import (
	"reflect"
	"server/global"
	//"server/module/sys_base/model/request"

	//sys_base2 "server/module/sys_base/model/request"
	"server/module/training_guide/model"
	//model2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/model"
)

// AddAndFindTrainingStage
// 添加 或是 查找用户的 培训流程记录 如果用户从没有参加培训则为他创建一个
func AddAndFindTrainingStage(userID string, userName string) (userExist model.PtUserTrainingGuide, err error) {

	err = global.GDb.Table("pt_user_training_guides").Where("user_id = ?", userID).Scan(&userExist).Error
	if err != nil {
		return
	}
	if reflect.DeepEqual(userExist, model.PtUserTrainingGuide{}) { //查询结果为空不存在,添加新的记录
		userExist = model.PtUserTrainingGuide{
			UserId:        userID,
			UserName:      userName,
			ProjectCode:   global.GConfig.System.ProCode,
			TrainingStage: 1,
		}
		if err = global.GDb.Create(&userExist).Error; err != nil {
			return
		}
	} else {
		return
	}
	return
}

// UpdateTrainingStage
// 更新培训阶段
func UpdateTrainingStage(userID string, updateStage int) (err error) {
	var userExist model.PtUserTrainingGuide
	err = global.GDb.Table("pt_user_training_guides").Where("user_id = ?", userID).Scan(&userExist).Error
	if userExist.TrainingStage >= updateStage { //修改状态小于或等于(流程倒序)，则不执行更新
		return
	}
	if userExist.TrainingStage == updateStage-1 {
		if err = global.GDb.Table("pt_user_training_guides").
			Where("user_id = ?", userID).Update("training_stage", updateStage).Error; err != nil {
			return
		}
	}
	return

}

func GetTrainingStage(stageID string, primaryKey struct {
	UserID      string
	ProjectCode string
}) (stage model.PtUserTrainingGuide, err error) {

	if stageID != "" {
		err = global.GDb.Table("pt_user_training_guides").Where("id = ?", stageID).Scan(&stage).Error
		return
	} else if primaryKey.ProjectCode != "" && primaryKey.UserID != "" {
		err = global.GDb.Table("pt_user_training_guides").
			Where("project_code = ? and user_id = ?", primaryKey.ProjectCode, primaryKey.UserID).Scan(&stage).Error
		return
	}

	return
}

func SetTrainingStage4Practice(proCode, code string) (err error) {

	userId := ""
	err = global.GDb.Table("sys_users").
		Select("id").
		Where("code = ?", code).Scan(&userId).Error

	if userId == "" {
		return
	}
	stage, err := GetTrainingStage("", struct {
		UserID      string
		ProjectCode string
	}{UserID: userId, ProjectCode: proCode})

	if err = stage.NextStage(true, 3, global.GDb); err != nil {
		return
	}

	return
}
