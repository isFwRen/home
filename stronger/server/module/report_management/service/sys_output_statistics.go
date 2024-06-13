package service

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"server/global"
	"server/module/report_management/model"
	"server/module/report_management/model/request"
	"server/module/report_management/model/response"
	modelBase "server/module/sys_base/model"
	u "server/utils"
	"strconv"
	"strings"
	"time"
	model2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/model"

	"github.com/shopspring/decimal"
)

func CreateTable() (err error) {
	db := global.ProDbMap["B0114"]
	if db == nil {
		return global.ProDbErr
	}
	err = db.Migrator().CreateTable(&model.OutputStatistics{})
	//err = global.GDb.Migrator().CreateTable(&model.SysCorrected{})
	//err = db.Migrator().CreateTable(&model.OutputStatisticsSummary{}, &model.Op1{}, &model.Op2{}, &model.Op0{}, &model.OpQ{}, &model.OutputStatistics{})
	return err
}

func GetOutputStatisticsTask(info request.OutPutStatisticsSearch, uid string) (err error, list interface{}, total int64, Top []string) {
	var U model2.SysUser
	err = global.GUserDb.Model(&model2.SysUser{}).Where("id = ? ", uid).Find(&U).Error
	if err != nil {
		return err, nil, 0, nil
	}

	ProCode := info.ProCode
	db := global.ProDbMap[ProCode]
	if db == nil {
		return global.ProDbErr, nil, 0, nil
	}

	info.Code = U.Code
	return GetOutPutStatisticsDetail(info)
}

func GetOcrStatistics(info request.GetOCROutPutStatisticsSearch) (err error, list []model.OcrStatistics, total int64) {
	ProCode := info.ProCode
	db := global.ProDbMap[ProCode]
	if db == nil {
		return global.ProDbErr, nil, 0
	}
	//直接查！
	db = db.Model(&model.OcrStatistics{}).Where("created_at >= ? AND created_at <= ? ", info.StartTime, info.EndTime)
	if info.FieldName != "" {
		db = db.Where("field_name = ? ", info.FieldName)

	}
	if info.BillNum != "" {
		db = db.Where("bill_num = ? ", info.BillNum)

	}
	if info.Disable != "" {
		db = db.Where("disable = ? ", info.Disable)

	}
	if info.Compare != "" {
		db = db.Where("compare = ? ", info.Compare)

	}
	if info.ResultValue != "" {
		db = db.Where("result_value like ? ", "%"+info.ResultValue+"%")

	}
	err = db.Count(&total).Error
	if err != nil {
		return err, nil, 0
	}
	// fmt.Println("---------------PageSize-------------------------", info.PageInfo.PageSize)
	if info.PageInfo.PageSize > 0 {
		limit := info.PageInfo.PageSize
		offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
		err = db.Order("created_at desc").Limit(limit).Offset(offset).Find(&list).Error
	} else {
		err = db.Order("created_at desc").Find(&list).Error
	}

	if err != nil {
		return err, nil, 0
	}
	return err, list, total
}

func GetOutputStatistics(info request.OutPutStatisticsSearch) (err error, list interface{}, total int64, Top []string) {

	if info.IsCheckAll == 1 {
		if info.Code != "" {
			return GetSomeBodyOutPutStatistics(info)
		} else {
			return GetOutPutStatisticsAll(info)
		}

	} else if info.IsCheckAll == 2 {
		return GetOutPutStatisticsDetail(info)
	}
	return errors.New("参数有误"), nil, 0, nil
}

func GetSomeBodyOutPutStatistics(info request.OutPutStatisticsSearch) (err error, list interface{}, total int64, Top []string) {
	reg1 := regexp.MustCompile("^([\u4e00-\u9fa5]+)$")
	StartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.StartTime, time.Local)
	EndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.EndTime, time.Local)
	var SummaryReturn []model.OutputStatisticsSummaryList
	var FirstSummaryList []model.OutputStatisticsSummary
	top := make([]string, 0)
	codes := make([]string, 0)
	if reg1.MatchString(info.Code) {
		fmt.Println("------------------info.Code12=", info.Code)

		err = global.GDb.Model(&model.OutputStatisticsSummary{}).Order("id desc").Where("submit_time >= ? AND submit_time <= ? AND nick_name like ? ", StartTime, EndTime, "%"+info.Code+"%").Find(&FirstSummaryList).Error
		if err != nil {
			return err, nil, 0, nil
		}
		err = global.GDb.Model(&model.OutputStatisticsSummary{}).Order("id desc").Where("submit_time >= ? AND submit_time <= ? AND nick_name like ? ", StartTime, EndTime, "%"+info.Code+"%").Distinct("pro_code").Count(&total).Error
		fmt.Println(total)
		if err != nil {
			return err, nil, 0, nil
		}
	} else {
		fmt.Println("------------------info.Code13=", info.Code)

		err = global.GDb.Model(&model.OutputStatisticsSummary{}).Order("id desc").Where("submit_time >= ? AND submit_time <= ? AND code like ? ", StartTime, EndTime, "%"+info.Code+"%").Find(&FirstSummaryList).Error
		if err != nil {
			return err, nil, 0, nil
		}
		err = global.GDb.Model(&model.OutputStatisticsSummary{}).Order("id desc").Where("submit_time >= ? AND submit_time <= ? AND code like ? ", StartTime, EndTime, "%"+info.Code+"%").Distinct("pro_code").Count(&total).Error
		fmt.Println(total)
		if err != nil {
			return err, nil, 0, nil
		}
	}
	if total == 0 {
		return errors.New("没有该用户"), nil, 0, nil
	}

	//防止查到不同日期但同一人的数据
	for _, v := range FirstSummaryList {
		existCode := false
		existPro := false
		//处理相同项目编码
		needAppend := false
		for _, v1 := range codes {
			if v1 == v.ProCode {
				needAppend = true
				break
			}
		}
		if !needAppend {
			codes = append(codes, v.ProCode)
		}
		//处理返回结构体
		for l, r := range SummaryReturn {
			if r.Code == v.Code {
				for k, ps := range r.ProSummary {
					if ps.ProCode == v.ProCode {
						r.ProSummary[k].Op0 += float64(v.Op0)
						r.ProSummary[k].Mary += float64(v.Mary)
						r.ProSummary[k].Op1ExpenseAccount += float64(v.Op1ExpenseAccount)
						r.ProSummary[k].Op1NotExpenseAccount += float64(v.Op1NotExpenseAccount)
						r.ProSummary[k].Op2ExpenseAccount += float64(v.Op2ExpenseAccount)
						r.ProSummary[k].Op2NotExpenseAccount += float64(v.Op2NotExpenseAccount)
						r.ProSummary[k].Question += float64(v.Question)
						SummaryReturn[l].AddUpToSomething += float64(v.Mary)
						existPro = true
					}
				}
				existCode = true
				if !existPro {
					SumMary := model.Summary{
						ProCode:              v.ProCode,
						Mary:                 float64(v.Mary),
						Op0:                  float64(v.Op0),
						Op1NotExpenseAccount: float64(v.Op1NotExpenseAccount),
						Op1ExpenseAccount:    float64(v.Op1ExpenseAccount),
						Op2NotExpenseAccount: float64(v.Op2NotExpenseAccount),
						Op2ExpenseAccount:    float64(v.Op2ExpenseAccount),
						Question:             float64(v.Question),
					}
					SummaryReturn[l].AddUpToSomething += float64(v.Mary)
					SummaryReturn[l].ProSummary = append(SummaryReturn[l].ProSummary, SumMary)
				}
			}
		}

		if !existCode {
			SumMary := model.Summary{
				ProCode:              v.ProCode,
				Mary:                 float64(v.Mary),
				Op0:                  float64(v.Op0),
				Op1NotExpenseAccount: float64(v.Op1NotExpenseAccount),
				Op1ExpenseAccount:    float64(v.Op1ExpenseAccount),
				Op2NotExpenseAccount: float64(v.Op2NotExpenseAccount),
				Op2ExpenseAccount:    float64(v.Op2ExpenseAccount),
				Question:             float64(v.Question),
			}

			SummaryItem := model.OutputStatisticsSummaryList{
				Code:             v.Code,
				NickName:         v.NickName,
				SubmitTime:       v.SubmitTime.Format("2006-01-02 15:04:05"),
				AddUpToSomething: float64(v.Mary),
			}
			SummaryItem.ProSummary = append(SummaryItem.ProSummary, SumMary)

			SummaryReturn = append(SummaryReturn, SummaryItem)
		}
	}

	for _, v := range SummaryReturn[0].ProSummary {
		top = append(top, v.ProCode)
	}
	fmt.Println("GetSomeBodyOutPutStatistics Finish!")
	return err, SummaryReturn, total, top
}

