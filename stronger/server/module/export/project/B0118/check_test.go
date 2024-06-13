/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/29 1:50 下午
 */

package B0118

import (
	"fmt"
	model2 "server/module/export/model"
	"server/module/pro_manager/model"
	model3 "server/module/sys_base/model"
	"testing"
)

func TestCheckXml(t *testing.T) {
	xmlValue := "<?xml version=\"1.0\" encoding=\"gb2312\"?>\n<bb1>\n    <bbb>3333</bbb>\n    <aaa>222</aaa>\n   <bbb>33331</bbb> <ccc></ccc>\n   <bbb/> <test1>111</test1>\n    <test2>124</test2>\n    <test5>123</test5>\n    <test6>1</test6>\n</bb1>"
	var resultDataBill = model2.ResultDataBill{
		Bill: model.ProjectBill{
			Model: model3.Model{
				ID: "900102601343762432",
			},
			ProCode: "B0114",
		},
	}
	err := CheckXml(resultDataBill, xmlValue)
	if err != nil {
		fmt.Printf(err.Error())
	}
}
