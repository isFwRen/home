<template>
	<div class="staff-yield">
		<!-- 全部 BEGIN -->
		<staff-total v-if="currentType === 1" @type="handleType"></staff-total>
		<!-- 全部 END -->

		<!-- 明细 BEGIN -->
		<staff-detail v-else @type="handleType"></staff-detail>
		<!-- 明细 END -->
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";

export default {
	name: "StaffYield",
	mixins: [TableMixins],

	data() {
		return {
			formId: "StaffYield",
			currentType: 1
		};
	},
	methods: {
		handleType(value) {
			this.currentType = value;
			this.params.currentType = value;
		}
	},

	components: {
		staffDetail: () => import("./staffDetail"),
		staffTotal: () => import("./staffTotal")
	}
};
</script>
