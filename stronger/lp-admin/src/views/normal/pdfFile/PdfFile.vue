<template>
	<div class="pdf_wrap" ref="scrollBox">
		<div class="pdfView">
			<iframe ref="iframe" id="ifream" width="100%" height="100%" :src="pdfUrl"></iframe>
		</div>
	</div>
</template>

<script>
import { tools as lpTools } from "@/libs/util";
import { sessionStorage, tools } from "vue-rocket";
import _ from "lodash";

export default {
	name: "PdfFile",
	data() {
		return {
			pdfUrl: ""
		};
	},
	created() {
		this.setViewPdfWidth();
	},
	mounted() {
		this.setListen();
	},
	unmounted() {
		this.removeListen();
	},
	methods: {
		setListen() {
			this.$refs.scrollBox.addEventListener("scroll", e => {
				const scrollTop = e;
				console.log(scrollTop, this.$refs.scrollBox.scrollHeight, "e");
			});
		},
		removeListen() {
			this.$refs.scrollBox.removeEventListener("scroll", e => {});
		},
		async setViewPdfWidth() {
			const path = sessionStorage.get("pdfUrl") || "";
			const newBase64 = await lpTools.getTokenImg(path);
			if (newBase64) {
				lpTools.getBase64(newBase64).then(base64String => {
					this.pdfUrl = base64String;
				});
			}
		}
	}
};
</script>

<style lang="scss" scoped>
.pdf_wrap {
	// width: 100%;
	// height: 95vh;
	// overflow-y: auto;
}
.pdfView {
	width: 100%;
	height: 95vh;
}
</style>
