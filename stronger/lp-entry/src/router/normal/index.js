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
			path: "doc-file",
			name: "DocFile",
			meta: {
				key: "doc-file",
				pKey: "normal",
				realm: "doc-file",
				title: "培训资料"
			},
			component: () => import("@/views/normal/docFile")
		}
	]
};

export default NormalRoutes;
