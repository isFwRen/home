<template>
	<div class="update-item">
		<lp-dialog ref="dialog" :title="title" width="500" @dialog="handleDialog">
			<div class="pt-6" slot="main">
				<v-row>
					<v-col
						v-for="(item, index) in cells.fields"
						:key="`updateItem_${index}`"
						class="py-0"
						:cols="item.cols"
					>
						<template v-if="item.inputType === 'input'">
							<z-text-field
								:formId="formId"
								:formKey="item.formKey"
								:disabled="item.disabled"
								:hideDetails="item.hideDetails"
								:hint="item.hint"
								:label="item.label"
								:suffix="item.suffix"
								:validation="item.validation"
								:defaultValue="detail[item.formKey]"
							>
							</z-text-field>
						</template>

						<template v-else>
							<z-btn-toggle
								:formId="formId"
								:formKey="item.formKey"
								color="primary"
								:disabled="item.disabled"
								mandatory
								:options="item.options"
								:defaultValue="detail[item.formKey]"
							></z-btn-toggle>
						</template>
					</v-col>
				</v-row>
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
import cells from "./cells";

export default {
	name: "UpdateItemDialog",
	mixins: [DialogMixins],

	data() {
		return {
			formId: "UpdateItemDialog",
			dispatchForm: "UPDATE_CONFIG_ITEM",
			cells
		};
	},

	methods: {
		onConfirm() {
			const form = this.forms[this.formId];

			const body = {
				id: this.rowInfo.ID,
				editVersion: this.rowInfo.editVersion,
				...form
			};

			this.submit(body);
		}
	},

	watch: {
		dialog: {
			handler(dialog) {
				if (dialog) {
					if (this.status === -1) {
						cells.fields[0].disabled = false;
						cells.fields[1].disabled = false;
					} else {
						cells.fields[0].disabled = true;
						cells.fields[1].disabled = true;
					}
				}
			},
			immediate: true
		}
	}
};
</script>
