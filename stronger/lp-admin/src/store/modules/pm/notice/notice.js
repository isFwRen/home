import { request } from "@/api/service";
import { tools } from "@/libs/util.tools.js";

const actions = {
	async GET_PM_NOTICE_LIST({}, params) {
		const data = {
			startTime:
				(params.dayRange && new Date(params.dayRange[0] + " 00:00:00")) ||
				new Date("1990-01-01"),
			endTime: (params.dayRange && new Date(params.dayRange[1] + " 23:59:59")) || new Date(),
			proCode: params.proCode || "",
			releaseType: params.releaseType || "",
			status: params.status || "",
			pageIndex: params.pageIndex || 0,
			pageSize: params.pageSize || 10,
			orderBy: params.orderBy || ""
		};

		const result = await request({
			url: "pro-manager/announcement-manager/page",
			params: data
		});

		return result;
	},
	async CHANGE_PM_NOTICE_LIST_ITEM_STATUS({}, body) {
		const data = {
			ids: body.ids,
			status: body.status
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/announcement-manager/change-status",
			data
		});

		return result;
	},

	async UPDATE_PM_NOTICE_LIST_ITEM({}, body) {
		const data = {
			id: body.ID,
			content: body.content || "",
			createdAt: body.createdAt,
			proCode: body.proCode,
			releaseDate: body.releaseDate,
			releaseType: +body.releaseType,
			releaseUserCode: body.releaseUserCode,
			releaseUserName: body.releaseUserName,
			status: body.status,
			title: body.title,
			updatedAt: body.updatedAt
		};

		const result = await request({
			method: "POST",
			url: "pro-manager/announcement-manager/edit",
			data
		});

		return result;
	},
	async ADD_A_NOTICE_ITEM({}, body) {
		const data = {
			content: body.content || "",
			proCode: body.proCode,
			releaseType: +body.releaseType,
			title: body.title
		};
		const result = await request({
			method: "POST",
			url: "pro-manager/announcement-manager/add",
			data
		});

		return result;
	},
	async NOTIC_EDIT_FILE_UPLOAD({}, body) {
		const formData = new FormData();
		formData.append("file", body.file);
		const result = await request({
			method: "POST",
			url: "file-upload-and-download/upload",
			data: formData,
			headers: {
				"Content-Type": "multipart/form-data"
			}
		});
		result.data.file.url = tools.baseURL().baseURLApi + result.data.file.url;
		return result;
	}
};

export default {
	actions
};
