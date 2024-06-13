<template>
	<lp-dialog
		ref="dialog"
		:title="`${task.proCode}${cellInfo.title}(已分配)`"
		transition="dialog-bottom-transition"
		width="1200"
		@dialog="handleDialog"
	>
		<div class="pt-6" slot="main">
			<div class="table">
				<vxe-table
					:data="desserts"
					:size="tableSize"
					:sort-config="{ defaultSort: { field: 'scanAt', order: 'asc' } }"
				>
					<template v-for="item in headers">
						<vxe-column
							v-if="item.value === 'options'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:width="item.width"
						>
							<template #default="{ row }">
								<v-tooltip bottom>
									<template v-slot:activator="{ on, attrs }">
										<v-icon
											color="primary"
											dark
											v-bind="attrs"
											v-on="on"
											@click="openStaffDialog(row)"
										>
											mdi-pencil-circle
										</v-icon>
									</template>
									<span>修改</span>
								</v-tooltip>

								<v-tooltip bottom>
									<template v-slot:activator="{ on, attrs }">
										<v-icon
											color="error"
											dark
											v-bind="attrs"
											v-on="on"
											@click="onRelease(row)"
										>
											mdi-minus-circle
										</v-icon>
									</template>
									<span>释放</span>
								</v-tooltip>
							</template>
						</vxe-column>

						<vxe-column
							v-else-if="item.value === 'op0ApplyAt'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:sortable="item.sortable"
						>
							<template #default="{ row }">
								{{ row.op0ApplyAt }}
							</template>
						</vxe-column>

						<vxe-column v-else :field="item.value" :title="item.text" :key="item.value">
							<template #default="{ row }">
								{{ row[item.value] | ifLousyValue }}
							</template>
						</vxe-column>
					</template>
				</vxe-table>

				<z-pagination
					:pageNum="params.pageIndex"
					:total="pagination.total"
					:options="pageSizes"
					@page="handlePage"
				></z-pagination>
			</div>

			<lp-dialog
				ref="staffDialog"
				title="任务分配"
				transition="dialog-bottom-transition"
				width="300"
			>
				<div slot="main">
					<z-autocomplete
						:formId="formId"
						formKey="staff"
						hide-details
						:options="staffItems"
						@change="handleStaff"
					></z-autocomplete>
				</div>
			</lp-dialog>
		</div>
	</lp-dialog>
</template>

<script>
import { tools } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import DialogMixins from "@/mixins/DialogMixins";
import cells from "./cells";

export default {
	name: "AssignedDialog",
	mixins: [TableMixins, DialogMixins],
	inject: ["task"],

	props: {
		cellInfo: {
			type: Object,
			default: () => {}
		}
	},

	data() {
		return {
			formId: "AssignedDialog",
			cells,
			dispatchList: "GET_TASK_DETAIL_LIST",
			manual: true,

			lineInfo: {},
			staffItems: []
		};
	},
	computed: {
		headers() {
			const type = this.cellInfo.type;
			const arr = [];
			if (type !== "any") {
				cells.headers.forEach(item => {
					if (item.type === type || item.type === "required") {
						arr.push(item);
					}
				});
				return arr;
			}
			return cells.headers;
		}
	},
	methods: {
		// 获取员工列表
		async getStaffList() {
			this.staffItems = [];
			const params = {
				pageIndex: 1,
				pageSize: 1500,
				proCode: this.task.proCode,
				op: this.cellInfo.op
			};

			const result = await this.$store.dispatch("TASK_GET_STAFF_LIST", params);

			if (result.code === 200) {
				if (tools.isYummy(result.data.list)) {
					result.data.list.map(staff => {
						this.staffItems.push({
							label: staff.name,
							value: staff.code
						});
					});
				}
			}
		},

		// 释放任务
		onRelease(row) {
			this.$modal({
				visible: true,
				title: "释放任务提醒",
				content: "请确认是否要释放任务？",
				confirm: async () => {
					const body = {
						id: row.blockId,
						op: this.cellInfo.op,
						code: "",
						project: this.task.project
					};


					const result = await this.$store.dispatch("TASK_ALLOCATION_TASK", body);

					this.toasted.dynamic(result.msg, result.code);
				}
			});
		},

		// 获取已分配列表
		getAssignedList() {
			this.params.pageIndex = 1;

			const { columnIndex, rowIndex, op } = this.cellInfo;
			const isExpenseAccount = [2, 4].includes(columnIndex)
				? "1"
				: [3, 5].includes(columnIndex)
				? "2"
				: "";

			this.effectParams = {
				opStage: `${rowIndex + 1}`,
				proCode: this.task.proCode,
				op,
				isExpenseAccount
			};

			this.getList();
			this.getStaffList();
		},

		async openStaffDialog(row) {
			this.lineInfo = row;
			this.$refs.staffDialog.onOpen();
		},

		// 分配任务
		async handleStaff(value) {
			const body = {
				id: this.lineInfo.blockId,
				op: this.cellInfo.op,
				code: value,
				project: this.task.project
			};

			const result = await this.$store.dispatch("TASK_ALLOCATION_TASK", body);

			this.toasted.dynamic(result.msg, result.code);
		}
	},

	watch: {
		dialog(dialog) {
			dialog && this.getAssignedList();
		}
	}
};
</script>
