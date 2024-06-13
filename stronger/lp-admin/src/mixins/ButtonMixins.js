export default {
	data() {
		return {
			color: "primary",

			disabled: false,
			loading: false,

			pending: false,

			text: "获取验证码",
			clock: null,
			interval: 60,
			counting: false
		};
	},

	methods: {
		isPending(pending) {
			this.pending = pending;
		},

		countdown() {
			this.interval--;
			if (this.interval >= 0) {
				this.counting = true;
				this.text = `${this.interval}秒后重新获取`;
			} else {
				clearInterval(this.clock);
				this.counting = false;
				this.interval = 60;
				this.text = "获取验证码";
			}
		}
	}
};
