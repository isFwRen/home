<template>
	<div class="project-report">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col cols="2">
					<z-date-picker
						:formId="searchFormId"
						formKey="reportDay"
						:hideDetails="true"
						label="日期"
						prepend-icon="mdi-calendar"
						:immediate="true"
						:defaultValue="cells.fields[0].defaultValue"
					></z-date-picker>
				</v-col>

				<div class="z-flex">
					<z-btn class="pb-3 pl-3" color="primary" @click="search">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>
					<z-btn class="pb-3 pl-3" color="primary" @click="edit"> 保存 </z-btn>
					<z-btn class="pb-3 pl-3" color="primary" @click="exportReport"> 导出 </z-btn>
				</div>
			</v-row>
		</div>

		<div class="table project-report-table">
			<h1 class="table_title">{{ this.chooseDay }}理赔项目日报</h1>
			<vxe-table
				:border="tableBorder"
				:data="desserts"
				:merge-cells="mergeCells"
				:size="tableSize"
			>
				<template v-for="item in cells.headers">
					<vxe-colgroup
						v-if="
							item.value === 'ouputValue' ||
							item.value === 'businessVolume' ||
							item.value === 'prescription' ||
							item.value === 'quality'
						"
						align="center"
						:title="item.text"
						:key="item.value"
					>
						<template v-for="record in item.children">
							<vxe-column :field="record.value" :title="record.text" align="center">
								<template
									#default="{ row }"
									v-if="
										record.value === 'predictValue' ||
										record.value === 'finishValue' ||
										record.value === 'dayCount' ||
										record.value === 'monthRightPercent'
									"
								>
									<template
										v-if="record.value === 'predictValue' && row.id == 'b-1'"
									>
										<z-text-field
											:formId="formId"
											formKey="otherMess"
											:defaultValue="row[record.value]"
											:disabled="ableDisabled"
										>
										</z-text-field>
									</template>
									<template
										v-else-if="
											record.value === 'predictValue' && !unEdit[row.name]
										"
									>
										<z-text-field
											:formId="PrformId"
											:formKey="row.proCode"
											:defaultValue="row.predictValue"
											:disabled="ableDisabled"
											:validation="[
												{
													rule: 'numeric',
													message: '允许字段为正整数.'
												}
											]"
										>
										</z-text-field>
									</template>

									<template
										v-else-if="
											record.value === 'finishValue' && row.id == 'b-2'
										"
									>
										<z-text-field
											:formId="formId"
											formKey="userCount"
											:defaultValue="row[record.value]"
											:disabled="ableDisabled"
											:validation="[
												{
													rule: 'numeric',
													message: '允许字段为正整数.'
												}
											]"
										>
										</z-text-field>
									</template>

									<template
										v-else-if="record.value === 'dayCount' && row.id == 'b-2'"
									>
										<z-text-field
											:formId="formId"
											formKey="activeUserCount"
											:defaultValue="row[record.value]"
											:disabled="ableDisabled"
											:validation="[
												{
													rule: 'numeric',
													message: '允许字段为正整数.'
												}
											]"
										>
										</z-text-field>
									</template>

									<template
										v-else-if="
											record.value === 'monthRightPercent' && row.id == 'b-2'
										"
									>
										<z-date-picker
											:formId="formId"
											formKey="closingTime"
											format="24hr"
											:immediate="false"
											mode="time"
											prepend-icon="mdi-alarm"
											time-use-seconds
											time-format="24hr"
											:defaultValue="row[record.value]"
											:disabled="ableDisabled"
										></z-date-picker>
									</template>
									<template v-else>
										<span>{{
											(item.output && item.output(row[record.value])) ||
											row[record.value]
										}}</span>
									</template>
								</template>
							</vxe-column>
						</template>
					</vxe-colgroup>

					<vxe-column
						v-else
						:field="item.value"
						:title="item.text"
						width="90px"
						align="center"
					>
						<template #default="{ row }">
							<span>{{
								(item.output && item.output(row[item.value])) || row[item.value]
							}}</span>
						</template>
					</vxe-column>
				</template>
			</vxe-table>
		</div>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import moment from "moment";
import cells from "./cells";

