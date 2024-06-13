import * as echarts from "echarts";

export default {
	methods: {
		setPieChart(ref, body) {
			const chartDom = this.$refs[ref];

			if (!chartDom) return;

			const chart = echarts.init(chartDom);

			const option = {
				tooltip: {
					trigger: "item",
					formatter: "{a} <br/>{b}: {c} ({d}%)"
				},

				legend: {
					orient: "vertical",
					left: "left"
				},

				series: [
					{
						label: {
							position: "inner",
							fontSize: 14
						},

						labelLine: {
							show: false
						},

						name: body.name,
						type: "pie",
						radius: "50%",

						data: body.items,

						emphasis: {
							itemStyle: {
								shadowBlur: 10,
								shadowOffsetX: 0,
								shadowColor: "rgba(0, 0, 0, 0.5)"
							}
						}
					}
				]
			};

			option && chart.setOption(option);
		},
		/**
		 *
		 * @param {String} ref 要渲染的节点名
		 * @param {
		 * dataNameList:Array 渲染数据列表名
		 * xAxis:Array X轴坐标名
		 * series:ArrayObj 渲染数据列表名
		 * max:Y轴最大值
		 * min:Y轴最小值
		 * } body
		 * @returns
		 */
		/**
		 * series tempalate
		 {
			name: e.project,
			type: keyType[e.project] ? 'bar' : 'line',
			tooltip: {
				valueFormatter: function(value) {
					return value + outValue
				}
			},
			data: []
		} 
		 */

		setLineAndHistogramChart(ref, body, interval = 50, chartType = '') {
			const chartDom = this.$refs[ref];
			if (!chartDom) return;
			chartDom.innerHtml = "";
			const chart = echarts.init(chartDom);

			let option = {};

			if (chartType === 'pie') {
				option = {
					tooltip: {
					},
					legend: {
						left: 'left',
					},
					series: body.series,
					emphasis: {
						itemStyle: {
							shadowBlur: 10,
							shadowOffsetX: 0,
							shadowColor: "rgba(0, 0, 0, 0.5)"
						}
					},
					avoidLabelOverlap: false,
					label: {
						show: true,
						formatter(param) {
							return param.name + "" + param.value;
						}
					},
					labelLine: {
						length: 2,
						length2: 5
					}
				};
			
			} else if (chartType === 'radar') {

				option = {
					color: ['#5189F8'],
					tooltip: {
						show: true
					},
					legend: {
						data: ['Actual Spending']
					},
					radar: {
						axisName: {
							fontSize: 14,
							color: '#6D7278',
						},
						axisLine: {
							lineStyle: {
								color: '#ebeef3',
							},
						},
						shape: 'circle',
						center: ['50%', '50%'],
						radius: '70%',
						triggerEvent: false,
						indicator: body.indicator || []
					},
					series: body.series || []
				};

			} else if (chartType === 'line') {
				option = {
					tooltip: {
						trigger: 'axis',
						axisPointer: {
							type: 'shadow'
						}
					},
					grid: {
						top: '15%',
						right: '3%',
						left: '5%',
						bottom: '12%'
					},
					xAxis: [{
						type: 'category',
						data: body.xAxis,
						axisLine: {
							lineStyle: {
								color: '#7e9fb6'
							}
						},
						axisLabel: {
							margin: 10,
							color: '#7e9fb6',
							textStyle: {
								fontSize: 14
							},
						},
					}],
					yAxis: [{
						axisLabel: {
							color: '#7e9fb6',
						},
						axisLine: {
							show: false,
							lineStyle: {
								color: '#7e9fb6'
							}
						},
						splitLine: {
							lineStyle: {
								color: 'rgba(126,159,182,0.12)'
							}
						}
					}],
					series: body.series
				};

			} else {
				option = {
					tooltip: {
						trigger: "axis",
						axisPointer: {
							type: "cross",
							crossStyle: {
								shadowColor: "rgba(0, 0, 0, 0.5)"
							}
						}
					},
					legend: {
						data: body.dataNameList || []
					},
					xAxis: [
						{
							type: "category",
							data: body.xAxis
						}
					],
					yAxis: body.yAxis ?? [
						{
							type: "value",
							name: "",
							min: body.min ?? 0,
							max: body.max ?? 300,
							interval,
							axisLabel: {
								formatter: "{value}"
							}
						}
					],
					series: body.series
				};
			}
			option && chart.setOption(option);
		}
	}
};
