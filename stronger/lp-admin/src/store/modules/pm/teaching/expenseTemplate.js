import { request } from "@/api/service";

const actions = {
	// 列表
	async GET_PM_TEACHING_EXPENSE_TEMPLATE_LIST({}, params) {
		const data = {
			proCode: params.proCode,
			name: params.name
		};

		const result = await request({
			url: "pro-manager/reimbursementFormTemplate/list",
			params: data
		});

		return result;
	},

	// 新增
	async POST_PM_TEACHING_EXPENSE_TEMPLATE_ITEM({}, body) {
		const formData = new FormData();
		formData.append("proCode", body.proCode);

		formData.append("isRequired", body.isRequired);
		for (let file of body.files) {
			formData.append("file", file);
		}

		const result = await request({
			method: "POST",
			url: "pro-manager/reimbursementFormTemplate/add",
			data: formData
		});

		return result;
	},

	// 删除
	async DELETE_PM_TEACHING_EXPENSE_TEMPLATE_ITEM({}, body) {
		const data = {
			proCode: body.proCode,
			ids: body.ids
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/reimbursementFormTemplate/delete",
			data
		});

		return result;
	},

	// 重命名
	async RENAME_PM_TEACHING_EXPENSE_TEMPLATE_ITEM({}, body) {
		const data = {
			proCode: body.proCode,
			id: body.id,
			name: body.name
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/reimbursementFormTemplate/rename",
			data
		});

		return result;
	}
};

export default {
	actions
};
