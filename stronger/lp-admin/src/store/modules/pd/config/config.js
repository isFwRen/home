// 项目开发(配置管理)
import { request } from "@/api/service";

const state = {
	config: {
		proId: "",
		pro: {}
	}
};

const getters = {
	config: () => state.config
};

const mutations = {
	UPDATE_CONFIG(state, data) {
		const { config } = state;
		state.config = Object.assign(config, data);
	}
};

const actions = {
	// 导入/导出模块
	async EXPORT_OR_IMPORT_EXPORT({ }, form) {
		const data = {
			proCode: form.proCode,
			mtype: form.mtype,
			templateId: form.templateId,
		};

		const result = await request({
			method: "POST",
			url: "pro-config/sys-project/move-conf",
			data
		});

		return result;
	},
	// Add or modify
	async UPDATE_CONFIG_ITEM({ }, body) {
		const { status } = body;

		const data = {
			id: body.id,
			editVersion: body.editVersion,
			name: body.name,
			code: body.code,
			type: body.type,
			cacheTime: +body.cacheTime,
			saveDate: +body.saveDate,
			restartAt: body.restartAt,
			autoReturn: body.autoReturn
		};

		const result = await request({
			method: "POST",
			url: `pro-config/sys-project/${status === -1 ? "add" : "edit"}`,
			data
		});

		return result;
	},

	// Table
	async GET_CONFIG_ITEM_LIST({ }) {
		const result = await request({
			url: "pro-config/sys-project/list"
		});
		return result;
	},

	// 添加新模板
	async ADD_CONFIG_ITEM_TEMPLATE({ }, form) {
		const data = {
			proId: form.proId,
			name: form.name
		};

		const result = await request({
			method: "POST",
			url: "pro-config/sys-template/add",
			data
		});

		return result;
	},

	// 设置 - 刷新录入配置
	async REFRESH_CONFIG_PRO_CONFIG({ }, body) {
		const ipPort = body.isIntranet
			? `${body.innerIp}:${body.inAppPort}`
			: `${body.outIp}:${body.outAppPort}`;

		const result = await request({
			method: "POST",
			url: `https://${ipPort}/api/task/conf/refresh-pro-conf`
		});

		return result;
	},

	// 设置 - 刷新管理配置
	async REFRESH_CONFIG_PRO_MANAGE_CONFIG({ }, body) {
		const result = await request({
			method: "POST",
			url: `/pro-config/refresh-pro-conf`,
			data: {
				proCode: body.proCode
			}
		});

		return result;
	}
};

export default {
	state,
	getters,
	mutations,
	actions
};
