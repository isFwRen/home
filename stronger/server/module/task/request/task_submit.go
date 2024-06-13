package request

import (
	"server/module/load/model"
	bModel "server/module/pro_manager/model"
)

type TaskSubmit struct {
	Bill   bModel.ProjectBill     `json:"bill"`
	Block  model.ProjectBlock     `json:"block"`
	Fields [][]model.ProjectField `json:"fields"`
	Op     string                 `json:"op"`
}

type KeyLog struct {
	BillNum string   `json:"billNum"`
	Log     []string `json:"log"`
}
