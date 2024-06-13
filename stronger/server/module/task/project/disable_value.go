package project

import (
	"server/global"
	"server/module/load/model"

	b0108 "server/module/task/project/B0108"
	"server/module/task/project/B0113"
	b0114 "server/module/task/project/B0114"
	b0118 "server/module/task/project/B0118"
	"server/module/task/project/B0121"
	"server/module/task/project/B0122"
)

func DisableValue(block model.ProjectBlock, fields []model.ProjectField, op string) []model.ProjectField {
	// op := block.Stage
	switch global.GConfig.System.ProCode {
	case "B0118":
		return b0118.DisableValue(block, fields, op)
	case "B0108":
		return b0108.DisableValue(block, fields, op)
	case "B0114":
		return b0114.DisableValue(block, fields, op)
	case "B0113":
		return B0113.DisableValue(block, fields, op)
	case "B0121":
		return B0121.DisableValue(block, fields, op)
	case "B0122":
		return B0122.DisableValue(block, fields, op)
	default:
		return fields
	}
}
