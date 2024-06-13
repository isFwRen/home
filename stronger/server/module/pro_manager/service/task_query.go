package service

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"server/global"
	bf "server/module/load/model"
	load "server/module/load/model"
	billmap "server/module/pro_manager/const_data"
	m "server/module/pro_manager/model"
	model2 "server/module/pro_manager/model"
	res "server/module/pro_manager/model/request"
	resp "server/module/pro_manager/model/response"
	"server/module/pro_manager/project/B0108"
	"server/module/pro_manager/project/B0113"
	"server/module/pro_manager/project/B0118"
	"server/module/sys_base/model"
	"server/utils"
	"sort"
	"strings"
	"time"
	model3 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/model"
)

const (
	AgingMsg = "请!滚回去配置对应项目的时效!"
)

func GetTaskListDetail(info res.GetTaskListDetail) (err error, list interface{}, total int64) {
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	db := global.ProDbMap[info.ProCode+"_task"]
	if db == nil {
		return global.ProDbErr, nil, 0
	}
	//判断传过来的数据是否符合要求
	//1、一码二码必带IsExpenseAccount(1|2)参数
	//2、初审问题件不能带IsExpenseAccount(1|2)参数
	if (info.Op == "op0" || info.Op == "opq") && info.IsExpenseAccount != 0 {
		fmt.Println("初审问题件不能带IsExpenseAccount(1|2)参数")
		return errors.New("获取失败"), nil, 0
	}
	if (info.Op == "op1" || info.Op == "op2") && info.IsExpenseAccount == 0 {
		fmt.Println("一码二码必带IsExpenseAccount(1|2)参数")
		return errors.New("获取失败"), nil, 0
	}
	//统计总数量, 构造sql  to_char(submit_time,'YYYY-MM-DD')
	query := ""
	if info.OpStage == 1 {
		//待分配
		if info.IsExpenseAccount == 1 {
			//一码 二码 报销单
			if info.Op == "op1" {
				query = "op1_stage = 'op1' AND (op1_submit_at is null or to_char(op1_submit_at, 'YYYY-MM-DD') = '0001-01-01') AND name like '%报销单%'"
			} else if info.Op == "op2" {
				query = "op2_stage = 'op2' AND (op2_submit_at is null or to_char(op2_submit_at, 'YYYY-MM-DD') = '0001-01-01') AND name like '%报销单%'"
			}
		} else if info.IsExpenseAccount == 2 {
			//一码 二码 非报销单
			if info.Op == "op1" {
				query = "op1_stage = 'op1' AND (op1_submit_at is null or to_char(op1_submit_at, 'YYYY-MM-DD') = '0001-01-01') AND name not like '%报销单%'"
			} else if info.Op == "op2" {
				query = "op2_stage = 'op2' AND (op2_submit_at is null or to_char(op2_submit_at, 'YYYY-MM-DD') = '0001-01-01') AND name not like '%报销单%'"
			}
		} else {
			//初审 问题件
			if info.Op == "op0" {
				query = "op0_stage = 'op0' AND (op0_submit_at is null or to_char(op0_submit_at, 'YYYY-MM-DD') = '0001-01-01')"
			} else if info.Op == "opq" {
				query = "opq_stage = 'opq' AND (opq_submit_at is null or to_char(opq_submit_at, 'YYYY-MM-DD') = '0001-01-01')"
			}
		}
	} else if info.OpStage == 2 {
		//已分配
		if info.IsExpenseAccount == 1 {
			if info.Op == "op1" {
				//一码 二码 报销单
				query = "op1_stage = 'op1Cache' AND (op1_submit_at is null or to_char(op1_submit_at, 'YYYY-MM-DD') = '0001-01-01') AND name like '%报销单%'"
			} else if info.Op == "op2" {
				query = "op2_stage = 'op2Cache' AND (op2_submit_at is null or to_char(op2_submit_at, 'YYYY-MM-DD') = '0001-01-01') AND name like '%报销单%'"
			}
		} else if info.IsExpenseAccount == 2 {
			//一码 二码 非报销单
			if info.Op == "op1" {
				query = "op1_stage = 'op1Cache' AND (op1_submit_at is null or to_char(op1_submit_at, 'YYYY-MM-DD') = '0001-01-01') AND name not like '%报销单%'"
			} else if info.Op == "op2" {
				query = "op2_stage = 'op2Cache' AND (op2_submit_at is null or to_char(op2_submit_at, 'YYYY-MM-DD') = '0001-01-01') AND name not like '%报销单%'"
			}
		} else {
			//初审 问题件
			if info.Op == "op0" {
				query = "op0_stage = 'op0Cache' AND (op0_submit_at is null or to_char(op0_submit_at, 'YYYY-MM-DD') = '0001-01-01')"
			} else if info.Op == "opq" {
				query = "opq_stage = 'opqCache' AND (opq_submit_at is null or to_char(opq_submit_at, 'YYYY-MM-DD') = '0001-01-01')"
			}
		}
	} else if info.OpStage == 3 {
		//缓存区
		if info.IsExpenseAccount == 1 {
			//一码 二码 报销单
			if info.Op == "op1" {
				query = "op1_stage = 'op1Cache' AND op1_submit_at is not null AND to_char(op1_submit_at, 'YYYY-MM-DD') != '0001-01-01' AND name like '%报销单%'"
			} else if info.Op == "op2" {
				query = "op2_stage = 'op2Cache' AND op2_submit_at is not null AND to_char(op2_submit_at, 'YYYY-MM-DD') != '0001-01-01' AND name like '%报销单%'"
			}
		} else if info.IsExpenseAccount == 2 {
			//一码 二码 非报销单
			if info.Op == "op1" {
				query = "op1_stage = 'op1Cache' AND op1_submit_at is not null AND to_char(op1_submit_at, 'YYYY-MM-DD') != '0001-01-01' AND name not like '%报销单%'"
			} else if info.Op == "op2" {
				query = "op2_stage = 'op2Cache' AND op2_submit_at is not null AND to_char(op2_submit_at, 'YYYY-MM-DD') != '0001-01-01' AND name not like '%报销单%'"
			}
		} else {
			//初审 问题件
			if info.Op == "op0" {
				query = "op0_stage = 'op0Cache' AND op0_submit_at is not null AND to_char(op0_submit_at, 'YYYY-MM-DD') != '0001-01-01'"
			} else if info.Op == "opq" {
				query = "opq_stage = 'opqCache' AND opq_submit_at is not null AND to_char(opq_submit_at, 'YYYY-MM-DD') != '0001-01-01'"
			}
		}
	}
	fmt.Println("query", query)

	//newQuery := strings.Replace(strings.Replace(query, info.Op+"_stage", "x."+info.Op+"_stage", -1), info.Op+"_submit_at", "x."+info.Op+"_submit_at", -1)
	//fmt.Println(newQuery)

	err = db.Model(&load.ProjectBlock{}).Where(query).Count(&total).Error
	if err != nil {
		return err, nil, 0
	}

	//查找数据 + 整理数据
	var detail []model2.TaskListDetailNotAssign
	var detailItem model2.TaskListDetailNotAssign

	var detailAssignAndCache []model2.TaskListDetailAssignAndCache
	var detailItemAssignAndCache model2.TaskListDetailAssignAndCache

	var blocks []load.ProjectBlock

	if info.OpStage == 1 {
		//待分配
		err = db.Model(&load.ProjectBlock{}).Order("id desc").Limit(limit).Offset(offset).Where(query).Find(&blocks).Error
		if err != nil {
			return err, nil, 0
		}

		//err = db.Raw("SELECT x.bill_id, x.name, y.bill_name, y.agency from project_blocks as x left join project_bills as y on x.bill_id = y.id where " + newQuery).Scan(&b).Error
		//if err != nil {
		//	return err, nil, 0
		//}

		for _, v := range blocks {
			//fmt.Println("v", v.BillID)
			var bill m.ProjectBill
			detailItem.BlockId = v.ID
			detailItem.TaskAssign = ""
			detailItem.BlockName = v.Name
			err = db.Model(&m.ProjectBill{}).Where("id = ? ", v.BillID).Find(&bill).Error
			if err != nil {
				return err, nil, 0
			}
			//fmt.Println(bill)
			detailItem.Agency = bill.Agency
			detailItem.BillNum = bill.BillNum
			if v.Name == "未定义" {
				detailItem.TempType = "MB001"
			} else {
				detailItem.TempType = "MB002"
			}
			detail = append(detail, detailItem)
		}
		return err, detail, total
	} else {
		//已分配 缓存区
		if info.Op == "op0" {
			err = db.Model(&load.ProjectBlock{}).Order("op0_apply_at asc").Limit(limit).Offset(offset).Where(query).Find(&blocks).Error

		}
		if info.Op == "op1" {
			err = db.Model(&load.ProjectBlock{}).Order("op1_apply_at asc").Limit(limit).Offset(offset).Where(query).Find(&blocks).Error

		}
		if info.Op == "op2" {
			err = db.Model(&load.ProjectBlock{}).Order("op2_apply_at asc").Limit(limit).Offset(offset).Where(query).Find(&blocks).Error

		}
		if info.Op == "opq" {
			err = db.Model(&load.ProjectBlock{}).Order("opq_apply_at asc").Limit(limit).Offset(offset).Where(query).Find(&blocks).Error

		}
		//err = db.Model(&load.ProjectBlock{}).Order("op0_apply_at desc").Order("op1_apply_at desc").Order("op2_apply_at desc").Order("opq_apply_at desc").Limit(limit).Offset(offset).Where(query).Find(&blocks).Error
		if err != nil {
			return err, nil, 0
		}
		zeroTime := "0001-01-01"
		for _, v := range blocks {
			var bill m.ProjectBill
			detailItemAssignAndCache.BlockId = v.ID
			detailItemAssignAndCache.BlockName = v.Name
			if v.Op1ApplyAt.Format("2006-01-02") != zeroTime {
				detailItemAssignAndCache.Op1ApplyAt = v.Op1ApplyAt.Format("2006-01-02 15:04:05")
				detailItemAssignAndCache.Op1SubmitAt = v.Op1SubmitAt.Format("2006-01-02 15:04:05")
				detailItemAssignAndCache.Op1Code = v.Op1Code + global.UserCodeName[v.Op1Code]
			}

			if v.Op2ApplyAt.Format("2006-01-02") != zeroTime {
				detailItemAssignAndCache.Op2ApplyAt = v.Op2ApplyAt.Format("2006-01-02 15:04:05")
				detailItemAssignAndCache.Op2SubmitAt = v.Op2SubmitAt.Format("2006-01-02 15:04:05")
				detailItemAssignAndCache.Op2Code = v.Op2Code + global.UserCodeName[v.Op2Code]
			}

			if v.OpqApplyAt.Format("2006-01-02") != zeroTime {
				detailItemAssignAndCache.OpqApplyAt = v.OpqApplyAt.Format("2006-01-02 15:04:05")
				detailItemAssignAndCache.OpqSubmitAt = v.OpqSubmitAt.Format("2006-01-02 15:04:05")
				detailItemAssignAndCache.OpqCode = v.OpqCode + global.UserCodeName[v.OpqCode]
			}

			if v.Op0ApplyAt.Format("2006-01-02") != zeroTime {
				detailItemAssignAndCache.Op0ApplyAt = v.Op0ApplyAt.Format("2006-01-02 15:04:05")
				detailItemAssignAndCache.Op0SubmitAt = v.Op0SubmitAt.Format("2006-01-02 15:04:05")
				detailItemAssignAndCache.Op0Code = v.Op0Code + global.UserCodeName[v.Op0Code]
			}

			err = db.Model(&m.ProjectBill{}).Where("id = ? ", v.BillID).Find(&bill).Error
			if err != nil {
				return err, nil, 0
			}
			detailItemAssignAndCache.BillNum = bill.BillNum
			detailItemAssignAndCache.Agency = bill.Agency

			detailAssignAndCache = append(detailAssignAndCache, detailItemAssignAndCache)
		}
		return err, detailAssignAndCache, total
	}
}

