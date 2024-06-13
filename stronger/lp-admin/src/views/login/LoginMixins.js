import { mapState } from "vuex";
import { tools } from "@/libs/util";
import { regex_phone } from "./cells";
import cells from "./cells";

const isIntranet = tools.isIntranet();

export default {
	data() {
		return {
			isIntranet: true,
			userField: isIntranet ? cells.usernameField : cells.phoneField,
			validAccount: false,
			captchaId: undefined
		};
	},

	created() {
		this.judgeIntranet();
	},

	methods: {
		// 获取验证码
		async sendCode() {
			const accountKey = this.isIntranet ? "username" : "phone";

			const form = {
				// isIntranet: this.isIntranet,
				phone: this.forms[this.formId].account
			};

			this.clock = setInterval(this.countdown, 1000);

			const result = await this.$store.dispatch("GET_DING_CODE", form);
			this.toasted.dynamic(result.msg, result.code);

			// if (result.code === 200) {
			// 	this.captchaId = result.data.captchaId;
			// }
		},

		onConfirm() {
			const accountKey = this.isIntranet ? "username" : "phone";

			const form = {
				isIntranet: this.isIntranet,
				accountKey,
				captchaId: this.captchaId,
				...this.forms[this.formId]
			};

			this.submit(form);
		},

		// 判断内/外网
		judgeIntranet() {
			const { usernameField, phoneField } = cells;

			this.isIntranet = tools.isIntranet();

			this.userField = this.isIntranet ? usernameField : phoneField;
		}
	},

	computed: {
		...mapState(["forms"])
	},

	watch: {
		forms: {
			handler() {
				if (this.forms[this.formId]) {
					const { account } = this.forms[this.formId];

					this.validAccount = this.isIntranet
						? !!account
						: account && regex_phone.test(account);
				}
			},
			immediate: true
		}
	}
};
