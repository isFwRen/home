<template>
	<div class="error-analysis">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col
					v-for="(item, index) in cells.fields"
					:key="`entry_filters_${index}`"
					:cols="item.cols"
				>
					<template v-if="item.inputType === 'input'">
						<z-text-field
							:formId="formId"
							:formKey="item.formKey"
							:clearable="item.clearable"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:suffix="item.suffix"
							:defaultValue="item.defaultValue"
						>
						</z-text-field>
					</template>

					<template v-else-if="item.inputType === 'date'">
						<z-date-picker
							:formId="formId"
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
							:formId="formId"
							:formKey="item.formKey"
							:clearable="item.clearable"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:options="item.options"
							:suffix="item.suffix"
							:defaultValue="item.defaultValue"
							@change="onChangeType($event, item)"
						></z-select>
					</template>
				</v-col>

				<div class="z-flex">
					<z-btn
						:formId="formId"
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
			<z-btn
				v-for="item of cells.btns"
				:key="item.icon"
				:class="item.class"
				:color="item.color"
				small
				outlined
				@click="onExport"
			>
				{{ item.text }}
			</z-btn>
		</div>

		<div class="table error-analysis-table">
			<vxe-table :border="tableBorder" :data="desserts" :size="tableSize">
				<template v-for="item in cells.headers">
					<vxe-column
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
import { mapGetters } from "vuex";
import { tools } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import SocketsMixins from "@/mixins/SocketsMixins";
import cells from "./cells";

export default {
	name: "ErrorAnalysis",
	mixins: [TableMixins, SocketsMixins],

	data() {
		return {
			formId: "ErrorAnalysis",
			cells,
			socketPath: "incorrectAnalysis"
		};
	},

	methods: {
		onChangeType(value, item) {
			console.log("asdasd", value, item);
		},

		onSearch() {
			const form = this.forms[this.formId];
			console.log("form", form);
			const params = {
				proCode: form.proCode,
				startTime: form.date[0] + " 00:00:00",
				endTime: form.date[1] + " 23:59:59",
				code: form.code,
				nickName: form.nickName
			};
			this.getErrorAnalysis(params);
		},

		//获取错误分析
		async getErrorAnalysis(params) {
			const result = await this.$store.dispatch("GET_ERRORANALYSIS_ITEM_LIST", params);
			console.log("错误分析-views-report-errorAnalysis", result);
			this.desserts = result.data.list;
		},
		//导出错误明细
		async onExport() {
			const form = this.forms[this.formId];
			console.log(form);
			if (!tools.isYummy(form.date) || !tools.isYummy(form.proCode)) {
				this.toasted.dynamic("沒有选择日期或者项目", 400);
				return;
			}
			const result = await this.$store.dispatch("EXPORT_ERRORANALYSIS_DETAIL", form);
			this.toasted.dynamic(`${result.msg}，正在导出`, result.code);
		}
	},

	computed: {
		...mapGetters(["auth"])
	},

	watch: {
		"auth.proItems": {
			handler(items) {
				this.cells.fields[0].options = items;
			},
			deep: true,
			immediate: true
		}
	}
};
</script>
