/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/15 2:54 下午
 */

package utils

import (
	"server/module/pro_conf/model"
	"testing"
)

func TestXlsxExport(t *testing.T) {

	//var m model.SysExportNodeExport
	//list, err1 := serviceProConf.ExportNodeByExportId("907280101312299008")
	//if err1 != nil {
	//	panic(err1)
	//}

	list := []model.SysExportNodeExport{
		model.SysExportNodeExport{
			Name: "11111",
		},
		model.SysExportNodeExport{
			Name: "33333",
		},
	}
	ExportBigExcel("../files/", "Book.xlsx", "export", list)
}
