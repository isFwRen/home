import { request } from "@/api/service";
import whitelist_list from "./whitelist_list";

const actions = {
	// 当前项目下的所有模板
	async GET_STAFF_WHITELIST_TEMP_OPTIONS({}, params) {
		const result = await request({
			url: "pro-config/sys-template/list",
			params: {
				proId: params.proId
			}
		});

		return result;
	},

	// 当前模板下的分块列表(表头，无分页)
	async GET_STAFF_WHITELIST_TEMP_CHUNK_LIST({}, params) {
		const data = {
			tempId: params.tempId
		};

		const result = await request({
			url: "pro-config/sys-block/get-info",
			params: data
		});

		return result;
	},

	// 当前模板下的分块列表的勾选人数(表头，无分页)
	async GET_STAFF_WHITELIST_Top_Sum({}, params) {
		const data = {
			proCode: params.proCode,
			tempName: params.tempCode
		};

		const result = await request({
			//baseURL: 'http://127.0.0.1:9999/',
			url: "sys-base/white-list/getBlockPeopleSum",
			params: data
		});

		return result;
	},

	// 获取白名单列表
	async GET_STAFF_WHITELIST({}, params) {
		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proCode: params.proCode,
			tempName: params.tempCode,
			code: params.code
		};

		const result = await request({
			//baseURL: 'http://127.0.0.1:9999/',
			url: "sys-base/white-list/list",
			params: data
		});

		return result;
	},

	// 复制
	async COPY_STAFF_WHITELIST_ITEM({}, body) {
		const data = {
			proCode: body.proCode,
			tempName: body.tempName,
			code: body.code,
			copyCode: body.copyCode
		};

		const result = await request({
			//baseURL: 'http://127.0.0.1:9999/',
			method: "POST",
			url: "sys-base/white-list/copy",
			data
		});

		return result;
	},

	// 导出
	async EXPORT_STAFF_WHITELIST({}, body) {
		const data = {
			proCode: body.proCode,
			tempName: body.tempCode
		};
		console.log("data", data);

		const result = await request({
			//baseURL: 'http://127.0.0.1:9999/',
			method: "POST",
			url: "sys-base/white-list/export",
			data
		});

		return result;
	},

	// 修改白名单
	async UPDATE_STAFF_WHITELIST_ITEM({}, body) {
		const data = {
			...body
		};

		const result = await request({
			//baseURL: 'http://127.0.0.1:9999/',
			method: "POST",
			url: "sys-base/white-list/edit",
			data
		});

		return result;
	}
};

export default {
	actions
};
