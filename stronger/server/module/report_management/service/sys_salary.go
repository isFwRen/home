package service

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"path"
	"server/global"
	"server/module/report_management/model"
	"server/module/report_management/model/request"
	modelBase "server/module/sys_base/model"
	"strconv"
	"strings"
	"time"
)

func GetPtSalaryTask(info request.SysSalarySearch, uid string) (err error, list interface{}, total int64) {
	var U modelBase.SysUser
	err = global.GDb.Model(&modelBase.SysUser{}).Where("id = ? ", uid).Find(&U).Error
	if err != nil {
		return err, nil, 0
	}
	Start, _ := time.ParseInLocation("2006-01-02 15:04:05", info.Start+"-02 00:00:00", time.Local)
	End, _ := time.ParseInLocation("2006-01-02 15:04:05", info.End+"-02 00:00:00", time.Local)
	var PtSalary []model.SysSalary
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	db := global.GDb.Model(&model.SysSalary{})
	db = db.Where("pay_day >= ? AND pay_day < ? AND code = ? ", Start, End.AddDate(0, 1, 0), U.Code)
	err = db.Count(&total).Error
	if total > 0 {
		err = db.Order("id desc").Limit(limit).Offset(offset).Find(&PtSalary).Error
		return err, PtSalary, total
	}
	return err, nil, 0
}

func GetPtSalary(info request.SysSalarySearch) (err error, list interface{}, total int64) {
	var PtSalary []model.SysSalary
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	db := global.GDb.Model(&model.SysSalary{})
	payDay, _ := time.ParseInLocation("2006-01-02 15:04:05", info.Ym[:4]+"-"+info.Ym[4:6]+"-02 00:00:00", time.Local)
	db = db.Where("pay_day = ? ", payDay)
	if info.Code != "" {
		db = db.Where("code = ? ", info.Code)
	}
	if info.Name != "" {
		db = db.Where("nick_name = ? ", info.Name)
	}
	err = db.Count(&total).Error
	if total > 0 {
		err = db.Order("id desc").Limit(limit).Offset(offset).Find(&PtSalary).Error
		return err, PtSalary, total
	}
	return err, nil, 0
}

