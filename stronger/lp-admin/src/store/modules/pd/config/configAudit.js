// 项目开发(配置管理 ——> 审核配置)
import { request } from "@/api/service";
import { R } from "vue-rocket";

const actions = {
	// Add or modify
	async UPDATE_CONFIG_AUDIT({ }, form) {
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
			xmlNodeCode: form.xmlNodeCode,
			xmlNodeName: form.xmlNodeName,
			maxLen: form.maxLen,
			minLen: form.minLen,
			maxVal: form.maxVal,
			minVal: form.minVal,
			notInput: form.notInput,
			onlyInput: form.onlyInput,
			proId: form.proId,
			validation: form.validation,
			proName: form.proName,
			msg: form.msg
		};

		const result = await request({
			method: "POST",
			url: `pro-config/sys-inspection/${status === -1 ? "add" : "edit"}`,
			data
		});
		return result;
	},

	// Update input process
	async UPDATE_CONFIG_AUDIT_INPUT_PROCESS({ }, form) {
		const data = {
			id: form.id,
			proId: form.proId,
			inputProcess: form.inputProcess
		};

		const result = await request({
			method: "POST",
			url: "pro-config/sys-field/edit",
			data
		});

		return result;
	},

	// Table
	async GET_CONFIG_AUDIT_LIST({ }, params) {

		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proId: params.proId,
		};

		if (params.xmlNodeName) {
			data.xmlNodeName = params.xmlNodeName;
		}
		if (params.xmlNodeCode) {
			data.xmlNodeCode = params.xmlNodeCode;
		}
	
		const result = await request({
			url: "pro-config/sys-inspection/page",
			params: data
		});

		return result;
	},

	// Delete
	async DELETE_CONFIG_AUDIT_ITEM({ }, ids) {
		const result = await request({
			method: "DELETE",
			url: "pro-config/sys-inspection/delete",
			data: {
				ids
			}
		});
		return result;
	},

	// Get validations
	async GET_VALIDATIONS() {
		const result = await request({
			url: "pro-config/sys-inspection/validation-type"
		});
		return result;
	}
};

export default {
	actions
};
