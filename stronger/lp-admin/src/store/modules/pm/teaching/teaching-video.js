import { request } from "@/api/service";

const actions = {
	// 列表
	async GET_PM_TEACHING_TEACHING_VIDEO_LIST({}, params) {
		const data = {
			pageSize: params.pageSize,
			pageIndex: params.pageIndex,
			proCode: params.proCode,
			blockName: params.blockName,
			rule: params.rule
		};

		const result = await request({
			url: "pro-manager/teachVideo/list",
			params: data
		});

		return result;
	},

	// 删除
	async DELETE_PM_TEACHING_TEACHING_VIDEO_ITEMS({}, ids) {
		const result = await request({
			method: "POST",
			url: "pro-manager/teachVideo/delete",
			data: {
				ids
			}
		});

		return result;
	},

	// 导出
	async EXPORT_PM_TEACHING_TEACHING_VIDEO({}, body) {
		const data = {
			proCode: body.proCode
		};

		const result = await request({
			url: "pro-manager/teachVideo/export",
			params: data
		});

		return result;
	},

	// 视频上传
	async UPLOAD_PM_TEACHING_TEACHING_VIDEO_VIDEOS({}, body) {
		const formData = new FormData();

		if (body.id) {
			for (let key in body) {
				formData.append(key, body[key]);
			}
		} else {
			for (let key in body) {
				if (key !== "file") {
					formData.append(key, body[key]);
				} else {
					for (let file of body[key]) {
						formData.append(key, file);
					}
				}
			}
		}

		const [addUrl, editUrl] = ["pro-manager/teachVideo/upload", "pro-manager/teachVideo/edit"];

		const result = await request({
			method: "POST",
			url: body.id ? editUrl : addUrl,
			data: formData
		});

		return result;
	}
};

export default {
	actions
};
