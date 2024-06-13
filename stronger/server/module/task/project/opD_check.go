package project

import (
	"server/global"
	"server/module/load/model"
)

func OpDCheck(block model.ProjectBlock, fields []model.ProjectField) model.ProjectBlock {
	// op := ""
	switch global.GConfig.System.ProCode {
	case "B0118":
		return block
	case "B0108":
		return block
	default:
		return block
	}
}
