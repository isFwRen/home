<template>
	<div class="click-dialog">
		<lp-dialog ref="dialog" title="查看" width="900" height="600" @dialog="handleDialog">
			<div slot="main">
				<div class="click-content"></div>
			</div>
		</lp-dialog>
	</div>
</template>

<script>
import DialogMixins from "@/mixins/DialogMixins";
export default {
	mixins: [DialogMixins],
	data() {
		return {};
	},
	props: ["content"],
	watch: {
		content: function (nextContent) {
			const timer = setTimeout(() => {
				let content = document.getElementsByClassName("click-content")[0];
				if (content) {
					content.innerHTML = nextContent;
					clearTimeout(timer);
				}
			});
		}
	}
};
</script>

<style lang="scss">
.click-content {
	padding: 10px;
}
</style>
