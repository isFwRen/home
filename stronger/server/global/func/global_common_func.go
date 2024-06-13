package _func

import (
	"fmt"
	"server/module/load/model"
	"time"
)

func StringToTime(t string) (T time.Time, err error) {
	local, err := time.LoadLocation("Local")
	if err != nil {
		return T, nil
	}
	T, err = time.ParseInLocation("20060102150405", t, local)
	if err != nil {
		return T, err
	}
	fmt.Println(T.Format("20060102150405"))
	return T, nil
}

func FieldsFormat(fields []model.ProjectField) [][]model.ProjectField {
	data := [][]model.ProjectField{}
	blockIndex := -1
	value := []model.ProjectField{}
	for _, field := range fields {
		if field.BlockIndex != blockIndex {
			if len(value) > 0 {
				data = append(data, value)
				value = []model.ProjectField{}
			}
			blockIndex = field.BlockIndex
		}
		value = append(value, field)
	}
	if len(value) > 0 {
		data = append(data, value)
	}
	return data
}
