package project

import (
	"server/module/load/model"
	"server/module/task/project/B0103"
	"server/module/task/project/B0106"
	b0108 "server/module/task/project/B0108"
	"server/module/task/project/B0110"
	"server/module/task/project/B0113"
	b0114 "server/module/task/project/B0114"
	"server/module/task/project/B0116"
	b0118 "server/module/task/project/B0118"
	"server/module/task/project/B0121"
	"server/module/task/project/B0122"
)

func BillValue(pro string, block model.ProjectBlock, fields []model.ProjectField, op string) map[string]interface{} {
	// op := block.Stage
	switch pro {
	case "B0118":
		return b0118.BillValue("B0118_task", block, fields, op)
	case "B0108":
		return b0108.BillValue("B0108_task", block, fields, op)
	case "B0114":
		return b0114.BillValue("B0114_task", block, fields, op)
	case "B0113":
		return B0113.BillValue("B0113_task", block, fields, op)
	case "B0121":
		return B0121.BillValue("B0121_task", block, fields, op)
	case "B0110":
		return B0110.BillValue("B0110_task", block, fields, op)
	case "B0122":
		return B0122.BillValue("B0122_task", block, fields, op)
	case "B0116":
		return B0116.BillValue("B0116_task", block, fields, op)
	case "B0106":
		return B0106.BillValue("B0106_task", block, fields, op)
	case "B0103":
		return B0103.BillValue("B0103_task", block, fields, op)
	case "MB0002":
		return make(map[string]interface{})
	default:
		return make(map[string]interface{})
	}
}
