import { mapState } from "vuex";
import { rocket } from "vue-rocket";
import { mapActions } from "vuex";

const titles = new Map([
	[-1, "新增"],
	[0, "详情"],
	[1, "编辑"]
]);

export default {
	props: {
		rowInfo: {
			type: Object,
			default: () => {}
		}
	},

	data() {
		return {
			closeDelay: 10000,
			dialog: false,
			status: null /* -1：新增，0：只读，1：编辑 */,
			title: null,
			detail: {}
		};
	},

	methods: {
		async submit(form) {
			form.status = this.status;

			const result = await this.$store.dispatch(this.dispatchForm, form);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.onClose();
				this.$emit("submitted", result);
			}

			return result;
		},

		setDetail() {
			if (!this.rowInfo) return;

			delete this.rowInfo._XID;

			this.detail = this.rowInfo;
		},

		onBack() {
			this.onClose();
		},

		onClose() {
			this.$refs.dialog.onClose();
		},

		...mapActions("configField", ["ADD_ISSUE_CONFIG"]),

		async onConfirm(data) {
			this.$refs.dialog.onClose();
			const result = await this.$store.dispatch("ADD_ISSUE_CONFIG", data);

			this.toasted.dynamic(result.msg);
		},

		onOpen(status) {
			if (typeof status === "object") {
				this.status = status.status;
				this.title = status.title;
			} else {
				this.status = status;
				this.title = titles.get(this.status) || status;
			}
			this.$refs.dialog.onOpen();
		},

		handleDialog(dialog) {
			this.dialog = dialog;

			if (dialog) {
				this.detail = {};
				this.$nextTick(() => {
					this.setDetail();
				});
			} else {
				rocket.emit("ZHT_RESET_FORM", this.formId);
			}
		}
	},

	computed: {
		...mapState(["forms"])
	}
};
