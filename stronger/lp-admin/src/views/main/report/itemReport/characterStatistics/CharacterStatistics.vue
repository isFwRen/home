<template>
	<div class="character-statistics">
		<div class="z-flex align-end mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col
					v-for="(item, index) in cells.fields"
					:key="`entry_filters_${index}`"
					:cols="item.cols"
				>
					<template v-if="item.inputType === 'date'">
						<z-date-picker
							:formId="searchFormId"
							:formKey="item.formKey"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:options="item.options"
							:range="item.range"
							:suffix="item.suffix"
							z-index="10"
							:defaultValue="item.defaultValue"
						></z-date-picker>
					</template>

					<template v-else>
						<z-select
							:formId="searchFormId"
							:formKey="item.formKey"
							:clearable="item.clearable"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:options="auth.proItems"
							:suffix="item.suffix"
							:defaultValue="auth.proItems[0].value"
						></z-select>
					</template>
				</v-col>

				<div class="z-flex">
					<z-btn
						:formId="searchFormId"
						btnType="validate"
						class="pb-3 pl-3"
						color="primary"
						@click="onSearch"
					>
						<v-icon class="text-h6">mdi-magnify</v-icon>
						查询
					</z-btn>
				</div>

				<div class="z-flex">
					<z-btn class="pb-3 pl-3" color="primary" @click="copyList"> 复制 </z-btn>
				</div>
				<div class="z-flex">
					<z-btn class="pb-3 pl-3" color="primary" @click="exportList"> 导出 </z-btn>
				</div>
			</v-row>
		</div>

		<div class="table character-statistics-table">
			<vxe-table
				:border="tableBorder"
				:data="desserts"
				:size="tableSize"
				:sort-config="{ defaultSort: { field: 'sumDate', order: 'desc' } }"
			>
				<template v-for="item in cells.headers">
					<vxe-colgroup
						v-if="
							item.value === 'summary' ||
							item.value === 'claimAverageSettlement' ||
							item.value === 'claimAverageEntry' ||
							item.value === 'documentAverageSettlement' ||
							item.value === 'documentAverageEntry'
						"
						align="center"
						:title="item.text"
						:key="item.value"
					>
						<template v-for="record in item.children">
							<vxe-column
								:field="record.value"
								:width="record.width"
								:title="record.text"
								align="center"
							>
								<template #default="{ row }">
									<span>{{
										(item.output && item.output(row[record.value])) ||
										row[record.value]
									}}</span>
								</template>
							</vxe-column>
						</template>
					</vxe-colgroup>

					<vxe-column
						v-else
						:width="item.width"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:sortable="item.sortable"
						align="center"
					>
						<template #default="{ row }">
							<div class="py-2 z-flex">
								{{
									(item.outPut && item.outPut(row[item.value])) || row[item.value]
								}}
							</div>
						</template>
					</vxe-column>
				</template>
			</vxe-table>

			<z-pagination
				:total="pagination.total"
				:options="pageSizes"
				:pageNum="page.pageIndex"
				@page="handlePage"
			></z-pagination>
		</div>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";
import { copy, copyText } from "clipboard-vue";
import { mapGetters } from "vuex";
import moment from "moment";

const today = [moment().format("YYYY-MM-DD"), moment().format("YYYY-MM-DD")];

export default {
	name: "CharacterStatistics",
	mixins: [TableMixins],

	data() {
		return {
			formId: "characterStatistics",
			searchFormId: "characterStatistics",
			cells,
			today,
			desserts: [],
			dispatchList: "GET_CHARACTER_STATISTICS_LIST"
		};
	},

	methods: {
		exportList() {
			const form = this.forms[this.searchFormId];
			this.$store.dispatch("EXPORT_CHARACTER_STATISTICS_LIST", form);
		},
		copyList() {
			var copyData = "";
			const keyArrr = [];

			cells.headers.map(header => {
				keyArrr.push(header.value);
				copyData += "\t" + header.text + "\t";
			});
			copyData += "\n";
			this.desserts?.forEach(element => {
				keyArrr.forEach(e => {
					copyData += "\t" + element[e] + "\t";
				});
				copyData += "\n";
			});
			copyText(copyData)
				.then(e => {
					this.toasted.success("复制成功");
				})
				.catch(e => {
					this.toasted.error("复制失败");
				});
		}
	},
	computed: {
		...mapGetters(["auth"])
	}
};
</script>
