<template>
	<lp-dialog ref="dialog" title="修改日志" transition="dialog-bottom-transition">
		<div class="pt-6 table" slot="main">
			<vxe-table :data="desserts" :border="tableBorder" :size="tableSize">
				<template v-for="item in cells.headers">
					<vxe-column
						v-if="item.value === 'options'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						align="center"
						min-width="150px"
					>
						<div class="z-flex">
							<z-btn class="pr-3" color="primary" outlined small @click="onEdit">
								编辑
							</z-btn>
							<z-btn class="pr-3" color="error" outlined small> 删除 </z-btn>
							<z-btn class="pr-3" color="primary" outlined small @click="onLog">
								日志
							</z-btn>
						</div>
					</vxe-column>

					<vxe-column
						v-else
						:field="item.value"
						:title="item.text"
						:key="item.value"
						align="center"
					>
					</vxe-column>
				</template>
			</vxe-table>
		</div>
	</lp-dialog>
</template>

<script>
import DialogMixins from "@/mixins/DialogMixins";
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";

export default {
	name: "UpdateRoleDialog",
	mixins: [DialogMixins, TableMixins],

	data() {
		return {
			formId: "updateRoleDialog",
			cells,
			desserts: []
		};
	},

	methods: {
		onConfirm() {}
	}
};
</script>