func GetOutPutStatisticsDetail(info request.OutPutStatisticsSearch) (err error, list interface{}, total int64, Top []string) {
	reg1 := regexp.MustCompile("^([\u4e00-\u9fa5]+)$")
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	StartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.StartTime, time.Local)
	EndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.EndTime, time.Local)
	//查询明细
	//连接数据库
	ProCode := info.ProCode
	db := global.ProDbMap[ProCode]
	if db == nil {
		return global.ProDbErr, nil, 0, nil
	}
	//直接查！
	db = db.Model(&model.OutputStatistics{}).Where("submit_time >= ? AND submit_time <= ? ", StartTime, EndTime)
	if info.Code != "" {
		if !reg1.MatchString(info.Code) {
			db = db.Where("code like ? ", "%"+info.Code+"%")
		} else {
			db = db.Where("nick_name like ? ", "%"+info.Code+"%")
		}
	}
	err = db.Order("id desc").Count(&total).Error
	if err != nil {
		return err, nil, 0, nil
	}
	var OutputStatisticsReturn []model.OutputStatistics
	err = db.Order("id desc").Limit(limit).Offset(offset).Find(&OutputStatisticsReturn).Error
	if err != nil {
		return err, nil, 0, nil
	}
	var OutputStatisticsRes []response.OutputStatisticsRes
	var OutputStatisticsResItem response.OutputStatisticsRes
	for _, OR := range OutputStatisticsReturn {
		InitDataTy := reflect.TypeOf(OR)
		InitDataVa := reflect.ValueOf(OR)

		ResItemTy := reflect.TypeOf(OutputStatisticsResItem)
		ResItemVa := reflect.ValueOf(&OutputStatisticsResItem).Elem()

		for i := 0; i < InitDataTy.NumField(); i++ {
			InitDataTyFieldName := InitDataTy.Field(i).Name
			InitDataTyFieldType := InitDataTy.Field(i).Type
			if InitDataTyFieldName == "Model" {
				continue
			}
			for j := 0; j < ResItemTy.NumField(); j++ {
				ResItemTyFieldName := ResItemTy.Field(j).Name
				if InitDataTyFieldName == ResItemTyFieldName && strings.Index(ResItemTyFieldName, "CostTime") == -1 {
					if InitDataTyFieldType.Kind() == 2 {
						ResItemVa.FieldByName(ResItemTyFieldName).SetInt(int64(InitDataVa.FieldByName(InitDataTyFieldName).Interface().(int)))
						break
					}
					if InitDataTyFieldType.Kind() == 14 {
						if strings.Index(InitDataTyFieldName, "QuestionMarkProportion") != -1 {
							ResItemVa.FieldByName(ResItemTyFieldName).SetString(strconv.FormatFloat(InitDataVa.FieldByName(InitDataTyFieldName).Interface().(float64), 'G', 5, 64) + "%")
						} else if strings.Index(InitDataTyFieldName, "AccuracyRate") != -1 {
							if InitDataVa.FieldByName(InitDataTyFieldName).Interface().(float64) == 100 {
								ResItemVa.FieldByName(ResItemTyFieldName).SetString(decimal.NewFromFloat(InitDataVa.FieldByName(InitDataTyFieldName).Interface().(float64)).RoundCeil(2).String() + "%")
							} else if InitDataVa.FieldByName(InitDataTyFieldName).Interface().(float64) == 0 && OR.NickName == "" {
								//2022.4.1-5774王泽如提出需求：A今天现在第一次录一张单, 然后A只录一个分块就不录了, 那么当前准确率赋为100%回显给前端
								ResItemVa.FieldByName(ResItemTyFieldName).SetString("100%")
							} else {
								ResItemVa.FieldByName(ResItemTyFieldName).SetString(decimal.NewFromFloat(InitDataVa.FieldByName(InitDataTyFieldName).Interface().(float64)).Mul(decimal.NewFromFloat(100)).RoundCeil(2).String() + "%")
							}
						} else {
							ResItemVa.FieldByName(ResItemTyFieldName).SetFloat(InitDataVa.FieldByName(InitDataTyFieldName).Interface().(float64))
						}
						break
					}
					if InitDataTyFieldType.Kind() == 24 {
						//ResItemVa.SetString()
						ResItemVa.FieldByName(ResItemTyFieldName).SetString(InitDataVa.FieldByName(InitDataTyFieldName).Interface().(string))
						break
					}
					if InitDataTyFieldName == "SubmitTime" {
						OutputStatisticsResItem.SubmitTime = OR.SubmitTime.Format("2006-01-02")
						break
					}
				} else if InitDataTyFieldName == ResItemTyFieldName && strings.Index(ResItemTyFieldName, "CostTime") != -1 {
					costTime := InitDataVa.FieldByName(InitDataTyFieldName).Interface().(int64)
					_, cost := J(int(costTime))
					ResItemVa.FieldByName(ResItemTyFieldName).SetString(cost)
					break
				}
			}
		}
		var user model2.SysUser
		err = global.GUserDb.Model(&model2.SysUser{}).Where("code = ? ", OR.Code).Find(&user).Error
		if err != nil {
			return err, OutputStatisticsRes, total, nil
		}
		OutputStatisticsResItem.NickName = user.Name
		OutputStatisticsRes = append(OutputStatisticsRes, OutputStatisticsResItem)
	}
	return err, OutputStatisticsRes, total, nil
}

// GetOutPutStatisticsAll 查询产量全部2.0
func GetOutPutStatisticsAll(info request.OutPutStatisticsSearch) (err error, list interface{}, total int64, Top []string) {
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	StartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.StartTime, time.Local)
	EndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.EndTime, time.Local)
	var SummaryReturn []model.OutputStatisticsSummaryList
	top := make([]string, 0)
	//info.ProCode = "B0118"
	//db := global.ProDbMap[info.ProCode]
	//if db == nil {
	//	return global.ProDbErr, nil, 0, nil
	//}
	fmt.Println("------------------info.Code21=", info.Code)
	fmt.Println("-----------------------------明细全部~~~~~~~~~~~~~~~~~~~~~~---------------------")
	var code []string
	err = global.GDb.Model(&model.OutputStatisticsSummary{}).Limit(limit).Offset(offset).Where("submit_time >= ? AND submit_time <= ? ", StartTime, EndTime).Pluck("DISTINCT(code)", &code).Error
	fmt.Println(code)
	err = global.GDb.Model(&model.OutputStatisticsSummary{}).Where("submit_time >= ? AND submit_time <= ? ", StartTime, EndTime).Distinct("code").Count(&total).Error
	fmt.Println(total)
	if err != nil {
		return err, nil, 0, nil
	}
	for _, c := range code {
		var item []model.OutputStatisticsSummary
		err = global.GDb.Model(&model.OutputStatisticsSummary{}).Where("submit_time >= ? AND submit_time <= ? AND code = ? ", StartTime, EndTime, c).Find(&item).Error
		if err != nil {
			return err, nil, 0, nil
		}
		for _, v := range item {
			existPro := false
			existCode := false

			for l, r := range SummaryReturn {
				if r.Code == v.Code {
					for k, ps := range r.ProSummary {
						if ps.ProCode == v.ProCode {
							r.ProSummary[k].Op0 += float64(v.Op0)
							r.ProSummary[k].Op0InvoiceNum += float64(v.Op0InvoiceNum)
							r.ProSummary[k].Mary += float64(v.Mary)
							r.ProSummary[k].Op1ExpenseAccount += float64(v.Op1ExpenseAccount)
							r.ProSummary[k].Op1NotExpenseAccount += float64(v.Op1NotExpenseAccount)
							r.ProSummary[k].Op2ExpenseAccount += float64(v.Op2ExpenseAccount)
							r.ProSummary[k].Op2NotExpenseAccount += float64(v.Op2NotExpenseAccount)
							r.ProSummary[k].Question += float64(v.Question)
							SummaryReturn[l].AddUpToSomething += float64(v.Mary)
							existPro = true
						}
					}
					existCode = true
					if !existPro {
						SumMary := model.Summary{
							ProCode:              v.ProCode,
							Mary:                 float64(v.Mary),
							Op0:                  float64(v.Op0),
							Op0InvoiceNum:        float64(v.Op0InvoiceNum),
							Op1NotExpenseAccount: float64(v.Op1NotExpenseAccount),
							Op1ExpenseAccount:    float64(v.Op1ExpenseAccount),
							Op2NotExpenseAccount: float64(v.Op2NotExpenseAccount),
							Op2ExpenseAccount:    float64(v.Op2ExpenseAccount),
							Question:             float64(v.Question),
						}
						SummaryReturn[l].ProSummary = append(SummaryReturn[l].ProSummary, SumMary)
						SummaryReturn[l].AddUpToSomething += float64(v.Mary)
					}
				}
			}

			if !existCode {
				SumMary := model.Summary{
					ProCode:              v.ProCode,
					Mary:                 float64(v.Mary),
					Op0:                  float64(v.Op0),
					Op0InvoiceNum:        float64(v.Op0InvoiceNum),
					Op1NotExpenseAccount: float64(v.Op1NotExpenseAccount),
					Op1ExpenseAccount:    float64(v.Op1ExpenseAccount),
					Op2NotExpenseAccount: float64(v.Op2NotExpenseAccount),
					Op2ExpenseAccount:    float64(v.Op2ExpenseAccount),
					Question:             float64(v.Question),
				}

				SummaryItem := model.OutputStatisticsSummaryList{
					Code:             v.Code,
					NickName:         v.NickName,
					SubmitTime:       v.SubmitTime.Format("2006-01-02 15:04:05"),
					AddUpToSomething: float64(v.Mary),
				}
				SummaryItem.ProSummary = append(SummaryItem.ProSummary, SumMary)

				SummaryReturn = append(SummaryReturn, SummaryItem)
			}

		}
	}
	//处理空缺的项目
	//0.[a], [b]
	//1.[a], [a,b]
	//2.[c], [a,b]
	//3.[a], [b,c]
	//4.[a,b], [c,d]
	//5.[a,b], [b,c]
	TopMap := make([][]string, 0)
	for _, v := range SummaryReturn {
		t := make([]string, 0)
		for _, v1 := range v.ProSummary {
			t = append(t, v1.ProCode)
		}
		TopMap = append(TopMap, t)
	}
	top = append(top, SL(TopMap)...)
	for i, v := range SummaryReturn {
		n := make([]string, len(top))
		copy(n, top)
		self := make([]string, 0)
		for _, v2 := range v.ProSummary {
			self = append(self, v2.ProCode)
		}

		fill := S(n, self)
		for _, v3 := range fill {
			SumMary := model.Summary{
				ProCode:              v3,
				Mary:                 0,
				Op0:                  0,
				Op0InvoiceNum:        0,
				Op1NotExpenseAccount: 0,
				Op1ExpenseAccount:    0,
				Op2NotExpenseAccount: 0,
				Op2ExpenseAccount:    0,
				Question:             0,
			}
			SummaryReturn[i].ProSummary = append(SummaryReturn[i].ProSummary, SumMary)
		}
	}

	//按top的项目顺序返回渲染
	var NewSummaryReturn []model.OutputStatisticsSummaryList
	for _, v := range SummaryReturn {

		NewSummaryReturnItem := model.OutputStatisticsSummaryList{
			Code:             v.Code,
			NickName:         v.NickName,
			SubmitTime:       v.SubmitTime,
			AddUpToSomething: v.AddUpToSomething,
		}

		for _, v1 := range top {
			for _, v2 := range v.ProSummary {
				if v1 == v2.ProCode {
					SumMary := model.Summary{
						ProCode:              v2.ProCode,
						Mary:                 v2.Mary,
						Op0:                  v2.Op0,
						Op0InvoiceNum:        v2.Op0InvoiceNum,
						Op1NotExpenseAccount: v2.Op1NotExpenseAccount,
						Op1ExpenseAccount:    v2.Op1ExpenseAccount,
						Op2NotExpenseAccount: v2.Op2NotExpenseAccount,
						Op2ExpenseAccount:    v2.Op2ExpenseAccount,
						Question:             v2.Question,
					}
					NewSummaryReturnItem.ProSummary = append(NewSummaryReturnItem.ProSummary, SumMary)
				}
			}
		}
		NewSummaryReturn = append(NewSummaryReturn, NewSummaryReturnItem)
	}

	return err, NewSummaryReturn, total, top
}

