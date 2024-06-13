package request

type QuaAnaRes struct {
	Types     string `json:"types" form:"types"`
	StartTime string `json:"startTime" form:"startTime"`
	EndTime   string `json:"endTime" form:"endTime"`
}
