import moment from "moment";
import { request } from "@/api/service";
import { Store } from "vuex";

const actions = {
	async GET_SALARY_INSIDE_LIST({}, params) {
		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			ym: moment(params.date).format("YYYYMM"),
			code: params.code,
			name: params.name
		};
		let result = await request({
			//baseURL: 'http://127.0.0.1:9999/',
			url: "report-management/internal/salary/list",
			params: data
		});

		return result;
	},

	async IMPORT_SALARY_INSIDE_UPLOAD({}, body) {
		const data = {
			file: body.file
		};
		const result = await request({
			methods: "POST",
			url: "report-management/internal/salary/upload",
			body: data
		});
		return result;
	},

	async EXPORT_SALARY_INSIDE_DOWNLOAD({}, ym) {
		const data = {
			ym: ym
		};
		const result = await request({
			//baseURL: 'http://127.0.0.1:9999/',
			url: "/report-management/internal/salary/download",
			params: data,
			responseType: "blob"
		});
		console.log("re", result);
		const blob = new Blob([result], { type: "application/vnd.ms-excel;" });
		const a = document.createElement("a");
		// 生成文件路径
		let href = window.URL.createObjectURL(blob);
		a.href = href;
		//let _fileName = result.headers['content-disposition'].split(';')[1].split('=')[1].split('.')[0]
		let _fileName = ym + "内部工资表.xlsx";
		// 文件名中有中文 则对文件名进行转码
		a.download = decodeURIComponent(_fileName);
		// 利用a标签做下载
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
		window.URL.revokeObjectURL(href);
	}
};

export default {
	actions
};
