<template>
	<div class="ruleList">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<div class="none">
					<v-col :cols="3"
						><z-select
							:formId="searchFormId"
							formKey="releaseType"
							label="项目"
							:options="[{ label: '公告列表', text: 2 }]"
							:defaultValue="1"
						></z-select>
					</v-col>
				</div>
				<v-col :cols="3"
					><z-select
						:formId="searchFormId"
						formKey="proCode"
						label="项目编码"
						:options="auth.proItems"
					></z-select>
				</v-col>
				<v-col :cols="3">
					<z-date-picker
						:formId="searchFormId"
						formKey="dayRange"
						label="日期范围"
						prepend-icon="mdi-calendar"
						range
						:defaultValue="DEFAULT_DATE"
					></z-date-picker>
				</v-col>
				<v-col :cols="3">
					<z-text-field :formId="searchFormId" formKey="title" label="标题">
					</z-text-field
				></v-col>
				<div class="z-flex">
					<z-btn class="pb-3 px-3" color="primary" @click="onSearch">
						<v-icon size="20">mdi-magnify</v-icon>
						查询
					</z-btn>
				</div>
			</v-row>
		</div>
		<div class="table">
			<vxe-table :data="desserts" :border="tableBorder" :size="tableSize">
				<vxe-column type="seq" title="序号" width="60"></vxe-column>

				<template v-for="item in cells.headers">
					<vxe-column
						v-if="item.value === 'opitions'"
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
									@click="isWatch(row)"
								>
									查看
								</z-btn>
							</div>
						</template>
					</vxe-column>
					<vxe-column v-else :field="item.value" :title="item.text" :width="item.width">
						<template #default="{ row }">
							<div class="py-2 z-flex">
								{{
									(item.output && item.output(row[item.value])) || row[item.value]
								}}
							</div>
						</template>
					</vxe-column>
				</template>
			</vxe-table>
		</div>
		<z-pagination
			:options="pageSizes"
			@page="handlePage"
			:total="pagination.total"
		></z-pagination>
	</div>
</template>

<script>
import moment from "moment";
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";
import { mapGetters } from "vuex";
const date = new Date();
const [year, month, day] = [date.getFullYear(), date.getMonth(), date.getDate()];
const DEFAULT_DATE = [
	moment(`${year}-${month}-${day}`).format("YYYY-MM-DD"),
	moment().format("YYYY-MM-DD")
];
export default {
	mixins: [TableMixins],

	data() {
		return {
			formId: "ruleList",
			cells,
			dispatchList: "HOME_GET_ANNOUNCEMENT",
			total: 0,
			DEFAULT_DATE
		};
	},
	computed: {
		...mapGetters(["auth"])
	},
	methods: {
		isWatch(row) {
			this.$emit("showDetial", row);
		}
	}
};
</script>

<style lang="scss" scoped>
.none {
	display: none;
}
</style>
