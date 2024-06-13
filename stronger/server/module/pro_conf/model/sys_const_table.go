package model

type ConstTable struct {
	Id  string        `json:"_id"`
	Rev string        `json:"_rev"`
	Arr []interface{} `json:"arr"`
}
