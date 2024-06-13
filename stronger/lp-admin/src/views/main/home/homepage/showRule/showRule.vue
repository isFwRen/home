<template>
	<div class="show-dialog">
		<lp-dialog
			ref="dialog"
			title="规则动态列表"
			width="900"
			height="700"
			@dialog="handleDialog"
		>
			<div class="main" slot="main">
				<RuleList @showDetial="showDetial"></RuleList>
			</div>
		</lp-dialog>
	</div>
</template>

<script>
import DialogMixins from "@/mixins/DialogMixins";
import RuleList from "./RuleList.vue";

export default {
	name: "ChangeHomeTarget",
	mixins: [DialogMixins],
	components: { RuleList },
	methods: {
		showDetial(content) {
			this.$emit("showDetial", content);
		}
	}
};
</script>

<style lang="scss" scoped></style>
