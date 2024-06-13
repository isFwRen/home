import { request } from "@/api/service";
import { tools } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";

const isIntranet = lpTools.isIntranet();

const actions = {
	// 获取录入通道列表(与录入系统录入通道模块的列表一致)
	async PM_TASK_GET_INPUT_CHANNEL_LIST() {
		const result = await request({
			url: "data-entry/channel/list"
		});

		return result;
	},

	// 紧急件 优先件的提交
	async SET_URGENCY_BILL_OR_PRIORITY_BILL_ITEM_LIST({}, body) {
		const data = {
			caseNumbers: body.list.split(","),
			proCode: body.proCode,
			stickLevel: body.item.formKey === "stickLevel2" ? 2 : 1
		};

		const result = await request({
			method: "POST",
			url: "pro-config/task/UrgencyBillOrPriorityBill/edit",
			data
		});

		return result;
	},

	// 案件明细紧急件提交
	async SET_URGENCY_BILL_OR_PRIORITY_BILL_ITEM_LIST2({}, body) {
		const data = {
			caseNumbers: body.billName,
			proCode: body.proCode,
			stickLevel: body.stickLevel
		};

		const result = await request({
			method: "POST",
			url: "pro-config/task/UrgencyBillOrPriorityBill/edit",
			data
		});
		return result;
	},

	// 机构号提交
	async SET_PRIORITY_ORGANIZATION_NUMBER_ITEM_LIST({}, body) {
		const data = {
			organizationNumber: body.list,
			proCode: body.proCode,
			stickLevel: body.stickLevel
		};

		const result = await request({
			method: "POST",
			url: "pro-config/task/PriorityOrganizationNumber/edit",
			data
		});

		return result;
	},

	// 获取任务列表
	async GET_TASK_LIST({}, body) {
		const data = {
			proCode: body.proCode
		};

		const result = await request({
			url: "pro-config/task/List",
			params: data
		});

		const clues = ["NotAssign", "Assign", "Cache"];
		const names = ["待分配", "已分配", "缓存区"];
		const list = [];

		if (result.code === 200 && tools.isYummy(result.data)) {
			clues.map((clue, index) => {
				list.push({
					name: names[index],
					op0: result.data[`op0${clue}`],
					op1: result.data[`op1${clue}ExpenseAccount`],
					op1No: result.data[`op1${clue}NotExpenseAccount`],
					op2: result.data[`op2${clue}ExpenseAccount`],
					op2No: result.data[`op2${clue}NotExpenseAccount`],
					opq: result.data[`opq${clue}`]
				});
			});
		}

		if (tools.getType(result.data) === "string") {
			result.data = {};
		}

		result.data.list = list;

		return result;
	},

	// 查询待分配-已分配-缓存区-紧急件-优先件 详细信息
	async GET_TASK_DETAIL_LIST({}, params) {
		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proCode: params.proCode,
			op: params.op,
			opStage: params.opStage,
			isExpenseAccount: params.isExpenseAccount
		};

		const result = await request({
			url: "pro-config/task/detail/List",
			params: data
		});

		return result;
	},

	// 获取员工列表
	async TASK_GET_STAFF_LIST({}, params) {
		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			proCode: params.proCode,
			op: params.op
		};

		const result = await request({
			url: "pro-config/task/getProPermissionPeople",
			params: data
		});

		return result;
	},

	// 任务分配/释放任务
	async TASK_ALLOCATION_TASK({}, body) {
		const data = {
			id: body.id,
			op: body.op,
			code: body.code
		};

		const { innerIp, inAppPort, outIp, outAppPort } = body.project;

		const origin = `https://${
			isIntranet ? `${innerIp}:${inAppPort}` : `${outIp}:${outAppPort}`
		}/`;

		const result = await request({
			method: "POST",
			baseURL: origin,
			url: "api/task/releaseBlock",
			data,
			headers: {
				process: ""
			},
			effect: { task: true }
		});

		return result;
	},

	// 设置单据紧急状态
	async TASK_SET_BILL_STICK_LEVEL({}, body) {
		const data = {
			proCode: body.proCode,
			caseNumbers: body.list,
			stickLevel: body.stickLevel
		};

		const result = await request({
			method: "POST",
			url: "pro-config/task/UrgencyBillOrPriorityBill/edit",
			data
		});

		return result;
	}
};

export default {
	actions
};
