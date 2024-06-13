<template>
	<div class="error-details">
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

				<v-col :cols="2">
					<z-text-field :formId="searchFormId" formKey="code" hide-details label="工号">
					</z-text-field>
				</v-col>

				<v-col :cols="2">
					<z-text-field :formId="searchFormId" formKey="name" hide-details label="姓名">
					</z-text-field>
				</v-col>

				<v-col :cols="2">
					<z-text-field
						:formId="searchFormId"
						formKey="fieldName"
						hide-details
						label="字段名称"
					>
					</z-text-field>
				</v-col>

				<v-col :cols="2">
					<z-text-field
						:formId="searchFormId"
						formKey="createdCode"
						hide-details
						label="抽检人工号"
					>
					</z-text-field>
				</v-col>

				<v-col :cols="2">
					<z-text-field
						:formId="searchFormId"
						formKey="createdName"
						hide-details
						label="抽检人姓名"
					>
					</z-text-field>
				</v-col>

				<div class="z-flex">
					<z-btn class="pb-3 pl-3" color="primary" @click="onSearch">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>
				</div>
			</v-row>
		</div>

		<div class="z-flex mt-n3 mb-4">
			<z-btn class="pr-3" color="primary" outlined>
				<v-icon class="text-h6">mdi-content-copy</v-icon>
				复制
			</z-btn>

			<z-btn color="primary" outlined>
				<v-icon class="text-h6">mdi-export-variant</v-icon>
				导出
			</z-btn>
		</div>

		<div class="table">
			<vxe-table :data="desserts" :size="tableSize">
				<template v-for="item in cells.headers">
					<vxe-column
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
	name: "ErrorDetails",
	mixins: [TableMixins],

	data() {
		return {
			formId: "errorDetails",
			cells,
			dispatchList: "QUALITY_ERROR_DETAILS_GET_LIST",
			manual: true
		};
	},

	computed: {
		...mapGetters(["auth"])
	}
};
</script>
