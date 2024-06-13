<template>
	<lp-dialog ref="dialog" title="查看日志" transition="dialog-bottom-transition" width="900">
		<div slot="main" class="pt-8 main">
			<vxe-table :border="tableBorder" :data="desserts" :size="tableSize">
				<template v-for="item in cells.headers">
					<vxe-column
						v-if="item.value === 'CreatedAt'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							<div class="py-2 z-flex">
								{{ row[item.value] | dateFormat("YYYY-MM-DD HH:mm:ss") }}
							</div>
						</template>
					</vxe-column>
					<vxe-column
						v-else
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
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
	name: "ChangeLogDialog",
	mixins: [TableMixins, DialogMixins],

	data() {
		return {
			cells,
			desserts: [
				{
					name: "",
					beforeVal: "",
					afterVal: "",
					CreatedAt: "",
					editCode: "",
					editName: ""
				}
			]
		};
	},
	created() {
		// this.getlogs()
	},
	methods: {
		// 获取日志列表
		async getLogs(data) {
			// debugger;
			const result = await this.$store.dispatch("GET_CASE_LOGS", {
				id: data.caseId,
				proCode: data.proCode
			});
			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.desserts = result.data || [];
			}
		}
	}
};
</script>
