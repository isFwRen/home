package task

import (
	"fmt"
	"server/global"
	"server/module/load/service"
	"server/module/pro_conf/model"
	"strings"
)

type ProConfData struct {
	MB001  []ApiTempData       `json:"mb001"`
	Fields []model.SysProField `json:"fields"`
}

type CodeWhiteList struct {
	Temp string   `json:"temp"`
	Code []string `json:"code"`
}

type ApiTempData struct {
	Code   string                 `json:"code"`
	Name   string                 `json:"name"`
	Fields []model.TempBFRelation `json:"fields"`
}

var CacheProConf ProConfData = ProConfData{}
var CacheCodeWhiteList = make(map[string][]CodeWhiteList)

// var CacheFieldConf = make(map[string]model.SysProField)

func Init() {
	// cacheProConf := make(map[string]interface{})
	// CacheFieldConf = make(map[string]model.SysProField)
	CacheCodeWhiteList = make(map[string][]CodeWhiteList)
	proCode := strings.ReplaceAll(TaskProCode, "_task", "")
	proID := global.ProCodeId[proCode]
	fmt.Println("proIDproID", proID)
	err, temp := service.GetSysProTempByProIdAndName(proID, "MB001")
	fmt.Println("temptemp", err, temp)
	err, blocks := service.GetSysProTempBlockByTempId(temp.ID)
	tBlocks := []ApiTempData{}
	for _, block := range blocks {
		_, fields := service.GetTempBFRelationByBId(block.ID)
		tBlock := ApiTempData{
			Code:   block.Code,
			Name:   block.Name,
			Fields: fields,
		}
		tBlocks = append(tBlocks, tBlock)
	}
	err, fields := service.GetSysFields(proID)
	err, sysProject := service.GetSysProject(proCode)
	NextTime = float64(sysProject.CacheTime)
	CacheProConf = ProConfData{
		Fields: fields,
		MB001:  tBlocks,
	}
	err, whiteLists := service.GetSysWhiteLists(proCode)
	for _, whiteList := range whiteLists {
		code := whiteList.UserCode
		_, isOK := CacheCodeWhiteList[code]
		if !isOK {
			CacheCodeWhiteList[code] = []CodeWhiteList{}
		}
		value := CodeWhiteList{
			Temp: whiteList.TempName,
			Code: whiteList.BlockPermissions,
		}
		CacheCodeWhiteList[code] = append(CacheCodeWhiteList[code], value)
	}

}
