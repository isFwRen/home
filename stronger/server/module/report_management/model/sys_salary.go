package model

import (
	"github.com/lib/pq"
	"server/module/sys_base/model"
	"time"
)

type SysSalary struct {
	model.Model
	Code                 string    `json:"code" gorm:"comment:'工号'"`                   //工号
	NickName             string    `json:"nickName" gorm:"comment:'姓名'"`               //姓名
	ProductionStatistics string    `json:"productionStatistics" gorm:"comment:'产量合计'"` //产量合计
	ProductionSalary     string    `json:"productionSalary" gorm:"comment:'产量工资'"`     //产量工资
	CustomerComplaints   string    `json:"customerComplaints" gorm:"comment:'客户投诉'"`   //客户投诉
	ReferralBonus        string    `json:"referralBonus" gorm:"comment:'推荐奖'"`         //推荐奖
	TotalWages           string    `json:"totalWages" gorm:"comment:'工资合计'"`           //工资合计
	Tax                  string    `json:"tax" gorm:"comment:'税额'"`                    //税额
	EventuallyPay        string    `json:"eventuallyPay" gorm:"comment:'最终工资'"`        //最终工资
	PayDay               time.Time `json:"payDay"`                                     //几月份的工资表
}

type SysInternalSalary struct {
	model.Model
	Code                   string    `json:"code" gorm:"comment:'工号'"`           //工号
	NickName               string    `json:"nickName" gorm:"comment:'姓名'"`       //姓名
	PayDay                 time.Time `json:"payDay" gorm:"comment:'姓名'"`         //几月份的工资表
	DateOfEntry            string    `json:"dateOfEntry" gorm:"入职日期"`            //入职日期x`
	DateOfMountGuard       string    `json:"dateOfMountGuard" gorm:"上岗日期"`       //上岗日期
	ProjectGroup           string    `json:"projectGroup" gorm:"项目群"`            //项目群
	Date                   string    `json:"date" gorm:"日期"`                     //日期
	YearsOfEmployment      string    `json:"yearsOfEmployment" gorm:"入职年限"`      //入职年限
	IsStayOutsideOverNight string    `json:"isStayOutsideOverNight" gorm:"是否外宿"` //是否外宿

	WorkLoad            pq.StringArray `json:"workLoad" gorm:"type:varchar(100)[] comment:'工作量'"`              //工作量
	ReducedProportion   pq.StringArray `json:"reducedProportion" gorm:"type:varchar(100)[] comment:'折算比例'"`    //折算比例
	TotalBusinessVolume pq.StringArray `json:"totalBusinessVolume" gorm:"type:varchar(100)[] comment:'业务量合计'"` //业务量合计

	OvertimeWorkLoad     string `json:"overtimeWorkLoad" gorm:"加班工作量"`    //加班工作量
	WorkOvertimeTogether string `json:"workOvertimeTogether" gorm:"加班合计"` //加班合计
	BasicWage            string `json:"basicWage" gorm:"基本工资"`            //基本工资
	PerformancePay       string `json:"performancePay" gorm:"绩效工资"`       //绩效工资
	WageJobs             string `json:"wageJobs" gorm:"岗位工资"`             //岗位工资
	OvertimePay          string `json:"overtimePay" gorm:"加班工资"`          //加班工资

	ExternalQuality pq.StringArray `json:"externalQuality" gorm:"type:varchar(100)[] comment:'外部质量'"` //外部质量
	InternalQuality pq.StringArray `json:"internalQuality" gorm:"type:varchar(100)[] comment:'内部质量'"` //内部质量

	//出勤
	StandardNumberOfDays string `json:"standardNumberOfDays" gorm:"标准天数"`  //标准天数
	NumberOfDaysOnTheJob string `json:"numberOfDaysOnTheJob" gorm:"上岗天数"`  //上岗天数
	BereavementLeave     string `json:"bereavementLeave" gorm:"丧假"`        //丧假
	Late                 string `json:"late" gorm:"迟到"`                    //迟到
	LeaveEarly           string `json:"leaveEarly" gorm:"早退"`              //早退
	PersonalLeave        string `json:"personalLeave" gorm:"事假"`           //事假
	Absenteeism          string `json:"absenteeism" gorm:"旷工"`             //旷工
	AnnualHoliday        string `json:"annualHoliday" gorm:"年休"`           //年休
	ResignationOrNewJob  string `json:"resignationOrNewJob" gorm:"离职/新上岗"` //离职/新上岗
	SickLeave            string `json:"sickLeave" gorm:"病假"`               //病假
	//福利
	AnnualWork             string `json:"annualWork" gorm:"年工"`              //年工
	PerfectAttendanceAward string `json:"perfectAttendanceAward" gorm:"全勤奖"` //全勤奖
	DispatchOrReserve      string `json:"dispatchOrReserve" gorm:"调度/储备"`    //调度/储备
	StayOutsideOverNight   string `json:"stayOutsideOverNight" gorm:"外宿"`    //外宿
	OtherOfWelfare         string `json:"otherOfWelfare" gorm:"其他"`          //其他福利
	//扣除
	Attendance        string `json:"attendance" gorm:"考勤"`          //考勤
	SocialSecurity    string `json:"socialSecurity" gorm:"社保"`      //社保
	Quality           string `json:"quality" gorm:"质量"`             //质量
	InsuranceDispatch string `json:"insuranceDispatch" gorm:"保险调度"` //保险调度
	CodeOfConduct     string `json:"codeOfConduct" gorm:"行为规范"`     //行为规范
	OtherOfDeduct     string `json:"otherOfDeduct" gorm:"其他"`       //其他扣除

	PieceRateWages    string `json:"pieceRateWages" gorm:"计件工资"`   //计件工资
	MinimumGuarantee  string `json:"minimumGuarantee" gorm:"最低保障"` //最低保障
	RealWages         string `json:"realWages" gorm:"实际工资"`        //实际工资
	TaxPay            string `json:"taxPay" gorm:"代扣个税"`           //代扣个税
	WithholdingTax    string `json:"withholdingTax" gorm:"本月发放"`   //本月发放
	CompanySupplement string `json:"companySupplement" gorm:"公司补"` //公司补
	Remark            string `json:"remark" gorm:"备注"`             //备注
}