func UpdateOutputStatistics(info request.OutPutStatisticsExport) (err error, list interface{}, total int64, Top []string) {
	StartTime, err := time.ParseInLocation("2006-01-02 15:04:05", info.StartTime, time.Local)
	fmt.Println("StartTime", err)
	EndTime, err := time.ParseInLocation("2006-01-02 15:04:05", info.EndTime, time.Local)
	fmt.Println("EndTime", err)
	UpdateTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.UpdateTime+" 00:00:00", time.Local)

	var all []model.OutputStatisticsSummaryList
	var NeedChange []model.OutputStatisticsSummaryList
	if StartTime.After(UpdateTime) && EndTime.After(UpdateTime) && StartTime.Before(EndTime) {
		infos := request.OutPutStatisticsSearch{
			ProCode:    info.ProCode,
			StartTime:  info.StartTime,
			EndTime:    info.EndTime,
			Code:       info.Code,
			IsCheckAll: 1,
			PageInfo:   info.PageInfo,
		}
		err, NeedChange, _, _ = GetExportOutputStatisticsNeedPaging(infos)
		if err != nil {
			return err, nil, 0, nil
		}
	} else if StartTime.Before(UpdateTime) && EndTime.Before(UpdateTime) && StartTime.Before(EndTime) {
		infos := request.OutPutStatisticsSearch{
			ProCode:    info.ProCode,
			StartTime:  info.StartTime,
			EndTime:    info.EndTime,
			Code:       info.Code,
			IsCheckAll: 1,
			PageInfo:   info.PageInfo,
		}
		err, all, _, _ = GetExportOutputStatisticsNeedPaging(infos)
		if err != nil {
			return err, nil, 0, nil
		}
	} else if StartTime.Format("2006-01-02") == EndTime.Format("2006-01-02") && StartTime.Format("2006-01-02") == UpdateTime.Format("2006-01-02") {
		infos := request.OutPutStatisticsSearch{
			ProCode:    info.ProCode,
			StartTime:  info.StartTime,
			EndTime:    info.EndTime,
			Code:       info.Code,
			IsCheckAll: 1,
			PageInfo:   info.PageInfo,
		}
		err, NeedChange, _, _ = GetExportOutputStatisticsNeedPaging(infos)
		if err != nil {
			return err, nil, 0, nil
		}
	} else {
		theDayBeforeTheUpdate := UpdateTime.Add(-1 * time.Minute)
		infos := request.OutPutStatisticsSearch{
			ProCode:    info.ProCode,
			StartTime:  info.StartTime,
			EndTime:    theDayBeforeTheUpdate.Format("2006-01-02 15:04:05"),
			Code:       info.Code,
			IsCheckAll: 1,
			PageInfo:   info.PageInfo,
		}
		err, all, _, _ = GetExportOutputStatisticsNeedPaging(infos)
		if err != nil {
			return err, nil, 0, nil
		}

		infoS := request.OutPutStatisticsSearch{
			ProCode:    info.ProCode,
			StartTime:  info.UpdateTime,
			EndTime:    info.EndTime,
			Code:       info.Code,
			IsCheckAll: 1,
			PageInfo:   info.PageInfo,
		}
		err, NeedChange, _, _ = GetExportOutputStatisticsNeedPaging(infoS)
		if err != nil {
			return err, nil, 0, nil
		}
	}

	//On 3*3*10*3 ✔
	//On 3*3 + 3*10*3*3
	for i, v1 := range NeedChange {
		NeedChange[i].AddUpToSomething = 0
		for j, v2 := range v1.ProSummary {
			var proCorrected model.SysCorrected
			var t int64
			err = global.GDb.Model(&model.SysCorrected{}).Where("pro_code = ? ", v2.ProCode).Find(&proCorrected).Count(&t).Error
			if err != nil {
				return err, nil, 0, nil
			}
			if t == 0 {
				continue
			}

			if proCorrected.Op0AsTheBlock != 0 && proCorrected.Op0AsTheInvoice == 0 {
				NeedChange[i].ProSummary[j].Op0 = Decimal(NeedChange[i].ProSummary[j].Op0 * proCorrected.Op0AsTheBlock)
			}
			if proCorrected.Op0AsTheInvoice != 0 && proCorrected.Op0AsTheBlock == 0 {
				//这里由于明细表发票数量不知为啥没有存进汇总表，这里先查明细的发票数量出来算(待解决)
				db := global.ProDbMap[NeedChange[i].ProSummary[j].ProCode]
				if db == nil {
					return global.ProDbErr, nil, 0, nil
				}
				var out model.OutputStatistics
				err = db.Model(&model.OutputStatistics{}).Where("code = ? AND nick_name = ? AND to_char(submit_time,'YYYY-MM-DD') = ?", NeedChange[i].Code, NeedChange[i].NickName, strings.Split(NeedChange[i].SubmitTime, " ")[0]).
					Find(&out).Error
				if err != nil {
					NeedChange[i].ProSummary[j].Op0 = 0
				}
				NeedChange[i].ProSummary[j].Op0 = Decimal(float64(out.Op0InvoiceNum) * proCorrected.Op0AsTheInvoice)
			}

			NeedChange[i].ProSummary[j].Op1ExpenseAccount = Decimal(NeedChange[i].ProSummary[j].Op1ExpenseAccount * proCorrected.Op1ExpenseAccount)
			NeedChange[i].ProSummary[j].Op1NotExpenseAccount = Decimal(NeedChange[i].ProSummary[j].Op1NotExpenseAccount * proCorrected.Op1NotExpenseAccount)
			NeedChange[i].ProSummary[j].Op2ExpenseAccount = Decimal(NeedChange[i].ProSummary[j].Op2ExpenseAccount * proCorrected.Op2ExpenseAccount)
			NeedChange[i].ProSummary[j].Op2NotExpenseAccount = Decimal(NeedChange[i].ProSummary[j].Op2NotExpenseAccount * proCorrected.Op2NotExpenseAccount)
			NeedChange[i].ProSummary[j].Question = Decimal(NeedChange[i].ProSummary[j].Question * proCorrected.Question)

			NeedChange[i].ProSummary[j].Mary = Decimal(v1.ProSummary[j].Op0 + v1.ProSummary[j].Op1ExpenseAccount + v1.ProSummary[j].Op1NotExpenseAccount + v1.ProSummary[j].Op2ExpenseAccount + v1.ProSummary[j].Op2NotExpenseAccount + v1.ProSummary[j].Question)
			NeedChange[i].AddUpToSomething = Decimal(v1.ProSummary[j].Op0 + v1.ProSummary[j].Op1ExpenseAccount + v1.ProSummary[j].Op1NotExpenseAccount + v1.ProSummary[j].Op2ExpenseAccount + v1.ProSummary[j].Op2NotExpenseAccount + v1.ProSummary[j].Question)

			if len(all) == 0 {
				continue
			}

			for k, v3 := range all {
				if v3.Code != v1.Code {
					continue
				}
				all[k].AddUpToSomething += NeedChange[i].AddUpToSomething
				for l, v4 := range v3.ProSummary {
					if v4.ProCode != v2.ProCode {
						continue
					}
					all[k].ProSummary[l].Op0 = Decimal(all[k].ProSummary[l].Op0 + NeedChange[i].ProSummary[j].Op0)
					all[k].ProSummary[l].Op1NotExpenseAccount = Decimal(all[k].ProSummary[l].Op1NotExpenseAccount + NeedChange[i].ProSummary[j].Op1NotExpenseAccount)
					all[k].ProSummary[l].Op1ExpenseAccount = Decimal(all[k].ProSummary[l].Op1ExpenseAccount + NeedChange[i].ProSummary[j].Op1ExpenseAccount)
					all[k].ProSummary[l].Op2NotExpenseAccount = Decimal(all[k].ProSummary[l].Op2NotExpenseAccount + NeedChange[i].ProSummary[j].Op2NotExpenseAccount)
					all[k].ProSummary[l].Op2ExpenseAccount = Decimal(all[k].ProSummary[l].Op2ExpenseAccount + NeedChange[i].ProSummary[j].Op2ExpenseAccount)
					all[k].ProSummary[l].Question = Decimal(all[k].ProSummary[l].Question + NeedChange[i].ProSummary[j].Question)
					all[k].ProSummary[l].Mary = Decimal(all[k].ProSummary[l].Op0 + all[k].ProSummary[l].Op1NotExpenseAccount + all[k].ProSummary[l].Op1ExpenseAccount + all[k].ProSummary[l].Op2NotExpenseAccount + all[k].ProSummary[l].Op2ExpenseAccount + all[k].ProSummary[l].Question)
				}
			}
		}
	}

	top := make([]string, 0)
	TopMap := make([][]string, 0)
	if len(all) != 0 {
		for _, v := range all {
			t := make([]string, 0)
			for _, v2 := range v.ProSummary {
				t = append(t, v2.ProCode)
			}
			TopMap = append(TopMap, t)
		}
	}
	top = append(top, SL(TopMap)...)

	TopMap2 := make([][]string, 0)
	if len(NeedChange) != 0 {
		for _, v := range NeedChange {
			t := make([]string, 0)
			for _, v2 := range v.ProSummary {
				t = append(t, v2.ProCode)
			}
			TopMap2 = append(TopMap2, t)
		}
	}
	top = append(top, SL(TopMap2)...)

	if len(all) == 0 {
		return err, NeedChange, total, top
	}
	return err, all, total, top
}

