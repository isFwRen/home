<template>
	<lp-dialog ref="dialog" title="添加新分块" width="500" @dialog="handleDialog">
		<div class="pt-4" slot="main">
			<ul>
				<li>
					<z-text-field
						:formId="formId"
						formKey="itemName"
						label="项目名称"
						disabled
						:defaultValue="config.pro.name"
					>
					</z-text-field>
				</li>

				<li>
					<z-text-field
						:formId="formId"
						formKey="_code"
						label="分块编码"
						disabled
						:defaultValue="code"
					>
					</z-text-field>
				</li>

				<li>
					<z-text-field
						:formId="formId"
						formKey="name"
						label="分块名称"
						:validation="[{ rule: 'required', message: '分块名称不能为空.' }]"
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
</template>

<script>
import { mapGetters } from "vuex";
import DialogMixins from "@/mixins/DialogMixins";
import cells from "../cells";

export default {
	name: "NewChunkDialog",
	mixins: [DialogMixins],

	props: {
		tempItem: {
			type: Object,
			default: () => {}
		}
	},

	data() {
		return {
			formId: "NewChunkDialog",
			dispatchForm: "UPDATE_CONFIG_TEMP_CHUNK",
			cells,
			code: "1"
		};
	},

	methods: {
		onConfirm() {
			const { label, value } = this.tempItem;
			const freeTime = label === cells.UNDEFINED ? 3600 : 1800;
			const form = {
				tempId: value,
				freeTime,
				myOrder: this.detail.myOrder,
				...this.forms[this.formId]
			};

			this.submit(form);
		}
	},

	computed: {
		...mapGetters(["config"])
	},

	watch: {
		detail: {
			handler(detail) {
				const { myOrder } = detail;
				if (myOrder) {
					let num = +myOrder;
					if (num < 10) {
						num = "00" + num;
					} else if (num < 100) {
						num = "0" + num;
					}

					this.code = "bc" + num;
				}
			},
			immediate: true
		}
	}
};
</script>
