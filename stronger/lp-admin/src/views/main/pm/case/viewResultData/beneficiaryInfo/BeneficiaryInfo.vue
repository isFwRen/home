<template>
	<div class="basic-info">
		<v-expansion-panels v-model="panel" mandatory multiple>
			<v-expansion-panel>
				<v-expansion-panel-header>
					<div class="z-flex align-center">
						<h6 class="mr-8 fw-bold">受益人信息</h6>
					</div>
				</v-expansion-panel-header>

				<v-expansion-panel-content class="pb-8">
					<vxe-table
						ref="billInfoTable"
						:data="desserts"
						resizable
						:border="tableBorder"
						:size="tableSize"
					>
						<vxe-column type="seq" title="序号" width="60"></vxe-column>

						<template v-for="item in headers">
							<vxe-column
								v-if="item.value === 'options'"
								:field="item.value"
								:title="item.text"
								:key="item.value"
							>
								<template #default="{ row }">
									<div class="py-2 z-flex">
										<z-btn
											color="primary"
											depressed
											:lockedTime="0"
											small
											@click="onDetail(row)"
										>
											查看详情
										</z-btn>
									</div>
								</template>
							</vxe-column>

							<vxe-column
								v-else
								:field="item.value"
								:fixed="item.fixed"
								:title="item.text"
								:key="item.value"
								:width="item.width"
							></vxe-column>
						</template>
					</vxe-table>
				</v-expansion-panel-content>
			</v-expansion-panel>

			<v-expansion-panel
				v-for="(item, index) in beneficiaryInfoList"
				:key="`BeneficiaryInfo${index}`"
			>
				<v-expansion-panel-header>
					<div class="z-flex align-center">
						<h6 class="mr-8 fw-bold">{{ titles[index] }}</h6>
					</div>
				</v-expansion-panel-header>

				<v-expansion-panel-content class="pb-8">
					<v-row>
						<v-col
							v-for="(cItem, cIndex) in item"
							:key="`BeneficiaryInfo-Child${cIndex}`"
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
import TableMixins from "@/mixins/TableMixins";
import { fieldTypes } from "../cells";
import { LPCaseFields } from "../components";

const titlesMap = new Map([
	["5", "受益人详细信息"],
	["6", "领款人详细信息"]
]);

export default {
	name: "BeneficiaryInfo",
	mixins: [TableMixins],

	props: {
		beneficiaryInfo: {
			type: Object,
			default: () => {}
		}
	},

	data() {
		return {
			panel: [0],
			titles: [],
			beneficiaryInfoList: []
		};
	},

	watch: {
		beneficiaryInfo: {
			handler(info) {
				let [titles, beneficiaryInfoList] = [[], []];

				if (R.isYummy(info)) {
					const keys = Object.keys(info);

					console.log(info);

					for (let i = 0; i < keys.length; i++) {
						titles = [...titles, titlesMap.get(keys[i])];
						beneficiaryInfoList[i] = [];

						for (let j = 0; j < info[keys[i]].length; j++) {
							beneficiaryInfoList[i][j] = [];

							for (let k = 0; k < info[keys[i]][j].xmlNodeVal.length; k++) {
								beneficiaryInfoList[i][j].push({
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
				this.beneficiaryInfoList = beneficiaryInfoList;
			},
			immediate: true
		},

		"beneficiaryInfo.5": {
			handler(list) {
				let [headers, desserts, formKeys] = [[], [], []];
				let [length, max] = [list.length, 0];

				for (let i = 0; i < length; i++) {
					headers = [...headers, { text: list[i].fieldName, value: list[i].xmlNodeName }];
					formKeys = [...formKeys, list[i].xmlNodeName];
					max =
						R.isYummy(list[i].xmlNodeVal) && list[i].xmlNodeVal.length > max
							? list[i].xmlNodeVal.length
							: max;
				}

				for (let i = 0; i < max; i++) {
					desserts[i] = { formKeys };

					for (let j = 0; j < length; j++) {
						(desserts[i][`${list[j].xmlNodeName}`] = list[j].xmlNodeVal[i]),
							(desserts[i][`${list[j].xmlNodeName}InputType`] = fieldTypes.get(
								list[j].inputType
							)),
							(desserts[i][`${list[j].xmlNodeName}FormKey`] = list[j].xmlNodeName),
							(desserts[i][`${list[j].xmlNodeName}FieldName`] = list[j].fieldName);
					}
				}

				this.headers = headers;

				this.desserts = desserts;
			},
			immediate: true
		}
	},

	components: {
		"lp-case-fields": LPCaseFields
	}
};
</script>
