/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/5/15 16:37
 */

package B0106

import (
	"server/module/export/model"
	"server/module/export/service"
	"server/utils"
	"strings"
)

//报表管理-特殊报表-机构抽取：（rcptInfoList大节点下存在多个以下节点时，节点值默认以;隔开）
//机构：机构号
//核心立案号：单号
//疾病：mIcd10Code节点值匹配常量库《B0106_陕西国寿理赔_ICD10疾病编码》中的“疾病代码”，转换成“疾病名称”输出
//是否匹配：输出“是”或“否”，判断疾病是否匹配上常量库《ICD10疾病编码》
//票据类型：expenMode（0020：门诊 0040：住院 根据节点值转成对应类型）
//医院名称：空
//总金额：rcptAmnt
//统筹金额：socialPayAmnt
//范围外金额：空
//范围内金额：空

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
		name, total := utils.FetchConst(obj.Bill.ProCode, "B0106_陕西国寿理赔_ICD10疾病编码", "疾病名称", map[string]string{"疾病代码": code})
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
