package response

type ProAnalysis struct {
	ProCode    string `json:"proCode" excel:"项目编码"`
	Nums       int64  `json:"nums" excel:"差错票数"`
	Percentage string `json:"percentage" excel:"差错占比"`
}

type FieldAnalysis struct {
	FiledName  string `json:"filedName" excel:"错误字段"`
	Nums       int64  `json:"nums" excel:"差错票数"`
	Percentage string `json:"percentage" excel:"差错占比"`
}

type PeopleAnalysis struct {
	People     string `json:"people" excel:"责任人"`
	Nums       int64  `json:"nums" excel:"差错票数"`
	Percentage string `json:"percentage" excel:"差错占比"`
}
