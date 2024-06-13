const DownloadRoutes = {
	path: "/download",
	name: "Download",
	meta: {
		title: '珠海汇流理赔2.0客户端下载网站',
		key: "download",
		path: "download"
	},
	component: () => import("@/views/download/")
};

export default DownloadRoutes;
