<template>
	<div class="income-analysis">
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
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:clearable="item.clearable"
							:options="item.options"
							:suffix="item.suffix"
							@change="onChangeType($event, item)"
						></z-select>
					</template>
				</v-col>

				<div class="z-flex">
					<z-btn class="pb-3 pl-3" color="primary" @click="onSearch">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>
				</div>
			</v-row>
		</div>

		<div class="pb-6 btns">
			<z-btn class="pr-3" color="primary" small outlined @click="onCopy"> 复制 </z-btn>
			<z-btn class="pr-3" color="primary" small outlined @click="onExport"> 导出 </z-btn>
		</div>

		<div class="table staff-total-table">
			<vxe-table :border="tableBorder" :data="desserts" :size="tableSize">
				<template v-for="item in this.headers">
					<vxe-column :field="item.value" :title="item.text" :key="item.value">
					</vxe-column>
				</template>
			</vxe-table>
		</div>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import SocketsMixins from "@/mixins/SocketsMixins";
import SpecialReportMixins from "../SpecialReportMixins";
import cells from "./cells";
import { R } from "vue-rocket";
import { copy, copyText } from "clipboard-vue";
import { tools as lpTools } from "@/libs/util";

export default {
	name: "IncomeAnalysis",
	mixins: [TableMixins, SocketsMixins, SpecialReportMixins],

	data() {
		return {
			formId: "IncomeAnalysis",
			project: "1",
			cells,
			dispatchList: "GET_INCOME_ANALYSIS_LIST",
			headers: [],
			socketPath: "BusinessDownloadAnalysis"
		};
	},

	methods: {
		onChangeType(value) {
			if (value === "全部") {
				this.headers = cells.headers;
			} else {
				this.headers = cells.headers2;
			}
		},
		// 复制到粘貼板
		onCopy() {
			console.log("123", this.headers);

			const removeKey = ["submitTime", "proCode", "_X_ID"];

			const dataArr = [];

			const rowTopOne = [];

			for (var i = 0; i < this.headers.length; i++) {
				rowTopOne.push('"' + this.headers[i].text + '"');
			}
			dataArr.push(rowTopOne.join("\t"));

			this.desserts.forEach(element => {
				const rowArr = [];
				for (const key of Object.keys(element)) {
					if (key == "proSummary" && element[key].constructor == Array) {
						element[key].forEach(item => {
							for (const k of Object.keys(item)) {
								if (removeKey.indexOf(k) == -1) {
									rowArr.push('"' + item[k] + '"');
								}
							}
						});
					} else {
						if (removeKey.indexOf(key) == -1) {
							if (element[key] == "download") {
								continue;
							}
							if (/^row/.test(element[key])) {
								continue;
							}
							rowArr.push('"' + element[key] + '"');
						}
					}
				}
				dataArr.push(rowArr.join("\t"));
			});
			if (dataArr.length == 0) {
				this.toasted.warning("没有查询到数据");
			}
			copyText(dataArr.join("\n"))
				.then(e => {
					this.toasted.success("复制成功");
				})
				.catch(e => {
					this.toasted.error("复制失败");
				});
		},

		//导出
		async onExport() {
			const form = this.forms[this.searchFormId];

			if (!R.isYummy(form.time)) {
				this.toasted.warning("沒有选择日期或者项目");
				return;
			}

			const result = await this.$store.dispatch("EXPORT_INCOME_ANALYSIS_DETAIL", form);
			lpTools.createExcelFun(result, `${form.proCode}来量分析${form.time[0]}-${form.time[1]}`)
			this.toasted.dynamic(`正在导出...`, result.code);
		}
	},

	directives: {
		copy
	}
};
</script>
