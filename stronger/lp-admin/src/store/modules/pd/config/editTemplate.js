// 项目开发(配置管理 ——> 修改模板)
import { request } from "@/api/service";

const state = {
	editTemp: {
		tempFieldList: [],
		frontChunkList: [],
		onCodeFrontChunkList: [],
		suggestChunkList: [],
		fieldList: [],
		chunkList: [],
		chunkId: ""
	}
};

const getters = {
	editTemp: () => state.editTemp
};

const mutations = {
	UPDATE_EDIT_TEMP(state, data) {
		state.editTemp = Object.assign(state.editTemp, data);

		console.log(state.editTemp);
	}
};

const actions = {
	// 当前项目下的所有模板
	async GET_CONFIG_TEMP_OPTIONS({}, params) {
		const result = await request({
			url: "pro-config/sys-template/list",
			params: {
				proId: params.proId
			}
		});

		return result;
	},

	// 复制模板
	async COPY_CONFIG_TEMP({}, form) {
		const data = {
			proId: form.proId,
			id: form.id,
			name: form.name
		};

		const result = await request({
			method: "POST",
			url: "pro-config/sys-template/copy-temp",
			data
		});

		return result;
	},

	// Add or modify(分块)
	async UPDATE_CONFIG_TEMP_CHUNK({}, form) {
		const { status } = form;

		const data = {
			proTempId: form.tempId,
			id: form.id,
			myOrder: +form.myOrder,
			code: form.code,
			name: form.name,
			fEight: form.fEight,
			ocr: form.ocr,
			freeTime: +form.freeTime,
			relation: form.relation,
			isCompetitive: form.isCompetitive,
			isLoop: form.isLoop,
			isMobile: form.isMobile
		};

		const result = await request({
			method: "POST",
			url: `pro-config/sys-block/${status === -1 ? "add" : "edit"}`,
			data
		});

		return result;
	},

	// 当前模板下的分块列表(无分页)
	async GET_CONFIG_TEMP_CHUNK_LIST({}, params) {
		const data = {
			tempId: params.tempId,
			blockName: params.blockName
		};

		const result = await request({
			url: "pro-config/sys-block/get-info",
			params: data
		});

		return result;
	},

	//  Delete
	async DEL_CONFIG_TEMP_CHUNK({}, ids) {
		const result = await request({
			method: "DELETE",
			url: "pro-config/sys-block/delete",
			data: {
				ids
			}
		});
		return result;
	},

	// Exchange
	async EXCHANGE_CONFIG_TEMP_CHUNK({}, form) {
		const data = {
			startId: form.startId,
			startOrder: +form.startOrder,
			endId: form.endId,
			endOrder: +form.endOrder
		};

		const result = await request({
			method: "POST",
			url: "pro-config/sys-block/change-order",
			data
		});

		return result;
	},

	// 字段配置
	async GET_CONFIG_TEMP_CHUNK_FIELD_LIST({}, params) {
		const data = {
			id: params.id
		};

		const result = await request({
			url: `pro-config/sys-block-relation/get-info/${data.id}`,
			params: data
		});

		return result;
	},

	// 当前项目下的所有字段
	async GET_CONFIG_ALL_FIELD_LIST({}, params) {
		const data = {
			proId: params.proId
		};

		const result = await request({
			url: "pro-config/sys-field/list",
			params: data
		});

		return result;
	},

	// 更新新的分块字段关系
	async UPDATE_CHUNK_FIELD_RELATION({}, form) {
		const data = form;

		const result = await request({
			method: "POST",
			url: "pro-config/sys-block-field-relation/delete-and-add-all",
			data
		});

		return result;
	},

	// 插入新的分块关系
	async UPDATE_CHUNK_CHUNK_RELATION({}, form) {
		const data = form;

		const result = await request({
			method: "POST",
			url: "pro-config/sys-block-relation/delete-and-add-all",
			data
		});

		return result;
	},

	// 更新模板图片
	async UPDATE_CONFIG_TEMP_IMAGES({}, form) {
		const data = {};

		const result = await request({
			method: "POST",
			url: "pro-config/sys-template/update-images",
			data
		});

		return result;
	},

	// 分块截图配置
	async UPDATE_CONFIG_TEMP_CHUNK_SCREENSHOT({}, form) {
		const data = {
			blockId: form.blockId,
			coordinate: form.coordinate,
			coordinateType: form.coordinateType,
			picPage: form.picPage
		};

		const result = await request({
			method: "POST",
			url: "pro-config/sys-block/edit-crop-coordinate",
			data
		});

		return result;
	},
};

export default {
	state,
	getters,
	mutations,
	actions
};
