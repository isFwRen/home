/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/25 9:35 上午
 */

package model

type QingDan struct {
	InvoiceNum  string `json:"invoiceNum" excel:"账单号"`     // 账单号
	MedicalType string `json:"medicalType" excel:"医疗项目分类"` // 医疗项目分类
	Name        string `json:"name" excel:"项目名称"`          // 项目名称
	Type        string `json:"type" excel:"项目类型"`          // 项目类型
	Price       string `json:"price" excel:"项目金额"`         // 项目金额
	Count       string `json:"count" excel:"数量"`           // 数量
	Percent     string `json:"percent" excel:"项目比例"`       // 项目比例
	Pay         string `json:"pay" excel:"项目自付"`           // 项目自付
}
