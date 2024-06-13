<template>
	<div class="task-details">
		<div class="z-flex align-center justify-between">
			<lp-tabs class="mb-6" :options="cells.tabsOptions" @change="changeTab"></lp-tabs>

			<div v-if="selectedTab !== 1" class="z-flex">
				<z-text-field
					:formId="formId"
					formKey="saleChannel"
					class="mt-n4 mr-4"
					hideDetails
					label="销售渠道"
				>
				</z-text-field>
				<z-text-field
					:formId="formId"
					formKey="billNum"
					class="mt-n4"
					hideDetails
					label="案件号"
				>
				</z-text-field>

				<z-btn class="pl-3" color="primary" depressed @click="handleSearch">
					<v-icon class="text-h6">mdi-magnify</v-icon>
					查询
				</z-btn>
			</div>
		</div>

		<template v-if="proCode">
			<!-- 任务管理 BEGIN -->
			<task-manage v-if="selectedTab === 1"></task-manage>
			<!-- 任务管理 END -->

			<!-- 案件明细 BEGIN -->
			<case-details v-else ref="details"></case-details>
			<!-- 案件明细 END -->
		</template>
	</div>
</template>

<script>
import { mapState } from "vuex";
import cells from "./cells";

export default {
	name: "TaskDetails",
	provide() {
		return {
			task: {
				proCode: this.proCode,
				project: this.project
			}
		};
	},

	props: {
		proCode: {
			type: String,
			required: true
		},

		project: {
			type: Object,
			requried: true
		}
	},

	computed: {
		...mapState(["forms"])
	},

	data() {
		return {
			formId: "TaskDetailsSearch",
			cells,
			selectedTab: 1
		};
	},

	methods: {
		changeTab({ value }) {
			this.selectedTab = value;
		},

		handleSearch() {

			const params = {
				billNum: this.forms[this.formId].billNum,
				saleChannel: this.forms[this.formId].saleChannel
			};

			this.$refs.details.onFilters(params);
		}
	},

	components: {
		"task-manage": () => import("./taskManage"),
		"case-details": () => import("./caseDetails")
	}
};
</script>
