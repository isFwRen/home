package model

import (
	"fmt"
	"gorm.io/gorm"
	"server/module/sys_base/model"
	"server/utils"
	"strconv"
	"time"
)

type PtUserTrainingGuide struct {
	model.Model
	UserId        string `json:"userId"` // 用户id
	UserName      string `json:"userName"`
	ProjectCode   string `json:"projectCode"`    //项目代码
	TrainingStage int    `json:"training_stage"` //培训阶段 1.培训指引流程

	TrainStage `json:"trainStage"` //培训阶段
}

type TrainStage struct {
	Stage1   bool      `json:"stage1"`   // "培训指引流程学习"
	Stage1At time.Time `json:"stage1At"` // "培训指引流程学习完成时间"
	Stage2   bool      `json:"stage2"`   // "教学文件学习"
	Stage2At time.Time `json:"stage2At"` // "教学文件学习完成时间"
	Stage3   bool      `json:"stage3"`   // "录入练习"
	Stage3At time.Time `json:"stage3At"` // "录入练习完成时间"
	Stage4   bool      `json:"stage4"`   // "上岗考核"
	Stage4At time.Time `json:"stage4At"` // "上岗考核完成时间"
	Stage5   bool      `json:"stage5"`   // "上岗审核"
	Stage5At time.Time `json:"stage5At"` // "上岗审核完成时间"
	Stage6   bool      `json:"stage6"`   // "正式上岗"
	Stage6At time.Time `json:"stage6At"` // "正式上岗完成时间"
}

// 获取当前阶段
func (t TrainStage) GetCurrentStage() (int, error) {
	trainMap := utils.StructToMap(t)
	for i := 1; ; i++ {
		value, ok := trainMap["Stage"+strconv.Itoa(i)]

		b, err := strconv.ParseBool(fmt.Sprint(value))
		if b && ok && err == nil {
			continue
		} else if !ok {
			return 0, err
		} else if err != nil {
			return 0, err
		} else if !b {
			return i, nil
		}
	}
	return 0, nil
}

// NextStage 进入到下一阶段 stronger
/**
isAchieve:是否进入下阶段
currentMap:使用本方法时所处阶段
*/
func (p *PtUserTrainingGuide) NextStage(isAchieve bool, currentMap int, db *gorm.DB) (err error) {
	currentStage, err := p.TrainStage.GetCurrentStage()
	if isAchieve && currentStage == currentMap {
		err = db.Table("pt_user_training_guides").Where("training_stage = ? and id = ? ", currentStage, p.ID).
			Updates(map[string]interface{}{
				"stage" + strconv.Itoa(currentStage):         true,
				"stage" + strconv.Itoa(currentStage) + "_at": time.Now(),
				"training_stage":                             currentStage + 1,
			}).Error
	}

	return
}
