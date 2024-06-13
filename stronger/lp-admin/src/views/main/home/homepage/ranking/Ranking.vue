<template>
	<div class="show-dialog">
		<lp-dialog ref="dialog" title="排行榜" width="900" height="700" @dialog="handleDialog">
			<div class="main" slot="main">
				<v-row class="z-flex align-end">
					<v-col :cols="10">
						<z-radios
							:formId="searchFormId"
							formKey="rankingType"
							label="类型"
							:rules="[{ required: true, message: '排行版了类型' }]"
							:options="[
								{ value: 0, label: '日排行榜' },
								{ value: 1, label: '月排行榜' }
							]"
							:defaultValue="0"
							@change="onSearch"
						>
						</z-radios> </v-col
				></v-row>
				<div class="table">
					<vxe-table :data="desserts" :border="tableBorder" :size="tableSize">
						<template v-for="item in headers">
							<vxe-column :field="item.value" :title="item.text" :width="item.width">
								<template #default="{ row, rowIndex }">
									<div class="py-2 z-flex">
										{{
											(item.output && item.output(rowIndex)) ||
											row[item.value]
										}}
									</div>
								</template>
							</vxe-column>
						</template>
					</vxe-table>
				</div>
				<z-pagination
					:pageNum="page.pageIndex"
					:options="pageSizes"
					@page="handlePage"
					:total="pagination.total"
				></z-pagination>
			</div>
		</lp-dialog>
	</div>
</template>

<script>
import DialogMixins from "@/mixins/DialogMixins";
import TableMixins from "@/mixins/TableMixins";

export default {
	name: "ChangeHomeTarget",
	mixins: [DialogMixins, TableMixins],

	data() {
		return {
			formId: "homeRanking",
			headers: [
				{
					text: "名次",
					value: "myOrder",
					output: e => {
						return (this.page.pageIndex - 1) * this.page.pageSize + e + 1;
					}
				},
				{ text: "工号", value: "userCode" },
				{ text: "姓名", value: "userName" },
				{ text: "产量", value: "value" }
			],
			desserts: [{}],
			dispatchList: "HOME_GET_RANK_YIELD"
		};
	},
	methods: {
		rankingTypeChange(e) {}
	}
};
</script>

<style lang="scss" scoped></style>