func GetExportOutputStatisticsNeedPaging(info request.OutPutStatisticsSearch) (err error, all []model.OutputStatisticsSummaryList, total int64, Top []string) {
	fmt.Println("GetExportOutputStatisticsNeedPaging")
	reg1 := regexp.MustCompile("^([\u4e00-\u9fa5]+)$")
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	StartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.StartTime, time.Local)
	EndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.EndTime, time.Local)
	var SummaryReturn []model.OutputStatisticsSummaryList
	var FirstSummaryList []model.OutputStatisticsSummary

	top := make([]string, 0)
	codes := make([]string, 0)
	if info.Code != "" {
		if reg1.MatchString(info.Code) {
			err = global.GDb.Model(&model.OutputStatisticsSummary{}).Order("id desc").Limit(limit).Offset(offset).Where("submit_time >= ? AND submit_time <= ? AND nick_name = ? ", StartTime, EndTime, info.Code).Find(&FirstSummaryList).Error
			if err != nil {
				return err, nil, 0, nil
			}
			err = global.GDb.Model(&model.OutputStatisticsSummary{}).Order("id desc").Where("submit_time >= ? AND submit_time <= ? AND nick_name = ? ", StartTime, EndTime, info.Code).Distinct("pro_code").Count(&total).Error
			fmt.Println(total)
			if err != nil {
				return err, nil, 0, nil
			}
		} else {
			err = global.GDb.Model(&model.OutputStatisticsSummary{}).Order("id desc").Limit(limit).Offset(offset).Where("submit_time >= ? AND submit_time <= ? AND code = ? ", StartTime, EndTime, info.Code).Find(&FirstSummaryList).Error
			if err != nil {
				return err, nil, 0, nil
			}
			err = global.GDb.Model(&model.OutputStatisticsSummary{}).Order("id desc").Where("submit_time >= ? AND submit_time <= ? AND code = ? ", StartTime, EndTime, info.Code).Distinct("pro_code").Count(&total).Error
			fmt.Println(total)
			if err != nil {
				return err, nil, 0, nil
			}
		}
		if total == 0 {
			return errors.New("没有该用户"), nil, 0, nil
		}

		//防止查到不同日期但同一人的数据
		for _, v := range FirstSummaryList {
			existCode := false
			existPro := false
			//处理相同项目编码
			needAppend := false
			for _, v1 := range codes {
				if v1 == v.ProCode {
					needAppend = true
					break
				}
			}
			if !needAppend {
				codes = append(codes, v.ProCode)
			}
			//处理返回结构体
			for l, r := range SummaryReturn {
				if r.Code == v.Code {
					for k, ps := range r.ProSummary {
						if ps.ProCode == v.ProCode {
							r.ProSummary[k].Op0 += float64(v.Op0)
							r.ProSummary[k].Op0InvoiceNum += float64(v.Op0InvoiceNum)
							r.ProSummary[k].Mary += float64(v.Mary)
							r.ProSummary[k].Op1ExpenseAccount += float64(v.Op1ExpenseAccount)
							r.ProSummary[k].Op1NotExpenseAccount += float64(v.Op1NotExpenseAccount)
							r.ProSummary[k].Op2ExpenseAccount += float64(v.Op2ExpenseAccount)
							r.ProSummary[k].Op2NotExpenseAccount += float64(v.Op2NotExpenseAccount)
							r.ProSummary[k].Question += float64(v.Question)
							SummaryReturn[l].AddUpToSomething += float64(v.Mary)
							existPro = true
						}
					}
					existCode = true
					if !existPro {
						SumMary := model.Summary{
							ProCode:              v.ProCode,
							Mary:                 float64(v.Mary),
							Op0:                  float64(v.Op0),
							Op0InvoiceNum:        float64(v.Op0InvoiceNum),
							Op1NotExpenseAccount: float64(v.Op1NotExpenseAccount),
							Op1ExpenseAccount:    float64(v.Op1ExpenseAccount),
							Op2NotExpenseAccount: float64(v.Op2NotExpenseAccount),
							Op2ExpenseAccount:    float64(v.Op2ExpenseAccount),
							Question:             float64(v.Question),
						}
						SummaryReturn[l].ProSummary = append(SummaryReturn[l].ProSummary, SumMary)
					}
				}
			}

			if !existCode {
				SumMary := model.Summary{
					ProCode:              v.ProCode,
					Mary:                 float64(v.Mary),
					Op0:                  float64(v.Op0),
					Op0InvoiceNum:        float64(v.Op0InvoiceNum),
					Op1NotExpenseAccount: float64(v.Op1NotExpenseAccount),
					Op1ExpenseAccount:    float64(v.Op1ExpenseAccount),
					Op2NotExpenseAccount: float64(v.Op2NotExpenseAccount),
					Op2ExpenseAccount:    float64(v.Op2ExpenseAccount),
					Question:             float64(v.Question),
				}

				SummaryItem := model.OutputStatisticsSummaryList{
					Code:             v.Code,
					NickName:         v.NickName,
					SubmitTime:       v.SubmitTime.Format("2006-01-02 15:04:05"),
					AddUpToSomething: float64(v.Mary),
				}
				SummaryItem.ProSummary = append(SummaryItem.ProSummary, SumMary)

				SummaryReturn = append(SummaryReturn, SummaryItem)
			}
		}
		//整理分页
		remaining := 0
		if int(total) > limit {
			remaining = limit - len(codes)
		} else {
			remaining = int(total) - len(codes)
		}

		pageIndex := info.PageInfo.PageIndex + 1
		k := 0
		for remaining > 0 {
			var ASummaryList []model.OutputStatisticsSummary
			newLimit := 0
			if k == 0 {
				newLimit = info.PageInfo.PageSize
			} else {
				newLimit = remaining
			}
			newOffset := newLimit * (pageIndex - 1)
			newLimit = remaining

			//第二次查询, 可能有重复数据
			if reg1.MatchString(info.Code) {
				err = global.GDb.Model(&model.OutputStatisticsSummary{}).Order("id desc").Limit(newLimit).Offset(newOffset).Where("submit_time >= ? AND submit_time <= ? AND nick_name = ? ", StartTime, EndTime, info.Code).Find(&FirstSummaryList).Error
				if err != nil {
					return err, nil, 0, nil
				}
			} else {
				err = global.GDb.Model(&model.OutputStatisticsSummary{}).Order("id desc").Limit(newLimit).Offset(newOffset).Where("submit_time >= ? AND submit_time <= ? AND code = ? ", StartTime, EndTime, info.Code).Find(&FirstSummaryList).Error
				if err != nil {
					return err, nil, 0, nil
				}
			}
			if len(FirstSummaryList) == 0 {
				break
			}

			for _, v := range ASummaryList {
				//处理相同项目编码
				existCode := false
				existPro := false
				needAppend := false
				for _, v1 := range codes {
					if v1 == v.ProCode {
						needAppend = true
						break
					}
				}
				if !needAppend {
					codes = append(codes, v.ProCode)
				}

				//处理返回的结构体
				for l, r := range SummaryReturn {
					if r.Code == v.Code {
						for k, ps := range r.ProSummary {
							if ps.ProCode == v.ProCode {
								r.ProSummary[k].Op0 += float64(v.Op0)
								r.ProSummary[k].Op0InvoiceNum += float64(v.Op0InvoiceNum)
								r.ProSummary[k].Mary += float64(v.Mary)
								r.ProSummary[k].Op1ExpenseAccount += float64(v.Op1ExpenseAccount)
								r.ProSummary[k].Op1NotExpenseAccount += float64(v.Op1NotExpenseAccount)
								r.ProSummary[k].Op2ExpenseAccount += float64(v.Op2ExpenseAccount)
								r.ProSummary[k].Op2NotExpenseAccount += float64(v.Op2NotExpenseAccount)
								r.ProSummary[k].Question += float64(v.Question)
								SummaryReturn[l].AddUpToSomething += float64(v.Mary)
								existPro = true
							}
						}
						existCode = true
						if !existPro {
							SumMary := model.Summary{
								ProCode:              v.ProCode,
								Mary:                 float64(v.Mary),
								Op0:                  float64(v.Op0),
								Op0InvoiceNum:        float64(v.Op0InvoiceNum),
								Op1NotExpenseAccount: float64(v.Op1NotExpenseAccount),
								Op1ExpenseAccount:    float64(v.Op1ExpenseAccount),
								Op2NotExpenseAccount: float64(v.Op2NotExpenseAccount),
								Op2ExpenseAccount:    float64(v.Op2ExpenseAccount),
								Question:             float64(v.Question),
							}
							SummaryReturn[l].ProSummary = append(SummaryReturn[l].ProSummary, SumMary)
						}
					}
				}

				if !existCode {
					SumMary := model.Summary{
						ProCode:              v.ProCode,
						Mary:                 float64(v.Mary),
						Op0:                  float64(v.Op0),
						Op0InvoiceNum:        float64(v.Op0InvoiceNum),
						Op1NotExpenseAccount: float64(v.Op1NotExpenseAccount),
						Op1ExpenseAccount:    float64(v.Op1ExpenseAccount),
						Op2NotExpenseAccount: float64(v.Op2NotExpenseAccount),
						Op2ExpenseAccount:    float64(v.Op2ExpenseAccount),
						Question:             float64(v.Question),
					}

					SummaryItem := model.OutputStatisticsSummaryList{
						Code:             v.Code,
						NickName:         v.NickName,
						SubmitTime:       v.SubmitTime.Format("2006-01-02 15:04:05"),
						AddUpToSomething: float64(v.Mary),
					}
					SummaryItem.ProSummary = append(SummaryItem.ProSummary, SumMary)

					SummaryReturn = append(SummaryReturn, SummaryItem)
				}
			}

			if int(total) > limit {
				remaining = limit - len(codes)
			} else {
				remaining = int(total) - len(codes)
			}
			pageIndex += 1
			k++
		}
		for _, v := range SummaryReturn[0].ProSummary {
			top = append(top, v.ProCode)
		}

	} else {
		//总数
		err = global.GDb.Model(&model.OutputStatisticsSummary{}).Where("submit_time >= ? AND submit_time <= ? ", StartTime, EndTime).Distinct("code").Count(&total).Error
		fmt.Println(total)
		if err != nil {
			return err, nil, 0, nil
		}

		//第一次查询, 可能有重复数据
		err = global.GDb.Model(&model.OutputStatisticsSummary{}).Order("id desc").Limit(limit).Offset(offset).Where("submit_time >= ? AND submit_time <= ? ", StartTime, EndTime).Find(&FirstSummaryList).Error
		if err != nil {
			return err, nil, 0, nil
		}
		//防止查到不同日期但同一人的数据
		for _, v := range FirstSummaryList {
			existCode := false
			existPro := false
			//处理相同工号
			needAppend := false
			for _, v1 := range codes {
				if v1 == v.Code {
					needAppend = true
					break
				}
			}
			if !needAppend {
				codes = append(codes, v.Code)
			}
			//处理返回结构体
			for l, r := range SummaryReturn {
				if r.Code == v.Code {
					for k, ps := range r.ProSummary {
						if ps.ProCode == v.ProCode {
							r.ProSummary[k].Op0 += float64(v.Op0)
							r.ProSummary[k].Op0InvoiceNum += float64(v.Op0InvoiceNum)
							r.ProSummary[k].Mary += float64(v.Mary)
							r.ProSummary[k].Op1ExpenseAccount += float64(v.Op1ExpenseAccount)
							r.ProSummary[k].Op1NotExpenseAccount += float64(v.Op1NotExpenseAccount)
							r.ProSummary[k].Op2ExpenseAccount += float64(v.Op2ExpenseAccount)
							r.ProSummary[k].Op2NotExpenseAccount += float64(v.Op2NotExpenseAccount)
							r.ProSummary[k].Question += float64(v.Question)
							SummaryReturn[l].AddUpToSomething += float64(v.Mary)
							existPro = true
						}
					}
					existCode = true
					if !existPro {
						SumMary := model.Summary{
							ProCode:              v.ProCode,
							Mary:                 float64(v.Mary),
							Op0:                  float64(v.Op0),
							Op0InvoiceNum:        float64(v.Op0InvoiceNum),
							Op1NotExpenseAccount: float64(v.Op1NotExpenseAccount),
							Op1ExpenseAccount:    float64(v.Op1ExpenseAccount),
							Op2NotExpenseAccount: float64(v.Op2NotExpenseAccount),
							Op2ExpenseAccount:    float64(v.Op2ExpenseAccount),
							Question:             float64(v.Question),
						}
						SummaryReturn[l].ProSummary = append(SummaryReturn[l].ProSummary, SumMary)
						SummaryReturn[l].AddUpToSomething += float64(v.Mary)
					}
				}
			}

			if !existCode {
				SumMary := model.Summary{
					ProCode:              v.ProCode,
					Mary:                 float64(v.Mary),
					Op0:                  float64(v.Op0),
					Op0InvoiceNum:        float64(v.Op0InvoiceNum),
					Op1NotExpenseAccount: float64(v.Op1NotExpenseAccount),
					Op1ExpenseAccount:    float64(v.Op1ExpenseAccount),
					Op2NotExpenseAccount: float64(v.Op2NotExpenseAccount),
					Op2ExpenseAccount:    float64(v.Op2ExpenseAccount),
					Question:             float64(v.Question),
				}

				SummaryItem := model.OutputStatisticsSummaryList{
					Code:             v.Code,
					NickName:         v.NickName,
					SubmitTime:       v.SubmitTime.Format("2006-01-02 15:04:05"),
					AddUpToSomething: float64(v.Mary),
				}
				SummaryItem.ProSummary = append(SummaryItem.ProSummary, SumMary)

				SummaryReturn = append(SummaryReturn, SummaryItem)
			}
		}
		//整理分页
		remaining := limit - len(codes)
		pageIndex := info.PageInfo.PageIndex + 1
		k := 0
		for remaining > 0 {
			var ASummaryList []model.OutputStatisticsSummary
			//limit := info.PageInfo.PageSize
			//offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
			newLimit := 0
			if k == 0 {
				newLimit = info.PageInfo.PageSize
			} else {
				newLimit = remaining
			}
			newOffset := newLimit * (pageIndex - 1)
			newLimit = remaining
			//第二次查询, 可能有重复数据
			err = global.GDb.Model(&model.OutputStatisticsSummary{}).Order("id desc").Limit(newLimit).Offset(newOffset).Where("submit_time >= ? AND submit_time <= ? ", StartTime, EndTime).Find(&ASummaryList).Error
			if err != nil {
				return err, nil, 0, nil
			}
			if len(ASummaryList) == 0 {
				break
			}

			for _, v := range ASummaryList {
				//处理相同工号
				existCode := false
				existPro := false
				needAppend := false
				for _, v1 := range codes {
					if v1 == v.Code {
						needAppend = true
						break
					}
				}
				if !needAppend {
					codes = append(codes, v.Code)
				}

				//处理返回的结构体
				for l, r := range SummaryReturn {
					if r.Code == v.Code {
						for k, ps := range r.ProSummary {
							if ps.ProCode == v.ProCode {
								r.ProSummary[k].Op0 += float64(v.Op0)
								r.ProSummary[k].Op0InvoiceNum += float64(v.Op0InvoiceNum)
								r.ProSummary[k].Mary += float64(v.Mary)
								r.ProSummary[k].Op1ExpenseAccount += float64(v.Op1ExpenseAccount)
								r.ProSummary[k].Op1NotExpenseAccount += float64(v.Op1NotExpenseAccount)
								r.ProSummary[k].Op2ExpenseAccount += float64(v.Op2ExpenseAccount)
								r.ProSummary[k].Op2NotExpenseAccount += float64(v.Op2NotExpenseAccount)
								r.ProSummary[k].Question += float64(v.Question)
								SummaryReturn[l].AddUpToSomething += float64(v.Mary)
								existPro = true
							}
						}
						existCode = true
						if !existPro {
							SumMary := model.Summary{
								ProCode:              v.ProCode,
								Mary:                 float64(v.Mary),
								Op0:                  float64(v.Op0),
								Op0InvoiceNum:        float64(v.Op0InvoiceNum),
								Op1NotExpenseAccount: float64(v.Op1NotExpenseAccount),
								Op1ExpenseAccount:    float64(v.Op1ExpenseAccount),
								Op2NotExpenseAccount: float64(v.Op2NotExpenseAccount),
								Op2ExpenseAccount:    float64(v.Op2ExpenseAccount),
								Question:             float64(v.Question),
							}
							SummaryReturn[l].AddUpToSomething += float64(v.Mary)
							SummaryReturn[l].ProSummary = append(SummaryReturn[l].ProSummary, SumMary)
						}
					}
				}

				if !existCode {
					SumMary := model.Summary{
						ProCode:              v.ProCode,
						Mary:                 float64(v.Mary),
						Op0:                  float64(v.Op0),
						Op0InvoiceNum:        float64(v.Op0InvoiceNum),
						Op1NotExpenseAccount: float64(v.Op1NotExpenseAccount),
						Op1ExpenseAccount:    float64(v.Op1ExpenseAccount),
						Op2NotExpenseAccount: float64(v.Op2NotExpenseAccount),
						Op2ExpenseAccount:    float64(v.Op2ExpenseAccount),
						Question:             float64(v.Question),
					}

					SummaryItem := model.OutputStatisticsSummaryList{
						Code:             v.Code,
						NickName:         v.NickName,
						SubmitTime:       v.SubmitTime.Format("2006-01-02 15:04:05"),
						AddUpToSomething: float64(v.Mary),
					}
					SummaryItem.ProSummary = append(SummaryItem.ProSummary, SumMary)

					SummaryReturn = append(SummaryReturn, SummaryItem)
				}
			}

			remaining = limit - len(codes)
			pageIndex += 1
			k++
		}

		//处理空缺的项目
		//0.[a], [b]
		//1.[a], [a,b]
		//2.[c], [a,b]
		//3.[a], [b,c]
		//4.[a,b], [c,d]
		//5.[a,b], [b,c]
		TopMap := make([][]string, 0)
		for _, v := range SummaryReturn {
			t := make([]string, 0)
			for _, v1 := range v.ProSummary {
				t = append(t, v1.ProCode)
			}
			TopMap = append(TopMap, t)
		}
		top = append(top, SL(TopMap)...)
		for i, v := range SummaryReturn {
			n := make([]string, len(top))
			copy(n, top)
			self := make([]string, 0)
			for _, v2 := range v.ProSummary {
				self = append(self, v2.ProCode)
			}

			fill := S(n, self)
			for _, v3 := range fill {
				SumMary := model.Summary{
					ProCode:              v3,
					Mary:                 0,
					Op0:                  0,
					Op0InvoiceNum:        0,
					Op1NotExpenseAccount: 0,
					Op1ExpenseAccount:    0,
					Op2NotExpenseAccount: 0,
					Op2ExpenseAccount:    0,
					Question:             0,
				}
				SummaryReturn[i].ProSummary = append(SummaryReturn[i].ProSummary, SumMary)
			}
		}
	}

	return err, SummaryReturn, total, top
}