export default {
	name: "ProjectReport",
	mixins: [TableMixins],

	data() {
		return {
			searchFormId: "dayReprotList",
			PrformId: "prformID",
			formId: "projectReport",
			cells,
			chooseDay: "",
			tableRow: 0,
			ableDisabled: false,
			unEdit: { 合计: 1, 人员情况: 1 },
			beEdit: {
				finishValue: 1,
				dayCount: 1,
				monthQuality: 1
			},
			desserts: [],
			defaultDesserts: [
				{
					id: "b-3",
					predictValue: "64000",
					name: "合计",
					finishValue: "0",
					timeProportion: "0%",
					completeProportion: "0%"
				},
				{
					id: "b-2",
					predictValue: "编制人数",
					name: "人员情况",
					finishValue: "0",
					monthCount: "实到人数",
					dayCount: "0",
					dayTimeoutCount: "下班时间",
					monthRightPercent: "0"
				},
				{
					id: "b-1",
					predictValue: ".",
					name: "其他运行情况"
				}
			],
			mergeCells: []
		};
	},
	watch: {
		desserts(val) {
			const deno = val.length > 0 ? val.length - 3 : 1;
			val.map(item => {
				if (item.name === "合计") {
					item.finishValue =
						item.finishValue == 0 ? 0 : (item.finishValue / deno).toFixed(2);
					const timePercentNum = Number(item.timePercent.split("%")[0]) / deno;
					const finishPercentNum = Number(item.finishPercent.split("%")[0]) / deno;
					const monthAgingPercentNum =
						Number(item.monthAgingPercent.split("%")[0]) / deno;
					const monthRightPercentNum =
						Number(item.monthRightPercent.split("%")[0]) / deno;

					item.monthAgingPercent =
						monthAgingPercentNum == 0
							? monthAgingPercentNum.toFixed(0) + "%"
							: monthAgingPercentNum.toFixed(2) + "%";

					item.monthRightPercent =
						monthRightPercentNum == 0
							? monthRightPercentNum.toFixed(0) + "%"
							: monthRightPercentNum.toFixed(2) + "%";
					item.finishPercent =
						finishPercentNum == 0
							? finishPercentNum.toFixed(0) + "%"
							: finishPercentNum.toFixed(2) + "%";
					item.timePercent =
						timePercentNum == 0
							? timePercentNum.toFixed(0) + "%"
							: timePercentNum.toFixed(2) + "%";
				}
			});
		}
	},
	created() {
		this.search();
	},
	methods: {
		async exportReport() {
			const form = this.forms[this.searchFormId];
			this.chooseDay =
				form?.reportDay.length == 1
					? form?.reportDay[0]
					: form?.reportDay || moment().subtract(1, "day").format("YYYY-MM-DD");
			const result = await this.$store.dispatch("EXPORT_COMMOUNT_FILE", {
				url: "report-management/project-report/report-info-export",
				data: {
					reportDay:
						(form?.reportDay && form.reportDay + "T16:00:00.000Z") ||
						moment().subtract(1, "day").format("YYYY-MM-DD") + "T16:00:00.000Z"
				},
				fileName: form.reportDay + "理赔项目日报.xlsx"
			});
		},
		async edit() {
			this.$modal({
				visible: true,
				title: "修改提示",
				content: "请确认是否要保存？",
				confirm: () => {
					const form = this.forms[this.formId];
					const proReportOtherInfo = {
						...form,
						reportDate: this.chooseDay + "T16:00:00.000Z",
						userCount: +form.userCount,
						activeUserCount: +form.activeUserCount
					};
					const pr = this.forms[this.PrformId];
					const proReport = [];
					for (let key in pr) {
						proReport.push({
							proCode: key,
							predictValue: +pr[key],
							reportDate: this.chooseDay + "T16:00:00.000Z"
						});
					}
					const data = {
						proReport,
						proReportOtherInfo
					};
					this.saveEdit(data);
				}
			});
		},
		async saveEdit(data) {
			const result = await this.$store.dispatch("SAVE_REPORT_PROJECT_EDIT", data);

			if (result.code == 200) {
				this.toasted.dynamic(result.msg, result.code);
				this.search();
			} else {
				this.toasted.dynamic("保存失败", result.code);
			}
		},
		async search() {
			const form = this.forms[this.searchFormId];
			this.chooseDay =
				form?.reportDay.length == 1
					? form?.reportDay[0]
					: form?.reportDay || moment().subtract(1, "day").format("YYYY-MM-DD");
			const result = await this.$store.dispatch("GET_HOMRPAGE_PRO_REPORT_LIST", {
				reportDay: this.chooseDay + "T16:00:00.000Z"
			});

			let mydesserts = [];
			let total = {};
			const precentArr = [
				"timePercent",
				"finishPercent",
				"monthAgingPercent",
				"monthRightPercent"
			];
			result.data.list.forEach(e => {
				e.name = e.proCode;
				precentArr.forEach(key => {
					e[key] *= 100;
				});
				for (let key in e) {
					if (typeof e[key] == "number") {
						if (total[key]) {
							total[key] += e[key];
						} else {
							total[key] = e[key];
						}
					}
				}

				precentArr.forEach(key => {
					e[key] = Number(e[key]).toFixed(2);
					let num = Number(e[key].split("%")[0]);
					//e[key] = e[key] / result.data.list.length + '%'
					e[key] = num == 0 ? num + "%" : Number(e[key].split("%")[0]).toFixed(2) + "%";
				});

				mydesserts.push(e);
			});

			precentArr.forEach(key => {
				total[key] = Number(total[key]).toFixed(2);
				let num = Number(total[key].split("%")[0]);
				total[key] =
					num == 0 ? num + "%" : Number(total[key].split("%")[0]).toFixed(2) + "%";
			});
			//这是一条生肉注释
			const otherInfo = result.data.otherInfo;
			this.defaultDesserts[1].finishValue = otherInfo.userCount;
			this.defaultDesserts[1].dayCount = otherInfo.activeUserCount;
			this.defaultDesserts[1].monthRightPercent = otherInfo.closingTime;
			this.defaultDesserts[2].predictValue = otherInfo.otherMess;
			this.defaultDesserts[0] = { name: "合计", ...total };
			this.desserts = [...mydesserts, ...this.defaultDesserts];
			this.getRow();
		},
		getRow() {
			this.tableRow = this.desserts.length;
			this.mergeCells = [
				{
					row: this.tableRow - 2,
					col: 2,
					rowspan: 1,
					colspan: 3
				},
				{
					row: this.tableRow - 2,
					col: 6,
					rowspan: 1,
					colspan: 3
				},
				{
					row: this.tableRow - 2,
					col: 10,
					rowspan: 1,
					colspan: 3
				},
				{
					row: this.tableRow - 1,
					col: 1,
					rowspan: 3,
					colspan: 12
				}
			];
		}
	}
};
</script>

<style scoped>
.vxe-table--render-default .vxe-cell {
	white-space: pre-line;
	word-break: break-all;
	padding: 0;
}
.table input {
	margin: 0;
	padding: 0;
	width: 100%;
	border-style: none;
}
.other_input {
	height: 100px;
}
.table_title {
	text-align: center;
	color: rgb(47, 110, 139);
	font-size: 20px;
	font-weight: bold;
}
</style>
