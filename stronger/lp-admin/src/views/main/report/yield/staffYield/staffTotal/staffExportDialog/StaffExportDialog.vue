<template>
	<lp-dialog ref="dialog" title="导出" transition="dialog-bottom-transition">
		<div class="pt-6" slot="main">
			<v-row>
				<v-col v-for="(item, index) in cells.fields" :key="`updateExport_${index}`">
					<template v-if="item.inputType === 'date'">
						<z-date-picker
							:formId="formId"
							:formKey="item.formKey"
							:hideDetails="item.hideDetails"
							:hint="item.hint"
							:label="item.label"
							:options="item.options"
							:suffix="item.suffix"
							:defaultValue="item.defaultValue"
							range
						></z-date-picker>
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
</template>

<script>
import { mapGetters } from "vuex";
import DialogMixins from "@/mixins/DialogMixins";
import cells from "./cells";

export default {
	name: "StaffExportDialog",
	mixins: [DialogMixins],

	data() {
		return {
			formId: "staffExportDialog",
			cells
		};
	},
	created() {
		this.sockets.subscribe("outputStatistics", result => {
			this.toasted.success(`${result}，导出成功(files/output-statistics)`);
		});
	},
	methods: {
		async onConfirm() {
			const form = this.forms[this.formId];
			console.log(form);
			const result = await this.$store.dispatch("EXPORT_STAFF_TOTAL", form);

			this.toasted.info(`${result.msg}，正在导出...`);

			if (result.code === 200) {
				this.onClose();
			}
		}
	}
};
</script>
