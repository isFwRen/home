<template>
	<div class="analysis">
		<div class="mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col cols="2">
					<z-select
						:formId="searchFormId"
						formKey="types"
						hideDetails
						label="类型"
						:options="cells.typeOptions"
						:defaultValue="type"
						@change="onSelect"
					></z-select>
				</v-col>

				<v-col cols="3">
					<z-date-picker
						:formId="searchFormId"
						formKey="date"
						hideDetails
						label="反馈日期"
						range
						:defaultValue="cells.DEFAULT_DATE"
					></z-date-picker>
				</v-col>

				<div class="z-flex">
					<z-btn class="pb-3 pl-3" color="primary" @click="onSearch">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>

					<z-btn class="pb-3 pl-3" color="primary" @click="onExport">
						<v-icon class="text-h6">mdi-export</v-icon>
						导出
					</z-btn>
				</div>
			</v-row>
		</div>

		<div>
			<!-- 按项目 BEGIN -->
			<analysis-project v-if="type === '1'" :list="desserts"></analysis-project>
			<!-- 按项目 END -->

			<!-- 按字段 BEGIN -->
			<analysis-field v-else-if="type === '2'" :list="desserts"></analysis-field>
			<!-- 按字段 END -->

			<!-- 按人员 BEGIN -->
			<analysis-staff v-else :list="desserts"></analysis-staff>
			<!-- 按人员 END -->
		</div>

		<lp-spinners :overlay="overlay">
			<h3 slot="tips">拼命导出中，请耐心等待...</h3>
		</lp-spinners>
	</div>
</template>

<script>
import { tools as lpTools } from "@/libs/util";
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";

const { baseURLApi } = lpTools.baseURL();

export default {
	name: "Analysis",
	mixins: [TableMixins],

	data() {
		return {
			formId: "Analysis",
			cells,
			dispatchList: "GET_PM_QUALITIES_ANALYSIS_LIST",
			type: "1",
			overlay: false
		};
	},

	methods: {
		onSelect(value) {
			this.type = value;
			this.getList();
		},

		async onExport() {
			this.overlay = true;

			const form = this.forms[this.searchFormId];
			const params = {
				date: form.date
			};

			const result = await this.$store.dispatch("EXPORT_PM_QUALITIES_ANALYSIS", params);

			if (result.code === 200) {
				location.href = `${baseURLApi}${result.data.list}`;
			}

			this.overlay = false;
		}
	},

	components: {
		"analysis-project": () => import("./project"),
		"analysis-field": () => import("./field"),
		"analysis-staff": () => import("./staff")
	}
};
</script>
