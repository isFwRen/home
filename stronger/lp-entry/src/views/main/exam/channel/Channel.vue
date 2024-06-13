<template>
	<div class="practice">
		<vxe-table :data="desserts.filter(el => el.projectCode == code.slice(1, -1))" :border="tableBorder" :max-height="tableMaxHeight" :size="tableSize" :stripe="tableStripe">
			<template v-for="item in cells.headers">
				<vxe-column
					v-if="item.value === 'options'"
					:field="item.value"
					:title="item.text"
					:key="item.value"
					:width="item.width"
				>
					<template #default="{ row }">
						<v-btn color="primary" depressed rounded @click="enterProject(row)">开始考核</v-btn>
					</template>
				</vxe-column>
				<vxe-column v-else :field="item.value" :title="item.text" :key="item.value" :width="item.width"> </vxe-column>
			</template>
		</vxe-table>

		<task-dialog ref="channel"></task-dialog>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";

export default {
	name: "ExamChannel",
	mixins: [TableMixins],

	data() {
		return {
			cells,
			dispatchList: "EXAM_GET_LIST",
		};
	},

	methods: {
		enterProject(row) {
			// this.$router.replace({ name: 'Channel', query: { proCode: code, op: -1 } })
			// this.$router.replace({ name: "ExamChannel", query: { proCode: "B0114", op: -1 } });
			console.log(row, "rr");
			this.$store.commit("UPDATE_EXAM", { topInfo: "true", config: "true" });
			this.$router.replace({ name: "ExamChannel", query: { proCode: row.projectCode, op: "ope" } });
		}
	},

	components: {
		"task-dialog": () => import("../../entry/channel/taskDialog")
	}
};
</script>

<style lang="scss" scoped>
.v-btn {
	height: 25px !important;
}
</style>