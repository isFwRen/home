const NormalRoutes = {
	path: "/normal",
	name: "Normal",
	meta: {
		key: "normal",
		title: ""
	},
	component: () => import("@/views/normal"),
	children: [
		{
			path: "/normal",
			redirect: "view-images"
		},

		{
			path: "view-images",
			name: "ViewImages",
			meta: {
				key: "view-images",
				pKey: "normal",
				realm: "view-images",
				title: "查看图片"
			},
			component: () => import("@/views/normal/viewImages")
		},
		{
			path: "pdf-file",
			name: "PdfFile",
			meta: {
				key: "pdf-file",
				pKey: "normal",
				realm: "pdf-file",
				title: "查看PDF"
			},
			component: () => import("@/views/normal/pdfFile")
		},
		{
			path: "des-report",
			name: "des-report",
			meta: {
				key: "des-report",
				pKey: "report",
				path: "des-report",
				realm: "specialReport",
				title: "导出销毁报告"
			},
			component: () => import("@/views/normal/desReport")
		}
	]
};

export default NormalRoutes;