/* 关于交付工资结构体
1、业务量合计
 TotalBusinessVolume 数组 每个下标对应的意思如下
	1、B0101
	2、B0101-T
	3、B0101-Q
	4、B0101-D
	5、B0103
	6、B0103-T
	7、B0103-Q
	8、B0103-D
	9、B0106
	10、B0106-T
	11、B0106-Q
	12、B0106-D
	13、B0108
	14、B0108-T
	15、B0108-Q
	16、B0108-D
	17、B0110
	18、B0110-T
	19、B0110-Q
	20、B0110-D
	21、B0111
	22、B0111-T
	23、B0111-Q
	24、B0111-D
	25、B0113
	26、B0113-T
	27、B0113-Q
	28、B0113-D
	29、B0114
	30、B0114-T
	31、B0114-Q
	32、B0114-D
	33、B0115
	34、B0115-T
	35、B0115-Q
	36、B0115-D
	37、B0116
	38、B0116-T
	39、B0116-Q
	40、B0116-D
	41、B0117
	42、B0117-T
	43、B0117-Q
	44、B0117-D
	45、B0118
	46、B0118-T
	47、B0118-Q
	48、B0118-D
	49、B0120
	50、B0120-T
	51、B0120-Q
	52、B0120-D
	53、B0153
	54、B0153-Q
	55、B0153-D
	56、B0121
	57、B0121-T
	58、B0121-Q
	59、B0121-D
	60、双录
	61、B0101-已导出
	62、B0103-已导出
	63、B0103-待审核
	64、B0106-已导出
	65、B0106-待审核
	66、B0108-已导出
	67、B0108-待审核
	68、B0110-已导出
	69、B0110-待审核
	70、B0113质检
	71、B0114质检
	72、B0115质检
	73、B0116质检
	74、B0117质检
	75、B0118质检
	76、B0120质检
	77、B0121质检
	78、其他字符补贴
	79、培训产量扣除
	80、合计
	81、加班工作量
	82、加班合计
2、外部/内部工资不固定内容
*/

