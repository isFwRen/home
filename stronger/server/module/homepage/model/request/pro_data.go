/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/8/3 14:54
 */

package request

type QueryDayReq struct {
	QueryDay string `json:"queryDay" form:"queryDay"`
}

type QueryBusinessReq struct {
	RankingType int `json:"rankingType" form:"rankingType"` //0:日排行，1：月排行，2：年排行
}
