<template>
	<v-expansion-panel>
		<v-expansion-panel-header>程序开启情况</v-expansion-panel-header>
		<v-expansion-panel-content>
			<div class="table program-table">
				<vxe-table
					:data="desserts"
					:border="tableBorder"
					:loading="loading"
					:size="tableSize"
					@checkbox-all="handleSelectAll"
					@checkbox-change="handleSelectChange"
				>
					<template v-for="item in cells.headers">
						<vxe-column
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:width="item.width"
						>
							<template #default="{ row }">
								<span
									:class="[row[item.value] ? 'success--text' : 'error--text']"
									>{{ row[item.value] ? "开启" : "关闭" }}</span
								>
							</template>
						</vxe-column>
					</template>
				</vxe-table>
			</div>
		</v-expansion-panel-content>
	</v-expansion-panel>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import ConfigMixins from "@/views/main/pd/config/ConfigMixins";
import cells from "./cells";

export default {
	name: "ProgramTable",
	mixins: [TableMixins, ConfigMixins],

	data() {
		return {
			dispatchList: "GET_CONFIG_PATH_ENABLE_LIST",
			cells
		};
	}
};
</script>

<style scoped lang="scss">
.program-table {
	width: 100%;
}
</style>
