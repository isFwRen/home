package request

type ReqSettingReport struct {
	ProjectCode string   `json:"projectCode"`
	TagsList    []string `json:"tagsList"`
}

type ReqGetUserTagList struct {
	ProjectCode string `json:"projectCode"`
}
