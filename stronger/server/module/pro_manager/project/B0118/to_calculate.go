package B0118

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"reflect"
	"regexp"
	"server/global"
	pf "server/module/pro_conf/model"
	M "server/module/pro_manager/model"
	pro108 "server/module/pro_manager/project/B0108"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	AgingMsg = "请!滚回去配置对应项目的时效!"
)

// CalculateBackTimeAndTimeRemaining 计算最晚回传时间和剩余时间
func CalculateBackTimeAndTimeRemaining(bill M.ProjectBill, billsCount int64, proCode string) (err error, backAtTheLatest, timeRemaining string, second float64) {
	var pro pf.SysProject
	err = global.GDb.Model(&pf.SysProject{}).Where("code = ? ", proCode).Find(&pro).Error
	if err != nil {
		return err, "", "", 0
	}
	//查时效 优先级: 下载download > 回传upload > 基础base
	var aging []pf.SysProjectConfigAging
	isNotAging := false
	if billsCount > 0 {
		var agingTotal int64
		err = global.GDb.Model(&pf.SysProjectConfigAging{}).Where("pro_id = ? AND config_type = 'base' ", pro.ID).Find(&aging).Count(&agingTotal).Error
		if err != nil {
			return err, "", "", 0
		}
		if agingTotal == 0 {
			isNotAging = true
		}
	}
	if isNotAging {
		return errors.New(AgingMsg), "", "", 0
	}
	//查节假日
	var agingHoliday pf.SysProjectConfigAgingHoliday
	err = global.GDb.Model(&pf.SysProjectConfigAgingHoliday{}).Where("date = ? ", time.Now().Format("20061")).Find(&agingHoliday).Error
	if err != nil {
		return err, "", "", 0
	}

	err, M2 := C(agingHoliday)
	if err != nil {
		return err, "", "", 0
	}

	//计算剩余时间和最晚回传时间
	re, _ := regexp.Compile("[h|m]")

	hms := bill.ScanAt.Format("15:04:05")
	f, in, out, index := BigS(hms, aging)

	//先判断是否为节假日
	num, err := strconv.Atoi(bill.ScanAt.Format("02"))
	if err != nil {
		return errors.New(bill.BillName + ":" + err.Error()), "", "", 0
	}
	fmt.Println("今天是：", num, "号")
	if M2[num] {
		fmt.Println("今天是节假日！")
		//找到下一个工作日
		nextWorkDay := G(M2, num)
		addMouth := 0
		nextWorkDay2 := 0
		if nextWorkDay == 0 {
			addMouth++
		C:
			//查下一个月节假日
			var agingHoliday2 pf.SysProjectConfigAgingHoliday
			err = global.GDb.Model(&pf.SysProjectConfigAgingHoliday{}).Where("date = ? ", bill.ScanAt.AddDate(0, addMouth, 0).Format("20061")).Find(&agingHoliday2).Error
			if err != nil {
				return err, "", "", 0
			}

			fmt.Println("节假日跨月份")

			err, M3 := C(agingHoliday)
			if err != nil {
				return err, "", "", 0
			}
			//找到下一个工作日
			nextWorkDay2 = G(M3, 0)
			if nextWorkDay2 == 0 {
				addMouth++
				goto C
			}
		}

		fmt.Println("下一个工作日为：", nextWorkDay)
		//找出下一个工作日的起始时间
		arr := make([]string, 0)
		for _, v := range aging {
			arr = append(arr, v.AgingStartTime)
		}
		err, s := D(arr)
		fmt.Println("下一个工作日的起始时间为：", s)
		if err != nil {
			return errors.New("单号-" + bill.BillName + ":" + err.Error()), "", "", 0
		}
		local, _ := time.LoadLocation("Local")
		sT, _ := time.ParseInLocation("15:04:05", s, local)
		t, _ := time.ParseDuration("-1s")
		t1, _ := time.ParseDuration(aging[index].RequirementsTime + "m")

		fmt.Println("+考核时效", sT.Add(t).Add(t1).Format("15:04:05"))
		BackAtTheLatest := ""
		if nextWorkDay == 0 {
			fmt.Println("+1天", bill.ScanAt.AddDate(0, addMouth, nextWorkDay2-num).Format("2006-01-02"))
			BackAtTheLatest = bill.ScanAt.AddDate(0, addMouth, nextWorkDay2-num).Format("2006-01-02") + " " + sT.Add(t).Add(t1).Format("15:04:05")
		} else {
			fmt.Println("+1天", bill.ScanAt.AddDate(0, addMouth, nextWorkDay2-num).Format("2006-01-02"))
			BackAtTheLatest = bill.ScanAt.AddDate(0, addMouth, nextWorkDay2-num).Format("2006-01-02") + " " + sT.Add(t).Add(t1).Format("15:04:05")
		}

		backAtTheLatest, _ := time.ParseInLocation("2006-01-02 15:04:05", BackAtTheLatest, time.Local)
		remaining := backAtTheLatest.Sub(time.Now()).Round(time.Second).String()
		fmt.Println("BackAtTheLatest", BackAtTheLatest)
		return nil, BackAtTheLatest, strings.Replace(re.ReplaceAllString(remaining, ":"), "s", "", -1), backAtTheLatest.Sub(time.Now()).Round(time.Second).Seconds()
	}
	//不是节假日
	if !f {
		//没有找到对应的时效
		return nil, "", "", 0
	}

	ag := aging[index]

	if in {
		t, _ := time.ParseDuration(ag.RequirementsTime + "m")
		backAtTheLatest, _ := time.ParseInLocation("2006-01-02 15:04:05", bill.ScanAt.Add(t).Format("2006-01-02 15:04:05"), time.Local)
		remaining := backAtTheLatest.Sub(time.Now()).Round(time.Second).String()
		return nil, bill.ScanAt.Add(t).Format("2006-01-02 15:04:05"), strings.Replace(re.ReplaceAllString(remaining, ":"), "s", "", -1), backAtTheLatest.Sub(time.Now()).Round(time.Second).Seconds()
	}

	//寻找下一个工作日起始时效
	if out {
		//找到下一个工作日
		nextWorkDay := G(M2, num)

		//找出下一个工作日的起始时间
		arr := make([]string, 0)
		for _, v := range aging {
			arr = append(arr, v.AgingStartTime)
		}
		err, s := D(arr)
		if err != nil {
			return errors.New(bill.BillName + ":" + err.Error()), "", "", 0
		}

		local, _ := time.LoadLocation("Local")
		sT, _ := time.ParseInLocation("15:04:05", s, local)
		t, _ := time.ParseDuration("-1s")
		t1, _ := time.ParseDuration(aging[index].RequirementsTime + "m")

		BackAtTheLatest := bill.ScanAt.AddDate(0, 0, nextWorkDay-num).Format("2006-01-02") + " " + sT.Add(t).Add(t1).Format("15:04:05")
		backAtTheLatest, err := time.ParseInLocation("2006-01-02 15:04:05", BackAtTheLatest, time.Local)
		fmt.Println(err)
		remaining := backAtTheLatest.Sub(time.Now()).Round(time.Second).String()
		return nil, BackAtTheLatest, strings.Replace(re.ReplaceAllString(remaining, ":"), "s", "", -1), backAtTheLatest.Sub(time.Now()).Round(time.Second).Seconds()
	}
	return errors.New("没有匹配到对应的时效"), "", "", 0
}

