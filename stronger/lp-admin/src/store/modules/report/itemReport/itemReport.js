import { request } from "@/api/service";
import { method } from "lodash";
import moment from "moment";

const actions = {
	// 列表
	async ITEM_REPORT_GET_LIST({}, params) {
		const data = {
			pageSize: params.pageSize,
			pageIndex: params.pageIndex,
			startTime: (params.date ? params.date[0] : moment().format("YYYY-MM-DD")) + " 00:00:00",
			endTime: (params.date ? params.date[1] : moment().format("YYYY-MM-DD")) + " 23:59:59",
			proCode: params.proCode
		};

		const result = await request({
			url: "report-management/project-report/business-details/list",
			params: data
		});

		return result;
	},

	// 导出
	async ITEM_REPORT_EXPORT_EXCEL({}, params) {
		const data = {
			startTime: (params.date ? params.date[0] : moment().format("YYYY-MM-DD")) + " 00:00:00",
			endTime: (params.date ? params.date[1] : moment().format("YYYY-MM-DD")) + " 23:59:59",
			proCode: params.proCode
		};

		const result = await request({
			url: "report-management/project-report/business-details/export",
			params: data,
			responseType: "blob"
		});

		return result;
	},
	//时效报表
	async GET_REPORT_AGING_LIST({}, params) {
		const data = {
			proCode: params.proCode ?? "all",
			type: params.type ?? 1,
			startTime: params.startTime ?? "",
			endTime: params.endTime ?? ""
		};
		const result = await request({
			url: "report-management/project-report/aging-report",
			params: data
		});
		if (result.code == 200) {
			for (let key in result.data.list) {
				result.data.list[key].map(e => {
					e.count = Number(e.count).toFixed(2) * 100;
					//e.count = Number(e.count).toFixed(2) * 100;
					return e;
				});
			}
		}
		return result;
	},
	//出处理时长报表
	async GET_REPORT_TIME_LIST({}, params) {
		const data = {
			proCode: params.proCode ?? "all",
			type: params.type ?? 1,
			startTime: params.startTime ?? "",
			endTime: params.endTime ?? ""
		};
		const result = await request({
			url: "report-management/project-report/deal-time-report",
			params: data
		});
		if (result.code == 200) {
			let data = result.data.list;
			for (let key in data) {
				data[key] = data[key].map(e => {
					e.count = Number(e.count / 3600).toFixed(2);
					return e;
				});
			}
			result.data.list = data;
		}
		return result;
	},
	//业务量报表
	async GET_REPORT_BUSINESS_LIST({}, params) {
		const data = {
			proCode: params.proCode ?? "all",
			type: params.type ?? 1,
			startTime:
				params.startTime ??
				moment().subtract(1, "month").format("YYYY-MM-DD") + "T00:00:00.000Z",
			endTime: params.endTime ?? moment().format("YYYY-MM-DD") + "T23:59:59.000Z"
		};
		const result = await request({
			url: "report-management/project-report/business-report",
			params: data
		});
		return result;
	},

	async GET_REPORT_ALL_LIST({}, params) {
		if (params.type === 3) {
			params.startTime = moment().startOf("year").format("YYYY-MM-DD") + "T23:56:59.000Z";
			params.endTime = moment().format("YYYY-MM-DD") + "T23:56:59.000Z";
		}

		const that = this._actions;
		const data = {
			proCode: params.proCode ?? "all",
			type: params.type ?? 1,
			startTime:
				params.startTime ||
				moment().subtract(1, "month").format("YYYY-MM-DD") + "T23:56:59.000Z",
			endTime: params.endTime || moment().format("YYYY-MM-DD") + "T23:56:59.000Z"
		};
		var result = {};
		(result["prescription"] = await that.GET_REPORT_AGING_LIST[0](data)),
			(result["volume"] = await that.GET_REPORT_BUSINESS_LIST[0](data)),
			(result["spendTime"] = await that.GET_REPORT_TIME_LIST[0](data));
		return result;
	},
	//导出报表
	async EXPORT_REPORT_LIST({}, params) {
		const data = {
			proCode: params.proCode ?? "all",
			type: params.type ?? 1,
			startTime:
				params.startTime ||
				moment().subtract(1, "month").format("YYYY-MM-DD") + "T00:00:00.000Z",
			endTime: params.endTime || moment().format("YYYY-MM-DD") + "T23:56:59.000Z"
		};
		const type = ["日", "周", "月", "年"];
		const result = await request({
			url: "report-management/project-report/report-export",
			params: data,
			responseType: "blob"
		});
		var anchor = document.createElement("a");
		anchor.style.display = "none";
		anchor.setAttribute("download", type[data.type - 1] + "报表.xlsx");
		anchor.href = URL.createObjectURL(result);
		document.body.appendChild(anchor);
		anchor.click();
		document.body.removeChild(anchor);
	},
	//项目报表
	async PROJECT_REPORT({}, params) {
		const data = {
			proCode: "B0118",
			startTime: moment().subtract(1, "month").format("YYYY-MM-DD") + "T00:00:00.000Z",
			endTime: moment().format("YYYY-MM-DD") + "T23:56:59.000Z",
			pageIndex: 1,
			pageSize: 10
		};
		const result = await request({
			url: "report-management/project-report/business-details/list",
			params: data
		});
	},
	// 项目报表配置
	async REPORT_SETTING({}, params) {
		const result = await request({
			url: "setting-report-tag/get-tag-list",
			params
		});
		return result;
	},
	// 获取报表表头
	async REPORT_SETTING_CELL({}, data) {
		const result = await request({
			url: "setting-report-tag/get-user-tags",
			data,
			method: "POST"
		});
		return result;
	},
	async REPORT_SETTING_SET({}, data) {
		const result = await request({
			url: "setting-report-tag/set-user-tags",
			data,
			method: "POST"
		});
		return result;
	}
};

export default {
	actions
};
