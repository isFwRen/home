import { request } from "@/api/service";

const actions = {
	async GET_NOTIC_BY_PRO_LIST({}, params) {
		const data = {
			groupId: params.proCode
		};
		const result = await request({
			url: "msg-manager/group-notice/get-by-group-id",
			params: data
		});
		return result;
	},
	//按照固定间隔在
	async ADD_NOTIC_NEW_REGULAR_BYTIME({}, params) {
		const data = {
			twos: params.twos || [],
			ones: params.ones || [],
			type: params.type
		};
		const result = await request({
			url: "msg-manager/group-notice/add",
			data,
			method: "POST"
		});
		return result;
	},
	//重置
	async RESET_NOTIC_REGULAR({}, params) {
		const data = {
			type: params.type,
			groupId: params.groupId
		};
		const result = await request({
			url: "msg-manager/group-notice/re",
			params: data,
			method: "POST"
		});
		return result;
	},
	//编辑
	async EDIT_NOTIC_NEW_REGULAR_BYTIME({}, { type, twos, ones }) {
		const data = {
			ones,
			twos,
			type
		};
		const result = await request({
			url: "msg-manager/group-notice/edit",
			data,
			method: "POST"
		});
		return result;
	}
};

export default {
	actions
};
