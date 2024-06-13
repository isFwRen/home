package response

import (
	"server/module/pro_conf/model/request"
)

type TableTop struct {
	Id       string        `json:"_id"`
	Rev      string        `json:"_rev"`
	Tabletop []interface{} `json:"tabletop"`
}

type TableTopResp struct {
	Id  string   `json:"_id"`
	Rev string   `json:"_rev"`
	Top []string `json:"top"`
}

type ConstTableNameByProjectResp struct {
	Id          string `json:"_id"`
	Rev         string `json:"_rev"`
	ProCode     string `json:"proCode"`
	FileName    string `json:"filename"`
	FilePath    string `json:"filePath"`
	Conf        string `json:"conf" example:"类型, 默认为const"`
	ChineseName string `json:"chineseName" example:"常量表的中文名字"`
	DbName      string `json:"dbName"`
}

type ConstTableNameByProject struct {
	Id          string `json:"_id"`
	Rev         string `json:"_rev"`
	ProCode     string `json:"proCode"`
	FileName    string `json:"filename"`
	FilePath    string `json:"filePath"`
	Conf        string `json:"conf"`
	ChineseName string `json:"chineseName"`
	DbName      string `json:"dbName"`
}

type ConstTablesTop struct {
	ConstTables []TableTopResp `json:"list"`
}

type ConstTablesWithName struct {
	ConstTables []ConstTableNameByProjectResp `json:"list"`
	Count       int                           `json:"count"`
}

type ConstTableResp struct {
	ConstTable request.UpdateConstTableStructWithItems `json:"list"`
}

type ConstResp struct {
	Ok     string
	Err    string
	Reason string
}

type UploadFile struct {
	RelationShip map[string]string `json:"dbname"`
	Err          []string          `json:"err"`
}
