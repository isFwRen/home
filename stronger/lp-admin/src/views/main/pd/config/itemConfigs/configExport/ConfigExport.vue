<template>
	<div class="config-export">
		<div class="mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col cols="2">
					<z-text-field
						:formId="searchFormId"
						formKey="name"
						hideDetails
						label="代码"
					></z-text-field>
				</v-col>

				<v-col cols="2">
					<z-text-field
						:formId="searchFormId"
						formKey="fieldLike"
						hideDetails
						label="字段"
					></z-text-field>
				</v-col>

				<div class="z-flex pb-3 btns">
					<z-btn class="pl-3" color="primary" @click="onSearch">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>

					<z-btn
						:formId="searchFormId"
						btnType="reset"
						class="pl-3"
						color="warning"
						@click="onSearch"
					>
						<v-icon class="text-h6">mdi-reload</v-icon>
						重置
					</z-btn>

					<z-btn
						class="pl-3"
						color="error"
						:disabled="!isDeleteMore"
						@click="onDeleteMore"
					>
						<v-icon class="text-h6">mdi-trash-can-outline</v-icon>
						批量删除
					</z-btn>

					<z-btn class="pl-3" color="primary" @click="onNew">
						<v-icon class="text-h6">mdi-plus</v-icon>
						新增
					</z-btn>

					<z-btn class="pl-3" color="primary" :lockedTime="0" @click="onExpand">
						<template v-if="!expanded">
							<v-icon class="text-h6">mdi-chevron-down</v-icon>
							更多
						</template>
						<template v-else>
							<v-icon class="text-h6">mdi-chevron-up</v-icon>
							收起
						</template>
					</z-btn>
				</div>

				<v-spacer></v-spacer>

				<v-col cols="2">
					<z-select
						:formId="searchFormId"
						formKey="xmlType"
						hideDetails
						label="XML编码"
						:options="cells.xmlOptions"
						:defaultValue="xmlType"
						@change="onChangeXmlTpye"
					></z-select>
				</v-col>
			</v-row>

			<div class="z-flex align-center pt-4 chunk" v-if="expanded">
				<label class="pt-4 mr-3 fw-bold">将代码</label>
				<z-text-field
					:formId="exchangeFormId"
					formKey="startIndex"
					class="mb-n6"
					:validation="[
						{ rule: 'required', message: '序号不能为空.' },
						{ rule: 'numeric', message: '序号只能为正整数.' }
					]"
					label="序号"
				>
				</z-text-field>
				<label class="pt-4 mx-3 fw-bold">插入到</label>
				<z-text-field
					:formId="exchangeFormId"
					formKey="endIndex"
					class="mb-n6"
					:validation="[
						{ rule: 'required', message: '序号不能为空.' },
						{ rule: 'numeric', message: '序号只能为正整数.' }
					]"
					label="序号"
				>
				</z-text-field>

				<z-btn
					:formId="exchangeFormId"
					btnType="validate"
					class="ml-3 mt-4 mr-3"
					color="primary"
					@click="onExchange"
				>
					确定插入</z-btn
				>

				<z-btn class="mt-4 mr-3" color="primary" @click="onExport">
					<v-icon class="text-h6">mdi-export-variant</v-icon>
					导出
				</z-btn>

				<z-btn class="mt-4" color="primary" @click="viewConfig"> 查看配置 </z-btn>

				<v-btn btnType="validate" class="ml-3 mt-4" color="success" @click="onExportField">
					同步导出配置</v-btn
				>
				
			</div>
		</div>

		<div class="table config-export-table">
			<vxe-table
				:data="desserts"
				:border="tableBorder"
				:size="tableSize"
				@checkbox-all="handleSelectAll"
				@checkbox-change="handleSelectChange"
			>
				<vxe-column type="checkbox" width="60"></vxe-column>

				<!-- 2022/10/21 泽如说排序改为累加排序 -->
				<!-- <vxe-column title="序号" width="60">
          <template #default="item">
            {{ increaseSeq(item) }}
          </template>
        </vxe-column> -->

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
								<z-btn color="primary" depressed small @click="onEdit(row)">
									编辑
								</z-btn>

								<lp-dropdown
									class="pl-3"
									color="primary"
									depressed
									offset-y
									small
									:options="cells.moreOptions"
									@click="onMore($event, row)"
								>
									更多
									<v-icon>mdi-chevron-down</v-icon>
								</lp-dropdown>
							</div>
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'threeFields'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							{{ showText(fieldList, row.threeFields) }}
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'myOrder'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							{{ row.myOrder }}
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'twoFields'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							{{ showText(fieldList, row.twoFields) }}
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'oneFields'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							{{ showText(fieldList, row.oneFields) }}
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'myType'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							{{ row.myType | myTypeText }}
						</template>
					</vxe-column>

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

			<z-pagination
				:total="pagination.total"
				:options="pageSizes"
				:pageNum="page.pageIndex"
				@page="handlePage"
			></z-pagination>
		</div>

		<!-- 新增/编辑 BEGIN -->
		<z-dynamic-form
			ref="dynamic"
			:detail="detailInfo"
			:fieldList="cells.fields"
			@confirm="handleConfirm"
		></z-dynamic-form>
		<!-- 新增/编辑 END -->

		<view-config-dialog ref="viewConfig" :xml="xmlValue"></view-config-dialog>
	</div>
</template>

