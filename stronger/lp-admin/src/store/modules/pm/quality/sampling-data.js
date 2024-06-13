import { request } from "@/api/service";

const actions = {
	async QUALITY_SAMPLING_DATA_GET_LIST({}, params) {
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
			code: params.code,
			name: params.name
		};

		const result = await request({
			url: "pro-manager/sys-spot-check-data/page",
			params: data
		});

		return result;
	}
};

export default {
	actions
};
