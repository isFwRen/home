<template>
	<v-dialog v-model="dialog" width="900">
		<v-card>
			<v-card-title class="text-h6 grey lighten-2"> 赔付设定 </v-card-title>

			<v-card-text>
				<div class="content">
					<div class="content_top">
						<div class="case_number">
							立案号：<span class="blue_text">98746616456464</span>
						</div>
						<div class="case_number">
							责任分类：<span class="blue_text">医疗费用</span>
						</div>
					</div>

					<div class="grid_table">
						<div class="item">
							<div class="theader">责任列表</div>
							<div class="tbody flex_justify">
								<div class="justify_wrap">
									<div class="cell">险种代码：1170</div>
									<div class="cell">责任名称：综合医疗个人部分：0免赔100%</div>
									<div class="cell">赔付：自付一+自付二+自费</div>
								</div>
							</div>
						</div>
						<div class="item">
							<div class="theader">案件条件</div>
							<div class="tbody flex_justify">
								<div class="justify_wrap">
									<div class="cell">任意</div>
									<div class="cell">医疗</div>
									<div class="cell">副被保人</div>
								</div>
							</div>
						</div>
						<div class="item">
							<div class="theader">特殊条件</div>
							<div class="tbody">
								<div class="cell_top">
									<div class="cell">门诊</div>
									<div class="cell">出险时医保报销：无关</div>
									<div class="cell">最小医院等级：无关</div>
								</div>
								<div class="cell_bottom border_top">
									<div class="cell">住院</div>
									<div class="cell">出院时医保报销：无关</div>
									<div class="cell">最小医院等级：无关</div>
									<div class="cell">医院等级范围：无关</div>
								</div>
							</div>
						</div>
						<div class="item">
							<div class="theader">技术计算方式</div>
							<div class="tbody flex_justify">
								<div class="justify_wrap">
									<div class="cell">
										标准费用纳入方式：核准费用，核准费用，自定义免赔款
									</div>
									<div class="cell">控制方式：不控制</div>
									<div class="cell">次控制方式：不控制</div>
								</div>
							</div>
						</div>
					</div>
				</div>
			</v-card-text>

			<v-divider></v-divider>

			<v-card-actions>
				<v-spacer></v-spacer>
				<v-btn color="primary" text @click="dialog = false"> 关闭 </v-btn>
			</v-card-actions>
		</v-card>
	</v-dialog>
</template>

<script>
export default {
	data() {
		return { dialog: false };
	},
	methods: {
		open() {
			this.dialog = true;
		}
	}
};
</script>

<style lang="scss" scoped>
.content {
	padding: 10px;
	font-size: 14px;
	.content_top {
		display: flex;
		gap: 20px;
		padding-bottom: 20px;
		.blue_text {
			color: #007aff;
		}
	}

	.grid_table {
		border: 2px solid #eee;
		border-right: transparent;
		display: flex;
		.item {
			width: 25%;
			border-right: 2px solid #eee;
			.theader {
				background-color: #007aff;
				padding: 10px 5px;
				color: #fff;
			}
			.tbody {
				padding: 6px;
			}
			.cell_top,
			.cell_bottom {
				padding: 10px 0;
			}
			.border_top {
				border-top: 2px solid #eee;
			}
			.flex_justify {
				display: flex;
				height: 220px;
				align-items: center;
			}
		}
	}
}
</style>
