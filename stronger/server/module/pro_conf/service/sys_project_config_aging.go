package service

import (
	"errors"
	"github.com/lib/pq"
	"server/global"
	"server/module/pro_conf/model"
	"server/module/pro_conf/model/request"
	"strings"
	"time"
)

func CreateAgingConfig() (err error) {
	var SysProjectConfigAging model.SysProjectConfigAgingReq
	err = global.GDb.Migrator().CreateTable(&SysProjectConfigAging)
	return err
}

func InsertAgingConfig(agingConfig model.SysProjectConfigAgingReq) (err error, configInter model.SysProjectConfigAgingReq) {
	var count int64
	var agingSave model.SysProjectConfigAging
	var agingCheck model.SysProjectConfigAging
	//时效内跟时效外配置的话不会出现交叉
	//配置多条时效设置的话，时效外配置需保持一致
	//时效外会跨天
	allowSave := true
	isFirstConfig := true
	if agingConfig.ConfigType == "download" {
		err = global.GDb.Model(&model.SysProjectConfigAging{}).Order("created_at ASC").Limit(1).
			Where("config_type = ? And pro_id = ? And node_name = ? And node_content = ? ", agingConfig.ConfigType, agingConfig.ProId, agingConfig.NodeName, agingConfig.NodeContent).
			Find(&agingCheck).
			Count(&count).Error
		if err != nil {
			return err, configInter
		}
	} else if agingConfig.ConfigType == "upload" {
		err = global.GDb.Model(&model.SysProjectConfigAging{}).Order("created_at ASC").Limit(1).
			Where("config_type = ? And pro_id = ? And field_name = ? And field_content = ? ", agingConfig.ConfigType, agingConfig.ProId, agingConfig.FieldName, agingConfig.FieldContent).
			Find(&agingCheck).
			Count(&count).Error
		if err != nil {
			return err, configInter
		}
	} else if agingConfig.ConfigType == "base" {
		err = global.GDb.Model(&model.SysProjectConfigAging{}).Order("created_at ASC").Limit(1).Where("config_type = ? And pro_id = ? ", agingConfig.ConfigType, agingConfig.ProId).
			Find(&agingCheck).Count(&count).Error
		if err != nil {
			return err, configInter
		}
	} else {
		return errors.New("请指明时效类型"), configInter
	}

	if count > 0 {
		isFirstConfig = false
		//判断时效内有无交叉
		err, includingStartTime := SVersionTwo(agingConfig.AgingStartTime, agingCheck.AgingStartTime, agingCheck.AgingEndTime)
		if err != nil {
			return err, model.SysProjectConfigAgingReq{}
		}
		err, includingEndTime := SVersionTwo(agingConfig.AgingEndTime, agingCheck.AgingStartTime, agingCheck.AgingEndTime)
		if err != nil {
			return err, model.SysProjectConfigAgingReq{}
		}
		if includingStartTime && includingEndTime {
			allowSave = false
		}
	}

	if allowSave && isFirstConfig {
		if agingConfig.AgingOutStartTime == "" || agingConfig.AgingOutEndTime == "" {
			return errors.New("该新增时效不符合要求! "), model.SysProjectConfigAgingReq{}
		}
		if agingConfig.AgingOutEndTime == "00:00:00" {
			agingConfig.AgingOutEndTime = "23:59:59"
		}
		if agingConfig.AgingEndTime == "00:00:00" {
			agingConfig.AgingEndTime = "23:59:59"
		}
		agingSave.ConfigType = agingConfig.ConfigType
		agingSave.FieldName = agingConfig.FieldName
		agingSave.FieldContent = agingConfig.FieldContent
		agingSave.NodeName = agingConfig.NodeName
		agingSave.NodeContent = agingConfig.NodeContent
		agingSave.AgingStartTime = agingConfig.AgingStartTime
		agingSave.AgingEndTime = agingConfig.AgingEndTime
		agingSave.RequirementsTime = agingConfig.RequirementsTime
		agingSave.ProId = agingConfig.ProId

		_, isSBigE := Compare(agingConfig.AgingOutStartTime, agingConfig.AgingOutEndTime)
		if isSBigE {
			agingSave.AgingOut = append(agingSave.AgingOut, agingConfig.AgingOutStartTime+"-23:59:59")
			agingSave.AgingOut = append(agingSave.AgingOut, "00:00:00-"+agingConfig.AgingOutEndTime)
		} else {
			agingSave.AgingOut = append(agingSave.AgingOut, agingConfig.AgingOutStartTime+"-"+agingConfig.AgingOutEndTime)
		}

		err = global.GDb.Model(&model.SysProjectConfigAging{}).Create(&agingSave).Error
		return err, agingConfig
	} else if allowSave && !isFirstConfig {
		agingSave.ConfigType = agingConfig.ConfigType
		agingSave.FieldName = agingConfig.FieldName
		agingSave.FieldContent = agingConfig.FieldContent
		agingSave.NodeName = agingConfig.NodeName
		agingSave.NodeContent = agingConfig.NodeContent
		agingSave.AgingStartTime = agingConfig.AgingStartTime
		agingSave.AgingEndTime = agingConfig.AgingEndTime
		agingSave.RequirementsTime = agingConfig.RequirementsTime
		agingSave.ProId = agingConfig.ProId
		agingSave.AgingOut = append(agingSave.AgingOut, agingCheck.AgingOut...)

		err = global.GDb.Model(&model.SysProjectConfigAging{}).Create(&agingSave).Error
		return err, agingConfig
	} else {
		return errors.New("该新增时效不符合要求! "), model.SysProjectConfigAgingReq{}
	}
}

