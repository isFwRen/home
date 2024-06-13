/**
 * @Author: 星期一
 * @Description:
 * @Date: 2020/12/23 下午2:02
 */

package request

type ChangeOrder struct {
	ExportId   string `json:"exportId"`
	StartOrder int    `json:"startOrder"`
	EndOrder   int    `json:"endOrder"`
	StartId    string `json:"startId"`
}

type SwitchOrder struct {
	EndId      string `json:"endId" form:"endId"`
	StartOrder int    `json:"startOrder" form:"startOrder"`
	EndOrder   int    `json:"endOrder" form:"endOrder"`
	StartId    string `json:"startId" form:"startId"`
}
