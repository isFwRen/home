<template>
	<div class="config-quality">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col cols="2">
					<z-autocomplete
						:formId="searchFormId"
						formKey="parentXmlNodeName"
						clearable
						hideDetails
						label="XML节点"
						:options="cells.parentXmlNodeNameOptions"
					></z-autocomplete>
				</v-col>

				<v-col cols="2">
					<z-autocomplete
						:formId="searchFormId"
						formKey="xmlNodeName"
						clearable
						hideDetails
						label="XML代码"
						:options="cells.xmlNodeNameOptions"
					></z-autocomplete>
				</v-col>

				<v-col cols="2">
					<z-text-field
						:formId="searchFormId"
						formKey="fieldName"
						hideDetails
						label="字段名称"
					></z-text-field>
				</v-col>

				<v-col cols="2">
					<z-select
						:formId="searchFormId"
						formKey="belongType"
						clearable
						label="所属信息"
						hideDetails
						:options="cells.belongTypeOptions"
					></z-select>
				</v-col>

				<div class="z-flex pb-3 btns">
					<z-btn class="pl-3" color="primary" @click="onSearch">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>
				</div>
			</v-row>
		</div>

		<div class="z-flex mb-6 btns">
			<z-btn color="primary" @click="onNew">
				<v-icon class="text-h6">mdi-plus</v-icon>
				新增
			</z-btn>

			<z-btn class="pl-3" color="primary" @click="onView">
				<v-icon class="text-h6">mdi-eye-outline</v-icon>
				查看
			</z-btn>

			<z-btn class="pl-3" color="error" :disabled="!isDeleteMore" @click="onDeleteMore">
				<v-icon class="text-h6">mdi-trash-can-outline</v-icon>
				批量删除
			</z-btn>
		</div>

		<div class="table config-field-table">
			<vxe-table
				:data="desserts"
				:border="tableBorder"
				:loading="loading"
				:size="tableSize"
				@checkbox-all="handleSelectAll"
				@checkbox-change="handleSelectChange"
			>
				<vxe-column type="checkbox" width="60"></vxe-column>

				<vxe-column title="序号" width="60">
					<template #default="item">
						{{ increaseSeq(item) }}
					</template>
				</vxe-column>

				<template v-for="item in cells.headers">
					<vxe-column
						v-if="item.value === 'options'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
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
						v-else-if="item.value === 'inputType'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							{{ chineseInputType(row[item.value]) }}
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'belongType'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							{{ chineseBelongType(row[item.value]) }}
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'billInfo'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							{{ chineseBillInfo(row[item.value]) }}
						</template>
					</vxe-column>

					<vxe-column
						v-else-if="item.value === 'beneficiary'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							{{ chineseBeneficiary(row[item.value]) }}
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
				@page="handlePage"
			></z-pagination>
		</div>

		<!-- 新增/编辑 BEGIN -->
		<z-dynamic-form
			ref="dynamic"
			:config="{
				belongType: {
					mutex: [
						{
							formKey: 'beneficiary',
							excludes: [5]
						},

						{
							formKey: 'billInfo',
							excludes: [7]
						}
					]
				}
			}"
			:detail="detailInfo"
			:fieldList="cells.fields"
			@confirm="handleConfirm"
		></z-dynamic-form>
		<!-- 新增/编辑 END -->

		<view-quality ref="viewQuality"></view-quality>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import { R } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import ConfigMixins from "@/views/main/pd/config/ConfigMixins";
import cells from "./cells";

