<template>
	<div class="business-rules">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col :cols="2">
					<z-text-field
						:formId="searchFormId"
						formKey="name"
						clearable
						hideDetails
						label="钉钉群名称"
					>
					</z-text-field>
				</v-col>

				<v-col :cols="2">
					<z-select
						:formId="searchFormId"
						formKey="proCode"
						clearable
						hideDetails
						label="所属项目编码"
						:options="auth.proItems"
					></z-select>
				</v-col>

				<v-col :cols="2">
					<z-select
						:formId="searchFormId"
						formKey="env"
						clearable
						hideDetails
						label="环境"
						:options="ENV"
					></z-select>
				</v-col>
				<div class="z-flex">
					<z-btn class="pb-3 px-3" color="primary" @click="onSearch">
						<v-icon size="20">mdi-magnify</v-icon>
						查询
					</z-btn>
					<z-btn :formId="searchFormId" btnType="reset" class="pr-3 pb-3" color="error">
						<v-icon class="text-h6">mdi-reload</v-icon>
						重置
					</z-btn>

					<z-btn color="primary" class="pr-3 pb-3" @click="onEditItem()">
						<v-icon>mdi-plus</v-icon>
						新增
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
				</div>
			</v-row>
		</div>

		<div class="table">
			<vxe-table
				:data="desserts"
				:border="tableBorder"
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
								<z-btn
									class="mr-2"
									color="primary"
									depressed
									outlined
									small
									@click="onEditItem(row)"
									v-if="row.status !== 3"
								>
									编辑
								</z-btn>

								<z-btn
									color="error"
									depressed
									outlined
									small
									@click="delectOneLine(row)"
								>
									{{ "删除" }}
								</z-btn>
							</div>
						</template>
					</vxe-column>
					<vxe-column
						v-else-if="item.value == 'env'"
						:key="item.value"
						:field="item.value"
						:title="item.text"
						:width="item.width"
					>
						<template #default="{ row }">
							{{ envObj[row.env] }}
						</template>
					</vxe-column>
					<vxe-column
						v-else
						:key="item.value"
						:field="item.value"
						:title="item.text"
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

		<PTEdit
			ref="showEdit"
			:ENV="{ obj: envObj, arr: ENV }"
			:myDetail="detailInfo"
			@submitted="onSearch"
		/>
	</div>
</template>

<script>
import moment from "moment";
import { mapGetters } from "vuex";
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";
import PTEdit from "./dio/PTEdit.vue";

const date = new Date();
const [year, month, day] = [date.getFullYear(), date.getMonth(), date.getDate()];
const DEFAULT_DATE = [
	moment(`${year}-${month}-${day}`).format("YYYY-MM-DD"),
	moment().format("YYYY-MM-DD")
];

export default {
	name: "Notice",
	mixins: [TableMixins],
	data() {
		return {
			cells,
			DEFAULT_DATE,
			dispatchList: "GET_PT_MESSAGE_TABLE_LIST",
			desserts: [],
			dispatchDelete: "DELETE_PT_MESSAGE_TABLE_ROW",
			ENV: [],
			envObj: {}
		};
	},
	created() {
		this.init();
	},
	methods: {
		init() {
			this.getDefaultValue();
		},
		delectOneLine(row) {
			this.detailInfo = row;
			this.deleteItem();
		},
		onEditItem(row) {
			if (row) {
				this.detailInfo = row;
			} else {
				this.detailInfo = {};
			}
			this.$refs.showEdit.onOpen();
		},
		//初始化常量
		async getDefaultValue() {
			const result = await this.$store.dispatch("MESSAGE_GET_CONSTANTS", {});
			if (result.code == 200) {
				let envs = result.data.env;
				this.envObj = Object.assign({}, envs);
				for (let key in envs) {
					this.ENV.push({ label: envs[key], value: key });
				}
			}
		}
	},
	computed: {
		...mapGetters(["auth"])
	},
	components: { PTEdit }
};
</script>
