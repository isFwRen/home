/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/12/13 11:20 上午
 */

package api

import (
	"fmt"
	"io/ioutil"
	"server/global"
	"server/module/export/project"
	"server/module/pro_conf/model"
	model2 "server/module/pro_manager/model"
	"server/utils"

	"go.uber.org/zap"
)

// formatBill 格式化查看结果数据xml质检
func formatBill(qualities []model.SysQuality, bill model2.ProjectBill) (err error, formatQuality map[int][]model.SysQuality, amount model2.BillAmount) {

	xmlFile := global.GConfig.LocalUpload.FilePath + bill.ProCode + "/upload_xml/" +
		fmt.Sprintf("%v/%v/%v/%v.xml",
			bill.CreatedAt.Year(), int(bill.CreatedAt.Month()),
			bill.CreatedAt.Day(), bill.BillNum)
	global.GLog.Info("xml file:::" + xmlFile)

	//拿到项目回传单据xml字符串
	data, err := ioutil.ReadFile(xmlFile)
	if err != nil {
		global.GLog.Error(xmlFile+" File reading error", zap.Error(err))
		return err, formatQuality, amount
	}
	//global.GLog.Info(string(data))

	for i, quality := range qualities {
		arr := utils.GetNode(string(data), quality.ParentXmlNodeName)
		qualities[i].XmlNodeVal = []string{}
		for _, s := range arr {
			dataArr := utils.GetNodeData(s, quality.XmlNodeName)
			qualities[i].XmlNodeVal = append(qualities[i].XmlNodeVal, dataArr...)
			//fmt.Println("111")
		}
	}

	formatQuality = make(map[int][]model.SysQuality, 0)
	for _, quality := range qualities {
		formatQuality[quality.BelongType] = append(formatQuality[quality.BelongType], quality)
	}

	err, formatQuality, amount = project.BillGetTotalFeeAdapter(bill.ProCode, formatQuality, string(data))

	return err, formatQuality, amount
}
