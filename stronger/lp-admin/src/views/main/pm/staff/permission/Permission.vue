<template>
	<div class="permission">
		<div class="table config-audit-table">
			<vxe-table
				:border="tableBorder"
				:max-height="tableMaxHeight"
				:data="desserts"
				:size="tableSize"
				:stripe="tableStripe"
			>
				<template v-for="item in cells.headers">
					<vxe-colgroup
						v-if="item.value === 'NOP'"
						align="center"
						:title="item.text"
						:key="item.value"
					>
						<template v-for="record in item.children">
							<vxe-column
								:field="record.value"
								:title="record.text"
								:key="record.value"
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
import cells from "./cells";

export default {
	name: "Permission",
	mixins: [TableMixins],

	data() {
		return {
			formId: "Permission",
			dispatchList: "GET_STAFF_PERMISSION_LIST",
			cells
		};
	}
};
</script>

<style scoped lang="scss">
.table-header {
	.tr2 {
		width: 75%;
		color: rgba(0, 0, 0, 0.6);
		text-align: center;
		font-size: 0.75rem;
		font-weight: bold;
	}
}
</style>