// CalculateBackTimeAndRemainder计算最晚回传时间和剩余时间(时效管理-剩余时间使用)
func CalculateBackTimeAndRemainder(bill M.ProjectBill, billsCount int64, proCode string) (err error, backAtTheLatest, timeRemaining string, second float64) {
	var pro pf.SysProject
	err = global.GDb.Model(&pf.SysProject{}).Where("code = ? ", proCode).Find(&pro).Error
	if err != nil {
		return err, "", "", 0
	}
	//查时效 优先级: 下载download > 回传upload > 基础base
	var aging []pf.SysProjectConfigAging
	isNotAging := false
	if billsCount > 0 {
		var agingTotal int64
		err = global.GDb.Model(&pf.SysProjectConfigAging{}).Where("pro_id = ? AND config_type = 'base' ", pro.ID).Find(&aging).Count(&agingTotal).Error
		if err != nil {
			return err, "", "", 0
		}
		if agingTotal == 0 {
			isNotAging = true
		}
	}
	if isNotAging {
		return errors.New(AgingMsg), "", "", 0
	}
	//查节假日
	var agingHoliday pf.SysProjectConfigAgingHoliday
	err = global.GDb.Model(&pf.SysProjectConfigAgingHoliday{}).Where("date = ? ", time.Now().Format("20061")).Find(&agingHoliday).Error
	if err != nil {
		return err, "", "", 0
	}

	err, M2 := C(agingHoliday)
	if err != nil {
		return err, "", "", 0
	}

	//计算剩余时间和最晚回传时间
	re, _ := regexp.Compile("[h|m]")

	hms := bill.ScanAt.Format("15:04:05")
	f, in, out, index := BigS(hms, aging)

	//先判断是否为节假日
	num, err := strconv.Atoi(bill.ScanAt.Format("02"))
	if err != nil {
		return errors.New(bill.BillName + ":" + err.Error()), "", "", 0
	}
	fmt.Println("今天是：", num, "号")
	if M2[num] {
		fmt.Println("今天是节假日！")
		//找到下一个工作日
		nextWorkDay := G(M2, num)
		addMouth := 0
		nextWorkDay2 := 0
		if nextWorkDay == 0 {
			addMouth++
		C:
			//查下一个月节假日
			var agingHoliday2 pf.SysProjectConfigAgingHoliday
			err = global.GDb.Model(&pf.SysProjectConfigAgingHoliday{}).Where("date = ? ", bill.ScanAt.AddDate(0, addMouth, 0).Format("20061")).Find(&agingHoliday2).Error
			if err != nil {
				return err, "", "", 0
			}

			fmt.Println("节假日跨月份")

			err, M3 := C(agingHoliday)
			if err != nil {
				return err, "", "", 0
			}
			//找到下一个工作日
			nextWorkDay2 = G(M3, 0)
			if nextWorkDay2 == 0 {
				addMouth++
				goto C
			}
		}

		fmt.Println("下一个工作日为：", nextWorkDay)
		//找出下一个工作日的起始时间
		arr := make([]string, 0)
		for _, v := range aging {
			arr = append(arr, v.AgingStartTime)
		}
		err, s := D(arr)
		fmt.Println("下一个工作日的起始时间为：", s)
		if err != nil {
			return errors.New("单号-" + bill.BillName + ":" + err.Error()), "", "", 0
		}
		local, _ := time.LoadLocation("Local")
		sT, _ := time.ParseInLocation("15:04:05", s, local)
		t, _ := time.ParseDuration("-1s")
		t1, _ := time.ParseDuration(aging[index].RequirementsTime + "m")

		fmt.Println("+考核时效", sT.Add(t).Add(t1).Format("15:04:05"))
		BackAtTheLatest := ""
		if nextWorkDay == 0 {
			fmt.Println("+1天", bill.ScanAt.AddDate(0, addMouth, nextWorkDay2-num).Format("2006-01-02"))
			BackAtTheLatest = bill.ScanAt.AddDate(0, addMouth, nextWorkDay2-num).Format("2006-01-02") + " " + sT.Add(t).Add(t1).Format("15:04:05")
		} else {
			fmt.Println("+1天", bill.ScanAt.AddDate(0, addMouth, nextWorkDay2-num).Format("2006-01-02"))
			BackAtTheLatest = bill.ScanAt.AddDate(0, addMouth, nextWorkDay2-num).Format("2006-01-02") + " " + sT.Add(t).Add(t1).Format("15:04:05")
		}
		backAtTheLatest, _ := time.ParseInLocation("2006-01-02 15:04:05", BackAtTheLatest, time.Local)
		remaining := backAtTheLatest.Sub(time.Now()).Round(time.Second).String()
		replace := strings.Replace(re.ReplaceAllString(remaining, ":"), "s", "", -1)
		timeRemaining = pro108.GetRemainder(replace)
		return nil, BackAtTheLatest, timeRemaining, backAtTheLatest.Sub(time.Now()).Round(time.Second).Seconds()
	}
	//不是节假日
	if !f {
		//没有找到对应的时效
		return nil, "", "", 0
	}

	ag := aging[index]

	if in {
		t, _ := time.ParseDuration(ag.RequirementsTime + "m")
		backAtTheLatest, _ := time.ParseInLocation("2006-01-02 15:04:05", bill.ScanAt.Add(t).Format("2006-01-02 15:04:05"), time.Local)
		remaining := backAtTheLatest.Sub(time.Now()).Round(time.Second).String()
		replace := strings.Replace(re.ReplaceAllString(remaining, ":"), "s", "", -1)
		timeRemaining = pro108.GetRemainder(replace)
		return nil, bill.ScanAt.Add(t).Format("2006-01-02 15:04:05"), timeRemaining, backAtTheLatest.Sub(time.Now()).Round(time.Second).Seconds()
	}
	//寻找下一个工作日起始时效
	if out {
		//找到下一个工作日
		nextWorkDay := G(M2, num)

		//找出下一个工作日的起始时间
		arr := make([]string, 0)
		for _, v := range aging {
			arr = append(arr, v.AgingStartTime)
		}
		err, s := D(arr)
		if err != nil {
			return errors.New(bill.BillName + ":" + err.Error()), "", "", 0
		}

		local, _ := time.LoadLocation("Local")
		sT, _ := time.ParseInLocation("15:04:05", s, local)
		t, _ := time.ParseDuration("-1s")
		t1, _ := time.ParseDuration(aging[index].RequirementsTime + "m")

		BackAtTheLatest := bill.ScanAt.AddDate(0, 0, nextWorkDay-num).Format("2006-01-02") + " " + sT.Add(t).Add(t1).Format("15:04:05")
		backAtTheLatest, err := time.ParseInLocation("2006-01-02 15:04:05", BackAtTheLatest, time.Local)
		fmt.Println(err)
		remaining := backAtTheLatest.Sub(time.Now()).Round(time.Second).String()
		replace := strings.Replace(re.ReplaceAllString(remaining, ":"), "s", "", -1)
		timeRemaining = pro108.GetRemainder(replace)
		return nil, BackAtTheLatest, timeRemaining, backAtTheLatest.Sub(time.Now()).Round(time.Second).Seconds()
	}
	return errors.New("没有匹配到对应的时效"), "", "", 0
}

