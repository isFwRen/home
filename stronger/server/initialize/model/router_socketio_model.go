package model

type ConstRequest struct {
	ProCode      string            `json:"proCode"`
	Relationship map[string]string `json:"relationship"`
}

type ConstReplyVersionTwo struct {
	Msg ConstMsgVersionTwo `json:"list"`
}

type ConstReply struct {
	Msg ConstMsg `json:"list"`
}

type ConstMsgVersionTwo struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type ConstMsg struct {
	ReadExcel        string `json:"ReadExcelError"`
	SaveToCouchdb    string `json:"SaveToCouchdbError"`
	SuccessOrFailure string `json:"SuccessOrFailure"`
}

type SalaryReply struct {
	Msg SalaryMsg `json:"list"`
}

type SalaryMsg struct {
	ErrMsg           string `json:"ErrMsg"`
	SuccessOrFailure string `json:"SuccessOrFailure"`
}
