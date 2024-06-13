package request

type FinishParam struct {
	FinishId    string `json:"finishId"`
	ProjectCode string `json:"projectCode"`
	TrainType   int    `json:"trainType"`
}