func Calculate(start, mid, end string) (startT, midT, endT time.Time, err error) {
	local, err := time.LoadLocation("Local")
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, err
	}

	if start != "" {
		startT, err = time.ParseInLocation("15:04:05", start, local)
		if err != nil {
			return time.Time{}, time.Time{}, time.Time{}, err
		}
	}

	if mid != "" {
		midT, err = time.ParseInLocation("15:04:05", mid, local)
		if err != nil {
			return time.Time{}, time.Time{}, time.Time{}, err
		}
	}

	if end != "" {
		endT, err = time.ParseInLocation("15:04:05", end, local)
		if err != nil {
			return time.Time{}, time.Time{}, time.Time{}, err
		}
	}

	return startT, midT, endT, err
}

// C 将节假日数据数据传换成map
func C(agingHoliday pf.SysProjectConfigAgingHoliday) (error, map[int]bool) {
	s := reflect.TypeOf(agingHoliday)
	sc := reflect.ValueOf(agingHoliday)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	M2 := make(map[int]bool, s.NumField()+1)
	for i := 0; i < s.NumField(); i++ {
		fieldName := s.Field(i).Name
		if strings.Index(fieldName, "Day") != -1 {
			num, err := strconv.Atoi(strings.Split(fieldName, "Day")[1])
			if err != nil {
				return err, nil
			}
			M2[num] = sc.FieldByName(fieldName).Interface().(bool)
		}
	}
	return nil, M2
}

