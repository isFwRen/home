<template>
	<div class="add-dialog">
		<lp-dialog ref="dialog" title="修改目标" width="500" height="500" @dialog="handleDialog">
			<div slot="main">
				<z-text-field
					:formId="fromID"
					formKey="target"
					label="目标"
					:autofocus="true"
					@enter="onConfirm"
					:validation="[{ rule: 'numeric', message: '允许字段为正整数.' }]"
				>
				</z-text-field>
			</div>
			<template slot="bottom"
				><z-btn class="mr-4 right" color="primary" :formId="fromID" @click="onConfirm"
					>确定</z-btn
				></template
			>
		</lp-dialog>
	</div>
</template>
<script>
import DialogMixins from "@/mixins/DialogMixins";

export default {
	name: "ChangeHomeTarget",
	mixins: [DialogMixins],
	data() {
		return {
			fromID: "changeTarget",
			dispatchForm: "HOME_SET_USER_YIELD"
		};
	},
	methods: {
		onConfirm() {
			const form = this.forms[this.fromID];
			const body = {
				...this.rowInfo,
				...form
			};
			const result = this.submit(body);
			const that = this;
			result.then(function () {
				that.$emit("submited");
			});
		}
	}
};
</script>

<style lang="scss" scoped>
.mr-4.right {
	display: flex !important;
	justify-content: right;
}
</style>
