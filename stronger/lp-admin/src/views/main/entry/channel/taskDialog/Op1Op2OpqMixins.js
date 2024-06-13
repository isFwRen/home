import { mapGetters } from "vuex";
import { R } from "vue-rocket";

const delKeys = ["CreatedAt", "UpdatedAt"];

// [合法field]: `${ this.op }Input` !== ('no' | 'no_if') 的 field

const defaultBill = {
	pictures: [],
	phoTotal: 0,
	pageIndex: 0
};

export default {
	data() {
		return {
			hasOp: false,
			valid: true,

			// config
			mb001Config: [],
			fieldsConfig: [],

			// bill
			bill: defaultBill,

			// block
			block: {},
			isLoop: false,

			// codeValues
			codeValues: {},

			// fieldsList
			fieldsListLength: 0,

			// fields
			fieldsList: [],
			tempFields: [],
			focusFieldsIndex: 0,
			focusFieldIndex: 0,

			// field
			prevGoodField: null,
			nextGoodField: null,

			oldTime: null
		};
	},

	created() {
		this.setConfigs();
		this.getTask({ status: "new" });
	},

	mounted() {
		window.addEventListener("keydown", this.fuckEvents);
	},

	beforeDestroy() {
		window.removeEventListener("keydown", this.fuckEvents);
	},

	methods: {
		// 配置信息
		setConfigs() {
			const { mb001, fields } = this.task.config;
			this.mb001Config = mb001;
			this.fieldsConfig = fields;
		},

		// 领任务
		async getTask({ status, prevNums }) {
			const user = this.storage.get("user");

			const form = {
				code: user.code,
				op: this.op
			};

			let result = null;

			switch (status) {
				case "new":
					delete form.prevNums;
					result = await this.$store.dispatch("GET_LP_TASK", form);
					break;

				case "modify":
					form.prevNums = prevNums;
					result = await this.$store.dispatch("GET_LP_TASK_MODIFY", form);
					break;
			}

			const {
				bill = defaultBill,
				block = {},
				codeValues = {},
				fields = []
			} = R.isYummy(result.data) ? result.data : {};

			if (result.code === 200) {
				this.setBill(bill);
				this.block = block;
				this.isLoop = this.block.isLoop;
				this.codeValues = codeValues;
				this.setFieldsList(fields);
				this.hasOp = true;
			} else {
				this.toasted.dynamic(result.msg, result.code, { position: "top" });
				this.hasOp = false;
			}

			this.$emit("gotTaskResponse", {
				code: result.code,
				bill: this.bill,
				block: this.block
			});

			this.svResetVariable();
			this.scrollToTop();
		},

		// bill
		setBill(bill) {
			if (R.isLousy(bill.pictures)) {
				this.toasted.warning("未获取到图片！", { position: "top" });
				bill.pictures = [];
			}

			bill.phoTotal = bill.pictures.length;
			bill.pageIndex = 0;

			this.bill = bill;
		},

		// fieldsList
		setFieldsList(fieldsList) {
			if (R.isYummy(fieldsList)) {
				this.fieldsList = R.deepClone(fieldsList);
			} else {
				this.fieldsList = [];
				return;
			}

			// 设置 fieldsList 默认信息
			this.fieldsList.map((fields, fieldsIndex) => {
				fields.map((field, fieldIndex) => {
					// resultInput(op1 | op2 | opq)：op[1|2|q]Input 的值不为 no 才需要赋值
					if (field[`${this.op}Input`] !== "no") {
						field.resultInput = field[`${this.op}Input`];
					}

					this.setFieldEffectKeyValue(fieldsIndex, field, fieldIndex);

					// 设置默认值
					if (this.op === "opq") {
						this.opqSetFieldEffectKeyValue(field);
					}

					for (let key in field) {
						// 删除不需要的键值对
						if (delKeys.includes(key)) {
							delete field[key];
						}
					}

					// 根据 config 添加字段
					{
						const configField = R.find(this.fieldsConfig, field.code);

						field.specChar = configField.specChar;
						field.prompt = configField.prompt;
					}

					// 校验规则
					field.rules = this.setValidateRules(field);
				});
			});

			this.svInit({
				fieldsList: this.fieldsList,
				codeValues: this.codeValues
			});

			this.findAndMarkFirstLastGoodField();

			// 若 isLoop === true，以 fieldsList 的第 0 个为基础创建模板
			if (this.isLoop) {
				this.tempFields = R.deepClone(this.fieldsList[0]);

				this.tempFields.map(field => {
					delete field.ID;
				});
			}

			// console.log({fieldsList: this.fieldsList})
			// console.log({tempFields: this.tempFields})
		},

		// 前端自定义 keyValue
		setFieldEffectKeyValue(fieldsIndex, field, fieldIndex, isTempField = false) {
			// 唯一id
			field.uniqueId = `${fieldsIndex}_${fieldIndex}`;

			// 唯一key
			field.uniqueKey = field.uniqueId;

			// 当前 fields下的下标
			field.fieldIndex = fieldIndex;

			// 若为 tempFields 不需要往下执行
			if (isTempField) {
				return;
			}

			// 根据`${ this.op }Input`的值设置 field 状态，yes || '' 显示 field，no 隐藏 field，no_if 禁用 field
			{
				if (field[`${this.op}Input`] === "no") {
					field.show = false;
				} else {
					field.show = true;
				}

				if (field[`${this.op}Input`] === "no_if") {
					field.disabled = true;
				} else {
					field.disabled = false;
				}
			}
		},

		// 找到第一个及最后一个[合法field]，并标记
		findAndMarkFirstLastGoodField() {
			this.fieldsList.map(fields => {
				let [firstGoodField, lastGoodField] = [null, null];

				// 默认不设置 firstGoodField lastGoodField
				fields.map(field => {
					if (field.hasOwnProperty("isFirstGoodField")) delete field.isFirstGoodField;
					if (field.hasOwnProperty("isLastGoodField")) delete field.isLastGoodField;
				});

				fields.map(field => {
					if (field.show !== false && field.disabled !== true) {
						// 设置 firstGoodField
						if (!firstGoodField) {
							firstGoodField = R.deepClone(field);
							firstGoodField.isFirstGoodField = true;
							// 默认第一个[合法field]自动聚焦
							firstGoodField.autofocus = true;
						}

						// 设置 lastGoodField
						lastGoodField = R.deepClone(field);
						lastGoodField.isLastGoodField = true;
					}
				});

				// 第一个及最后一个合法的field为同一个field
				if (firstGoodField?.fieldIndex === lastGoodField?.fieldIndex) {
					const sameGoodField = { ...firstGoodField, ...lastGoodField };

					// 不为{}
					if (R.isYummy(sameGoodField)) {
						fields[lastGoodField.fieldIndex] = sameGoodField;
					}
				} else {
					if (R.isYummy(firstGoodField)) {
						fields[firstGoodField.fieldIndex] = firstGoodField;
					}

					if (R.isYummy(lastGoodField)) {
						fields[lastGoodField.fieldIndex] = lastGoodField;
					}
				}
			});

			this.fieldsList = [...this.fieldsList];

			console.log(this.fieldsList);
		},

		// 聚焦
		onFocusField(event, field, fieldsIndex, fieldIndex) {
			const { customValue: value } = event;

			this.focusFieldsIndex = fieldsIndex;
			this.focusFieldIndex = fieldIndex;

			// 默认都不聚焦
			{
				const flatFieldsList = R.flatArray(this.fieldsList);

				flatFieldsList.map(field => {
					field.autofocus = false;
				});
			}

			this.fieldsList[fieldsIndex][field.fieldIndex].autofocus = true;

			// 聚焦即更新底部 prompt
			this.$store.commit("UPDATE_CHANNEL", { prompt: field.prompt });

			// 找到一码二码第一个不同的值，并在 field 选中
			if (this.op === "opq") {
				this.opqGetFieldFirstDiffIndex(field);
			}

			this.svSearchConstants({ value, field });
		},

		// 回车
		onEnterField(event, field, fieldsIndex) {
			const { ctrlKey } = event;

			// 校验当前字段是否匹配所有校验规则(ctrlKey为true则强制通过)
			const isValidField =
				ctrlKey ||
				this.validateField({ event, field, fieldsList: this.fieldsList, fieldsIndex });

			if (!isValidField) {
				return;
			}

			// fieldsList
			this.fieldsListLength = this.fieldsList.length;

			this.svEnterUpdateField({
				field,
				fieldsList: this.fieldsList,
				focusFieldsIndex: fieldsIndex
			});

			this.findAndMarkFirstLastGoodField();

			const sameField = this.fieldsList[fieldsIndex][field.fieldIndex];

			this.ifIsLoopTrueAlwaysBeOneTempFields({ sameField, fieldsIndex });

			// this.findAndMarkScrollUpDnField()

			this.autofocusToNextField({ sameField, fieldsIndex });

			this.scrollUpDn({ field: this.nextGoodField });
		},

		// 用户输入值
		onInputField(value, field, fieldsIndex, fieldIndex) {
			this.fieldsList[fieldsIndex][fieldIndex][`${this.op}Value`] = value;
			this.fieldsList[fieldsIndex][fieldIndex].resultValue = value;

			this.svSearchConstants({ value, field });
		},

		// 若 isLoop === true，尾部永远有一个 tempFields
		ifIsLoopTrueAlwaysBeOneTempFields({ sameField, fieldsIndex }) {
			if (this.isLoop) {
				const fieldsListLastIndex = this.fieldsListLength - 1;

				if (fieldsIndex === fieldsListLastIndex && sameField.isLastGoodField) {
					const tempFields = R.deepClone(this.tempFields);

					tempFields.map((tempField, tempFieldIndex) => {
						this.setFieldEffectKeyValue(
							fieldsIndex + 1,
							tempField,
							tempFieldIndex,
							true
						);
					});

					this.fieldsList = [...this.fieldsList, tempFields];

					this.focusFieldsIndex = fieldsIndex + 1;

					this.svInit({
						fieldsList: this.fieldsList,
						codeValues: this.codeValues
					});

					this.scrollToTop();
				}
			}
		},

		// 回车后聚焦到下一个 field
		autofocusToNextField({ sameField, fieldsIndex }) {
			this.fieldsListLength = this.fieldsList.length;
			const fieldsListLastIndex = this.fieldsListLength - 1;

			// 默认均为false
			{
				const flatFieldsList = R.flatArray(this.fieldsList);

				flatFieldsList.map(field => {
					field.autofocus = false;
					field.uniqueKey = `enter_${field.uniqueId}_${Date.now()}`;
				});
			}

			const surplusFields = this.fieldsList[fieldsIndex].slice(sameField.fieldIndex + 1);
			this.nextGoodField = R.find(surplusFields, { show: true, disabled: false });

			// 当前 fields 的下一个[合法field]
			if (this.nextGoodField) {
				this.nextGoodField.autofocus = true;
			} else {
				// 下一个分块
				const focusToNextFields = () => {
					this.focusFieldsIndex = fieldsIndex + 1;
					this.nextGoodField = R.find(this.fieldsList[this.focusFieldsIndex], {
						isFirstGoodField: true
					});
					this.nextGoodField.autofocus = true;
				};

				// 最后一个 fields
				if (fieldsIndex === fieldsListLastIndex) {
					// isLoop === true: 循环分块
					if (this.isLoop) {
						focusToNextFields();
					} else {
						this.nextGoodField = sameField;
						this.nextGoodField.autofocus = true;

						this.limitSubmitTask();
					}
				} else {
					focusToNextFields();
				}
			}
		},

		// 按下向上键
		onDnKey({ ctrlKey }, field, fieldsIndex) {
			if (!ctrlKey) {
				// fieldsList
				this.fieldsListLength = this.fieldsList.length;

				const sameField = this.fieldsList[fieldsIndex][field.fieldIndex];

				this.autofocusToPrevField({ sameField, fieldsIndex });

				this.scrollUpDn({ field: this.prevGoodField });
			}
		},

		// 按下向上键后聚焦到上一个 field
		autofocusToPrevField({ sameField, fieldsIndex }) {
			// 默认均为false
			{
				const flatFieldsList = R.flatArray(this.fieldsList);

				flatFieldsList.map(field => {
					field.autofocus = false;
					field.uniqueKey = `up_${field.uniqueId}_${Date.now()}`;
				});
			}

			const reverseSurplusFields = this.fieldsList[fieldsIndex]
				.slice(0, sameField.fieldIndex)
				.reverse();
			this.prevGoodField = R.find(reverseSurplusFields, { show: true, disabled: false });

			if (this.prevGoodField) {
				this.prevGoodField.autofocus = true;
			} else {
				if (fieldsIndex > 0) {
					this.focusFieldsIndex = fieldsIndex - 1;
					this.prevGoodField = R.find(this.fieldsList[this.focusFieldsIndex], {
						isLastGoodField: true
					});
					this.prevGoodField.autofocus = true;
				} else {
					this.focusFieldsIndex = 0;
					this.prevGoodField = R.find(this.fieldsList[this.focusFieldsIndex], {
						isFirstGoodField: true
					});
					this.prevGoodField.autofocus = true;
				}
			}

			this.fieldsList = [...this.fieldsList];
		},

		// 提交(F8)
		async submitTask() {
			const data = {
				bill: this.bill,
				block: this.block,
				fields: this.fieldsList,
				op: this.op
			};

			console.log(data);

			// return

			const result = await this.$store.dispatch("UPDATE_LP_TASK", data);

			this.toasted.dynamic(result.msg, result.code, { position: "top" });

			if (result.code === 200) {
				this.prevNums = 0;
				this.getTask({ status: "new" });
			}
		},

		// 限制提交次数
		limitSubmitTask() {
			if (!this.oldTime) {
				this.submitTask();
				this.oldTime = Date.now();
			} else {
				if (Date.now() - this.oldTime > 1000) {
					this.submitTask();
					this.oldTime = Date.now();
				}
			}
		},

		// 提交(F8)
		async fEightSubmit() {
			const svValid = this.svValidateFields({
				bill: this.bill,
				block: this.block,
				fieldsList: this.fieldsList,
				op: this.op
			});
			const valid = this.$refs[`${this.op}Form`].validate();

			svValid && valid && this.limitSubmitTask();
		},

		// 聚焦field或者提交
		fuckEvents(event) {
			event = event || window.event;

			switch (event.keyCode) {
				// focus 到下一个未被屏蔽的 field(F4)
				case 115:
					event.preventDefault();

					this.svDisableFields({
						fieldsList: this.fieldsList,
						focusFieldsIndex: this.focusFieldsIndex
					});

					this.fieldsList = [...this.fieldsList];
					break;
			}
		}
	},

	computed: {
		...mapGetters(["task"])
	}
};