// BigS 计算时间在哪个范围
func BigS(hms string, aging []pf.SysProjectConfigAging) (f, in, out bool, index int) {
	for i, v := range aging {
		//计算在哪个时间范围
		//例如
		//时效内：8 - 15 时效外：18 - 24, 0 - 7
		//时效内：15 - 17 时效外：18 - 24, 0 - 7
		err, hasAging, O, I := S(hms, v.AgingStartTime, v.AgingEndTime, v.AgingOut)
		if err != nil {
			return
		}
		if hasAging {
			index = i
			f = true
			if O {
				out = true
			}
			if I {
				in = true
			}
			break
		}
	}
	return
}

// S 计算时间在哪个范围
func S(c, i, j string, k []string) (err error, h, out, in bool) {
	//i c j
	iT, cT, jT, err := Calculate(i, c, j)
	if err != nil {
		return err, h, false, false
	}
	if cT.After(iT) && cT.Before(jT) {
		h = true
		return nil, h, false, true
	}
	//c k
	for _, v := range k {
		arr := strings.Split(v, "-")
		kT, cT, lT, err := Calculate(arr[0], c, arr[1])
		if err != nil {
			return err, h, false, false
		}
		if cT.After(kT) && cT.Before(lT) {
			h = true
			break
		}
	}
	return nil, h, true, false
}

