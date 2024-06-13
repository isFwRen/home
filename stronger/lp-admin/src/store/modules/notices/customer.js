import Vue from "vue";
import { request } from "@/api/service";

const state = {
	globalNotice: {
		count: 0,
		proCode: "",
		fileName: "",
		sendTime: "",
		content: ""
	},
};

const getters = {
	globalNotice: () => state.globalNotice
};

const mutations = {
	GLOBAL_NOTIFICATION_UPDATE_ITEM(state, data) {
		const newData = { ...state.globalNotice, ...data };
		Vue.set(state, "globalNotice", newData);
	}
};

const actions = {
	// 列表
	async CUSTOMER_GET_LIST({}, params) {
		let [time0, time1] = params.time;

		const [startTime, endTime] = [new Date(`${time0} 00:00:00`), new Date(`${time1} 23:59:59`)];

		const data = {
			proCode: params.proCode,
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			startTime,
			endTime,
			msgType: params.msgType,
			status: params.status
		};

		const result = await request({
			url: "msg-manager/customer-notice/page",
			params: data
		});

		return result;
	},

	// 回复
	async CUSTOMER_POST_REPLY({}, body) {
		const data = {
			proCode: body.proCode,
			id: body.id,
			isReply: body.isReply,
			expectNum: +body.expectNum
		};

		const result = await request({
			method: "POST",
			url: "msg-manager/customer-notice/reply",
			data
		});

		return result;
	}
};

export default {
	state,
	getters,
	mutations,
	actions
};
