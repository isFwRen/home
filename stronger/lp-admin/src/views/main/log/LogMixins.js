import { mapGetters } from "vuex";

export default {
	data() {
		return {
			manual: true
		};
	},

	created() {
		this.setProOptions();
	},

	beforeDestroy() {
		this.$store.commit("SET_PROJECT_INFO");
	},

	methods: {
		// 设置项目下拉
		async setProOptions() {
			const proCodeOptions = [{ label: "全部", value: "" }];

			for (let item of this.auth.perm) {
				if (item.hasPm) {
					proCodeOptions.push({
						label: item.proCode,
						value: item.proCode
					});
				}
			}

			this.cells.fields[1].options = proCodeOptions;
		},

		handleChange(value, item) {
			const { formKey } = item;

			if (formKey === "proCode") {
				this.$store.commit("SET_PROJECT_INFO", { code: value });
			}
		}
	},

	computed: {
		...mapGetters(["auth"])
	}
};
