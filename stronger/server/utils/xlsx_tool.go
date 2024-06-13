/**
 * @Author: 星期一
 * @Description:
 * @Date: 2020/12/28 下午3:49
 */

package utils

import (
	"errors"
	"github.com/tealeg/xlsx"
	"reflect"
	"server/module/pro_conf/model"
	"strings"
	"time"
)

func ExportExcel(interfase interface{}, sysExportNodeList []model.SysExportNodeExport) (err error, file *xlsx.File) {
	if len(sysExportNodeList) <= 0 {
		return errors.New("no data to generate xlsx!"), nil
	}

	//生成excel
	file = xlsx.NewFile()
	sheet, err := file.AddSheet("sheet1")
	if err != nil {
		panic(err)
	}

	//写入头
	row := sheet.AddRow()
	row.SetHeightCM(1) //设置每行的高度
	t := reflect.TypeOf(interfase)
	for i := 0; i < t.NumField(); i++ {
		h := t.Field(i).Tag.Get("gorm")
		cell := row.AddCell()
		val := strings.Replace(h, "comment:", "", -1)
		cell.Value = val
	}

	//写入每一行
	for _, item := range sysExportNodeList {
		row := sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度

		t := reflect.TypeOf(item)
		v := reflect.ValueOf(item)
		for i := 0; i < t.NumField(); i++ {
			//f := t.Field(i)
			val := v.Field(i)
			cell := row.AddCell()


			valInterface := val.Interface()
			//if i == 0 {
			//global.G_LOG.Info("++++++++++++++++++++++++++")
			//fmt.Println(reflect.TypeOf(valInterface))
			//fmt.Println(val.Type())
			//fmt.Println(val.Kind())
			switch valInterface.(type) {
			case time.Time:
				//fmt.Println(valInterface)
				timeLocationUTC, _ := time.LoadLocation("UTC")
				DefaultDateOptions := xlsx.DateTimeOptions{
					Location:        timeLocationUTC,
					ExcelTimeFormat: "m/d/yy h:mm",
				}
				cell.SetDateWithOptions(valInterface.(time.Time),DefaultDateOptions)
			default:
				cell.SetValue(val)
			}
			//global.G_LOG.Info("++++++++++++++++++++++++++")
			//}
		}
	}
	return err, file
}
