/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/5/15 16:37
 */

package B0110

import (
	"server/module/export/model"
	"server/module/export/service"
	"server/utils"
	"strings"
)

var expendModeMap = map[string]string{"0020": "门诊", "0040": "住院"}
var isMatchMap = map[bool]string{true: "是", false: "否"}

func FetchAgency(obj model.ResultDataBill, xmlValue string) error {
	//常量
	//constMap := constDeal(obj.Bill.ProCode)

	mIcd10CodeName := ""
	isMatch := ""
	mIcd10Codes := utils.GetNodeData(xmlValue, "mIcd10Code")
	for _, code := range mIcd10Codes {
		//name, ok := constMap["jiBingBianMaCodeToNameMap"][code]
		name, total := utils.FetchConst(obj.Bill.ProCode, "B0110_新疆国寿理赔_ICD10疾病编码", "疾病名称", map[string]string{"疾病代码": code})
		mIcd10CodeName += name + ";"
		//isMatch += isMatchMap[ok] + ";"
		isMatch += isMatchMap[total != 0] + ";"
	}

	expenModeName := ""
	expenModes := utils.GetNodeData(xmlValue, "expenMode")
	for _, mode := range expenModes {
		expenModeName += expendModeMap[mode] + ";"
	}

	rcptAmnt := strings.Join(utils.GetNodeData(xmlValue, "rcptAmnt"), ";")
	socialPayAmnt := strings.Join(utils.GetNodeData(xmlValue, "socialPayAmnt"), ";")

	agency := model.Agency{
		BillId:         obj.Bill.ID,
		BillNum:        obj.Bill.BillNum,
		MIcd10Code:     mIcd10CodeName,
		Agency:         obj.Bill.Agency,
		IsMatch:        isMatch,
		ExpenMode:      expenModeName,
		CountMoney:     rcptAmnt,
		SocialPayMoney: socialPayAmnt,
	}

	return service.UpdateAgency(obj.Bill.ProCode, agency)
}
