package model

import (
	"github.com/lib/pq"
	"reflect"
	"server/module/sys_base/model"
)

type ReportTag struct {
	model.Model
	ProjectCode string         `json:"projectCode"`                //项目编码
	UserId      string         `json:"userId"`                     //用户id
	TagList     pq.StringArray `json:"tagList" gorm:"type:text[]"` //存放BusinessDetailsExport json数组
}

func (r ReportTag) TableName() string {
	return "report_tags"
}

type TagLabel struct {
	Value string `json:"value"` //BusinessDetailsExport json
	Label string `json:"label"` //BusinessDetailsExport excel
}

func (r *ReportTag) GetTagLabelList() (list []TagLabel) {
	list = []TagLabel{}
	r.RefrushTagLabelMap()
	for _, s := range r.TagList {
		label, ok := TagLabelMap[s]
		if !ok {
			label = "未知字段"
		}
		list = append(list, TagLabel{Label: label, Value: s})
	}
	//list = TagLabelList
	return
}

func (r ReportTag) RefrushTagLabelMap() {
	bTmp := BusinessDetailsExport{}
	var getReflect interface{} = bTmp
	t := reflect.TypeOf(getReflect)
	if len(TagLabelMap) != t.NumField() {
		for i := 0; i < t.NumField(); i++ {
			TagLabelMap[t.Field(i).Tag.Get("json")] = t.Field(i).Tag.Get("excel")
			TagLabelList = append(TagLabelList, TagLabel{
				Label: t.Field(i).Tag.Get("excel"),
				Value: t.Field(i).Tag.Get("json")})
		}
	}
}

func (r *ReportTag) SortTagList() {
	if len(r.TagList) == 0 {
		return
	}
	mapTmp := make(map[string]string)
	for _, s := range r.TagList {
		mapTmp[s] = s
	}
	listTmp := []string{}
	for _, label := range TagLabelList {
		if value, ok := mapTmp[label.Value]; ok {
			listTmp = append(listTmp, value)
		}
	}
	r.TagList = listTmp
}

var TagLabelMap map[string]string = make(map[string]string)
var TagLabelList []TagLabel = []TagLabel{}
