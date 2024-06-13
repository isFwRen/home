package model

// Paging common input parameter structure
type PageInfo struct {
	PageIndex int    `json:"pageIndex" form:"pageIndex"`
	PageSize  int    `json:"pageSize" form:"pageSize"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Team      string `json:"team"`
	State     string `json:"state"`
}

// Find by id structure
type GetById struct {
	Id string `json:"id" form:"id"`
}

//type IdsReq struct {
//	Ids []string `json:"ids" form:"ids"`
//}
//
//type IdsIntReq struct {
//	Ids []string `json:"ids" form:"ids"`
//}

type ProId struct {
	ProId string `json:"proId" form:"proId"`
}

type PageInfoSearch struct {
	FieldLike string `json:"fieldLike" form:"fieldLike"`
	BasePageInfo
}

type InfoSearch struct {
	BlockName string `json:"blockName" form:"blockName"`
	TempId    string `json:"tempId" form:"tempId" binding:"required"`
}