func GetInternalSalary(info request.SysSalarySearch, uid string) (err error, list interface{}, total int64) {
	var U modelBase.SysUser
	err = global.GDb.Model(&modelBase.SysUser{}).Where("id = ? ", uid).Find(&U).Error
	if err != nil {
		return err, nil, 0
	}
	var Role modelBase.SysRoles
	err = global.GDb.Model(&modelBase.SysRoles{}).Where("id = ? ", U.RoleId).Find(&Role).Error
	if err != nil {
		return err, nil, 0
	}
	fmt.Println("role", Role.Name)
	var InternalSalary []model.SysInternalSalary
	var SysInternalSalaryVersionTwo []model.SysInternalSalaryVersionTwo

	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	db := global.GDb.Model(&model.SysInternalSalary{})
	payDay, _ := time.ParseInLocation("2006-01-02 15:04:05", info.Ym[:4]+"-"+info.Ym[4:6]+"-02 00:00:00", time.Local)
	db = db.Where("pay_day = ? ", payDay)

	if Role.Name != "操作员" && Role.Name != "调度员" && Role.Name != "质检员" {
		if info.Code != "" {
			db = db.Where("code like ? ", "%"+info.Code+"%")
		}
		if info.Name != "" {
			db = db.Where("nick_name like ? ", "%"+info.Name+"%")
		}
	} else {
		db = db.Where("code = ? ", U.Code)
	}
	err = db.Count(&total).Error
	if total > 0 {
		err = db.Order("id desc").Limit(limit).Offset(offset).Find(&InternalSalary).Error
		for _, v := range InternalSalary {
			var item model.SysInternalSalaryVersionTwo
			item.Code = v.Code
			item.NickName = v.NickName
			item.DateOfEntry = v.DateOfEntry
			item.DateOfMountGuard = v.DateOfMountGuard
			item.ProjectGroup = v.ProjectGroup
			item.Date = v.Date
			item.YearsOfEmployment = v.YearsOfEmployment
			item.IsStayOutsideOverNight = v.IsStayOutsideOverNight
			item.BasicWage = v.BasicWage
			item.PerformancePay = v.PerformancePay
			item.WageJobs = v.WageJobs
			item.OvertimePay = v.OvertimePay
			item.OvertimeWorkLoad = v.OvertimeWorkLoad
			item.WorkOvertimeTogether = v.WorkOvertimeTogether
			//出勤
			item.StandardNumberOfDays = v.StandardNumberOfDays
			item.NumberOfDaysOnTheJob = v.NumberOfDaysOnTheJob
			item.BereavementLeave = v.BereavementLeave
			item.Late = v.Late
			item.LeaveEarly = v.LeaveEarly
			item.PersonalLeave = v.PersonalLeave
			item.Absenteeism = v.Absenteeism
			item.AnnualHoliday = v.AnnualHoliday
			item.ResignationOrNewJob = v.ResignationOrNewJob
			item.SickLeave = v.SickLeave
			//福利
			item.AnnualWork = v.AnnualWork
			item.PerfectAttendanceAward = v.PerfectAttendanceAward
			item.DispatchOrReserve = v.DispatchOrReserve
			item.StayOutsideOverNight = v.StayOutsideOverNight
			item.OtherOfWelfare = v.OtherOfWelfare
			//扣除
			item.Attendance = v.Attendance
			item.SocialSecurity = v.SocialSecurity
			item.Quality = v.Quality
			item.InsuranceDispatch = v.InsuranceDispatch
			item.CodeOfConduct = v.CodeOfConduct
			item.OtherOfDeduct = v.OtherOfDeduct
			//
			item.PieceRateWages = v.PieceRateWages
			item.MinimumGuarantee = v.MinimumGuarantee
			item.RealWages = v.RealWages
			item.WithholdingTax = v.WithholdingTax
			item.TaxPay = v.TaxPay
			item.CompanySupplement = v.CompanySupplement
			item.Remark = v.Remark
			//payday
			item.PayDay = v.PayDay.Format("2006-01")
			//业务量合计、外部质量、内部质量、折算比例、工作量由于其的不确定内容，因此做成数组形式存进数据库
			//数组格式 : [ 标题,标题所对应的内容, 标题,标题所对应的内容.......]
			//取出来的时候, 把他们转换成map形式返回给前端, 以便取值
			//业务量合计
			TotalBusinessVolume := make(map[string]string, 0)
			ExternalQuality := make(map[string]string, 0)
			InternalQuality := make(map[string]string, 0)
			WorkLoad := make(map[string]string, 0)
			ReducedProportion := make(map[string]string, 0)
			for k, l := range v.TotalBusinessVolume {
				if k == 0 {
					TotalBusinessVolume[l] = ""
				}
				if k%2 != 0 {
					TotalBusinessVolume[v.TotalBusinessVolume[k-1]] = l
				}
				if k%2 == 0 {
					TotalBusinessVolume[l] = ""
				}
			}
			//外部质量
			for k, l := range v.ExternalQuality {
				if k == 0 {
					ExternalQuality[l] = ""
				}
				if k%2 != 0 {
					ExternalQuality[v.ExternalQuality[k-1]] = l
				}
				if k%2 == 0 {
					ExternalQuality[l] = ""
				}
			}
			//内部质量
			for k, l := range v.InternalQuality {
				if k == 0 {
					InternalQuality[l] = ""
				}
				if k%2 != 0 {
					InternalQuality[v.InternalQuality[k-1]] = l
				}
				if k%2 == 0 {
					InternalQuality[l] = ""
				}
			}
			//工作量
			for k, l := range v.WorkLoad {
				if k == 0 {
					WorkLoad[l] = ""
				}
				if k%2 != 0 {
					WorkLoad[v.WorkLoad[k-1]] = l
				}
				if k%2 == 0 {
					WorkLoad[l] = ""
				}
			}
			//折算比例
			for k, l := range v.ReducedProportion {
				if k == 0 {
					ReducedProportion[l] = ""
				}
				if k%2 != 0 {
					ReducedProportion[v.ReducedProportion[k-1]] = l
				}
				if k%2 == 0 {
					ReducedProportion[l] = ""
				}
			}
			item.TotalBusinessVolume = TotalBusinessVolume
			item.ExternalQuality = ExternalQuality
			item.InternalQuality = InternalQuality
			item.WorkLoad = WorkLoad
			item.ReducedProportion = ReducedProportion
			SysInternalSalaryVersionTwo = append(SysInternalSalaryVersionTwo, item)
		}
		return err, SysInternalSalaryVersionTwo, total
	}
	return err, nil, 0
}

