<template>
	<v-card>
		<v-toolbar class="z-toolbar" color="primary" dark>
			<z-btn icon dark @click="onBack">
				<v-icon>mdi-arrow-left</v-icon>
			</z-btn>

			<lp-tooltip-btn
				bottom
				btnIcon
				fab
				icon="mdi-swap-horizontal"
				small
				tip="切换到修改录入数据"
				@click="onSwitch"
			>
			</lp-tooltip-btn>

			<span class="pl-2">案件号 {{ cases.caseInfo.billNum }}</span>

			<v-spacer></v-spacer>

			<z-btn icon dark @click="onClose">
				<v-icon>mdi-close</v-icon>
			</z-btn>
		</v-toolbar>

		<v-card-text class="pt-16">
			<lp-tabs class="mb-6" :options="tabsOptions" @change="onTab"></lp-tabs>

			<basic-info
				v-if="currentTab === 'basicInfo' && R.isYummy(basicInfo)"
				:basicInfo="basicInfo"
			></basic-info>

			<beneficiary-info
				v-else-if="currentTab === 'beneficiaryInfo' && R.isYummy(beneficiaryInfo)"
				:beneficiaryInfo="beneficiaryInfo"
			></beneficiary-info>

			<bill-info
				v-else-if="currentTab === 'billInfo' && R.isYummy(billInfoList)"
				:billInfoDetail="billInfoDetail"
				:billInfoList="billInfoList"
			></bill-info>

			<risk-info
				v-else-if="currentTab === 'riskInfo' && R.isYummy(riskInfoList)"
				:riskInfoList="riskInfoList"
			></risk-info>
		</v-card-text>
	</v-card>
</template>

<script>
import { mapGetters } from "vuex";
import { R } from "vue-rocket";
import DialogMixins from "@/mixins/DialogMixins";
import CaseMixins from "../CaseMixins";
import cells from "./cells";

const basicInfoKeys = ["1", "2", "3", "4"];
const beneficiaryInfoKeys = ["5", "6"];
const billInfoKeys = ["7"];
const riskInfoKeys = ["8"];

export default {
	name: "ViewResultData",
	mixins: [DialogMixins, CaseMixins],

	data() {
		return {
			R,
			formId: "viewResultData",
			cells,
			currentTab: "",

			basicInfo: {},
			beneficiaryInfo: {},
			billInfoList: [],
			riskInfoList: [],
			tabsOptions: []
		};
	},

	created() {
		this.getResultData();
	},

	methods: {
		onSwitch() {
			this.$router.push({ path: "/main/PM/case/update-entry-data" });
		},

		onTab(item) {
			this.currentTab = item.value;
		},

		async getResultData() {
			const { proCode, caseId } = this.cases.caseInfo;

			const form = {
				proCode,
				id: caseId
			};

			let [basicInfo, beneficiaryInfo, billInfoList, riskInfoList] = [{}, {}, [], []];
			let tabsOptions = [];

			const result = await this.$store.dispatch("GET_CASE_RESULT_DATA", form);

			if (result.code !== 200) {
				this.toasted.dynamic(result.msg, result.code);
				return;
			}

			this.billInfoDetail = result.data.amount;

			for (let key in result.data.data) {
				// 基础信息
				if (basicInfoKeys.includes(key)) basicInfo[key] = result.data.data[key];
				// 受益人信息
				if (beneficiaryInfoKeys.includes(key)) beneficiaryInfo[key] = result.data.data[key];
				// 账单信息
				if (billInfoKeys.includes(key)) billInfoList = result.data.data[key];
				// 出险信息
				if (riskInfoKeys.includes(key)) riskInfoList = result.data.data[key];
			}

			// tab 基础信息
			if (R.isYummy(basicInfo)) {
				tabsOptions.push(cells.tab1);
			}
			// tab 受益人信息
			if (R.isYummy(beneficiaryInfo)) {
				tabsOptions.push(cells.tab2);
			}
			// tab 账单信息
			if (R.isYummy(billInfoList)) {
				tabsOptions.push(cells.tab3);
			}
			// tab 出险信息
			if (R.isYummy(riskInfoList)) {
				tabsOptions.push(cells.tab4);
			}

			this.basicInfo = basicInfo;
			this.beneficiaryInfo = beneficiaryInfo;
			this.billInfoList = billInfoList.sort((el1, el2) => el1.billInfo - el2.billInfo);
			this.riskInfoList = riskInfoList;
			this.tabsOptions = tabsOptions;
			this.currentTab = this.tabsOptions[0].value;
		}
	},

	computed: {
		...mapGetters(["cases"])
	},

	components: {
		"basic-info": () => import("./basicInfo"),
		"beneficiary-info": () => import("./beneficiaryInfo"),
		"bill-info": () => import("./billInfo"),
		"risk-info": () => import("./riskInfo")
	}
};
</script>
