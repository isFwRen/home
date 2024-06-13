<template>
	<div class="dialog-wrap">
		<v-dialog
			v-model="dialog"
			:light="true"
			:persistent="true"
			:hide-overlay="true"
			width="500px"
		>
			<v-card>
				<v-card-title class="justify-center text-h5"> 更新备注 </v-card-title>
				<v-col cols="12">
					<v-textarea
						v-model.trim="remark"
						@keyup.ctrl.enter="comFirm"
						label="请输入内容"
					></v-textarea>
				</v-col>
				<v-card-actions>
					<v-spacer></v-spacer>
					<v-btn color="blue darken-1" text @click="closeOut"> 关闭 </v-btn>
					<v-btn color="blue darken-1" text @click="comFirm"> 确认 </v-btn>
				</v-card-actions>
			</v-card>
		</v-dialog>
	</div>
</template>

<script>
export default {
	name: "RemarkDialog",
	data() {
		return {
			dialog: false,
			remark: ""
		};
	},
	props: {
		value: {
			type: [String],
			default: ""
		}
	},
	watch: {
		value(val) {
			this.remark = val;
		}
	},
	methods: {
		comFirm() {
			if (this.value === this.remark) {
				this.onClose();
			} else {
				this.$emit("remarksEmit", this.remark);
				this.clear("confirm");
				this.onClose();
			}
		},
		closeOut() {
			this.clear("close");
			this.onClose();
		},
		onOpen() {
			this.dialog = true;
		},
		onClose() {
			this.dialog = false;
		},
		clear(remarks) {
			this.remark = "";
			this.$emit("remarksClose", remarks);
		}
	}
};
</script>

<style scoped lang="scss">
.text-h5 {
	text-align: center;
}
</style>
