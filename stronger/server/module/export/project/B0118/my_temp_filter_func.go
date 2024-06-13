/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/26 5:12 下午
 */

package B0118

import (
	"fmt"
	"github.com/flosch/pongo2/v4"
	"reflect"
	"server/module/export/utils"
	model2 "server/module/load/model"
	"server/module/pro_conf/model"
	"strings"
)

func b0118TempFilterFunc(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	//global.GLog.Info("B0118:::MyTempFilterFunc")
	fCode := strings.Split(in.String(), ",")
	val := ""
	var fields []model2.ProjectField
	fmt.Println("118reflect.TypeOf(param.Interface()).Name():::", reflect.TypeOf(param.Interface()).Name())
	if reflect.TypeOf(param.Interface()) == nil {
		return pongo2.AsValue(val), nil
	}
	if reflect.TypeOf(param.Interface()).Name() == "FieldsMap" {
		fieldsMap := param.Interface().(FieldsMap)
		fields = fieldsMap.Fields
	} else {
		fieldsMap := param.Interface().(FormatObj)
		fields = fieldsMap.Fields
	}
	for _, item := range fCode {
		if item != "" {
			//todo 根据编码去拿结果值
			value, err := utils.GetFieldValue(fields, item)
			if err != nil {
				continue
			}
			val = value
		}
	}
	//global.GLog.Info(in.Interface())
	//global.GLog.Info(param.String())
	//return pongo2.AsValue(in.Interface()), nil
	return pongo2.AsValue(val), nil
}

func TempRender(export model.SysExport) (err error, tpl *pongo2.Template) {
	if !pongo2.FilterExists("B0118FilterFunc") {
		err = pongo2.RegisterFilter("B0118FilterFunc", b0118TempFilterFunc)
		if err != nil {
			return err, tpl
		}
	}
	tpl, err = pongo2.FromString(export.TempVal)
	if err != nil {
		return err, tpl
	}
	return err, tpl
}
