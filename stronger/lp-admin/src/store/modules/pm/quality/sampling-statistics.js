import { request } from "@/api/service";

const actions = {
	async QUALITY_SAMPLING_STATISTICS_GET_LIST({}, params) {
		const [startTime, endTime] = [
			new Date(`${params.date?.[0]} 00:00:00`),
			new Date(`${params.date?.[1]} 23:59:59`)
		];

		const data = {
			startTime,
			endTime,
			proCode: params.proCode,
			type: params.type
		};

		const result = await request({
			url: "pro-manager/sys-spot-check-statistic/find",
			params: data
		});

		return result;
	}
};

export default {
	actions
};
