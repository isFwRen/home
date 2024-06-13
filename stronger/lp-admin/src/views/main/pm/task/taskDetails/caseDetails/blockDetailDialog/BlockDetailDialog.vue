<template>
	<lp-dialog
		ref="dialog"
		title="分块明细"
		transition="dialog-bottom-transition"
		width="1400"
		max-width="1400"
	>
		<div class="pt-6" slot="main">
			<div class="table">
				<vxe-table :data="dataSource" :size="tableSize" align="center" @cell-click="onCell">
					<template v-for="item in cells.headers">
						<vxe-column
							v-if="item.value === 'op0Code'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							width="100"
						>
						</vxe-column>

						<vxe-column
							v-else-if="item.value === 'op1Code'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							width="125"
						>
							<template #default="{ row }">
								{{ row.op1Code }}
								<v-tooltip
									bottom
									v-if="
										(row.op1SubmitAt === '-' || row.op1SubmitAt === '') &&
										row.op1Code
									"
								>
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
							v-else-if="item.value === 'op2Code'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							width="125"
						>
							<template #default="{ row }">
								{{ row.op2Code }}
								<v-tooltip
									bottom
									v-if="
										(row.op2SubmitAt === '-' || row.op2SubmitAt === '') &&
										row.op2Code
									"
								>
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
							v-else-if="item.value === 'opqCode'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							width="100"
						>
						</vxe-column>

						<vxe-column
							v-else-if="item.value === 'op0SubmitAt'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							width="150"
						>
						</vxe-column>

						<vxe-column
							v-else-if="item.value === 'op1SubmitAt'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							width="150"
						>
						</vxe-column>

						<vxe-column
							v-else-if="item.value === 'op2SubmitAt'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							width="150"
						>
						</vxe-column>

						<vxe-column
							v-else-if="item.value === 'opqSubmitAt'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							width="150"
						>
						</vxe-column>

						<vxe-column
							v-else-if="item.value === 'status'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							width="110"
						>
						</vxe-column>

						<vxe-column
							v-else-if="item.value === 'options'"
							:field="item.value"
							:title="item.text"
							:key="item.value"
							width="60"
						>
							<template #default="{ row }">
								<v-icon class="mr-1" color="primary" @click="onDetail(row)"
									>mdi-text-box</v-icon
								>
							</template>
						</vxe-column>

						<vxe-column
							v-else-if="item.value.indexOf('SubmitAt') != -1"
							:field="item.value"
							:title="item.text"
							:key="item.value"
						>
							<template #default="{ row }">
								{{ row[item.value] | dateFormat() }}
							</template>
						</vxe-column>

						<vxe-column
							v-else
							:field="item.value"
							:title="item.text"
							:width="item.width"
							:key="item.value"
						>
						</vxe-column>
					</template>
				</vxe-table>
			</div>

			<field-detail-dialog ref="field"></field-detail-dialog>
		</div>
	</lp-dialog>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import DialogMixins from "@/mixins/DialogMixins";
import cells from "./cells";
import _ from "lodash";

export default {
	name: "AssignedDialog",
	mixins: [TableMixins, DialogMixins],
	inject: ["task"],
	data() {
		return {
			formId: "AssignedDialog",
			cells,
			dispatchList: "GET_CASE_DETAIL_BLOCK_LIST",
			manual: true,
			dataSource: [],
			cellInfo: {}
		};
	},
	watch: {
		desserts(val) {
			const onRE = /^000[0-9]/;
			val.map(item => {
				cells.handleDate.map(key => {
					// 匹配00开头非法日期
					if (onRE.test(item[key])) {
						item[key] = "-";
					}
				});
			});

			console.log(val,'val')
			const sortArrs = _.orderBy(val, ["status"]);
			this.dataSource = _.cloneDeep(sortArrs);
		}
	},
	methods: {
		// 打开弹窗
		onCell(cell) {
			const opFields = ["op1Code", "op2Code"];
			// 数量为0
			if (!cell.row[cell.column.field]) {
				return;
			}
			if (!opFields.includes(cell.column.field)) {
				return;
			}

			this.cellInfo = {
				op: cell.column.field.substr(0, 3)
			};
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
					if (result.code === 200) {
						this.getList();
					}
					this.toasted.dynamic(result.msg, result.code);
				}
			});
		},
		onDetail(row) {
			this.$refs.field.params["proCode"] = this.task.proCode;
			this.$refs.field.params["billId"] = row.billID;
			this.$refs.field.params["blockId"] = row.blockId;
			this.$refs.field.onOpen();
			this.$refs.field.getList();
		}
	},

	components: {
		"field-detail-dialog": () => import("./fieldDetailDialog")
	}
};
</script>
