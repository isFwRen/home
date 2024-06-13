package request

import modelBase "server/module/sys_base/model"

type GetTransactionRule struct {
	PageInfo modelBase.BasePageInfo `json:"page_info"`
	ProCode  string                 `json:"proCode" form:"proCode"`
	RuleType string                 `json:"ruleType" form:"ruleType"`
}

type AddTransactionRule struct {
	Id          string `json:"id" form:"id"`
	ProCode     string `json:"proCode" form:"proCode"`
	RuleName    string `json:"ruleName" form:"ruleName"`
	RuleType    string `json:"ruleType" form:"ruleType"`
	DocsPath    string `json:"docsPath" form:"docsPath"`
	UpdatedName string `json:"updatedName" form:"updatedName"`
	IsRequired  int    `json:"isRequired" form:"isRequired"`
}

type Rm struct {
	ProCode string   `json:"proCode" form:"proCode"`
	Ids     []string `json:"ids" form:"ids"`
}
