package response

type ProPermissionSum struct {
	List  interface{} `json:"list"`
	Total int         `json:"total"`
}
