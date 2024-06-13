<template>
	<div class="ocr-detail">
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

				<v-col cols="2">
					<z-select
						:formId="searchFormId"
						formKey="proCode"
						hideDetails
						label="项目"
						:options="auth.proItems"
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

				<v-col cols="2">
					<z-text-field
						:formId="searchFormId"
						formKey="code"
						hideDetails
						label="字段名称"
					>
					</z-text-field>
				</v-col>

				<div class="z-flex">
					<z-btn class="pb-3 pl-3" color="primary" @click="onSearch()">
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
			>
				<v-icon class="text-h6">{{ item.icon }}</v-icon>
				{{ item.text }}
			</z-btn>
		</div>

		<div class="table ocr-detail-table">
			<vxe-table :border="tableBorder" :data="desserts" :size="tableSize" align="center">
				<template v-for="item in this.headers">
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
							<vxe-column
								:field="record.value"
								:title="record.text"
								:key="record.value"
								width="70px"
							></vxe-column>
						</template>
					</vxe-colgroup>

					<vxe-column
						v-else
						:field="item.value"
						:title="item.text"
						:key="item.value"
						width="70px"
					></vxe-column>
				</template>
			</vxe-table>
		</div>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";
import { typeOptions } from "../cells";

export default {
	name: "OcrDetail",
	mixins: [TableMixins],

	data() {
		return {
			formId: "ocrDetail",
			dispatchList: "GET_OCR_TOTAL",
			cells,
			typeOptions,
			defaultType: 2,
			headers: cells.headers
		};
	},

	methods: {
		changeType(value) {
			this.$emit("type", value);
		}
	},

	computed: {
		...mapGetters(["auth"])
	}
};
</script>
