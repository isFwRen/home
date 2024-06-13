<template>
	<div class="manage">
		<div class="z-flex align-center">
			<z-btn class="mr-3" color="primary" small unlocked @click="selectAll">
				{{ selectedPro.length === auth.proItems.length ? "全不选" : "全选" }}
			</z-btn>

			<z-checkboxs
				:formId="searchFormId"
				formKey="proCode"
				ref="proCode"
				:options="auth.proItems"
				:defaultValue="selectedPro"
				@change="selectProjects"
			></z-checkboxs>
		</div>

		<v-row class="z-flex align-end mb-3">
			<v-col :cols="2">
				<z-text-field :formId="searchFormId" formKey="billName" hideDetails label="案件号">
				</z-text-field>
			</v-col>

			<v-col :cols="2">
				<z-text-field
					:formId="searchFormId"
					formKey="wrongFieldName"
					hideDetails
					label="错误字段"
				>
				</z-text-field>
			</v-col>

			<v-col :cols="2">
				<z-text-field
					:formId="searchFormId"
					formKey="responsibleName"
					hideDetails
					label="责任人"
				>
				</z-text-field>
			</v-col>

			<v-col :cols="2">
				<z-date-picker
					:formId="searchFormId"
					formKey="month"
					hideDetails
					label="月份"
					picker-type="month"
					:defaultValue="cells.DEFAULT_MONTH"
				></z-date-picker>
			</v-col>

			<div class="z-flex">
				<z-btn class="pb-3 pl-3" color="primary" @click="onSearch">
					<v-icon class="text-h6">mdi-magnify</v-icon>
					查询
				</z-btn>
			</div>
		</v-row>

		<div class="z-flex justify-between mb-6">
			<div class="z-flex">
				<z-file-input
					:action="`${baseURLApi}pro-manager/quality/management/push`"
					accept="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
					class="mt-n3 mr-4"
					clearable
					:deleteIcon="false"
					:headers="fileHeaders"
					label="导入数据"
					:multiple="false"
					name="file"
					parcel
					prependIcon="mdi-file-excel-outline"
					width="260"
					@response="handleResponse"
					@click="getUserInfo"
					
				>
				</z-file-input>

				<z-btn color="primary" @click="onExport">
					<v-icon>mdi-export</v-icon>
					导出
				</z-btn>
			</div>

			<z-btn color="primary" @click="onNews">
				<v-icon>mdi-plus</v-icon>
				新增
			</z-btn>
		</div>

		<div class="table mt-3">
			<vxe-table :data="desserts" :max-height="tableMaxHeight" :size="tableSize">
				<template v-for="item in cells.headers">
					<!-- 操作 BEGIN -->
					<vxe-column
						v-if="item.value === 'options'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							<div class="z-flex">
								<z-btn
									class="mr-2"
									color="primary"
									outlined
									small
									@click="onEdit(row)"
									>编辑</z-btn
								>

								<z-btn color="error" outlined small @click="onDelete(row)"
									>删除</z-btn
								>
							</div>
						</template>
					</vxe-column>
					<!-- 操作 END -->

					<!-- 影像 BEGIN -->
					<vxe-column
						v-else-if="item.value === 'imagePath'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row, rowIndex }">
							<z-upload
								formId="paths"
								:formKey="`path${rowIndex}`"
								show-only
								:defaultValue="setImages(row.imagePath)"
							></z-upload>
						</template>
					</vxe-column>
					<!-- 影像 END -->

					<!-- 录入日期 BEGIN -->
					<vxe-column
						v-else-if="item.value === 'entryDate'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							{{ row[item.value] | dateFormat("YYYY-MM-DD") }}
						</template>
					</vxe-column>
					<!-- 录入日期 END -->

					<!-- 反馈日期 BEGIN -->
					<vxe-column
						v-else-if="item.value === 'feedbackDate'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							{{ row[item.value] | dateFormat("YYYY-MM-DD") }}
						</template>
					</vxe-column>
					<!-- 反馈日期 END -->

					<vxe-column
						v-else
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
					</vxe-column>
				</template>
			</vxe-table>
		</div>

		<z-pagination
			:options="pageSizes"
			@page="handlePage"
			:total="pagination.total"
		></z-pagination>

		<z-dynamic-form
			ref="dynamic"
			:formId="formId"
			:title="title"
			:fieldList="cells.fields"
			:config="config"
			:detail="detailInfo"
			:width="600"
			@change:file="handleFileUpload"
			@confirm="handleConfirm"
			@enter="handleEnter"
		>
		</z-dynamic-form>

		<lp-spinners :overlay="overlay">
			<h3 slot="tips">拼命导出中，请耐心等待...</h3>
		</lp-spinners>
	</div>
