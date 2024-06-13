package model

type ProPermissionSum struct {
	Procode string `json:"proCode"` //项目编码
	Proname string `json:"proName"` //项目名称
	Op0     int64  `json:"op0"`     //初审人数
	Op1     int64  `json:"op1"`     //一码人数
	Op2     int64  `json:"op2"`     //二码人数
	Opq     int64  `json:"opq"`     //问题件人数
	Innet   int64  `json:"innet"`   //内网人数
	Outnet  int64  `json:"outnet"`  //外网人数
}
