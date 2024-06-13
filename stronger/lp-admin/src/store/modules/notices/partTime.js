import { request } from "@/api/service";

const actions = {
	// 常量
	async MESSAGE_GET_CONSTANTS() {
		const result = await request({
			url: "msg-manager/dingtalk-group/get-dict-const"
		});

		return result;
	},

	// 列表
	async GET_PT_MESSAGE_TABLE_LIST({}, data) {
		const params = {
			name: data.name || "",
			proCode: data.proCode || "",
			pageIndex: data.pageIndex || 1,
			pageSize: data.pageSize || 10,
			env: data.env || "",
			orderBy: data.orderBy || ""
		};
		const result = await request({
			url: "msg-manager/dingtalk-group/page",
			params
		});
		return result;
	},
	async ADD_PT_MESSAGE_TABLE_LIST({}, params) {
		const data = {
			accessToken: params.accessToken,
			env: +params.env,
			name: params.name,
			proCode: params.proCode,
			secret: params.secret
		};
		const result = await request({
			url: "msg-manager/dingtalk-group/add",
			data,
			method: "POST"
		});

		return result;
	},
	async EDIT_PT_MESSAGE_TABLE_LIST({}, params) {
		const data = {
			id: params.ID,
			accessToken: params.accessToken,
			env: +params.env,
			name: params.name,
			proCode: params.proCode,
			secret: params.secret
		};
		const result = await request({
			url: "msg-manager/dingtalk-group/edit",
			data,
			method: "POST"
		});
		return result;
	},
	async DELETE_PT_MESSAGE_TABLE_ROW({}, params) {
		const data = {
			ids: params
		};
		const result = await request({
			url: "msg-manager/dingtalk-group/delete",
			data,
			method: "DELETE"
		});
		return result;
	}
};

export default {
	actions
};
