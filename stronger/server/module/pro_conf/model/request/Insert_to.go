/**
 * @Author: 星期一
 * @Description:
 * @Date: 2020/12/23 下午2:02
 */

package request

type InsertTo struct {
	ExportId   int `json:"exportId"`
	StartOrder int `json:"startOrder"`
	EndOrder   int `json:"endOrder"`
	StartId    int `json:"startId"`
}
