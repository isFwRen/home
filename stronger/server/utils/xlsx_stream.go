/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/15 2:10 下午
 */

package utils

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"reflect"
	"strconv"
	"time"
)

// ExportBigExcel
// @Summary 导出接口
// @Param path		存放路径
// @Param bookName  文件名字
// @Param sheetName sheet名字
// @Param list      导出列表 （结构体切片，数组）

func ExportBigExcel(path string, bookName string, sheetName string, list interface{}) (err error) {

	//文件夹是否存在
	exists, err := PathExists(path)
	if err != nil {
		return err
	}
	if !exists {
		err = CreateDir(path)
		if err != nil {
			return err
		}
	}

	if reflect.ValueOf(list).IsNil() || reflect.ValueOf(list).IsZero() {
		return errors.New("数据为空")
	}

	file := excelize.NewFile()
	//file.NewSheet("hello") 新建sheet
	streamWriter, err := file.NewStreamWriter("Sheet1")
	if err != nil {
		return err
	}
	styleID, err := file.NewStyle(`{"font":{"color":"#000000"}}`)
	if err != nil {
		return err
	}

	switch reflect.TypeOf(list).Kind() {
	case reflect.Slice, reflect.Array:
		dataLen := reflect.ValueOf(list).Len()
		fmt.Println("数据长度", dataLen)
		for rowID := 0; rowID < dataLen; rowID++ {
			ele := reflect.ValueOf(list).Index(rowID)
			eleLens := ele.NumField()
			subLen := 0
			for i := 0; i < eleLens; i++ {
				//获取list的元素的第i个字段的tag的excel的值
				val := reflect.TypeOf(list).Elem().Field(i).Tag.Get("excel")
				if val == "" {
					subLen++
				}
			}
			//fmt.Println("1", subLen)
			eleLen := eleLens - subLen
			//fmt.Println("2", eleLen)
			//设置头
			headRow := make([]interface{}, eleLen)
			if rowID == 0 {
				fmt.Println("元素长度", eleLen)
				j := 0
				for i := 0; i < eleLens; i++ {
					//获取list的元素的第i个字段的tag的excel的值
					val := reflect.TypeOf(list).Elem().Field(i).Tag.Get("excel")
					if val == "" {
						continue
					}
					headRow[j] = excelize.Cell{StyleID: styleID, Value: val}
					j++
				}
				headCell, _ := excelize.CoordinatesToCellName(1, 1)
				if err = streamWriter.SetRow(headCell, headRow); err != nil {
					return err
				}
			}
			//写其他数据
			row := make([]interface{}, eleLen)
			k := 0
			for i := 0; i < eleLens; i++ {
				values := reflect.TypeOf(list).Elem().Field(i).Tag.Get("excel")
				excelFormat := reflect.TypeOf(list).Elem().Field(i).Tag.Get("excelFormat")
				if values == "" {
					continue
				}
				val := ele.Field(i)
				var valStr interface{}
				//valStr := val.String()
				//fmt.Println(reflect.TypeOf(list).Elem().Field(i).Type.String())
				//时间格式化
				switch reflect.TypeOf(list).Elem().Field(i).Type.String() {
				case "time.Time":
					if excelFormat == "" {
						//默认日期格式
						excelFormat = "2006-01-02 15-04-05"
					}
					valStr = val.Interface().(time.Time).Format(excelFormat)
				case "float64":
					valStr, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", val.Float()), 64)
				//fmt.Println(valStr)
				default:
					valStr = val
				}
				row[k] = valStr
				k++
			}
			cellData, _ := excelize.CoordinatesToCellName(1, rowID+2)
			if err = streamWriter.SetRow(cellData, row); err != nil {
				return err
			}
		}
	default:
		return errors.New("数据不是数组或切片")
	}

	//重名名sheet
	if sheetName != "" {
		file.SetSheetName("Sheet1", sheetName)
	}

	if err = streamWriter.Flush(); err != nil {
		return err
	}
	if err = file.SaveAs(path + bookName); err != nil {
		return err
	}
	return nil
}

