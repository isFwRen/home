package response

type BasePageResult struct {
	List      interface{} `json:"list"`
	Total     int64       `json:"total"`
	PageIndex int         `json:"pageIndex"`
	PageSize  int         `json:"pageSize"`
	MaxOrder   int64      `json:"maxOrder"`
}

type PageResult struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	//Total2   int64       `json:"total2"`
	PageIndex int `json:"pageIndex"`
	PageSize  int `json:"pageSize"`
}

type HasMaxPageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	MaxCode  string      `json:"maxCode"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

type RowResult struct {
	Row int64 `json:"row"`
}

type Authority struct {
	List interface{} `json:"list"`
}
