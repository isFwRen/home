const menus = [
	{
		pId: "-1",
		id: "home",
		key: "home",
		icon: "mdi-home-outline",
		title: "首页",
		link: "/main/home",
		leaf: false
	},

	{
		pId: "-1",
		id: "PM",
		key: "PM",
		icon: "mdi-folder-star-multiple-outline",
		title: "项目管理",
		link: "/main/PM",
		leaf: true,
		children: [
			{
				pId: "PM",
				key: "case",
				realm: "case",
				icon: "",
				title: "案件列表",
				link: "/main/PM/case"
			},

			{
				pId: "PM",
				key: "task",
				realm: "task",
				icon: "",
				title: "任务管理",
				link: "/main/PM/task"
			},

			{
				pId: "PM",
				key: "prescription",
				realm: "prescription",
				icon: "",
				title: "时效管理",
				link: "/main/PM/prescription"
			},

			{
				pId: "PM",
				key: "qualities",
				realm: "qualities",
				icon: "",
				title: "质量管理",
				link: "/main/PM/qualities"
			},

			{
				pId: "PM",
				key: "staff",
				realm: "staff",
				icon: "",
				title: "人员管理",
				link: "/main/PM/staff"
			},

			{
				pId: "PM",
				key: "quality",
				realm: "quality",
				icon: "",
				title: "质检管理",
				link: "/main/PM/quality"
			},

			{
				pId: "PM",
				key: "teaching",
				realm: "teaching",
				icon: "",
				title: "教学管理",
				link: "/main/PM/teaching"
			},

			{
				pId: "PM",
				key: "notice",
				realm: "notice",
				icon: "",
				title: "公告管理",
				link: "/main/PM/notice"
			}
		]
	},

	{
		pId: "-1",
		id: "report",
		key: "report",
		icon: "mdi-chart-bell-curve-cumulative",
		title: "报表管理",
		link: "/main/report",
		leaf: true,
		children: [
			{
				pId: "report",
				key: "yield",
				realm: "yield",
				icon: "",
				title: "产量统计",
				link: "/main/report/yield"
			},

			{
				pId: "report",
				key: "error",
				realm: "error",
				icon: "",
				title: "错误查询",
				link: "/main/report/error"
			},

			{
				pId: "report",
				key: "salary",
				realm: "salary",
				icon: "",
				title: "工资数据",
				link: "/main/report/salary"
			},

			{
				pId: "report",
				key: "itemReport",
				realm: "itemReport",
				icon: "",
				title: "项目报表",
				link: "/main/report/itemReport"
			},

			{
				pId: "report",
				key: "specialReport",
				realm: "specialReport",
				icon: "",
				title: "特殊报表",
				link: "/main/report/specialReport"
			}
		]
	},

	{
		pId: "-1",
		id: "entry",
		key: "entry",
		icon: "mdi-circle-edit-outline",
		title: "录入管理",
		link: "/main/entry",
		leaf: true,
		children: [
			{
				pId: "entry",
				key: "channel",
				realm: "channel",
				icon: "",
				title: "录入通道",
				link: "/main/entry/channel"
			}
		]
	},

	{
		pId: "-1",
		id: "PD",
		key: "PD",
		icon: "mdi-cellphone-link",
		title: "项目开发",
		link: "/main/PD",
		leaf: true,
		children: [
			{
				pId: "PD",
				key: "config",
				realm: "config",
				icon: "",
				title: "配置管理",
				link: "/main/PD/config"
			},

			{
				pId: "PD",
				key: "constant",
				realm: "constant",
				icon: "",
				title: "常量管理",
				link: "/main/PD/constant"
			}
		]
	},

	{
		pId: "-1",
		id: "notices",
		key: "notices",
		icon: "mdi-bulletin-board",
		title: "消息管理",
		link: "/main/notices",
		leaf: true,
		children: [
			{
				pId: "notices",
				key: "partTime",
				realm: "partTime",
				icon: "",
				title: "PT群管理",
				link: "/main/notices/partTime"
			},

			{
				pId: "notices",
				key: "customer",
				realm: "customer",
				icon: "",
				title: "客户通知",
				link: "/main/notices/customer"
			},

			{
				pId: "notices",
				key: "input",
				realm: "input",
				icon: "",
				title: "录入通知",
				link: "/main/notices/input"
			},

			{
				pId: "notices",
				key: "realTime",
				realm: "realTime",
				icon: "",
				title: "实时通知",
				link: "/main/notices/realTime"
			},

			{
				pId: "notices",
				key: "regular",
				realm: "regular",
				icon: "",
				title: "固定通知",
				link: "/main/notices/regular"
			},
			{
				pId: "notices",
				key: "business",
				realm: "business",
				icon: "",
				title: "业务通知",
				link: "/main/notices/business"
			}
		]
	},

	{
		pId: "-1",
		id: "log",
		key: "log",
		icon: "mdi-notebook-outline",
		title: "日志管理",
		link: "/main/log",
		leaf: true,
		children: [
			{
				pId: "log",
				key: "PM",
				realm: "PM",
				icon: "",
				title: "项目管理",
				link: "/main/log/PM"
			},

			{
				pId: "log",
				key: "PD",
				realm: "PD",
				icon: "",
				title: "项目开发",
				link: "/main/log/PD"
			},

			{
				pId: "log",
				key: "restart",
				realm: "restart",
				icon: "",
				title: "系统重启",
				link: "/main/log/restart"
			},

			{
				pId: "log",
				key: "partTime",
				realm: "partTime",
				icon: "",
				title: "PT管理",
				link: "/main/log/partTime"
			},

			{
				pId: "log",
				key: "entry",
				realm: "entry",
				icon: "",
				title: "录入管理",
				link: "/main/log/entry"
			}
		]
	}
];

export default menus;
