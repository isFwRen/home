/**
 * @Author: 星期一
 * @Description:
 * @Date: 2020/11/23 下午5:07
 */

package model

import (
	"github.com/lib/pq"
	"server/module/sys_base/model"
)

type SysInspection struct {
	model.Model
	ProName     string        `json:"proName" form:"proName" gorm:"comment:项目名称"`
	ProId       string        `json:"proId" form:"proId" gorm:"comment:项目id"`
	MyOrder     int           `json:"myOrder" form:"myOrder" gorm:"comment:排序"`
	XmlNodeName string        `json:"xmlNodeName" form:"xmlNodeName" gorm:"xml节点名字"`
	XmlNodeCode string        `json:"xmlNodeCode" form:"xmlNodeCode" gorm:"xml节点代码"`
	Msg         string        `json:"msg" form:"msg" gorm:"自定义描述"`
	OnlyInput   string        `json:"onlyInput" form:"onlyInput" gorm:"只可输入"`
	NotInput    string        `json:"notInput" form:"notInput" gorm:"不可输入"`
	MaxLen      string        `json:"maxLen" form:"maxLen" gorm:"最大长度"`
	MinLen      string        `json:"minLen" form:"minLen" gorm:"最小长度"`
	MaxVal      string        `json:"maxVal" form:"maxVal" gorm:"最大值"`
	MinVal      string        `json:"minVal" form:"minVal" gorm:"最小值"`
	Validation  pq.Int32Array `json:"validation" form:"validation" gorm:"type:int4[]  comment:'校验数组'"`
}
