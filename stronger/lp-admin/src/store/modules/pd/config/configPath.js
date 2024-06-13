// 项目开发(配置管理 ——> 路径配置)
import { request } from "@/api/service";

const actions = {
	// 路径配置
	async GET_CONFIG_PATH_PATH_LIST({}, params) {
		const data = {
			proId: params.proId
		};

		const result = await request({
			url: "pro-config/sys-pro-download-paths/list",
			params: data
		});

		return result;
	},

	// 更改下载路径
	async UPDATE_CONFIG_PATH_ITEM({}, data) {
		const result = await request({
			method: "POST",
			url: "pro-config/sys-pro-download-paths/set-available",
			data: {
				sysProPathsList: data
			}
		});

		return result;
	},

	// 程序开启情况
	async GET_CONFIG_PATH_ENABLE_LIST({}, params) {
		const data = {
			proCode: params.proCode
		};

		const result = await request({
			url: "pro-config/sys-pro-download-paths/process/list",
			params: data
		});

		return result;
	}
};

export default {
	actions
};