// detail 产量明细
func ExportBigExcelDetail(path string, bookName string, sheetName string, list interface{}) (err error) {

	//文件夹是否存在
	exists, err := PathExists(path)
	if err != nil {
		return err
	}
	if !exists {
		err = CreateDir(path)
		if err != nil {
			return err
		}
	}

	if reflect.ValueOf(list).IsNil() || reflect.ValueOf(list).IsZero() {
		return errors.New("数据为空")
	}

	file := excelize.NewFile()
	//file.NewSheet("hello") 新建sheet
	streamWriter, err := file.NewStreamWriter("Sheet1")
	if err != nil {
		return err
	}
	styleID, err := file.NewStyle(`{"font":{"color":"#000000"}}`)
	if err != nil {
		return err
	}
	//合并单元格
	cellsArr := [][]string{
		{"A1", "A2"},
		{"B1", "B2"},
		{"C1", "C2"},
	}
	for _, elem := range cellsArr {
		streamWriter.MergeCell(elem[0], elem[1])
	}
	style, err := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	streamWriter.SetRow("A1", []interface{}{
		excelize.Cell{Value: "日期", StyleID: style},
	})
	streamWriter.SetRow("B1", []interface{}{
		excelize.Cell{Value: "工号", StyleID: style},
	})
	streamWriter.SetRow("C1", []interface{}{
		excelize.Cell{Value: "姓名", StyleID: style},
	})

	mergeCell := [][]string{
		{"D1", "L1", "汇总"},
		{"M1", "R1", "初审"},
		{"S1", "AB1", "一码"},
		{"AC1", "AL1", "二码"},
		{"AM1", "AV1", "问题件"},
	}
	for _, cell := range mergeCell {
		streamWriter.MergeCell(cell[0], cell[1])
		streamWriter.SetRow(cell[0], []interface{}{
			excelize.Cell{Value: cell[2], StyleID: style},
		})
	}
	//写入数据
	switch reflect.TypeOf(list).Kind() {
	case reflect.Slice, reflect.Array:
		dataLen := reflect.ValueOf(list).Len()
		fmt.Println("数据长度", dataLen)
		for rowID := 0; rowID < dataLen; rowID++ {
			ele := reflect.ValueOf(list).Index(rowID)
			eleLens := ele.NumField()
			subLen := 0
			for i := 0; i < eleLens; i++ {
				//获取list的元素的第i个字段的tag的excel的值
				val := reflect.TypeOf(list).Elem().Field(i).Tag.Get("excel")
				if val == "" {
					subLen++
				}
			}
			//fmt.Println("1", subLen)
			eleLen := eleLens - subLen
			//fmt.Println("2", eleLen)
			//设置头
			headRow := make([]interface{}, eleLen)
			if rowID == 0 {
				fmt.Println("元素长度", eleLen)
				j := 0
				for i := 0; i < eleLens; i++ {
					//获取list的元素的第i个字段的tag的excel的值
					val := reflect.TypeOf(list).Elem().Field(i).Tag.Get("excel")
					if val == "" {
						continue
					}
					headRow[j] = excelize.Cell{StyleID: styleID, Value: val}
					j++
				}
				headCell, _ := excelize.CoordinatesToCellName(1, 2)
				if err = streamWriter.SetRow(headCell, headRow); err != nil {
					return err
				}
			}
			//写其他数据
			row := make([]interface{}, eleLen)
			k := 0
			for i := 0; i < eleLens; i++ {
				values := reflect.TypeOf(list).Elem().Field(i).Tag.Get("excel")
				excelFormat := reflect.TypeOf(list).Elem().Field(i).Tag.Get("excelFormat")
				if values == "" {
					continue
				}
				val := ele.Field(i)
				var valStr interface{}
				//valStr := val.String()
				//fmt.Println(reflect.TypeOf(list).Elem().Field(i).Type.String())
				//时间格式化
				switch reflect.TypeOf(list).Elem().Field(i).Type.String() {
				case "time.Time":
					if excelFormat == "" {
						//默认日期格式
						excelFormat = "2006-01-02 15-04-05"
					}
					valStr = val.Interface().(time.Time).Format(excelFormat)
				case "float64":
					valStr, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", val.Float()), 64)
				case "int64":
					valStr = val.Int()
				//fmt.Println(valStr)
				default:
					valStr = val
				}
				row[k] = valStr
				k++
			}
			cellData, _ := excelize.CoordinatesToCellName(1, rowID+3)
			if err = streamWriter.SetRow(cellData, row); err != nil {
				return err
			}
		}
	default:
		return errors.New("数据不是数组或切片")
	}

	//重名名sheet
	if sheetName != "" {
		file.SetSheetName("Sheet1", sheetName)
	}

	if err = streamWriter.Flush(); err != nil {
		return err
	}
	if err = file.SaveAs(path + bookName); err != nil {
		return err
	}
	return nil
}

