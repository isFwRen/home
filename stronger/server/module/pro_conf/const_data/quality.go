/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/12/13 3:54 下午
 */

package const_data

var QualityInputType = map[int]string{
	1: "输入框",
	2: "常量",
	3: "下拉",
}

var QualityBelongType = map[int]string{
	1: "申请人信息",
	2: "被保人信息",
	3: "受托人信息",
	4: "其他信息",
	5: "受益人信息",
	6: "领款人信息",
	7: "账单信息",
	8: "出险信息",
}

var QualityBeneficiary = map[int]string{
	1: "受益人姓名",
	2: "领款人姓名",
	3: "其他信息",
}

var QualityBillInfo = map[int]string{
	1:  "账单号",
	2:  "账单类型",
	3:  "医院名称",
	4:  "开始时间",
	5:  "结束时间",
	6:  "天数",
	7:  "账单金额",
	8:  "调整金额",
	9:  "自费金额",
	10: "自付金额",
	11: "报销金额",
}
