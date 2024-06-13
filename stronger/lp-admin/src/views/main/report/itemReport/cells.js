const prefixPath = "/main/report/itemReport/";

const tabsOptions = [
	{
		label: "业务明细",
		key: " business-detail",
		path: `${prefixPath}business-detail`
	},

	{
		label: "字符统计",
		key: "character-statistics",
		path: `${prefixPath}character-statistics`
	},
	{
		label: "日报表",
		key: "day-report",
		path: `${prefixPath}day-report`
	},
	{
		label: "周报表",
		key: "week-report",
		path: `${prefixPath}week-report`
	},
	{
		label: "月报表",
		key: "month-report",
		path: `${prefixPath}month-report`
	},

	{
		label: "年报表",
		key: "year-report",
		path: `${prefixPath}year-report`
	},
	{
		label: "项目报表",
		key: "project-report",
		path: `${prefixPath}project-report`
	},
	{
		label: "项目结算",
		key: "project-settlement",
		path: `${prefixPath}project-settlement`
	}
];

export default {
	tabsOptions
};