func ExportBigExcel2(path string, bookName string, sheetName string, list interface{}, sheet string) (err error) {

	//文件夹是否存在
	exists, err := PathExists(path)
	if err != nil {
		return err
	}
	if !exists {
		err = CreateDir(path)
		if err != nil {
			return err
		}
	}

	if reflect.ValueOf(list).IsNil() || reflect.ValueOf(list).IsZero() {
		return errors.New("数据为空")
	}

	var file *excelize.File
	//判断文件是否存在
	_, err = os.Stat(path + bookName)
	if os.IsNotExist(err) {
		file = excelize.NewFile()
	} else {
		file, err = excelize.OpenFile(path + bookName)
	}

	var streamWriter *excelize.StreamWriter
	fmt.Println("1", file.GetSheetIndex(sheetName))
	if file.GetSheetIndex(sheetName) == -1 {
		if file.GetSheetIndex(sheet) == -1 {
			_ = file.NewSheet(sheet)
		}
		streamWriter, err = file.NewStreamWriter(sheet)
		if err != nil {
			return err
		}
	} else {
		streamWriter, err = file.NewStreamWriter(sheetName)
		if err != nil {
			return err
		}
	}
	//file.NewSheet("hello") 新建sheet

	styleID, err := file.NewStyle(`{"font":{"color":"#000000"}}`)
	if err != nil {
		return err
	}

	//写入数据
	switch reflect.TypeOf(list).Kind() {
	case reflect.Slice, reflect.Array:
		dataLen := reflect.ValueOf(list).Len()
		fmt.Println("数据长度", dataLen)
		for rowID := 0; rowID < dataLen; rowID++ {
			ele := reflect.ValueOf(list).Index(rowID)
			eleLen := ele.NumField()

			//设置头
			headRow := make([]interface{}, eleLen)
			if rowID == 0 {
				fmt.Println("元素长度", eleLen)
				for i := 0; i < eleLen; i++ {
					//获取list的元素的第i个字段的tag的excel的值
					val := reflect.TypeOf(list).Elem().Field(i).Tag.Get("excel")
					headRow[i] = excelize.Cell{StyleID: styleID, Value: val}
				}
				headCell, _ := excelize.CoordinatesToCellName(1, 1)
				if err = streamWriter.SetRow(headCell, headRow); err != nil {
					return err
				}
			}

			//写其他数据
			row := make([]interface{}, eleLen)
			for i := 0; i < eleLen; i++ {
				val := ele.Field(i)
				row[i] = val
			}
			cellData, _ := excelize.CoordinatesToCellName(1, rowID+2)
			if err = streamWriter.SetRow(cellData, row); err != nil {
				return err
			}
		}
	default:
		return errors.New("数据不是数组或切片")
	}

	//重名名sheet
	if sheetName != "" {
		file.SetSheetName(sheet, sheetName)
	}

	if err = streamWriter.Flush(); err != nil {
		return err
	}
	if err = file.SaveAs(path + bookName); err != nil {
		return err
	}
	return nil
}

func ExportBigExcel3(path string, bookName string, sheetName string, list interface{}) (err error) {

	//文件夹是否存在
	exists, err := PathExists(path)
	if err != nil {
		return err
	}
	if !exists {
		err = CreateDir(path)
		if err != nil {
			return err
		}
	}

	if reflect.ValueOf(list).IsNil() || reflect.ValueOf(list).IsZero() {
		return errors.New("数据为空")
	}

	file := excelize.NewFile()
	//file.NewSheet("hello") 新建sheet
	streamWriter, err := file.NewStreamWriter("Sheet1")
	if err != nil {
		return err
	}
	styleID, err := file.NewStyle(`{"font":{"color":"#000000"}}`)
	if err != nil {
		return err
	}

	//写入数据
	switch reflect.TypeOf(list).Kind() {
	case reflect.Slice, reflect.Array:
		dataLen := reflect.ValueOf(list).Len()
		fmt.Println("数据长度", dataLen)
		for rowID := 0; rowID < dataLen; rowID++ {
			ele := reflect.ValueOf(list).Index(rowID)
			eleLens := ele.NumField()
			subLen := 0
			for i := 0; i < eleLens; i++ {
				//获取list的元素的第i个字段的tag的excel的值
				val := reflect.TypeOf(list).Elem().Field(i).Tag.Get("excel")
				if val == "" {
					subLen++
				}
			}
			//fmt.Println("1", subLen)
			eleLen := eleLens - subLen
			//fmt.Println("2", eleLen)
			//设置头
			headRow := make([]interface{}, eleLen)
			if rowID == 0 {
				fmt.Println("元素长度", eleLen)
				j := 0
				for i := 0; i < eleLens; i++ {
					//获取list的元素的第i个字段的tag的excel的值
					val := reflect.TypeOf(list).Elem().Field(i).Tag.Get("excel")
					if val == "" {
						continue
					}
					headRow[j] = excelize.Cell{StyleID: styleID, Value: val}
					j++
				}
				headCell, _ := excelize.CoordinatesToCellName(1, 1)
				if err = streamWriter.SetRow(headCell, headRow); err != nil {
					return err
				}
			}

			//写其他数据
			row := make([]interface{}, eleLen)
			k := 0
			for i := 0; i < eleLens; i++ {
				vals := reflect.TypeOf(list).Elem().Field(i).Tag.Get("excel")
				istime := reflect.TypeOf(list).Elem().Field(i).Tag.Get("time")
				if istime == "y" {
					val := ele.Field(i).Interface().(time.Time).Format("2006-01-02")
					if vals == "" {
						continue
					}
					row[k] = val
					k++
				} else {
					if vals == "" {
						continue
					}
					val := ele.Field(i)
					row[k] = val
					k++
				}
			}
			cellData, _ := excelize.CoordinatesToCellName(1, rowID+2)
			if err = streamWriter.SetRow(cellData, row); err != nil {
				return err
			}
		}
	default:
		return errors.New("数据不是数组或切片")
	}

	//重名名sheet
	if sheetName != "" {
		file.SetSheetName("Sheet1", sheetName)
	}

	if err = streamWriter.Flush(); err != nil {
		return err
	}
	if err = file.SaveAs(path + bookName); err != nil {
		return err
	}
	return nil
}