func ExportOutputStatistics(info request.OutPutStatisticsExport) (err error, path, name string) {
	r := request.OutPutStatisticsSearch{
		//ProCode:    info.ProCode,
		StartTime:  info.StartTime,
		EndTime:    info.EndTime,
		IsCheckAll: 1,
	}
	err, Summary, _, top := GetExportOutputStatistics(r)
	if err != nil {
		return err, "", ""
	}
	excelLine := make([][]interface{}, 0)
	//表头
	topItem := make([]interface{}, len(top)*7+3)
	topItem[0] = "项目/工号"
	for i, proCode := range top {
		if i == 0 {
			topItem[i+3] = proCode
		} else {
			topItem[i+2+7] = proCode
		}
	}
	excelLine = append(excelLine, topItem)
	topNameItem := make([]interface{}, 0)
	topNameItem = append(topNameItem, "工号", "姓名", "合计")
	for _, t := range top {
		fmt.Println(t)
		topNameItem = append(topNameItem, "汇总", "初审", "一码非报销单", "一码报销单", "二码非报销单", "二码报销单", "问题件")
	}
	excelLine = append(excelLine, topNameItem)
	//表格
	for _, v := range Summary {
		item := make([]interface{}, 0)
		item = append(item, v.Code, v.NickName, v.AddUpToSomething)
		for _, k := range v.ProSummary {
			item = append(item, k.Mary, k.Op0, k.Op1ExpenseAccount, k.Op1NotExpenseAccount, k.Op2ExpenseAccount, k.Op2NotExpenseAccount, k.Question)
		}
		excelLine = append(excelLine, item)
	}
	StartTime := strings.Replace(info.StartTime, " 00:00:00", "", -1)
	EndTime := strings.Replace(info.EndTime, " 23:59:59", "", -1)
	bookName := "人员产量统计(全部)" + StartTime + "-" + EndTime
	err = u.ExportExcelTheMainEntrance(excelLine, bookName, info.ProCode, "output-statistics/")

	return err, "files/人员(全部)产量统计导出/" + info.ProCode + "/" + bookName + ".xlsx", bookName
}

