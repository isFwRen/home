<template>
	<div class="teaching-video">
		<div class="mt-n4">
			<div class="z-flex align-center">
				<v-row class="z-flex align-end">
					<v-col :cols="2">
						<z-select
							:formId="searchFormId"
							formKey="proCode"
							hideDetails
							label="项目编码"
							:options="auth.proItems"
						></z-select>
					</v-col>

					<v-col :cols="2">
						<z-text-field
							:formId="searchFormId"
							formKey="blockName"
							hideDetails
							label="分块名称"
						></z-text-field>
					</v-col>

					<v-col :cols="2">
						<z-select
							:formId="searchFormId"
							formKey="rule"
							clearable
							hideDetails
							label="教学视频"
							:options="cells.teachOptions"
						></z-select>
					</v-col>

					<div class="z-flex pb-2">
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

		<div class="mt-8 table">
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
								:disabled="!tools.isYummy(row.video)"
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
				:pageNum="params.pageIndex"
				@page="handlePage"
				:total="pagination.total"
			></z-pagination>
		</div>

		<view-dialog ref="viewDialog" :proCode="proCode" :video="video"></view-dialog>

		<import-dialog
			ref="importDialog"
			:id="id"
			:limit="limit"
			:proCode="proCode"
			:sysBlockName="sysBlockName"
			:video="video"
		></import-dialog>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import { tools } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import { tools as lpTools } from "@/libs/util";
import cells from "./cells";

const { baseURLApi } = lpTools.baseURL();

export default {
	name: "TeachingVideo",
	mixins: [TableMixins],

	tools,

	data() {
		return {
			formId: "TeachingVideo",
			cells,
			dispatchList: "GET_PM_TEACHING_TEACHING_VIDEO_LIST",
			dispatchDelete: "DELETE_PM_TEACHING_TEACHING_VIDEO_ITEMS",
			manual: true,
			proCode: "",
			id: "",
			sysBlockName: "",
			limit: 1,

			video: []
		};
	},

	methods: {
		// 导入
		onImport() {
			this.proCode = this.forms?.[this.searchFormId]?.proCode;
			this.limit = 5;
			this.video = null;

			this.$refs.importDialog.onOpen({ title: "批量导入" });
		},

		// 导出
		async onExport() {
			const body = {
				proCode: this.forms?.[this.searchFormId]?.proCode
			};

			const result = await this.$store.dispatch("EXPORT_PM_TEACHING_TEACHING_VIDEO", body);

			if (result.code === 200) {
				const url = `${baseURLApi}${result.data.list[0]}`;
				const res = await this.$store.dispatch("GET_EXPORT_DATA", url);
				lpTools.createExcelFun(res, `${body.proCode}-教学视频`);
			}
		},

		// 编辑
		onEditItem(row) {
			this.id = row.model.ID;
			console.log(row, "row");
			this.sysBlockName = row.sysBlockName;
			this.limit = 1;
			this.video = row.video;

			this.$refs.importDialog.onOpen({ title: "编辑" });
		},

		onViewItem(row) {
			this.proCode = this.forms?.[this.searchFormId]?.proCode;
			this.video = row.video;

			this.$refs.viewDialog.onOpen({ title: row.video[0].name });
		},

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
