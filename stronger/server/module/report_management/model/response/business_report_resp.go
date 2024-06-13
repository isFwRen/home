/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/28 15:27
 */

package response

type BusinessReportResp struct {
	CountDate string `json:"countDate"`
	Count     int64  `json:"count"`
}

type AgingReportResp struct {
	CountDate string  `json:"countDate"`
	Count     float64 `json:"count"`
}
type DealTimeReportResp struct {
	CountDate string `json:"countDate"`
	Count     string `json:"count"`
}
