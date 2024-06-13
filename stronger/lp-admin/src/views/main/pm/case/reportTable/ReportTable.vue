<template>
	<v-card>
		<v-toolbar class="z-toolbar" color="primary" dark>
			<z-btn icon dark @click="onBack">
				<v-icon>mdi-arrow-left</v-icon>
			</z-btn>

			<span class="ml-3">报表</span>

			<v-spacer></v-spacer>

			<z-btn icon dark @click="onClose">
				<v-icon>mdi-close</v-icon>
			</z-btn>
		</v-toolbar>

		<v-card-text class="pt-16">
			<div class="pt-4 pb-8 lp-filters">
				<v-row class="z-flex align-end">
					<v-col
						v-for="(item, index) in cells.fields"
						:key="`entry_filters_${index}`"
						:cols="2"
					>
						<template v-if="item.inputType === 'input'">
							<z-text-field
								:formId="searchFormId"
								:formKey="item.formKey"
								:hideDetails="item.hideDetails"
								:hint="item.hint"
								:label="item.label"
								:suffix="item.suffix"
								:defaultValue="item.defaultValue"
							>
							</z-text-field>
						</template>
					</v-col>

					<div class="z-flex btns">
						<z-btn class="px-3 pb-3" color="primary" @click="onSearch">
							<v-icon class="text-h6">mdi-magnify</v-icon>
							搜索
						</z-btn>

						<z-btn class="pb-3" color="primary" @click="onDownload">
							<v-icon class="text-h6">mdi-download</v-icon>
							下载
						</z-btn>
					</div>
				</v-row>
			</div>

			<v-row class="mb-n1">
				<v-col :cols="2">
					项目金额合计：<span class="error--text">{{ totalPrice }}</span>
				</v-col>

				<v-col :cols="2">
					选定区域项目金额：<span class="error--text">{{ selectPrice }}</span>
				</v-col>

				<v-col :cols="2">
					项目自付金额合计：<span class="error--text">{{ totalPay }}</span>
				</v-col>

				<v-col :cols="2">
					选定区域自付金额：<span class="error--text">{{ selectPay }}</span>
				</v-col>

				<v-col :cols="4">
					<z-text-field
						:formId="searchFormId"
						formKey="indexArea"
						class="mt-n7"
						label="选择序号区域"
						hint="输入选定范围，区间用-隔开，多个用半角逗号隔开，如：5-10,12,15-20"
						:defaultValue="indexArea"
						@input="inputSelectChange"
					>
					</z-text-field>
				</v-col>
			</v-row>

			<vxe-table
				:border="tableBorder"
				:data="desserts"
				ref="reportTable"
				:size="tableSize"
				@checkbox-all="selectAll"
				@checkbox-change="selectChange"
			>
				<vxe-column type="checkbox" width="60" title=""></vxe-column>
				<vxe-column type="seq" width="60" title="序号"></vxe-column>
				<template v-for="item in cells.headers">
					<vxe-column
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					></vxe-column>
				</template>
			</vxe-table>
		</v-card-text>
	</v-card>
</template>

<script>
import { mapState, mapGetters } from "vuex";
import TableMixins from "@/mixins/TableMixins";
import { tools as lpTools } from "@/libs/util";
import CaseMixins from "../CaseMixins";
import cells from "./cells";

const { baseURLApi } = lpTools.baseURL();

export default {
	name: "ReportTable",
	mixins: [TableMixins, CaseMixins],

	data() {
		return {
			formId: "reportTable",
			dispatchList: "GET_REPORTS",
			cells,
			rotate: 0,
			desserts: [],
			indexArea: "",
			totalPrice: 0,
			selectPrice: 0,
			totalPay: 0,
			selectPay: 0
		};
	},

	methods: {
		onBack() {
			this.$emit("back");
		},

		async onDownload() {
			location.href = `${baseURLApi}pro-manager/bill-list/qing-dan/export?id=${this.effectParams.id}&proCode=${this.effectParams.proCode}`;
		},
		async getQingDan() {
			// debugger;
			const result = await this.$store.dispatch("GET_REPORTS", {});
			this.toasted.dynamic(result.msg, result.code);

			if (result.code !== 200) {
			}
		},

		// 全选/全不选
		selectAll(event) {
			this.indexArea = "";
			const { records } = event;
			this.selected = records;
			// this.indexArea = records.length == 0 ? "" : "1-" + records.length;
		},

		// 选中/不选中
		selectChange(event) {
			this.indexArea = "";
			// console.log(event);
			const { records, seq, checked } = event;
			this.selected = records;
			// checked
			//   ? (this.indexArea += seq + ",")
			//   : (this.indexArea = this.indexArea.replace(seq + ",", ""));
		},

		inputSelectChange(val) {
			// console.log(this.desserts[+val - 1]);

			this.clearPrice();
			this.indexArea = val;
			console.log("%chahahahahahahahaha", "color:red;font-size:90px");
			// 5-10,12,15-20
			// console.log(val);
			this.$refs.reportTable.clearCheckboxRow();
			this.selected = [];
			const arrIndex = [];
			if (val == "" || this.desserts.length < +val) {
				return;
			}
			const strVal = val.split(",");
			// console.log(strVal);
			for (let index = 0; index < strVal.length; index++) {
				const element = strVal[index];
				if (element.indexOf("-") != -1) {
					for (let j = element.split("-")[0]; j <= element.split("-")[1]; j++) {
						arrIndex.push(+j);
						// debugger;
						this.selected.push(this.desserts[+j - 1]);
						// console.log(this.selected);
					}
				} else {
					// debugger;
					arrIndex.push(+element);
					this.selected.push(this.desserts[+element - 1]);
					// console.log(this.selected);
				}
			}
			// console.log(arrIndex);
			// console.log(this.selected);
			this.$refs.reportTable.setCheckboxRow(
				this.desserts.filter((val, index) => arrIndex.includes(index + 1)),
				true
			);
		},

		clearPrice() {
			this.selectPrice = 0;
			this.selectPay = 0;
		}
	},

	computed: {
		...mapState(["forms"]),
		...mapGetters(["cases"])
	},

	watch: {
		$route: {
			handler() {
				this.effectParams = {
					id: this.cases.caseInfo.caseId,
					proCode: this.cases.caseInfo.proCode
				};
			},
			immediate: true
		},

		selected(val) {
			this.clearPrice();
			val.forEach(element => {
				// console.log(element);
				this.selectPrice += +element.price;
				this.selectPay += +element.pay;
			});
		},

		indexArea() {
			this.clearPrice();
			// console.log("this.selected");
			// console.log(this.selected);
			this.selected.forEach(element => {
				this.selectPrice += +element.price;
				this.selectPay += +element.pay;
			});
		},

		desserts(val) {
			this.desserts.forEach(element => {
				this.totalPrice += +element.price;
				this.totalPay += +element.pay;
			});
		}
	}
};
</script>
