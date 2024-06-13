<template>
	<div class="field-rules">
		<div class="mt-0">
			<div class="z-flex align-center">
				<v-row class="z-flex align-end">
					<v-col :cols="2">
						<z-select
							:formId="searchFormId"
							formKey="proCode"
							label="项目编码"
							:options="auth.proItems"
							:validation="[{ rule: 'required', message: '请选择项目编码.' }]"
						></z-select>
					</v-col>

					<v-col :cols="2">
						<z-text-field
							:formId="searchFormId"
							formKey="fieldsName"
							label="字段名称"
						></z-text-field>
					</v-col>

					<v-col :cols="2">
						<z-select
							:formId="searchFormId"
							formKey="rule"
							clearable
							label="录入规则"
							:options="cells.ruleOptions"
						></z-select>
					</v-col>

					<div class="z-flex pb-8">
						<z-btn class="mr-3" color="primary" @click="onSearch">
							<v-icon class="text-h6">mdi-magnify</v-icon>
							查询
						</z-btn>

						<z-btn class="mr-3" color="primary" @click="onImport">
							<v-icon>mdi-import</v-icon>
							批量导入
						</z-btn>

						<z-btn class="mr-3" color="primary" @click="onExport">
							<v-icon>mdi-export</v-icon>
							导出
						</z-btn>

						<z-btn
							class="mr-3"
							color="error"
							:disabled="!isDeleteMore"
							@click="onDeleteMore"
						>
							<v-icon>mdi-trash-can-outline</v-icon>
							批量删除
						</z-btn>
					</div>
				</v-row>
			</div>
		</div>

		<div class="mt-4 table">
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
					<vxe-column
						v-if="item.value === 'options'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<z-btn
								class="mr-2"
								color="primary"
								outlined
								small
								@click="onEditItem(row)"
								>编辑</z-btn
							>

							<z-btn
								class="mr-2"
								color="primary"
								outlined
								small
								@click="onViewItem(row)"
								>查看</z-btn
							>

							<z-btn color="error" outlined small @click="onDeleteItem(row)"
								>删除</z-btn
							>
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
				:options="pageSizes"
				:pageNum="page.pageIndex"
				@page="handlePage"
				:total="pagination.total"
			></z-pagination>
		</div>

		<view-dialog ref="viewDialog" :imagePath="imagePath"></view-dialog>

		<import-dialog
			ref="importDialog"
			:id="id"
			:imagePath="imagePath"
			:limit="limit"
			:path="path"
			:proCode="proCode"
		></import-dialog>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import TableMixins from "@/mixins/TableMixins";
import { tools as lpTools } from "@/libs/util";
import cells from "./cells";

const { baseURLApi } = lpTools.baseURL();

export default {
	name: "FieldRules",
	mixins: [TableMixins],

	data() {
		return {
			formId: "FieldRules",
			cells,
			dispatchList: "GET_PM_TEACHING_FIELD_RULES_LIST",
			dispatchDelete: "DELETE_PM_TEACHING_FIELD_RULES_ITEMS",
			baseURLApi,
			manual: true,
			proCode: "",
			path: "",
			id: "",
			limit: 1,

			imagePath: ""
		};
	},

	methods: {
		// 导入
		onImport() {
			this.proCode = this.forms?.[this.searchFormId]?.proCode;
			this.path = "pro-manager/fieldsRule/upload";
			this.limit = 100;

			this.$refs.importDialog.onOpen({ title: "批量导入" });
		},

		// 导出
		async onExport() {
			const form = this.forms?.[this.searchFormId];

			const body = {
				proCode: form?.proCode,
				rule: form?.rule
			};

			const result = await this.$store.dispatch("EXPORT_PM_TEACHING_FIELD_RULES", body);
			if (result.code === 200) {
				const url = `${baseURLApi}${result.data.list}`;
				const res = await this.$store.dispatch("GET_EXPORT_DATA", url);
				lpTools.createExcelFun(res, `${body.proCode}-字段规则`);
			}
		},

		// 编辑
		onEditItem(row) {
			this.id = row.model.ID;
			this.path = "pro-manager/fieldsRule/edit";
			this.limit = 1;
			this.imagePath = row.rulePicture[0] || "";

			this.$refs.importDialog.onOpen({ title: "编辑" });
		},

		// 查看
		onViewItem(row) {
			this.imagePath = row.rulePicture[0] || "";

			this.$refs.viewDialog.onOpen(row.sysFieldName);
		},

		// 删除
		onDeleteItem(row) {
			this.getDetail(row);
			this.deleteItem();
		}
	},

	computed: {
		...mapGetters(["auth"])
	},

	components: {
		"view-dialog": () => import("./viewDialog"),
		"import-dialog": () => import("./importDialog")
	}
};
</script>