export default {
	name: "ConfigQuality",
	mixins: [TableMixins, ConfigMixins],

	data() {
		return {
			formId: "configQuality",
			dispatchList: "GET_CONFIG_QUALITY_LIST",
			dispatchDelete: "DELETE_CONFIG_QUALITY_ITEM",
			cells,
			parentXmlNodeNameOptions: [],
			xmlNodeNameOptions: [],
			fieldNameOptions: [],

			// 常量
			pending: true,

			belongTypeOptions: [],
			beneficiaryOptions: [],
			billInfoOptions: [],
			inputTypeOptions: [],

			belongTypeItem: {},
			beneficiaryItem: {},
			billInfoItem: {},
			inputTypeItem: {}
		};
	},

	created() {
		this.getConstList();
		this.getXMLList();
		this.getFieldList();
	},

	methods: {
		onView() {
			this.$refs.viewQuality.getConfigQualityFormat();
			this.$refs.viewQuality.onOpen();
		},

		onNew() {
			const row = {
				myOrder: this.sabayon.data.maxOrder || 1
			};
			this.getDetail(row);
			this.$refs.dynamic.open({ title: "新增", status: -1 });
		},

		onEdit(row) {
			this.getDetail(row);
			this.$refs.dynamic.open({ title: "编辑", status: 1 });
		},

		handleConfirm(effect, form) {
			let data = {
				myOrder: this.detailInfo.myOrder,
				id: this.detailInfo.ID,
				proId: this.config.proId,
				proName: this.config.pro.name,
				status: effect.status,
				...form
			};

			const { parentXmlNodeName: parentXmlNodeId, xmlNodeName: xmlNodeId, fieldCode } = data;

			const _parentXmlNodeName = R.find(this.parentXmlNodeNameOptions, parentXmlNodeId).label;
			const _xmlNodeName = R.find(this.xmlNodeNameOptions, xmlNodeId).label;
			const fieldName = R.find(this.fieldNameOptions, fieldCode).label;

			data = {
				...data,
				parentXmlNodeId,
				xmlNodeId,
				parentXmlNodeName: _parentXmlNodeName,
				xmlNodeName: _xmlNodeName,
				fieldName
			};

			this.updateListItem(data, "UPDATE_CONFIG_QUALITY");
		},

		onMore({ customValue }, row) {
			this.getDetail(row);
			customValue === "delete" && this.deleteItem();
		},

		async getFieldList() {
			const result = await this.$store.dispatch("GET_CONFIG_QUALITY_FIELD_LIST", {
				proId: this.storage.get("config").proId
			});

			if (result.code !== 200) return;

			const target = R.find(this.cells.fields, "fieldCode");
			const index = this.cells.fields.indexOf(target);

			for (let item of result.data) {
				this.cells.fields[index].options.push({
					label: item.name,
					value: item.code
				});
			}

			this.fieldNameOptions = this.cells.fields[index].options;
		},

		async getXMLList() {
			const result = await this.$store.dispatch("GET_CONFIG_QUALITY_XML_LIST", {
				proId: this.storage.get("config").proId
			});

			if (result.code !== 200) return;

			const target1 = R.find(this.cells.fields, "parentXmlNodeName");
			const index1 = this.cells.fields.indexOf(target1);

			const target2 = R.find(this.cells.fields, "xmlNodeName");
			const index2 = this.cells.fields.indexOf(target2);

			for (let item of result.data.list.nodeList) {
				this.cells.fields[index1].options.push({
					label: item.name,
					value: item.name
				});

				this.cells.fields[index2].options.push({
					label: item.name,
					value: item.name
				});
			}

			this.parentXmlNodeNameOptions = this.cells.fields[index1].options;
			this.xmlNodeNameOptions = this.cells.fields[index2].options;
		},

		async getConstList() {
			this.pending = true;

			this.belongTypeOptions = [];
			this.beneficiaryOptions = [];
			this.billInfoOptions = [];
			this.inputTypeOptions = [];

			this.belongTypeItem = {};
			this.beneficiaryItem = {};
			this.billInfoItem = {};
			this.inputTypeItem = {};

			const result = await this.$store.dispatch("GET_CONFIG_QUALITY_CONST_LIST");

			if (result.code === 200) {
				const { qualityBelongType, qualityBeneficiary, qualityBillInfo, qualityInputType } =
					result.data;
				this.belongTypeItem = qualityBelongType;
				this.beneficiaryItem = qualityBeneficiary;
				this.billInfoItem = qualityBillInfo;
				this.inputTypeItem = qualityInputType;

				// 所属信息
				for (let key in this.belongTypeItem) {
					this.belongTypeOptions = [
						...this.belongTypeOptions,
						{ label: this.belongTypeItem[key], value: +key }
					];
				}

				// 受益人
				for (let key in this.beneficiaryItem) {
					this.beneficiaryOptions = [
						...this.beneficiaryOptions,
						{ label: this.beneficiaryItem[key], value: +key }
					];
				}

				// 账单
				for (let key in this.billInfoItem) {
					this.billInfoOptions = [
						...this.billInfoOptions,
						{ label: this.billInfoItem[key], value: +key }
					];
				}

				// 输入方式
				for (let key in this.inputTypeItem) {
					this.inputTypeOptions = [
						...this.inputTypeOptions,
						{ label: this.inputTypeItem[key], value: +key }
					];
				}
			}

			this.cells.fields[3].options = this.inputTypeOptions;
			this.cells.fields[4].options = this.belongTypeOptions;
			this.cells.fields[5].options = this.billInfoOptions;
			this.cells.fields[6].options = this.beneficiaryOptions;

			this.pending = false;
		},

		// 所属信息
		chineseBelongType(value) {
			return this.belongTypeItem[value];
		},

		// 受益人
		chineseBeneficiary(value) {
			return this.beneficiaryItem[value];
		},

		// 账单
		chineseBillInfo(value) {
			return this.billInfoItem[value];
		},

		// 输入方式
		chineseInputType(value) {
			return this.inputTypeItem[value];
		}
	},

	computed: {
		...mapGetters(["config"])
	},

	components: {
		"view-quality": () => import("./viewQuality")
	}
};
</script>
