package model

type EntryChannel struct {
	ProCode     string `json:"proCode"`
	Op0Num      int64  `json:"op0Num"`
	Op1Num      int64  `json:"op1Num"`
	Op2Num      int64  `json:"op2Num"`
	OpqNum      int64  `json:"opqNum"`
	InnerIp     string `json:"innerIp"`
	OutIp       string `json:"outIp"`
	OutAppPort  int    `json:"outAppPort"`
	InAppPort   int    `json:"inAppPort"`
	BackEndPort int    `json:"backEndPort"`
}
