import moment from "moment";

const DEFAULT_DATE = [moment().format("YYYY-MM-DD"), moment().format("YYYY-MM-DD")];

// 项目
const proCodeOptions = [];

// 案件状态
const statusOptions = [];

// 医保类型
const insuranceTypeOptions = [];

// 理赔类型
const claimTypeOptions = [];

// 加急件
const stickLevelOptions = [];

// 问题件
const isQuestionOptions = [
	{
		label: "是",
		value: 1
	},
	{
		label: "否",
		value: 2
	}
];

// 录入状态
const stageOptions = [];

const fields = [
	{
		cols: 2,
		formKey: "proCode",
		inputType: "select ",
		hideDetails: true,
		label: "项目",
		options: [],
		hasDefaultValue: true
	},

	{
		cols: 3,
		formKey: "time",
		inputType: "date",
		hideDetails: true,
		label: "时间",
		clearable: false,
		defaultValue: DEFAULT_DATE
	},
	{
		cols: 2,
		formKey: "batchNum",
		inputType: "input",
		hideDetails: true,
		label: "批次号"
	},

	{
		cols: 2,
		formKey: "billCode",
		inputType: "input",
		hideDetails: true,
		label: "案件号"
	},

	{
		cols: 2,
		formKey: "stage",
		inputType: "select",
		hideDetails: true,
		label: "录入状态",
		clearable: true,
		options: []
	}
];

const moreFields = [
	{
		cols: 2,
		formKey: "saleChannel",
		inputType: "input",
		hideDetails: true,
		label: "销售渠道"
	},

	{
		cols: 2,
		formKey: "agency",
		inputType: "input",
		hideDetails: true,
		label: "机构号"
	},

	{
		cols: 2,
		formKey: "insuranceType",
		inputType: "input",
		hideDetails: true,
		label: "单证类型",
		clearable: true,
		options: []
	},

	{
		cols: 2,
		formKey: "claimType",
		inputType: "select",
		hideDetails: true,
		label: "理赔类型",
		clearable: true,
		options: []
	},

	{
		cols: 2,
		formKey: "stickLevel",
		inputType: "select",
		hideDetails: true,
		label: "加急件",
		clearable: true,
		options: []
	},

	{
		cols: 2,
		formKey: "minCountMoney",
		inputType: "input",
		hideDetails: true,
		label: "账单金额(最低额)"
	},

	{
		appendOuter: "-",
		class: "ml-n4",
		cols: 2,
		formKey: "maxCountMoney",
		inputType: "input",
		hideDetails: true,
		label: "账单金额(最高额)"
	},

	{
		cols: 2,
		formKey: "isQuestion",
		inputType: "select",
		hideDetails: true,
		label: "问题件",
		clearable: true,
		options: []
	},

	{
		cols: 2,
		formKey: "invoiceNum",
		inputType: "input",
		hideDetails: true,
		label: "发票数量"
	},

	{
		cols: 2,
		formKey: "qualityUser",
		inputType: "input",
		hideDetails: true,
		label: "质检人"
	},
	{
		cols: 2,
		formKey: "status",
		inputType: "select",
		hideDetails: true,
		label: "案件状态",
		clearable: true,
		options: []
	}
];

const headers = [
	{ text: "时间", value: "CreatedAt", width: 120, sortable: true, fixed: "left" },
	{ text: "案件号", value: "billNum", width: 190, sortable: true, fixed: "left" },
	{ text: "批次号", value: "batchNum", width: 150, sortable: true },
	{ text: "机构", value: "agency", width: 75, sortable: true },
	{ text: "扫描时间", value: "scanAt", width: 75, sortable: true },
	{ text: "导出时间", value: "exportAt", width: 70, sortable: true },
	{ text: "回传时间", value: "lastUploadAt", width: 70, sortable: true },
	{ text: "案件状态", value: "status", width: 80, sortable: true },
	{ text: "录入状态", value: "stage", width: 80, sortable: true },
	{ text: "销售渠道", value: "saleChannel", width: 80, sortable: true },
	{ text: "单证类型", value: "insuranceType", width: 80, sortable: true },
	{ text: "理赔类型", value: "claimType", width: 80, sortable: true },
	{ text: "账单金额", value: "countMoney", width: 80, sortable: true },
	{ text: "发票数量", value: "invoiceNum", width: 80, sortable: true },
	{ text: "问题件", value: "questionNum", width: 60, sortable: true },
	{ text: "加急件", value: "stickLevel", width: 60, sortable: true },
	{ text: "质检人", value: "qualityUserCode", width: 80, fixed: "right", sortable: true },
	// { text: "质检状态", value: "exportStage", width: 80, fixed: "right"},
	{ text: "备注", value: "remark", width: 100, fixed: "right", sortable: true },
	{ text: "处理", value: "dealWith", width: 100, fixed: "right" },
	{ text: "手动回传", value: "isAutoUpload", width: 120, fixed: "right" },
	{ text: "操作", value: "options", width: 130, fixed: "right" }
];

// 手动回传
const autoUploadOptions = [
	{
		value: false,
		label: "是"
	},

	{
		value: true,
		label: "否"
	}
];

// 更多
const moreOptions = [
	{
		value: 1,
		label: "删除"
	},

	{
		value: 2,
		label: "重加载",
		title: "重加载提示",
		content: "是否重加载案件："
	},

	{
		value: 3,
		label: "恢复",
		title: "恢复提示",
		content: "是否恢复案件："
	},

	{
		value: 4,
		label: "导出异常",
		title: "导出异常提示",
		content: "是否导出异常案件："
	},

	{
		value: 5,
		label: "强制导出",
		title: "强制导出提示",
		content: "是否强制导出案件："
	}
];

// 操作
const optionsOptions = [
	{
		value: 1,
		label: "查看结果数据",
		icon: "mdi-magnify",
		class: "mr-1",
		path: "/main/PM/case/view-result-data"
	},

	{
		value: 2,
		label: "修改录入数据",
		icon: "mdi-pencil-circle",
		class: "mr-1",
		path: "/main/PM/case/update-entry-data"
	},

	{
		value: 3,
		label: "修改结果XML",
		icon: "mdi-file",
		class: "mr-1",
		path: "/main/PM/case/update-result-xml"
	},

	{
		value: 4,
		label: "回传",
		icon: "mdi-arrow-left-top-bold",
		disabled: true
	}
];

const deleteFields = [
	{
		formKey: "delRemarks",
		inputType: "textarea",
		hideDetails: false,
		label: "备注",
		prependOuter: "*",
		prependOuterClass: "error--text",
		validation: [{ rule: "required", message: "备注不能为空." }]
	}
];

export default {
	proCodeOptions,
	statusOptions,
	insuranceTypeOptions,
	claimTypeOptions,
	stickLevelOptions,
	isQuestionOptions,
	stageOptions,
	fields,
	moreFields,
	moreOptions,
	headers,
	autoUploadOptions,
	optionsOptions,
	deleteFields
};