func GetTaskList(info res.GetTaskList) (err error, list interface{}, total int64) {
	db := global.ProDbMap[info.ProCode+"_task"]
	if db == nil {
		return global.ProDbErr, nil, 0
	}
	var bills []m.ProjectBill
	err = db.Model(&m.ProjectBill{}).Find(&bills).Count(&total).Error
	fmt.Println("GetTaskList 单量", total)
	if err != nil {
		return err, nil, 0
	}
	var taskList m.TaskList
	taskList.Num = total

	//2.0
	//	1: "紧急件",
	var UrgentBill []string
	err = db.Model(&m.ProjectBill{}).Where("stick_level = '1'").Pluck("bill_num", &UrgentBill).Count(&total).Error
	fmt.Println("GetTaskList 紧急单", UrgentBill, total)
	if err != nil {
		return err, nil, 0
	}
	//	2: "优先件",
	var PriorityBill []string
	err = db.Model(&m.ProjectBill{}).Where("stick_level = '2'").Pluck("bill_num", &PriorityBill).Count(&total).Error
	fmt.Println("GetTaskList 优先单", PriorityBill, total)
	if err != nil {
		return err, nil, 0
	}

	err, agency := utils.GetRedisAgency(info.ProCode)
	if err != nil {
		return err, nil, 0
	}
	taskList.Agency = append(taskList.Agency, agency...)
	taskList.Urgent = append(taskList.Urgent, UrgentBill...)
	taskList.Priority = append(taskList.Priority, PriorityBill...)

	Op0Query := "SELECT count(case when op0_stage = 'op0' AND (op0_submit_at is null or to_char(op0_submit_at, 'YYYY-MM-DD') = '0001-01-01') then 0 end ) AS Op0notassign,\n       " +
		"count(case  when op0_stage = 'op0Cache' AND (op0_submit_at is null or to_char(op0_submit_at, 'YYYY-MM-DD') = '0001-01-01') then 0 end) AS Op0assign,\n       " +
		"count(case when op0_stage = 'op0Cache' AND op0_submit_at is not null AND to_char(op0_submit_at, 'YYYY-MM-DD') != '0001-01-01' then 0 end ) AS Op0cache\n       " +
		"from project_blocks\n"
	err = db.Model(&load.ProjectBlock{}).Raw(Op0Query).Scan(&taskList).Error
	if err != nil {
		return err, nil, 0
	}

	Op1Query := "SELECT count(case when op1_stage = 'op1' AND (op1_submit_at is null or to_char(op1_submit_at, 'YYYY-MM-DD') = '0001-01-01') AND name like '%报销单%' then 0 end ) AS Op1notassignexpenseaccount,\n       " +
		"count(case  when op1_stage = 'op1' AND (op1_submit_at is null or to_char(op1_submit_at, 'YYYY-MM-DD') = '0001-01-01') AND name not like '%报销单%' then 0 end) AS Op1notassignnotexpenseaccount,\n       " +
		"count(case when op1_stage = 'op1Cache' AND (op1_submit_at is null or to_char(op1_submit_at, 'YYYY-MM-DD') = '0001-01-01') AND name like '%报销单%' then 0 end ) AS Op1assignexpenseaccount,\n       " +
		"count(case when op1_stage = 'op1Cache' AND (op1_submit_at is null or to_char(op1_submit_at, 'YYYY-MM-DD') = '0001-01-01') AND name not like '%报销单%' then 0 end ) AS Op1assignnotexpenseaccount,\n       " +
		"count(case when op1_stage = 'op1Cache' AND op1_submit_at is not null AND to_char(op1_submit_at, 'YYYY-MM-DD') != '0001-01-01' AND name like '%报销单%' then 0 end) AS Op1cacheexpenseaccount,\n       " +
		"count(case when op1_stage = 'op1Cache' AND op1_submit_at is not null AND to_char(op1_submit_at, 'YYYY-MM-DD') != '0001-01-01' AND name not like '%报销单%' then 0 end ) AS Op1cachenotexpenseaccount\n" +
		"from project_blocks\n"
	err = db.Model(&load.ProjectBlock{}).Raw(Op1Query).Scan(&taskList).Error
	if err != nil {
		return err, nil, 0
	}

	Op2Query := "SELECT count(case when op2_stage = 'op2' AND (op2_submit_at is null or to_char(op2_submit_at, 'YYYY-MM-DD') = '0001-01-01') AND name like '%报销单%' then 0 end ) AS Op2notassignexpenseaccount,\n       " +
		"count(case  when op2_stage = 'op2' AND (op2_submit_at is null or to_char(op2_submit_at, 'YYYY-MM-DD') = '0001-01-01') AND name not like '%报销单%' then 0 end) AS Op2notassignnotexpenseaccount,\n       " +
		"count(case when op2_stage = 'op2Cache' AND (op2_submit_at is null or to_char(op2_submit_at, 'YYYY-MM-DD') = '0001-01-01') AND name like '%报销单%' then 0 end ) AS Op2assignexpenseaccount,\n       " +
		"count(case when op2_stage = 'op2Cache' AND (op2_submit_at is null or to_char(op2_submit_at, 'YYYY-MM-DD') = '0001-01-01') AND name not like '%报销单%' then 0 end ) AS Op2assignnotexpenseaccount,\n       " +
		"count(case when op2_stage = 'op2Cache' AND op2_submit_at is not null AND to_char(op2_submit_at, 'YYYY-MM-DD') != '0001-01-01' AND name like '%报销单%' then 0 end) AS Op2cacheexpenseaccount,\n       " +
		"count(case when op2_stage = 'op2Cache' AND op2_submit_at is not null AND to_char(op2_submit_at, 'YYYY-MM-DD') != '0001-01-01' AND name not like '%报销单%' then 0 end ) AS Op2cachenotexpenseaccount\n" +
		"from project_blocks\n"
	err = db.Model(&load.ProjectBlock{}).Raw(Op2Query).Scan(&taskList).Error
	if err != nil {
		return err, nil, 0
	}

	OpQQuery := "SELECT count(case when opq_stage = 'opq' AND (opq_submit_at is null or to_char(opq_submit_at, 'YYYY-MM-DD') = '0001-01-01') then 0 end ) AS Opqnotassign,\n       " +
		"count(case  when opq_stage = 'opqCache' AND (opq_submit_at is null or to_char(opq_submit_at, 'YYYY-MM-DD') = '0001-01-01') then 0 end) AS Opqassign,\n       " +
		"count(case when opq_stage = 'opqCache' AND opq_submit_at is not null AND to_char(opq_submit_at, 'YYYY-MM-DD') != '0001-01-01' then 0 end ) AS Opqcache\n       " +
		"from project_blocks\n"
	err = db.Model(&load.ProjectBlock{}).Raw(OpQQuery).Scan(&taskList).Error
	if err != nil {
		return err, nil, 0
	}

	//1.0
	//var blocks []load.ProjectBlock
	//for _, bill := range bills {
	//	//	1: "紧急件",
	//	//	2: "优先件",
	//	//if bill.StickLevel == 1 {
	//	//	taskList.Urgent = append(taskList.Urgent, bill.BillName)
	//	//}
	//	//if bill.StickLevel == 2 {
	//	//	taskList.Priority = append(taskList.Priority, bill.BillName)
	//	//	taskList.Agency = append(taskList.Agency, bill.Agency)
	//	//}
	//	err = db.Model(&load.ProjectBlock{}).Where("bill_id = ? ", bill.ID).Find(&blocks).Error
	//	if err != nil {
	//		return err, nil, 0
	//	}
	//	for _, block := range blocks {
	//		//----------------------------------------------待分配---------------------------------------------------------
	//		//op0
	//		if block.Op0Stage == "op0" {
	//			taskList.Op0NotAssign += 1
	//		}
	//		//op1
	//		if block.Op1Stage == "op1" {
	//			if strings.Index(block.Name, "报销单") == -1 {
	//				//非报销单
	//				taskList.Op1NotAssignNotExpenseAccount += 1
	//			} else {
	//				taskList.Op1NotAssignExpenseAccount += 1
	//			}
	//		}
	//		//op2
	//		if block.Op2Stage == "op2" {
	//			if strings.Index(block.Name, "报销单") == -1 {
	//				//非报销单
	//				taskList.Op2NotAssignNotExpenseAccount += 1
	//			} else {
	//				taskList.Op2NotAssignExpenseAccount += 1
	//			}
	//		}
	//		//opq
	//		if block.OpqStage == "opq" {
	//			taskList.OpqNotAssign += 1
	//		}
	//
	//		//------------------------------------------------已分配|缓存区---------------------------------------
	//		//op0
	//		if block.Op0Stage == "op0Cache" {
	//			fmt.Println("a", block.Op0SubmitAt)
	//			if block.Op0SubmitAt.IsZero() {
	//				//已分配
	//				taskList.Op0Assign += 1
	//			} else {
	//				//缓存区
	//				taskList.Op0Cache += 1
	//			}
	//		}
	//		//op1
	//		if block.Op1Stage == "op1Cache" {
	//			//一码 报销单 or 非报销单 - 已分配 or 缓存区
	//			if strings.Index(block.Name, "报销单") == -1 {
	//				//非报销单
	//				//一码 已分配 or 缓存区
	//				if block.Op1SubmitAt.IsZero() {
	//					//已分配
	//					taskList.Op1AssignNotExpenseAccount += 1
	//				} else {
	//					//缓存区
	//					taskList.Op1CacheNotExpenseAccount += 1
	//				}
	//			} else {
	//				//报销单
	//				//一码 已分配 or 缓存区
	//				if block.Op1SubmitAt.IsZero() {
	//					//已分配
	//					taskList.Op1AssignExpenseAccount += 1
	//				} else {
	//					//缓存区
	//					taskList.Op1CacheExpenseAccount += 1
	//				}
	//			}
	//		}
	//		//op2
	//		if block.Op2Stage == "op2Cache" {
	//			//二码 报销单 or 非报销单 - 已分配 or 缓存区
	//			if strings.Index(block.Name, "报销单") == -1 {
	//				//非报销单
	//				//二码 已分配 or 缓存区
	//				if block.Op2SubmitAt.IsZero() {
	//					//已分配
	//					taskList.Op2AssignNotExpenseAccount += 1
	//				} else {
	//					//缓存区
	//					taskList.Op2CacheNotExpenseAccount += 1
	//				}
	//			} else {
	//				//报销单
	//				//二码 已分配 or 缓存区
	//				if block.Op2SubmitAt.IsZero() {
	//					//已分配
	//					taskList.Op2AssignExpenseAccount += 1
	//				} else {
	//					//缓存区
	//					taskList.Op2CacheExpenseAccount += 1
	//				}
	//			}
	//		}
	//		//opq
	//		if block.OpqStage == "opqCache" {
	//			//问题件 已分配 or 缓存区
	//			if block.OpqSubmitAt.IsZero() {
	//				//已分配
	//				taskList.OpqAssign += 1
	//			} else {
	//				//缓存区
	//				taskList.OpqCache += 1
	//			}
	//		}
	//	}
	//}
	return err, taskList, 0
}

