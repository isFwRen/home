/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/10 11:16 上午
 */

package B0110

import (
	"server/global"
	"server/module/pro_manager/model"
	"server/utils"
	"strings"
)

func DealErrXml(xml string, bill model.ProjectBill) (e error, x string) {
	global.GLog.Info("B0110:::DealErrXml")
	//constSpecialMap := constSpecialDeal(bill.ProCode)
	//删除案件在保单列表《处理》一栏提供《导出异常》按钮，点击后按照附件《异常件导出模板》导出json，
	//回传时json文件名需改为案件号；
	//bpoAbnormalReason：默认为F9删除的备注信息
	//hospitalCode：默认为6500000000000349
	//mIcd10Code：默认为R52.9
	//删除bankAccInfoList、errorRcptList整个节点，包含里面的小节点
	//accNum：0
	//catalogCode：用机构号匹配常量表《B0110_新疆国寿理赔_数据库编码对应表》中“本级机构代码”（第二列）对应的“医疗目录编码”（第四列）
	xml = SetNodeValue(xml, "claimNo", strings.Split(bill.BillNum, "_")[0])
	xml = SetNodeValue(xml, "claimTpaId", bill.BatchNum)
	xml = SetNodeValue(xml, "bpoAbnormalReason", strings.Split(utils.RegReplace(bill.DelRemarks, ".*：", ""), "-")[0])
	xml = SetNodeValue(xml, "hospitalCode", "6500000000000349")
	xml = SetNodeValue(xml, "mIcd10Code", "R52.9")
	xml = SetNodeValue(xml, "accNum", "0")
	v, _ := utils.FetchConst(bill.ProCode, "B0110_新疆国寿理赔_数据库编码对应表", "医疗目录编码", map[string]string{"本级机构代码": bill.Agency})
	xml = utils.SetNodeValue(xml, "catalogCode", v)
	//xml = utils.SetNodeValue(xml, "catalogCode", constSpecialMap["shuJuKuBianMa"][bill.Agency][3])

	// if item, ok := constSpecialMap["shuJuKuBianMa"]["医疗目录"+bill.Agency]; ok {
	// 	xml = SetNodeValue(xml, "catalogCode", item[3])
	// }
	return e, xml
}
