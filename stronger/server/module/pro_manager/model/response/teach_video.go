package response

type TVReq struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}
