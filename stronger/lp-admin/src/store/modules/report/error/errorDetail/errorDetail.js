import { request } from "@/api/service";

const actions = {
	// 错误明细列表
	async GET_REPORT_ERROR_DETAIL_LIST_ITEM({}, params) {
		const data = {
			pageSize: params.pageSize,
			pageIndex: params.pageIndex,
			proCode: params.proCode,
			startTime: `${params.date[0]} 00:00:00`,
			endTime: `${params.date[1]} 23:59:59`,
			code: params.code,
			nickName: params.nickName,
			fieldName: params.fieldName,
			op: params.op,
			complaint: params.complaint,
			confirm: params.confirm,
			isAudit: params.isAudit
		};

		const result = await request({
			url: "report-management/error-statistics/list",
			params: data
		});

		return result;
	},

	// 申诉
	async APPEAL_REPORT_ERROR_DETAIL_LIST_ITEM({}, body) {
		const data = {
			proCode: body.proCode,
			id: body.id,
			complainConfirm: body.complainConfirm
		};

		const result = await request({
			method: "POST",
			url: "report-management/error-statistics/complain",
			data
		});

		return result;
	},

	// 差错审核
	async REVIEW_REPORT_ERROR_DETAIL_LIST_ITEM({}, body) {
		const data = {
			proCode: body.proCode,
			// right: body.right,
			list: body.list,
			wrongConfirm: body.wrongConfirm
		};

		const result = await request({
			method: "POST",
			url: "report-management/error-statistics/wrong-confirm",
			data
		});
		return result;
	},

	// 导出
	async EXPORT_ERROR_DETAIL_LIST_ITEM({}, body) {
		const data = {
			endTime: body.endTime + "T12:00:00.000Z",
			proCode: body.proCode,
			startTime: body.startTime + "T00:00:00.000Z"
		};

		const result = await request({
			method: "POST",
			url: "report-management/error-statistics/export",
			data,
			responseType: 'blob'
		});
		return result;
	}
};

export default {
	actions
};
