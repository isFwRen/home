<template>
	<div>
		<h1 class="he">Hello， {{ name }}</h1>
		<p class="p">
			为了让您达到正式上岗的必备知识和技能，请按照以下培训流程完成学习、练习，完成练习考核目标后即可正式开始录入兼职啦~
		</p>
		<p class="p">当前进度：<span>培训指引流程学习</span></p>
		<p class="p1">培训流程：</p>
		<div class="guide_wrap">
			<div class="steps_wrap">
				<div class="step">
					<h5 :class="['step_l', steps[1] ? 'c_blue' : '']">培训指引流程学习</h5>
					<div class="step_c"></div>
					<div class="step_r">
						<span class="c_green" v-if="!steps[1]" @click="goGuidance">前往学习</span>
						<span class="status" v-else>已完成 <v-icon class="c_blue">mdi-check-circle</v-icon></span>
					</div>
				</div>
				<div class="line_wrap">
					<div class="line"></div>
				</div>

				<div class="step">
					<h5 :class="['step_l', steps[2] ? 'c_blue' : '']">教学文件学习</h5>
					<div class="step_c" v-if="current >= 2"></div>
					<div class="step_r" v-if="current >= 2">
						<span class="c_green" v-if="!steps[2]" @click="goBusiness">前往学习</span>
						<span class="status" v-else>已完成 <v-icon style="color: #3894ff">mdi-check-circle</v-icon></span>
					</div>
				</div>
				<div class="line_wrap">
					<div class="line"></div>
				</div>

				<div class="step">
					<h5 :class="['step_l', steps[3] ? 'c_blue' : '']">录入练习</h5>
					<div class="step_c" v-if="current >= 3"></div>
					<div class="step_r" v-if="current >= 3">
						<span class="c_green" v-if="!steps[3]" @click="goEntry">前往练习</span>
						<span class="status" v-else>已完成 <v-icon style="color: #3894ff">mdi-check-circle</v-icon></span>
					</div>
				</div>
				<div class="line_wrap">
					<div class="line"></div>
				</div>

				<div class="step">
					<h5 :class="['step_l', steps[4] ? 'c_blue' : '']">上岗考核</h5>
					<div class="step_c" v-if="current >= 4"></div>
					<div class="step_r" v-if="current >= 4">
						<span class="c_green" v-if="!steps[4]">前往考核</span>
						<span class="status" v-else>完成考核 <v-icon style="color: #3894ff">mdi-check-circle</v-icon></span>
					</div>
				</div>
				<div class="line_wrap">
					<div class="line"></div>
				</div>

				<div class="step">
					<h5 :class="['step_l', steps[5] ? 'c_blue' : '']">上岗审核</h5>
					<div class="step_c" v-if="current >= 5"></div>
					<div class="step_r" v-if="current >= 5">
						<span class="c_green" v-if="!steps[5]">前往审核</span>
						<span class="c_yellow" v-else>审核通过 <v-icon style="color: #3894ff">mdi-check-circle</v-icon></span>
					</div>
				</div>
				<div class="line_wrap">
					<div class="line"></div>
				</div>

				<div class="step">
					<h5 :class="['step_l', steps[6] ? 'c_blue' : '']">正式上岗</h5>
					<div class="step_c" v-if="current >= 6"></div>
					<div class="step_r" v-if="current >= 6">
						<span class="c_green" v-if="!steps[6]">前往上岗</span>
						<span class="c_red" v-else>正式上岗</span>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import { mapState } from "vuex";
import { localStorage } from "vue-rocket";

export default {
	name: "Train",
	mixins: [],

	data() {
		return {
			name: localStorage.get("user").name,
			current: 1,
			dispatchList: "TRAINING_GUIDE",
			steps: {
				1: false,
				2: false,
				3: false,
				4: false,
				5: false,
				6: false
			},
			broadListen: null
		};
	},
	async created() {
		this.getStatus();
		this.broadListen = new BroadcastChannel("updateRead");
		this.broadListen.addEventListener("message", event => {
			this.getStatus();
		});
	},
	computed: {
		...mapState(["train"])
	},
	methods: {
		async getStatus() {
			const result = await this.$store.dispatch(this.dispatchList);
			if (result.code === 200) {
				this.current = result.data.userTrainingStage;
				this.changeStatus(result.data.userTrainingStage);
			}
		},
		changeStatus(number) {
			for (let i = 1; i < number; i++) {
				this.steps[i] = true;
			}
		},
		goGuidance() {
			window.open(`${location.origin}/normal/doc-file`, "_blank", "toolbar=yes, scrollbars=yes, resizable=yes");
		},
		goBusiness() {
			this.$router.push("/main/rule/business");
		},
		goEntry() {
			this.$router.push("/main/practice/practiceChannel");
		},
		next() {
			if (this.current < 6) {
				this.current++;
			}
		}
	},
	destroy() {
		this.broadListen.close();
	},
	components: {}
};
</script>

<style scoped lang="scss">
.he {
	font-size: 22px;
	margin-bottom: 20px;
}

.p {
	font-size: 18px;
	margin-bottom: 20px;

	span {
		cursor: pointer;
		color: orangered;
	}
}

.p1 {
	font-size: 18px;
}

.guide_wrap {
	width: 100%;
	display: flex;
	justify-content: center;

	.steps_wrap {
		width: 42%;

		.c_green {
			color: #58a55c;
		}

		.c_yellow {
			color: #ecb51c;
		}

		.c_blue {
			color: #3894ff;
		}

		.c_red {
			color: #fb1010;
		}

		.step {
			min-width: 590px;
			display: flex;
			gap: 20px;
			align-items: center;

			.step_l {
				width: 24%;
				min-width: 150px;
				font-size: 18px;
			}

			.step_c {
				width: 36%;
				min-width: 250px;
				border: 1px dashed #eee;
			}

			.step_r {
				width: 20%;
				display: flex;
				font-weight: 600;
				font-size: 18px;

				span {
					cursor: pointer;
					user-select: none;
				}
			}
		}

		.line_wrap {
			padding: 10px 30px;

			.line {
				background-color: #e0e3ea;
				width: 6px;
				height: 50px;
				border-radius: 10px;
			}
		}
	}
}

.steps {
	display: flex;
	justify-content: flex-start;

	.step {
		width: 800px;
		margin-left: 250px;

		.c {
			display: flex;
			justify-content: center;
			flex-direction: column;
			height: 100px;
			border-radius: 10px;
			margin-bottom: 15px;
			background-color: #e3eef9;

			.btn_wrap {
				height: 25px;
				line-height: 25px;
				font-size: 16px;
				font-weight: bold;
				text-align: center;
				cursor: pointer;
			}

			.btn_wrap:hover {
				text-decoration: underline;
				color: #00a9ff !important;
			}
		}

		.r {
			margin-right: 10px;
		}
	}
}
</style>
