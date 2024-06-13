import { request } from "@/api/service";

const actions = {
	// 获取权限统计列表
	async GET_STAFF_PERMISSION_LIST({}, params) {
		// const data = {
		//   pageIndex: params.pageIndex,
		//   pageSize: params.pageSize,
		//   proCode: params.proCode,
		//   startTime: params.startTime,
		//   endTime: params.endTime,
		//   code: params.code,
		//   isCheckAll: params.isCheckAll
		// }

		const result = await request({
			url: "sys-base/pro-permission-check/list"
		});

		return result;
	}
};

export default {
	actions
};
