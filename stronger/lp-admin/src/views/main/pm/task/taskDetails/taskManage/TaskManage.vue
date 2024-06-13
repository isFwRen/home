<template>
	<div class="task-manage">
		<div class="z-flex justify-between pb-4">
			<p class="mb-0 tips">
				<span class="error--text">
					{{ `${task.proCode}: ${proBillNumber || 0}` }}
				</span>

				<v-icon class="pb-1 ml-1">mdi-volume-high</v-icon>
			</p>

			<z-btn color="primary" depressed outlined @click="onRefresh">
				<v-icon size="18">mdi-autorenew</v-icon>
				刷新
			</z-btn>
		</div>

		<div class="table manage-table">
			<vxe-table border :data="desserts" :size="tableSize" @cell-click="onCell">
				<template v-for="item in cells.headers">
					<vxe-column
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
						:sortable="item.sortable"
					></vxe-column>
				</template>
			</vxe-table>
		</div>

		<v-row>
			<v-col v-for="item in cells.textareas" :cols="4" :key="item.formKey">
				<v-textarea
					:formId="`${task.proCode}${item.formId}`"
					:formKey="item.formKey"
					:label="item.label"
					:placeholder="item.placeholder"
					v-model="taskInfo[item.val]"
					@focus="handleFocus(item.val)"
					@blur="handleBlur(item.val)"
				></v-textarea>

				<div class="z-flex justify-end btns">
					<z-btn class="mr-3" color="primary" small @click="onShowDialog(item)">
						新增
					</z-btn>

					<z-btn
						:formId="`${task.proCode}${item.formId}`"
						color="error"
						small
						@click="onClear(item)"
					>
						清空
					</z-btn>
				</div>
			</v-col>
		</v-row>

		<!-- 待分配 BEGIN -->
		<unassigned-dialog ref="unassigned" :cellInfo="cellInfo"></unassigned-dialog>
		<!-- 待分配 END -->

		<!-- 已分配 BEGIN -->
		<assigned-dialog ref="assigned" :cellInfo="cellInfo"></assigned-dialog>
		<!-- 已分配 END -->

		<!-- 缓存区 BEGIN -->
		<buffer-dialog ref="buffer" :cellInfo="cellInfo"></buffer-dialog>
		<!-- 缓存区 END -->

		<add-dialog ref="add" :firstInfo="firstInfo" :taskInfos="taskInfo"></add-dialog>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";

const REF_NAMES = ["unassigned", "assigned", "buffer"];

export default {
	name: "TaskManage",
	mixins: [TableMixins],
	inject: ["task"],

	data() {
		return {
			formId: "TaskManage",
			dispatchList: "GET_TASK_LIST",
			cells,
			proCode: "B0118",
			proBillNumber: "0",

			cellInfo: {},
			firstInfo: {},
			taskInfo: {
				urgent: "",
				priority: "",
				agency: ""
			}
		};
	},

	mounted() {
		this.params["proCode"] = this.task.proCode;
	},

	methods: {
		// 刷新
		onRefresh() {
			this.params["proCode"] = this.task.proCode;
			this.proBillNumber = this.sabayon.data.num || 0;
			this.getList();
		},

		// 清空
		async onClear(item) {
			// const list = this.forms[`${this.task.proCode}${item.formId}`][item.formKey];

			var result = [];
			if (item.formKey == "organizationNumber") {
				result = await this.$store.dispatch("SET_PRIORITY_ORGANIZATION_NUMBER_ITEM_LIST", {
					list: this.taskInfo[item.val].split(","),
					proCode: this.params["proCode"],
					stickLevel: 99
				});
			} else {
				console.log("非机构号urgent, priority");
				result = await this.$store.dispatch("TASK_SET_BILL_STICK_LEVEL", {
					list: this.taskInfo[item.val].split("\n"),
					proCode: this.params["proCode"],
					stickLevel: 99
				});
			}
			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.onRefresh();
			}
		},

		handleFocus(name) {
			if (this.taskInfo[name]) {
				const regex = /\n$/;
				const match = this.taskInfo[name].match(regex);
				if (!match) {
					this.taskInfo[name] = this.taskInfo[name] + "\n";
				}
			}
		},
		handleBlur(name) {
			if (this.taskInfo[name]) {
				const regex = /\n$/;
				const match = this.taskInfo[name].match(regex);
				if (match) {
					this.taskInfo[name] = this.taskInfo[name].replace(regex, "");
				}
			}
		},

		// 打开弹窗
		onCell(cell) {
			// 数量为0
			if (!cell.row[cell.column.field]) {
				return;
			}

			let type = "any";
			if (cell.row.name === "已分配") {
				type = cell.column.title ? cell.column.title.slice(0, 2) : false;
			}

			this.cellInfo = {
				columnIndex: cell.$columnIndex,
				rowIndex: cell.$rowIndex,
				title: cell.column.title,
				op: cell.column.field.substr(0, 3),
				type
			};

			this.$refs[REF_NAMES[this.cellInfo.rowIndex]].onOpen();
		},

		// 打开表单弹窗
		onShowDialog(field) {
			this.firstInfo = {
				val: field.val,
				label: field.label,
				formId: field.formId,
				formKey: field.formKey,
				placeholder: field.placeholder
			};
			this.$refs.add.onSubmit(this.taskInfo[field.val], this.firstInfo);
		}
	},

	components: {
		"unassigned-dialog": () => import("./unassignedDialog"),
		"assigned-dialog": () => import("./assignedDialog"),
		"buffer-dialog": () => import("./bufferDialog"),
		"add-dialog": () => import("./addDialog")
	},
	watch: {
		"params.proCode": {
			handler() {
				this.params["proCode"] = this.task.proCode;
			},
			immediate: true,
			deep: true
		},
		// \r 回车  \n 换行
		"sabayon.data.num": {
			handler() {
				this.proBillNumber = this.sabayon.data.num;
				this.taskInfo.urgent = this.sabayon.data["urgent"]?.join("\n");
				this.taskInfo.priority = this.sabayon.data["priority"]?.join("\n");
				this.taskInfo.agency = this.sabayon.data["agency"]?.join(",");
			},
			immediate: true,
			deep: true
		}
	}
};
</script>
