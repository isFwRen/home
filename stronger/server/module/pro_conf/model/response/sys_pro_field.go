/**
 * @Author: 星期一
 * @Description:
 * @Date: 2020/11/3 4:20 下午
 */

package response

import (
	"server/module/sys_base/model"
)

type SysProField struct {
	model.Model
	Name string `json:"name" gorm:"comment:'字段名字'"`
	Code string `json:"code" gorm:"comment:'字段编码'"`
}
