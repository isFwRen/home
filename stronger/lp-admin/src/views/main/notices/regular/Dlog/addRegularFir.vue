<template>
	<div class="add-dialog">
		<lp-dialog ref="dialog" :title="title" width="900" height="600" @dialog="handleDialog">
			<div class="pt-6" slot="main">
				<div>
					<v-row class="z-flex align-center" v-for="e in 7">
						<v-col cols="1"
							><span>周{{ numberArr[e - 1] }}</span>
						</v-col>
						<v-col cols="3">
							<z-date-picker
								:formId="formId"
								:formKey="'0startTime' + e"
								label="起始时间"
								format="24hr"
								mode="time"
								prepend-icon="mdi-alarm"
								time-use-seconds
								time-format="24hr"
								defaultValue=""
								:validation="
									rule
										? [
												{
													rule: 'required',
													message: '字段不允许为空.'
												}
										  ]
										: []
								"
							></z-date-picker>
						</v-col>
						<v-col cols="3">
							<z-date-picker
								:formId="formId"
								:formKey="'1endTime' + e"
								label="结束时间"
								format="24hr"
								mode="time"
								prepend-icon="mdi-alarm"
								time-use-seconds
								time-format="24hr"
								defaultValue=""
								:validation="
									rule
										? [
												{
													rule: 'required',
													message: '字段不允许为空.'
												}
										  ]
										: []
								"
							></z-date-picker>
						</v-col>
						<v-col cols="3">
							<z-text-field
								:formId="formId"
								:formKey="'2interval' + e"
								label="时间间隔(min)"
								:validation="[
									{
										regex: /^[1-9]([0-9]{0,2})$/,
										message: '必须是大于0的数字且不能大于3位'
									}
								]"
							>
							</z-text-field>
						</v-col>
					</v-row>
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

export default {
	name: "NoticeDialog",
	mixins: [DialogMixins],
	data() {
		return {
			formId: "setSendTime",
			dispatchForm: "ADD_NOTIC_NEW_REGULAR_BYTIME",
			title: "新增",
			ruleList: {},
			numberArr: ["一", "二", "三", "四", "五", "六", "日"],
			rule: false
		};
	},
	props: ["id"],
	methods: {
		onConfirm() {
			const form = this.forms[this.formId];

			const objKeyArr = ["startTime", "endTime", "interval"];
			let oneObj = {};
			for (let key in form) {
				let keyArr = key.split("");
				if (!oneObj[keyArr[keyArr.length - 1]]) {
					oneObj[keyArr[keyArr.length - 1]] = {};
				}
				oneObj[keyArr[keyArr.length - 1]][objKeyArr[keyArr[0]]] = form[key];
			}
			let ones = [];
			for (let key in oneObj) {
				ones.push({
					...oneObj[key],
					dayOfWeek: +key,
					groupId: this.id,
					interval: +oneObj[key]["interval"]
				});
			}
			ones = ones.filter(e => {
				return e.startTime && e.endTime && e.interval;
			});
			const body = {
				ones,
				type: 1
			};
			this.submit(body);
		},
		timeBeclick(e) {
			const form = this.forms[this.formId];
			if (form["startTime" + e] || form["endTime" + e] || form["interval" + e]) {
				this.ruleList[e] = [{ rule: "required", message: "字段不允许为空." }];
			} else {
				this.ruleList[e] = [];
			}
			console.log(this.ruleList);
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
