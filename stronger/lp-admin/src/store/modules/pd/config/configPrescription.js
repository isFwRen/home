// 项目开发(配置管理 ——> 时效配置)
import { request } from "@/api/service";

const actions = {
	// Add or modify
	async UPDATE_CONFIG_PRESCRIPTION({ }, form) {
		const { status } = form;

		const data = {
			proId: form.proId,
			id: form.id,
			configType: form.configType,
			agingStartTime: form.agingStartTime,
			agingEndTime: form.agingEndTime,
			agingOutStartTime: form.agingOutStartTime,
			agingOutEndTime: form.agingOutEndTime,
			nodeName: form.nodeName,
			nodeContent: form.nodeContent,
			fieldName: form.fieldName,
			fieldContent: form.fieldContent,
			requirementsTime: form.requirementsTime
		};

		const result = await request({
			method: "POST",
			url: `pro-config/project-config-aging/${status === -1 ? "add" : "edit"}`,
			data
		});

		return result;
	},

	// Table
	async GET_CONFIG_PRESCRIPTION_LIST({ }, params) {
		const data = {
			proId: params.proId,
			configType: params.configType
		};

		const result = await request({
			url: "pro-config/project-config-aging/list",
			params: data
		});

		return result;
	},

	// Delete
	async DELETE_CONFIG_PRESCRIPTION_ITEM({ }, ids) {
		const result = await request({
			method: "DELETE",
			url: "pro-config/project-config-aging/delete",
			data: {
				ids
			}
		});
		return result;
	},

	// Add or modify(节假日设置)
	async UPDATE_CONFIG_PRESCRIPTION_HOLIDAY({ }, form) {
		const { id } = form;

		const data = {
			...form
		};

		const result = await request({
			method: "POST",
			url: `pro-config/project-config-aging-holiday/${!id ? "add" : "edit"}`,
			data
		});

		return result;
	},

	// Holiday
	async GET_CONFIG_PRESCRIPTION_HOLIDAY({ }, params) {
		const data = {
			InquireStartDate: params.startDate,
			InquireEndDate: params.endDate
		};

		const result = await request({
			url: "pro-config/project-config-aging-holiday/list",
			params: data
		});

		return result;
	},

	// Contract
	async GET_CONFIG_CONTRACT_LIST({ }, params) {
		const result = await request({
			url: "pro-config/project-config-aging-contract/list",
			params
		});

		return result;
	},
	async ADD_CONFIG_CONTRACT({ }, data) {
		const result = await request({
			url: "pro-config/project-config-aging-contract/add",
			data,
			method: 'post'
		});

		return result;
	},
	async EDIT_CONFIG_CONTRACT({ }, data) {
		const result = await request({
			url: "pro-config/project-config-aging-contract/edit",
			data,
			method: 'post'
		});

		return result;
	},
	async DELETE_CONFIG_CONTRACT({ }, data) {
		const result = await request({
			url: "pro-config/project-config-aging-contract/delete",
			data,
			method: 'delete'
		});

		return result;
	},

};

export default {
	actions
};