// D 找出下一个工作日的起始时间
func D(arr []string) (error, string) {
	if len(arr) == 0 {
		return errors.New("请!滚回去配置对应项目的时效! "), ""
	}
	for i, v := range arr {
		if i == 0 {
			continue
		}
		before, after, _, err := Calculate(arr[i-1], v, "")
		if err != nil {
			return err, ""
		}
		if after.Before(before) {
			arr[i-1], arr[i] = arr[i], arr[i-1]
		}
	}
	return nil, arr[0]
}

// G 由于生成M2时的无序写入, 因此该方法为寻找节假日的下一个工作日
func G(M2 map[int]bool, holiday int) int {
	nextWorkDay := 32
	for i, v := range M2 {
		if i < holiday {
			continue
		}
		if nextWorkDay == 32 && i != holiday && !v {
			nextWorkDay = i
		}
		if !v && i < nextWorkDay && nextWorkDay != 32 && i != holiday {
			nextWorkDay = i
		}
	}

	//if nextWorkDay == 32 {
	//	nextWorkDay = 0
	//}

	return nextWorkDay
}

// CalculateRequirementOfAging 获取回传时效要求
func CalculateRequirementOfAging(bill M.ProjectBill) string {
	var pro pf.SysProject
	err := global.GDb.Model(&pf.SysProject{}).Where("code = ? ", bill.ProCode).Find(&pro).Error
	if err != nil {
		global.GLog.Error("CalculateRequirementOfAging : "+global.ProDbErr.Error(), zap.Error(errors.New("CalculateRequirementOfAging : +"+global.ProDbErr.Error())))
		return ""
	}
	var agingTotal int64
	var aging []pf.SysProjectConfigAging
	err = global.GDb.Model(&pf.SysProjectConfigAging{}).Where("pro_id = ? AND config_type = 'base' ", pro.ID).Find(&aging).Count(&agingTotal).Error
	if err != nil {
		global.GLog.Error("CalculateRequirementOfAging : "+global.ProDbErr.Error(), zap.Error(errors.New("CalculateRequirementOfAging : +"+global.ProDbErr.Error())))
		return ""
	}
	if len(aging) == 0 {
		return ""
	}
	hms := bill.ScanAt.Format("15:04:05")
	_, _, _, index := BigS(hms, aging)
	return aging[index].RequirementsTime
}

