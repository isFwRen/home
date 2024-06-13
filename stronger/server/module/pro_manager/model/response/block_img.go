/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/2/25 4:10 下午
 */

package response

import (
	"github.com/lib/pq"
	"server/module/sys_base/model"
)

type BlockImg struct {
	model.Model
	BillID string `json:"billID" gorm:"comment:'单ID'"` //单ID
	Name   string `json:"name" gorm:"模板分块名字"`          //模板分块名字
	Code   string `json:"code" gorm:"模板分块编码"`
	//PicPage   int            `json:"picPage" gorm:"图片页码，第几张图片"`                   //图片页码，第几张图片
	//MPicPage  int            `json:"mPicPage" gorm:"移动端图片页码，第几张图片"`               //移动端图片页码，第几张图片
	//PreBCode  pq.StringArray `json:"preBCode" gorm:"type:varchar(100)[];前置分块编码"`  //前置分块编码
	LinkBCode pq.StringArray `json:"linkBCode" gorm:"type:varchar(100)[];参考分块编码"` //参考分块编码
	Status    int            `json:"status" gorm:"comment:'1(初审)|2'"`             //1(初审)|2
	Picture   string         `json:"picture" gorm:"图片"`
}
