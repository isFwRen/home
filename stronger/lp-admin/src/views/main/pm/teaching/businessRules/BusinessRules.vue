<template>
	<div class="business-rules">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col :cols="2">
					<z-select
						:formId="searchFormId"
						formKey="proCode"
						clearable
						hideDetails
						label="项目编码"
						:options="auth.proItems"
					></z-select>
				</v-col>

				<v-col :cols="2">
					<z-select
						:formId="searchFormId"
						formKey="ruleType"
						clearable
						hideDetails
						label="规则类型"
						:options="cells.ruleTypes"
					></z-select>
				</v-col>

				<div class="z-flex">
					<z-btn class="pb-3 px-3" color="primary" @click="onSearch">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>

					<z-btn color="primary" @click="onAddItem">
						<v-icon>mdi-plus</v-icon>
						新增
					</z-btn>
				</div>
			</v-row>
		</div>

		<div class="table">
			<vxe-table
				:data="desserts"
				:border="tableBorder"
				:max-height="tableMaxHeight"
				:size="tableSize"
			>
				<vxe-column type="seq" title="序号" width="60"></vxe-column>

				<template v-for="item in cells.headers">
					<vxe-column
						v-if="item.value === 'options'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							<div class="py-2 z-flex">
								<z-btn
									class="mr-2"
									color="primary"
									depressed
									outlined
									small
									@click="onEditItem(row)"
								>
									编辑
								</z-btn>

								<z-btn
									class="mr-2"
									color="primary"
									depressed
									outlined
									small
									@click="onView(row)"
								>
									查看
								</z-btn>

								<z-btn
									class="mr-2"
									color="primary"
									depressed
									outlined
									small
									@click="onDownload(row)"
								>
									下载
								</z-btn>

								<z-btn
									color="error"
									depressed
									outlined
									small
									@click="onDeleteItem(row)"
								>
									删除
								</z-btn>
							</div>
						</template>
					</vxe-column>

					<!-- 更新日期 BEGIN -->
					<vxe-column
						v-else-if="item.value === 'UpdatedAt'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							{{ row[item.value] | dateFormat("YYYY-MM-DD") }}
						</template>
					</vxe-column>
					<!-- 更新日期 END -->

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
		>
			<div v-if="detailInfo.docsPath" class="mt-4" slot="bottom">
				<a :href="`${baseURLApi}${detailInfo.docsPath}`" target="_blank">{{ fileName }}</a>
			</div>
		</z-dynamic-form>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import TableMixins from "@/mixins/TableMixins";
import { tools as lpTools } from "@/libs/util";
import cells from "./cells";
import { sessionStorage, tools } from "vue-rocket";
const { baseURLApi } = lpTools.baseURL();

export default {
	name: "TeachingBusinessRules",
	mixins: [TableMixins],

	data() {
		return {
			formId: "TeachingBusinessRules",
			cells,
			baseURLApi,
			manual: true,
			dispatchList: "GET_PM_TEACHING_BUSINESS_RULES_LIST",
			title: "新增",
			config: {}
		};
	},

	created() {
		this.config = {
			proCode: {
				items: this.auth.proItems
			},

			ruleType: {
				items: this.cells.ruleTypes
			}
		};

		this.detailInfo = {
			file: []
		};
	},

	methods: {
		onAddItem() {
			this.getDetail();
			this.title = "新增";
			this.$refs.dynamic.open();
		},

		onEditItem(row) {
			this.getDetail(row);

			this.detailInfo = {
				...this.detailInfo
				// file: [
				//   {
				//     url: `${ baseURLApi }${ row.docsPath }`,
				//     label: row.ruleName
				//   }
				// ]
			};

			this.title = "编辑";
			this.$refs.dynamic.open();
		},

		async onView(row) {
			// window.open(base64String);
			sessionStorage.set("pdfUrl", `${baseURLApi}${row.docsPath}`);
		
			window.open(
				`${location.origin}/normal/pdf-file`,
				"_blank",
				"toolbar=yes, scrollbars=yes, resizable=yes"
			);
		},

		onDownload(row) {
			location.href = `${baseURLApi}${row.docsPath}`;
		},

		onDeleteItem(row) {
			console.log(row,'row')
			this.$modal({
				visible: true,
				title: "删除提示",
				content: `请确认是否要删除？`,
				confirm: async () => {
					const body = {
						proCode: row.proCode,
						ids: [row.model.ID]
					};

					const result = await this.$store.dispatch(
						"DELETE_PM_TEACHING_BUSINESS_RULES_ITEM",
						body
					);

					this.toasted.dynamic(result.msg, result.code);

					if (result.code === 200) {
						this.getList();
					}
				}
			});
		},

		handleFileUpload(file) {
			const maxSize = 5 * 1024;
			const size = file.size / 1024;

			if (size > maxSize) {
				this.toasted.warning("不能超过5M!");
				return;
			}

			this.detailInfo = {
				...this.detailInfo,
				file0: file
				// file: [
				//   {
				//     url: file,
				//     label: file.name
				//   }
				// ]
			};

			console.log(this.detailInfo);
		},

		async handleConfirm(effect, form) {
			const body = {
				...effect,
				...form,
				id: this.detailInfo.ID,
				file: this.detailInfo.file0
			};

			const result = await this.$store.dispatch("POST_PM_TEACHING_BUSINESS_RULES_ITEM", body);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.$refs.dynamic.close();
				this.getList();
			}
		}
	},

	computed: {
		...mapGetters(["auth"]),

		fileName() {
			const name = this.detailInfo.docsPath?.split("/").reverse()[0];
			return name;
		}
	}
};
</script>
