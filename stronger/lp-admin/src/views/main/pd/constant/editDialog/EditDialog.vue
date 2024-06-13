<template>
	<div class="edit-dialog">
		<lp-dialog ref="dialog" :title="title" :width="600" @dialog="handleDialog">
			<div slot="main">
				<ul class="mt-4">
					<li v-for="(item, index) in fields" :key="index" class="mb-4">
						<z-text-field
							:formId="formId"
							:formKey="item.formKey"
							:label="item.label"
							:defaultValue="rowInfo[item.label]"
						>
						</z-text-field>
					</li>
				</ul>
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
import { R } from "vue-rocket";

export default {
	name: "EditDialog",
	mixins: [DialogMixins],

	props: {
		headers: {
			type: Array,
			required: true
		}
	},

	data() {
		return {
			formId: "EditDialog",

			fields: []
		};
	},

	methods: {
		onConfirm() {
			const form = this.forms[this.formId];

			const item = {};

			Object.keys(form).forEach(key => {
				const f = this.fields.find(h => h.formKey === key);
				if (f) {
					item[f.label] = form[key];
				}
			});

			console.log(this.rowInfo,'rowInfo')

			const info = {
				id: this.rowInfo._id,
				// _X_ROW_KEY: this.rowInfo._X_ROW_KEY,
				item
			};

			this.$emit("submitted", info);
		}
	},

	watch: {
		dialog: {
			handler(dialog) {
				if (dialog) {
					const headers = [...this.headers];

					this.fields = [];

					for (let index in headers) {
						this.fields = [
							...this.fields,
							{
								formKey: `name${index}`,
								inputType: "text",
								hideDetails: false,
								label: headers[index].text
							}
						];
					}
				}
			},
			immediate: true
		}
	}
};
</script>
