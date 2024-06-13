<template>
	<div class="abnormal-part">
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
							@change="handleChange($event, item)"
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
			<z-btn class="pr-3" color="primary" outlined small @click="onCopy"> 复制 </z-btn>

			<z-btn class="pr-3" color="primary" outlined small @click="onExport"> 导出 </z-btn>
		</div>

		<div class="table abnormal-part-table">
			<vxe-table :border="tableBorder" :data="desserts" :size="tableSize">
				<template v-for="item in cells.headers">
					<vxe-column
						:field="item.value"
						:title="item.text"
						:key="item.value"
					></vxe-column>
				</template>
			</vxe-table>

			<z-pagination
				:total="pagination.total"
				:options="pageSizes"
				@page="handlePage"
			></z-pagination>
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
import { mapState } from "vuex";

export default {
	name: "AbnormalPart",
	mixins: [TableMixins, SocketsMixins, SpecialReportMixins],

	data() {
		return {
			formId: "AbnormalPart",
			dispatchList: "GET_ABNORMAL_PART_LIST",
			cells,
			socketPath: "exportAbnormalBill"
		};
	},
	computed: {
		...mapState(["forms"])
	},
	methods: {
		onCopy() {
			const copyHeader = [];
			this.cells.headers.forEach(element => {
				if (element.children) {
					element.children.forEach(item => {
						copyHeader.push(item.value);
					});
				} else {
					copyHeader.push(element.value);
				}
			});

			const dataArr = [];

			const rowTopOne = [];

			for (var i = 0; i < this.cells.headers.length; i++) {
				rowTopOne.push('"' + this.cells.headers[i].text + '"');
			}

			dataArr.push(rowTopOne.join("\t"));

			this.desserts.forEach(element => {
				const rowArr = [];
				for (var i = 0; i < this.cells.headers.length; i++) {
					rowArr.push('"' + element[this.cells.headers[i].value] + '"');
				}
				dataArr.push(rowArr.join("\t"));
			});

			copyText(dataArr.join("\n"))
				.then(e => {
					this.toasted.success("复制成功");
				})
				.catch(e => {
					this.toasted.error("复制失败");
				});
		},

		// 导出
		async onExport() {
			const form = this.forms[this.searchFormId];

			if (!R.isYummy(form.time) || !R.isYummy(form.proCode)) {
				this.toasted.warning("沒有选择日期或者项目");
				return;
			}

			const result = await this.$store.dispatch("EXPORT_ABNORMAL_PART_BILL", form);
			lpTools.createExcelFun(result[0], `${form.proCode}异常件数据${form.time[0]}-${form.time[1]}`)
			this.toasted.dynamic(`正在导出...`, result.code);
		}
	}
};
</script>
