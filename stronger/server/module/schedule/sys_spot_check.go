package schedule

import (
	"fmt"
	proManagerModel "server/module/pro_manager/model"
	proManagerService "server/module/pro_manager/service"
	"time"

	"server/module/schedule/service"
)

func InitSysSpotCheckSum() {
	fmt.Println("--------------------InitSysSpotCheckSum-----------------------")
	err := UserCheck()
	err = BlockCheck()
	fmt.Println("-----------------InitSysSpotCheckSum---err-----------------------", err)
}

func UserCheck() (err error) {
	err, lists := proManagerService.GetSysSpotCheckByType(1)
	if err != nil {
		return
	}
	for _, list := range lists {
		pro_code := list.ProCode
		ratio := list.Ratio
		code := list.Code
		// fmt.Println("--------------------list-----------------------", pro_code, code, ratio)
		if pro_code != "" && code != "" {
			_, ids := service.SelectBlockByUser(pro_code, code, ratio)
			data := proManagerModel.SysSpotCheckData{}
			data.Code = code
			data.Name = list.Name
			data.BlockId = ids
			data.ProCode = pro_code
			data.Num = len(ids)
			data.DoneNum = 0
			data.UndoneNum = len(ids)
			data.WrongNum = 0
			data.Type = 1
			t := time.Now()
			submitDay, _ := time.Parse("2006-01-02 00:00:00", t.AddDate(0, 0, -1).Format("2006-01-02 00:00:00"))
			data.SubmitDay = submitDay
			// fmt.Println("--------------------data-----------------------", data)
			err = proManagerService.InsertSysSpotCheckData(data)
			fmt.Println("-------------UserCheck-------err-----------------------", err)

		}

	}
	return err
}

func BlockCheck() (err error) {
	err, lists := proManagerService.GetSysSpotCheckByType(2)
	if err != nil {
		return
	}
	for _, list := range lists {
		pro_code := list.ProCode
		ratio := list.Ratio
		code := list.Code
		// fmt.Println("--------------------list-----------------------", pro_code, code, ratio)
		if pro_code != "" && code != "" {
			_, ids := service.SelectBlockByBlock(pro_code, code, ratio)
			data := proManagerModel.SysSpotCheckData{}
			data.Code = code
			data.Name = list.Name
			data.BlockId = ids
			data.ProCode = pro_code
			data.Num = len(ids)
			data.DoneNum = 0
			data.UndoneNum = len(ids)
			data.WrongNum = 0
			data.Type = 2
			t := time.Now()
			submitDay, _ := time.Parse("2006-01-02 00:00:00", t.AddDate(0, 0, -1).Format("2006-01-02 00:00:00"))
			data.SubmitDay = submitDay
			// fmt.Println("--------------------data-----------------------", data)
			err = proManagerService.InsertSysSpotCheckData(data)
			fmt.Println("------------BlockCheck--------err-----------------------", err)

		}

	}
	return err
}

// func StatisticsCheck() (err error) {
// 	err, lists := proManagerService.GetSysSpotCheckByType(2)
// 	if err != nil {
// 		return
// 	}
// 	for _, list := range lists {
// 		pro_code := list.ProCode
// 		ratio := list.Ratio
// 		code := list.Code
// 		// fmt.Println("--------------------list-----------------------", pro_code, code, ratio)
// 		if pro_code != "" && code != "" {
// 			_, num := service.CountProBlock(pro_code)
// 			data := proManagerModel.SysSpotCheckStatistic{}
// 			data.ProCode = pro_code
// 			data.Sum = num

// 			// fmt.Println("--------------------data-----------------------", data)
// 			err = proManagerService.InsertSysSpotCheckData(data)
// 			fmt.Println("------------BlockCheck--------err-----------------------", err)

// 		}

// 	}
// 	return err
// }
