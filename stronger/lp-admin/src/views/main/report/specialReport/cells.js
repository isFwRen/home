const prefixPath = "/main/report/specialReport/";

const tabsOptions = [
	{
		label: "来量分析",
		key: "income-analysis",
		path: `${prefixPath}income-analysis`
	},

	{
		label: "回传分析",
		key: "return-analysis",
		path: `${prefixPath}return-analysis`
	},
	{
		label: "目录外数据",
		key: "directory-out",
		path: `${prefixPath}directory-out`
	},
	{
		label: "异常件数据",
		key: "abnormal-part",
		path: `${prefixPath}abnormal-part`
	},
	{
		label: "机构抽取",
		key: "institutional-extraction",
		path: `${prefixPath}institutional-extraction`
	},
	{
		label: "ocr识别统计",
		key: "identify-statistics",
		path: `${prefixPath}identify-statistics`
	},
	{
		label: "销毁报告",
		key: "destruction-report",
		path: `${prefixPath}destruction-report`
	},
	{
		label: "扣费明细",
		key: "deduction-details",
		path: `${prefixPath}deduction-details`
	},
];

export default {
	tabsOptions
};
