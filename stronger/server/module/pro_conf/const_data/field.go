/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/10/27 4:59 下午
 */

package const_data

var FieldValidation = map[int]string{
	1:  "必录",
	2:  "允许空格",
	3:  "显示长度",
	4:  "数字",
	5:  "金额",
	6:  "中文",
	7:  "不能录负数",
	8:  "字母",
	9:  "整数",
	10: "年月日",
	11: "身份证",
	12: "邮件",
	13: "手机",
	14: "姓名",
}

var FieldProcess = map[int]string{
	1: "不录",
	2: "一码",
	3: "二码",
}

var DateLimit = map[int]string{
	//0: "不限",
	1: "早于今天",
	2: "不早于今天",
	3: "晚于今天",
	4: "不晚于今天",
}

var IssueMap = map[string]string{
	"01": "未填写内容",
	"02": "无法辨识填写内容",
	"03": "填写多项内容",
	"04": "填写不完整",
	"05": "填写错误",
	"06": "无法进行转换",
}
