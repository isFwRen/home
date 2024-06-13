<template>
	<div class="business-detail">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col :cols="2">
					<z-select
						:formId="searchFormId"
						formKey="proCode"
						:hideDetails="true"
						label="项目"
						:options="auth.proItems"
						@change="changeCode"
					>
					</z-select>
				</v-col>

				<v-col :cols="3">
					<z-date-picker
						:formId="searchFormId"
						formKey="date"
						:hideDetails="true"
						label="日期"
						:range="true"
						z-index="10"
						:defaultValue="today"
					></z-date-picker>
				</v-col>

				<div class="z-flex">
					<z-btn class="pb-3 pl-3" color="primary" @click="onSearch">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>
					<z-btn class="pb-3 pl-3" color="primary" @click="onSetting">
						<v-icon class="text-h6">mdi-cog</v-icon>
						设置
					</z-btn>
				</div>
			</v-row>
		</div>

		<div class="pb-6 btns">
			<z-btn class="pr-3" color="primary" small outlined @click="onCopy"> 复制 </z-btn>

			<z-btn class="pr-3" color="primary" small outlined @click="onExport"> 导出 </z-btn>
		</div>

		<div class="table business-detail-table">
			<vxe-table :border="tableBorder" :data="desserts" :size="tableSize">
				<template v-for="item in columns">
					<vxe-column
						v-if="item.value === 'billType'"
						:field="item.value"
						:title="item.label"
						:key="item.value"
						align="center"
						:width="item.width"
					>
						<template #default="{ row }">
							{{ billTypes[row.billType] }}
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'firstExportAt'"
						:field="item.value"
						:title="item.label"
						:key="item.value"
						align="center"
						:width="item.width"
					>
						<template #default="{ row }">
							{{ row.firstExportAt | dateFormat("YYYY-MM-DD HH:mm:ss") }}
						</template>
					</vxe-column>

					<vxe-column
						v-else
						:field="item.value"
						:title="item.label"
						:key="item.value"
						align="center"
						:width="item.width"
					>
					</vxe-column>
				</template>
			</vxe-table>

			<z-pagination
				:options="pageSizes"
				@page="handlePage"
				:total="pagination.total"
			></z-pagination>
		</div>
		<SetColumn ref="setColumnRef" :procode="projectCode" @updateColumn="getColumn" />
	</div>
</template>

<script>
import SetColumn from "./setColumn/index.vue";
import { tools } from "vue-rocket";
import { mapGetters } from "vuex";
import moment from "moment";
import TableMixins from "@/mixins/TableMixins";
import SocketsMixins from "@/mixins/SocketsMixins";
import { copy, copyText } from "clipboard-vue";
import cells from "./cells";

const today = [moment().format("YYYY-MM-DD"), moment().format("YYYY-MM-DD")];

export default {
	name: "BusinessDetail",
	mixins: [TableMixins, SocketsMixins],

	data() {
		return {
			cells,
			today,
			formId: "BusinessDetail",
			socketPath: "projectReport",
			dispatchList: "ITEM_REPORT_GET_LIST",
			manual: true,
			billTypes: cells.billType,
			projectCode: "",
			columns: []
		};
	},
	methods: {
		onSearch() {
			this.page.pageIndex = 1;

			this.params = {
				...this.params,
				...this.page
			};

			this.getColumn();
			this.getList();
		},
		async getColumn() {
			const result = await this.$store.dispatch("REPORT_SETTING_CELL", {
				projectCode: this.projectCode
			});
			if (result.code === 200) {
				this.columns = result.data.map(item => {
					if (item.value === "billNum" || item.value === "diseaseDiagnosis") {
						item.width = 250;
					} else if (item.value === "createAt") {
						item.width = 120;
					} else if (item.value === "batchNum") {
						item.width = 120;
					} else {
						item.width = item.label.length * 23;
					}

					return {
						...item
					};
				});
			}
		},
		async getList() {
			if (this.dispatchList) {
				const params = {
					...this.effectParams,
					...this.params,
					...this.forms[this.searchFormId]
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
		changeCode(value) {
			this.projectCode = value;
		},
		onSetting() {
			if (this.projectCode === "") {
				this.toasted.warning("请先选择设置报表的项目");
			} else {
				this.$refs.setColumnRef.openDialog();
			}
		},
		// 复制
		onCopy() {
			var copyData = "";
			const keyArrr = [];

			cells.headers.map(header => {
				keyArrr.push(header.value);
				copyData += '"\t' + header.text + '"\t';
			});
			copyData += "\n";
			this.desserts?.forEach(element => {
				keyArrr.forEach(e => {
					copyData += '"\t' + element[e] + '"\t';
				});
				copyData += "\n";
			});
			copyText(copyData)
				.then(e => {
					this.toasted.success("复制成功");
				})
				.catch(e => {
					this.toasted.error("复制失败");
				});
		},

		// 导出
		async onExport() {
			const form = this.forms[this.searchFormId];

			if (tools.isYummy(form.date) && tools.isYummy(form.proCode)) {
				const result = await this.$store.dispatch("ITEM_REPORT_EXPORT_EXCEL", form);
				this.downloadFile(form.proCode, result);
			}
		},

		downloadFile(proCode, file) {
			var anchor = document.createElement("a");
			anchor.download = `${proCode}项目明细表.xlsx`;
			anchor.style.display = "none";

			anchor.href = URL.createObjectURL(file);
			document.body.appendChild(anchor);
			anchor.click();
			document.body.removeChild(anchor);
		}
	},

	computed: {
		...mapGetters(["auth"])
	},

	components: {
		SetColumn
	}
};
</script>
