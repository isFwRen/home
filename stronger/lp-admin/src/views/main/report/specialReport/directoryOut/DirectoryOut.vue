<template>
	<div class="directory-out">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col
					v-for="(item, index) in cells.fields"
					:key="`entry_filters_${index}`"
					:cols="item.cols"
				>
					<template v-if="item.inputType === 'date'">
						<z-date-picker
							:formId="searchFormId"
							:formKey="item.formKey"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:options="item.options"
							:range="item.range"
							:suffix="item.suffix"
							z-index="10"
							:defaultValue="item.defaultValue"
						></z-date-picker>
					</template>

					<template v-else>
						<z-select
							:formId="searchFormId"
							:formKey="item.formKey"
							:clearable="item.clearable"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:options="item.options"
							:suffix="item.suffix"
							:defaultValue="item.defaultValue"
							@change="handleChange($event, item)"
						></z-select>
					</template>
				</v-col>

				<div class="z-flex">
					<z-btn
						:formId="searchFormId"
						btnType="validate"
						class="pb-3 pl-3"
						color="primary"
						@click="onSearch"
					>
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>
				</div>
			</v-row>
		</div>

		<div class="pb-6 btns">
			<div class="pb-6 btns">
				<z-btn class="pr-3" color="primary" small outlined> 复制 </z-btn>
				<z-btn class="pr-3" color="primary" small outlined @click="onExport"> 导出 </z-btn>
			</div>
		</div>

		<div class="table directory-out-table">
			<vxe-table :border="tableBorder" :data="desserts" :size="tableSize">
				<template v-for="item in this.headers">
					<vxe-column
						v-if="item.value === 'CreatedAt'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
						:sortable="item.sortable"
						:formatter="formatDate"
					>
					</vxe-column>

					<vxe-column
						v-else
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
						:sortable="item.sortable"
					></vxe-column>
				</template>
			</vxe-table>

			<z-pagination
				:total="pagination.total"
				:options="pageSizes"
				:pageNum="page.pageIndex"
				@page="handlePage"
			></z-pagination>
		</div>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import SpecialReportMixins from "../SpecialReportMixins";
import cells from "./cells";
import { R } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
import moment from "moment";
import { mapState } from "vuex";

export default {
	name: "DirectoryOut",
	mixins: [TableMixins, SpecialReportMixins],

	data() {
		return {
			formId: "DirectoryOut",
			dispatchList: "GET_DIRECTORY_OUT_LIST",
			cells,
			headers: cells.headers
		};
	},
	created() {
		this.setConstOptions();
	},
	computed: {
		...mapState(["forms"])
	},
	methods: {
		onSearch() {
			const type = this.forms.DirectoryOutSearch.type
			this.headers[3].text = type === "1" ? "目录外医院名称" : "目录外清单名称"
			this.page.pageIndex = 1;

			this.params = {
				...this.params,
				...this.page
			};

			this.getList();
		},
		formatDate({cellValue}) {
			const date = moment(cellValue);
			return date.format("YYYY-MM-DD");
		},
		//导出
		async onExport() {
			const form = this.forms[this.searchFormId];
			if (!R.isYummy(form.time)) {
				this.toasted.warning("沒有选择日期或者项目");
				return;
			}

			const result = await this.$store.dispatch("EXPORT_DIRECTORY_OUT", form);
			lpTools.createExcelFun(
				result,
				`${form.proCode}目录外数据${form.time[0]}-${form.time[1]}`
			);
			this.toasted.dynamic(`正在导出...`, result.code);
		},
		// 设置除项目外的下拉options
		async setConstOptions() {
			const result = await this.$store.dispatch("GET_CASE_CONST_LIST");
			if (result.code === 200) {
				const { constType } = result.data;
				const arr = [];
				for (let key in constType) {
					arr.push({
						label: constType[key],
						value: key
					});
				}
				this.cells.fields[2].options = arr;
			}
		}
	}
};
</script>
