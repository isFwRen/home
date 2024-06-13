package model

import (
	modelBase "server/module/sys_base/model"
	"time"
)

type DingtalkGroup struct {
	modelBase.Model
	Name        string    `json:"name"`
	ProCode     string    `json:"proCode"`
	AccessToken string    `json:"accessToken"`
	Env         int       `json:"env"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	MyOrder     int       `json:"myOrder"`
	Secret      string    `json:"secret"`
}