// CalculateProjectReturnTime 计算项目的最晚回传时间 B0118.
// 最晚回传时间:除了工作日时效内时间要 加 考核要求的时间  其它全部直接取 时效外最晚时间
// 时效内工作日最晚回传时间 = 扫描时间 + 考核要求的时间
// 时效外工作日内最晚回传时间 = 下一个工作日的时效外最晚时间
// 时效外节假日最晚回传时间 = 下一个工作日的时效外最晚时间
func CalculateProjectReturnTime(bill M.ProjectBill, proCode string) (err error, backAtTheLatest string) {
	// 查看当前时间  是不是节假日
	//查节假日
	var agingHoliday pf.SysProjectConfigAgingHoliday
	err = global.GDb.Model(&pf.SysProjectConfigAgingHoliday{}).Where("date = ? ", bill.CreatedAt.Format("20061")).Order("created_at desc").Find(&agingHoliday).Error
	if err != nil {
		return err, ""
	}
	// 节假日转换成map
	err, M2 := C(agingHoliday)
	//今天日期
	num, err := strconv.Atoi(bill.CreatedAt.Format("02"))
	//是工作日 肯定是时效内
	// 获取合同时效
	db := global.GDb.Model(&pf.SysProjectConfigAgingContract{}).Order("created_at ASC")
	var configs pf.SysProjectConfigAgingContract
	var contractCount int64
	err = global.GDb.Model(&pf.SysProjectConfigAgingContract{}).Where("code = ? And claim_type = ?", proCode, bill.ClaimType).Order("created_at ASC").Count(&contractCount).Error
	if contractCount == 0 {
		err = db.Where("code = ? And claim_type = 8", proCode, bill.ClaimType).Find(&configs).Error
	} else {
		err = db.Where("code = ? And claim_type = ?", proCode, bill.ClaimType).Find(&configs).Error
	}
	billScanAt := bill.CreatedAt.Format("15:04:05")     //扫描时间
	err, ist, count := AgeingInner(configs, billScanAt) //判断是时效内还是时效外
	if err != nil {
		return err, ""
	}
	if !M2[num] {
		// 开始计算最晚回传时间
		// 查看是否在时效内 工作日
		//1.在时效内工作日
		if ist && count == 1 {
			requirementsTime := configs.RequirementsTime + "m" //考核要求时间
			duration, err := time.ParseDuration(requirementsTime)
			if err != nil {
				return err, ""
			}
			latestTime := bill.CreatedAt.Add(duration) // 时效内最晚回传时间
			backAtTheLatest = latestTime.Format("2006-01-02 15:04:05")
			return err, backAtTheLatest
		} else if ist && count == 2 {
			err, backAtTheLatest = AgeingOutTimeDefault(M2, num, bill.CreatedAt, configs, bill)

			return err, backAtTheLatest
		}
	} else if M2[num] {
		//节假日 - 时效外
		err, backAtTheLatest = AgeingOutTimeDefault(M2, num, bill.CreatedAt, configs, bill)
		if err != nil {
			return err, ""
		}
		return err, backAtTheLatest
	}
	return err, ""
}

// 判断时间是否在时效内 ，工作日
// 1 = 工作日 时效内
// 2 = 工作日 非时效内
// 3 = 报错 || 不是时效内或不是时效外
func AgeingInner(configs pf.SysProjectConfigAgingContract, billScanAt string) (err error, ist bool, count int64) {
	if configs.ContractStartTime == "" || configs.ContractEndTime == "" || billScanAt == "" {
		return errors.New("请检查时效开始，结束时间，扫描时间是否正确"), false, 3
	}
	startTime, _ := time.Parse("15:04:05", configs.ContractStartTime)
	endTime, _ := time.Parse("15:04:05", configs.ContractEndTime)
	scanTime, _ := time.Parse("15:04:05", billScanAt)
	//时效内的案件
	if scanTime.After(startTime) && scanTime.Before(endTime) {
		return err, true, 1
	}
	//非工作日时效内
	return err, true, 2
}

