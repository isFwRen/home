/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/10/20 11:36
 */

package check_is_upload

import (
	"errors"
	"server/global"
	"server/module/check_is_upload/project/B0108"
	"server/module/check_is_upload/project/B0114"
	"server/module/check_is_upload/project/B0116"
	"server/module/check_is_upload/project/B0118"
	"server/module/pro_manager/model"
)

//CheckIsUploadAdapter 检查客户是否收到
func CheckIsUploadAdapter(proCode string, bill model.ProjectBill) (err error, isUploadSuccess bool) {
	if proCode == "" || global.ProCodeId[proCode] == "" {
		return errors.New("该单据项目编码错误,id::" + proCode), false
	}
	switch proCode {
	case "B0108":
		return B0108.CheckIsUploadSuccess(bill)
	case "B0114":
		return B0114.CheckIsUploadSuccess(bill)
	case "B0118":
		return B0118.CheckIsUploadSuccess(bill)
	case "B0116":
		return B0116.CheckIsUploadSuccess(bill)
	default:
		return errors.New("该项目没有自定义检查客户是否收到函数"), false
	}
}