<script>
import { tools } from "vue-rocket";
import { mapGetters } from "vuex";
import TableMixins from "@/mixins/TableMixins";
import ConfigMixins from "@/views/main/pd/config/ConfigMixins";
import cells from "./cells";
import { tools as lpTools } from "@/libs/util";

const fieldKeys = ["oneFields", "twoFields", "threeFields"];

export default {
	name: "ConfigExport",
	mixins: [TableMixins, ConfigMixins],

	data() {
		return {
			formId: "ConfigExport",
			exchangeFormId: "ConfigExportExchange",
			dispatchList: "GET_CONFIG_EXPORT_LIST",
			dispatchDelete: "DELETE_CONFIG_EXPORT_ITEM",
			cells,
			fieldList: [],
			exportId: null,
			xmlValue: "",
			xmlType: "utf-8"
		};
	},

	created() {
		this.getFields();
	},

	methods: {
		async onExportField() {
			const body = {
				proCode: this.config.pro.code,
				mtype: "export",
				templateId: ""
			};

			const result = await this.$store.dispatch("EXPORT_OR_IMPORT_EXPORT", body);
			if (result.code === 200) {
				this.toasted.dynamic(result.msg, result.code);
			}
		},
		async getInsertDatalist(pageIndex) {
			const params = {
				...this.effectParams,
				pageIndex: pageIndex,
				pageSize: this.page.pageSize
			};
			const result = await this.$store.dispatch(this.dispatchList, params);

			if (result.code === 200) {
				return result.data.list.nodeList;
			} else {
				return [];
			}
		},

		async onChangeXmlTpye(xmlType) {
			const { proId } = this.effectParams;
			const { ID: id, proName, tempVal } = this.sabayon.data.list;

			const form = {
				xmlType,
				proId,
				id,
				proName,
				tempVal
			};

			const result = await this.$store.dispatch("UPDATE_CONFIG_EXPORT_TEMPLATE", form);
			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.getList();
			}
		},

		onNew() {
			const row = {
				myOrder: this.sabayon.data.total + 1 || 1,
				exportId: this.exportId
			};

			this.getDetail(row);

			this.$refs.dynamic.open({ title: "新增", status: -1 });
		},

		onEdit(row) {
			this.getDetail(row);
			this.$refs.dynamic.open({ title: "编辑", status: 1 });
		},

		handleConfirm(effect, form) {
			const oneFieldsName = tools.find(this.fieldList, { value: form.oneFields })?.label;
			const twoFieldsName = tools.find(this.fieldList, { value: form.twoFields })?.label;
			const threeFieldsName = tools.find(this.fieldList, { value: form.threeFields })?.label;

			const body = {
				myOrder: this.detailInfo.myOrder,
				id: this.detailInfo.ID,
				exportId: this.detailInfo.exportId,
				proId: this.config.proId,
				status: effect.status,
				oneFieldsName,
				twoFieldsName,
				threeFieldsName,
				...form
			};

			this.updateListItem(body, "UPDATE_CONFIG_EXPORT");
		},

		onMore({ customValue }, row) {
			this.getDetail(row);
			customValue === "delete" && this.deleteItem();
		},

		async onExchange() {
			const { startIndex, endIndex } = this.forms[this.exchangeFormId];
			let row = null;
			const startOrder = +startIndex;
			const endOrder = +endIndex;

			if (startIndex > this.page.pageSize) {
				let pageIndex = Math.ceil(startIndex / this.page.pageSize);
				let desserts = await this.getInsertDatalist(pageIndex);
				row = desserts.find(d => d.myOrder === +startOrder);
			} else {
				row = this.desserts.find(d => d.myOrder === +startOrder);
			}

			const form = {
				startOrder,
				endOrder,
				exportId: this.exportId,
				startId: row.ID
			};

			const result = await this.$store.dispatch("INSERT_CONFIG_EXPORT_ITEM", form);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.getList();
			}
		},

		async onExport() {
			const res = await this.$store.dispatch("EXPORT_CONFIG_EXPORT", this.exportId);
			console.log(res);
			// lpTools.createExcelFun(res,'导出配置')
			this.downloadFile(res);
		},
		downloadFile(file) {
			var anchor = document.createElement("a");
			anchor.download = "配置.xlsx";
			anchor.style.display = "none";

			anchor.href = URL.createObjectURL(file);
			document.body.appendChild(anchor);
			anchor.click();
			document.body.removeChild(anchor);
		},
		viewConfig() {
			this.$refs.viewConfig.onOpen();
		},

		async getFields() {
			const form = {
				proId: this.effectParams.proId
			};

			const list = await this.$store.dispatch("GET_CONFIG_EXPORT_FIELDS", form);

			this.fieldList = list;

			for (let item of this.cells.fields) {
				if (~fieldKeys.indexOf(item.formKey)) {
					item.options = this.fieldList;
				}
			}
		}
	},

	computed: {
		...mapGetters(["config"])
	},

	watch: {
		"sabayon.data.list": {
			handler(list) {
				if (list) {
					this.desserts = list.nodeList;
					this.exportId = list.ID;
					this.xmlValue = list.tempVal;
					this.xmlType = list.xmlType;
				}
			},
			immediate: true,
			deep: true
		}
	},

	filters: {
		myTypeText(value) {
			const myTypeCN = ["开始", "结束", "固定值", "录入值"];
			return myTypeCN[value - 1];
		}
	},

	components: {
		"view-config-dialog": () => import("./viewConfigDialog")
	}
};
</script>
