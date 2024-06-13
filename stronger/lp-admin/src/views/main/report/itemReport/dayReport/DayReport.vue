<template>
	<div class="day-report">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col
					v-for="(item, index) in cells.fields"
					:key="`entry_filters_${index}`"
					:cols="item.cols"
				>
					<template v-if="item.inputType === 'date'">
						<z-date-picker
							:formId="formId"
							:formKey="item.formKey"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
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
							:options="[{ label: '全部', value: 'all' }, ...auth.proItems]"
							:suffix="item.suffix"
						></z-select>
					</template>
				</v-col>

				<div class="z-flex">
					<z-btn
						:formId="formId"
						btnType="validate"
						class="pb-3 pl-3"
						color="primary"
						@click="search"
					>
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>
					<z-btn class="pb-3 pl-3" color="primary" @click="exportList"> 导出 </z-btn>
				</div>
			</v-row>
		</div>

		<div class="table day-report-table special" v-if="this.project === '1'">
			<p>业务量:</p>

			<vxe-grid v-if="this.type === 4" class="reverse-table" v-bind="volumeGrid"></vxe-grid>

			<vxe-table v-else :border="tableBorder" :data="desserts.volume" :size="tableSize">
				<template v-for="item in this.headers['volume']">
					<vxe-column
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
						align="center"
					></vxe-column>
				</template>
			</vxe-table>

			<div ref="volumeChart" class="main"></div>
		</div>

		<div class="table day-report-table special" v-if="this.project === '1'">
			<p>时效保障率:</p>

			<vxe-grid
				v-if="this.type === 4"
				class="reverse-table"
				v-bind="prescriptionGrid"
			></vxe-grid>

			<vxe-table v-else :border="tableBorder" :data="desserts.prescription" :size="tableSize">
				<template v-for="(item, i) in this.headers.prescription">
					<vxe-column
						:field="item.value"
						:title="item.text"
						:key="i"
						align="center"
						:width="item.width"
					>
					</vxe-column>
				</template>
			</vxe-table>
			<div ref="prescriptionChart" class="main"></div>
		</div>

		<div class="table day-report-table special" v-if="this.project === '1'">
			<p>处理时长:</p>
			<vxe-grid
				v-if="this.type === 4"
				class="reverse-table"
				v-bind="spendTimeGrid"
			></vxe-grid>

			<vxe-table v-else :border="tableBorder" :data="desserts.spendTime" :size="tableSize">
				<template v-for="(item, i) in this.headers.spendTime">
					<vxe-column
						:field="item.value"
						:title="item.text"
						:key="i"
						:width="item.width"
						align="center"
					></vxe-column>
				</template>
			</vxe-table>
			<div ref="spendTimeChart" class="main"></div>
		</div>

		<div class="table day-report-table" v-if="this.project !== '1'">
			<vxe-table :border="tableBorder" :data="dessertsDetail" :size="tableSize">
				<template v-for="(item, i) in this.headers">
					<vxe-column
						:field="item.value"
						:title="item.text"
						:key="i"
						align="center"
					></vxe-column>
				</template>
			</vxe-table>
			<div ref="QualityAccuracyChart" class="main"></div>
		</div>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import ChartMixins from "@/mixins/ChartMixins";
import { mapGetters } from "vuex";
import cells from "./cells";
import moment from "moment";

