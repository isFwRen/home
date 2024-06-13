<template>
	<lp-dialog ref="dialog" title="忘记密码" width="600" @dialog="handleDialog">
		<div slot="main">
			<ul class="mt-4">
				<li class="mb-1">
					<z-text-field
						:formId="formId"
						:formKey="userField.formKey"
						:label="userField.label"
						:validation="userField.validation"
					>
					</z-text-field>
				</li>

				<li class="z-flex align-end mb-4" v-if="!isIntranet">
					<z-text-field
						:formId="formId"
						formKey="captcha"
						class="flex-grow-1"
						label="验证码"
						:validation="[{ rule: 'required', message: '请输入验证码!' }]"
					>
					</z-text-field>
					<z-btn
						class="pb-5 ml-4"
						:color="color"
						:disabled="!validAccount || counting"
						:lockedTime="2500"
						@click="sendCode"
						>{{ text }}</z-btn
					>
				</li>

				<li class="mb-0">
					<z-text-field
						:formId="formId"
						formKey="password"
						label="重置密码"
						disabled
						defaultValue="123456"
					>
					</z-text-field>
				</li>
			</ul>
		</div>

		<div slot="actions">
			<z-btn
				:formId="formId"
				btnType="validate"
				class="pb-4"
				:color="color"
				block
				@click="onConfirm"
				>确认重置密码</z-btn
			>
		</div>
	</lp-dialog>
</template>

<script>
import ButtonMixins from "@/mixins/ButtonMixins";
import DialogMixins from "@/mixins/DialogMixins";
import LoginMixins from "../LoginMixins";

export default {
	name: "ForgotPassword",
	mixins: [ButtonMixins, DialogMixins, LoginMixins],

	data() {
		return {
			formId: "ForgotPassword",
			dispatchForm: "RESET_PASSWORD"
		};
	}
};
</script>
