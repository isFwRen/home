<template>
	<div class="pdf_wrap overflow-y-auto" ref="scrollBox" id="scroll-target" style="max-height: 100vh">
		<v-card
			style="padding-bottom: 50px; box-shadow: 0 0 0 0 !important; position: relative"
			class="scrollContainer"
			v-scroll:#scroll-target="onScroll"
		>
			<vue-office-pdf :src="pdfUrl" />
			<div class="btn_box" v-if="postForms.isRequired === 1">
				<div class="z-flex justify-center">
					<v-btn color="primary" @click="studyDone"> 完成学习 </v-btn>
				</div>
			</div>
		</v-card>
	</div>
</template>

<script>
import VueOfficePdf from "@vue-office/pdf";
import { sessionStorage, tools } from "vue-rocket";
import _ from "lodash";

export default {
	name: "PdfFile",
	data() {
		return {
			pdfUrl: "",
			dispatchList: "FILISH_FILE_READ",
			postForms: {},
			broadListen: null
		};
	},
	created() {
		this.setViewPdfWidth();
		this.broadListen = new BroadcastChannel("updateLearn");
	},

	methods: {
		onScroll() {},
		async studyDone() {
			const res = await this.$store.dispatch(this.dispatchList, {
				finishId: this.postForms.filishId,
				projectCode: this.postForms.projectCode,
				trainType: this.postForms.trainType
			});

			this.toasted.dynamic(res.msg, res.code);
			if (res.code === 200) {
				this.broadListen.postMessage({ id: this.postForms.filishId });
				window.close();
			}
		},
		setViewPdfWidth() {
			this.pdfUrl = sessionStorage.get("pdfUrl") || "";
			this.postForms = sessionStorage.get("post_data");
		}
	},
	destroy() {
		this.broadListen.close();
	},
	components: {
		VueOfficePdf
	}
};
</script>

<style lang="scss" scoped>
.pdf_wrap {
	width: 100%;
	position: relative;
	.btn_box {
		width: 100%;
		position: absolute;
		bottom: 5%;
	}
}
</style>
