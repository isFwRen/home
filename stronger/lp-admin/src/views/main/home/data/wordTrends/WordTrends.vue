<template>
	<div class="businessTrends">
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
		<div class="treeBox">
			<HomeMinTitle text=" ">
				<a
					:href="
						rankingType == 0
							? '../report/itemReport/day-report'
							: '../report/itemReport/month-report'
					"
				>
					<div class="more" @click="showMore('more')">
						<span>详情>></span>
					</div>
				</a>
			</HomeMinTitle>
			<div class="mb-4"></div>
			<div class="statistics">
				<div ref="treeView" class="treeView"></div>
			</div>
		</div>
		<div class="tableBox">
			<HomeMinTitle text="保障排行榜" />
			<div class="mb-4"></div>
			<div class="table">
				<vxe-table :data="desserts" :border="tableBorder" :size="tableSize">
					<template v-for="item in cells.headers">
						<vxe-column :field="item.value" :title="item.text" :width="item.width">
							<template #default="{ row }">
								<span>{{
									(item.output && item.output(row[item.value])) || row[item.value]
								}}</span>
							</template>
						</vxe-column>
					</template>
				</vxe-table>
			</div>
			<z-pagination
				:options="pageSizes"
				@page="handlePage"
				:total="pagination.total"
			></z-pagination>
		</div>
	</div>
</template>

<script>
import HomeMinTitle from "../../homepage/HomeMinTitle.vue";
import cells from "./cells";
import TableMixins from "@/mixins/TableMixins";
import charMixins from "@/mixins/ChartMixins";
export default {
	mixins: [TableMixins, charMixins],
	components: { HomeMinTitle },
	data() {
		return {
			cells,
			rankingType: 0,
			desserts: [],
			chartBody: {
				xAxis: [],
				series: [],
				yAxis: [
					{
						type: "value",
						name: "",
						min: 0,
						max: 300,
						interval: 50,
						axisLabel: {
							formatter: "{value}"
						}
					},
					{
						type: "value",
						name: "",
						min: 0,
						max: 100,
						interval: 50,
						axisLabel: {
							formatter: "{value}"
						}
					}
				]
			}
		};
	},
	computed: {
		watchDay() {
			return this.rankingType == 0 ? ["今日", "昨日"] : ["本月", "上月"];
		}
	},
	methods: {
		async search() {
			const result = await this.$store.dispatch("GET_HOME_AGING_TREND_LIST", {
				rankingType: this.rankingType
			});
			if (result.code == 200) {
				this.chartBody.dataNameist = [
					this.watchDay[0] + "超时票数",
					this.watchDay[1] + "超时票数",
					this.watchDay[0] + "保障率",
					this.watchDay[1] + "保障率"
				];
				result.data.list.sort(function (a, b) {
					return (
						Math.max(b.timeoutCount, b.yesterdayTimeoutCount) -
						Math.max(a.timeoutCount, a.yesterdayTimeoutCount)
					);
				});
				let todayCount = {
					name: this.watchDay[0] + "超时票数",
					type: "bar",
					tooltip: {
						valueFormatter: function (value) {
							return value;
						}
					},
					data: []
				};
				let todayPercent = {
						name: this.watchDay[0] + "保障率",
						yAxisIndex: 1,
						type: "line",
						tooltip: {
							valueFormatter: function (value) {
								return value + "%";
							}
						},
						data: []
					},
					yesterdayBillCount = {
						name: this.watchDay[1] + "超时票数",
						yAxisIndex: 0,
						type: "bar",
						tooltip: {
							valueFormatter: function (value) {
								return value;
							}
						},
						data: []
					},
					yesterdayBillPercent = {
						name: this.watchDay[1] + "保障率",
						yAxisIndex: 1,
						type: "line",
						tooltip: {
							valueFormatter: function (value) {
								return value + "%";
							}
						},
						data: []
					};
				this.chartBody.yAxis[0].max = Math.max(
					result.data.list[0].timeoutCount,
					result.data.list[0].yesterdayTimeoutCount
				);
				this.chartBody.xAxis = [];
				this.desserts = result.data.list.map((e, i) => {
					this.chartBody.xAxis.push(e.proCode);
					todayCount.data.push(e.timeoutCount);
					todayPercent.data.push(e.billPercent * 100);
					yesterdayBillCount.data.push(e.yesterdayTimeoutCount);
					yesterdayBillPercent.data.push(e.yesterdayBillPercent * 100);
					return { ...e, rank: "TOP" + (i + 1) };
				});
				this.chartBody.series = [
					todayCount,
					todayPercent,
					yesterdayBillCount,
					yesterdayBillPercent
				];
				this.setLineAndHistogramChart("treeView", this.chartBody);
			}
		},
		changeType(e) {
			this.rankingType = e;
			this.search();
		}
	},
	mounted() {
		this.search();
	}
};
</script>

<style lang="scss" scoped>
.businessTrends {
	padding: 0 16px;
	display: flex;
	justify-content: space-between;
	.statistics,
	.tableBox {
		display: inline-block;
	}
	.statistics {
		width: 60vw;
		min-width: 700px;
		height: 300px;
		border: 1px solid #999;
		.treeView {
			width: 100%;
			height: 100%;
		}
	}
	.tableBox {
		width: 35vw;
		min-width: 300px;
		margin-left: 20px;
	}
	.more {
		display: inline-block;
		position: absolute;
		color: rgb(91, 107, 115);
		font-size: 0.4em;
		right: 0;
		cursor: pointer;
		font-weight: 400;
	}
	.menu {
		position: absolute;
		top: 0;
		right: 0;
	}
	::v-deep .minTitle {
		font-size: 1em;
		color: red;
	}
}
</style>
