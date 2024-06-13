<template>
	<div class="config-audit">
		<div class="z-flex btns align-center">
			<z-text-field
				class="pb-3 mr-4"
				:formId="searchFormId"
				formKey="xmlNodeCode"
				hideDetails
				label="XML代码"
			></z-text-field>
			<z-text-field
				class="pb-3 mr-4"
				:formId="searchFormId"
				formKey="xmlNodeName"
				hideDetails
				label="名称"
			></z-text-field>
			<z-btn class="pl-3" color="primary" @click="onSearch">
				<v-icon class="text-h6">mdi-magnify</v-icon>
				查询
			</z-btn>
			<z-btn class="pl-3" color="primary" @click="onNew">
				<v-icon class="text-h6">mdi-plus</v-icon>
				新增
			</z-btn>

			<z-btn class="pl-3" color="error" :disabled="!isDeleteMore" @click="onDeleteMore">
				<v-icon class="text-h6">mdi-trash-can-outline</v-icon>
				批量删除
			</z-btn>
		</div>

		<div class="table config-audit-table">
			<vxe-table
				:data="desserts"
				:border="tableBorder"
				:row-config="{ height: 54 }"
				:size="tableSize"
				@checkbox-all="handleSelectAll"
				@checkbox-change="handleSelectChange"
			>
				<vxe-column type="checkbox" width="60"></vxe-column>
				<vxe-column type="seq" title="序号" width="60"></vxe-column>

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

					<!--- 只能录入begin--->
					<vxe-column
						v-else-if="item.value === 'onlyInput'"
						:min-width="item.width"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<div class="auto-two-line">
								{{ row.onlyInput }}
							</div>
						</template>
					</vxe-column>
					<!--- 只能录入end--->

					<!--- 只能录入begin--->
					<vxe-column
						v-else-if="item.value === 'notInput'"
						:min-width="item.width"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<div class="auto-two-line">
								{{ row.notInput }}
							</div>
						</template>
					</vxe-column>
					<!--- 只能录入end--->

					<vxe-column
						v-else-if="item.value === 'validations'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
					>
						<template #default="{ row }">
							<span v-for="validationItem in row.validationLabel">
								{{ validationItem }}
							</span>
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
			:detail="detailInfo"
			:fieldList="cells.fields"
			@confirm="handleConfirm"
		></z-dynamic-form>
		<!-- 新增/编辑 END -->
	</div>
</template>

<script>
import { mapGetters, mapState } from "vuex";
import { R } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import ConfigMixins from "@/views/main/pd/config/ConfigMixins";
import cells from "./cells";

export default {
	name: "ConfigAudit",
	mixins: [TableMixins, ConfigMixins],

	data() {
		return {
			formId: "configAudit",
			dispatchList: "GET_CONFIG_AUDIT_LIST",
			dispatchDelete: "DELETE_CONFIG_AUDIT_ITEM",
			cells,
			stage: "",
			validation: "",
			validationLabel: ""
		};
	},

	created() {
		this.getValidations();
	},
	watch: {
		sabayon(newVal) {
			newVal.data.list.forEach(el => {
				let arr = [];
				el.validation.forEach((item, index) => {
					arr[index] = this.validation[item];
				});
				el.validationLabel = arr;
			});
		}
	},
	methods: {
		onNew() {
			const { selectItem } = this.$store.state.auth.project;
			const row = {
				code: this.sabayon.data.maxCode,
				myOrder: this.sabayon.data.total || 1,
				proName: selectItem
			};
			this.getDetail(row);

			this.detailInfo["validation"] = [];

			this.$refs.dynamic.open({ title: "新增", status: -1 });
		},

		onEdit(row) {
			this.getDetail(row);
			this.$refs.dynamic.open({ title: "编辑", status: 1 });
		},

		handleConfirm(effect, form) {
			form = {
				myOrder: this.detailInfo.myOrder,
				id: this.detailInfo.ID,
				proId: this.config.proId,
				status: effect.status,
				...form
			};

			this.updateListItem(form, "UPDATE_CONFIG_AUDIT");
		},
		onMore({ customValue }, row) {
			this.getDetail(row);
			customValue === "delete" && this.deleteItem();
		},

		// 获取校验规则
		async getValidations() {
			const item = this.cells.fields.find(item => item.formKey === "validation");
			item.options = [];

			const result = await this.$store.dispatch("GET_VALIDATIONS");

			if (result.code !== 200) return;

			const validations = result.data;
			this.validation = result.data;
			for (let key in validations) {
				item.options.push({
					label: validations[key],
					value: +key
				});
			}
		}
	},

	computed: {
		...mapGetters(["config"])
	}
};
</script>
<style lang="scss" scoped>
.auto-two-line {
	height: 50px;
	overflow: hidden;
	text-overflow: ellipsis;
	display: -webkit-box;
	-webkit-line-clamp: 2;
	-webkit-box-orient: vertical;
}
</style>
