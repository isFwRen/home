import moment from "moment";

const headers = [
	// { text: '序号', value: 'orderNumber' },
	{ text: "工号", value: "code" },
	{ text: "姓名", value: "nickName" },
	{ text: "入职日期", value: "dateOfEntry", width: "80px" },
	{ text: "上岗日期", value: "dateOfMountGuard", width: "80px" },
	{ text: "项目群", value: "projectGroup", width: "80px" },
	{ text: "日期", value: "date", width: "80px" },
	{ text: "入职年限", value: "yearsOfEmployment", width: "70px" },
	{ text: "是否外宿", value: "isStayOutsideOverNight", width: "70px" },
	// { text: '岗位级别', value: 'postLevel', width: '70px' },
	// { text: '岗位最低产量', value: 'minYield', width: '100px' },
	// { text: '是否达标', value: 'getStandard', width: '70px' },

	{
		text: "工作量",
		value: "workLoad",
		children: []
	},
	{
		text: "折算比例",
		value: "reducedProportion",
		children: []
	},
	{
		text: "业务量合计",
		value: "totalBusinessVolume",
		children: []
	},

	{
		text: "",
		value: "workOvertime",
		children: [
			{ text: "加班工作量", value: "overtimeWorkLoad", width: "80px" },
			{ text: "加班合计", value: "workOvertimeTogether" }
		]
	},
	{ text: "基本工资", value: "basicWage", width: "70px" },
	{ text: "绩效工资", value: "performancePay", width: "75px" },
	{ text: "岗位工资", value: "wageJobs", width: "75px" },
	{ text: "加班工资", value: "overtimePay", width: "90px" },
	{
		text: "外部质量",
		value: "externalQuality",
		children: []
	},
	{
		text: "内部质量",
		value: "internalQuality",
		children: []
	},
	{
		text: "出勤",
		value: "attend",
		children: [
			{ text: "标准天数", value: "standardNumberOfDays" },
			{ text: "上岗天数", value: "numberOfDaysOnTheJob" },
			{ text: "丧假", value: "bereavementLeave" },
			{ text: "迟到", value: "late" },
			{ text: "早退", value: "leaveEarly" },
			{ text: "事假", value: "personalLeave" },
			{ text: "旷工", value: "absenteeism" },
			{ text: "年休", value: "annualHoliday" },
			{ text: "离职/新入职", value: "resignationOrNewJob", width: "90px" },
			{ text: "病假", value: "sickLeave" }
		]
	},
	{
		text: "福利",
		value: "welfare",
		children: [
			{ text: "年工", value: "annualWork" },
			{ text: "全勤奖", value: "perfectAttendanceAward" },
			{ text: "调度/储备", value: "dispatchOrReserve", width: "80px" },
			{ text: "外宿", value: "stayOutsideOverNight" },
			{ text: "其他", value: "otherOfWelfare" }
		]
	},
	{
		text: "扣除",
		value: "deduction",
		children: [
			{ text: "考勤", value: "attendance" },
			{ text: "社保", value: "socialSecurity" },
			{ text: "质量", value: "quality" },
			{ text: "保险调度", value: "insuranceDispatch" },
			{ text: "行为规范", value: "codeOfConduct" },
			{ text: "其他", value: "otherOfDeduct" }
		]
	},
	{ text: "计件工资", value: "pieceRateWages", width: "70px" },
	{ text: "最低保障", value: "minimumGuarantee", width: "70px" },
	{ text: "实际工资", value: "realWages", width: "70px" },
	{ text: "代扣个税", value: "taxPay", width: "70px" },
	{ text: "本月发放", value: "withholdingTax", width: "70px" },
	{ text: "公司补", value: "companySupplement", width: "70px" },
	{ text: "备注", value: "remark", width: "140px" }
];

const fields = [
	{
		cols: 3,
		formKey: "date",
		inputType: "date",
		hideDetails: true,
		label: "日期",
		clearable: true,
		pickerType: "month",
		defaultValue: moment().format("YYYY-MM")
	},

	{
		cols: 2,
		formKey: "code",
		inputType: "input",
		hideDetails: true,
		label: "工号",
		defaultValue: undefined
	},
	{
		cols: 2,
		formKey: "name",
		inputType: "input",
		hideDetails: true,
		label: "姓名",
		defaultValue: undefined
	}
];

const btns = [
	{
		class: "pr-3",
		key: 2,
		color: "primary",
		text: "导出"
	}
];

export default {
	headers,
	fields,
	btns
};
