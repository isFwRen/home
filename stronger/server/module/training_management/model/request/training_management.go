package request

type ReqTraining struct {
	PageIndex   int    `json:"pageIndex" form:"pageIndex,default=1" binding:"required,min=1"`
	PageSize    int    `json:"pageSize" form:"pageSize,default=10" binding:"required,max=1500,min=1"`
	ProjectCode string `json:"projectCode" form:"projectCode"` //项目编码
	IsAt        string `json:"isAt" form:"isAt"`               //日期
	UserName    string `json:"userName" form:"userName"`       //姓名
	UserCode    string `json:"userCode" form:"userCode"`       //工号
	AuditStatus string `json:"auditStatus" form:"auditStatus"` //审核状态 1.待审核 2.审核通过  3.审核未通过
}
type PageResult struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	//Total2   int64       `json:"total2"`
	PageIndex int `json:"pageIndex"`
	PageSize  int `json:"pageSize"`
}

// 审核
type ReqAuditStatus struct {
	Ids         []string `json:"ids" form:"ids"`
	AuditStatus string   `json:"auditStatus" form:"auditStatus"` //审核状态 1.待审核 2.审核通过  3.审核未通过
}
