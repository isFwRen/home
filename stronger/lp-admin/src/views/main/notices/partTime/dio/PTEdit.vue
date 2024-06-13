<template>
	<div class="add-dialog">
		<lp-dialog ref="dialog" :title="title" width="500" height="600" @dialog="handleDialog">
			<div class="pt-6" slot="main">
				<z-text-field
					:formId="formId"
					formKey="name"
					label="钉钉群名称"
					:valdetailidation="[{ rule: 'required', message: '标题不能为空' }]"
					:defaultValue="myDetail.name || ''"
				>
				</z-text-field>
				<z-select
					:formId="formId"
					formKey="proCode"
					label="所属项目编码"
					:options="auth.proItems"
					:defaultValue="myDetail.proCode || ''"
					:validation="[
						{
							rule: 'required',
							message: '项目编码不允许为空.'
						}
					]"
				>
				</z-select>
				<z-select
					:formId="formId"
					formKey="env"
					label="环境"
					:options="ENV.arr"
					:defaultValue="myDetail.env + '' || ''"
					:validation="[
						{
							rule: 'required',
							message: '规则类型不允许为空.'
						}
					]"
				>
				</z-select>
				<z-text-field
					:formId="formId"
					formKey="accessToken"
					label="令牌"
					:defaultValue="myDetail.accessToken || ''"
					:validation="[{ rule: 'required', message: '标题不能为空' }]"
				>
				</z-text-field>
				<z-text-field
					:formId="formId"
					formKey="secret"
					label="秘钥"
					:defaultValue="myDetail.secret || ''"
					:validation="[{ rule: 'required', message: '标题不能为空' }]"
				>
				</z-text-field>

				<div></div>
			</div>

			<div class="z-flex" slot="actions">
				<z-btn class="mr-3" color="normal" @click="onClose">取消</z-btn>

				<z-btn :formId="formId" btnType="validate" color="primary" @click="onConfirm"
					>确认</z-btn
				>
			</div>
		</lp-dialog>
	</div>
</template>

<script>
import DialogMixins from "@/mixins/DialogMixins";
import { mapGetters } from "vuex";

export default {
	name: "NoticeDialog",
	mixins: [DialogMixins],
	data() {
		return {
			formId: "PT_MESSAGE_EDIT",
			dispatchForm: "ADD_PT_MESSAGE_TABLE_LIST",
			title: "编辑"
		};
	},
	methods: {
		onConfirm() {
			const form = this.forms[this.formId];
			const body = {
				...this.myDetail,
				...form
			};
			if (body.ID) {
				this.dispatchForm = "EDIT_PT_MESSAGE_TABLE_LIST";
			} else {
				this.dispatchForm = "ADD_PT_MESSAGE_TABLE_LIST";
			}
			this.submit(body);
		}
	},
	props: ["ENV", "myDetail"],

	computed: {
		...mapGetters(["auth"])
	}
};
</script>
