/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/11 10:45
 */

package model

import (
	modelBase "server/module/sys_base/model"
	"time"
)

type YieldTarget struct {
	modelBase.Model
	Target     float64   `json:"target"`
	UserName   string    `json:"userName"`
	UserCode   string    `json:"userCode"`
	TargetDate time.Time `json:"releaseDate"`
}

type Yield struct {
	modelBase.Model
	Value     float64   `json:"value"`
	UserName  string    `json:"userName"`
	UserCode  string    `json:"userCode"`
	YieldDate time.Time `json:"yieldDate"`
}
