<template>
	<div class="identify-statistics">
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

					<template v-else-if="item.inputType === 'input'">
						<z-text-field
							:formId="searchFormId"
							:formKey="item.formKey"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:suffix="item.suffix"
						>
						</z-text-field>
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

					<z-btn :formId="formId" class="pb-3 pl-3" color="error" @click="onExport">
						导出
					</z-btn>
				</div>
			</v-row>
		</div>

		<!--table-->
		<div class="table staff-total-table" id="tableList">
			<vxe-table
				:border="tableBorder"
				:data="desserts"
				:size="tableSize"
				:scroll-y="{ enabled: false }"
				align="center"
			>
				<template v-for="item in cells.headers">
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
						v-else-if="item.value === 'result_value'"
						:field="item.value"
						:title="item.text"
						type="html"
						:key="item.value + Math.random()"
						:width="item.width"
						:sortable="item.sortable"
					>
						<template #default="{ row }">
							<span v-html="diffStr(row, item.value)"></span>
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'rate'"
						:field="item.value"
						:title="item.text"
						:key="item.value + Math.random()"
						:width="item.width"
						:sortable="item.sortable"
					>
						<template #default="{ row }"> {{ row.rate }}% </template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'disable'"
						:field="item.value"
						:title="item.text"
						:key="item.value + Math.random()"
						:width="item.width"
						:sortable="item.sortable"
					>
						<template #default="{ row }">
							{{ row.disable === "1" ? "是" : "否" }}
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'compare'"
						:field="item.value"
						:title="item.text"
						:key="item.value + Math.random()"
						:width="item.width"
						:sortable="item.sortable"
					>
						<template #default="{ row }">
							{{ row.compare === "1" ? "是" : "否" }}
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'name'"
						:field="item.value"
						:title="item.text"
						:key="item.value + Math.random() * Math.random()"
						:width="item.width"
						:sortable="item.sortable"
					>
						<template #default="{ row }">
							<span @click="onViewImage(row)"> {{ row.name }}</span>
						</template>
					</vxe-column>

					<vxe-column
						v-else
						:field="item.value"
						:title="item.text"
						:width="item.width"
						:key="item.value + Math.random() + 0.1"
					>
					</vxe-column>
				</template>
			</vxe-table>

			<z-pagination
				:total="pagination.total"
				:options="pageSizes"
				@page="handlePage"
				:pageNum="page.pageIndex"
			></z-pagination>
		</div>
	</div>
</template>
<script>
import cells from "./cells";
import SpecialReportMixins from "../SpecialReportMixins";
import TableMixins from "@/mixins/TableMixins";
import { mapState } from "vuex";
import moment from "moment";
import { tools as lpTools } from "@/libs/util";
import { R,sessionStorage } from "vue-rocket";
const { baseURLApi } = lpTools.baseURL();

export default {
	name: "identifyStatistics",
	mixins: [SpecialReportMixins, TableMixins],
	data() {
		return {
			cells,
			dispatchList: "GET_IDENTIFY_STATISTICS_LIST",
			formId: "identifyStatistics"
		};
	},
	computed: {
		...mapState(["forms"])
	},
	methods: {
		formatDate({ cellValue }) {
			const date = moment(cellValue).format("YYYY-MM-DD");
			return date;
		},
		onViewImage(row) {
			let prefixUrl = `${baseURLApi}files/${row.pic}`;
			const thumbs = [];

			thumbs.push({
				thumbPath: prefixUrl,
				path: prefixUrl
			});
			sessionStorage.set("thumbs", thumbs);
			window.open(
				`${location.origin}/normal/view-images`,
				"_blank",
				"toolbar=yes, scrollbars=yes, resizable=yes"
			);
		},
		diffStr(row, column) {
			const ocr = row.value;
			const upload = row.result_value;
			if (!upload) {
				return "";
			}
			const compare = {
				len: ocr.length,
				diff: 0,
				newStr: []
			};
			for (let i = 0; i < ocr.length; i++) {
				if (upload[i]) {
					if (ocr[i] !== upload[i]) {
						compare.newStr.push(`<span style="color: red">${upload[i]}</span>`);
					} else {
						compare.newStr.push(upload[i]);
					}
				}
			}
			if (ocr.length < upload.length) {
				const len = upload.length - ocr.length;
				for (let i = 0; i < len; i++) {
					compare.newStr.push(
						`<span style="color: red">${upload[ocr.length + i]}</span>`
					);
				}
			}

			if (column === "result_value") {
				let str = "";
				compare.newStr.map(item => (str += item));
				return str;
			}
		},
		async onExport() {
			const form = this.forms[this.searchFormId];
			if (!R.isYummy(form.time) || !R.isYummy(form.proCode)) {
				this.toasted.warning("沒有选择日期或者项目");
				return;
			}
			const result = await this.$store.dispatch("EXPORT_IDENTIFY_STATISTICS", form);
			lpTools.createExcelFun(
				result,
				`${form.proCode}OCR识别统计${form.time[0]}-${form.time[1]}`
			);
			this.toasted.dynamic(`正在导出...`, result.code);
		}
	}
};
</script>
<style scoped lang="scss"></style>
