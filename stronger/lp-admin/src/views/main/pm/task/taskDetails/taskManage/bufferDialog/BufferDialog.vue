<template>
	<lp-dialog
		ref="dialog"
		:title="`${task.proCode}${cellInfo.title}(缓存区)`"
		transition="dialog-bottom-transition"
		width="1200"
		@dialog="handleDialog"
	>
		<div class="pt-6" slot="main">
			<div class="table">
				<vxe-table :data="desserts" :size="tableSize">
					<template v-for="item in cells.headers">
						<vxe-column :field="item.value" :title="item.text" :key="item.value">
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
import TableMixins from "@/mixins/TableMixins";
import DialogMixins from "@/mixins/DialogMixins";
import cells from "./cells";

export default {
	name: "BufferDialog",
	mixins: [TableMixins, DialogMixins],
	inject: ["task"],
	data() {
		return {};
	},
	props: {
		cellInfo: {
			type: Object,
			default: () => {}
		}
	},
	data() {
		return {
			formId: "BufferDialog",
			cells,
			dispatchList: "GET_TASK_DETAIL_LIST",
			manual: true
		};
	},

	methods: {
		// 获取缓存区列表
		getBufferList() {
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
		}
	},

	watch: {
		dialog(dialog) {
			dialog && this.getBufferList();
		},
		desserts(val) {
			const onRE = /^000[0-9]/;
			val.forEach(item => {
				if (onRE.test(item["op1SubmitAt"])) {
					item["op1SubmitAt"] = "-";
				}
				if (onRE.test(item["op2SubmitAt"])) {
					item["op2SubmitAt"] = "-";
				}
			});
		}
	}
};
</script>
