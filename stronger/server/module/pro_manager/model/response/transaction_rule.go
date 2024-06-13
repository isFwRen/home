package response

type TransactionRuleRes struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}