export default {
	name: "DayReport",
	mixins: [TableMixins, ChartMixins],

	data() {
		return {
			formId: "dayReport",
			headers: {
				volume: [],
				spendTime: [],
				prescription: []
			},
			volumeYearData: [],
			project: "1",
			cells,
			desserts: { volume: [{ project: "1" }] },
			dessertsDetail: [{ name: "业务量" }],
			chartRefList: [
				"volumeChart",
				"prescriptionChart",
				"spendTimeChart",
				"QualityAccuracyChart"
			],
			data: {},
			title: {
				prescription: "平均",
				spendTime: "平均",
				volume: "总量"
			},
			outValue: {
				prescription: "%",
				spendTime: "",
				volume: ""
			},
			volumeGrid: {
				border: true,
				showOverflow: true,
				showHeader: false,
				columns: [],
				data: []
			},
			prescriptionGrid: {
				border: true,
				showOverflow: true,
				showHeader: false,
				columns: [],
				data: []
			},
			spendTimeGrid: {
				border: true,
				showOverflow: true,
				showHeader: false,
				columns: [],
				data: []
			}
		};
	},
	mounted() {
		if (this.type === 3) {
			this.cells.fields[1].defaultValue = [
				moment().startOf("year").format("YYYY-MM-DD"),
				moment().format("YYYY-MM-DD")
			];
		}
		this.forms[this.formId].proCode = "all";
		//设置cartSet
		this.search();
	},
	props: ["type"],

	methods: {
		async exportList() {
			const form = this.forms[this.formId];
			const result = await this.$store.dispatch("EXPORT_REPORT_LIST", {
				startTime: form.date ? form.date[0] + "T00:00:00.000Z" : "",
				type: this.type || 1,
				endTime: form.date ? form.date[1] + "T23:59:59.000Z" : "",
				proCode: form.proCode || "all"
			});
		},
		async search() {
			const form = this.forms[this.formId];
			const result = await this.$store.dispatch("GET_REPORT_ALL_LIST", {
				startTime: form.date ? form.date[0] + "T00:00:00.000Z" : "",
				type: this.type || 1,
				endTime: form.date ? form.date[1] + "T23:59:59.000Z" : "",
				proCode: form.proCode || "all"
			});
			
			for (let key in result) {
				this.deal(key, result[key].data.list);
			}
		},
		deal(type, data) {
			this.getHeaders(type, data);
			let dataArr = [];
			let total = { project: this.title[type] };
			for (let key in data) {
				let dataObj = { project: key };
				data[key]?.forEach((e, i) => {
					if (total[e.countDate] !== undefined) {
						total[e.countDate] = Number(total[e.countDate]) + Number(e.count);
						if (type != "volume") {
							total[e.countDate] = total[e.countDate] / 2;
						}
					} else {
						total[e.countDate] = e.count;
					}
					dataObj[e.countDate] = e.count;
				});
				dataArr.push(dataObj);
			}

			console.log(dataArr,total,'ttt',type)

			dataArr.push(total);

			//反转年,业务表格
			if (this.type === 4) {
				this.reverseTable(dataArr, type);
			}

			this.setChart(type, dataArr);
			for (let key in dataArr[dataArr.length - 1]) {
				if (key != "project") dataArr[dataArr.length - 1][key] += this.outValue[type];
			}

			if (type === "spendTime") {
				dataArr.map(item => {
					if (item.project === "平均") {
						for (let key in item) {
							if (key !== "_X_ROW_KEY" && key !== "project") {
								//item[key] = Number(item[key]) >= 0 ? Number(item[key]).toFixed(0) : Number(item[key]).toFixed(2);
								item[key] = Number(item[key]).toFixed(2);
							}
						}
					}
				});
			}
			
			if (type === "prescription") {
				dataArr.map(item => {
					if (item.project !== "平均") {
						for (let key in item) {
							if (key !== "_X_ROW_KEY" && key !== "project") {
								item[key] = item[key] + '%';
							}
						}
					}
				});
			}
			this.desserts[type] = dataArr;
		},
		setChart(type, data) {
			const keyType = { 总量: 1, 平均: 1 };
			const outValue = this.outValue[type];
			const body = {
				dataNameList: [],
				xAxis: [],
				series: [],
				max: 0
			};
			const interval = this.switchType();
			let chartType = "";

			// 饼状图
			if (this.type === 4 && type === "volume") {
				chartType = "pie";
				let series = {
					type: "pie",
					radius: "50%",
					data: []
				};
				data.forEach(e => {
					if (keyType[e.project]) {
						return;
					}
					series.name = e["project"];
					series.data.push({
						name: e["project"]
					});

					for (let key in e) {
						if (key !== "project") {
							series.data.push(e[key]);
						}
					}

					body.series.push(series);
				});
			} else if (this.type === 4 && type === "prescription") {
				// 雷达图
				chartType = "radar";
				const body = {
					indicator: [],
					max: 100,
					min: 0,
					series: {
						name: "时效保障率",
						areaStyle: {
							color: {
								type: "radial",
								x: 0.5,
								y: 0.5,
								r: 0.5,
								colorStops: [
									{
										offset: 0,
										color: "rgba(137, 190, 252, 0.05)" // 0% 处的颜色
									},
									{
										offset: 1,
										color: "rgba(81, 137, 248, 0.3)" // 100% 处的颜色
									}
								]
							}
						},
						data: [
							{
								value: []
							}
						],
						type: chartType
					}
				};
				data.forEach(e => {
					if (keyType[e.project]) {
						return;
					}
					for (let key in e) {
						if (key !== "project") {
							body.max = Math.max(body.max, e[key]);
							body.series.data[0].value.push(e[key]);
						}
					}
					body.indicator.push({
						name: e.project,
						max: body.max,
						min: body.min
					});
				});

				this.setLineAndHistogramChart(type + "Chart", body, interval, chartType);
				return;
			} else if (this.type === 4 && type === "spendTime") {
				// 年报表处理时长
				chartType = "line";
				body.series = [
					{
						type: "line",
						data: [],
						barWidth: "2px",
						itemStyle: {
							color: "#6EC3C9"
						}
					}
				];
				data.forEach(e => {
					body.dataNameList.push(e.project);
					body.xAxis.push(e["project"]);
					for (let key in e) {
						if (key != "project") {
							body.series[0].data.push(e[key]);
						}
					}
				});
			} else {
				chartType = "";

				for (let key in data[0]) {
					if (key != "project") {
						body.xAxis.push(key);
					}
				}
				data.forEach(e => {
					body.dataNameList.push(e.project);
					let series = {
						name: e.project,
						type: keyType[e.project] ? "bar" : "line",
						tooltip: {
							valueFormatter: function (value) {
								return value + outValue;
							}
						},
						data: []
					};

					for (let key in e) {
						if (key !== "project") {
							body.max = Math.max(body.max, e[key]);
							series.data.push(e[key]);
						}
					}

					// 处理时长统计图
					if (type === "spendTime") {
						series.type = keyType[e.project] ? "line" : "bar";
					}

					//月报业务量时效保障率
					if (this.type === 3) {
						series.type = keyType[e.project] ? "line" : "bar";
					}

					// 日报时效保障换成折线图
					if (!this.type && type === "prescription") {
						series.smooth = keyType[e.project] ? false : true;
						series.type = "line";
					}

					// 日报处理时长增加数轴坐标
					if (!this.type && type === "spendTime") {
						body.yAxis = [
							{
								type: "value",
								name: "",
								min: body.min ?? 0,
								max: body.max ?? 300,
								interval: 0.5,
								axisLabel: {
									formatter: "{value}"
								}
							}
						];
					}

					body.series.push(series);
				});
			}
			this.setLineAndHistogramChart(type + "Chart", body, interval, chartType);
		},
		getHeaders(headerkey, data) {
			let defaultArr = [{ text: "项目", value: "project", width: 100 }];
			for (let key in data) {
				let newHeader = data[key]?.map(e => {
					return {
						text: e.countDate,
						value: e.countDate,
						width: 100
					};
				});
				if (newHeader.length > 0) {
					defaultArr = [...defaultArr, ...newHeader];
					this.headers[headerkey] = defaultArr;
				}
				break;
			}
		},
		reverseTable(myData, type) {
			const columns = this.headers[type];
			const buildData = columns.map(column => {
				const item = { col0: column.text };
				myData.forEach((row, index) => {
					item[`col${index + 1}`] = row[column.value];
				});
				return item;
			});
			const buildColumns = [
				{
					field: "col0",
					fixed: "left",
					width: 100
				}
			];
			myData.forEach((item, index) => {
				buildColumns.push({
					field: `col${index + 1}`,
					width: 100
				});
			});

			this[type + "Grid"].data = buildData;
			this[type + "Grid"].columns = buildColumns;
		},
		// 设置纵坐标间隔
		switchType() {
			let interval = 50;
			switch (this.type) {
				case 2:
					interval = 100;
					break;
				case 3:
					interval = 500;
					break;
				case 4:
					interval = 1000;
					break;
				default:
					interval = 50;
			}
			return interval;
		}
	},
	computed: {
		...mapGetters(["auth"])
	}
};
</script>

<style lang="scss">
.special {
	.reverse-table .vxe-body--row .vxe-body--column:first-child {
		background-color: #f8f8f9;
	}
}

.v-application p {
	color: red;
}

.main {
	min-width: 600px;
	min-height: 400px;
}
</style>
