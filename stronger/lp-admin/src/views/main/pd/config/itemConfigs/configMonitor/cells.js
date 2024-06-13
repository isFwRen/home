const headers = [
	{ text: "监测频率", value: "frequency" },
	{ text: "异常提示描述", value: "wrongMsg" },
	{ text: "创建时间", value: "CreatedAt", },
	{ text: "创建人", value: "createdName" },
	{ text: "更新时间", value: "UpdatedAt" }
];

const moreOptions = [
	{
		value: "delete",
		label: "删除"
	}
];

const frequencyOptions = [
	{
		value: 1,
		label: "5/分钟"
	},
	{
		value: 2,
		label: "10/分钟"
	},
	{
		value: 3,
		label: "20/分钟"
	},
	{
		value: 4,
		label: "30/分钟"
	}
];

export default {
	headers,
	moreOptions,
	frequencyOptions
};
