// 项目开发(配置管理 ——> 质检配置)
import { request } from "@/api/service";

const actions = {
	// Add or modify
	async UPDATE_CONFIG_QUALITY({}, form) {
		const { status } = form;

		const data = {
			proId: form.proId,
			proName: form.proName,
			id: form.id,
			// parentXmlNodeId: form.parentXmlNodeId,
			parentXmlNodeName: form.parentXmlNodeName,
			// xmlNodeId: form.xmlNodeId,
			xmlNodeName: form.xmlNodeName,
			fieldName: form.fieldName,
			fieldCode: form.fieldCode,
			inputType: form.inputType,
			belongType: form.belongType,
			billInfo: form.billInfo,
			beneficiary: form.beneficiary,
			widthPercent: +form.widthPercent,
			myOrder: +form.myOrder
		};

		const result = await request({
			method: "POST",
			url: `pro-config/sys-quality/${status === -1 ? "add" : "edit"}`,
			data
		});

		return result;
	},

	// Table
	async GET_CONFIG_QUALITY_LIST({}, params) {
		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proId: params.proId,
			parentXmlNodeName: params.parentXmlNodeName,
			xmlNodeName: params.xmlNodeName,
			fieldName: params.fieldName,
			belongType: params.belongType
		};

		const result = await request({
			url: "pro-config/sys-quality/page",
			params: data
		});

		return result;
	},

	// format
	async GET_CONFIG_QUALITY_FORMAT({}, params) {
		const data = {
			proId: params.proId
		};

		const result = await request({
			url: "/pro-config/sys-quality/format",
			params: data
		});
		return result;
	},

	// Delete
	async DELETE_CONFIG_QUALITY_ITEM({}, ids) {
		const result = await request({
			method: "DELETE",
			url: "pro-config/sys-quality/delete",
			data: {
				ids
			}
		});

		return result;
	},

	// XML节点/代码
	// async GET_CONFIG_QUALITY_XML_LIST() {
	//   const result = await request({
	//     url: ''
	//   })
	// },

	// 字段名称
	async GET_CONFIG_QUALITY_FIELD_LIST({}, params) {
		const data = {
			proId: params.proId
		};

		const result = await request({
			url: "pro-config/sys-field/list",
			params: data
		});

		return result;
	},

	// XML节点/XML代码
	async GET_CONFIG_QUALITY_XML_LIST({}, params) {
		const data = {
			proId: params.proId,
			pageIndex: 1,
			pageSize: 1000
		};

		const result = await request({
			url: "pro-config/sys-export/get-info-page",
			params: data
		});

		return result;
	},

	// 常量表(所属信息/受益人/账单/输入方式)
	async GET_CONFIG_QUALITY_CONST_LIST() {
		const result = await request({
			url: "pro-config/sys-quality/get-const"
		});

		return result;
	}
};

export default {
	actions
};
