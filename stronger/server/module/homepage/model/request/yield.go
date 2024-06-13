/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/11 10:42
 */

package request

import "server/module/sys_base/model"

type Target struct {
	Target float64 `json:"target" form:"target" binding:"min=0"`
}

type YieldRankingReq struct {
	model.BasePageInfo
	RankingType int `json:"rankingType" form:"rankingType"` //0:日排行，1：月排行
}
