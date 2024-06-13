package request

type GetRFT struct {
	ProCode string `json:"proCode" form:"proCode"`
	Name    string `json:"name" form:"name"`
}

type RFTRename struct {
	Id      string `json:"id" form:"id"`
	ProCode string `json:"proCode" form:"proCode"`
	Name    string `json:"name" form:"name"`
}
