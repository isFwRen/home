<template>
	<div class="sampling-setting">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col :cols="3">
					<z-date-picker
						:formId="searchFormId"
						formKey="date"
						hideDetails
						label="日期"
						range
						z-index="10"
						:defaultValue="cells.DEFAULT_DATE"
					></z-date-picker>
				</v-col>

				<v-col :cols="2">
					<z-select
						:formId="searchFormId"
						formKey="type"
						hideDetails
						label="类型"
						:options="cells.typeOptions"
						:defaultValue="1"
						@change="onChangeType"
					></z-select>
				</v-col>

				<v-col :cols="2">
					<z-select
						:formId="searchFormId"
						formKey="proCode"
						hideDetails
						label="项目"
						:options="auth.proItems"
					></z-select>
				</v-col>

				<v-col :cols="2">
					<z-select
						:formId="searchFormId"
						formKey="status"
						hideDetails
						label="状态"
						:options="cells.statusOptions"
					></z-select>
				</v-col>

				<div class="z-flex">
					<z-btn class="pb-3 pl-3" color="primary" @click="onSearch">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>
				</div>
			</v-row>
		</div>

		<div class="pb-6 btns">
			<z-btn
				v-for="item of cells.btns"
				:key="item.icon"
				:class="item.class"
				:color="item.color"
				small
				outlined
				@click="handleOptions(item)"
			>
				<v-icon class="text-h6">{{ item.icon }}</v-icon>
				{{ item.text }}
			</z-btn>
		</div>

		<div class="table">
			<vxe-table :data="desserts" :size="tableSize">
				<vxe-column type="checkbox" width="60"></vxe-column>

				<template v-for="item in headers">
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
									class="pr-3"
									color="primary"
									outlined
									small
									@click="onEditItem(row)"
								>
									编辑
								</z-btn>

								<z-btn color="error" outlined small> 停用 </z-btn>
							</div>
						</template>
					</vxe-column>

					<!-- 添加时间 BEGIN -->
					<vxe-column
						v-else-if="item.value === 'CreatedAt'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							{{ row[item.value] | dateFormat("YYYY-MM-DD") }}
						</template>
					</vxe-column>
					<!-- 添加时间 END -->

					<vxe-column
						v-else-if="item.value === 'index'"
						type="seq"
						title="序号"
						:key="item.value"
						width="60"
					></vxe-column>

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

		<!-- 新增/编辑 BEGIN -->
		<z-dynamic-form
			ref="dynamic"
			:config="config"
			:detail="detailInfo"
			:fieldList="fields"
			:formId="formId"
			@change:proCode="handleProCodeChange"
			@dialog="handleDialog"
			@confirm="handleConfirm"
		></z-dynamic-form>
		<!-- 新增/编辑 END -->
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import { tools } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";

export default {
	name: "SamplingSetting",
	mixins: [TableMixins],

	data() {
		return {
			formId: "samplingSetting",
			cells,
			dispatchList: "QUALITY_SAMPLING_SETTING_GET_LIST",
			currentType: 1,
			headers: cells.headers1,
			fields: cells.staffFields,
			config: {},
			staffItems: [],
			manual: true
		};
	},

	created() {
		this.config = {
			proCode: {
				items: this.auth.proItems
			},

			code: {
				disabled: true
			}
		};
	},

	methods: {
		async handleProCodeChange(proCode) {
			const params = {
				pageIndex: 1,
				pageSize: 1500,
				proCode,
				op: "op0"
			};

			const result = await this.$store.dispatch("TASK_GET_STAFF_LIST", params);

			if (result.code === 200) {
				if (tools.isYummy(result.data.list)) {
					result.data.list.map(staff => {
						this.staffItems.push({
							label: staff.name,
							value: staff.code
						});
					});

					this.config = {
						...this.config,
						code: {
							items: this.staffItems
						}
					};
				}
			}
		},

		handleConfirm({}, form) {
			const body = {
				id: this.detailInfo.ID,
				...form
			};

			this.updateListItem(form, "QUALITY_SAMPLING_SETTING_UPDATE_ITEM", body);
		},

		handleDialog(dialog) {
			if (!dialog) return;

			this.$nextTick(() => {
				const { proCode } = this.forms[this.formId];

				proCode && this.handleProCodeChange(proCode);

				this.config = {
					...this.config,
					code: {
						disabled: proCode ? false : true
					}
				};
			});
		},

		// 类型
		onChangeType(value) {
			this.currentType = value;
			this.headers = cells[`headers${value}`];
			this.fields = value === 1 ? cells.staffFields : cells.thunkFields;
		},

		handleOptions({ value }) {
			switch (value) {
				case "new":
					this.getDetail({ currentType: this.currentType });
					this.$refs.dynamic.open({ title: "新增", status: -1 });
					break;
			}
		},

		// 编辑
		onEditItem(row) {
			this.getDetail({ ...row, currentType: this.currentType });
			this.$refs.dynamic.open({ title: "编辑", status: 1 });
		}
	},

	computed: {
		...mapGetters(["auth"])
	}
};
</script>
