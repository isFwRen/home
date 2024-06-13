import { request } from "@/api/service";
import { tools, sessionStorage } from "vue-rocket";

import config from "./config";
import mock_op0 from "./mock_op0";
import mock_op1 from "./mock_op1";
import mock_op2 from "./mock_op2";
import mock_opq from "./mock_opq";
import mock_modify from "./mock_modify";

const state = {
	task: {
		baseURL: "",
		topInfo: {},
		rowInfo: {},
		config: {},
		prompt: "",
		mainHeight: 0
	}
};

const getters = {
	task: () => state.task
};

const mutations = {
	UPDATE_CHANNEL(state, data) {
		const localTask = sessionStorage.get("task") || {};
		state.task = Object.assign({ ...state.task, ...localTask }, data);
		sessionStorage.set("task", state.task);
	}
};

const actions = {
	// 列表
	async GET_CHANNEL_LIST() {
		const result = await request({
			url: "data-entry/channel/list"
		});

		return result;
	},

	// 顶部信息
	async GET_LP_TASKS_INFO({}, params) {
		const data = {
			code: params.code
		};

		const result = await request({
			baseURL: state.task.baseURL,
			url: "task/opNum",
			params: data
		});

		return result;
	},

	// 理赔领单
	async GET_LP_TASK({}, params) {
		const data = {
			code: params.code,
			op: params.op
		};

		const result = await request({
			baseURL: state.task.baseURL,
			url: "task/op",
			params: data,
			headers: {
				process: params.op
			}
		});

		return result;

		// return mock_op0
	},

	// 返回重录
	async GET_LP_TASK_MODIFY({}, params) {
		const data = {
			code: params.code,
			op: params.op,
			num: params.prevNums
		};

		const result = await request({
			baseURL: state.task.baseURL,
			url: "task/modifyBlock",
			params: data,
			headers: {
				process: params.op
			}
		});

		return result;
	},

	// 配置
	async GET_LP_TASK_CONFIG() {
		const result = await request({
			baseURL: state.task.baseURL,
			url: "task/conf"
		});

		return result;

		// return config
	},

	// 提交操作后的图片
	async UPDATE_LP_TASK_FIELD_IMAGE({}, body) {
		const formData = new FormData();

		for (let key in body) {
			if (key !== "op") formData.append(key, body[key]);
		}

		const result = await request({
			method: "POST",
			baseURL: state.task.baseURL,
			url: "task/uploadImage",
			data: formData,
			headers: {
				process: body.op
			}
		});

		return result;
	},

	// 提交
	async UPDATE_LP_TASK({}, body) {
		const data = {
			bill: body.bill,
			block: body.block,
			fields: body.fields,
			op: body.op
		};

		const result = await request({
			method: "POST",
			baseURL: state.task.baseURL,
			url: "task/submit",
			data,
			headers: {
				process: body.op
			}
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
