<template>
	<div class="business-rules" style="min-width: 1090px">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col :cols="3">
					<z-date-picker
						:formId="searchFormId"
						formKey="dayRange"
						clearable
						hideDetails
						label="日期范围"
						prepend-icon="mdi-calendar"
						range
						:defaultValue="DEFAULT_DATE"
					></z-date-picker>
				</v-col>

				<v-col :cols="2">
					<z-select
						:formId="searchFormId"
						formKey="releaseType"
						clearable
						hideDetails
						label="发布类型"
						:options="cells.types"
					></z-select>
				</v-col>

				<v-col :cols="2">
					<z-select
						:formId="searchFormId"
						formKey="status"
						clearable
						hideDetails
						label="发布状态"
						:options="cells.status"
					></z-select>
				</v-col>

				<v-col :cols="2">
					<z-select
						:formId="searchFormId"
						formKey="proCode"
						clearable
						hideDetails
						label="项目编码"
						:options="auth.proItems"
					></z-select>
				</v-col>

				<div class="z-flex">
					<z-btn class="pb-3 px-3" color="primary" @click="onSearch">
						<v-icon size="20">mdi-magnify</v-icon>
						查询
					</z-btn>

					<z-btn color="primary" @click="onAddItem">
						<v-icon>mdi-plus</v-icon>
						新增
					</z-btn>
				</div>
			</v-row>
		</div>

		<div class="table">
			<vxe-table :data="desserts" :border="tableBorder" :size="tableSize">
				<vxe-column type="seq" title="序号" width="60"></vxe-column>

				<template v-for="item in cells.headers">
					<vxe-column
						v-if="item.value === 'options'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<div class="py-2 z-flex">
								<z-btn
									class="mr-2"
									color="primary"
									depressed
									outlined
									small
									@click="onEditItem(row)"
									v-if="row.status !== 3"
								>
									编辑
								</z-btn>

								<z-btn
									class="mr-2"
									color="primary"
									depressed
									outlined
									small
									@click="isWatch(row)"
								>
									查看
								</z-btn>

								<z-btn
									class="mr-2"
									color="primary"
									depressed
									outlined
									small
									@click="isRelease(row, row.status == 2 ? 1 : 2)"
									v-if="row.status !== 3"
								>
									{{ row.status == 2 ? "取消发布" : "发布" }}
								</z-btn>

								<z-btn
									color="error"
									depressed
									outlined
									small
									@click="onDeleteItem(row, row.status == 3 ? 1 : 3)"
								>
									{{ row.status == 3 ? "恢复" : "删除" }}
								</z-btn>
							</div>
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value == 'status'"
						:field="item.value"
						:title="item.text"
						:width="item.width"
						:key="item.value"
						><template #default="{ row }">
							{{
								row.status
									? cells.status.filter(e => {
											return e.value == row.status;
									  })[0].label
									: ""
							}}
						</template>
					</vxe-column>
					<vxe-column
						v-else-if="item.value == 'releaseDate'"
						:field="item.value"
						:title="item.text"
						:width="item.width"
						:key="item.value"
						><template #default="{ row }">
							{{
								row.status == 2
									? row.releaseDate &&
									  row.releaseDate.substr(0, 10) +
											" " +
											row.releaseDate.substr(11, 5)
									: "-"
							}}
						</template>
					</vxe-column>
					<vxe-column
						v-else-if="item.value == 'releaseType'"
						:field="item.value"
						:title="item.text"
						:width="item.width"
						:key="item.value"
						><template #default="{ row }">
							{{
								row.releaseType
									? cells.types.filter(e => {
											return e.value == row.releaseType;
									  })[0].label
									: ""
							}}
						</template>
					</vxe-column>
					<vxe-column
						v-else-if="item.value == 'releaseUserName'"
						:field="item.value"
						:title="item.text"
						:width="item.width"
						:key="item.value"
						><template #default="{ row }">
							{{
								row.status == 2
									? row.releaseUserName + "(" + row.releaseUserCode + ")"
									: "-"
							}}
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
		<z-pagination
			:options="pageSizes"
			:pageNum="page.pageIndex"
			@page="handlePage"
			:total="pagination.total"
		></z-pagination>
		<update-dialog
			ref="updateDialog"
			:rowInfo="detailInfo"
			@submitted="getList"
		></update-dialog>
		<WatchDialog ref="watchDialog" :content="detailInfo.content" />
	</div>
</template>

<script>
import moment from "moment";
import { mapGetters } from "vuex";
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";
import WatchDialog from "./updateDialog/WatchDialog.vue";

const date = new Date();
const [year, month, day] = [date.getFullYear(), date.getMonth(), date.getDate()];
const DEFAULT_DATE = [
	moment(`${year}-${month}-${day}`).format("YYYY-MM-DD"),
	moment().format("YYYY-MM-DD")
];

export default {
	name: "Notice",
	mixins: [TableMixins],

	data() {
		return {
			formId: "Notice",
			dispatchList: "GET_PM_NOTICE_LIST",
			DEFAULT_DATE,
			cells,
			desserts: [{}],
			total: 0
		};
	},

	methods: {
		async changeStatus(row, status) {
			const body = {
				ids: [row.ID],
				status: status
			};
			const result = await this.$store.dispatch("CHANGE_PM_NOTICE_LIST_ITEM_STATUS", body);

			if (result.code == 200 && result.data == 1) {
				this.onSearch();
				this.toasted.success(result.msg);
			} else {
				this.toasted.error(result.msg);
			}
		},

		isRelease(row, status) {
			let opition = status == 2 ? "发布" : "取消发布";
			this.$modal({
				visible: true,
				title: opition + "提示",
				content: `请确认是否要${opition}？`,
				confirm: () => {
					this.changeStatus(row, status);
				}
			});
		},
		isWatch(row) {
			this.detailInfo = row;
			this.$refs.watchDialog.onOpen(1);
		},
		onAddItem() {
			this.detailInfo = {
				content: " "
			};
			this.$refs.updateDialog.onOpen(-1);
		},

		onEditItem(row) {
			this.detailInfo = {};
			this.getDetail(row);
			this.$refs.updateDialog.onOpen(1);
		},

		onDeleteItem(row, status) {
			let opition = status == 3 ? "删除" : "恢复";
			this.$modal({
				visible: true,
				title: opition + "提示",
				content: `请确认是否要${opition}？`,
				confirm: () => {
					this.changeStatus(row, status);
				}
			});
		}
	},

	computed: {
		...mapGetters(["auth"])
	},

	components: {
		"update-dialog": () => import("./updateDialog"),
		WatchDialog
	}
};
</script>
