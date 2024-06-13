<template>
	<div class="add-dialog">
		<lp-dialog ref="dialog" :title="title" width="900" height="600" @dialog="handleDialog">
			<div class="pt-6" slot="main">
				<z-select
					:formId="formId"
					formKey="proCode"
					label="项目编码"
					:options="auth.proItems"
					:defaultValue="rowInfo.proCode && rowInfo.proCode + ''"
					:validation="[
						{
							rule: 'required',
							message: '项目编码不允许为空.'
						}
					]"
				>
				</z-select>

				<z-text-field
					:formId="formId"
					formKey="title"
					label="标题"
					:defaultValue="rowInfo.title"
					:validation="[{ rule: 'required', message: '标题不能为空' }]"
				>
				</z-text-field>

				<z-select
					:formId="formId"
					formKey="releaseType"
					label="规则类型"
					:validation="[
						{
							rule: 'required',
							message: '规则类型不允许为空.'
						}
					]"
					:defaultValue="rowInfo.releaseType && rowInfo.releaseType + ''"
					:options="cells.types"
				>
				</z-select>
				<div>
					<Edit @getMsg="getMsg" :content="this.rowInfo.content" />
				</div>
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
import Edit from "./Edit.vue";
import cells from "../cells";

export default {
	name: "NoticeDialog",
	mixins: [DialogMixins],
	data() {
		return {
			formId: "NoticeUpdateDialog",
			dispatchForm: "UPDATE_PM_NOTICE_LIST_ITEM",
			cells
		};
	},
	methods: {
		onConfirm() {
			const form = this.forms[this.formId];
			const body = {
				...this.rowInfo,
				...form
			};
			if (!body.proCode) {
				return;
			}
			this.submit(body);
		},
		getMsg(msg) {
			this.rowInfo.content = msg;
		}
	},

	computed: {
		...mapGetters(["auth"])
	},

	watch: {
		dialog: {
			handler(dialog) {
				if (!dialog) {
					// this.rowInfo.content = ' '
					this.dispatchForm = "UPDATE_PM_NOTICE_LIST_ITEM";
				}
				if (dialog && !this.rowInfo.ID) {
					this.dispatchForm = "ADD_A_NOTICE_ITEM";
				}
			},
			immediate: true
		}
	},

	components: { Edit }
};
</script>
<style scoped>
::v-deep p {
	text-indent: 0 !important;
}
::v-deep .w-e-full-screen-container {
	z-index: 99;
}
</style>
