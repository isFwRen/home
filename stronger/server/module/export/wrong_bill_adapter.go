/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/11 9:14 上午
 */

package export

import (
	"errors"
	"server/module/export/project/B0103"
	"server/module/export/project/B0106"
	"server/module/export/project/B0110"
	"server/module/export/project/B0118"
	"server/module/pro_manager/model"
)

//WrongBillExportAdapter 自定义导出异常中转站
func WrongBillExportAdapter(xml string, bill model.ProjectBill, code string) (e error, x string) {
	switch code {
	case "B0118":
		return B0118.DealErrXml(xml, bill)
	case "B0110":
		return B0110.DealErrXml(xml, bill)
	case "B0106":
		return B0106.DealErrXml(xml, bill)
	case "B0103":
		return B0103.DealErrXml(xml, bill)
	default:
		return errors.New("该项目没有自定义导出异常单处理"), xml
	}
}
