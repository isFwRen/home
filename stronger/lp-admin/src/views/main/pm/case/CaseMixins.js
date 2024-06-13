import storage from "@/libs/util.storage";

const defaultCaseInfo = {
	createAt: "",
	caseId: "",
	billNum: "",
	proCode: "",
	billInfo: {}
};

export default {
	data() {
		return {};
	},

	methods: {
		onBack() {
			this.$emit("back");
		},

		onClose() {
			this.$emit("close");
		},

		rememberCaseInfo(caseInfo) {
			let storageCaseInfo = storage.get("caseInfo") || defaultCaseInfo;

			storageCaseInfo = Object.assign(storageCaseInfo, caseInfo);

			storage.set("caseInfo", storageCaseInfo);
		}
	}
};
