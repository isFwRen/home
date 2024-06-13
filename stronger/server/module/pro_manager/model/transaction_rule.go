package model

import modelbase "server/module/sys_base/model"

type TransactionRule struct {
	modelbase.Model `json:"model"`
	ProCode         string `json:"proCode"`
	RuleName        string `json:"ruleName"`
	RuleType        string `json:"ruleType"`
	DocsPath        string `json:"docsPath"`
	UpdatedName     string `json:"updatedName"`
	IsRequired      int    `json:"isRequired"`
}
