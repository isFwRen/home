<template>
	<lp-dialog
		ref="dialog"
		contentClass="lp-export-validation-tip"
		hide-overlay
		transition="scroll-x-transition"
		width="450"
	>
		<v-card slot="card" class="pa-4 pb-0">
			<v-card-text>
				<div class="width:100%;" v-html="wrongHtml"></div>
			</v-card-text>

			<v-card-actions class="z-card-actions justify-end">
				<z-btn color="primary" outlined small @click="onClose"> 关闭 </z-btn>
			</v-card-actions>
		</v-card>
	</lp-dialog>
</template>

<script>
import DialogMixins from "@/mixins/DialogMixins";
export default {
	name: "ExportValidation",
	mixins: [DialogMixins],
	props: {
		wrongNote: {
			type: String,
			required: false
		}
	},
	computed: {
		wrongHtml() {
			if (!this.wrongNote) {
				return "";
			}
			const { rowInfo } = this.$store.state["pm/case/case"].cases;
			let remark = rowInfo.remark ? rowInfo.remark : " 无";
			let remarkElement = `<span style="display:block;">备注信息：${remark}；</span>`;
			const listNode = this.wrongNote?.replace(/;|；/g, ";")?.split(";");
			let errorElement = "";
			if (listNode.length && listNode.length > 0) {
				listNode.forEach(element => {
					if (element !== "") {
						errorElement += `<span style="display:block;">错误内容：${element}；</span>`;
					}
				});
			}

			return remarkElement + errorElement;
		}
	}
};
</script>

<style lang="scss">
.lp-export-validation-tip {
	position: absolute;
	top: 34px;
	right: 0;
}
</style>
