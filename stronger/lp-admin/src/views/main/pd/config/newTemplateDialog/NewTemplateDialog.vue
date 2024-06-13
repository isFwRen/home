<template>
	<div class="new-template-dialog">
		<lp-dialog ref="dialog" :title="title" width="500" @dialog="handleDialog">
			<div class="pt-4" slot="main">
				<z-text-field
					:formId="formId"
					formKey="name"
					label="模板名称"
					:validation="[{ rule: 'required', message: '模板名称不能为空.' }]"
					:defaultValue="tempItem.label ? `${tempItem.label}副本` : undefined"
				>
				</z-text-field>
			</div>

			<div class="z-flex mt-n6" slot="actions">
				<z-btn class="mr-3" color="normal" @click="onClose">取消</z-btn>

				<z-btn :formId="formId" btnType="validate" color="primary" @click="onConfirm"
					>确认</z-btn
				>
			</div>
		</lp-dialog>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import DialogMixins from "@/mixins/DialogMixins";

export default {
	name: "NewTempDialog",
	mixins: [DialogMixins],

	props: {
		tempType: {
			validator(value) {
				return ~["new", "copy"].indexOf(value);
			},
			default: "new"
		},

		tempItem: {
			type: Object,
			default: () => ({ label: undefined })
		}
	},

	data() {
		return {
			formId: "NewTempDialog",
			dispatchForm: ""
		};
	},

	methods: {
		onConfirm() {
			if (this.tempType === "new") {
				this.newTemp();
			} else {
				this.copyTemp();
			}
		},

		// 新增模板
		newTemp() {
			this.dispatchForm = "ADD_CONFIG_ITEM_TEMPLATE";

			const form = {
				proId: this.rowInfo.ID,
				...this.forms[this.formId]
			};

			this.submit(form);
		},

		// 复制模板
		async copyTemp() {
			this.dispatchForm = "COPY_CONFIG_TEMP";

			const form = {
				proId: this.config.proId,
				id: this.tempItem.value,
				...this.forms[this.formId]
			};

			this.submit(form);

			// const result = await this.submit(form)

			// console.log(result)

			// if(result.code === 200) {
			//   this.$emit('submitted', result.data)
			// }
		}
	},

	computed: {
		...mapGetters(["config"])
	}
};
</script>
