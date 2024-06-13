const NoticesRoutes = {
	path: "notices",
	name: "Notices",
	meta: {
		key: "notices",
		title: "消息管理"
	},
	component: () => import("@/views/main/notices"),

	children: [
		{
			path: "/main/notices",
			redirect: "partTime"
		},

		{
			path: "partTime",
			name: "PartTime",
			meta: {
				key: "partTime",
				pKey: "notices",
				realm: "partTime",
				title: "PT群管理"
			},
			component: () => import("@/views/main/notices/partTime")
		},

		{
			path: "input",
			name: "Input",
			meta: {
				key: "input",
				pKey: "notices",
				realm: "input",
				title: "录入通知"
			},
			component: () => import("@/views/main/notices/input")
		},

		{
			path: "customer",
			name: "Customer",
			meta: {
				key: "customer",
				pKey: "notices",
				realm: "customer",
				title: "客户通知"
			},
			component: () => import("@/views/main/notices/customer")
		},

		{
			path: "business",
			name: "business",
			meta: {
				key: "business",
				pKey: "notices",
				realm: "business",
				title: "业务通知"
			},
			component: () => import("@/views/main/notices/business")
		},

		{
			path: "realTime",
			name: "RealTime",
			meta: {
				key: "realTime",
				pKey: "notices",
				realm: "realTime",
				title: "实时通知"
			},
			component: () => import("@/views/main/notices/realTime")
		},

		{
			path: "regular",
			name: "Regular",
			meta: {
				key: "regular",
				pKey: "notices",
				realm: "regular",
				title: "固定通知"
			},
			component: () => import("@/views/main/notices/regular")
		}
	]
};

export default NoticesRoutes;
