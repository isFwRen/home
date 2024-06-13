/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/10 9:56 上午
 */

package export

import (
	"errors"
	"fmt"
	"io/ioutil"
	"server/global"
	"server/module/export/project"
	"server/module/pro_manager/model"
	"server/module/pro_manager/service"
	"time"

	"go.uber.org/zap"
)

// WrongBill 导出异常单
func WrongBill(reqParam model.ProCodeAndId) (err error) {
	//获取单据
	err, bill := service.GetProBillById(reqParam)
	if err != nil {
		return err
	}
	if bill.DelRemarks == "" {
		return errors.New("异常件删除原因为空时，不能导出异常，需重新加载；")
	}

	//拿到项目异常单模板
	data, err := ioutil.ReadFile("resource/err_template/" + reqParam.ProCode + "/err_bill.xml")
	if err != nil {
		global.GLog.Error("err_bill.xml File reading error", zap.Error(err))
		return err
	}
	fmt.Println(string(data))

	//根据单证信息写入异常单模板项目自定义处理
	err, xml := WrongBillExportAdapter(string(data), bill, reqParam.ProCode)
	if err != nil {
		return err
	}
	err = project.SaveXml(&bill, nil, xml, "utf-8", true)
	if err != nil {
		return err
	}

	//保存单据信息,设置单据为异常 录入状态为待审核 手动回传
	db := global.ProDbMap[reqParam.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	return db.Where("id = ?", bill.ID).UpdateColumns(model.ProjectBill{
		Status:       3,
		Stage:        3,
		IsAutoUpload: false,
		ExportAt:     time.Now(),
	}).Error
}