func GetExportOutputStatistics(info request.OutPutStatisticsSearch) (err error, all []model.OutputStatisticsSummaryList, total int64, Top []string) {
	//无需分页
	reg1 := regexp.MustCompile("^([\u4e00-\u9fa5]+)$")
	fmt.Println("GetExportOutputStatistics")
	StartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.StartTime, time.Local)
	EndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.EndTime, time.Local)
	var SummaryReturn []model.OutputStatisticsSummaryList
	var FirstSummaryList []model.OutputStatisticsSummary

	top := make([]string, 0)
	codes := make([]string, 0)

	if info.Code != "" {
		if reg1.MatchString(info.Code) {
			err = global.GDb.Model(&model.OutputStatisticsSummary{}).Order("id desc").Where("submit_time >= ? AND submit_time <= ? AND nick_name = ? ", StartTime, EndTime, info.Code).Find(&FirstSummaryList).Error
			if err != nil {
				return err, nil, 0, nil
			}
			err = global.GDb.Model(&model.OutputStatisticsSummary{}).Order("id desc").Where("submit_time >= ? AND submit_time <= ? AND nick_name = ? ", StartTime, EndTime, info.Code).Distinct("pro_code").Count(&total).Error
			fmt.Println(total)
			if err != nil {
				return err, nil, 0, nil
			}
		} else {
			err = global.GDb.Model(&model.OutputStatisticsSummary{}).Order("id desc").Where("submit_time >= ? AND submit_time <= ? AND code = ? ", StartTime, EndTime, info.Code).Find(&FirstSummaryList).Error
			if err != nil {
				return err, nil, 0, nil
			}
			err = global.GDb.Model(&model.OutputStatisticsSummary{}).Order("id desc").Where("submit_time >= ? AND submit_time <= ? AND code = ? ", StartTime, EndTime, info.Code).Distinct("pro_code").Count(&total).Error
			fmt.Println(total)
			if err != nil {
				return err, nil, 0, nil
			}
		}
		if total == 0 {
			return errors.New("没有该用户"), nil, 0, nil
		}

		//防止查到不同日期但同一人的数据
		for _, v := range FirstSummaryList {
			existCode := false
			existPro := false
			//处理相同项目编码
			needAppend := false
			for _, v1 := range codes {
				if v1 == v.ProCode {
					needAppend = true
					break
				}
			}
			if !needAppend {
				codes = append(codes, v.ProCode)
			}
			//处理返回结构体
			for l, r := range SummaryReturn {
				if r.Code == v.Code {
					for k, ps := range r.ProSummary {
						if ps.ProCode == v.ProCode {
							r.ProSummary[k].Op0 += float64(v.Op0)
							r.ProSummary[k].Op0InvoiceNum += float64(v.Op0InvoiceNum)
							r.ProSummary[k].Mary += float64(v.Mary)
							r.ProSummary[k].Op1ExpenseAccount += float64(v.Op1ExpenseAccount)
							r.ProSummary[k].Op1NotExpenseAccount += float64(v.Op1NotExpenseAccount)
							r.ProSummary[k].Op2ExpenseAccount += float64(v.Op2ExpenseAccount)
							r.ProSummary[k].Op2NotExpenseAccount += float64(v.Op2NotExpenseAccount)
							r.ProSummary[k].Question += float64(v.Question)
							SummaryReturn[l].AddUpToSomething += float64(v.Mary)
							existPro = true
						}
					}
					existCode = true
					if !existPro {
						SumMary := model.Summary{
							ProCode:              v.ProCode,
							Mary:                 float64(v.Mary),
							Op0:                  float64(v.Op0),
							Op0InvoiceNum:        float64(v.Op0InvoiceNum),
							Op1NotExpenseAccount: float64(v.Op1NotExpenseAccount),
							Op1ExpenseAccount:    float64(v.Op1ExpenseAccount),
							Op2NotExpenseAccount: float64(v.Op2NotExpenseAccount),
							Op2ExpenseAccount:    float64(v.Op2ExpenseAccount),
							Question:             float64(v.Question),
						}
						SummaryReturn[l].ProSummary = append(SummaryReturn[l].ProSummary, SumMary)
					}
				}
			}

			if !existCode {
				SumMary := model.Summary{
					ProCode:              v.ProCode,
					Mary:                 float64(v.Mary),
					Op0:                  float64(v.Op0),
					Op0InvoiceNum:        float64(v.Op0InvoiceNum),
					Op1NotExpenseAccount: float64(v.Op1NotExpenseAccount),
					Op1ExpenseAccount:    float64(v.Op1ExpenseAccount),
					Op2NotExpenseAccount: float64(v.Op2NotExpenseAccount),
					Op2ExpenseAccount:    float64(v.Op2ExpenseAccount),
					Question:             float64(v.Question),
				}

				SummaryItem := model.OutputStatisticsSummaryList{
					Code:             v.Code,
					NickName:         v.NickName,
					SubmitTime:       v.SubmitTime.Format("2006-01-02 15:04:05"),
					AddUpToSomething: float64(v.Mary),
				}
				SummaryItem.ProSummary = append(SummaryItem.ProSummary, SumMary)

				SummaryReturn = append(SummaryReturn, SummaryItem)
			}
		}
		//整理分页
		for _, v := range SummaryReturn[0].ProSummary {
			top = append(top, v.ProCode)
		}

	} else {
		//总数
		err = global.GDb.Model(&model.OutputStatisticsSummary{}).Where("submit_time >= ? AND submit_time <= ? ", StartTime, EndTime).Distinct("code").Count(&total).Error
		fmt.Println(total)
		if err != nil {
			return err, nil, 0, nil
		}

		//第一次查询, 可能有重复数据
		err = global.GDb.Model(&model.OutputStatisticsSummary{}).Order("id desc").Where("submit_time >= ? AND submit_time <= ? ", StartTime, EndTime).Find(&FirstSummaryList).Error
		if err != nil {
			return err, nil, 0, nil
		}
		//防止查到不同日期但同一人的数据
		for _, v := range FirstSummaryList {
			existCode := false
			existPro := false
			//处理相同工号
			needAppend := false
			for _, v1 := range codes {
				if v1 == v.Code {
					needAppend = true
					break
				}
			}
			if !needAppend {
				codes = append(codes, v.Code)
			}
			//处理返回结构体
			for l, r := range SummaryReturn {
				if r.Code == v.Code {
					for k, ps := range r.ProSummary {
						if ps.ProCode == v.ProCode {
							r.ProSummary[k].Op0 += float64(v.Op0)
							r.ProSummary[k].Op0InvoiceNum += float64(v.Op0InvoiceNum)
							r.ProSummary[k].Mary += float64(v.Mary)
							r.ProSummary[k].Op1ExpenseAccount += float64(v.Op1ExpenseAccount)
							r.ProSummary[k].Op1NotExpenseAccount += float64(v.Op1NotExpenseAccount)
							r.ProSummary[k].Op2ExpenseAccount += float64(v.Op2ExpenseAccount)
							r.ProSummary[k].Op2NotExpenseAccount += float64(v.Op2NotExpenseAccount)
							r.ProSummary[k].Question += float64(v.Question)
							SummaryReturn[l].AddUpToSomething += float64(v.Mary)
							existPro = true
						}
					}
					existCode = true
					if !existPro {
						SumMary := model.Summary{
							ProCode:              v.ProCode,
							Mary:                 float64(v.Mary),
							Op0:                  float64(v.Op0),
							Op0InvoiceNum:        float64(v.Op0InvoiceNum),
							Op1NotExpenseAccount: float64(v.Op1NotExpenseAccount),
							Op1ExpenseAccount:    float64(v.Op1ExpenseAccount),
							Op2NotExpenseAccount: float64(v.Op2NotExpenseAccount),
							Op2ExpenseAccount:    float64(v.Op2ExpenseAccount),
							Question:             float64(v.Question),
						}
						SummaryReturn[l].ProSummary = append(SummaryReturn[l].ProSummary, SumMary)
						SummaryReturn[l].AddUpToSomething += float64(v.Mary)
					}
				}
			}

			if !existCode {
				SumMary := model.Summary{
					ProCode:              v.ProCode,
					Mary:                 float64(v.Mary),
					Op0:                  float64(v.Op0),
					Op0InvoiceNum:        float64(v.Op0InvoiceNum),
					Op1NotExpenseAccount: float64(v.Op1NotExpenseAccount),
					Op1ExpenseAccount:    float64(v.Op1ExpenseAccount),
					Op2NotExpenseAccount: float64(v.Op2NotExpenseAccount),
					Op2ExpenseAccount:    float64(v.Op2ExpenseAccount),
					Question:             float64(v.Question),
				}

				SummaryItem := model.OutputStatisticsSummaryList{
					Code:             v.Code,
					NickName:         v.NickName,
					SubmitTime:       v.SubmitTime.Format("2006-01-02 15:04:05"),
					AddUpToSomething: float64(v.Mary),
				}
				SummaryItem.ProSummary = append(SummaryItem.ProSummary, SumMary)

				SummaryReturn = append(SummaryReturn, SummaryItem)
			}
		}

		//处理空缺的项目
		//0.[a], [b]
		//1.[a], [a,b]
		//2.[c], [a,b]
		//3.[a], [b,c]
		//4.[a,b], [c,d]
		//5.[a,b], [b,c]
		TopMap := make([][]string, 0)
		for _, v := range SummaryReturn {
			t := make([]string, 0)
			for _, v1 := range v.ProSummary {
				t = append(t, v1.ProCode)
			}
			TopMap = append(TopMap, t)
		}
		top = append(top, SL(TopMap)...)
		for i, v := range SummaryReturn {
			n := make([]string, len(top))
			copy(n, top)
			self := make([]string, 0)
			for _, v2 := range v.ProSummary {
				self = append(self, v2.ProCode)
			}

			fill := S(n, self)
			for _, v3 := range fill {
				SumMary := model.Summary{
					ProCode:              v3,
					Mary:                 0,
					Op0:                  0,
					Op0InvoiceNum:        0,
					Op1NotExpenseAccount: 0,
					Op1ExpenseAccount:    0,
					Op2NotExpenseAccount: 0,
					Op2ExpenseAccount:    0,
					Question:             0,
				}
				SummaryReturn[i].ProSummary = append(SummaryReturn[i].ProSummary, SumMary)
			}
		}
	}

	return err, SummaryReturn, total, top

}

