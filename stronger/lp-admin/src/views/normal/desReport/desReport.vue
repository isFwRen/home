<template>
	<div class="out">
		<div class="table">
			<div class="row1">珠海汇流数据信息销毁报告</div>
			<div class="row2">
				<div class="row2l">销毁文件内容</div>
				<div class="row2r">
					{{ projectMap.get(proCode) }}录入项目{{ year }}年{{
						month
					}}月生产环境原始图像、生产数据库数据
				</div>
			</div>
			<div class="row3">
				<div class="row3l">销毁文件原因</div>
				<div class="row3r">
					按照合同要求汇流保留六十天的数据备份，数据备份期满汇流须销毁数据
				</div>
			</div>
			<div class="row4">
				<div class="row4l">文件删除时间范围</div>
				<div class="row4r">自{{ startTime }} 00：00 至{{ endTime }}24：00</div>
			</div>
			<div class="row5">
				<div class="row5l">文件数量</div>
				<div class="row5r">
					<div class="row5rt">
						<div class="row5rtl">图像文件 （单位：任务包）</div>
						<div class="row5rtr">
							<p class="a1" style="border-bottom: 1px solid black">理赔：{{ wordContent.dNum }}</p>
							<p class="a1">秒赔：{{ wordContent.dNum_2 }}</p>
						</div>
					</div>
					<div class="row5rm">
						<div class="row5rtl">生产数据 （单位：表单）</div>
						<div class="row5rtr">
							<p class="a2" style="border-bottom: 1px solid black">理赔：{{ wordContent.uNum }}</p>
							<p class="a2">秒赔：{{ wordContent.uNum_2 }}</p>
						</div>
					</div>
					<div class="row5rb">
						<div class="row5rtl">图片大小 （单位：MB）</div>
						<div class="row5rtr">
							<p class="a3" style="border-bottom: 1px solid black">理赔：{{ (wordContent.total / 1024).toFixed(2) }}</p>
							<p class="a3">秒赔：{{ (wordContent.total_2 / 1024).toFixed(2) }}</p>
						</div>
					</div>
				</div>
			</div>
			<div class="row6">
				<div class="row6l">删除审核日期</div>
				<div class="row6r">________年____月____日</div>
			</div>
			<div class="row7">
				<div class="row7l">软件执行人</div>
				<div class="row7ll"></div>
				<div class="row7r">网络确认人</div>
				<div class="row7rr"></div>
			</div>
			<div class="row8">
				<div class="row8l">运营经理确认</div>
				<div class="row8r">
					<p>{{ projectMap.get(proCode) }}录入项目{{ year }}年{{ month }}月份生产环境</p>
					<p>原始图像及生产数据库数据已经确认全部进行删除。</p>
					<span>确认人签名:</span>
				</div>
			</div>
		</div>
	</div>
</template>
<script>
import { sessionStorage } from "vue-rocket";
import cells from "./cells";

export default {
	name: "desReport",

	data() {
		return {
			proCode: sessionStorage.get("proCode"),
			year: sessionStorage.get("year"),
			month: sessionStorage.get("month"),
			day: sessionStorage.get("day"),
			startTime: sessionStorage.get("startTime"),
			endTime: sessionStorage.get("endTime"),
			wordContent: sessionStorage.get("wordContent"),
			projectMap: cells.projectMap
		};
	},

	mounted() {
		window.print();
	}
};
</script>
<style scoped lang="scss">
.out {
	display: flex;
	justify-content: center;
	margin-top: 40px;
	.table {
		height: 900px;
		width: 800px;
		border: 1px solid black;

		.row1 {
			width: 100%;
			height: 100px;
			display: flex;
			line-height: 100px;
			justify-content: center;
			font-weight: 800;
			font-size: 25px;
			border-bottom: 1px solid black;
		}

		.row2 {
			width: 100%;
			height: 50px;
			display: flex;
			line-height: 50px;
			border-bottom: 1px solid black;

			.row2l {
				width: 25%;
				text-align: center;
				border-right: 1px solid black;
			}

			.row2r {
				width: 75%;
				text-align: center;
			}
		}

		.row3 {
			width: 100%;
			height: 50px;
			display: flex;
			line-height: 50px;
			border-bottom: 1px solid black;

			.row3l {
				width: 25%;
				text-align: center;
				border-right: 1px solid black;
			}

			.row3r {
				width: 75%;
				text-align: center;
			}
		}

		.row4 {
			width: 100%;
			height: 50px;
			display: flex;
			line-height: 50px;
			border-bottom: 1px solid black;

			.row4l {
				width: 25%;
				text-align: center;
				border-right: 1px solid black;
			}

			.row4r {
				width: 75%;
				text-align: center;
			}
		}

		.row5 {
			width: 100%;
			height: 300px;
			display: flex;
			border-bottom: 1px solid black;

			.row5l {
				width: 25%;
				line-height: 300px;
				text-align: center;
				border-right: 1px solid black;
			}

			.row5r {
				width: 75%;
				.row5rt {
					height: 100px;
					line-height: 100px;
					text-align: center;
					display: flex;
					border-bottom: 1px solid black;
				}

				.row5rm {
					height: 100px;
					line-height: 100px;
					text-align: center;
					display: flex;
					border-bottom: 1px solid black;
				}

				.row5rb {
					height: 98px;
					line-height: 98px;
					text-align: center;
					display: flex;
				}

				.row5rtl {
					width: 45%;
					border-right: 1px solid black;
				}

				.row5rtr {
					width: 55%;

					.a1 {
						width: 100%;
						height: 50px;
						line-height: 50px;
						margin-bottom: 0px;
					}

					.a2 {
						width: 100%;
						height: 50px;
						line-height: 50px;
						margin-bottom: 0px;
					}

					.a3 {
						width: 100%;
						height: 49px;
						line-height: 49px;
						margin-bottom: 0px;
					}
				}
			}
		}

		.row6 {
			width: 100%;
			height: 50px;
			display: flex;
			line-height: 50px;
			border-bottom: 1px solid black;

			.row6l {
				width: 25%;
				text-align: center;
				border-right: 1px solid black;
			}

			.row6r {
				width: 75%;
				text-align: center;
			}
		}

		.row7 {
			width: 100%;
			height: 50px;
			display: flex;
			line-height: 50px;
			border-bottom: 1px solid black;

			.row7l {
				width: 25%;
				text-align: center;
				border-right: 1px solid black;
			}

			.row7ll {
				width: 25%;
				text-align: center;
				border-right: 1px solid black;
			}

			.row7r {
				width: 25%;
				text-align: center;
				border-right: 1px solid black;
			}

			.row7rr {
				width: 25%;
				text-align: center;
			}
		}

		.row8 {
			width: 100%;
			height: 249px;
			display: flex;

			.row8l {
				width: 25%;
				line-height: 249px;
				text-align: center;
				border-right: 1px solid black;
			}

			.row8r {
				width: 75%;
				margin-top: 60px;
				text-align: center;

				span {
					display: block;
					margin-top: 30px;
					margin-left: 150px;
				}
			}
		}
	}
}
</style>
