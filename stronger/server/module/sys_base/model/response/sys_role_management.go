package response

type RoleListRes struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}
