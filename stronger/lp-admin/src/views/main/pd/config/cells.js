const headers = [
	{ text: "所属行业", value: "type", width: 160 },
	{ text: "模板", value: "template" },
	{ text: "缓存时间(秒)", value: "cacheTime", width: 90 },
	{ text: "回传方式", value: "autoReturn", width: 70 },
	{ text: "数据保存天数", value: "saveDate", width: 60 },
	{ text: "自动重启", value: "restartAt", width: 70 },
	{ text: "配置", value: "options" },
	{ text: "设置", value: "settings", width: 115 }
];

const options = [
	{
		class: "pr-3",
		icon: "mdi-plus",
		key: 1,
		text: "创建新项目"
	}
];

const configs = [
	{
		class: "mb-1 mr-1",
		key: 1,
		text: "字段配置",
		path: "/main/PD/config/field"
	},

	{
		class: "mb-1 mr-1",
		key: 2,
		text: "导出配置",
		path: "/main/PD/config/export"
	},

	{
		class: "mb-1 mr-1",
		key: 3,
		text: "审核配置",
		path: "/main/PD/config/audit"
	},

	{
		class: "mb-1 mr-1",
		key: 4,
		text: "质检配置",
		path: "/main/PD/config/quality"
	},

	{
		class: "mb-1 mr-1",
		key: 5,
		text: "时效配置",
		path: "/main/PD/config/prescription"
	},

	{
		key: 6,
		text: "路径配置",
		path: "/main/PD/config/path"
	},

	{
		key: 7,
		text: "监控配置",
		path: "/main/PD/config/monitor"
	}
];

const settings = [
	{
		class: "pr-3",
		icon: "mdi-autorenew",
		key: "entry",
		text: "刷新录入配置",
		disabled: true
	},

	{
		class: "pr-3",
		icon: "mdi-autorenew",
		key: "manage",
		text: "刷新管理配置",
		disabled: true
	}

	// {
	//   icon: 'mdi-autorenew',
	//   key: 'all',
	//   text: '刷新全部配置',
	//   disabled: true
	// }
];

export default {
	headers,
	options,
	configs,
	settings
};
