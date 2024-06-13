<template>
	<div class="error-detail">
		<div class="z-flex align-end mb-6 lp-filters">
			<v-row class="z-flex align-end">
				<v-col
					v-for="(item, index) in cells.fields"
					:key="`entry_filters_${index}`"
					:cols="item.cols"
				>
					<template v-if="item.inputType === 'text'">
						<z-text-field
							:formId="searchFormId"
							:formKey="item.formKey"
							:clearable="item.clearable"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:suffix="item.suffix"
							:defaultValue="item.defaultValue"
						>
						</z-text-field>
					</template>

					<template v-else-if="item.inputType === 'date'">
						<z-date-picker
							:formId="searchFormId"
							:formKey="item.formKey"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:options="item.options"
							:range="item.range"
							:suffix="item.suffix"
							z-index="10"
							:defaultValue="item.defaultValue"
						></z-date-picker>
					</template>

					<template v-else>
						<z-select
							:formId="searchFormId"
							:formKey="item.formKey"
							:clearable="item.clearable"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:options="item.options"
							:suffix="item.suffix"
							:defaultValue="item.defaultValue"
						></z-select>
					</template>
				</v-col>

				<div class="z-flex">
					<z-btn class="pb-3 pl-3" color="primary" @click="onSearch">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>
					<z-btn class="pb-3 pl-3" color="error" @click="onExport"> 导出 </z-btn>
				</div>
			</v-row>
		</div>

		<div class="pb-3 btns">
			<z-btn
				class="mr-3"
				color="primary"
				:disabled="!selected.length"
				outlined
				small
				@click="onBatchAcceptOrReject('1')"
			>
				批量通过
			</z-btn>

			<z-btn
				class="mr-3"
				color="primary"
				:disabled="!selected.length"
				outlined
				small
				@click="onBatchAcceptOrReject('2')"
			>
				批量不通过
			</z-btn>
		</div>

		<div class="pb-3">
			<v-tabs v-model="showType">
				<v-tab>全部</v-tab>
				<v-tab>待审核</v-tab>
			</v-tabs>
		</div>

		<div class="table error-detail-table">
			<vxe-table
				:data="desserts"
				:border="tableBorder"
				:max-height="tableMaxHeight"
				:size="tableSize"
				:stripe="tableStripe"
				@checkbox-all="handleSelectAll"
				@checkbox-change="handleSelectChange"
			>
				<vxe-column type="checkbox" width="40"></vxe-column>

				<template v-for="item in cells.headers">
					<!-- 错误数据 BEGIN -->
					<vxe-column
						v-if="item.value === 'wrong'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<span
								v-if="
									_tools.compareString(row.wrong, row.right, 'error--text')
										.targetHtml
								"
								v-html="
									_tools.compareString(row.wrong, row.right, 'error--text')
										.targetHtml
								"
							></span>
							<span v-else>{{ row.wrong }}</span>
						</template>
					</vxe-column>
					<!-- 错误数据 END -->

					<!-- 正确数据 BEGIN -->
					<vxe-column
						v-else-if="item.value === 'right'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<span
								v-if="
									_tools.compareString(row.right, row.wrong, 'error--text')
										.targetHtml
								"
								v-html="
									_tools.compareString(row.right, row.wrong, 'error--text')
										.targetHtml
								"
							></span>
							<span v-else>{{ row.right }}</span>
						</template>
					</vxe-column>
					<!-- 正确数据 END -->

					<!-- 解析 BEGIN -->
					<vxe-column
						v-else-if="item.value === 'analysis'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template>
							<span class="primary--text"> 规则解析 </span>
						</template>
					</vxe-column>
					<!-- 解析 END -->

	

					<!-- 差错审核 BEGIN -->
					<vxe-column
						v-else-if="item.value === 'isWrongConfirm'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							<div v-if="row.isComplain && !row.isWrongConfirm" class="py-2 z-flex">
								<z-btn-toggle
									:formId="formId"
									:formKey="`review_${row.id}`"
									color="primary"
									:options="cells.reviewOptions"
									:defaultValue="row.isWrongConfirm"
									@change="onReview($event, row)"
								>
								</z-btn-toggle>
							</div>
						</template>
					</vxe-column>
					<!-- 差错审核 END -->

					<vxe-column
						v-else-if="item.value === 'fieldName'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<span @click="onShow(row)">{{ row[item.value] }}</span>
						</template>
					</vxe-column>

					<vxe-column
						v-else
						:field="item.value"
						:title="item.text"
						:key="item.value"
					></vxe-column>
				</template>
			</vxe-table>

			<z-pagination
				class="mt-4"
				:total="pagination.total"
				:options="pageSizes"
				:pageNum="page.pageIndex"
				@page="handlePage"
			></z-pagination>
		</div>

		<lp-dialog ref="showImg" :fullscreen="true">
			<div slot="main">
				<!-- <img :src="newBase64" /> -->
				<div class="img-wrapper">
					<lp-images class="preview-img" :src="newBase64" />
				</div>
			</div>
		</lp-dialog>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import { tools } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import { tools as lpTools } from "@/libs/util";
import * as cells from "./cells";

