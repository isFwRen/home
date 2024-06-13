import { request } from "@/api/service";
import moment from "moment";
import { R } from "vue-rocket";
const TODAY = moment().format("YYYY-MM-DD");

const actions = {
	// 获取回传分析
	async GET_RETRURN_ANALYSIS_LIST({}, params) {
		const proCodeSplit = params.proCode ? params.proCode.split("/") : ["全部", "整体"];
		const proCodeFront = proCodeSplit[0];
		const proCodeBehind = proCodeSplit[1];
		const startTime = R.isYummy(params.time) ? params.time[0] : TODAY;
		const endTime = R.isYummy(params.time) ? params.time[1] : TODAY;
		const isCheckAll = proCodeBehind === "整体" ? true : false;

		const data = {
			proCode: proCodeFront,
			startTime,
			endTime,
			types: "upload",
			isCheckAll
		};

		const result = await request({
			url: "report-management/business-analysis/upload/list",
			params: data
		});

		return result;
	},

	// 导出
	async EXPORT_RETURNANALYSIS_DETAIL({}, param) {
		var IsCheckAll = false;
		if (param.proCode.split("/")[1] == "整体") {
			IsCheckAll = true;
		}

		const data = {
			startTime: param.time ? param.time[0] : moment().format("yyyy-MM-DD"),
			endTime: param.time ? param.time[1] : moment().format("yyyy-MM-DD"),
			proCode: param.proCode.split("/")[0],
			types: "upload",
			isCheckAll: IsCheckAll
		};

		const result = await request({
			url: "/report-management/business-analysis/upload/export",
			params: data,
			responseType: "blob"
		});

		return result;
	}
};

export default {
	actions
};
