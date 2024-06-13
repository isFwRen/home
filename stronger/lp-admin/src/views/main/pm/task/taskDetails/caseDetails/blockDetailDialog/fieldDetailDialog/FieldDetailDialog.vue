<template>
	<lp-dialog ref="dialog" title="字段明细" transition="dialog-bottom-transition" max-width="900">
		<div class="pt-6" slot="main">
			<vxe-table :data="desserts" :size="tableSize">
				<template v-for="item in cells.headers">
					<vxe-column
						:field="item.value"
						:title="item.text"
						:key="item.value"
					></vxe-column>
				</template>
			</vxe-table>
		</div>
	</lp-dialog>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import DialogMixins from "@/mixins/DialogMixins";
import cells from "./cells";

export default {
	name: "FieldDetailDialog",
	mixins: [TableMixins, DialogMixins],

	data() {
		return {
			formId: "FieldDetailDialog",
			cells,
			dispatchList: "GET_CASE_DETAIL_FEILD_LIST",
			manual: true
		};
	},

	methods: {
		onDetail() {}
	}
};
</script>
