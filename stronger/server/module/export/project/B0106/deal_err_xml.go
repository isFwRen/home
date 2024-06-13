/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/10 11:16 上午
 */

package B0106

import (
	"encoding/json"
	"server/global"
	"server/module/pro_manager/model"
	"server/utils"
	"strings"

	"go.uber.org/zap"
)

func DealErrXml(xmlVal string, bill model.ProjectBill) (e error, x string) {
	global.GLog.Info("B0106:::DealErrXml")
	bill_num := strings.Split(bill.BillNum, "_")
	xmlVal = SetNodeValue(xmlVal, "claimNo", bill_num[0])
	var otherInfo OtherInfo
	err := json.Unmarshal([]byte(bill.OtherInfo), &otherInfo)
	if err != nil {
		global.GLog.Error("otherInfo", zap.Error(err))
	}
	//utils.GetNodeValue(obj.Bill.OtherInfo, "claimTpaId")
	xmlVal = SetNodeValue(xmlVal, "claimTpaId", otherInfo.ClaimTpaId)

	// xmlVal = utils.SetNodeValue(xmlVal, "claimNo", bill.BillNum)
	xmlVal = utils.SetNodeValue(xmlVal, "bpoAbnormalReason", utils.RegReplace(bill.DelRemarks, ".*：", ""))
	bpoAbnormalReason := GetNodeValue(xmlVal, "bpoAbnormalReason")
	bpoAbnormalReason = utils.RegReplace(bpoAbnormalReason, `-.*$`, "")
	xmlVal = SetNodeValue(xmlVal, "bpoAbnormalReason", bpoAbnormalReason)

	// xmlVal = SetNodeValue(xmlVal, "accNum", "0")
	xmlVal = RegReplace(xmlVal, `<accNum \/>`, "<accNum>0</accNum>")

	xmlVal = RegReplace(xmlVal, `<uedUnitNo \/>`, "")
	xmlVal = RegReplace(xmlVal, `<otherSpiltAmnt \/>`, "")
	xmlVal = RegReplace(xmlVal, `<medicalCategoryName \/>`, "")
	xmlVal = RegReplace(xmlVal, `<classifiedAmount \/>`, "")
	//constSpecialMap := constSpecialDeal(bill.ProCode)
	v, _ := utils.FetchConst(bill.ProCode, "B0106_陕西国寿理赔_数据库编码对应表", "医疗目录编码", map[string]string{"本级机构代码": bill.Agency})
	//"删除案件在保单列表《处理》一栏提供《导出异常》按钮，点击后按照附件《异常件导出模板》导出xml、json，回传时json文件名需改为案件号（单号、问题件描述、项目名称需要代码单独定制）；
	//bpoAbnormalReason：默认为录入界面按F10键删除的备注信息（B0110新疆国寿有同样需求可参考）
	//catalogCode：匹配常量表《B0106_陕西国寿理赔_数据库编码对应表》中机构号对应的医疗目录编码"
	// xml = utils.SetNodeValue(xml, "bpoAbnormalReason", utils.RegReplace(bill.DelRemarks, ".*：", ""))
	//xmlVal = utils.SetNodeValue(xmlVal, "catalogCode", constSpecialMap["shuJuKuBianMa"][bill.Agency][3])
	xmlVal = utils.SetNodeValue(xmlVal, "catalogCode", v)
	return e, xmlVal
}
