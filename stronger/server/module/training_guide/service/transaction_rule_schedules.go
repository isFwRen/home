package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"server/global"
	needReadModel "server/module/pro_manager/model"
	baseModel "server/module/sys_base/model"
	"server/module/training_guide/model"
	"server/module/training_guide/request"
	"server/utils"
	"strings"
	model3 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/model"
	//model2 "server/module/sys_base/model"
)

func FinishStep(param request.FinishParam, userInfo *model3.CustomClaims) (err error) {

	db, ok := global.ProDbMap[param.ProjectCode]
	if !ok {
		err = errors.New("不存在" + param.ProjectCode + "项目数据库连接")
		return
	}
	var ruleSchedule model.RuleSchedule
	err = db.Table("rule_schedules").
		Where("user_id = ? and project_code = ? and train_type = ?", userInfo.ID, param.ProjectCode, param.TrainType).
		Scan(&ruleSchedule).Error
	if err != nil {
		return
	}
	if reflect.DeepEqual(ruleSchedule, model.RuleSchedule{}) {
		//不存在记录 第一次点击进入
		db2GetList := db
		if param.TrainType == 3 {
			db2GetList = global.GDb
		}

		list, err1 := GetNeedReadList(param.TrainType, db2GetList)
		if err1 != nil {
			err = err1
			return
		}
		newRuleSchedule := model.RuleSchedule{
			UserId:            userInfo.ID,
			ProjectCode:       param.ProjectCode,
			TrainType:         param.TrainType,
			RequiredLearnList: list,
			NeedLearnList:     list,
		}
		err = db.Model(model.RuleSchedule{}).Create(&newRuleSchedule).Error

		err = db.Table("rule_schedules").
			Where("user_id = ? and project_code = ? and train_type = ?", userInfo.ID, param.ProjectCode, param.TrainType).
			Scan(&ruleSchedule).Error
		if err != nil {
			return
		}
		err = errors.New("用户不存在记录")
		//return
	}

	if isExist, newSlice := utils.RemoveElem(ruleSchedule.NeedLearnList, param.FinishId); isExist {

		updateString := `{` + strings.Join(newSlice, ", ") + `}`
		err = db.Table("rule_schedules").Where("id = ?", ruleSchedule.ID).
			Update("need_learn_list", updateString).Error

	} else {
		err = errors.New("已完成该任务")
		return
	}

	return
}

func GetNeedReadList(TranType int, db *gorm.DB) (list []string, err error) {
	switch TranType {
	case 1:
		err = db.Model(&needReadModel.TransactionRule{}).Select("id").
			Where("is_required = 1").Scan(&list).Error
	case 2:
		err = db.Model(&needReadModel.PicInformation{}).Select("id").
			Where("is_required = 1").Scan(&list).Error
	case 3:
		err = db.Model(&needReadModel.TeachVideo{}).Select("id").
			Where("is_required = 1").Scan(&list).Error
	default:
		list = []string{}
	}
	return
}

// 检查是否所有文件都学习完成

func CheckFinishAll(projectCode string, userInfo *model3.CustomClaims) (isFinish bool, err error) {
	db, ok := global.ProDbMap[projectCode]
	if !ok {
		err = errors.New("不存在" + projectCode + "项目数据库连接")
		return
	}

	isFinish = true

	ruleArray := []model.RuleSchedule{}
	err = db.Table("rule_schedules").
		Where("user_id = ? and train_type in (1,2,3)", userInfo.ID).Scan(&ruleArray).Error

	for _, schedule := range ruleArray {
		isFinish = isFinish && len(schedule.NeedLearnList) == 0
	}

	//for i, _ := range model.TrainTypeMap {
	//	var ruleItem model.RuleSchedule
	//	err = db.Table("rule_schedules").
	//		Where("user_id = ? and train_type = ?", userInfo.ID, i).Scan(&ruleItem).Error
	//	if err != nil {
	//		return
	//	}
	//	isFinish = isFinish && ruleItem.NoNeedLearn()
	//}
	return
}

