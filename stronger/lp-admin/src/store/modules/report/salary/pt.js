import moment from "moment";
import { request } from "@/api/service";

const actions = {
	async GET_SALARY_PT_LIST({}, params) {
		console.log("params", params);
		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			pageCode: params.pageCode,
			ym: moment(params.date).format("YYYYMM"),
			code: params.code,
			name: params.name
		};

		console.log(data.ym);
		const result = await request({
			//baseURL: 'http://127.0.0.1:9999/',
			url: "report-management/pt/salary/list",
			params: data
		});
		return result;
	},

	async EXPORT_SALARY_PT_DOWNLOAD({}, ym) {
		//location.href = `${url}${ api }/report-management/pt/salary/download?params=${ ym }`

		const data = {
			ym: ym
		};
		const result = await request({
			//baseURL: 'http://127.0.0.1:9999/',
			url: "/report-management/pt/salary/download",
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
		let _fileName = ym + "PT工资表.xlsx";
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
