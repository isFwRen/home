export default () => ({
	path: "case",
	name: "Case",
	meta: {
		key: "case",
		pKey: "PM",
		realm: "case",
		title: "案件列表"
	},
	component: () => import("@/views/main/pm/case"),

	children: [
		{
			path: "view-result-data",
			name: "ViewResultData",
			meta: {
				key: "viewResultData",
				pKey: "PM",
				path: "view-result-data",
				realm: "case",
				title: "查看结果数据"
			},
			component: () => import("@/views/main/pm/case/viewResultData")
		},

		{
			path: "update-entry-data",
			name: "UpdateEntryData",
			meta: {
				key: "updateEntryData",
				pKey: "PM",
				path: "update-entry-data",
				realm: "case",
				title: "修改录入数据"
			},
			component: () => import("@/views/main/pm/case/updateEntryData")
		},

		{
			path: "update-result-xml",
			name: "UpdateResultXML",
			meta: {
				key: "UpdateResultXML",
				pKey: "PM",
				path: "update-result-xml",
				realm: "case",
				title: "修改结果XML"
			},
			component: () => import("@/views/main/pm/case/updateResultXML")
		},

		{
			path: "report-table",
			name: "ReportTable",
			meta: {
				key: "reportTable",
				pKey: "PM",
				path: "report-table",
				realm: "case",
				title: "报表"
			},
			component: () => import("@/views/main/pm/case/reportTable")
		},

		{
			path: "view-image/:blockId",
			name: "ViewImage",
			meta: {
				key: "viewImage",
				pKey: "PM",
				path: "viewImage",
				realm: "case",
				title: "查看图片"
			},
			component: () => import("@/views/main/pm/case/updateEntryData/viewImage")
		}
	]
});