const { baseURLApi } = lpTools.baseURL();

export default {
	name: "ErrorDetail",
	mixins: [TableMixins],

	data() {
		return {
			formId: "ErrorDetail",
			cells,
			dispatchList: "GET_REPORT_ERROR_DETAIL_LIST_ITEM",
			showType: 0,
			manual: true,
			imgUrl: "",
			newBase64: ""
		};
	},

	methods: {
		geneAppealText(row) {

			if (!row.isComplain) {
				return "申诉";
			}
			if (row.isComplain && !row.isWrongConfirm) {
				return "再次申诉";
			}
			if (row.isComplain && row.isWrongConfirm) {
				return "已申诉";
			}

		},
		async onExport() {
			const data = {
				proCode: this.forms[this.searchFormId].proCode,
				startTime: this.forms[this.searchFormId].date[0],
				endTime: this.forms[this.searchFormId].date[1]
			};

			const result = await this.$store.dispatch("EXPORT_ERROR_DETAIL_LIST_ITEM", data);
			lpTools.createExcelFun(
				result,
				`${data.proCode}错误明细${data.startTime}-${data.endTime}`
			);
		},
		async onSearch() {
			this.params = {
				...this.params,
				...this.page
			};

			const total = await this.getList();
			if (typeof total !== "number") {
				return;
			}
			const index = this.page.pageIndex - 1;
			if (index * this.page.pageSize > total) {
				this.page.pageIndex = 1;
			}
		},
		async getList() {
			if (this.dispatchList) {
				const params = {
					...this.effectParams,
					...this.params,
					...this.forms[this.searchFormId],
					isAudit: this.showType === 0 ? true : false
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
		// 批量通过/批量不通过
		onBatchAcceptOrReject(status) {
			const form = this.forms[this.searchFormId];
			const tips =
				status === "1" ? cells.BATCH_ACCEPT : status === "2" ? cells.BATCH_REJECT : {};
			const list = [];

			this.selected.map(item => {
				list.push({ id: item.id || item.ID });
			});

			this.$modal({
				visible: true,
				...tips,
				confirm: async () => {
					const body = {
						proCode: form.proCode,
						wrongConfirm: status,
						list
					};

					const result = await this.$store.dispatch(
						"REVIEW_REPORT_ERROR_DETAIL_LIST_ITEM",
						body
					);

					this.toasted.dynamic(result.msg, result.code);

					if (result.code === 200) {
						this.desserts = [];
						this.getList();
					}
				}
			});
		},

		// 差错审核
		onReview(value, row) {
			if (tools.getType(value) === "undefined") {
				return;
			}

			const form = this.forms[this.searchFormId];
			const tips = value === "1" ? cells.ACCEPT : value === "2" ? cells.REJECT : {};

			this.$modal({
				visible: true,
				...tips,
				confirm: async () => {
					const body = {
						proCode: form.proCode,
						wrongConfirm: value,
						list: [{ id: row.id || row.ID }],
						startTime: this.forms.ErrorDetailSearch.date[0],
						endTime: this.forms.ErrorDetailSearch.date[1]
					};

					const result = await this.$store.dispatch(
						"REVIEW_REPORT_ERROR_DETAIL_LIST_ITEM",
						body
					);

					this.toasted.dynamic(result.msg, result.code);

					this.desserts = [];
					this.getList();
				},

				cancel: async () => {
					this.desserts = [];
					this.getList();
				}
			});
		},

		async onShow(row) {
			this.imgUrl = baseURLApi + "files/" + row.path + row.picture[0];
			let reg = new RegExp("/files/files/", "g");
			let convert = new RegExp("/convert_", "g");

			this.imgUrl = this.imgUrl.replace(reg, "/files/").replace(convert, "/");

			const newBase64 = await lpTools.getTokenImg(this.imgUrl);
			if (newBase64) {
				lpTools.getBase64(newBase64).then(base64String => {
					this.newBase64 = base64String;
				});
				this.$refs.showImg.onOpen();
			}
		},

		// 申诉
		async onAppeal(row) {
			const form = this.forms[this.searchFormId];

			const body = {
				proCode: form.proCode,
				id: row.id || row.ID,
				complainConfirm: true
			};

			const result = await this.$store.dispatch("APPEAL_REPORT_ERROR_DETAIL_LIST_ITEM", body);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.onSearch();
			}
		}
	},

	computed: {
		...mapGetters(["auth"])
	},

	watch: {
		"auth.proItems": {
			handler(items) {
				this.cells.fields[0].options = items;
			},
			deep: true,
			immediate: true
		},
		showType(value) {
			this.onSearch();
		}
	},
	components: {
		"lp-images": () => import("@/components/lp-images")
	}
};
</script>

<style scoped lang="scss">
.col-2 {
	flex: 0 0 13.6666666667%;
	max-width: 13.666667%;
}
.col-3 {
	flex: 0 0 21.6666666667%;
	max-width: 21.666667%;
}
.error-detail {
	width: 100%;
	height: 100%;
}
.img-wrapper {
	width: 80vw;
	height: 88vh;
	margin: 0 auto;
}
.preview-img {
	width: 100%;
	height: 80%;
}
</style>
