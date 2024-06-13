/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021年11月08日11:23:27
 */

package const_data

////////////////////案件列表////////////////////

var BillStatus = map[int]string{
	1: "正常",
	2: "重复",
	3: "异常",
	4: "删除",
}

var BillInsuranceType = map[int]string{
	1: "医保",
	2: "非医保",
	3: "混合型",
}

var BillClaimType = map[int]string{
	-1: "",
	0:  "未定义",
	// 1: "医疗",
	// 2: "非医疗",
	3: "无发票",
	4: "无报销",
	5: "有报销",
	6: "混合型",
	7: "简易",
}

var BillStage = map[int]string{
	1: "待加载",
	2: "录入中",
	3: "已导出",
	4: "待审核",
	5: "已回传",
	6: "已完成",
	7: "已接收",
	8: "待下载",
}

var BillStickLevel = map[int]string{
	1: "紧急件",
	2: "优先件",
}

var OrderBy = map[string]string{
	"CreatedAt":    "created_at",
	"billNum":      "billNum",
	"agency":       "agency",
	"scanAt":       "scan_at",
	"exportAt":     "export_at",
	"lastUploadAt": "last_upload_at",
	"status":       "status",
	"stage":        "stage",
}

////////////////////案件列表////////////////////

////////////////////公告管理////////////////////

var AnnouncementType = map[int]string{
	1: "公告通知",
	2: "规则动态",
}

var AnnouncementStatus = map[int]string{
	1: "待发布",
	2: "已发布",
	3: "已删除",
}

////////////////////公告管理////////////////////

////////////////////特殊报表////////////////////

var ConstType = map[int]string{
	1: "医疗机构",
	2: "医疗目录",
}

////////////////////特殊报表////////////////////

////////////////////项目报表-业务明细////////////////////

var BillType = map[int]string{
	1: "门诊",
	2: "住院",
	3: "门诊，住院",
}

////////////////////项目报表-业务明细////////////////////

////////////////////消息管理--业务通知////////////////////

var BusinessMsgType = map[int]string{
	1: "下载",
	2: "回传",
}

////////////////////消息管理--业务通知////////////////////

// ContractClaimType ///////////////// 时效配置--新增合同时效 ////////////////
var ContractClaimType = map[int]string{
	-1: "",
	0:  "未定义",
	3:  "无发票",
	4:  "无报销",
	5:  "有报销",
	6:  "混合型",
	7:  "简易",
	8:  "全部",
}

/////////////////// 时效配置--新增合同时效 ////////////////
