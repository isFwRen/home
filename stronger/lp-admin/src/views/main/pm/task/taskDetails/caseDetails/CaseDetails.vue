<template>
	<div class="case-details">
		<vxe-table
			:data="desserts"
			resizable
			:border="tableBorder"
			:max-height="tableMaxHeight"
			:size="tableSize"
			:sort-config="{ defaultSort: { field: 'scanAt', order: 'asc' } }"
			:cell-class-name="cellClassName"
		>
			<template v-for="item in cells.headers">
				<vxe-column
					v-if="item.value === 'scanAt'"
					:field="item.value"
					:title="item.text"
					:key="item.value"
					:sortable="item.sortable"
				>
					<template #default="{ row }">
						{{ row.scanAt }}
					</template>
				</vxe-column>
				<vxe-column
					v-else-if="item.value === 'options'"
					:field="item.value"
					:title="item.text"
					:key="item.value"
				>
					<template #default="{ row }">
						<v-icon class="mr-1" color="primary" @click="onDetail(row)"
							>mdi-text-box</v-icon
						>

						<v-icon
							:color="row.stickLevel === 1 ? 'primary' : ''"
							@click="onUrgent(row)"
							>mdi-lightning-bolt-circle</v-icon
						>
					</template>
				</vxe-column>

				<vxe-column
					v-else
					:field="item.value"
					:title="item.text"
					:sortable="item.sortable"
					:key="item.value"
				>
				</vxe-column>
			</template>
		</vxe-table>

		<block-detail-dialog ref="block"></block-detail-dialog>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import TableMixins from "@/mixins/TableMixins";
import DialogMixins from "@/mixins/DialogMixins";
import cells from "./cells";

const [half_hour, one_hour] = [1800, 3600];

export default {
	name: "CaseDetails",
	mixins: [TableMixins, DialogMixins],
	inject: ["task"],

	data() {
		return {
			cells,
			manual: true,
			dispatchList: "GET_CASE_DETAIL_LIST"
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
					if (item.stage !== "已完成") {
						item["appCompleteAt"] = "-";
					}
		
				});
			});
		}
	},

	created() {
		this.setEffectParams();
		//console.log(this.task)
	},

	// mounted() {
	//   this.params['proCode'] = this.task.proCode
	// },

	methods: {
		onFilters({ billNum, saleChannel }) {
			this.params = {
				...this.params,
				billNum,
				saleChannel
			};


			this.onSearch();
		},

		onDetail(row) {
			this.$refs.block.params["proCode"] = this.task.proCode;
			this.$refs.block.params["billId"] = row.billId;
			this.$refs.block.getList();
			this.$refs.block.onOpen();
		},

		async onUrgent(row) {
			this.$refs.block.params["proCode"] = this.task.proCode;
			this.$refs.block.params["stickLevel"] = row.stickLevel === 1 ? 99 : 1;
			this.$refs.block.params["list"] = [row.billNum];

			const result = await this.$store.dispatch(
				"TASK_SET_BILL_STICK_LEVEL",
				this.$refs.block.params
			);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.getList();
			}
		},

		setEffectParams() {
			this.effectParams = {
				proCode: this.task.proCode
			};

			this.getList();
		},

		cellClassName({ row, column }) {
			if (column.property === "remainderAt") {
				if (row.second > one_hour) {
					return "";
				} else if (row.second >= half_hour && row.second < one_hour) {
					return "warning-bg";
				} else {
					return "error-bg";
				}
			}
		}
	},

	computed: {
		...mapGetters(["project"])
	},

	components: {
		"block-detail-dialog": () => import("./blockDetailDialog")
	}
};
</script>
