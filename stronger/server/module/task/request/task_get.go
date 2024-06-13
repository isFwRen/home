package request

type TaskGet struct {
	Code string `json:"code" form:"code" gorm:"工号"`
	Op   string `json:"op" form:"op" gorm:"工序"`
	Num  int    `json:"num" form:"op" gorm:"返回数字"`
}

type BillRelease struct {
	ID string `json:"id" form:"id" gorm:"单号"`
}

type ImageData struct {
	Path string `json:"path" form:"path" gorm:"路径"`
}

type TaskRelease struct {
	Code string `json:"code" form:"code" gorm:"工号"`
	Op   string `json:"op" form:"op" gorm:"工序"`
	ID   string `json:"id" form:"id" gorm:"ID"`
}

type TaskImageThumbnail struct {
	PageIndex int    `json:"pageIndex" form:"pageIndex"`
	BillId    string `json:"billId" form:"billId"`
	ProCode   string `json:"proCode" form:"proCode"`
}

type TaskImage struct {
	Index   int    `json:"index" form:"index"`
	BillId  string `json:"billId" form:"billId"`
	ProCode string `json:"proCode" form:"proCode"`
	BlockId string `json:"blockId" form:"blockId"`
}