func DelAgingConfig(ids []string) (err error) {
	var configAging model.SysProjectConfigAging
	for _, id := range ids {
		err = global.GDb.Where("id = ?", id).Delete(&configAging).Error
		if err != nil {
			return err
		}
	}
	return err
}

func UpdateProjectConfigAging(configAging model.SysProjectConfigAgingReq, id string) (err error) {
	var count int64
	if configAging.ConfigType != "download" && configAging.ConfigType != "upload" && configAging.ConfigType != "base" {
		return errors.New("请指明时效类型")
	}
	var total int64
	var firstProjectConfigAging model.SysProjectConfigAging

	db := global.GDb.Model(&model.SysProjectConfigAging{}).Order("created_at ASC").Limit(1)
	dbs := global.GDb.Model(&model.SysProjectConfigAging{})

	if configAging.ConfigType == "base" {
		db = db.Where("config_type = ? And pro_id = ? ", configAging.ConfigType, configAging.ProId)
		dbs = dbs.Where("config_type = ? And pro_id = ? ", configAging.ConfigType, configAging.ProId)
	} else if configAging.ConfigType == "download" {
		db = db.Where("config_type = ? And pro_id = ? And node_name = ? And node_content = ? ", configAging.ConfigType, configAging.ProId, configAging.NodeName, configAging.NodeContent)
		dbs = dbs.Where("config_type = ? And pro_id = ? And node_name = ? And node_content = ? ", configAging.ConfigType, configAging.ProId, configAging.NodeName, configAging.NodeContent)
	} else if configAging.ConfigType == "upload" {
		db = db.Where("config_type = ? And pro_id = ? And field_name = ? And field_content = ? ", configAging.ConfigType, configAging.ProId, configAging.FieldName, configAging.FieldContent)
		dbs = dbs.Where("config_type = ? And pro_id = ? And field_name = ? And field_content = ? ", configAging.ConfigType, configAging.ProId, configAging.FieldName, configAging.FieldContent)
	}
	err = db.Find(&firstProjectConfigAging).Count(&total).Error
	if err != nil {
		return err
	}
	if id == firstProjectConfigAging.ID {
		total = 0
	}
	agingOut := make(pq.StringArray, 0)
	if total > 0 {
		agingOut = append(agingOut, firstProjectConfigAging.AgingOut...)
	} else {
		if configAging.AgingOutEndTime == "00:00:00" {
			configAging.AgingOutEndTime = "23:59:59"
		}

		_, isSBigE := Compare(configAging.AgingOutStartTime, configAging.AgingOutEndTime)
		if isSBigE {
			agingOut = append(agingOut, configAging.AgingOutStartTime+"-23:59:59")
			agingOut = append(agingOut, "00:00:00-"+configAging.AgingOutEndTime)
		} else {
			agingOut = append(agingOut, configAging.AgingOutStartTime+"-"+configAging.AgingOutEndTime)
		}
		//如果修改的是最初那条时效配置, 那么就要更新同一类型的时效配置的时效外的时间
		var ag []model.SysProjectConfigAging
		err = dbs.Find(&ag).Error
		if err != nil {
			return err
		}
		for i, _ := range ag {
			ag[i].AgingOut = agingOut
			err = global.GDb.Model(&model.SysProjectConfigAging{}).Where("id = ? ", ag[i].ID).Updates(ag[i]).Error
			if err != nil {
				return err
			}
		}

	}

	err = global.GDb.Model(&model.SysProjectConfigAging{}).Where("id = ? And config_type = ? And pro_id = ? ", id, configAging.ConfigType, configAging.ProId).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("该时效设置不存在,无法编辑")
	}
	if configAging.AgingEndTime == "00:00:00" {
		configAging.AgingEndTime = "23:59:59"
	}

	err = global.GDb.Model(&model.SysProjectConfigAging{}).Where("id = ? And config_type = ? And pro_id = ? ", id, configAging.ConfigType, configAging.ProId).
		Updates(map[string]interface{}{
			"aging_start_time":  configAging.AgingStartTime,
			"aging_end_time":    configAging.AgingEndTime,
			"aging_out":         agingOut,
			"node_name":         configAging.NodeName,
			"node_content":      configAging.NodeContent,
			"field_name":        configAging.FieldName,
			"field_content":     configAging.FieldContent,
			"requirements_time": configAging.RequirementsTime,
		}).Error

	return err
}

