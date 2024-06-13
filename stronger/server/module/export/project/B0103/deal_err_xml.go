/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/10 11:16 上午
 */

package B0103

import (
	"encoding/json"
	"server/global"
	"server/module/pro_manager/model"
	"server/utils"
	"strings"

	"go.uber.org/zap"
)

func DealErrXml(xmlVal string, bill model.ProjectBill) (e error, x string) {
	global.GLog.Info("B0103:::DealErrXml")
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

	xmlVal = RegReplace(xmlVal, `<uedUnitNo \/>`, "")
	xmlVal = RegReplace(xmlVal, `<otherSpiltAmnt \/>`, "")
	xmlVal = RegReplace(xmlVal, `<medicalCategoryName \/>`, "")
	xmlVal = RegReplace(xmlVal, `<classifiedAmount \/>`, "")

	//catalogCode：用机构号匹配常量表《B0103_广西贵州国寿理赔_数据库编码对应表》中“本级机构代码”（第二列）对应的“医疗目录编码”（第四列）
	v, _ := utils.FetchConst(bill.ProCode, "B0103_广西贵州国寿理赔_数据库编码对应表", "医疗目录编码", map[string]string{"本级机构代码": bill.Agency})
	//constSpecialMap := constSpecialDeal(bill.ProCode)
	//bianMaItems := constSpecialMap["shuJuKuBianMaMap"][bill.Agency]
	//xmlVal = utils.SetNodeValue(xmlVal, "catalogCode", bianMaItems[3])
	xmlVal = utils.SetNodeValue(xmlVal, "catalogCode", v)
	return e, xmlVal
}
