package response

type PageResult struct {
	List      interface{} `json:"list"`
	Total     int64       `json:"total"`
	Msg       string      `json:"msg"`
	PageIndex int         `json:"pageIndex"`
	PageSize  int         `json:"pageSize"`
}
