/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/10 11:16 上午
 */

package B0118

import (
	"server/global"
	"server/module/pro_manager/model"
	"server/utils"
	"strings"
)

func DealErrXml(xml string, bill model.ProjectBill) (e error, x string) {
	global.GLog.Info("B0118:::DealErrXml")
	xml = strings.Replace(xml, "<claimNum></claimNum>", "<claimNum>"+bill.BillNum+"</claimNum>", 1)
	xml = strings.Replace(xml, "<errorDesc></errorDesc>", "<errorDesc>"+utils.RegReplace(RegReplace(bill.DelRemarks, ".*：", ""), `-.*$`, "")+"</errorDesc>", 1)
	return e, xml
}
