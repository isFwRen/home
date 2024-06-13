// 项目开发(配置管理 ——> 导出配置)
import { request } from "@/api/service";
import { tools as lpTools } from "@/libs/util";

const { baseURL, baseURLApi } = lpTools.baseURL();

const actions = {
	// Update export template
	async UPDATE_CONFIG_EXPORT_TEMPLATE({ }, form) {
		const data = {
			xmlType: form.xmlType,
			proId: form.proId,
			id: form.id,
			proName: form.proName,
			tempVal: form.tempVal
		};

		const result = await request({
			method: "POST",
			url: "pro-config/sys-export/edit",
			data
		});

		return result;
	},

	// Add or modify
	async UPDATE_CONFIG_EXPORT({ }, body) {
		const { status } = body;
		let wink = {};

		if (status === 1) {
			wink = {
				id: body.id
			};
		}

		const data = {
			...wink,
			myOrder: +body.myOrder,
			exportId: body.exportId,
			proId: body.proId,
			name: body.name,
			oneFields: body.oneFields,
			twoFields: body.twoFields,
			threeFields: body.threeFields,
			oneFieldsName: body.oneFieldsName,
			twoFieldsName: body.twoFieldsName,
			threeFieldsName: body.threeFieldsName,
			myType: body.myType,
			fixedValue: body.fixedValue,
			remark: body.remark
		};

		const result = await request({
			method: "POST",
			url: `pro-config/sys-export/${status === -1 ? "add" : "edit"}-node`,
			data
		});

		return result;
	},

	// Table
	async GET_CONFIG_EXPORT_LIST({ }, params) {
		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proId: params.proId,
			xmlType: params.xmlType,
			name: params.name,
			fieldLike: params.fieldLike
		};

		const result = await request({
			url: "pro-config/sys-export/get-info-page",
			params: data
		});

		return result;
	},

	// Delete
	async DELETE_CONFIG_EXPORT_ITEM({ }, ids) {
		const result = await request({
			method: "DELETE",
			url: "pro-config/sys-export/delete",
			data: {
				ids
			}
		});
		return result;
	},

	// Insert
	async INSERT_CONFIG_EXPORT_ITEM({ }, form) {
		const data = {
			startOrder: +form.startOrder,
			endOrder: +form.endOrder,
			startId: form.startId,
			exportId: form.exportId
		};

		const result = await request({
			method: "POST",
			url: "pro-config/sys-export/change-order",
			data
		});

		return result;
	},

	// Export
	async EXPORT_CONFIG_EXPORT({ }, id) {
		console.log(id, 'id')
		const result = await request({
			url: `${baseURLApi}pro-config/sys-export/export?exportId=${id}`,
			responseType: "blob"
		});

		return result;
	},

	// Fields
	async GET_CONFIG_EXPORT_FIELDS({ }, form) {
		const data = {
			proId: form.proId
		};

		const result = await request({
			url: "pro-config/sys-field/list",
			params: data
		});

		if (result.code !== 200) return [];

		const fieldsList = [];

		for (let item of result.data) {
			fieldsList.push({
				label: item.name,
				value: item.code
			});
		}

		return fieldsList;
	},
};

export default {
	actions
};
