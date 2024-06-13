package response

type FRReq struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}
