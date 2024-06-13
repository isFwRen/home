package response

type QualityRes struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}
