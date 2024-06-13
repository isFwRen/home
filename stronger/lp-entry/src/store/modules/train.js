import { request } from "@/api/service";

const state = {
	trainSteps: {
		1: false,
		2: false,
		3: false,
		4: false,
		5: false,
		6: false
	}
};
const getters = {};

const mutations = {
	SET_TRAINSTEPS(state, { key, value }) {
		state.trainSteps[key] = value;
	}
};

const actions = {
	async GET_DOC_FILE({}, params) {
		const result = await request({
			url: "/files/common/pxzy.pdf",
			responseType: "blob"
		});
		return result;
	},
	// 培训指引
	async TRAINING_GUIDE({}, data) {
		const reuslt = await request({
			url: "/training-guide/training-stage/find",
			method: "POST"
		});
		return reuslt;
	},
	async UPDATE_GUIDE_STAGE({}, params) {
		const reuslt = await request({
			url: "/training-guide/finish-read-doc",
			method: "GET"
		});
		return reuslt;
	},
	// 完成一个文件学习
	async FILISH_FILE_READ({}, data) {
		const reuslt = await request({
			url: "/training-guide/finish-read",
			method: "POST",
			data
		});
		return reuslt;
	}
};

export default {
	state,
	getters,
	mutations,
	actions
};
