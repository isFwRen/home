<template>
	<div class="return-analysis">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col
					v-for="(item, index) in cells.fields"
					:key="`entry_filters_${index}`"
					:cols="item.cols"
				>
					<template v-if="item.inputType === 'input'">
						<z-text-field
							:formId="searchFormId"
							:formKey="item.formKey"
							:clearable="item.clearable"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:suffix="item.suffix"
						>
						</z-text-field>
					</template>

					<template v-else-if="item.inputType === 'date'">
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
					<vxe-colgroup
						v-if="
							item.value === 'ReturnSituation' ||
							item.value === 'ReturnProportion' ||
							item.value === 'project2'
						"
						align="center"
						:title="item.text"
						:key="item.value"
					>
						<template v-for="record in item.children">
							<vxe-column
								:field="record.value"
								:title="record.text"
								:key="record.value"
							></vxe-column>
						</template>
					</vxe-colgroup>

					<vxe-column
						v-else
						:field="item.value"
						:title="item.text"
						:key="item.value"
					></vxe-column>
				</template>
			</vxe-table>
		</div>
	</div>
</template>

<script>
import { tools } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import SocketsMixins from "@/mixins/SocketsMixins";
import SpecialReportMixins from "../SpecialReportMixins";
import cells from "./cells";
import { copy, copyText } from "clipboard-vue";
import { tools as lpTools } from "@/libs/util";

export default {
	name: "ReturnAnalysis",
	mixins: [TableMixins, SocketsMixins, SpecialReportMixins],

	data() {
		return {
			formId: "ReturnAnalysis",
			project: "1",
			cells,
			dispatchList: "GET_RETRURN_ANALYSIS_LIST",
			headers: cells.headers1,
			socketPath: "ExportBusinessUploadAnalysis"
		};
	},

	methods: {
		// 复制到粘貼板
		onCopy() {
			var removeKey = [];

			if (this.params.proCode.indexOf("B") !== -1) {
				removeKey = ["_X_ID"];
			} else {
				removeKey = ["submitTime", "createAt", "_X_ID"];
			}

			const dataArr = [];
			const rowTopOne = [];
			const rowTopTwo = [];

			for (var i = 0; i < this.headers.length; i++) {
				if (this.headers[i].children == undefined) {
					rowTopOne.push('"' + this.headers[i].text + '"');
					rowTopTwo.push('"' + this.headers[i].text + '"');
				} else {
					for (var j = 0; j < this.headers[i].children.length; j++) {
						if (j == 0) {
							rowTopOne.push('"' + this.headers[i].text + '"');
							rowTopTwo.push('"' + this.headers[i].children[j].text + '"');
							continue;
						}
						rowTopOne.push('"' + this.headers[i].children[j].text + '"');
						rowTopTwo.push('"' + this.headers[i].children[j].text + '"');
					}
				}
			}

			if (rowTopOne.length != 0) {
				for (var k = 0; k < rowTopOne.length; k++) {
					if (rowTopOne[k] == rowTopTwo[k]) {
						rowTopOne[k] = "";
					}
				}
				dataArr.push(rowTopOne.join("\t"));
			}
			dataArr.push(rowTopTwo.join("\t"));

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
							if (element[key] == "upload") {
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

		onChangeType(value) {
			if (value === "全部/整体") {
				this.headers = cells.headers1;
			} else if (value === "全部/明细") {
				this.headers = cells.headers2;
			} else if (value.indexOf("B") !== -1 && value.indexOf("整体") !== -1) {
				this.headers = cells.headers3;
			} else {
				this.headers = cells.headers4;
			}
		},

		//导出
		async onExport() {
			const form = this.forms[this.searchFormId];

			if (!tools.isYummy(form.time) || !tools.isYummy(form.proCode)) {
				this.toasted.warning("请选择日期或者项目");
				return;
			}

			const result = await this.$store.dispatch("EXPORT_RETURNANALYSIS_DETAIL", form);
			lpTools.createExcelFun(result, `${form.proCode}回传分析${form.time[0]}-${form.time[1]}`)
			this.toasted.dynamic(`正在导出...`, result.code);
		}
	},

	directives: {
		copy
	}
};
</script>