func ExportOutputStatisticsDetail(info request.OutPutStatisticsDetailExport) (err error, path, name string) {
	fmt.Println("ExportOutputStatisticsDetail", info)
	ProCode := info.ProCode
	//info.ProCode = "B01182"
	r := request.OutPutStatisticsSearch{
		ProCode:    info.ProCode,
		StartTime:  info.StartTime,
		EndTime:    info.EndTime,
		IsCheckAll: 2,
	}
	err, OutputStatisticsReturn, _, _ := GetExportOutputStatisticsDetail(r)
	if err != nil {
		return err, "", ""
	}

	s := strings.Replace(info.StartTime, " 00:00:00", "", -1)
	e := strings.Replace(info.EndTime, " 23:59:59", "", -1)
	bookName := "人员产量统计(明细)" + s + "-" + e + ".xlsx"
	// 线上 uploads/file
	basicPath := global.GConfig.LocalUpload.FilePath + "人员(明细)产量统计导出/" + ProCode + "/"
	// 本地
	//basicPath := "./"
	err = u.ExportBigExcelDetail(basicPath, bookName, "", OutputStatisticsReturn)
	return err, basicPath + bookName, bookName
}

func ExportOcrOutput(info request.GetOCROutPutStatisticsSearch) (err error, path, name string) {
	fmt.Println("ExportOcrOutput", info)
	ProCode := info.ProCode
	//info.ProCode = "B01182"
	// r := request.OutPutStatisticsSearch{
	// 	ProCode:   info.ProCode,
	// 	StartTime: info.StartTime,
	// 	EndTime:   info.EndTime,
	// }
	info.PageInfo.PageSize = 0
	err, OutputStatisticsReturn, _ := GetOcrStatistics(info)
	if err != nil {
		return err, "", ""
	}
	for ii, ocrStatics := range OutputStatisticsReturn {
		if ocrStatics.Compare == "2" {
			OutputStatisticsReturn[ii].Compare = "否"
		} else {
			OutputStatisticsReturn[ii].Compare = "是"
		}
		if ocrStatics.Disable == "2" {
			OutputStatisticsReturn[ii].Disable = "否"
		} else {
			OutputStatisticsReturn[ii].Disable = "是"
		}
	}

	s := info.StartTime.Format("2006-01-02")
	e := info.EndTime.Format("2006-01-02")
	bookName := "ocr统计统计" + s + "-" + e + ".xlsx"
	// 线上 uploads/file
	basicPath := global.GConfig.LocalUpload.FilePath + "ocr统计导出/" + ProCode + "/"
	// 本地
	//basicPath := "./"
	err = u.ExportBigExcel(basicPath, bookName, "", OutputStatisticsReturn)
	return err, basicPath + bookName, bookName
}

func GetExportOutputStatisticsDetail(info request.OutPutStatisticsSearch) (err error, all []response.OutputStatisticsRes, total int64, Top []string) {
	var OutputStatisticsRes []response.OutputStatisticsRes
	var OutputStatisticsResItem response.OutputStatisticsRes
	StartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.StartTime, time.Local)
	EndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.EndTime, time.Local)
	reg1 := regexp.MustCompile("^([\u4e00-\u9fa5]+)$")
	ProCode := info.ProCode
	db := global.ProDbMap[ProCode]
	if db == nil {
		return global.ProDbErr, nil, 0, nil
	}
	db = db.Model(&model.OutputStatistics{}).Where("submit_time >= ? AND submit_time <= ? ", StartTime, EndTime)
	if info.Code != "" {
		if reg1.MatchString(info.Code) {
			db = db.Where("nick_name = ? ", info.Code)
		} else {
			db = db.Where("code = ? ", info.Code)
		}
	}
	err = db.Order("id desc").Count(&total).Error
	if err != nil {
		return err, nil, 0, nil
	}
	var OutputStatisticsReturn []model.OutputStatistics
	err = db.Order("id desc").Find(&OutputStatisticsReturn).Error
	if err != nil {
		return err, nil, 0, nil
	}
	for _, OR := range OutputStatisticsReturn {
		InitDataTy := reflect.TypeOf(OR)
		InitDataVa := reflect.ValueOf(OR)

		ResItemTy := reflect.TypeOf(OutputStatisticsResItem)
		ResItemVa := reflect.ValueOf(&OutputStatisticsResItem).Elem()

		for i := 0; i < InitDataTy.NumField(); i++ {
			InitDataTyFieldName := InitDataTy.Field(i).Name
			InitDataTyFieldType := InitDataTy.Field(i).Type
			if InitDataTyFieldName == "Model" {
				continue
			}
			for j := 0; j < ResItemTy.NumField(); j++ {
				ResItemTyFieldName := ResItemTy.Field(j).Name
				if InitDataTyFieldName == ResItemTyFieldName && strings.Index(ResItemTyFieldName, "CostTime") == -1 {
					if InitDataTyFieldType.Kind() == 2 {
						ResItemVa.FieldByName(ResItemTyFieldName).SetInt(int64(InitDataVa.FieldByName(InitDataTyFieldName).Interface().(int)))
						break
					}
					if InitDataTyFieldType.Kind() == 14 {
						if strings.Index(InitDataTyFieldName, "AccuracyRate") == -1 && strings.Index(InitDataTyFieldName, "QuestionMarkProportion") == -1 {
							fmt.Println(InitDataTyFieldName)
							ResItemVa.FieldByName(ResItemTyFieldName).SetFloat(InitDataVa.FieldByName(InitDataTyFieldName).Interface().(float64))
						} else {
							if InitDataVa.FieldByName(InitDataTyFieldName).Interface().(float64) == 100 {
								ResItemVa.FieldByName(ResItemTyFieldName).SetString(decimal.NewFromFloat(InitDataVa.FieldByName(InitDataTyFieldName).Interface().(float64)).RoundCeil(2).String() + "%")
							} else {
								ResItemVa.FieldByName(ResItemTyFieldName).SetString(decimal.NewFromFloat(InitDataVa.FieldByName(InitDataTyFieldName).Interface().(float64)).Mul(decimal.NewFromFloat(100)).RoundCeil(2).String() + "%")
							}

						}
						break
					}
					if InitDataTyFieldType.Kind() == 24 {
						//ResItemVa.SetString()
						ResItemVa.FieldByName(ResItemTyFieldName).SetString(InitDataVa.FieldByName(InitDataTyFieldName).Interface().(string))
						break
					}
					if InitDataTyFieldName == "SubmitTime" {
						OutputStatisticsResItem.SubmitTime = OR.SubmitTime.Format("2006-01-02")
						break
					}
				} else if InitDataTyFieldName == ResItemTyFieldName && strings.Index(ResItemTyFieldName, "CostTime") != -1 {
					costTime := InitDataVa.FieldByName(InitDataTyFieldName).Interface().(int64)
					_, cost := J(int(costTime))
					ResItemVa.FieldByName(ResItemTyFieldName).SetString(cost)
					break
				}
			}
		}
		var user modelBase.SysUser
		err = global.GDb.Model(&modelBase.SysUser{}).Where("code = ? ", OR.Code).Find(&user).Error
		if err != nil {
			return err, OutputStatisticsRes, total, nil
		}
		OutputStatisticsResItem.NickName = user.NickName
		//removalTime := RemovalTime(OutputStatisticsResItem.Op0CostTime)
		//OutputStatisticsResItem.Op0CostTime = removalTime //初审花费时间
		//op1CostTime := RemovalTime(OutputStatisticsResItem.Op1CostTime)
		//OutputStatisticsResItem.Op1CostTime = op1CostTime //一码花费时间
		//op2CostTime := RemovalTime(OutputStatisticsResItem.Op2CostTime)
		//OutputStatisticsResItem.Op2CostTime = op2CostTime //二码花费时间
		//opQCostTime := RemovalTime(OutputStatisticsResItem.OpQCostTime)
		//OutputStatisticsResItem.OpQCostTime = opQCostTime //问题件花费时间
		//summaryCostTime := RemovalTime(OutputStatisticsResItem.SummaryCostTime)
		//OutputStatisticsResItem.SummaryCostTime = summaryCostTime //汇总花费时间
		OutputStatisticsRes = append(OutputStatisticsRes, OutputStatisticsResItem)
	}
	return err, OutputStatisticsRes, total, nil
}

