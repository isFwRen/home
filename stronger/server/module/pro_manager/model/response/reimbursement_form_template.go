package response

type RFT struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}