type SysInternalSalaryVersionTwo struct {
	model.Model
	Code                   string            `json:"code" gorm:"comment:'工号'"`                   //工号
	NickName               string            `json:"nickName" gorm:"comment:'姓名'"`               //姓名
	PayDay                 string            `json:"payDay" gorm:"comment:'姓名'"`                 //发薪日
	DateOfEntry            string            `json:"dateOfEntry" gorm:"入职日期"`                    //入职日期x`
	DateOfMountGuard       string            `json:"dateOfMountGuard" gorm:"上岗日期"`               //上岗日期
	ProjectGroup           string            `json:"projectGroup" gorm:"项目群"`                    //项目群
	Date                   string            `json:"date" gorm:"日期"`                             //日期
	YearsOfEmployment      string            `json:"yearsOfEmployment" gorm:"入职年限"`              //入职年限
	IsStayOutsideOverNight string            `json:"isStayOutsideOverNight" gorm:"是否外宿"`         //是否外宿
	TotalBusinessVolume    map[string]string `json:"totalBusinessVolume" gorm:"comment:'业务量合计'"` //业务量合计
	BasicWage              string            `json:"basicWage" gorm:"基本工资"`                      //基本工资
	PerformancePay         string            `json:"performancePay" gorm:"绩效工资"`                 //绩效工资
	WageJobs               string            `json:"wageJobs" gorm:"岗位工资"`                       //岗位工资
	OvertimePay            string            `json:"overtimePay" gorm:"加班工资"`                    //加班工资
	ExternalQuality        map[string]string `json:"externalQuality" gorm:"comment:'外部质量'"`      //外部质量
	InternalQuality        map[string]string `json:"internalQuality" gorm:"comment:'内部质量'"`      //内部质量

	WorkLoad          map[string]string `json:"workLoad" gorm:"comment:'工作量'"`           //工作量
	ReducedProportion map[string]string `json:"reducedProportion" gorm:"comment:'折算比例'"` //折算比例

	OvertimeWorkLoad     string `json:"overtimeWorkLoad" gorm:"加班工作量"`    //加班工作量
	WorkOvertimeTogether string `json:"workOvertimeTogether" gorm:"加班合计"` //加班合计
	//出勤
	StandardNumberOfDays string `json:"standardNumberOfDays" gorm:"标准天数"`  //标准天数
	NumberOfDaysOnTheJob string `json:"numberOfDaysOnTheJob" gorm:"上岗天数"`  //上岗天数
	BereavementLeave     string `json:"bereavementLeave" gorm:"丧假"`        //丧假
	Late                 string `json:"late" gorm:"迟到"`                    //迟到
	LeaveEarly           string `json:"leaveEarly" gorm:"早退"`              //早退
	PersonalLeave        string `json:"personalLeave" gorm:"事假"`           //事假
	Absenteeism          string `json:"absenteeism" gorm:"旷工"`             //旷工
	AnnualHoliday        string `json:"annualHoliday" gorm:"年休"`           //年休
	ResignationOrNewJob  string `json:"resignationOrNewJob" gorm:"离职/新上岗"` //离职/新上岗
	SickLeave            string `json:"sickLeave" gorm:"病假"`               //病假
	//福利
	AnnualWork             string `json:"annualWork" gorm:"年工"`              //年工
	PerfectAttendanceAward string `json:"perfectAttendanceAward" gorm:"全勤奖"` //全勤奖
	DispatchOrReserve      string `json:"dispatchOrReserve" gorm:"调度/储备"`    //调度/储备
	StayOutsideOverNight   string `json:"stayOutsideOverNight" gorm:"外宿"`    //外宿
	OtherOfWelfare         string `json:"otherOfWelfare" gorm:"其他"`          //其他福利
	//扣除
	Attendance        string `json:"attendance" gorm:"考勤"`          //考勤
	SocialSecurity    string `json:"socialSecurity" gorm:"社保"`      //社保
	Quality           string `json:"quality" gorm:"质量"`             //质量
	InsuranceDispatch string `json:"insuranceDispatch" gorm:"保险调度"` //保险调度
	CodeOfConduct     string `json:"codeOfConduct" gorm:"行为规范"`     //行为规范
	OtherOfDeduct     string `json:"otherOfDeduct" gorm:"其他"`       //其他扣除

	PieceRateWages    string `json:"pieceRateWages" gorm:"计件工资"`   //计件工资
	MinimumGuarantee  string `json:"minimumGuarantee" gorm:"最低保障"` //最低保障
	RealWages         string `json:"realWages" gorm:"实际工资"`        //实际工资
	WithholdingTax    string `json:"withholdingTax" gorm:"本月发放"`   //本月发放
	TaxPay            string `json:"taxPay" gorm:"代扣个税"`           //代扣个税
	CompanySupplement string `json:"companySupplement" gorm:"公司补"` //公司补
	Remark            string `json:"remark" gorm:"备注"`             //备注
}
