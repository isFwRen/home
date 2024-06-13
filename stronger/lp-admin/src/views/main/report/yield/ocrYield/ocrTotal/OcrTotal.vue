<template>
	<div class="ocr-total">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col cols="2">
					<z-select
						:formId="searchFormId"
						formKey="type"
						hideDetails
						label="类型"
						:options="typeOptions"
						:defaultValue="defaultType"
						@change="changeType"
					>
					</z-select>
				</v-col>

				<v-col cols="3">
					<z-date-picker
						:formId="searchFormId"
						formKey="date"
						hideDetails
						label="日期"
						range
						z-index="10"
					></z-date-picker>
				</v-col>

				<!-- <v-col cols="2">
          <z-text-field
            :formId="searchFormId"
            formKey="code"
            hideDetails
            label="工号/姓名"
          >
          </z-text-field>
        </v-col> -->

				<div class="z-flex">
					<z-btn class="pb-3 pl-3" color="primary" @click="onSearch()">
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>
				</div>
			</v-row>
		</div>

		<div class="pb-6 btns">
			<z-btn class="pr-3" color="primary" small outlined @click="onUpdate"> 复制 </z-btn>
			<z-btn class="pr-3" color="primary" small outlined @click="onUpdate"> 导出 </z-btn>
		</div>

		<div class="table ocr-total-table">
			<vxe-table :border="tableBorder" :data="desserts" :size="tableSize" align="center">
				<template v-for="item in cells.headers">
					<vxe-colgroup
						v-if="
							item.value === 'project' ||
							item.value === 'project1' ||
							item.value === 'project2'
						"
						align="center"
						:title="item.text"
						:key="item.value"
					>
						<template v-for="record in item.children">
							<vxe-column
								:field="record.value"
								:title="record.text"
								:key="record.value"
								:width="record.width"
							></vxe-column>
						</template>
					</vxe-colgroup>

					<vxe-column
						v-else
						:field="item.value"
						:title="item.text"
						:key="item.value"
					></vxe-column>
				</template>
			</vxe-table>
		</div>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import { typeOptions } from "../cells";
import cells from "./cells";

export default {
	name: "OcrTotal",
	mixins: [TableMixins],

	data() {
		return {
			formId: "ocrTotal",
			cells,
			typeOptions,
			defaultType: 1,
			dispatchList: "GET_OCR_TOTAL"
		};
	},
	methods: {
		changeType(value) {
			this.$emit("type", value);
		},

		onUpdate() {
			this.$refs.update.onOpen();
		}
	}
};
</script>
