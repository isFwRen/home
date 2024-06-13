import { request } from "@/api/service";
import moment from "moment";

const actions = {
	async GET_HOMRPAGE_PRO_REPORT_LIST({}, params) {
		const result = await request({
			url: "homepage/pro-report/list",
			params: {
				reportDay: params.reportDay || moment().format("YYYY-MM-DD") + "T16:00:00.000Z"
			}
		});
		return result;
	},
	async EXPORT_COMMOUNT_FILE({}, { url, data, fileName }) {
		const result = await request({
			url: url,
			params: data,
			responseType: "blob"
		});
		var anchor = document.createElement("a");
		anchor.style.display = "none";
		anchor.setAttribute("download", fileName);
		anchor.href = URL.createObjectURL(result);
		document.body.appendChild(anchor);
		anchor.click();
		document.body.removeChild(anchor);
	},
	async SAVE_REPORT_PROJECT_EDIT({}, params) {
		const data = params;
		const result = await request({
			url: "report-management/project-report/set-report-info",
			method: "POST",
			data
		});
		return result;
	}
};

export default {
	actions
};
