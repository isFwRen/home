package response

type BusinessRes struct {
	List   interface{} `json:"list"`
	Total  int64       `json:"total"`
	ErrMsg string      `json:"errMsg"`
}