func PtSalaryFileUpload(filename []string) (errMsg string) {
	fmt.Println("PtSalaryFileUpload")
	for _, f := range filename {
		//查询工资数据excel是否重复保存
		payDay, _ := time.ParseInLocation("2006-01-02 15:04:05", f[:4]+"-"+f[4:6]+"-02 00:00:00", time.Local)
		var total int64
		err := global.GDb.Model(&model.SysSalary{}).Where("pay_day = ? ", payDay).Count(&total).Error
		if err != nil {
			return err.Error()
		}
		if total != 0 {
			//删除旧数据
			err = global.GDb.Model(&model.SysSalary{}).Where("pay_day = ? ", payDay).Delete(&model.SysSalary{}).Error
			if err != nil {
				return err.Error()
			}
		}
		var Salaries []model.SysSalary
		var item model.SysSalary
		// 线上 uploads/file
		basicPath := global.GConfig.LocalUpload.FilePath + "/Salary/" + "pt/" + f[:4] + "/"
		// 本地
		//basicPath := "./Salary/" + f[:4] + "/"
		dst := path.Join(basicPath, f)
		xlsx, err := excelize.OpenFile(dst)
		if err != nil {
			errMsg += err.Error()
			return errMsg
		}
		// 获取excel中具体的列的值
		rows := xlsx.GetRows("Sheet" + "1")
		for i, row := range rows {
			if i != 0 {
				item.Code = row[0]
				item.NickName = row[1]
				item.ProductionStatistics = row[2]
				item.ProductionSalary = row[3]
				item.CustomerComplaints = row[4]
				item.ReferralBonus = row[5]
				item.TotalWages = row[6]
				item.Tax = row[7]
				item.EventuallyPay = row[8]
				item.PayDay = payDay
				Salaries = append(Salaries, item)
			}
		}
		err = global.GDb.Model(&model.SysSalary{}).Create(&Salaries).Error
		if err != nil {
			errMsg += err.Error()
			return errMsg
		}

	}
	return errMsg
}

