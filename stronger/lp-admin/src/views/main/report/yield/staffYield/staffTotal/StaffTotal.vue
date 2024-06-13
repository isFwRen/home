<template>
	<div class="staff-total">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col cols="2">
					<z-select
						:formId="searchFormId"
						formKey="type"
						hideDetails
						label="类型"
						:options="typeOptions"
						:defaultValue="defaultType"
						@change="changeType"
					>
					</z-select>
				</v-col>

				<v-col cols="3">
					<z-date-picker
						:formId="searchFormId"
						formKey="date"
						hideDetails
						label="日期"
						range
						z-index="10"
						:defaultValue="DEFAULT_DATE"
					></z-date-picker>
				</v-col>

				<v-col cols="2">
					<z-text-field
						:formId="searchFormId"
						formKey="code"
						hideDetails
						label="工号/姓名"
					>
					</z-text-field>
				</v-col>

				<div class="z-flex">
					<z-btn class="pb-3 pl-3" color="primary" @click="onSearch()">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>
				</div>
			</v-row>
		</div>

		<div class="pb-6 btns">
			<z-btn class="pr-3" color="primary" small outlined @click="onCopy"> 复制 </z-btn>
			<z-btn class="pr-3" color="primary" small outlined @click="onExport"> 导出 </z-btn>
			<z-btn class="pr-3" color="primary" small outlined @click="onUpdate"> 更新 </z-btn>
			<z-btn class="pr-3" color="primary" small outlined @click="onSet"> 设置 </z-btn>
		</div>

		<div class="table staff-total-table" id="tableList">
			<vxe-table
				:border="tableBorder"
				:data="desserts"
				:max-height="tableMaxHeight"
				:size="tableSize"
				:sort-config="{ multiple: true, trigger: 'cell' }"
				align="center"
			>
				<template v-for="(item, i) in headers">
					<vxe-colgroup
						v-if="i === 0"
						align="center"
						:title="item.text"
						:key="item.value"
					>
						<template v-for="record in item.children">
							<vxe-column
								:sortable="record.sortable"
								:width="record.width"
								:field="record.value"
								:title="record.text"
								:key="record.value"
							></vxe-column>
						</template>
					</vxe-colgroup>

					<vxe-colgroup align="center" :title="item.text" :key="item.value" v-else>
						<template v-for="record in item.children">
							<vxe-column
								:title="record.text"
								:width="record.width"
								:key="record.value"
							>
								<template #default="{ row }">
									<!-- {{ row.proSummary[i - 1]['proCode'] ==  item.value ? row.proSummary[i - 1][record.value] : '000'}} -->
									{{
										row.proSummary.find(val => val.proCode == item.text)[
											record.value
										]
									}}
								</template>
							</vxe-column>
						</template>
					</vxe-colgroup>
				</template>
			</vxe-table>

			<z-pagination
				:total="pagination.total"
				:options="pageSizes"
				@page="handlePage"
			></z-pagination>
		</div>

		<staff-update-dialog
			ref="update"
			:rowInfo="detailInfo"
			@updateTime="handleUpdateTime"
		></staff-update-dialog>

		<staff-set-dialog ref="set"></staff-set-dialog>

		<router-view></router-view>
	</div>
</template>

<script>
import moment from "moment";
import TableMixins from "@/mixins/TableMixins";
import SocketsMixins from "@/mixins/SocketsMixins";
import { typeOptions } from "../cells";
import cells from "./cells";
import { copy, copyText } from "clipboard-vue";
import { tools } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
const DEFAULT_DATE = [moment().format("YYYY-MM-DD"), moment().format("YYYY-MM-DD")];

export default {
	name: "StaffTotal",
	mixins: [TableMixins, SocketsMixins],

	data() {
		return {
			DEFAULT_DATE,
			formId: "StaffTotal",
			dispatchList: "GET_STAFF_TOTAL",
			cells,
			typeOptions,
			defaultType: 1,
			headers: [],
			socketPath: "outputStatistics"
		};
	},
	methods: {
		// 复制到粘貼板
		onCopy() {
			const removeKey = ["submitTime", "proCode"];
			const dataArr = [];

			const rowTopOne = [];
			const rowTopTwo = [];

			for (var i = 0; i < this.headers.length; i++) {
				rowTopOne.push('"' + this.headers[i].text + '"');

				for (var j = 0; j < this.headers[i].children.length; j++) {
					rowTopOne.push('"' + " " + '"');
					rowTopTwo.push('"' + this.headers[i].children[j].text + '"');
				}
			}

			dataArr.push(rowTopOne.join("\t"));
			dataArr.push(rowTopTwo.join("\t"));

			this.desserts.forEach(element => {
				const rowArr = [];

				for (const key of Object.keys(element)) {
					if (key == "proSummary" && element[key].constructor == Array) {
						element[key].forEach(item => {
							for (const k of Object.keys(item)) {
								console.log("k", k);
								if (k == "op0InvoiceNum") {
									continue;
								}
								if (removeKey.indexOf(k) == -1) {
									console.log("item[k]", item[k]);
									rowArr.push('"' + item[k] + '"');
								}
							}
						});
					} else {
						if (removeKey.indexOf(key) == -1) {
							if (/^row/.test(element[key])) {
								console.log("row");
							} else {
								rowArr.push('"' + element[key] + '"');
							}
						}
					}
				}
				dataArr.push(rowArr.join("\t"));
			});

			if (dataArr.length == 0) {
				this.toasted.warning("没有查询到数据");
			}

			console.log(dataArr);

			copyText(dataArr.join("\n"))
				.then(e => {
					this.toasted.success("复制成功");
				})
				.catch(e => {
					this.toasted.error("复制失败");
				});
		},

		buildHeaders(proHeader) {
			this.headers = tools.deepClone(this.cells.headers);
			for (const iterator of proHeader) {
				this.headers.push({
					text: iterator,
					value: iterator,
					children: this.cells.children
				});
			}
		},

		changeType(value) {
			this.$emit("type", value);
		},

		onUpdate() {
			// console.log("form", this.forms[this.searchFormId])
			// var form = this.forms[this.searchFormId]
			// this.getDetail(form);

			this.$refs.update.onOpen();
		},

		handleUpdateTime({ updateTime }) {
			this.effectParams.updateTime = updateTime;

			// console.log(this.effectParams);
			// console.log(updateTime);

			this.onSearch();
		},

		onSet() {
			this.$refs.set.onOpen();
		},

		async onExport() {
			const form = this.forms[this.searchFormId];
			if (!tools.isYummy(form.date)) {
				this.toasted.warning("沒有选择日期");
				return;
			}

			const fileName = `人员产量统计${form.date[0]}-${form.date[1]}.xlsx`;
			lpTools.exportExcel(fileName, "#tableList");
			this.toasted.info("正在导出...");
		}
	},

	watch: {
		"sabayon.data.top": {
			handler(top) {
				if (top) {
					this.buildHeaders(top);
				}
			},
			immediate: true,
			deep: true
		},
		desserts(val) {
			val.forEach(item => {
				item["proSummary"].forEach(record => {
					record.op1 =
						Number(record.op1NotExpenseAccount) + Number(record.op1ExpenseAccount);
					record.op2 =
						Number(record.op2NotExpenseAccount) + Number(record.op2ExpenseAccount);
				});
			});
		}
	},

	directives: {
		copy
	},

	components: {
		"staff-update-dialog": () => import("./staffUpdateDialog"),
		"staff-set-dialog": () => import("./staffSetDialog")
		// "staff-export-dialog": () => import("./staffExportDialog"),
	}
};
</script>
