export default prefixPath => ({
	path: "quality",
	name: "Quality",
	meta: {
		key: "quality",
		realm: "quality",
		pKey: "PM",
		title: "质检管理"
	},
	component: () => import("@/views/main/pm/quality"),

	children: [
		{
			path: `${prefixPath}quality`,
			redirect: "sampling-data"
		},

		{
			path: "sampling-data",
			name: "SamplingData",
			meta: {
				key: "sampling-data",
				pKey: "PM",
				path: "sampling-data",
				realm: "quality",
				title: "抽检数据"
			},
			component: () => import("@/views/main/pm/quality/samplingData")
		},

		{
			path: "error-details",
			name: "ErrorDetails",
			meta: {
				key: "error-details",
				pKey: "PM",
				path: "error-details",
				realm: "quality",
				title: "错误明细"
			},
			component: () => import("@/views/main/pm/quality/errorDetails")
		},

		{
			path: "sampling-statistics",
			name: "SamplingStatistics",
			meta: {
				key: "sampling-statistics",
				title: "抽检统计",
				pKey: "PM",
				path: "sampling-statistics",
				realm: "quality"
			},
			component: () => import("@/views/main/pm/quality/samplingStatistics")
		},

		{
			path: "sampling-setting",
			name: "SamplingSetting",
			meta: {
				key: "sampling-setting",
				pKey: "PM",
				path: "sampling-setting",
				realm: "quality",
				title: "抽检设置"
			},
			component: () => import("@/views/main/pm/quality/samplingSetting")
		}
	]
});
