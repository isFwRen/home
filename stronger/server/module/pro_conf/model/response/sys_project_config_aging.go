package response

type ProjectConfigAgingResponse struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}
