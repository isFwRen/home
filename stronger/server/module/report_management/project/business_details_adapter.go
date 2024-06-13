package project

import (
	"errors"
	m "server/module/pro_manager/model"
	"server/module/pro_manager/project/B0108"
	"server/module/pro_manager/project/B0118"
)

//FindAndCalculateBackAtTheLatestAdapter 获取最晚回传时间函数
func FindAndCalculateBackAtTheLatestAdapter(bill m.ProjectBill, proCode string) (err error, backAtTheLatest string) {
	switch proCode {
	case "B0118":
		err, backAtTheLatest, _, _ = B0118.CalculateBackTimeAndTimeRemaining(bill, 1, proCode)
		if err != nil {
			return err, backAtTheLatest
		}
		return err, backAtTheLatest
	case "B0108":
		backAtTheLatest, _, _ = B0108.CalculateBackTimeAndTimeRemaining(bill)
		return nil, backAtTheLatest
	default:
		return errors.New("该项目没有自定义获取最晚回传时间函数函数"), backAtTheLatest
	}
}

//CalculateRequirementOfAgingAdapter 获取回传时效要求
func CalculateRequirementOfAgingAdapter(bill m.ProjectBill, proCode string) (err error, requirementOfAging string) {
	switch proCode {
	case "B0118":
		requirementOfAging = B0118.CalculateRequirementOfAging(bill)
		return err, requirementOfAging
	case "B0108":
		requirementOfAging = B0108.CalculateRequirementOfAging(bill)
		return err, requirementOfAging
	default:
		return errors.New("该项目没有自定义获取回传时效要求函数"), requirementOfAging
	}
}

// 计算最晚回传时间函数
func CalculateAllReturnTime(bill m.ProjectBill, proCode string) (err error, backAtTheLatest string) {
	err, backAtTheLatest = B0118.CalculateProjectReturnTime(bill, proCode)
	if err != nil {
		return err, ""
	}
	return err, backAtTheLatest
	//switch proCode {
	//case "B0118":
	//	fmt.Println("======CalculateAllReturnTime=========bill.CreatedAt===============", bill.CreatedAt)
	//	err, backAtTheLatest = B0118.CalculateProjectReturnTime(bill, proCode)
	//	fmt.Println(backAtTheLatest)
	//	if err != nil {
	//		return err, backAtTheLatest
	//	}
	//	return err, backAtTheLatest
	//default:
	//	return errors.New("该项目没有自定义获取最晚回传时间函数函数"), backAtTheLatest
	//}
}