func InternalSalaryFileUpload(filename []string) (errMsg string) {
	fmt.Println("InternalSalaryFileUpload")
	for _, f := range filename {
		//查询工资数据excel是否重复保存
		fmt.Println(f[:4] + "-" + f[4:6] + "-01 00:00:00")
		//选择每个月的2号是为了防止时区不同
		payDay, _ := time.ParseInLocation("2006-01-02 15:04:05", f[:4]+"-"+f[4:6]+"-02 00:00:00", time.Local)
		var total int64
		err := global.GDb.Model(&model.SysInternalSalary{}).Where("pay_day = ? ", payDay).Count(&total).Error
		if err != nil {
			return err.Error()
		}
		if total != 0 {
			//删除旧数据
			err = global.GDb.Model(&model.SysInternalSalary{}).Where("pay_day = ? ", payDay).Delete(&model.SysInternalSalary{}).Error
			if err != nil {
				return err.Error()
			}
		}
		var InternalSalarySalaries []model.SysInternalSalary

		// 线上 uploads/file
		basicPath := global.GConfig.LocalUpload.FilePath + "/Salary/" + "internal/" + f[:4] + "/"
		// 本地
		//basicPath := "./Salary/" + f[:4] + "/"
		dst := path.Join(basicPath, f)
		xlsx, err := excelize.OpenFile(dst)
		if err != nil {
			errMsg += err.Error()
			return errMsg
		}

		// 获取excel中具体的列的值
		rows := xlsx.GetRows("Sheet" + "1")
		mergeCells := xlsx.GetMergeCells("Sheet" + "1")

		TotalBusinessVolumeArr := make([]string, 0)
		ExternalQualityArr := make([]string, 0)
		InternalQualityArr := make([]string, 0)
		WorkLoadArr := make([]string, 0)
		ReducedProportionArr := make([]string, 0)

		TotalBusinessVolumeRange := make(map[string][]int, 0)
		ExternalQualityRange := make(map[string][]int, 0)
		InternalQualityRange := make(map[string][]int, 0)
		AttendanceRange := make(map[string][]int, 0)
		WelfareRange := make(map[string][]int, 0)
		DeductRange := make(map[string][]int, 0)
		WorkLoadRange := make(map[string][]int, 0)
		ReducedProportionRange := make(map[string][]int, 0)
		for _, v := range mergeCells {

			arr := strings.Split(v[0], ":")

			if v[1] == "业务量合计" {
				TotalBusinessVolumeRange["业务量合计"] = append(TotalBusinessVolumeRange["业务量合计"], ChangeAxisToIndex(arr[0]), ChangeAxisToIndex(arr[1]))
			}
			if v[1] == "外部质量" {
				ExternalQualityRange["外部质量"] = append(ExternalQualityRange["外部质量"], ChangeAxisToIndex(arr[0]), ChangeAxisToIndex(arr[1]))
			}
			if v[1] == "内部质量" {
				InternalQualityRange["内部质量"] = append(InternalQualityRange["内部质量"], ChangeAxisToIndex(arr[0]), ChangeAxisToIndex(arr[1]))
			}
			if v[1] == "出勤" {
				AttendanceRange["出勤"] = append(AttendanceRange["出勤"], ChangeAxisToIndex(arr[0]), ChangeAxisToIndex(arr[1]))
			}
			if v[1] == "福利" {
				WelfareRange["福利"] = append(WelfareRange["福利"], ChangeAxisToIndex(arr[0]), ChangeAxisToIndex(arr[1]))
			}
			if v[1] == "扣除" {
				DeductRange["扣除"] = append(DeductRange["扣除"], ChangeAxisToIndex(arr[0]), ChangeAxisToIndex(arr[1]))
			}
			if v[1] == "工作量" {
				WorkLoadRange["工作量"] = append(WorkLoadRange["工作量"], ChangeAxisToIndex(arr[0]), ChangeAxisToIndex(arr[1]))
			}
			if v[1] == "折算比例" {
				ReducedProportionRange["折算比例"] = append(ReducedProportionRange["折算比例"], ChangeAxisToIndex(arr[0]), ChangeAxisToIndex(arr[1]))
			}
		}

		for i, row := range rows {
			var item model.SysInternalSalary
			if i == 1 {
				for k, v := range row {
					//记录不唯一内容里面的标题
					if k >= TotalBusinessVolumeRange["业务量合计"][0]-1 && k < TotalBusinessVolumeRange["业务量合计"][1] {
						//业务量合计下的标题
						TotalBusinessVolumeArr = append(TotalBusinessVolumeArr, v)
					}
					if k >= ExternalQualityRange["外部质量"][0]-1 && k < ExternalQualityRange["外部质量"][1] {
						//外部质量下的标题
						ExternalQualityArr = append(ExternalQualityArr, v)
					}
					if k >= InternalQualityRange["内部质量"][0]-1 && k < InternalQualityRange["内部质量"][1] {
						//内部质量下的标题
						InternalQualityArr = append(InternalQualityArr, v)
					}
					if k >= WorkLoadRange["工作量"][0]-1 && k < WorkLoadRange["工作量"][1] {
						//工作量下的标题
						WorkLoadArr = append(WorkLoadArr, v)
					}
					if k >= ReducedProportionRange["折算比例"][0]-1 && k < ReducedProportionRange["折算比例"][1] {
						//折算比例下的标题
						ReducedProportionArr = append(ReducedProportionArr, v)
					}
				}
			}

			if i != 0 && i != 1 {
				item.Code = row[1]
				item.NickName = row[2]
				item.DateOfEntry = "20" + strings.Split(row[3], "-")[2] + "/" + strings.Split(row[3], "-")[0] + "/" + strings.Split(row[3], "-")[1]
				item.DateOfMountGuard = "20" + strings.Split(row[4], "-")[2] + "/" + strings.Split(row[4], "-")[0] + "/" + strings.Split(row[4], "-")[1]
				item.ProjectGroup = row[5]
				item.Date = "20" + strings.Split(row[6], "-")[2] + "/" + strings.Split(row[6], "-")[0] + "/" + strings.Split(row[6], "-")[1]
				item.YearsOfEmployment = row[7]
				item.IsStayOutsideOverNight = row[8]
				//业务量合计、外部质量、内部质量、折算比例、工作量由于其的不确定内容，因此做成数组形式存进数据库
				//数组格式 : [ 标题,标题所对应的内容, 标题,标题所对应的内容.......]
				//业务量合计
				L := 0
				for k := TotalBusinessVolumeRange["业务量合计"][0] - 1; k < TotalBusinessVolumeRange["业务量合计"][1]; k++ {
					item.TotalBusinessVolume = append(item.TotalBusinessVolume, TotalBusinessVolumeArr[L])
					item.TotalBusinessVolume = append(item.TotalBusinessVolume, row[k])
					L += 1
				}
				item.OvertimeWorkLoad = row[TotalBusinessVolumeRange["业务量合计"][1]]
				item.WorkOvertimeTogether = row[TotalBusinessVolumeRange["业务量合计"][1]+1]
				item.BasicWage = row[TotalBusinessVolumeRange["业务量合计"][1]+2]
				item.PerformancePay = row[TotalBusinessVolumeRange["业务量合计"][1]+3]
				item.WageJobs = row[TotalBusinessVolumeRange["业务量合计"][1]+4]
				item.OvertimePay = row[TotalBusinessVolumeRange["业务量合计"][1]+5]
				//外部质量
				L = 0
				for k := ExternalQualityRange["外部质量"][0] - 1; k < ExternalQualityRange["外部质量"][1]; k++ {
					item.ExternalQuality = append(item.ExternalQuality, ExternalQualityArr[L])
					item.ExternalQuality = append(item.ExternalQuality, row[k])
					L += 1
				}
				//内部质量
				L = 0
				for k := InternalQualityRange["内部质量"][0] - 1; k < InternalQualityRange["内部质量"][1]; k++ {
					item.InternalQuality = append(item.InternalQuality, InternalQualityArr[L])
					item.InternalQuality = append(item.InternalQuality, row[k])
					L++
				}
				//工作量
				L = 0
				for k := WorkLoadRange["工作量"][0] - 1; k < WorkLoadRange["工作量"][1]; k++ {
					item.WorkLoad = append(item.WorkLoad, WorkLoadArr[L])
					item.WorkLoad = append(item.WorkLoad, row[k])
					L++
				}
				//折算比例
				L = 0
				for k := ReducedProportionRange["折算比例"][0] - 1; k < ReducedProportionRange["折算比例"][1]; k++ {
					item.ReducedProportion = append(item.ReducedProportion, ReducedProportionArr[L])
					item.ReducedProportion = append(item.ReducedProportion, row[k])
					L++
				}

				//出勤
				item.StandardNumberOfDays = row[AttendanceRange["出勤"][0]-1]
				item.NumberOfDaysOnTheJob = row[AttendanceRange["出勤"][0]]
				item.BereavementLeave = row[AttendanceRange["出勤"][0]+1]
				item.Late = row[AttendanceRange["出勤"][0]+2]
				item.LeaveEarly = row[AttendanceRange["出勤"][0]+3]
				item.PersonalLeave = row[AttendanceRange["出勤"][0]+4]
				item.Absenteeism = row[AttendanceRange["出勤"][0]+5]
				item.AnnualHoliday = row[AttendanceRange["出勤"][0]+6]
				item.ResignationOrNewJob = row[AttendanceRange["出勤"][0]+7]
				item.SickLeave = row[AttendanceRange["出勤"][0]+8]
				//福利
				item.AnnualWork = row[WelfareRange["福利"][0]-1]
				item.PerfectAttendanceAward = row[WelfareRange["福利"][0]]
				item.DispatchOrReserve = row[WelfareRange["福利"][0]+1]
				item.StayOutsideOverNight = row[WelfareRange["福利"][0]+2]
				item.OtherOfWelfare = row[WelfareRange["福利"][0]+3]
				//扣除
				item.Attendance = row[DeductRange["扣除"][0]-1]
				item.SocialSecurity = row[DeductRange["扣除"][0]]
				item.Quality = row[DeductRange["扣除"][0]+1]
				item.InsuranceDispatch = row[DeductRange["扣除"][0]+2]
				item.CodeOfConduct = row[DeductRange["扣除"][0]+3]
				item.OtherOfDeduct = row[DeductRange["扣除"][0]+4]
				//
				item.PieceRateWages = row[DeductRange["扣除"][1]]
				item.MinimumGuarantee = row[DeductRange["扣除"][1]+1]
				item.RealWages = row[DeductRange["扣除"][1]+2]
				item.TaxPay = row[DeductRange["扣除"][1]+3]
				item.WithholdingTax = row[DeductRange["扣除"][1]+4]
				item.CompanySupplement = row[DeductRange["扣除"][1]+5]
				if DeductRange["扣除"][1]+6 >= len(row) {
					item.Remark = ""
				} else {
					item.Remark = row[DeductRange["扣除"][1]+6]
				}

				//payday
				item.PayDay = payDay
				InternalSalarySalaries = append(InternalSalarySalaries, item)
			}

		}

		err = global.GDb.Model(&model.SysInternalSalary{}).Create(&InternalSalarySalaries).Error
		if err != nil {
			errMsg += err.Error()
			return errMsg
		}
	}
	return errMsg
}

