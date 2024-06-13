<template>
	<div class="staff-detail">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col cols="2">
					<z-select :formId="searchFormId" formKey="type" hideDetails label="类型" :options="typeOptions"
						:defaultValue="defaultType" @change="changeType">
					</z-select>
				</v-col>
				<v-col cols="2">
					<z-select :formId="searchFormId" formKey="proCode" hideDetails label="项目" :options="auth.proItems">
					</z-select>
				</v-col>

				<v-col cols="3">
					<z-date-picker :formId="searchFormId" formKey="date" hideDetails label="日期" range z-index="10"
						:defaultValue="DEFAULT_DATE"></z-date-picker>
				</v-col>

				<v-col cols="2">
					<z-text-field :formId="searchFormId" formKey="code" hideDetails label="工号/姓名">
					</z-text-field>
				</v-col>

				<div class="z-flex">
					<z-btn class="pb-3 pl-3" color="primary" @click="Search">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>
					<!-- <z-btn class="pb-3 pl-3" color="primary" @click="ClearMySelf">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						清空
					</z-btn> -->
				</div>
			</v-row>
		</div>

		<div class="pb-6 btns">
			<z-btn class="pr-3" color="primary" small outlined @click="onCopy"> 复制 </z-btn>
			<z-btn class="pr-3" color="primary" small outlined @click="onExport"> 导出 </z-btn>
		</div>

		<div class="table staff-detail-table">
			<vxe-table :border="tableBorder" :data="desserts" :size="tableSize" :max-height="tableMaxHeight"
				:sort-config="{ multiple: true, trigger: 'cell' }" align="center">
				<template v-for="item in cells.headers">
					<vxe-colgroup v-if="item.value === 'Summary' ||
						item.value === 'first' ||
						item.value === 'one' ||
						item.value === 'two' ||
						item.value === 'problem'
						" align="center" :title="item.text" :key="item.value">
						<template v-for="record in item.children">
							<vxe-column :field="record.value" :title="record.text" :key="record.value" width="70px"></vxe-column>
						</template>
					</vxe-colgroup>
					<vxe-column v-else-if="item.value === 'CreatedAt'" :field="item.value" :title="item.text" :key="item.value"
						width="150px">
						<template #default="{ row }">
							{{ row.CreatedAt | dateFormat("YYYY-MM-DD HH:mm:ss") }}
						</template>
					</vxe-column>

					<vxe-column v-else :sortable="item.sortable" :field="item.value" :title="item.text" :key="item.value"
						width="80px"></vxe-column>
				</template>
			</vxe-table>
		</div>

		<z-pagination :total="pagination.total" :options="pageSizes" @page="handlePage"></z-pagination>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import moment from "moment";
import TableMixins from "@/mixins/TableMixins";
import SocketsMixins from "@/mixins/SocketsMixins";
import { typeOptions } from "../cells";
import cells from "./cells";
import { copy, copyText } from "clipboard-vue";
import { R } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
import _ from "lodash";

const DEFAULT_DATE = [moment().format("YYYY-MM-DD"), moment().format("YYYY-MM-DD")];

