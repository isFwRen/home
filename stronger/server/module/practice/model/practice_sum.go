package model

import (
	"server/module/sys_base/model"
	"time"
)

// "github.com/lib/pq"
type PracticeSum struct {
	model.Model
	Name                           string    `json:"name" from:"name"`               //用户名
	Code                           string    `json:"code" from:"code"`               //工号
	ProCode                        string    `json:"proCode" from:"proCode"`         //项目编号
	ApplyAt                        time.Time `json:"applyAt" gorm:"comment:'领取时间'"`  //1码领取时间
	SubmitAt                       time.Time `json:"submitAt" gorm:"comment:'提交时间'"` //1码提交时间
	SummaryFieldNum                int       `json:"summaryFieldNum" gorm:"comment:'汇总-字段数量'" excel:"汇总-字段数量"`
	SummaryCostTime                int64     `json:"summaryCostTime" gorm:"comment:'汇总-花费时间'" excel:"汇总-花费时间"`
	CostTime                       string    `json:"costTime" from:"-" gorm:"-"`
	SummaryFieldCharacter          int       `json:"summaryFieldCharacter" gorm:"comment:'总字符总量'" excel:"汇总-总字符总量"`
	SummaryFieldEffectiveCharacter int       `json:"summaryFieldEffectiveCharacter" gorm:"comment:'总字符有效总量'" excel:"汇总-总字符有效总量"`
	SummaryAccuracyRate            float64   `json:"summaryAccuracyRate" gorm:"comment:'总准确率'" excel:"汇总-总准确率"`
	SummaryBlockNum                int       `json:"summaryBlockNum" gorm:"comment:'总分块数量'" excel:"汇总-总分块数量"`
	SummaryBlockEfficiency         float64   `json:"summaryBlockEfficiency" gorm:"comment:'总分块效率'" excel:"汇总-总分块效率"`
	SummaryFieldEfficiency         float64   `json:"summaryFieldEfficiency" gorm:"comment:'总字符效率'" excel:"汇总-总字符效率"`
	SummaryQuestionMarkNumber      int       `json:"summaryQuestionMarkNumber" gorm:"comment:'总录入?数量'" excel:"汇总-总录入?数量"`
	SummaryQuestionMarkProportion  float64   `json:"summaryQuestionMarkProportion" gorm:"comment:'总录入?比例'" excel:"汇总-总录入?比例"`
}
