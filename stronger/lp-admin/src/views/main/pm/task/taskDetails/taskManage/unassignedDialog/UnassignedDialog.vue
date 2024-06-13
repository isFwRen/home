<template>
	<lp-dialog
		ref="dialog"
		:title="`${task.proCode}${cellInfo.title}(待分配)`"
		transition="dialog-bottom-transition"
		width="900"
		@dialog="handleDialog"
	>
		<div class="pt-6" slot="main">
			<div class="table">
				<vxe-table :data="desserts" :size="tableSize">
					<template v-for="item in cells.headers">
						<vxe-column
							v-if="item.value === 'taskAssign'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
						>
							<template #default="{ row }">
								<z-autocomplete
									:formId="formId"
									:formKey="row.blockId"
									class="mt-n4"
									hide-details
									:options="staffItems"
									@change="handleStaff($event, row)"
								></z-autocomplete>
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
		</div>
	</lp-dialog>
</template>

<script>
import { tools } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import DialogMixins from "@/mixins/DialogMixins";
import cells from "./cells";

export default {
	name: "UnassignedDialog",
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
			formId: "UnassignedDialog",
			cells,
			dispatchList: "GET_TASK_DETAIL_LIST",
			manual: true,

			staffItems: []
		};
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

		// 获取待分配列表
		getUnassignedList() {
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

		// 分配任务
		async handleStaff(value, row) {
			const body = {
				id: row.blockId,
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
			dialog && this.getUnassignedList();
		}
	}
};
</script>
