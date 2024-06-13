<template>
	<div class="basic-info">
		<v-expansion-panels v-model="panel" mandatory multiple>
			<v-expansion-panel v-for="(item, index) in basicInfoList" :key="`BasicInfo${index}`">
				<v-expansion-panel-header>
					<div class="z-flex align-center">
						<h6 class="mr-8 fw-bold">{{ titles[index] }}</h6>
					</div>
				</v-expansion-panel-header>

				<v-expansion-panel-content class="pb-8">
					<v-row>
						<v-col
							v-for="(cItem, cIndex) in item"
							:key="`BasicInfo-Child${cIndex}`"
							:cols="4"
						>
							<lp-case-fields :fieldList="cItem"></lp-case-fields>
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

const titlesMap = new Map([
	["1", "申请人信息"],
	["2", "被保人信息"],
	["3", "受托人信息"],
	["4", "其它信息"]
]);

export default {
	name: "BasicInfo",

	props: {
		basicInfo: {
			type: Object,
			default: () => {}
		}
	},

	data() {
		return {
			panel: [0],
			titles: [],
			basicInfoList: []
		};
	},

	watch: {
		basicInfo: {
			handler(info) {
				let [titles, basicInfoList] = [[], []];

				if (R.isYummy(info)) {
					const keys = Object.keys(info);

					console.log(info);

					for (let i = 0; i < keys.length; i++) {
						titles = [...titles, titlesMap.get(keys[i])];
						basicInfoList[i] = [];

						for (let j = 0; j < info[keys[i]].length; j++) {
							basicInfoList[i][j] = [];

							for (let k = 0; k < info[keys[i]][j].xmlNodeVal.length; k++) {
								basicInfoList[i][j].push({
									formKey: info[keys[i]][j].xmlNodeVal[k],
									inputType: fieldTypes.get(info[keys[i]][j].inputType),
									label: info[keys[i]][j].fieldName,
									hideDetails: true,
									defaultValue: info[keys[i]][j].xmlNodeVal[k]
								});
							}
						}
					}
				}

				this.titles = titles;
				this.basicInfoList = basicInfoList;
			},
			immediate: true
		}
	},

	components: {
		"lp-case-fields": LPCaseFields
	}
};
</script>
