<template>
	<div class="countInfo">
		<div class="count_header">理算信息</div>
		<div class="grid">
			<div class="cell" v-for="(value, key) in formsCess" :key="key">
				<input type="text" v-model="formsCess[key]" />
			</div>
			<div class="cell blue_color" @click="openDialog">详见赔付设定</div>
		</div>

		<div class="btn_group">
			<v-btn depressed style="border: 1px solid #709eff; color: #709eff"> 本次影像 </v-btn>
			<v-btn depressed style="border: 1px solid #709eff; color: #709eff"> 历时赔付 </v-btn>
			<v-btn depressed style="border: 1px solid #709eff; color: #709eff"> 责任配置 </v-btn>
			<v-btn depressed style="border: 1px solid #709eff; color: #709eff"> 风险因子 </v-btn>
			<v-btn depressed style="border: 1px solid #709eff; color: #709eff"> 反洗钱 </v-btn>
			<v-btn depressed style="border: 1px solid #709eff; color: #709eff"> 不予处理 </v-btn>
			<v-btn depressed color="primary" @click="onSubmit"> 提交 </v-btn>
		</div>

		<div class="tabs">
			<v-tabs v-model="tab">
				<v-tab>费用责任汇总</v-tab>
				<v-tab>津贴责任汇总</v-tab>
				<v-tab>定额责任汇总</v-tab>
			</v-tabs>
			<div class="tab_content">
				<Cost v-if="tab === 0" :options="costForms" @updateForm="onUpdateForm" />
				<Quota v-if="tab === 1" :options="quotaForms" @updateForm="onUpdateForm" />
				<Allowance v-if="tab === 2" :options="allowanceForms" @updateForm="onUpdateForm" />
			</div>
		</div>
		<Compensation ref="compensation" />
	</div>
</template>

<script>
import { sessionStorage } from "vue-rocket";
export default {
	props: {
		CountInfo: {
			type: Object,
			default: () => {}
		}
	},
	data() {
		return {
			content: {},
			tab: 0,
			flag: [0, 1, 2, 3, 4],
			formsCess: {
				name1: "立案号",
				name2: "险种",
				name3: "给付类型",
				name4: "理算金额",
				name5: "是否拒赔",
				name6: "拒赔金额",
				name7: "实际赔付金额",
				name8: "处理意见",
				name9: "责任详情",
				name10: "9876543234567899",
				name11: "1190",
				name12: "医疗费用",
				name13: "5432.11",
				name14: "否",
				name15: "0.00",
				name16: "1234.11",
				name17: "正常赔付"
			},
			allowanceForms: {
				name1: "1234.11",
				name2: "1234.11",
				name3: "1234.11",
				name4: "1234.11",
				name5: "1234.11",
				name6: "1234.11",
				name7: "1234.11",
				name8: "1234.11",
				name9: "1234.11",
				name10: "1234.11",
				name11: "1234.11",
				reason: "5437474"
			},
			costForms: {
				name1: "1234.11",
				name2: "1234.11",
				name3: "1234.11",
				name4: "1234.11",
				name5: "1234.11",
				name6: "1234.11",
				name7: "1234.11",
				name8: "1234.11",
				name9: "1234.11",
				name10: "1234.11",
				name11: "1234.11",
				reason: "5437474"
			},
			quotaForms: {
				name1: "1234.11",
				name2: "1234.11",
				name3: "1234.11",
				name4: "1234.11",
				name5: "1234.11",
				name6: "1234.11",
				name7: "1234.11",
				name8: "1234.11",
				name9: "1234.11",
				name10: "1234.11",
				name11: "1234.11",
				reason: "5437474"
			}
		};
	},
	created() {
		this.initData();

		this.content = sessionStorage.get("checkForm");
		this.content.countInfo.formsCess = this.formsCess;
		this.content.countInfo.allowanceForms = this.allowanceForms;
		this.content.countInfo.costForms = this.costForms;
		this.content.countInfo.quotaForms = this.quotaForms;
	},
	methods: {
		onUpdateForm({ forms, key }) {
			this.content.countInfo[key] = forms;
			this.updateSession();
		},
		updateSession() {
			sessionStorage.set("checkForm", this.content);
		},
		initData() {
			console.log(this.CountInfo, "CountInfo");

			if (this.CountInfo.formsData1 && Object.keys(this.CountInfo.formsCess).length > 0) {
				this.formsCess = this.CountInfo.formsCess;
			}

			if (
				this.CountInfo.allowanceForms &&
				Object.keys(this.CountInfo.allowanceForms).length > 0
			) {
				this.allowanceForms = this.CountInfo.allowanceForms;
			}

			if (this.CountInfo.costForms && Object.keys(this.CountInfo.costForms).length > 0) {
				this.costForms = this.CountInfo.costForms;
			}

			if (this.CountInfo.quotaForms && Object.keys(this.CountInfo.quotaForms).length > 0) {
				this.quotaForms = this.CountInfo.quotaForms;
			}
		},
		openDialog() {
			this.$refs.compensation.open();
		},
		async onSubmit() {
			this.$EventBus.$emit("submitData");
		}
	},
	components: {
		Cost: () => import("./cost/index.vue"),
		Quota: () => import("./quota/index.vue"),
		Allowance: () => import("./allowance/index.vue"),
		Compensation: () => import("./compensation/index.vue")
	}
};
</script>

<style lang="scss">
.blue_color {
	color: #007aff;
	font-weight: 600;
	cursor: pointer;
}
.countInfo {
	padding: 10px;
	background-color: #fff;
	.count_header {
		background-color: #e9f6ff;
		padding: 10px 10px;
		color: #007aff;
		font-weight: 600;
		font-size: 15px;
	}

	.grid {
		display: grid;
		margin-top: 10px;
		grid-template-columns: repeat(9, 1fr);
		font-size: 14px;
		border-bottom: 1px solid #eee;
		border-right: 1px solid #eee;
		.cell {
			border-top: 1px solid #eee;
			border-left: 1px solid #eee;
			padding: 5px 5px;
			input {
				display: block;
				width: 100%;
				height: 100%;
				padding: 5px;
				border: 1px solid transparent;
				border-radius: 10px;
			}
			input:focus {
				outline: #eee;
				border: 1px solid #eee;
			}
		}
	}

	.btn_group {
		display: flex;
		justify-content: end;
		padding: 10px 0;
		gap: 10px;
	}

	.tab_content {
		padding: 10px 20px;
	}
}
</style>