func SelectAgingConfigByConfigType(req request.ReqProjectConfigAgingWithConfigType) (err error, config []model.SysProjectConfigAgingReq) {
	db := global.GDb.Model(&model.SysProjectConfigAging{}).Order("created_at ASC")
	var configs []model.SysProjectConfigAging
	var resitem model.SysProjectConfigAgingReq
	var resitems []model.SysProjectConfigAgingReq
	if req.ProId != "" && req.ConfigType != "" {
		err = db.Where("pro_id = ? And config_type = ?", req.ProId, req.ConfigType).Find(&configs).Error
		//beforeNodeOrFieldName := ""
		//beforeNodeOrFieldContent := ""
		for _, v := range configs {
			resitem.ID = v.ID
			resitem.ProId = v.ProId
			resitem.AgingStartTime = v.AgingStartTime
			resitem.AgingEndTime = v.AgingEndTime
			resitem.RequirementsTime = v.RequirementsTime
			resitem.ConfigType = v.ConfigType
			resitem.NodeName = v.NodeName
			resitem.NodeContent = v.NodeContent
			resitem.FieldName = v.FieldName
			resitem.FieldContent = v.FieldContent

			for j, ao := range v.AgingOut {
				arr := strings.Split(ao, "-")
				if j == 0 {
					resitem.AgingOutStartTime = arr[0]
				}
				resitem.AgingOutEndTime = arr[1]
			}

			//if i == 0 && req.ConfigType == "base" {
			//	for j, ao := range v.AgingOut {
			//		arr := strings.Split(ao, "-")
			//		if j == 0 {
			//			resitem.AgingOutStartTime = arr[0]
			//		}
			//		resitem.AgingOutEndTime = arr[1]
			//	}
			//}

			//if req.ConfigType == "download" {
			//	if beforeNodeOrFieldName != v.NodeName && beforeNodeOrFieldContent != v.NodeContent {
			//		for j, ao := range v.AgingOut {
			//			arr := strings.Split(ao, "-")
			//			if j == 0 {
			//				resitem.AgingOutStartTime = arr[0]
			//			}
			//			resitem.AgingOutEndTime = arr[1]
			//		}
			//		beforeNodeOrFieldName = v.NodeName
			//		beforeNodeOrFieldContent = v.NodeContent
			//	}
			//}

			//if req.ConfigType == "upload" {
			//	if beforeNodeOrFieldName != v.FieldName && beforeNodeOrFieldContent != v.FieldContent {
			//		for j, ao := range v.AgingOut {
			//			arr := strings.Split(ao, "-")
			//			if j == 0 {
			//				resitem.AgingOutStartTime = arr[0]
			//			}
			//			resitem.AgingOutEndTime = arr[1]
			//		}
			//		beforeNodeOrFieldName = v.FieldName
			//		beforeNodeOrFieldContent = v.FieldContent
			//	}
			//}

			resitems = append(resitems, resitem)
			//清脏数据
			//resitem.AgingOutStartTime = ""
			//resitem.AgingOutEndTime = ""
		}
		return err, resitems
	}
	return err, nil
}

// SVersionTwo 计算时间在哪个范围
func SVersionTwo(c, i, j string) (error, bool) {
	iT, cT, jT, err := Calculate(i, c, j)
	if err != nil {
		return err, false
	}
	if cT.After(iT) && cT.Before(jT) {
		return nil, true
	}
	return nil, false
}

func Calculate(start, mid, end string) (startT, midT, endT time.Time, err error) {
	local, err := time.LoadLocation("Local")
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, err
	}

	if start != "" {
		startT, err = time.ParseInLocation("15:04:05", start, local)
		if err != nil {
			return time.Time{}, time.Time{}, time.Time{}, err
		}
	}

	if mid != "" {
		midT, err = time.ParseInLocation("15:04:05", mid, local)
		if err != nil {
			return time.Time{}, time.Time{}, time.Time{}, err
		}
	}

	if end != "" {
		endT, err = time.ParseInLocation("15:04:05", end, local)
		if err != nil {
			return time.Time{}, time.Time{}, time.Time{}, err
		}
	}

	return startT, midT, endT, err
}

func Compare(c, i string) (error, bool) {
	cT, iT, _, err := Calculate(c, i, "")
	if err != nil {
		return err, false
	}
	if cT.Before(iT) {
		//c < i
		return nil, false
	}
	//c > i
	return nil, true
}
