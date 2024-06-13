<template>
	<div class="basic-info">
		<v-expansion-panels v-model="panel" mandatory multiple>
			<v-expansion-panel>
				<v-expansion-panel-header>
					<div class="z-flex align-center">
						<h6 class="mr-8 fw-bold">出险信息</h6>
					</div>
				</v-expansion-panel-header>

				<v-expansion-panel-content class="pb-8">
					<v-row>
						<v-col
							v-for="(item, index) in fieldList"
							:key="`RiskInfo-${index}`"
							:cols="4"
						>
							<lp-case-fields :fieldList="item"></lp-case-fields>
						</v-col>
					</v-row>
				</v-expansion-panel-content>
			</v-expansion-panel>
		</v-expansion-panels>
	</div>
</template>

<script>
import { R } from "vue-rocket";
import { fieldTypes } from "../cells";
import { LPCaseFields } from "../components";

export default {
	name: "RiskInfo",

	props: {
		riskInfoList: {
			type: Array,
			default: () => []
		}
	},

	data() {
		return {
			panel: [0],
			fieldList: []
		};
	},

	watch: {
		riskInfoList: {
			handler(list) {
				console.log(list);

				if (R.isLousy(list)) return;

				let fieldList = [];

				const length = list.length;

				for (let i = 0; i < length; i++) {
					fieldList[i] = [];
					for (let value of list[i].xmlNodeVal) {
						fieldList[i].push({
							formKey: value,
							inputType: fieldTypes.get(list[i].inputType),
							label: list[i].fieldName,
							hideDetails: true,
							defaultValue: value
						});
					}
				}

				this.fieldList = fieldList;

				console.log(fieldList);
			},
			immediate: true
		}
	},

	components: {
		"lp-case-fields": LPCaseFields
	}
};
</script>
