const prefixPath = "/main/home/";

const tabsOptions = [
	{
		index: 0,
		key: "daily",
		label: "项目日报",
		path: `${prefixPath}daily`
	},

	{
		index: 1,
		key: "data",
		label: "项目数据",
		path: `${prefixPath}data`
	},

	{
		index: 2,
		key: "homepage",
		label: "我的主页",
		path: `${prefixPath}homepage`
	}
];

export default {
	tabsOptions
};
