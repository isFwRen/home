package utils

/**
 * @Author: 6727
 * @Description: 数组形式 excel导出
 * @Date: 2021/11/15 2:10 下午
 */

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"os"
	"server/global"
	"strconv"
)

func ExportExcelTheMainEntrance(excelDate [][]interface{}, ExcelName, ProCode, path string) (err error) {
	xlsxNew := excelize.NewFile()
	newSheet := "Sheet1"
	xlsxNew.NewSheet(newSheet)

	for i, row := range excelDate {
		for j, colCell := range row {
			zoophobia := ChangIndexToAxis(i, j)
			ModifyExcelCellByAxis(xlsxNew, newSheet, zoophobia, colCell)
		}
	}

	// 线上 uploads/file
	basicPath := global.GConfig.LocalUpload.FilePath + path + "/" + ProCode + "/"
	isExist, err := Exists(basicPath)
	if err != nil {
		return err
	}
	if !isExist {
		err = os.MkdirAll(basicPath, 0777)
		if err != nil {
			return err
		}
		err = os.Chmod(basicPath, 0777)
		if err != nil {
			return err
		}
	}

	err = xlsxNew.SaveAs(basicPath + "/" + ExcelName + ".xlsx")
	// 本地
	//err = xlsxNew.SaveAs("./" + ExcelName + ".xlsx")
	return err
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

/*
 * Title : 数组下标转换成excel坐标
 * Author : weiweiabc109@163.com
 * Date : 2018-04-06
 */
func ChangIndexToAxis(intIndexX int, intIndexY int) string {
	var arr = [...]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	intIndexY = intIndexY + 1
	resultY := ""
	for true {
		if intIndexY <= 26 {
			resultY = resultY + arr[intIndexY-1]
			break
		}
		mo := intIndexY % 26
		resultY = arr[mo-1] + resultY
		shang := intIndexY / 26
		if shang <= 26 {
			resultY = arr[shang-1] + resultY
			break
		}
		intIndexY = shang
	}
	return resultY + strconv.Itoa(intIndexX+1)
}

/*
 * Title : 修改excel表格里的值
 * Author : weiweiabc109@163.com
 * Date : 2018-04-06
 */
func ModifyExcelCellByAxis(xlsx *excelize.File, sheet string, axis string, value interface{}) int {
	xlsx.SetCellValue(sheet, axis, value)
	return 0
}

//func main() {
//	//go get github.com/360EntSecGroup-Skylar/excelize
//	xlsxNew := excelize.NewFile()
//	newSheet := "Sheet1"
//	xlsxNew.NewSheet(newSheet)
//
//	result := [5][5]int{{1, 2, 3, 4, 5}, {1, 2, 3, 4, 5}, {1, 2, 3, 4, 5}, {1, 2, 3, 4, 5}, {1, 2, 3, 4, 5}}
//
//	for i, row := range result {
//		//fmt.Printf("i:%d,",i)
//		for j, colCell := range row {
//			//fmt.Printf("j:%d",j)
//			zuobiao := ChangIndexToAxis(i, j)
//			ModifyExcelCellByAxis(xlsxNew, newSheet, zuobiao, colCell)
//			//fmt.Print(colCell, ",")
//		}
//	}
//	xlsxNew.SaveAs("./你要的表格.xlsx")
//}
