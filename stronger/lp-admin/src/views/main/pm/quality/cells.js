const prefixPath = "/main/PM/quality/";

const tabsOptions = [
	{
		key: "sampling-data",
		label: "抽检数据",
		path: `${prefixPath}sampling-data`
	},

	{
		key: "error-details",
		label: "抽检明细",
		path: `${prefixPath}error-details`
	},

	{
		key: "sampling-statistics",
		label: "抽检统计",
		path: `${prefixPath}sampling-statistics`
	},

	{
		key: "sampling-setting",
		label: "抽检设置",
		path: `${prefixPath}sampling-setting`
	}
];

export default {
	tabsOptions
};
