const LogRoutes = {
	path: "log",
	name: "Log",
	meta: {
		key: "log",
		title: "日志管理"
	},
	component: () => import("@/views/main/log"),
	children: [
		{
			path: "/main/log",
			redirect: "PM"
		},

		{
			path: "PM",
			name: "PM",
			meta: {
				key: "PM",
				pKey: "log",
				realm: "PM",
				title: "项目管理"
			},
			component: () => import("@/views/main/log/pm")
		},

		{
			path: "PD",
			name: "PD",
			meta: {
				key: "PD",
				pKey: "log",
				realm: "PD",
				title: "项目开发"
			},
			component: () => import("@/views/main/log/pd")
		},

		{
			path: "restart",
			name: "Restart",
			meta: {
				key: "restart",
				pKey: "log",
				realm: "restart",
				title: "系统重启"
			},
			component: () => import("@/views/main/log/restart")
		},

		{
			path: "partTime",
			name: "PartTime",
			meta: {
				key: "partTime",
				pKey: "log",
				realm: "partTime",
				title: "PT管理"
			},
			component: () => import("@/views/main/log/partTime")
		},

		{
			path: "entry",
			name: "Entry",
			meta: {
				key: "entry",
				pKey: "log",
				realm: "entry",
				title: "录入管理"
			},
			component: () => import("@/views/main/log/entry")
		}
	]
};

export default LogRoutes;
