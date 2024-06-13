import { request } from "@/api/service";

const actions = {
	async QUALITY_ERROR_DETAILS_GET_LIST({}, params) {
		const [startTime, endTime] = [
			new Date(`${params.date?.[0]} 00:00:00`),
			new Date(`${params.date?.[1]} 23:59:59`)
		];

		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			orderBy: params.orderBy,
			startTime,
			endTime,
			proCode: params.proCode,
			code: params.code,
			name: params.name,
			field_name: params.fieldName,
			created_code: params.createdCode,
			created_name: params.createdName
		};

		const result = await request({
			url: "pro-manager/sys-spot-check-wrong/page",
			params: data
		});

		return result;
	}
};

export default {
	actions
};