func InsertIsReadStatus(needReadList []string, list []map[string]interface{}) (convMap []map[string]interface{}, err error) {
	convMap = []map[string]interface{}{}
	for _, m := range list {
		modelId, isExist := m["model"]
		value, ok := m["isRequired"]
		m["isLearned"] = 0
		if ok && isExist && len(needReadList) != 0 {
			modelStruct, convSuccess := modelId.(baseModel.Model)
			valueInt, convInt := value.(int)
			if convSuccess && convInt && valueInt == 1 {
				if utils.IsContain(needReadList, modelStruct.ID) {
					m["isLearned"] = 0
				} else {
					m["isLearned"] = 1
				}
			}

		}
		convMap = append(convMap, m)
	}

	return
}

func InsertReadStatus(needReadList []string, list interface{}) (convMap []map[string]interface{}, err error) {
	//a := list.([]interface{})
	//for _, m := range a {
	//	fmt.Println(m)
	//}

	convMap = []map[string]interface{}{}
	tList := reflect.TypeOf(list)

	switch tList.Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(list)
		rr, ok := list.([]struct{})
		if ok {
			global.GLog.Info(fmt.Sprintf("%v", rr))
			for _, s2 := range rr {
				itemMaprr := utils.StructToMap(s2)
				global.GLog.Info(fmt.Sprintf("%v", itemMaprr))
			}

		}
		//for _, item := range list.([]interface{}) {
		for i := 0; i < s.Len(); i++ {
			item := s.Index(i)

			v1 := item.FieldByName("ProCode").Interface()
			fmt.Println(v1)
			itemMap := make(map[string]interface{})
			//v := reflect.ValueOf(item)
			//t := reflect.TypeOf(item)
			//itemMap = utils.StructToMap(item)
			//for i := 0; i < item.NumField(); i++ {
			//	//itemMap[t.Field(i).Tag.Get("json")] = v.Field(i)
			//	global.GLog.Info(fmt.Sprintf("%v", v.FieldByName(t.Field(i).Name).Int()))
			//	itemMap[t.Field(i).Name] = v.FieldByName(t.Field(i).Name).Interface()
			//
			//	if t.Field(i).Tag.Get("json") == "isRequired" &&
			//		v.FieldByName(t.Field(i).Name).Int() == 1 { //找到必学字段 判断该文件是否必学 如果必学查看是否已学 添加已学字段
			//
			//		global.GLog.Info(v.FieldByName("ID").String())
			//		if utils.IsContain(needReadList, v.FieldByName("ID").String()) {
			//			itemMap["isLearned"] = 0
			//		} else {
			//			itemMap["isLearned"] = 1
			//		}
			//	}
			//}
			convMap = append(convMap, itemMap)

		}
	default:
		global.GLog.Info("no array!!!!")
		return nil, err
	}

	//if tList.Kind() == reflect.Array {
	//	for _, item := range list {
	//		itemMap := make(map[string]interface{})
	//		v := reflect.ValueOf(item)
	//		t := reflect.TypeOf(item)
	//		itemMap = utils.StructToMap(item)
	//		for i := 0; i < t.NumField(); i++ {
	//			//itemMap[t.Field(i).Tag.Get("json")] = v.Field(i)
	//			if t.Field(i).Tag.Get("json") == "isRequired" &&
	//				v.FieldByName(t.Field(i).Name).Int() == 1 { //找到必学字段 判断该文件是否必学 如果必学查看是否已学 添加已学字段
	//
	//				if utils.IsContain(needReadList, v.FieldByName("ID").String()) {
	//					itemMap["isLearned"] = 0
	//				} else {
	//					itemMap["isLearned"] = 1
	//				}
	//			}
	//		}
	//		convMap = append(convMap, itemMap)
	//	}
	//}

	return
}

func GetNeedReadList2(ruleType int, userId string) (list []string, err error) {

	db, ok := global.ProDbMap[global.GConfig.System.ProCode]

	list = []string{}
	if ok {
		ruleSchedule := model.RuleSchedule{}
		err = db.Table("rule_schedules").
			Where("user_id = ? and project_code = ? and train_type = ?", userId, global.GConfig.System.ProCode, ruleType).
			Scan(&ruleSchedule).Error

		list = ruleSchedule.NeedLearnList
	} else {
		global.GLog.Error("数据库" + global.GConfig.System.ProCode + "连接不存在")
	}

	return
}
