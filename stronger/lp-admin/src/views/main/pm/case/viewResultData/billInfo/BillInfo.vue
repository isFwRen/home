<template>
	<div class="basic-info">
		<v-expansion-panels v-model="panel" mandatory multiple>
			<v-expansion-panel>
				<v-expansion-panel-header>
					<div class="z-flex align-center">
						<h6 class="mr-8 fw-bold">账单详细信息</h6>

						<ul class="z-flex px-2">
							<li class="mr-4">
								<label>账单金额总和</label>
								<span class="ml-4 error--text">{{
									billInfoDetail.billAmount
								}}</span>
							</li>

							<li class="mr-4">
								<label>调整金额总和</label>
								<span class="ml-4 error--text">{{
									billInfoDetail.adjustmentAmount
								}}</span>
							</li>

							<li class="mr-4">
								<label>扣费金额总和</label>
								<span class="ml-4 error--text">{{ billInfoDetail.deductionAmount }}</span>
							</li>

							<li>
								<label>报销金额总和</label>
								<span class="ml-4 error--text">{{
									billInfoDetail.feeAmount
								}}</span>
							</li>
						</ul>
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

					<!-- <v-row>
            <v-col 
              v-for="(item, index) in fieldList"
              :key="index"
              :cols="4"
            >  -->
					<lp-case-fields
						v-if="isExpand"
						:cols="4"
						:fieldList="fieldList"
					></lp-case-fields>
					<!-- </v-col>
          </v-row> -->
				</v-expansion-panel-content>
			</v-expansion-panel>
		</v-expansion-panels>
	</div>
</template>

<script>
import { R } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import { fieldTypes } from "../cells";
import cells from "./cells";
import { LPCaseFields } from "../components";

export default {
	name: "BillInfo",
	mixins: [TableMixins],

	props: {
		billInfoDetail: {
			type: Object,
			default: () => {}
		},

		billInfoList: {
			type: Array,
			default: () => []
		}
	},

	data() {
		return {
			cells,
			panel: [0],
			headers: [],
			fieldList: [],
			rowId: null,
			isExpand: false
		};
	},

	methods: {
		onDetail(row) {
			const fieldList = [];

			for (let key of row.formKeys) {
				fieldList.push({
					formKey: key,
					inputType: row[`${key}InputType`],
					label: row[`${key}FieldName`],
					hideDetails: true,
					defaultValue: row[key]
				});
			}

			this.fieldList = fieldList;

			if (this.rowId !== row._XID) {
				this.isExpand = true;
				this.rowId = row._XID;
			} else {
				this.isExpand = !this.isExpand;
			}
		}
	},

	watch: {
		billInfoList: {
			async handler(list) {
				const result = await this.$store.dispatch("GET_CONFIG_QUALITY_CONST_LIST");
				let mapHeader = result.data.qualityBillInfo;
				// console.log("list", list);
				// console.log("mapHeader", mapHeader);
				list.forEach(el => (el.fieldName = mapHeader[el.billInfo]));
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

				this.headers = [...headers, { text: "操作", value: "options" }];
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
