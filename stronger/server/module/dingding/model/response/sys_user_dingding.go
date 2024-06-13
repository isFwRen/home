package response

import "encoding/json"

type SendDingMessageResponse struct {
	RequestId string    `json:"request_id"`
	ErrMsg    string `json:"errmsg"`
	ErrCode   json.Number `json:"errcode"`
	TaskId    json.Number   `json:"task_id"`
}

type DingUserIDResponse struct {
	ErrCode string    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Result  struct {
		UserId string `json:"userid"`
	} `json:"result"`
	RequestId string `json:"request_id"`
}