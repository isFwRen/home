import { tools } from "vue-rocket";
import _ from "lodash";

import B0118 from "./specificValidations/B0118/index";

export default {
	data() {
		return {
			specificProject: { ...B0118 },

			// 记录最后一次存储的合法field
			svMemoFields: {},

			// fields 的值从 targets 里的值选择
			svDropdownFields: {},

			svConstantsDD: {}
		};
	},

	methods: {
		// 重置变量
		svResetVariable() {
			this.svMemoFields = {};
			this.svDropdownFields = {};
			this.svConstantsDD = {};
		},

		// 记录最后一次存储的合法field
		svUpdateMemoFields({ value, code, items }) {
			const { op0, op1op2opq } = this.specificProject;
			const memoFields = [];

			if (op0.memoFields) {
				memoFields.push(...op0.memoFields);
			}

			if (op1op2opq.memoFields) {
				memoFields.push(...op1op2opq.memoFields);
			}

			// 记录
			if (memoFields.includes(code)) {
				_.set(this.svMemoFields, `${code}.value`, value);
				_.set(this.svMemoFields, `${code}.items`, items);
			}
		},

		// 以某个字段的所有值为下拉
		svUpdateDropdown() {
			const { op0, op1op2opq } = this.specificProject;
			const dropdownFields = [];

			if (op0.dropdownFields) {
				dropdownFields.push(...op0.dropdownFields);
			}

			if (op1op2opq.dropdownFields) {
				dropdownFields.push(...op1op2opq.dropdownFields);
			}

			for (let record of dropdownFields) {
				const [fields, items, flatFieldsList, flatInitFieldsList] = [[], [], [], []];

				for (let key in this.fieldsObject) {
					const sessionStorage = this.fieldsObject[key].sessionStorage;
					flatFieldsList.push(...tools.flatArray(this.fieldsObject[key].fieldsList));
					flatInitFieldsList.push(
						...tools.flatArray(this.fieldsObject[key].initFieldsList)
					);

					// 临时存储或为当前页面
					if (sessionStorage || this.thumbIndex === +key) {
						flatFieldsList.map(field => {
							if (record.targets.includes(field.code)) {
								fields.push(tools.deepClone(field));
							}
						});
					} else {
						flatInitFieldsList.map(field => {
							if (record.targets.includes(field.code)) {
								fields.push(tools.deepClone(field));
							}
						});
					}
				}

				// 找到需要将值设置为 dropdown 的 field
				fields.map(field => {
					if (record.targets.includes(field.code)) {
						if (!items.includes(field.resultValue) && field.resultValue) {
							items.push(field.resultValue);
						}
					}
				});

				// 给符合条件的 field 设置 dropdown
				record.fields.map(fieldCode => {
					_.set(this.svConstantsDD, `${fieldCode}.desserts`, items);
				});
			}
		},

		// 字段已生成
		svInit({ fieldsList, codeValues }) {
			const { op0, op1op2opq } = this.specificProject;
			let methods = {};

			if (op0.init?.methods) {
				methods = { ...methods, ...op0.init.methods };
			}

			if (op1op2opq.init?.methods) {
				methods = { ...methods, ...op1op2opq.init.methods };
			}

			// 执行函数
			if (tools.isYummy(methods)) {
				for (let funcName in methods) {
					methods[funcName]({
						op: this.op,
						bill: this.bill,
						focusFieldsIndex: this.focusFieldsIndex,
						fieldsList,
						flatFieldsList: tools.flatArray(fieldsList),
						codeValues
					});
				}
			}
		},

		// 回车后操作字段
		svEnterUpdateField({ field, fieldsList, focusFieldsIndex, memoFields }) {
			const { op0, op1op2opq } = this.specificProject;
			let methods = {};

			if (op0.enter?.methods) {
				methods = { ...methods, ...op0.enter.methods };
			}

			if (op1op2opq.enter?.methods) {
				methods = { ...methods, ...op1op2opq.enter.methods };
			}

			// 执行函数
			if (tools.isYummy(methods)) {
				for (let funcName in methods) {
					methods[funcName]({
						op: this.op,
						field,
						fieldsList,
						focusFieldsIndex,
						memoFields,
						flatFieldsList: tools.flatArray(fieldsList)
					});
				}
			}
		},

		// 按F4后禁用字段
		svDisableFields({ fieldsList, focusFieldsIndex }) {
			const flatFieldsList = tools.flatArray(fieldsList);

			const { op0, op1op2opq } = this.specificProject;
			let methods = {};

			if (op0.sessionSave?.methods) {
				methods = { ...methods, ...op0.sessionSave.methods };
			}

			if (op1op2opq.sessionSave?.methods) {
				methods = { ...methods, ...op1op2opq.sessionSave.methods };
			}

			// 执行禁用函数
			if (tools.isYummy(methods)) {
				for (let funcName in methods) {
					methods[funcName]({ fieldsList, focusFieldsIndex, flatFieldsList });
				}
			}
		},

		/**
		 * @description 提交前校验
		 * @param fieldsList
		 */
		svValidateFields({ fieldsList }) {
			const { op0, op1op2opq } = this.specificProject;
			let methods = {};
			const errors = [];

			if (op0.beforeSubmit.methods) {
				methods = { ...methods, ...op0.beforeSubmit.methods };
			}

			if (op1op2opq.beforeSubmit.methods) {
				methods = { ...methods, ...op1op2opq.beforeSubmit.methods };
			}

			// 执行校验函数
			if (tools.isYummy(methods)) {
				for (let funcName in methods) {
					const result = methods[funcName]({
						bill: this.bill,
						block: this.block,
						fieldsList,
						flatFieldsList: tools.flatArray(fieldsList),
						op: this.op
					});

					if (result !== true) {
						errors.push(result);
					}
				}
			}

			if (!errors.length) {
				return true;
			} else {
				for (let error of errors) {
					alert(error.errorMessage);
				}

				return false;
			}
		},

		// 从常量库查询匹配字段
		async svSearchConstants({ value, field }) {
			let items = [];

			// 需要查询常量
			if (field.table) {
				const { name: tableName, query, targets } = field.table;
				const tableInfo = window.constantsDB?.[tableName];

				// 当前常量存在
				if (tableInfo) {
					this.svDropdownFields[field.uniqueId] = {
						items: []
					};

					const targetIndexes = [];

					if (targets) {
						tableInfo.headers.map((header, headerIndex) => {
							if (targets.includes(header)) {
								targetIndexes.push(headerIndex);
							}
						});
					}

					// 获取过滤名称下标
					const queryIndex = tableInfo.headers.indexOf(query);

					if (value) {
						if (!targets) {
							for (let dessert of tableInfo.desserts) {
								const text = dessert?.[queryIndex] || "";

								if (text.includes(value) && !items.includes(text)) {
									items.push(text);
								}
							}
						} else {
							for (let dessert of tableInfo.desserts) {
								let texts = "";

								dessert.map((text, textIndex) => {
									if (targetIndexes.includes(textIndex)) {
										const lastTargetIndex =
											targetIndexes[targetIndexes.length - 1];

										texts += textIndex !== lastTargetIndex ? `${text}-` : text;
									}
								});

								if (texts.includes(value) && !items.includes(texts)) {
									items.push(texts);
								}
							}
						}
					} else {
						items = [];
					}

					field.items = items.slice(0, 30);
					this.svDropdownFields[field.uniqueId].items = items.slice(0, 30);
				}
			}

			// 当前字段的值从某个字段的所有值中选择
			if (this.svConstantsDD?.[field.code]) {
				const DD = this.svConstantsDD?.[field.code];

				this.svDropdownFields[field.uniqueId] = {
					items: []
				};

				if (value) {
					for (let text of DD.desserts) {
						if (text.includes(value) && !items.includes(text)) {
							items.push(text);
						}
					}
				} else {
					items = [];
				}

				field.items = items.slice(0, 30);
				this.svDropdownFields[field.uniqueId].items = items.slice(0, 30);
			}

			this.svDropdownFields = { ...this.svDropdownFields };
		}
	}
};
