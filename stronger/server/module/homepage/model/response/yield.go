/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/11 16:22
 */

package response

import "server/module/sys_base/model/response"

type YieldRanking struct {
	Value     float64 `json:"value" gorm:"value"`
	UserName  string  `json:"userName" gorm:"userName"`
	UserCode  string  `json:"userCode" gorm:"userCode"`
	YieldDate string  `json:"yieldDate" gorm:"yieldDate"`
	MyOrder   int     `json:"myOrder" gorm:"myOrder"`
}

type RankingResult struct {
	response.BasePageResult
	UserYieldRanking YieldRanking `json:"userYieldRanking"`
}
