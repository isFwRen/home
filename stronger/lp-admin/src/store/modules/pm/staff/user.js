import moment from "moment";
import { request } from "@/api/service";

const actions = {
	// 同步
	async STAFF_SYNC_USER_LIST({ }) {
		const result = await request({
			url: "sys-user/sync"
		});

		return result;
	},

	// 导出
	async STAFF_EXPORT_PERMISSION({ }, body) {
		const data = {
			proCode: body.proCode
		};

		const result = await request({
			url: "sys-base/user-management/sys-pro-permission/export",
			params: data
		});

		return result;
	},
	async STAFF_EXPORT_EXCEL({ }, url) {
		const result = await request({
			url,
			responseType: "blob"
		});

		return result;
	},
	//导出

	// 获取用户列表
	async STAFF_GET_USER_LIST({ }, params) {
		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			code: params.code,
			name: params.name,
			phone: params.phone,
			role: params.role,
			status: params.status,
			startTime: params.date ? params.date[0] : "",
			endTime: params.date ? params.date[1] : ""
		};

		const result = await request({
			url: "sys-base/user-management/list",
			params: data
		});

		return result;
	},

	// 获取用户项目权限列表
	async STAFF_GET_USER_PERMISSION_LIST({ }, params) {
		const data = {
			userId: params.userId
		};

		const result = await request({
			url: "sys-base/user-management/user-pro-permission/list",
			params: data
		});

		return result;
	},

	// 更新用户权限
	async STAFF_UPDATE_USER_PERMISSION({ }, body) {
		const data = {
			sysProPermission: body
		};

		const result = await request({
			method: "POST",
			url: "sys-base/user-management/update-user-pro-permission",
			data
		});

		return result;
	},

	// 修改用户角色
	async STAFF_UPDATE_USER({ }, data) {
		const result = await request({
			method: "POST",
			url: "sys-base/user-management/sys-user/change-role",
			data
		});
		return result;
	}
};

export default {
	actions
};
