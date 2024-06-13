// const prefixPath = '/main/PM/'

export default prefixPath => ({
	path: "staff",
	name: "Staff",
	meta: {
		key: "staff",
		realm: "staff",
		pKey: "PM",
		title: "人员管理"
	},
	component: () => import("@/views/main/pm/staff"),

	children: [
		{
			path: `${prefixPath}staff`,
			redirect: "user"
		},

		{
			path: "user",
			name: "user",
			meta: {
				key: "user",
				pKey: "PM",
				path: "user",
				realm: "staff",
				title: "角色管理"
			},
			component: () => import("@/views/main/pm/staff/user")
		},

		{
			path: "user",
			name: "User",
			meta: {
				key: "user",
				pKey: "PM",
				path: "user",
				realm: "staff",
				title: "用户管理"
			},
			component: () => import("@/views/main/pm/staff/user")
		},

		{
			path: "whiteList",
			name: "WhiteList",
			meta: {
				key: "whiteList",
				pKey: "PM",
				path: "whiteList",
				realm: "staff",
				title: "白名单"
			},
			component: () => import("@/views/main/pm/staff/whiteList")
		},

		{
			path: "permission",
			name: "Permission",
			meta: {
				key: "permission",
				pKey: "PM",
				path: "permission",
				realm: "staff",
				title: "项目权限"
			},
			component: () => import("@/views/main/pm/staff/permission")
		},

		{
			path: "menuManage",
			name: "MenuManage",
			meta: {
				key: "menuManage",
				pKey: "PM",
				path: "menuManage",
				realm: "staff",
				title: "菜单管理"
			},
			component: () => import("@/views/main/pm/staff/menuManage")
		}
	]
});
