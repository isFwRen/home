package response

type CC struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}
