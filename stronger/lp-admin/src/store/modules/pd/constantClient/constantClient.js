import { request } from "@/api/service";
import { tools } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";

import PocketDB from "@/api/pocketdb";

const { baseURL: constBaseURL } = lpTools.constBaseURL();

const actions = {
	// 上传常量文件
	//*客户端*
	async UPLOAD_CONSTANT_FILE_CLIENT({}, body) {
		const { proCode, formData } = body;
		if (!proCode) {
			return {
				msg: "参数错误",
				code: 502
			};
		}
		const result = await request({
			url: `sys-const/import/${proCode}`,
			method: "post",
			data: formData
		});
		return result;
	},

	//*客户端*

	// 常量表信息
	//*客户端*
	async GET_CONSTANT_OPTIONS_CLIENT({}, proCode) {
		const result = await request({
			url: `sys-const/info-list/${proCode}`,
			methods: "get"
		});
		return result;
	},
	//*客户端*

	//*客户端*
	// 获取表头
	async GET_CONSTANT_HEADERS_CLIENT({}, data) {
		const params = {
			pageSize: 10,
			pageIndex: 1
		};
		const body = Object.assign(params, data);
		const result = await request({
			method: "post",
			url: "sys-const/page",
			data: body
		});

		return result;
	},
	//*客户端*

	// 更新当前常量表
	//**客户端 */
	async UPDATE_CONSTANT_EXCEL_ITEM_CLIENT({}, body) {
		const result = await request({
			method: "post",
			url: "sys-const/edit",
			data: body
		});
		return result;
	},
	//**客户端 */

	// 导出常量表
	//*客户端*
	async EXPORT_CONSTANT_EXCEL_CLIENT({}, { proCode, dbName }) {
		const data = {
			proCode,
			name: dbName
		};

		const result = await request({
			url: "sys-const/export",
			method: "post",
			data,
			responseType: "blob"
		});

		return result;
	},

	//*客户端*

	//*客户端*
	async DESTROY_CONSTANT_EXCEL_CLIENT({}, data) {
		const result = await request({
			method: "POST",
			url: "sys-const/del-tables",
			data: data
		});
		return result;
	},
	//*客户端*

	// 批量删除(删除当前常量表下的数据)
	//*客户端*
	async DELETE_CONSTANT_EXCEL_ITEM_CLIENT({}, data) {
		const result = await request({
			method: "POST",
			url: "sys-const/del-lines",
			data
		});
		return result;
	},
	//*客户端*

	// 批量添加常量
	//**客户端 */
	async BATCH_ADD_EXCEL_CLIENT({}, data) {
		const result = await request({
			method: "post",
			url: "sys-const/insert",
			data
		});
		return result;
	},
	//**客户端 */

	async CONSTANT_PUBLISH_CLIENT({}, data) {
		const { proCode } = data;

		const result = await request({
			method: "patch",
			url: `sys-const/release/${proCode}`,
			data
		});
		return result;
	},

	async CONSTANT_PUBLISH_CLIENT({}, data) {
		const { proCode } = data;

		const result = await request({
			method: "patch",
			url: `sys-const/release/${proCode}`,
			data
		});
		return result;
	},

	async CONSTANT_OPERATION_LOG({}, data) {
		const body = {
			pageIndex: data.pageIndex,
			pageSize: data.pageSize
		};
		if (data.range && data.range.length > 0) {
			body.startTime = data.range && data.range[0] + "T00:00:00.000Z";
			body.endTime = data.range && data.range[1] + "T23:59:59.000Z";
		}
		if (data.type) {
			body.type = data.type;
		}

		const result = await request({
			method: "post",
			url: "pro-config/const/operation-log/page",
			data: body
		});
		return result;
	}
};

export default {
	actions
};
