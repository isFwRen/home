import { R } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
import { mapRules, validations } from "./components/opTextField/rules";

const { baseURLApi } = lpTools.baseURL();

export default {
	data() {
		return {
			fileUrl: `${baseURLApi}files/`,

			// 返回重录(1：上一单，2：前一单，以此类推...)
			prevNums: 0
		};
	},

	mounted() {
		window.addEventListener("keydown", this.fuckOpShortcut);
	},

	beforeDestroy() {
		window.removeEventListener("keydown", this.fuckOpShortcut);
	},

	methods: {
		// 生成校验规则
		setValidateRules(field) {
			const rules = [];

			// 只有[合法field]才需生成校验规则
			if (field.show === true && field.disabled === false) {
				const result = R.find(this.fieldsConfig, { code: field.code });

				// key
				{
					for (let key in result) {
						if (mapRules[key] && R.isYummy(result[key])) {
							if (key === "fixValue") {
								if (result.fixValue) {
									const values = result.fixValue.split(";");
									rules.push({
										key: mapRules.fixValue,
										rule: values,
										message: `必须为${result.fixValue}中的一项.`
									});
								}
							} else if (key === "specChar") {
								if (result.specChar) {
									const values = result.specChar.split(";");
									rules.push({
										key: mapRules.specChar,
										rule: values,
										message: `只有${result.specChar}为可通过字符.`
									});
								}
							} else {
								rules.push({ key: mapRules[key], rule: result[key] });
							}
						}
					}
				}

				// validations
				if (R.isYummy(result.validations)) {
					for (let value of result.validations) {
						if (mapRules[value]) {
							if (mapRules[value] === "required") {
								rules.unshift({ key: mapRules[value], rule: "NO" });
							} else {
								rules.push({ key: mapRules[value], rule: "NO" });
							}
						}
					}
				}

				// Specific validations
				if (R.isYummy(field.includes)) {
					rules.push({
						key: mapRules.includes,
						rule: field.includes.items,
						message: field.includes.message
					});
				}

				// 问题件
				if (this.op === "opq") {
					rules.push({
						key: "included",
						rule: [field.op1Value, field.op2Value],
						message: "请检查录入数据，问题件与1、2码录入数据不一致!"
					});
				}
			}

			return rules;
		},

		// 校验当前字段
		validateField({ event, field, fieldsList, fieldsIndex }) {
			const normalResult = this.normalValidate({
				value: event.customValue,
				rules: field.rules
			});
			const specificResult = this.specificValidate({ event, field, fieldsList, fieldsIndex });

			if (normalResult && specificResult) {
				return true;
			}

			return false;
		},

		// 校验规则从当前字段的属性获取
		normalValidate({ value, rules }) {
			const normalRules = rules || [];

			if (R.isLousy(normalRules)) {
				return true;
			}

			const goodFieldRules = [];

			normalRules.map(rule => {
				if (rule.key !== "free_value") {
					goodFieldRules.push(rule);
				}
			});

			{
				const normalRule = R.find(normalRules, { key: "free_value" });
				const freeValues = normalRule?.rule;

				if (R.isYummy(freeValues)) {
					if (freeValues.includes(value)) {
						return true;
					}
				}
			}

			console.log({ goodFieldRules });

			return goodFieldRules?.every(
				item => validations[item.key]({ op: this.op, value, rule: item.rule }) === true
			);
		},

		// 校验规则从特殊校验规则获取
		specificValidate({ event, field, fieldsList, fieldsIndex }) {
			const { customValue: value, customItems: items } = event;
			const { op0, op1op2opq } = this.specificProject;
			const rules = [];

			if (this.op === "op0") {
				op0.rules?.map(rule => {
					if (rule.fields.includes(field.code)) {
						rules.push(rule);
					}
				});
			} else {
				op1op2opq.rules?.map(rule => {
					if (rule.fields.includes(field.code)) {
						rules.push(rule);
					}
				});
			}

			if (R.isLousy(rules)) {
				return true;
			}

			for (let rule of rules) {
				for (let funcName in rule) {
					if (R.getType(rule[funcName]) === "function") {
						const result = rule[funcName]({
							op: this.op,
							value,
							items,
							field,
							fieldsList,
							fieldsIndex
						});

						if (result !== true) {
							return false;
						}
					}
				}
			}

			return true;
		},

		fuckOpShortcut(event) {
			const { keyCode } = event || window.event;

			switch (keyCode) {
				// 返回修改(F3)
				case 114:
					event.preventDefault();
					const prevNums = ++this.prevNums;

					this.getTask({
						status: "modify",
						prevNums
					});
					break;

				// 提交(F8)
				case 119:
					event.preventDefault();

					if (this.block.fEight) {
						this.fEightSubmit();
					}

					console.warn({ fEight: this.block.fEight });
					break;
			}
		}
	}
};
