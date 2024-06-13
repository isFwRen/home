<template>
	<div class="ocr-yield">
		<!-- 全部 BEGIN -->
		<ocr-total v-if="currentType === 1" @type="handleType"></ocr-total>
		<!-- 全部 END -->

		<!-- 明细 BEGIN -->
		<ocr-detail v-else @type="handleType"></ocr-detail>
		<!-- 明细 END -->
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";

export default {
	name: "OcrYield",
	mixins: [TableMixins],

	data() {
		return {
			formId: "ocrYield",
			currentType: 1,
			desserts: []
		};
	},
	methods: {
		handleType(value) {
			this.currentType = value;
			this.params.currentType = value;
		}
	},

	components: {
		ocrDetail: () => import("./ocrDetail"),
		ocrTotal: () => import("./ocrTotal")
	}
};
</script>
