<template>
	<lp-dialog ref="dialog" title="更新" transition="dialog-bottom-transition">
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
	name: "StaffUpdateDialog",
	mixins: [DialogMixins],

	data() {
		return {
			formId: "staffUpdateDialog",
			cells
		};
	},

	methods: {
		async onConfirm() {
			// console.log("11",this.rowInfo)
			// console.log(this.forms[this.formId]);
			// const param = {
			//   StartTime:this.rowInfo.date[0],
			//   EndTime:this.rowInfo.date[1],
			//   Code:  this.rowInfo.code,
			//   IsCheckAll:   this.rowInfo.type,
			//   IsCheckUpdate: "",
			//   UpdateTime:this.forms[this.formId].upDate,
			//   PageInfo  :"",
			// }

			const form = this.forms[this.formId];

			this.$emit("updateTime", { updateTime: form.updateTime });
			// const result = await this.$store.dispatch("UPDATE_OUTPUT_STATISTICS", param);
			// if (result.code === 200) {}
			this.onClose();
		}
	},
	computed: {
		...mapGetters(["pro"])
	}
};
</script>
