/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/19 3:18 下午
 */

package model

type XmlDataAndPath struct {
	Data    string `json:"data" from:"data"`       //xml数据
	Url     string `json:"url" from:"url"`         //xml路径
	ProCode string `json:"proCode" from:"proCode"` //项目编码
}