export default {
	name: "StafflDetail",
	mixins: [TableMixins, SocketsMixins],

	data() {
		return {
			DEFAULT_DATE,
			formId: "StafflDetail",
			dispatchList: "GET_STAFF_TOTAL",
			cells,
			typeOptions,
			defaultType: 2,
			manual: true,
			socketPath: "outputStatisticsDetail"
		};
	},
	watch: {
		desserts(val) {
			val.forEach(item => {
				item.op1FieldCharacter =
					Number(item.op1NotExpenseAccountFieldCharacter) +
					Number(item.op1ExpenseAccountFieldCharacter);

				item.op1FieldEffectiveCharacter =
					Number(item.op1NotExpenseAccountFieldEffectiveCharacter) +
					Number(item.op1ExpenseAccountFieldEffectiveCharacter);

				item.op2FieldCharacter =
					Number(item.op2NotExpenseAccountFieldCharacter) +
					Number(item.op2ExpenseAccountFieldCharacter);

				item.op2FieldEffectiveCharacter =
					Number(item.op2NotExpenseAccountFieldEffectiveCharacter) +
					Number(item.op2ExpenseAccountFieldEffectiveCharacter);
			});
		}
	},
	methods: {
		async ClearMySelf() {
			const user = this.storage.get("user");
			var copyDesserts = [];
			for (let i = 0; i < this.desserts.length; i++) {
				if (user.nickName != this.desserts[i].nickName) {
					copyDesserts.push(this.desserts[i]);
				}
			}

			this.desserts = null;
			this.desserts = copyDesserts;
			const param = this.forms[this.searchFormId];

			const s = {
				proCode: param.proCode,
				startTime:
					(param.date ? param.date[0] : moment().format("yyyy-MM-DD")) + " 00:00:00",
				endTime: (param.date ? param.date[1] : moment().format("yyyy-MM-DD")) + " 23:59:59",
				code: user.code
			};

			const result = await this.$store.dispatch("Delete_STAFF_TOTAL", s);
			if (result.code == 200) {
				this.toasted.success("操作成功");
			}
		},
		async Search() {
			if (this.dispatchList) {
				const forms = this.forms[this.searchFormId];
				const data = _.cloneDeep(forms);

				const reg = /^[0-9]/;
				if (forms.code) {
					const bool = reg.test(forms.code);
					if (!bool) {
						Reflect.deleteProperty(data, "code");
						data.code = forms.code;
					}
				}
				const params = {
					...this.effectParams,
					...this.params,
					...data
				};
				const result = await this.$store.dispatch(this.dispatchList, params);
				const { list, total } = result.data;
				if (result.code === 200) {
					if (typeof list === "object") {
						if (list instanceof Array) {
							this.desserts = list;
						} else {
							this.desserts = [];
						}
						this.pagination.total = total;
					} else {
						this.desserts = result.data;
						this.pagination.total = this.desserts.length;
					}
				} else {
					this.toasted.error(result.msg);

					this.desserts = [];
					this.pagination.total = 0;
				}

				this.sabayon = result;
			}

			this.loading = false;

			return this.sabayon;
		},
		changeType(value) {
			this.$emit("type", value);
			// this.onSearch()
		},
		// onExport(value) {
		//   this.forms[this.formId];
		// },
		// 复制到粘貼板
		onCopy() {
			// console.log(this.desserts);
			const dataArr = [];

			const rowTopOne = [];
			const rowTopTwo = [];
			for (var i = 0; i < this.cells.headers.length; i++) {
				if (this.cells.headers[i].children == undefined) {
					rowTopOne.push('"' + this.cells.headers[i].text + '"');
					rowTopTwo.push('"' + this.cells.headers[i].text + '"');
				} else {
					for (var j = 0; j < this.cells.headers[i].children.length; j++) {
						if (j == 0) {
							rowTopOne.push('"' + this.cells.headers[i].text + '"');
							rowTopTwo.push('"' + this.cells.headers[i].children[j].text + '"');
							continue;
						}
						rowTopOne.push('"' + this.cells.headers[i].children[j].text + '"');
						rowTopTwo.push('"' + this.cells.headers[i].children[j].text + '"');
					}
				}
			}

			if (rowTopOne.length != 0) {
				for (var k = 0; k < rowTopOne.length; k++) {
					if (rowTopOne[k] == rowTopTwo[k]) {
						rowTopOne[k] = "";
					}
				}
				dataArr.push(rowTopOne.join("\t"));
			}
			dataArr.push(rowTopTwo.join("\t"));

			const copyHeader = [];
			this.cells.headers.forEach(element => {
				if (element.children) {
					element.children.forEach(item => {
						copyHeader.push(item.value);
					});
				} else {
					copyHeader.push(element.value);
				}
			});
			this.desserts.forEach(element => {
				const rowArr = [];
				// for (const key of Object.keys(element)) {
				for (var i = 0; i < this.cells.headers.length; i++) {
					if (this.cells.headers[i].children != undefined) {
						for (var j = 0; j < this.cells.headers[i].children.length; j++) {
							rowArr.push(
								'"' + element[this.cells.headers[i].children[j].value] + '"'
							);
						}
					} else {
						rowArr.push('"' + element[this.cells.headers[i].value] + '"');
					}
				}
				dataArr.push(rowArr.join("\t"));
			});
			// console.log(dataArr.join("\n"));
			copyText(dataArr.join("\n"))
				.then(e => {
					this.toasted.success("复制成功");
				})
				.catch(e => {
					this.toasted.error("复制失败");
				});
		},
		async onExport() {
			const form = this.forms[this.searchFormId];

			if (!R.isYummy(form.date) || !R.isYummy(form.proCode)) {
				this.toasted.warning("沒有选择日期或者项目");
				return;
			}
			const result = await this.$store.dispatch("EXPORT_STAFF_DETAIL", form);
			let res = lpTools.createExcelFun(result, form.proCode);
			this.toasted.info(res.msg);
			if (result.code === 200) {
				// this.onClose();
			}
		}
	},

	computed: {
		...mapGetters(["auth"])
	}
};
</script>