</template>

<script>
import moment from "moment";
import { mapGetters } from "vuex";
import { localStorage, tools } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import { tools as lpTools } from "@/libs/util";
import cells from "./cells";

const { baseURLApi } = lpTools.baseURL();

const codeToName = new Map([
	["op0ResponsibleCode", "op0ResponsibleName"],
	["op1ResponsibleCode", "op1ResponsibleName"],
	["op2ResponsibleCode", "op2ResponsibleName"],
	["opqResponsibleCode", "opqResponsibleName"]
]);

export default {
	name: "Manage",
	mixins: [TableMixins],

	data() {
		return {
			formId: "manage",
			cells,
			baseURLApi,
			dispatchList: "GET_PM_QUALITIES_MANAGE_LIST",
			config: {},
			title: "新增",
			selectedPro: [],
			overlay: false,
			fileHeaders: {}
		};
	},

	created() {
		this.getUserInfo();

		this.config = {
			proCode: {
				items: this.auth.proItems
			}
		};
	},

	methods: {
		getUserInfo() {
			const { token, user } = localStorage.get(["token", "user"]);
			const secret = localStorage.get("secret");
			let code = "";
			if (secret) {
				code = lpTools.GetCode(secret);
			}
			this.fileHeaders = {
				"x-token": token,
				"x-user-id": user.id,
				"x-code": String(code)
			};
		},

		// 新增
		onNews() {
			this.getDetail();
			this.title = "新增";
			this.$refs.dynamic.open();
		},

		// 导入数据
		handleResponse({ result }) {
			this.toasted.dynamic(result.msg, result.code);
		},

		// 导出
		async onExport() {
			this.overlay = true;

			const form = this.forms[this.searchFormId];
			const params = {
				proCode: form.proCode,
				month: form.month
			};

			const result = await this.$store.dispatch("EXPORT_PM_QUALITIES_MANAGE_EXCEL", params);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				location.href = `${baseURLApi}${result.data.list}`;
			}

			this.overlay = false;
		},

		// 编辑
		onEdit(row) {
			console.log(row);
			row.feedbackDate = moment(row.feedbackDate).format("YYYY-MM-DD");
			row.entryDate = moment(row.entryDate).format("YYYY-MM-DD");
			if (row.imagePath) row.file = [{ url: `${baseURLApi}${row.imagePath[0]}` }];

			this.getDetail(row);
			this.title = "编辑";
			this.$refs.dynamic.open();
		},

		// 删除
		onDelete(row) {
			this.$modal({
				visible: true,
				title: "删除提示",
				content: `请确认是否要删除？`,
				confirm: async () => {
					const body = {
						ids: [row.ID]
					};

					const result = await this.$store.dispatch(
						"DELETE_PM_QUALITIES_MANAGE_ITEM",
						body
					);

					this.toasted.dynamic(result.msg, result.code);

					if (result.code === 200) {
						this.getList();
					}
				}
			});
		},

		selectProjects(values) {
			this.selectedPro = values;
		},

		selectAll() {
			this.$refs.proCode.onSelectAll();
		},

		async handleFileUpload(files) {
			this.detailInfo = {
				...this.detailInfo,
				file: [{ url: await tools.fileToBase64(files[0]) }],
				file0: files[0]
			};
		},

		async handleConfirm(effect, form) {
			const body = {
				...effect,
				...form,
				id: this.detailInfo.ID,
				file: this.detailInfo.file0
			};

			const result = await this.$store.dispatch("POST_PM_QUALITIES_MANAGE_ITEM", body);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.$refs.dynamic.close();
				this.getList();
			}
		},

		async handleEnter({ formKey, event }) {
			if (formKey.includes("ResponsibleCode")) {
				const params = {
					code: event.customValue
				};

				const result = await this.$store.dispatch("GET_PM_QUALITIES_MANAGE_STAFF", params);

				this.detailInfo = {
					...this.detailInfo,
					[codeToName.get(formKey)]: result.data?.list?.[0]
				};
			}
		},

		setImages(paths) {
			return paths?.map(path => ({ url: `${baseURLApi}${path}` }));
		}
	},

	computed: {
		...mapGetters(["auth"])
	}
};
</script>
