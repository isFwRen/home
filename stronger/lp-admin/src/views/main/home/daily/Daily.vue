<template>
	<div class="daily">
		<div class="menu">
			<z-btn-toggle
				class="pa-0 ma-0"
				formId="sexuals"
				formKey="sexual"
				color="primary"
				:group="true"
				mandatory
				:options="cells.menuList"
				@click="changeType"
			></z-btn-toggle>
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
								<template #default="{ row }">
									<span>{{
										(item.output && item.output(row[record.value])) ||
										row[record.value]
									}}</span>
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
			formId: "projectReport",
			cells,
			chooseDay: moment().subtract(1, "day").format("YYYY-MM-DD"),
			time: moment().subtract(1, "day").format("YYYY-MM-DD"),
			tableRow: 0,
			desserts: [
				{
					pridictValue: "5525",
					name: "B0101民生理赔",
					actualValue: "1722.0",
					timeProportion: "3.33%",
					completeProportion: "4.65%"
				},
				{
					pridictValue: "64000",
					name: "合计",
					actualValue: "1715.0",
					timeProportion: "3.33%",
					completeProportion: "2.68%"
				},
				{
					pridictValue: "编制人数",
					name: "人员情况",
					actualValue: "20",
					monthVolume: "实到人数",
					dayVolume: "10",
					dayOvertimeNum: "下班时间",
					monthAccuracy: "24:00"
				},
				{
					pridictValue: "2525",
					name: "其他运行情况",
					actualValue: "153.2522525"
				}
			],
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
					monthQuality: "0"
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
			this.desserts.forEach(ele => {
				let temp = ele.timePercent?.toString().split("%")[0];
				let tempFinishPercent = ele.finishPercent?.toString().split("%")[0];
				let tempMonthPercent = ele.monthAgingPercent?.toString().split("%")[0];
				let tempMonthRightPercent = ele.monthRightPercent?.toString().split("%")[0];
				ele.timePercent = Number(temp).toFixed(2) + "%";
				ele.finishPercent =
					Number(tempFinishPercent) === 0
						? Number(tempFinishPercent).toFixed(0) + "%"
						: Number(tempFinishPercent).toFixed(2) + "%";
				ele.monthAgingPercent =
					Number(tempMonthPercent) === 0
						? Number(tempMonthPercent).toFixed(0) + "%"
						: Number(tempMonthPercent).toFixed(2) + "%";
				ele.monthRightPercent =
					Number(tempMonthRightPercent) === 0
						? Number(tempMonthRightPercent).toFixed(0) + "%"
						: Number(tempMonthRightPercent).toFixed(2) + "%";
			});
		}
	},
	created() {
		this.search();
	},
	methods: {
		async search() {
			const result = await this.$store.dispatch("GET_HOMRPAGE_PRO_REPORT_LIST", {
				reportDay: this.time + "T00:00:00.000Z"
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
					e[key] = e[key] / result.data.list.length + "%";
				});

				mydesserts.push(e);
			});
			precentArr.forEach(key => {
				total[key] = Number(total[key]).toFixed(2);
				total[key] += "%";
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
		changeType(e) {
			this.time =
				e == 0
					? moment().subtract(1, "day").format("YYYY-MM-DD")
					: moment().subtract(2, "day").format("YYYY-MM-DD");
			this.search();
			this.chooseDay = this.time;
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
	border-style: solid;
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
