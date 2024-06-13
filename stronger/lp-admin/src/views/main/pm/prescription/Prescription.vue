<template>
	<div class="prescription">
		<div class="z-flex align-center">
			<z-btn class="mr-3" color="primary" small unlocked @click="selectAll">
				{{ selectedPro.length === auth.proItems.length ? "全不选" : "全选" }}
			</z-btn>

			<z-checkboxs
				:formId="formId"
				formKey="items"
				ref="items"
				:options="auth.proItems"
				:defaultValue="selectedPro"
				@change="selectProjects"
			></z-checkboxs>
		</div>

		<v-divider></v-divider>

		<v-subheader class="pl-0">所选项目({{ selectedPro.length }}个)：</v-subheader>

		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col
					v-for="(item, index) in cells.fields"
					:key="`prescriptionFilters_${index}`"
					:cols="item.cols"
				>
					<template v-if="item.inputType === 'input'">
						<z-text-field
							:formId="searchFormId"
							:formKey="item.formKey"
							:hideDetails="item.hideDetails"
							:label="item.label"
						>
						</z-text-field>
					</template>

					<template v-else-if="item.inputType === 'select'">
						<z-select
							:formId="searchFormId"
							:formKey="item.formKey"
							:hideDetails="item.hideDetails"
							:label="item.label"
							:options="item.options"
						></z-select>
					</template>

					<template v-else>
						<z-date-picker
							:formId="searchFormId"
							:formKey="item.formKey"
							:hideDetails="item.hideDetails"
							:label="item.label"
							:range="item.range"
							:defaultValue="item.defaultValue"
						></z-date-picker>
					</template>
				</v-col>

				<z-btn class="pb-3" color="primary" @click="onSearch">
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
				:cell-class-name="cellClassName"
				:sort-config="{
					multiple: true,
					trigger: 'cell',
					defaultSort: { field: 'createdAt', order: 'asc' }
				}"
				@sort-change="handleSort"
			>
				<template v-for="item in cells.headers">
					<vxe-column
						v-if="
							item.value === 'proCode' ||
							item.value === 'backAtTheLatest' ||
							item.value === 'caseStatus' ||
							item.value === 'caseNumber'
						"
						:field="item.value"
						:fixed="item.fixed"
						:title="item.text"
						:key="item.value"
						:width="item.width"
						:sortable="item.sortable"
					></vxe-column>

					<vxe-column
						v-else-if="item.value === 'timeRemaining'"
						:field="item.value"
						:fixed="item.fixed"
						:title="item.text"
						:key="item.value"
						:width="item.width"
						sort-type="string"
						sortable
					></vxe-column>

					<vxe-column
						v-else-if="item.value === 'claimType'"
						:field="item.value"
						:fixed="item.fixed"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							{{ computedType(row.claimType) }}
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
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";

const [half_hour, one_hour] = [1800, 3600];

export default {
	name: "Prescription",
	mixins: [TableMixins],
	data() {
		return {
			formId: "prescription",
			cells,
			manual: true,
			dispatchList: "GET_PERSCRIPTION_LIST",
			selectedPro: []
		};
	},
	methods: {
		computedType(type) {
			let label = "";
			cells.claimTypes.forEach(item => {
				if (item.value === type) {
					label = item.label;
				}
			});
			return label;
		},
		async onSearch() {
			this.params = {
				...this.params,
				...this.page
			};

			const total = await this.getList();
			if (typeof total !== "number") {
				return;
			}
			const index = this.page.pageIndex - 1;
			if (index * this.page.pageSize > total) {
				this.page.pageIndex = 1;
			}
		},
		async getList() {
			if (this.dispatchList) {
				const params = {
					...this.effectParams,
					...this.params,
					...this.forms[this.searchFormId]
				};
				const result = await this.$store.dispatch(this.dispatchList, params);
				const { list, total } = result.data;
				if (result.code === 200) {
					if (typeof list === "object") {
						if (list instanceof Array) {
							this.desserts = list;
						} else {
							this.desserts = [];
						}

						this.pagination.total = total;
					} else {
						this.desserts = result.data;
						this.pagination.total = this.desserts.length;
					}
				} else {
					this.toasted.error(result.msg);

					this.desserts = [];
					this.pagination.total = 0;
				}

				this.sabayon = result;
				return total;
			}

			this.loading = false;

			return this.sabayon;
		},
		selectProjects(values) {
			this.selectedPro = values;
			this.effectParams = {
				proCode: values.toString()
			};
		},

		selectAll() {
			this.$refs.items.onSelectAll();
		},

		onAll() {
			this.$refs.items.onSelectAll();
		},

		selectPro(values) {
			this.selectedPro = values;
		},

		cellClassName({ row, column }) {
			if (column.property === "timeRemaining") {
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
		...mapGetters(["auth"])
	}
};
</script>
