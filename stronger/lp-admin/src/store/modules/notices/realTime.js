import { request } from "@/api/service";
import moment from "moment";

const actions = {
	async GET_NOTICE_REALTIME_LIST({}, params) {
		const data = {
			pageIndex: params.pageIndex || 1,
			pageSize: params.pageSize || 10
		};
		const result = await request({
			url: "msg-manager/dingtalk-notice-msg/page",
			params: data
		});
		result.data.list = result.data.list.map(e => {
			e.CreatedAt = moment(e.CreatedAt).format("YYYY-MM-DD hh:mm:ss");
			return e;
		});
		return result;
	},
	async SEND_NOTICE_REALTIME({}, params) {
		const data = {
			groupId: params.groupId,
			msg: params.msg
		};
		const result = await request({
			url: "msg-manager/dingtalk-notice-msg/send",
			data,
			method: "POST"
		});
		return result;
	}
};

export default {
	actions
};
