package model

type ProjectAgingManagement struct {
	ProCode         string  `json:"proCode"`
	CaseNumber      string  `json:"caseNumber"`      //单号
	Agency          string  `json:"agency"`          //机构号
	ScanAt          string  `json:"scanAt"`          //扫描时间
	CreatedAt       string  `json:"createdAt"`       //案件时间
	BackAtTheLatest string  `json:"backAtTheLatest"` //最晚回传时间
	TimeRemaining   string  `json:"timeRemaining"`   //剩余时间
	Second          float64 `json:"second"`          //剩余时间(秒), 给前端排序
	Stage           string  `json:"status"`          //录入状态
	Status          string  `json:"caseStatus"`      //案件状态
	ClaimType       int     `json:"claimType"`       //理赔类型
}
