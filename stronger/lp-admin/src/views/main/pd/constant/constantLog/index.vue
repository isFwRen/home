<template>
	<div class="constant-log">
		<div class="text-center">
			<v-dialog v-model="dialog" width="1200">
				<v-card>
					<v-card-title class="text grey lighten-2"> 常量更新日志 </v-card-title>

					<v-card-text>
						<div class="z-flex mt-3 lp-filters">
							<v-row class="z-flex">
								<v-col cols="3">
									<z-date-picker
										:formId="searchFormId"
										formKey="range"
										label="日期"
										range
										:defaultValue="today"
									>
									</z-date-picker>
								</v-col>

								<v-col cols="2">
									<z-select
										:formId="searchFormId"
										formKey="type"
										hideDetails
										label="操作类型"
										:options="cells.options"
									></z-select>
								</v-col>
								<z-btn class="pt-6 pl-3" color="primary" @click="onSearch">
									<v-icon class="text-h6">mdi-magnify</v-icon>
									查询
								</z-btn>
							</v-row>
						</div>

						<div class="table config-audit-table">
							<vxe-table
								:data="desserts"
								resizable
								:border="tableBorder"
								:max-height="tableMaxHeight"
								:size="tableSize"
								:sort-config="{
									multiple: true,
									trigger: 'cell',
									defaultSort: { field: 'createdAt', order: 'asc' }
								}"
								@sort-change="handleSort"
							>
								<vxe-column type="seq" title="序号" width="60"></vxe-column>

								<template v-for="item in cells.headers">
									<vxe-column
										v-if="item.value === 'UpdatedAt'"
										:field="item.value"
										:fixed="item.fixed"
										:title="item.text"
										:key="item.value"
										:width="item.width"
									>
										<template #default="{ row }">
											{{ row.UpdatedAt | dateFormat("YYYY-MM-DD HH:mm:ss") }}
										</template>
									</vxe-column>

									<vxe-column
										v-else
										:field="item.value"
										:fixed="item.fixed"
										:title="item.text"
										:key="item.value"
										:width="item.width"
										sortable
									></vxe-column>
								</template>
							</vxe-table>

							<z-pagination
								class="mt-3"
								:total="pagination.total"
								:options="pageSizes"
								:pageNum="page.pageIndex"
								@page="handlePage"
							></z-pagination>
						</div>
					</v-card-text>

					<v-card-actions>
						<v-spacer></v-spacer>
						<v-btn color="primary" text @click="dialog = false">关闭 </v-btn>
					</v-card-actions>
				</v-card>
			</v-dialog>
		</div>
	</div>
</template>

<script>
import { mapState } from "vuex";
import cells from "./cells";
import TableMixins from "@/mixins/TableMixins";
import moment from "moment";
const today = [moment().format("YYYY-MM-DD"), moment().format("YYYY-MM-DD")];
export default {
	name: "constant-log",
	mixins: [TableMixins],
	data() {
		return {
			today,
			dialog: false,
			searchFormId: "constantLog",
			cells,
			dispatchList: "CONSTANT_OPERATION_LOG",
			manual: true,
			formId: "constantLog"
		};
	},
	methods: {
		openDialog() {
			this.dialog = true;
		}
	},
	computed: {
		...mapState(["forms"])
	},
	watch: {
		dialog(value) {
			if (value) {
				this.forms[this.searchFormId] = {
					range: today
				};

				this.onSearch();
			}
		}
	}
};
</script>

<style lang="scss" scoped></style>
