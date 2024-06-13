<template>
	<div class="institutional-extraction">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col v-for="(item, index) in cells.fields" :key="`entry_filters_${index}`" :cols="item.cols">
					<template v-if="item.inputType === 'date'">
						<z-date-picker :formId="searchFormId" :formKey="item.formKey" :hideDetails="item.hideDetails"
							:hint="item.hint" :label="item.label" :options="item.options" :range="item.range" :suffix="item.suffix"
							z-index="10" :defaultValue="item.defaultValue"></z-date-picker>
					</template>

					<template v-else>
						<z-select :formId="searchFormId" :formKey="item.formKey" :hideDetails="item.hideDetails" :hint="item.hint"
							:label="item.label" :clearable="item.clearable" :options="item.options" :suffix="item.suffix"
							@change="handleChange($event, item)"></z-select>
					</template>
				</v-col>

				<div class="z-flex">
					<z-btn class="pb-3 pl-3" color="primary" @click="onSearch" :formId="searchFormId">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>
					<z-btn class="pb-3 pl-3" color="success" @click="onExport('EXPORT_INCOME_ANALYSIS')"> 导出 </z-btn>
					<z-btn class="pb-3 pl-3" color="error" @click="onExport('EXPORT_HISTORY_ANALYSIS')"> 获取历史数据 </z-btn>

				</div>
			</v-row>
		</div>

		<div class="table staff-total-table">
			<vxe-table :border="tableBorder" :data="desserts" :size="tableSize">
				<template v-for="item in cells.headers">
					<vxe-column :field="item.value" :title="item.text" :key="item.value">
					</vxe-column>
				</template>
			</vxe-table>

			<z-pagination :total="pagination.total" :options="pageSizes" :pageNum="page.pageIndex"
				@page="handlePage"></z-pagination>

		</div>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import SpecialReportMixins from "../SpecialReportMixins";
import cells from "./cells";
import { R } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
import { mapState } from "vuex";

export default {
	name: "institutionalExtraction",
	mixins: [TableMixins, SpecialReportMixins],

	data() {
		return {
			formId: "institutionalExtraction",
			project: "1",
			cells,
			dispatchList: "GET_INSTITUTION_EXTRACTION_LIST"
		};
	},
	computed: {
		...mapState(["forms"])
	},
	methods: {
		//导出
		async onExport(name) {
			const form = this.forms[this.searchFormId];

			if (!R.isYummy(form.time)) {
				this.toasted.warning("沒有选择日期或者项目");
				return;
			}

			const result = await this.$store.dispatch(name, form);
			lpTools.createExcelFun(result, `${form.proCode}机构抽取${form.time[0]}-${form.time[1]}`);
			this.toasted.dynamic(`正在导出...`, result.code);
		}
	}
};
</script>
