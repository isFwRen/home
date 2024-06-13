import { mapGetters } from "vuex";
const selectProCodes = [
	"identifyStatistics",
	"institutionalExtraction",
	"destructionReport",
	"DirectoryOut",
	"AbnormalPart"
];

export default {
	data() {
		return {
			manual: true
		};
	},

	created() {
		if (this.$options.name !== "ReturnAnalysis") {
			this.setProOptions();
		} else {
			this.setReturnAnalysisProOptions();
		}
	},

	beforeDestroy() {
		this.$store.commit("SET_PROJECT_INFO");
	},

	methods: {
		// 设置项目下拉
		async setProOptions() {
			if (selectProCodes.indexOf(this.$options.name) !== -1) {
				this.cells.fields[0].options = [...this.auth.proItems];
			} else {
				this.cells.fields[0].options = [
					{ label: "全部", value: "" },
					...this.auth.proItems
				];
			}
		},

		// 设置回传分析项目下拉
		async setReturnAnalysisProOptions() {
			const items = [
				{ label: "全部/整体", value: "全部/整体" },
				{ label: "全部/明细", value: "全部/明细" }
			];

			for (let item of this.auth.proItems) {
				items.push(
					{ label: `${item.value}/整体`, value: `${item.value}/整体` },
					{ label: `${item.value}/明细`, value: `${item.value}/明细` }
				);
			}

			this.cells.fields[0].options = items;
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
