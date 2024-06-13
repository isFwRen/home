const HomeRoutes = {
	path: "home",
	name: "Home",
	meta: {
		key: "home",
		title: "首页"
	},
	component: () => import("@/views/main/home"),

	children: [
		{
			path: "/main/home",
			redirect: "daily"
		},

		{
			path: "daily",
			name: "Daily",
			meta: {
				key: "daily",
				pKey: "home",
				realm: "daily",
				title: "项目日报"
			},
			component: () => import("@/views/main/home/daily")
		},

		{
			path: "data",
			name: "Data",
			meta: {
				key: "data",
				pKey: "home",
				realm: "data",
				title: "项目数据"
			},
			component: () => import("@/views/main/home/data")
		},

		{
			path: "homepage",
			name: "Homepage",
			meta: {
				key: "index",
				pKey: "home",
				realm: "homepage",
				title: "我的主页"
			},
			component: () => import("@/views/main/home/homepage")
		},
		{
			path: "ranking",
			name: "Ranking",
			meta: {
				key: "index",
				pKey: "home",
				realm: "ranking",
				title: "近日排行"
			},
			component: () => import("@/views/main/home/homepage/ranking")
		}
	]
};

export default HomeRoutes;
