<template>
	<div class="sampling-statistics">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col :cols="3">
					<z-date-picker
						:formId="searchFormId"
						formKey="date"
						hide-details
						label="日期"
						range
						:defaultValue="cells.DEFAULT_DATE"
					>
					</z-date-picker>
				</v-col>

				<v-col :cols="2">
					<z-select
						:formId="searchFormId"
						formKey="proCode"
						clearable
						hide-details
						label="项目"
						:options="auth.proItems"
					></z-select>
				</v-col>

				<div class="z-flex">
					<z-btn class="pb-3 pl-3" color="primary" @click="onSearch">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>

					<z-btn class="pb-3 pl-3" color="primary" outlined>
						<v-icon class="text-h6">mdi-content-copy</v-icon>
						复制
					</z-btn>

					<z-btn class="pb-3 pl-3" color="primary" outlined>
						<v-icon class="text-h6">mdi-export-variant</v-icon>
						导出
					</z-btn>
				</div>
			</v-row>
		</div>

		<div class="table">
			<vxe-table :data="desserts" :size="tableSize">
				<template v-for="item in cells.headers">
					<!-- 更新日期 BEGIN -->
					<vxe-column
						v-if="item.value === 'UpdatedAt'"
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
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";

export default {
	name: "SamplingStatistics",
	mixins: [TableMixins],

	data() {
		return {
			formId: "samplingStatistics",
			cells,
			dispatchList: "QUALITY_SAMPLING_STATISTICS_GET_LIST",
			manual: true
		};
	},

	computed: {
		...mapGetters(["auth"])
	}
};
</script>
