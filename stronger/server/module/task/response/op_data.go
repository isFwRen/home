package response

import (
	"server/module/load/model"
	bModel "server/module/pro_manager/model"
)

type OpData struct {
	Bill       bModel.ProjectBill     `json:"bill"`
	Block      model.ProjectBlock     `json:"block"`
	Fields     [][]model.ProjectField `json:"fields"`
	Bins       BillAndBlockBin        `json:"bins"`
	CodeValues map[string]interface{} `json:"codeValues"`
	CacheTime  int                    `json:"cacheTime"`
}

type TaskOpNum struct {
	Op0        int64   `json:"op0" from:"op0"`
	Op1        int64   `json:"op1" from:"op1"`
	Op2        int64   `json:"op2" from:"op2"`
	Opq        int64   `json:"opq" from:"opq"`
	Character  int     `json:"character" from:"character"`
	Accuracy   float64 `json:"accuracy" from:"accuracy"`
	BlockSpeed int32   `json:"blockSpeed" from:"blockSpeed"`
}

type BillAndBlockBin struct {
	BillBin  [][]byte
	BlockBin []byte
}
