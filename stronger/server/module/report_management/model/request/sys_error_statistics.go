package request

type Complain struct {
	Id              string `json:"id" example:"id"`
	ComplainConfirm bool   `json:"complainConfirm"` //是否申诉
	ProCode         string `json:"proCode" example:"项目编码"`
}

type ComplainTask struct {
	List            []WrongConfirms `json:"list" form:"list"`
	ComplainConfirm bool            `json:"complainConfirm"` //是否申诉
	ProCode         string          `json:"proCode" example:"项目编码"`
}

type WrongConfirmArray struct {
	List         []WrongConfirms `json:"list" form:"list"`
	ProCode      string          `json:"proCode" form:"proCode"`
	Right        string          `json:"right" form:"right"`
	WrongConfirm string          `json:"wrongConfirm" form:"wrongConfirm"` //是否通过
}

type WrongConfirms struct {
	Id string `json:"id" form:"id"`
}

type IncorrectAnalysisExport struct {
	ProCode   string `json:"proCode" form:"proCode"`
	StartTime string `json:"startTime" form:"startTime"`
	EndTime   string `json:"endTime" form:"endTime"`
}
