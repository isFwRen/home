<template>
	<div class="destruction-report">
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
						class="pb-3 pl-3"
						color="primary"
						@click="onSearch"
					>
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>
					<z-btn
						:formId="searchFormId"
						class="pb-3 pl-3"
						color="success"
						@click="onExcel"
					>
						导出明细
					</z-btn>

					<z-btn :formId="searchFormId" class="pb-3 pl-3" color="error" @click="onWord">
						导出销毁报告
					</z-btn>

					<z-btn
						v-if="permission.includes(userCode)"
						:formId="searchFormId"
						class="pb-3 pl-3"
						color="orange"
						@click="getHistory"
					>
						获取历史明细
					</z-btn>
					<span v-if="isShow" class="tip">获取历史明细成功</span>
				</div>
			</v-row>
		</div>

		<!--table-->
		<div class="table staff-total-table">
			<vxe-table :border="tableBorder" :data="desserts" :size="tableSize" align="center">
				<vxe-table-column type="seq" width="60" title="序号"></vxe-table-column>
				<template v-for="item in cells.headers">
					<vxe-column
						v-if="item.value === 'scanAt'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
						:sortable="item.sortable"
						:fixed="item.fixed"
					>
						<template #default="{ row }">
							{{ transformTime(row.scanAt) }}
						</template>
					</vxe-column>
					<vxe-column
						v-else-if="item.value === 'delAt'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
						:sortable="item.sortable"
						:fixed="item.fixed"
					>
						<template #default="{ row }">
							{{ transformTime(row.delAt) }}
						</template>
					</vxe-column>
					<vxe-column v-else :field="item.value" :title="item.text" :key="item.value">
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
import { R, sessionStorage, localStorage } from "vue-rocket";
import { mapState } from "vuex";
import { tools as lpTools } from "@/libs/util";
import moment from "moment";

export default {
	name: "destructionReport",
	mixins: [SpecialReportMixins, TableMixins],
	data() {
		return {
			cells,
			dispatchList: "GET_DESTRUCTION_REPORT_LIST",
			formId: "destructionReport",
			dataSource: cells.dataSource,
			dialog: false,
			tableData: [],
			permission: ["6433", "5774", "7265"],
			userCode: "",
			isShow: false,
		};
	},

	async created() {
		this.userCode = localStorage.get("user").code;
	},

	computed: {
		...mapState(["forms"]),
		transformTime() {
			return times => moment(times).format("YYYY-MM-DD HH:mm:ss");
		}
	},
	methods: {
		async onExcel() {
			const form = this.forms[this.searchFormId];

			if (!R.isYummy(form.time) || !R.isYummy(form.proCode)) {
				this.toasted.warning("沒有选择日期或者项目");
				return;
			}

			const result = await this.$store.dispatch("EXPORT_DESTRUCTION_REPORT_EXCEL", form);

			lpTools.createExcelFun(
				result,
				`${form.proCode}销毁明细${form.time[0]}-${form.time[1]}`
			);
		},

		async onWord() {
			let wordContent = await this.getWord();
			const form = this.forms[this.searchFormId];
			let year = form.time[0].slice(0, 4);
			let month = form.time[0].slice(5, 7);
			let startTime = form.time[0].replace("-", "年").replace("-", "月") + "日";
			let endTime = form.time[1].replace("-", "年").replace("-", "月") + "日";
			sessionStorage.set("wordContent", wordContent);
			sessionStorage.set("year", year);
			sessionStorage.set("month", month);
			sessionStorage.set("startTime", startTime);
			sessionStorage.set("endTime", endTime);
			sessionStorage.set("proCode", form.proCode);
			// cells.tableData[1].col2 = `自${time1} 00 ： 00 至 ${time2} 24 ： 00`;
			// cells.tableData[2].col2 = this.wordContent.dNum;
			// cells.tableData[3].col2 = this.wordContent.uNum;
			// cells.tableData[4].col2 = (this.wordContent.total / 1024).toFixed(2);

			// this.tableData = cells.tableData;
			// this.dialog = true;
			window.open(
				`${location.origin}/normal/des-report`,
				"_blank",
				"width=" +
					(window.screen.availWidth - 10) +
					",height=" +
					(window.screen.availHeight - 30) +
					",top=0,left=0,toolbar=no,menubar=no,scrollbars=no, resizable=no,location=no, status=no"
			);
		},

		async getHistory() {
			this.isShow = false;
			const form = this.forms[this.searchFormId];
			const result = await this.$store.dispatch("GET_HISTORY_DETAIL", form);
			if (result && result.size && result.type == "application/json") {
				this.isShow = true;
				this.toasted.warning("操作成功");
			}
		},

		//获取后端blog内容
		async getWord() {
			const form = this.forms[this.searchFormId];

			if (!R.isYummy(form.time) || !R.isYummy(form.proCode)) {
				this.toasted.warning("沒有选择日期或者项目");
				return;
			}
			const result = await this.$store.dispatch("EXPORT_DESTRUCTION_REPORT_WORD", form);

			return new Promise((resolve, reject) => {
				const file = new FileReader();
				file.readAsText(result, "utf-8");
				file.onload = function () {
					const obj = JSON.parse(file.result);
					resolve(obj.data);
				};
			});
		},

		printWord() {
			this.dialog = false;
			this.$refs.xTable.print();
		}
	}
};
</script>
<style scoped lang="scss">
::v-deep .theme--light.v-btn {
	color: white !important;
}

.tip {
	display: inline-block;
	height: 36px;
	line-height: 36px;
	margin-left: 20px;
	color: red;
	font-weight: 900;
	font-size: 20px;
}
</style>
