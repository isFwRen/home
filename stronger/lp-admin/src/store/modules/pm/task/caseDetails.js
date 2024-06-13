import { request } from "@/api/service";

const actions = {
	// table 案件明细
	async GET_CASE_DETAIL_LIST({}, params) {
		const data = {
			billNum: params.billNum,
			proCode: params.proCode,
			saleChannel: params.saleChannel
		};

		const result = await request({
			url: "pro-config/CaseDetails/list",
			params: data
		});

		return result;
	},

	// dialog 分块明细
	async GET_CASE_DETAIL_BLOCK_LIST({}, body) {
		const data = {
			proCode: body.proCode,
			id: body.billId
		};

		const result = await request({
			url: "/pro-config/CaseDetails/block/list",
			params: data
		});

		return result;
	},

	// dialog 分块明细 下的字段明细
	async GET_CASE_DETAIL_FEILD_LIST({}, body) {
		const data = {
			proCode: body.proCode,
			billId: body.billId,
			blockId: body.blockId
		};

		const result = await request({
			url: "/pro-config/CaseDetails/field/list",
			params: data
		});

		return result;
	}
};

export default {
	actions
};
