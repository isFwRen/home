import { request } from "@/api/service";

const actions = {
	// 列表
	async QUALITY_SAMPLING_SETTING_GET_LIST({}, params) {
		const [startTime, endTime] = [
			new Date(`${params.date?.[0]} 00:00:00`),
			new Date(`${params.date?.[1]} 23:59:59`)
		];

		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			startTime,
			endTime,
			type: params.type,
			proCode: params.proCode,
			status: params.status
		};

		const result = await request({
			url: "pro-manager/sys-spot-check/page",
			params: data
		});

		return result;
	},

	// 更新
	async QUALITY_SAMPLING_SETTING_UPDATE_ITEM({}, body) {
		const data = {
			id: body.id,
			proCode: body.proCode,
			ratio: body.ratio,
			code: body.code
			// name: body.name
		};

		const result = await request({
			method: "POST",
			url: `pro-manager/sys-spot-check/${body.id ? "edit" : "add"}`,
			data
		});

		return result;
	}
};

export default {
	actions
};