// 获取最晚回传时间 时效外 - 直接取 时效外最晚时间
func AgeingOutTimeDefault(M2 map[int]bool, num int, scanAt time.Time, configs pf.SysProjectConfigAgingContract, bill M.ProjectBill) (err error, backAtTheLatest string) {
	//找到下一个工作日
	nextWorkDay := G(M2, num)
	bigNext := G(M2, nextWorkDay)
	_, isExit := M2[bigNext]
	inYearDay := bill.CreatedAt.Year()
	inMonthDay := bill.CreatedAt.Month()
	//当月的天数
	inDay := time.Date(inYearDay, inMonthDay+1, 0, 0, 0, 0, 0, time.UTC).Day()
	//fmt.Printf("当前月份：%d %d 共有 %d 天\n", inYearDay, int64(inMonthDay), inDay)
	var newYear string
	var newMonth string
	var newDay string
	//是工作日 且 大于数据库最大天数 或 下个工作日 大于这个月天数 如 2023 年 8月
	// NextJobDay(32) > DBMax(31) || NextJobDay(32) > ThisMonth(30)
	if !isExit && nextWorkDay > len(M2) || nextWorkDay > inDay {
		//获取下个月第一个工作日
		//查节假日
		//获取下个月
		nextMonth := bill.CreatedAt.Format("20061")
		//下个月
		var useMonth string
		if len(nextMonth) == 5 {
			lastStr := nextMonth[len(nextMonth)-1]
			atoi, _ := strconv.Atoi(string(lastStr))
			atoi = atoi + 1
			if atoi > 12 {
				atoi = 1
			}
			useMonth = strconv.Itoa(atoi)
		} else if len(nextMonth) == 6 {
			lastStr := nextMonth[len(nextMonth)-2:]
			atoi, _ := strconv.Atoi(string(lastStr))
			atoi = atoi + 1
			if atoi > 12 {
				atoi = 1
			}
			useMonth = strconv.Itoa(atoi)
		}
		//记录下来第一个工作日
		dayJob := 0
		year := nextMonth[:4] + useMonth
		//年 月 日
		if len(useMonth) == 1 {
			useMonth = "0" + useMonth
		}
		newYear = nextMonth[:4]
		newMonth = useMonth

		var agingHoliday pf.SysProjectConfigAgingHoliday
		err = global.GDb.Model(&pf.SysProjectConfigAgingHoliday{}).Where("date = ? ", year).Find(&agingHoliday).Error
		if err != nil {
			return err, ""
		}
		// 节假日转换成map
		_, HolidayMap := C(agingHoliday)

		//格式化年月日
		//给map排序
		arrays := []int{}
		for key := range HolidayMap {
			arrays = append(arrays, key)
		}
		sort.Ints(arrays)
		for _, k := range arrays {
			if !HolidayMap[k] {
				if dayJob != 0 {
					continue
				}
				dayJob = k
				continue
			}
		}
		newDay = strconv.Itoa(dayJob)
		if len(newDay) == 1 {
			newDay = "0" + newDay
		}
		// 扫描时间年月 + nextWorkDay  + 时效外开始时间
		backAtTheLatest = newYear + "-" + newMonth + "-" + newDay + " " + configs.ContractOutsideEndTime
		return err, backAtTheLatest
	}

	// 扫描时间年月 + nextWorkDay  + 时效外开始时间  nextMonth[:4] + useMonth
	itoaYear := strconv.Itoa(scanAt.Year())
	itoaMonth := strconv.Itoa(int(scanAt.Month()))
	itoaNextWorkDay := strconv.Itoa(nextWorkDay)
	if len(itoaMonth) == 1 {
		itoaMonth = "0" + itoaMonth
	}
	if len(itoaNextWorkDay) == 1 {
		itoaNextWorkDay = "0" + itoaNextWorkDay
	}
	backAtTheLatest = itoaYear + "-" + itoaMonth + "-" + itoaNextWorkDay + " " + configs.ContractOutsideEndTime
	return err, backAtTheLatest
}
