<template>
	<div class="sampling-data">
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
						formKey="type"
						clearable
						hide-details
						label="类型"
						:options="cells.typeOptions"
						:defaultValue="1"
						@change="onChangeType"
					></z-select>
				</v-col>

				<v-col :cols="2">
					<z-text-field :formId="searchFormId" formKey="code" hide-details label="工号">
					</z-text-field>
				</v-col>

				<v-col :cols="2">
					<z-text-field :formId="searchFormId" formKey="name" hide-details label="姓名">
					</z-text-field>
				</v-col>

				<div class="z-flex">
					<z-btn class="pb-3 pl-3" color="primary" @click="onSearch">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>

					<z-btn class="pb-3 pl-3" color="primary" outlined @click="onCopy">
						<v-icon class="text-h6">mdi-content-copy</v-icon>
						复制
					</z-btn>

					<z-btn class="pb-3 pl-3" color="primary" outlined @click="onExport">
						<v-icon class="text-h6">mdi-export-variant</v-icon>
						导出
					</z-btn>
				</div>
			</v-row>
		</div>

		<div class="table">
			<vxe-table :data="desserts" :size="tableSize">
				<template v-for="item in headers">
					<vxe-column
						v-if="item.value === 'options'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{}">
							<div class="z-flex">
								<z-btn class="pr-3" color="primary" outlined small>
									开始抽查
								</z-btn>

								<z-btn color="error" outlined small> 查看 </z-btn>
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

		<z-pagination
			:options="pageSizes"
			:pageNum="page.pageIndex"
			@page="handlePage"
			:total="pagination.total"
		></z-pagination>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";

export default {
	name: "SamplingData",
	mixins: [TableMixins],

	data() {
		return {
			formId: "samplingData",
			cells,
			dispatchList: "QUALITY_SAMPLING_DATA_GET_LIST",
			manual: true,
			headers: cells.headers1
		};
	},

	methods: {
		onChangeType(value) {
			this.headers = cells[`headers${value}`];
		},

		onCopy() {},

		onExport() {}
	}
};
</script>
