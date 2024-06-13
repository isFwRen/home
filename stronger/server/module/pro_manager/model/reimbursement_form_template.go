package model

import modelbase "server/module/sys_base/model"

type PicInformation struct {
	modelbase.Model `json:"model"`
	ProCode         string  `json:"proCode"`
	Name            string  `json:"name"`
	EditTime        string  `json:"editTime"`
	EditName        string  `json:"editName"`
	Types           string  `json:"types"`
	Size            float64 `json:"size"`
	Path            string  `json:"path"`
	IsRequired      int     `json:"isRequired"`
}
