<template>
	<div class="restart">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col
					v-for="(item, index) in cells.fields"
					:key="`entry_filters_${index}`"
					:cols="2"
				>
					<template v-if="item.inputType === 'input'">
						<z-text-field
							:formId="searchFormId"
							:formKey="item.formKey"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:suffix="item.suffix"
						>
						</z-text-field>
					</template>

					<template v-else-if="item.inputType === 'date'">
						<z-date-picker
							:formId="searchFormId"
							:formKey="item.formKey"
							:clearable="item.clearable"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:options="item.options"
							range
							:suffix="item.suffix"
							z-index="10"
						></z-date-picker>
					</template>

					<template v-else>
						<z-select
							:formId="searchFormId"
							:formKey="item.formKey"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:clearable="item.clearable"
							:options="item.options"
							:suffix="item.suffix"
							@change="handleChange($event, item)"
						></z-select>
					</template>
				</v-col>
			</v-row>

			<z-btn class="pl-3" color="primary" @click="onSearch">
				<v-icon class="text-h6">mdi-magnify</v-icon>
				查询
			</z-btn>
		</div>

		<vxe-table :data="desserts" :size="tableSize">
			<template v-for="item in cells.headers">
				<vxe-column
					v-if="item.value === 'CreatedAt'"
					:field="item.value"
					:title="item.text"
					:key="item.value"
				>
					<template #default="{ row }">
						{{ row.CreatedAt | dateFormat("YYYY-MM-DD HH:mm:ss") }}
					</template>
				</vxe-column>

				<vxe-column v-else :field="item.value" :title="item.text" :key="item.value">
				</vxe-column>
			</template>
		</vxe-table>

		<z-pagination
			:total="pagination.total"
			:options="pageSizes"
			@page="handlePage"
		></z-pagination>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import LogMixins from "../LogMixins";
import cells from "./cells";

export default {
	name: "Restart",
	mixins: [TableMixins, LogMixins],

	data() {
		return {
			formId: "Restart",
			cells,
			dispatchList: "GET_LOG_RESTART_LIST"
		};
	}
};
</script>
