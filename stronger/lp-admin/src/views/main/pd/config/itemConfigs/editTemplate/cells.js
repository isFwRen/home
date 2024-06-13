const UNDEFINED = "未定义";

const headers = [
	// { text: '序号', value: 'myOrder', width: 60 },
	{ text: "分块编码", value: "code", minWidth: 90 },
	{ text: "分块名称", value: "name", minWidth: 150 },
	{ text: "F8", value: "fEight", minWidth: 120 },
	{ text: "OCR", value: "ocr", minWidth: 120 },
	{ text: "释放时间(秒)", value: "freeTime", width: 100 },
	{ text: "关联", value: "relation", minWidth: 100 },
	{ text: "截图配置", value: "screenshot", coordinateType: true, width: 70 },
	{ text: "流程(仅限只需录入一码的分块)", value: "isCompetitive", minWidth: 140 },
	{ text: "循环分块", value: "isLoop", minWidth: 120 },
	{ text: "手机录入", value: "isMobile", minWidth: 120 },
	{ text: "手机截图", value: "mobileScreenshot", coordinateType: false, width: 70 },
	{ text: "删除", value: "delete", width: 70 }
];

const configs = [
	{
		class: "mr-2",
		fn: "copyTemplate",
		label: "复制模板",
		value: 1
	},

	{
		class: "mr-2",
		fn: "modifyTemplateImg",
		label: "修改模板图片",
		value: 2
	},

	{
		label: "脱敏配置",
		fn: "desensitization",
		value: 3,
		disabled: true
	}

	// {
	//   value: 4,
	//   label: '保存配置'
	// }
];

const codes = [
	{
		value: true,
		label: "一码"
	},

	{
		value: false,
		label: "二码"
	}
];

const competitive = [
	{
		value: true,
		label: "通用"
	},

	{
		value: false,
		label: "不通用"
	}
];

const choice = [
	{
		value: true,
		label: "是"
	},

	{
		value: false,
		label: "否"
	}
];

export default {
	UNDEFINED,
	headers,
	configs,
	codes,
	competitive,
	choice
};
