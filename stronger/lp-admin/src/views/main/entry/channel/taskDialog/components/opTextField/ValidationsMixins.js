import { R } from "vue-rocket";
import { validations } from "./rules";

export default {
	props: {
		bill: {
			type: Object,
			default: () => ({})
		},

		block: {
			type: Object,
			default: () => ({})
		},

		svHints: {
			type: Array,
			required: false
		},

		includes: {
			type: Array,
			default: () => []
		},

		svValidations: {
			type: Array,
			required: false
		},

		validations: {
			type: Array,
			required: false
		}
	},

	data() {
		return {
			error: false,
			errorMessages: [],
			rules: [],
			sessionRules: [],
			svHintList: [],

			billObject: {},
			blockObject: {}
		};
	},

	created() {
		this.billObject = {
			billNum: this.bill.billNum
		};

		this.blockObject = {
			code: this.block.code
		};

		this.setNormalValidations();
		this.setSpecificValidations();
		this.setSpecificHints();

		if (this.field.disabled || !this.field.show) {
			this.rules = [];
		}
	},

	methods: {
		// 设置普通校验
		setNormalValidations() {
			if (R.isYummy(this.validations)) {
				this.rules = [];

				this.showLength = R.find(this.validations, { key: "show_length" }) ? true : false;
				this.allowSpace = R.find(this.validations, { key: "spaced" }) ? true : false;

				for (let validation of this.validations) {
					if (R.isYummy(validation["rule"])) {
						if (validation["key"] !== "free_value") {
							this.rules.push(() =>
								validations[validation["key"]]({
									label: this.label,
									op: this.op,
									value: this.value,
									rule: validation["rule"],
									message: validation["message"],
									field: this.field,
									fieldsList: this.fieldsList
								})
							);
						}
					}
				}
			}
		},

		// 设置特殊校验
		setSpecificValidations() {
			const { code } = this.field;

			if (R.isYummy(this.svValidations)) {
				for (let validation of this.svValidations) {
					if (validation["fields"]?.includes(code)) {
						for (let funcName in validation) {
							if (R.getType(validation[funcName]) === "function") {
								this.rules.push(() =>
									validation[funcName]({
										bill: this.billObject,
										block: this.blockObject,
										field: this.field,
										fieldsIndex: this.fieldsIndex,
										fieldsList: this.fieldsList,
										label: this.label,
										op: this.op,
										value: this.value,
										items: this.items
									})
								);
							}
						}
					}
				}
			}

			this.sessionRules = R.deepClone(this.rules);
		},

		// 设置特殊提示文本
		setSpecificHints() {
			const { code } = this.field;
			this.svHintList = [];

			if (R.isYummy(this.svHints)) {
				for (let hint of this.svHints) {
					if (hint["fields"]?.includes(code)) {
						for (let funcName in hint) {
							if (R.getType(hint[funcName]) === "function") {
								const message = hint[funcName]({
									bill: this.billObject,
									block: this.blockObject,
									field: this.field,
									fieldsIndex: this.fieldsIndex,
									fieldsList: this.fieldsList,
									label: this.label,
									op: this.op,
									value: this.value
								});

								if (R.getType(message) === "string") {
									this.svHintList.push(message);
								}
							}
						}
					}
				}
			}
		}
	},

	watch: {
		value: {
			handler(value) {
				this.$nextTick(() => {
					const specChar = this.field.specChar;

					if (R.isYummy(specChar)) {
						const includes = specChar.split(";");

						if (includes.includes(value)) {
							this.rules = [];
						} else {
							this.rules = R.deepClone(this.sessionRules);
						}
					}

					this.rules = [...this.rules];
				});
			}
		},

		// 若当前 field 禁用，则去掉校验规则
		fieldsList: {
			handler(fieldsList) {
				fieldsList[this.fieldsIndex]?.map(field => {
					if (field.code === this.field.code) {
						if (this.field.disabled) {
							this.rules = [];
						}
					}
				});
			},
			deep: true,
			immediate: true
		}
	}
};
