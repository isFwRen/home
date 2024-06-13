/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/5/23 17:39
 */

package model

type TableInfo struct {
	Id       string   `json:"_id"`
	Rev      string   `json:"_rev"`
	Content  string   `json:"content"`
	TableTop []string `json:"tabletop"`
	Arr      []string `json:"arr"`
}
