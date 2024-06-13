// 项目开发(配置管理 ——> 字段配置)
import { request } from "@/api/service";

const actions = {
	async GET_FTP_MONITOR({}, params) {
		const result = await request({
			url: "pro-config/sys-ftp-monitor/info",
			params
		});

		return result;
	},

	async EDIT_FTP_MONITOR({}, data) {
		const result = await request({
			method: "POST",
			url: "pro-config/sys-ftp-monitor/edit",
			data
		});

		return result;
	}
};

export default {
	actions
};
