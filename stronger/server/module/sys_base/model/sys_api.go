/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/12/28 9:44 上午
 */

package model

type SysApi struct {
	Model
	Handle string `json:"handle" gorm:"size:128;comment:handle"`
	Title  string `json:"title" gorm:"size:128;comment:标题"`
	Path   string `json:"path" gorm:"size:128;comment:地址"`
	Action string `json:"action" gorm:"size:16;comment:请求类型"`
	Type   string `json:"type" gorm:"size:16;comment:接口类型"`
}

func (SysApi) TableName() string {
	return "sys_apis"
}
