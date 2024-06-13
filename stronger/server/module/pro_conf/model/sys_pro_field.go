/**
 * @Author: 星期一
 * @Description:
 * @Date: 2020/11/3 4:20 下午
 */

package model

import (
	"server/module/sys_base/model"

	"github.com/lib/pq"
)

type SysProField struct {
	model.Model
	Name              string             `json:"name" form:"name" gorm:"comment:'字段名字'"`
	Code              string             `json:"code" form:"code" gorm:"comment:'字段编码'"`
	MyOrder           int                `json:"myOrder" form:"myOrder" gorm:"comment:'排序'"`
	FixValue          string             `json:"fixValue" form:"fixValue" gorm:"comment:'固定值'"`
	SpecChar          string             `json:"specChar" form:"specChar" gorm:"comment:'可通过字符'"`
	DefaultVal        string             `json:"defaultVal" form:"defaultVal" gorm:"comment:'默认值'"`
	CheckDate         int                `json:"checkDate" form:"checkDate" gorm:"comment:'日期期限'"`
	ValChange         string             `json:"valChange" form:"valChange" gorm:"comment:'数值转换'"`
	QuestionChange    string             `json:"questionChange" form:"questionChange" gorm:"comment:'问题件转换'"`
	ValInsert         string             `json:"valInsert" form:"valInsert" gorm:"comment:'数据插入'"`
	IgnoreIf          string             `json:"ignoreIf" form:"ignoreIf" gorm:"comment:'不录条件'"`
	Prompt            string             `json:"prompt" form:"prompt" gorm:"comment:'录入提示'"`
	MaxLen            string             `json:"maxLen" form:"maxLen" gorm:"comment:'最大长度'"`
	MinLen            string             `json:"minLen" form:"minLen" gorm:"comment:'最小长度'"`
	FixLen            string             `json:"fixLen" form:"fixLen" gorm:"comment:'固定长度'"`
	MaxVal            string             `json:"maxVal" form:"maxVal" gorm:"comment:'最大值'"`
	MinVal            string             `json:"minVal" form:"minVal" gorm:"comment:'最小值'"`
	Validations       pq.Int32Array      `json:"validations" form:"validations" gorm:"type:int4[] comment:'一系列校验的数组'"`
	ProName           string             `json:"proName" form:"proName" gorm:"comment:'项目名字'"`
	ProId             string             `json:"proId" form:"proId" gorm:"comment:'项目id'"`
	InputProcess      int                `json:"inputProcess" form:"inputProcess" gorm:"comment:'录入工序'"`
	SysIssues         []SysIssue         `json:"sysIssues" form:"sysIssues" gorm:"foreignKey:FId;references:id;comment:'问题件配置'"`
	SysProFieldChecks []SysProFieldCheck `json:"sysProFieldChecks" form:"sysProFieldChecks" gorm:"foreignKey:FId;references:id;comment:'复核配置'"`
}
