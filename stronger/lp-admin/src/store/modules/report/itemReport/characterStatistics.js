import { request } from "@/api/service";
import moment from "moment";

const actions = {
	//字符统计列表
	async GET_CHARACTER_STATISTICS_LIST({}, params) {
		const data = {
			startTime:
				(params.date ? params.date[0] : moment().format("YYYY-MM-DD")) + "T00:00:00.000Z",
			endTime:
				(params.date ? params.date[1] : moment().format("YYYY-MM-DD")) + "T23:59:59.000Z",
			proCode: params.proCode,
			pageIndex: params.pageIndex,
			pageSize: params.pageSize
		};
		const result = await request({
			url: "report-management/project-report/get-char-sum",
			params: data
		});
		if (result.code == 200) {
			result.data.list = result.data.list.map(e => {
				e.averageCharCount = Number(e.averageCharCount).toFixed(2);
				e.averageInputCharCount = Number(e.averageInputCharCount).toFixed(2);
				e.charPercent = Number(e.charPercent * 100).toFixed(2) + "%";
				e.settleStaffPercent = Number(e.settleStaffPercent * 100).toFixed(2) + "%";
				return e;
			});
		}
		console.log(result);
		return result;
	},

	//导出数据
	async EXPORT_CHARACTER_STATISTICS_LIST({}, params) {
		const data = {
			startTime:
				(params.date ? params.date[0] : moment().format("YYYY-MM-DD")) + "T00:00:00.000Z",
			endTime:
				(params.date ? params.date[1] : moment().format("YYYY-MM-DD")) + "T23:59:59.000Z",
			proCode: params.proCode
		};

		const result = await request({
			url: "report-management/project-report/export-char-sum",
			params: data,
			responseType: "blob"
		});

		var anchor = document.createElement("a");
		anchor.style.display = "none";
		anchor.setAttribute("download", data.proCode + "字符统计表.xlsx");
		anchor.href = URL.createObjectURL(result);
		document.body.appendChild(anchor);
		anchor.click();
		document.body.removeChild(anchor);
	}
};

export default {
	actions
};
