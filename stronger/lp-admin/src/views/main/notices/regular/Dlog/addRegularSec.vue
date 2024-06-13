<template>
	<div class="add-dialog">
		<lp-dialog ref="dialog" title="新增" width="500" height="600" @dialog="handleDialog">
			<div class="pt-6" slot="main">
				<div>
					<v-row class="z-flex align-center" v-for="e in 7">
						<v-col cols="2"
							><span>周{{ e }}</span>
						</v-col>
						<v-col>
							<z-date-picker
								:formId="formId"
								:formKey="e.toString()"
								label="固定时间"
								format="24hr"
								:immediate="false"
								mode="time"
								prepend-icon="mdi-alarm"
								time-use-seconds
								time-format="24hr"
								defaultValue=""
							></z-date-picker>
						</v-col>
					</v-row>
				</div>
			</div>

			<div class="z-flex" slot="actions">
				<z-btn class="mr-3" color="normal" @click="onClose">取消</z-btn>

				<z-btn :formId="formId" color="primary" @click="onConfirm">确认</z-btn>
			</div>
		</lp-dialog>
	</div>
</template>

<script>
import DialogMixins from "@/mixins/DialogMixins";

export default {
	name: "NoticeREDialog",
	mixins: [DialogMixins],
	data() {
		return {
			formId: "NoticeUpdateSecDialog",
			dispatchForm: "ADD_NOTIC_NEW_REGULAR_BYTIME",
			title: "12"
		};
	},
	props: ["id", "proCode"],
	methods: {
		onConfirm() {
			const form = this.forms[this.formId];
			const twos = form
				.map((e, i) => {
					return {
						dayOfWeek: i % 7,
						groupId: this.id,
						sendTime: e,
						proCode: this.proCode
					};
				})
				.filter(e => {
					return e && e.sendTime;
				});
			const body = {
				twos,
				type: 2
			};
			this.submit(body);
		}
	}
};
</script>

<style scoped>
::v-deep .v-text-field,
::v-deep .v-input__slot {
	padding: 0;
	margin: 0;
}
</style>