func ChangIndexToAxis(intIndexX int, intIndexY int) string {
	var arr = [...]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	intIndexY = intIndexY + 1
	resultY := ""
	for true {
		if intIndexY <= 26 {
			resultY = resultY + arr[intIndexY-1]
			break
		}
		mo := intIndexY % 26
		resultY = arr[mo-1] + resultY
		shang := intIndexY / 26
		if shang <= 26 {
			resultY = arr[shang-1] + resultY
			break
		}
		intIndexY = shang
	}
	return resultY + strconv.Itoa(intIndexX+1)
}

func ChangeAxisToIndex(Axis string) (index int) {
	var arr = map[string]int{"A": 1, "B": 2, "C": 3, "D": 4, "E": 5, "F": 6, "G": 7, "H": 8, "I": 9, "J": 10, "K": 11,
		"L": 12, "M": 13, "N": 14, "O": 15, "P": 16, "Q": 17, "R": 18, "S": 19, "T": 20, "U": 21, "V": 22, "W": 23, "X": 24, "Y": 25, "Z": 26}
	lens := len(Axis) - 2
	if lens == 0 {
		index += arr[string(Axis[0])]
	} else {
		for i, v := range Axis {
			if i != lens {
				if _, ok := arr[string(v)]; ok {
					index += arr[string(v)] * 26
				}
			} else {
				if num, ok := arr[string(v)]; ok {
					index += num
				}
			}
		}
	}

	return
}
