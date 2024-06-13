<template>
	<div class="lp-yield">
		<div class="mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col cols="3">
					<z-select
						:formId="searchFormId"
						formKey="proCode"
						hide-details
						label="项目"
						:options="auth.proItems"
					></z-select>
				</v-col>

				<v-col cols="3">
					<z-date-picker :formId="searchFormId" formKey="date" hide-details label="日期" range z-index="10">
					</z-date-picker>
				</v-col>

				<z-btn class="pl-3" color="primary" @click="onSearch">
					<v-icon class="text-h6">mdi-magnify</v-icon>
					查询
				</z-btn>
			</v-row>
		</div>

		<div class="table staff-detail-table">
			<vxe-table :border="tableBorder" :data="desserts" :size="tableSize" :max-height="tableMaxHeight" align="center">
				<template v-for="item in cells.headers">
					<vxe-colgroup
						v-if="
							item.value === 'Summary' ||
							item.value === 'first' ||
							item.value === 'one' ||
							item.value === 'two' ||
							item.value === 'problem'
						"
						align="center"
						:title="item.text"
						:key="item.value"
					>
						<template v-for="record in item.children">
							<vxe-column :field="record.value" :title="record.text" :key="record.value" width="70px"></vxe-column>
						</template>
					</vxe-colgroup>
					<vxe-column
						v-else-if="item.value === 'CreatedAt'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						width="150px"
					>
						<template #default="{ row }">
							{{ row.CreatedAt | dateFormat("YYYY-MM-DD HH:mm:ss") }}
						</template>
					</vxe-column>

					<vxe-column v-else :field="item.value" :title="item.text" :key="item.value" width="80px"></vxe-column>
				</template>
			</vxe-table>
		</div>

		<z-pagination :total="pagination.total" :options="pageSizes" @page="handlePage"></z-pagination>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";

export default {
	name: "luYield",
	mixins: [TableMixins],

	data() {
		return {
			formId: "luYield",
			cells,
			dispatchList: "GET_STAFF_TOTAL"
		};
	},

	computed: {
		...mapGetters(["auth"])
	},

	methods: {}
};
</script>