func GetUrgencyBillOrPriorityBill(info res.GetVariousStateBill) (err error, list interface{}, total int64) {
	db := global.ProDbMap[info.ProCode+"_task"]
	if db == nil {
		return global.ProDbErr, nil, 0
	}
	var bills []m.ProjectBill
	err = db.Model(&m.ProjectBill{}).Where("stick_level = ? ", info.StickLevel).Find(&bills).Count(&total).Error
	if err != nil {
		return err, nil, 0
	}
	var C res.GetVariousStateBillNum
	for _, v := range bills {
		C.BillNum = append(C.BillNum, v.BillNum)
	}
	return err, C, int64(len(C.BillNum))
}

func SetUrgencyBillOrPriorityBill(info res.SetVariousStateBill) error {
	db := global.ProDbMap[info.ProCode+"_task"]
	if db == nil {
		return global.ProDbErr
	}
	var bills []m.ProjectBill
	for _, num := range info.CaseNumbers {
		err := db.Model(&m.ProjectBill{}).Where("bill_num = ? ", num).Updates(map[string]interface{}{
			"stick_level": info.StickLevel,
		}).Find(&bills).Error
		if err != nil {
			return err
		}
		for _, bill := range bills {
			err = db.Model(&load.ProjectBlock{}).Where("bill_id = ? ", bill.ID).Updates(map[string]interface{}{
				"level": info.StickLevel,
			}).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func SetPriorityOrganizationNumber(info res.SetPriorityOrganizationNumber) error {
	db := global.ProDbMap[info.ProCode+"_task"]
	if db == nil {
		return global.ProDbErr
	}
	var bills []m.ProjectBill
	for _, agencyNum := range info.OrganizationNumber {
		err := db.Model(&m.ProjectBill{}).Where("agency = ? ", agencyNum).Updates(map[string]interface{}{
			"stick_level": info.StickLevel,
		}).Find(&bills).Error
		if err != nil {
			return err
		}
		for _, bill := range bills {
			err = db.Model(&load.ProjectBlock{}).Where("bill_id = ? ", bill.ID).Updates(map[string]interface{}{
				"level": info.StickLevel,
			}).Error
			if err != nil {
				return err
			}
		}
		if info.StickLevel == 0 || info.StickLevel == 99 {
			fmt.Println("123")
			err = utils.DelRedisAgency(info.ProCode)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("345")
			err = utils.SetRedisAgency(info.ProCode, agencyNum)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetCaseDetails(info res.GetCaseDetails) (err error, list interface{}, total int) {
	db := global.ProDbMap[info.ProCode+"_task"]
	if db == nil {
		return global.ProDbErr, nil, 0
	}
	var bills []m.ProjectBill
	var t int64
	if info.SaleChannel != "" {
		db.Where("sale_channel like ?", "%"+info.SaleChannel+"%")
	}
	err = db.Order("id desc").Where("pro_code = ? and bill_num like ?", info.ProCode, "%"+info.BillNum+"%").Find(&bills).Count(&t).Error
	if err != nil {
		return err, nil, 0
	}

	var Cm m.CaseMessage
	var Cd []m.CaseMessage

	for _, bill := range bills {
		//新增录入完成时间
		var infos res.GetCaseDetailsBlock
		infos.Id = bill.ID
		infos.ProCode = info.ProCode
		err, l := GetCaseDetailsBlock(infos)
		if err != nil {
			return err, nil, 0
		}
		var opqSubmitAt []time.Time
		for _, item := range l {
			if item.OpqSubmitAt != "" && item.OpqSubmitAt != "0001-01-01 08:00:00" {
				parse, err := time.Parse("2006-01-02 15:04:05", item.OpqSubmitAt)
				if err != nil {
					return errors.New("转换失败"), nil, 0
				}
				opqSubmitAt = append(opqSubmitAt, parse)
			}
		}
		sort.Slice(opqSubmitAt, func(i, j int) bool {
			return opqSubmitAt[i].After(opqSubmitAt[j])
		})
		if len(opqSubmitAt) > 0 {
			endAt := opqSubmitAt[0].Format("2006-01-02 15:04:05")
			Cm.AppCompleteAt = endAt
		} else if len(opqSubmitAt) == 0 {
			Cm.AppCompleteAt = ""
		}
		//新增录入完成时间
		Cm.BillId = bill.ID
		Cm.ProCode = bill.ProCode
		Cm.BillName = bill.BillName
		Cm.BillNum = bill.BillNum
		Cm.ScanAt = bill.CreatedAt.Format("2006-01-02 15:04:05")
		Cm.Stage = billmap.BillStage[bill.Stage]
		Cm.StickLevel = bill.StickLevel
		Cm.Agency = bill.Agency
		Cm.SaleChannel = bill.SaleChannel
		backAtTheLatest := ""
		timeRemaining := ""
		second := 0.0
		switch info.ProCode {
		case "B0108":
			backAtTheLatest, timeRemaining, second = B0108.CalculateBackTimeAndTimeRemaining(bill)
			if err != nil {
				global.GLog.Error("", zap.Error(err))
			}
		case "B0113":
			zipNameArr := strings.Split(bill.BillName, "_")
			if len(zipNameArr) != 6 {
				global.GLog.Error("案件号压缩包名字有误", zap.Any("", bill.BillName))
			}
			if zipNameArr[len(zipNameArr)-1] == "1" {
				//CSB0113RC0081000
				backAtTheLatest, timeRemaining, second = B0113.CalculateBackTimeAndTimeRemaining(bill)
			} else {
				err, backAtTheLatest, timeRemaining, second = B0118.CalculateBackTimeAndTimeRemaining(bill, t, info.ProCode)
				if err != nil {
					global.GLog.Error("", zap.Error(err))
				}
			}
		default:
			err, backAtTheLatest, timeRemaining, second = B0118.CalculateBackTimeAndTimeRemaining(bill, t, info.ProCode)
			if err != nil {
				global.GLog.Error("", zap.Error(err))
				//return err, Cm, 0
			}
		}

		Cm.LastAt = backAtTheLatest
		Cm.RemainderAt = timeRemaining
		Cm.Second = second
		Cd = append(Cd, Cm)
	}

	return err, Cd, len(Cd)
}

//func GetCaseDetails(info res.GetCaseDetails) (err error, list interface{}, total int) {
//	var Cm m.CaseMessage
//	var Cd []m.CaseMessage
//	db := global.ProDbMap[info.ProCode+"_task"]
//	if db == nil {
//		return global.ProDbErr, nil, 0
//	}
//	var bills []m.ProjectBill
//	err = db.Order("id desc").Where("pro_code = ? ", info.ProCode).Find(&bills).Error
//	if err != nil {
//		return err, nil, 0
//	}
//	for _, bill := range bills {
//		Cm.BillId = bill.ID
//		Cm.ProCode = bill.ProCode
//		Cm.BillName = bill.BillName
//		Cm.BillNum = bill.BillNum
//		Cm.ScanAt = bill.ScanAt.Format("2006-01-02 15:04:05")
//		Cm.Stage = billmap.BillStage[bill.Stage]
//		Cm.LastAt = ""
//		Cm.RemainderAt = ""
//		Cm.StickLevel = bill.StickLevel
//		Cd = append(Cd, Cm)
//	}
//	return err, Cd, len(Cd)
//}

func GetCaseDetailsBlock(info res.GetCaseDetailsBlock) (err error, list []resp.GetCaseDetailsBlockResponse) {
	var blocks []bf.ProjectBlock
	var CaseDetailsBlock []resp.GetCaseDetailsBlockResponse
	db := global.ProDbMap[info.ProCode+"_task"]
	if db == nil {
		return global.ProDbErr, nil
	}
	err = db.Model(&bf.ProjectBlock{}).Order("status").Where("bill_id = ? ", info.Id).Find(&blocks).Error
	if err != nil {
		return err, nil
	}
	for _, block := range blocks {
		var CaseDetailsBlockItem resp.GetCaseDetailsBlockResponse
		CaseDetailsBlockItem.ProCode = info.ProCode
		CaseDetailsBlockItem.BlockId = block.ID
		CaseDetailsBlockItem.BillID = block.BillID
		CaseDetailsBlockItem.Name = block.Name
		CaseDetailsBlockItem.Code = block.Code
		CaseDetailsBlockItem.Stage = ""
		CaseDetailsBlockItem.Status = JudgeBlockStatus(block)

		if block.Op1Code != "0" {
			CaseDetailsBlockItem.Op1Code = block.Op1Code + global.UserCodeName[block.Op1Code]
		} else {
			CaseDetailsBlockItem.Op1Code = block.Op1Code
		}
		CaseDetailsBlockItem.Op1ApplyAt = block.Op1ApplyAt.Format("2006-01-02 15:04:05")
		CaseDetailsBlockItem.Op1SubmitAt = block.Op1SubmitAt.Format("2006-01-02 15:04:05")

		if block.Op2Code != "0" {
			CaseDetailsBlockItem.Op2Code = block.Op2Code + global.UserCodeName[block.Op2Code]
		} else {
			CaseDetailsBlockItem.Op2Code = block.Op2Code
		}
		CaseDetailsBlockItem.Op2ApplyAt = block.Op2ApplyAt.Format("2006-01-02 15:04:05")
		CaseDetailsBlockItem.Op2SubmitAt = block.Op2SubmitAt.Format("2006-01-02 15:04:05")

		if block.OpqCode != "0" {
			CaseDetailsBlockItem.OpqCode = block.OpqCode + global.UserCodeName[block.OpqCode]
		} else {
			CaseDetailsBlockItem.OpqCode = block.OpqCode
		}
		CaseDetailsBlockItem.OpqApplyAt = block.OpqApplyAt.Format("2006-01-02 15:04:05")
		CaseDetailsBlockItem.OpqSubmitAt = block.OpqSubmitAt.Format("2006-01-02 15:04:05")

		if block.Op0Code != "0" {
			CaseDetailsBlockItem.Op0Code = block.Op0Code + global.UserCodeName[block.Op0Code]
		} else {
			CaseDetailsBlockItem.Op0Code = block.Op0Code
		}
		CaseDetailsBlockItem.Op0ApplyAt = block.Op0ApplyAt.Format("2006-01-02 15:04:05")
		CaseDetailsBlockItem.Op0SubmitAt = block.Op0SubmitAt.Format("2006-01-02 15:04:05")
		CaseDetailsBlock = append(CaseDetailsBlock, CaseDetailsBlockItem)
	}
	return err, CaseDetailsBlock
}

func GetCaseDetailsField(info res.GetCaseDetailsField) (err error, list interface{}) {
	var fields []bf.ProjectField
	var CaseDetailsField []resp.GetCaseDetailsFieldResponse
	var CaseDetailsFieldItem resp.GetCaseDetailsFieldResponse
	db := global.ProDbMap[info.ProCode+"_task"]
	if db == nil {
		return global.ProDbErr, nil
	}
	err = db.Model(&bf.ProjectField{}).Order("block_index asc").Order("field_index asc").Where("bill_id = ? AND block_id = ? ", info.BillId, info.BlockId).Find(&fields).Error
	if err != nil {
		return err, nil
	}
	for _, field := range fields {
		CaseDetailsFieldItem.ProCode = info.ProCode
		CaseDetailsFieldItem.BillId = info.BillId
		CaseDetailsFieldItem.BlockId = info.BlockId
		CaseDetailsFieldItem.Name = field.Name
		CaseDetailsFieldItem.Code = field.Code
		CaseDetailsFieldItem.Op0Value = field.Op0Value
		CaseDetailsFieldItem.Op0Input = field.Op0Input
		CaseDetailsFieldItem.Op1Value = field.Op1Value
		CaseDetailsFieldItem.Op1Input = field.Op1Input
		CaseDetailsFieldItem.Op2Value = field.Op2Value
		CaseDetailsFieldItem.Op2Input = field.Op2Input
		CaseDetailsFieldItem.OpqValue = field.OpqValue
		CaseDetailsFieldItem.OpqInput = field.OpqInput
		CaseDetailsField = append(CaseDetailsField, CaseDetailsFieldItem)
	}
	return err, CaseDetailsField
}

// SetUrgencyPriorityBill 将单据紧急优先状态
func SetUrgencyPriorityBill(param res.SetVariousStateBill, taskSuffix string) error {
	db := global.ProDbMap[param.ProCode+taskSuffix]
	if db == nil {
		return global.ProDbErr
	}
	var billIds []string
	tx := db.Begin()
	err := tx.Model(&model2.ProjectBill{}).Where("bill_name in (?) ", param.CaseNumbers).Updates(map[string]interface{}{
		"stick_level": param.StickLevel,
	}).Select("id").Find(&billIds).Error
	if db == nil {
		tx.Rollback()
		return err
	}
	err = db.Model(&load.ProjectBlock{}).Where("bill_id in (?) ", billIds).Updates(map[string]interface{}{
		"level": param.StickLevel,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

type Un struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type UA struct {
	Count int64 `json:"count"`
}

func GetProPermissionPeople(search res.Search) (err error, list interface{}, total int) {
	//limit := search.PageInfo.PageSize
	//offset := search.PageInfo.PageSize * (search.PageInfo.PageIndex - 1)
	m2 := make(map[string]string, 4)
	m2["op0"] = "has_op0"
	m2["op1"] = "has_op1"
	m2["op2"] = "has_op2"
	m2["opq"] = "has_opq"
	//var U []Un
	//query := "SELECT x.code as code, x.nick_name AS name from sys_users as x " +
	//	"where (x.id) in (SELECT user_id from sys_pro_permissions where pro_code = '" + search.ProCode + "' AND " + m2[search.Op] + " = true) " +
	//	"AND x.status = true " +
	//	"limit " + strconv.Itoa(limit) + " " +
	//	"offset " + strconv.Itoa(offset)
	//fmt.Println(query)
	//err = global.GDb.Raw(query).Scan(&U).Error
	//if err != nil {
	//	return err, list, 0
	//}
	//var a UA
	//query2 := "SELECT count(*) from sys_users as x where (x.id) in (SELECT user_id from sys_pro_permissions where pro_code = '" + search.ProCode + "' AND " + m2[search.Op] + " = true) AND x.status = true"
	//err = global.GDb.Raw(query2).Scan(&a).Error
	//if err != nil {
	//	return err, list, 0
	//}
	//for i, _ := range U {
	//	U[i].Name = U[i].Code + U[i].Name
	//}

	var uidArr []string
	err = global.GDb.Model(model.SysProPermission{}).Select("user_id").Where("pro_code = ? and "+m2[search.Op]+" = true", search.ProCode).Find(&uidArr).Error
	if err != nil {
		return err, list, 0
	}
	var userArr []Un
	err = global.GUserDb.Model(model3.SysUser{}).Select("code,name").Where("id in (?)", uidArr).Find(&userArr).Error
	if err != nil {
		return err, list, 0
	}
	for i, _ := range userArr {
		userArr[i].Name = userArr[i].Code + userArr[i].Name
	}
	return err, userArr, len(userArr)
}

func JudgeBlockStatus(block bf.ProjectBlock) string {

	if block.Name == "未定义" {
		if block.Op0Stage == "op0" {
			return "初审待分配"
		}
		if block.Op0Stage == "op0Cache" {
			return "初审录入中"
		}
		if block.Op0Stage == "done" {
			return "录入完毕"
		}
	} else {

		if block.OpqStage == "opq" {
			return "问题件待分配"
		}
		if block.OpqStage == "opqCache" {
			return "问题件录入中"
		}

		if block.Op1Stage == "done" && block.Op2Stage == "done" {
			if block.OpqStage == "done" {
				return "录入完毕"
			} else {
				return "一二码录入完毕"
			}
		}

		if block.Op1Stage == "op1" && block.Op2Stage == "op2" {
			return "一二码待分配"
		} else {
			return "一二码录入中"
		}
	}
	return ""
}
