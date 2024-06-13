<template>
	<lp-dialog
		ref="dialog"
		:title="`新增${firstInfo.label}`"
		transition="dialog-bottom-transition"
		width="500"
		@dialog="handleDialog"
	>
		<div class="pt-6" slot="main">
			<z-textarea
				:formId="firstInfo.formId"
				:formKey="firstInfo.formKey"
				:placeholder="firstInfo.placeholder"
			></z-textarea>

			<div class="z-flex justify-end btns">
				<z-btn class="mr-3" color="primary" small @click="onSubmit(values1, values2)">
					提交
				</z-btn>
			</div>
		</div>
	</lp-dialog>
</template>

<script>
import { mapState } from "vuex";
import { rocket } from "vue-rocket";
import TableMixins from "@/mixins/TableMixins";
import DialogMixins from "@/mixins/DialogMixins";
import cells from "./cells";

export default {
	name: "AdddDialog",
	mixins: [TableMixins, DialogMixins],
	inject: ["task"],

	props: {
		firstInfo: {
			type: Object,
			default: () => {}
		},
		taskInfos: {
			type: String | Object | Array
		}
	},

	data() {
		return {
			cells
		};
	},

	methods: {
		async onSubmit(values1, values2) {
			const { formKey } = values2;

			// const form = this.forms[formId];
			const dispatchName =
				formKey === "organizationNumber"
					? "SET_PRIORITY_ORGANIZATION_NUMBER_ITEM_LIST"
					: "TASK_SET_BILL_STICK_LEVEL";
			const body = {
				proCode: this.task.proCode,
				list: values1.split("\n"),
				stickLevel: formKey === "stickLevel1" ? 1 : 2
			};

			const result = await this.$store.dispatch(dispatchName, body);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.onClose();
				this.$parent.onRefresh();
			}
		},

		handleDialog(dialog) {
			if (!dialog) {
				rocket.emit("ZHT_CLEAR_FORM", this.firstInfo.formId);
			}
		}
	},

	computed: {
		...mapState(["forms"])
	}
};
</script>
