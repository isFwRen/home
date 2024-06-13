import { request } from "@/api/service";

const actions = {
	// 查看日志
	async GET_CASE_LOGS({}, form) {
		const data = {
			proCode: form.proCode,
			id: form.id
		};

		const result = await request({
			url: "pro-manager/bill-list/get-log",
			params: data
		});

		return result;
	},

	// 查询/修改录入数据
	async GET_OR_UPDATE_CASE_ENTRY_DATA({}, body) {
		const data = {
			editType: body.editType,
			billId: body.billId,
			fields: body.fields,
			editFields: body.modifiedFields,
			proCode: body.proCode
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/bill-list/edit-bill-result-data",
			data
		});

		return result;
	},
	//
	async GET_OR_UPDATE_CASE_ENTRY_DATAS({}, body) {
		const data = {
			editType: body.editType,
			billId: body.billId,
			fields: body.fields,
			editFields: body.modifiedFields,
			proCode: body.proCode
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/bill-list/zbj-bill-result-data",
			data
		});

		return result;
	},

	// 获取字段、分块和字段配置信息
	async GET_CASE_FIELD_INFO({}, body) {
		const data = {
			proCode: body.proCode,
			id: body.fieldId
		};

		const result = await request({
			url: "pro-manager/bill-list/get-field-info",
			params: data
		});

		return result;
	},

	// 设置分块是否练习
	async UPDATE_CASE_BLOCK_PRACTICE({}, form) {
		const data = {
			blockIds: form.blockIds,
			isPractice: form.isPractice,
			proCode: form.proCode
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/bill-list/set-practice",
			data
		});

		return result;
	},

	// 查看图片
	async GET_CASE_FIELD_IMAGE({}, params) {
		const data = {
			proCode: params.proCode,
			id: params.blockId,
			fieldId: params.fieldId
		};

		const result = await request({
			url: "pro-manager/bill-list/get-block-img",
			params: data
		});

		return result;
	},

	// 填写正确数据
	async CASE_POST_FEEDBACK_VALUE({}, body) {
		const data = {
			fieldValue: body.fieldValue,
			editDate: body.editDate,
			month: body.month,
			responsibleCode: body.responsibleCode,
			responsibleName: body.responsibleName,
			op: body.op
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/bill-list/edit-feedback-val",
			data
		});

		return result;
	},

	// 保存为xml
	async SAVE_CASE_ENTRY_DATA_AS_XML({}, body) {
		const data = {
			editType: body.editType,
			proCode: body.proCode,
			billId: body.billId,
			fields: body.fields,
			editFields: body.modifiedFields
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/bill-list/save-bill-result-data",
			data
		});

		return result;
	}
};

export default {
	actions
};
