// 项目开发(配置管理 ——> 字段配置)
import { request } from "@/api/service";

const actions = {
	// Add or modify
	async UPDATE_CONFIG_FIELD({}, form) {
		const { status } = form;
		let wink = {};

		if (status === 1) {
			wink = {
				id: form.id
			};
		}

		const data = {
			...wink,
			myOrder: form.myOrder,
			proId: form.proId,
			code: form.code,
			name: form.name,
			fixValue: form.fixValue,
			specChar: form.specChar,
			inputProcess: form.inputProcess,
			checkDate: form.checkDate,
			questionChange: form.questionChange,
			valChange: form.valChange,
			valInsert: form.valInsert,
			ignoreIf: form.ignoreIf,
			prompt: form.prompt,

			maxLen: form.maxLen,
			minLen: form.minLen,
			fixLen: form.fixLen,
			maxVal: form.maxVal,
			minVal: form.minVal,

			defaultVal: form.defaultVal,
			validations: form.validations
		};

		const result = await request({
			method: "POST",
			url: `pro-config/sys-field/${status === -1 ? "add" : "edit"}`,
			data
		});

		return result;
	},

	// Update input process
	async UPDATE_CONFIG_FIELD_INPUT_PROCESS({}, body) {
		const data = {
			...body,
			id: body.id,
			inputProcess: body.inputProcess
		};

		const result = await request({
			method: "POST",
			url: "pro-config/sys-field/edit",
			data
		});

		return result;
	},

	// Table
	async GET_CONFIG_FIELD_LIST({}, params) {
		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proId: params.proId,
			name: params.name,
			inputProcess: params.inputProcess
		};

		const result = await request({
			url: "pro-config/sys-field/page",
			params: data
		});
		return result;
	},

	// Delete
	async DELETE_CONFIG_FIELD_ITEM({}, ids) {
		const result = await request({
			method: "DELETE",
			url: "pro-config/sys-field/delete",
			data: {
				ids
			}
		});
		return result;
	},

	// Get const
	async GET_CONFIG_CONSTANT() {
		const result = await request({
			url: "pro-config/sys-field/get-const"
		});
		return result;
	},

	// Issue config
	async ADD_ISSUE_CONFIG({}, list) {
		const data = {
			...list
		};

		const result = await request({
			method: "POST",
			url: "pro-config/sys-issue/edit",
			data
		});

		return result;
	},

	// Issue config list
	async GET_ISSUE_CONFIG_LIST({}, params) {
		const data = {
			...params
		};

		const result = await request({
			url: "pro-config/sys-issue/list",
			params: data
		});

		return result;
	},

	// export config list
	async GET_EXPORT_CONFIG_LIST({}, params) {
		const data = {
			...params
		};

		const result = await request({
			url: "pro-config/sys-field-check/list",
			params: data
		});

		return result;
	},

	// export config
	async ADD_EXPORT_CONFIG({}, list) {
		const data = {
			...list
		};

		const result = await request({
			method: "POST",
			url: "pro-config/sys-field-check/edit",
			data
		});

		return result;
	},
};

export default {
	actions
};
