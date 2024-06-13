package service

import (
	"reflect"
	"server/global"
	"server/module/report_management/model"
	reqModel "server/module/report_management/model/request"
)

func GetTagList() (list []model.TagLabel) {
	list = []model.TagLabel{}
	model.ReportTag{}.RefrushTagLabelMap()
	//for key, value := range model.TagLabelMap {
	//	list = append(list, model.TagLabel{
	//		Value: key,
	//		Label: value,
	//	})
	//}
	list = model.TagLabelList
	return
}

func GetUserTags(userId string, param reqModel.ReqGetUserTagList) (list []model.TagLabel, err error) {
	list = []model.TagLabel{}
	var reportTag model.ReportTag
	err = global.GDb.Table("report_tags").
		Where("user_id = ? and project_code = ?", userId, param.ProjectCode).
		First(&reportTag).Error
	if err != nil {
		return
	}
	if !reflect.DeepEqual(reportTag, model.ReportTag{}) {
		list = reportTag.GetTagLabelList()
	}
	return
}

func SetUserTags(userId string, param reqModel.ReqSettingReport) (err error) {
	reportTag := model.ReportTag{UserId: userId, ProjectCode: param.ProjectCode}
	reportTag.TagList = param.TagsList
	reportTag.SortTagList() //将新添加的表头按顺序插入
	recordId := ""
	err = global.GDb.Table("report_tags").Select("id").
		Where("project_code = ? and user_id = ?", param.ProjectCode, userId).Scan(&recordId).Error
	if err != nil {
		return
	}
	if recordId != "" {
		reportTag.ID = recordId
	}
	err = global.GDb.Table("report_tags").
		Where("project_code = ? and user_id = ?", param.ProjectCode, userId).Save(&reportTag).Error

	return
}