func DeleteOutputStatisticsDetail(info request.OutPutStatisticsSearch) (err error) {
	db := global.ProDbMap[info.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	StartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.StartTime, time.Local)
	EndTime, _ := time.ParseInLocation("2006-01-02 15:04:05", info.EndTime, time.Local)
	err = tx.Where("code = ? AND submit_time >= ? AND submit_time <= ? ", info.Code, StartTime, EndTime).Delete(&model.OutputStatistics{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Where("code = ? AND submit_day >= ? AND submit_day <= ? ", info.Code, StartTime, EndTime).Delete(&model.Wrong{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func GetCorrected(info modelBase.BasePageInfo) (CorrectedArr []model.SysCorrectedRep, count int64, err error) {
	var CorrectedItem []model.SysCorrected
	limit := info.PageSize
	offset := info.PageSize * (info.PageIndex - 1)
	err = global.GDb.Model(&model.SysCorrected{}).Order("id desc").Limit(limit).Offset(offset).Find(&CorrectedItem).Count(&count).Error
	var SysCorrectedRepArr []model.SysCorrectedRep
	var SysCorrectedRepItem model.SysCorrectedRep
	for _, v := range CorrectedItem {
		SysCorrectedRepItem.ID = v.ID
		SysCorrectedRepItem.CreatedAt = v.CreatedAt
		SysCorrectedRepItem.UpdatedAt = v.UpdatedAt
		SysCorrectedRepItem.ProCode = v.ProCode
		SysCorrectedRepItem.Op0AsTheBlock = v.Op0AsTheBlock
		SysCorrectedRepItem.Op0AsTheInvoice = v.Op0AsTheInvoice
		SysCorrectedRepItem.Op1NotExpenseAccount = v.Op1NotExpenseAccount
		SysCorrectedRepItem.Op1ExpenseAccount = v.Op1ExpenseAccount
		SysCorrectedRepItem.Op2NotExpenseAccount = v.Op2NotExpenseAccount
		SysCorrectedRepItem.Op2ExpenseAccount = v.Op2ExpenseAccount
		SysCorrectedRepItem.Question = v.Question
		SysCorrectedRepItem.StartTime = v.StartTime.Format("2006-01-02")

		SysCorrectedRepArr = append(SysCorrectedRepArr, SysCorrectedRepItem)
	}
	return SysCorrectedRepArr, count, err
}

func InsertCorrected(InsertCorrected model.SysCorrected) (err error) {
	if InsertCorrected.Op0AsTheBlock > 0 && InsertCorrected.Op0AsTheInvoice > 0 {
		return errors.New("配置规则：初审（按分块）和初审（按发票）只能配置其中一个")
	}
	if (InsertCorrected.Op0AsTheBlock < 0 && InsertCorrected.Op0AsTheInvoice == 0) || (InsertCorrected.Op0AsTheInvoice < 0 && InsertCorrected.Op0AsTheBlock == 0) {
		return errors.New("折算比配置不符合规则")
	}
	if InsertCorrected.Op0AsTheBlock < 0 && InsertCorrected.Op0AsTheInvoice > 0 {
		return errors.New("配置规则：初审（按分块）不配置默认为0")
	}
	if InsertCorrected.Op0AsTheInvoice < 0 && InsertCorrected.Op0AsTheBlock > 0 {
		return errors.New("配置规则：初审（按发票）不配置默认为0")
	}
	if InsertCorrected.Op0AsTheBlock == 0 && InsertCorrected.Op0AsTheInvoice == 0 {
		return errors.New("配置规则：初审（按分块）和初审（按发票）必须配置其中一个, 并且不能配置为0")
	}
	var count int64
	err = global.GDb.Model(&model.SysCorrected{}).Where("pro_code = ? ", InsertCorrected.ProCode).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("已有项目配置")
	}
	err = global.GDb.Model(&model.SysCorrected{}).Create(&InsertCorrected).Error
	return err
}

func UpdateCorrected(UpdateCorrected request.UpdateCorrected) (err error) {
	var count int64
	var oldCorrected model.SysCorrected
	err = global.GDb.Model(&model.SysCorrected{}).Where("id = ? ", UpdateCorrected.Id).Find(&oldCorrected).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("该折算比例不存在，无法更新")
	}

	obj1 := reflect.TypeOf(oldCorrected)
	obj1V := reflect.ValueOf(oldCorrected)
	obj1TypeOfT := obj1V.Type()
	obj2 := reflect.TypeOf(UpdateCorrected)
	obj2V := reflect.ValueOf(UpdateCorrected)

	editLog := make(map[string]interface{})
	update := make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		for j := 0; j < obj2.NumField(); j++ {
			if obj1.Field(i).Name == obj2.Field(j).Name {
				fmt.Println("1", obj1.Field(i).Name)
				fmt.Println("2", obj2.Field(j).Name)
				if obj1.Field(i).Name == "StartTime" {
					if obj1V.Field(i).Interface().(time.Time).Format("2006-01-02") != obj2V.Field(j).Interface().(string) {
						editLog[obj1TypeOfT.Field(i).Tag.Get("editName")] = obj1V.Field(i).Interface().(time.Time).Format("2006-01-02") + "," + obj2V.Field(j).Interface().(string)
						update[obj1.Field(i).Name] = obj2V.Field(j).Interface().(string)
					}
					continue
				}

				if obj1.Field(i).Name == "ProCode" {
					if obj1V.Field(i).Interface().(string) != obj2V.Field(j).Interface().(string) {
						editLog[obj1TypeOfT.Field(i).Tag.Get("editName")] = obj1V.Field(i).Interface().(string) + "," + obj2V.Field(j).Interface().(string)
						update[obj1.Field(i).Name] = obj2V.Field(j).Interface().(string)
					}
					continue
				}
				if obj1.Field(i).Name != "Id" || obj1.Field(i).Name != "Code" || obj1.Field(i).Name != "Name" {
					if obj1V.Field(i).Interface().(float64) != obj2V.Field(j).Interface().(float64) {
						editLog[obj1TypeOfT.Field(i).Tag.Get("editName")] = fmt.Sprintf("%f", obj1V.Field(i).Interface().(float64)) + "," + fmt.Sprintf("%f", Decimal(obj2V.Field(j).Interface().(float64)))
						update[obj1.Field(i).Name] = Decimal(obj2V.Field(j).Interface().(float64))
					}
				}

			}
		}

	}

	err = global.GDb.Model(&model.SysCorrected{}).Where("id = ? ", UpdateCorrected.Id).
		Updates(update).Error
	var Logs []model.SysCorrectedEditLog
	var Log model.SysCorrectedEditLog
	for i, v := range editLog {
		fmt.Println(i)
		arr := strings.Split(v.(string), ",")
		Log.CorrectedID = oldCorrected.ID
		Log.EditName = i
		A1, _ := strconv.ParseFloat(arr[0], 64)
		Log.BeforeEdit = fmt.Sprintf("%.2f", A1)
		A2, _ := strconv.ParseFloat(arr[1], 64)
		Log.AfterEdit = fmt.Sprintf("%.2f", A2)
		Log.EditedCode = UpdateCorrected.Code
		Log.EditedName = UpdateCorrected.Name

		Logs = append(Logs, Log)
	}
	err = global.GDb.Model(&model.SysCorrectedEditLog{}).Create(&Logs).Error
	if err != nil {
		return err
	}
	return err
}

func DeleteCorrected(deleteCorrectedArr request.DeleteCorrectedArr) (err error) {
	var Corrected model.SysCorrected
	var CorrectedEditLog model.SysCorrectedEditLog
	for _, id := range deleteCorrectedArr.Ids {
		err = global.GDb.Model(&model.SysCorrected{}).Where("id = ? ", id).Delete(&Corrected).Error
		if err != nil {
			return err
		}
		err = global.GDb.Model(&model.SysCorrectedEditLog{}).Where("corrected_id = ? ", id).Delete(&CorrectedEditLog).Error
		if err != nil {
			return err
		}
	}
	return err
}

func GetUpdateCorrectedLog(CorrectedID string) (Loga []model.SysCorrectedEditLogReq, count int64, err error) {
	var Log []model.SysCorrectedEditLog
	err = global.GDb.Model(&model.SysCorrectedEditLog{}).Where("corrected_id = ? ", CorrectedID).Find(&Log).Count(&count).Error
	var SysCorrectedEditLogReqArr []model.SysCorrectedEditLogReq
	var SysCorrectedEditLogReqItem model.SysCorrectedEditLogReq

	for _, v := range Log {
		SysCorrectedEditLogReqItem.ID = v.ID
		SysCorrectedEditLogReqItem.CreatedAt = v.CreatedAt.Format("2006-01-02 15:04:05")
		SysCorrectedEditLogReqItem.UpdatedAt = v.UpdatedAt.Format("2006-01-02 15:04:05")
		SysCorrectedEditLogReqItem.CorrectedID = v.CorrectedID
		SysCorrectedEditLogReqItem.EditName = v.EditName
		SysCorrectedEditLogReqItem.BeforeEdit = v.BeforeEdit
		SysCorrectedEditLogReqItem.AfterEdit = v.AfterEdit
		SysCorrectedEditLogReqItem.EditedCode = v.EditedCode
		SysCorrectedEditLogReqItem.EditedName = v.EditedName

		SysCorrectedEditLogReqArr = append(SysCorrectedEditLogReqArr, SysCorrectedEditLogReqItem)
	}
	return SysCorrectedEditLogReqArr, count, err
}

func J(second int) (error, string) {
	if second < 60 {
		return nil, strconv.Itoa(second) + "秒"
	}
	if second < 3600 && second >= 60 {
		m := second / 60
		remaining := second - m*60
		return nil, strconv.Itoa(m) + "分" + strconv.Itoa(remaining) + "秒"
	}
	if second >= 3600 {
		h := second / 3600
		mRemaining := second - h*3600
		m := mRemaining / 60
		s := mRemaining - m*60
		return nil, strconv.Itoa(h) + "小时" + strconv.Itoa(m) + "分" + strconv.Itoa(s) + "秒"
	}
	return nil, ""
}

func SL(arr [][]string) (top []string) {
	//0.[a], [b]
	//1.[a], [a,b]
	//2.[c], [a,b]
	//3.[a], [b,c]
	//4.[a,b], [c,d]
	//5.[a,b], [b,c]
	m := make(map[string]bool, 0)
	for _, v := range arr {
		for _, v1 := range v {
			if _, ok := m[v1]; !ok {
				m[v1] = true
			}
		}
	}

	for i, _ := range m {
		top = append(top, i)
	}

	return top
}

func S(t []string, value []string) (fill []string) {
	//0.[a], [b]
	//1.[a], [a,b]
	//2.[c], [a,b]
	//3.[a], [b,c]
	//4.[a,b], [c,d]
	//5.[a,b], [b,c]
	for _, v := range t {
		q := true
		for _, v1 := range value {
			if v == v1 {
				q = false
			}
		}
		if q {
			fill = append(fill, v)
		}
	}
	return fill
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

// 时分秒转换成 1.2.2
func RemovalTime(str string) string {
	hours := regexp.MustCompile(`\d+小时`).MatchString(str)
	part := regexp.MustCompile(`\d+分`).MatchString(str)

	if hours {
		re := regexp.MustCompile(`(\d+)小时(\d+)分(\d+)秒`)
		matches := re.FindStringSubmatch(str)
		if len(matches) == 0 {
			fmt.Println("无法匹配时间格式")
		}
		return matches[1] + "." + matches[2]
	} else if part {
		res := regexp.MustCompile(`(\d+)分(\d+)秒`)
		matches := res.FindStringSubmatch(str)
		if len(matches) == 0 {
			fmt.Println("无法匹配时间格式")

		}
		return matches[1]
	}
	return "0"
}
