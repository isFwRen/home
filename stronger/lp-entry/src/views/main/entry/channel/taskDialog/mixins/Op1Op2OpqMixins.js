import { mapGetters } from "vuex";
import nifty from "nifty-util";
import { tools, localStorage, sessionStorage } from "vue-rocket";
import { toastedOptions } from "../../cells";

const delKeys = ["CreatedAt", "UpdatedAt"];

// [合法field]: `${ this.op }Input` !== ('no' | 'no_if') 的 field

const defaultBill = {
	pictures: [],
	phoTotal: 0,
	pageIndex: 0
};

const defaultFocusFieldIndex = -1;

export default {
	data() {
		return {
			hasOp: false,
			valid: true,

			// config
			mb001Config: [],
			fieldsConfig: [],

			// bill
			bill: { ...defaultBill },

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
			focusFieldIndex: defaultFocusFieldIndex,

			// field
			prevGoodField: null,
			nextGoodField: null,

			clearNVs: false,

			// 计时器
			times: 15,
			// 是否开始计时
			start: false,
			// watch监听bill.ID
			billID: '',

			// 修复108一二码问题件F4光标不能自动聚焦
			fieldsIndex_s: '',
			fields_Index: ''
		};
	},

	created() {
		this.setConfigs();
		this.getTask({ status: "new" });
		// setTimeout(() => {
		// 	this.starting();
		// },1000)
	},

	mounted() {
		window.addEventListener("keydown", this.fuckEvents);
	},

	beforeDestroy() {
		window.removeEventListener("keydown", this.fuckEvents);
	},

	watch: {
		billID(newVal) {
			if (newVal) {
				this.starting();
			}
		}
	},

	methods: {
		// 重置数据
		resetData() {
			this.bill = { ...defaultBill };
			this.block = {};

			this.codeValues = {};

			this.fieldsListLength = 0;

			this.fieldsList = [];
			this.tempFields = [];
			this.focusFieldsIndex = 0;
			this.focusFieldIndex = defaultFocusFieldIndex;

			this.prevGoodField = null;
			this.nextGoodField = null;

			this.times = 15
		},

		// 配置信息
		setConfigs() {
			const { mb001, fields } = this.task.config;
			this.mb001Config = mb001;
			this.fieldsConfig = fields;
		},

		// 领任务
		async getTask({ status, prevNums }) {
			await this.reloadConstants();
			if (sessionStorage.get('isApp')?.isApp === 'true') {
				await this.reqConstants();
			}
			this.resetData();

			const user = localStorage.get("user");

			const form = {
				code: user.code,
				op: this.op
			};

			let result = null;

			switch (status) {
				case "new":
					delete form.prevNums;
					result = await this.$store.dispatch("INPUT_GET_TASK", form);
					break;

				case "modify":
					form.prevNums = prevNums;
					result = await this.$store.dispatch("INPUT_GET_PREVIOUS_TASK", form);
					break;
			}

			const {
				bill = defaultBill,
				block = {},
				codeValues = {},
				fields = []
			} = tools.isYummy(result.data) ? result.data : {};

			if (result.code === 200) {
				if (tools.isLousy(fields)) {
					this.toasted.warning("fields为空!", toastedOptions);
					return;
				}
				this.setBill(bill);
				this.block = block;
				this.isLoop = this.block.isLoop;
				this.codeValues = codeValues;
				this.setFieldsList(fields);
				this.hasOp = true;
			} else {
				this.hasOp = false;
				this.toasted.warning(result.msg, toastedOptions);
			}

			this.focusFieldIndex = 0
			this.$emit("gotTaskResponse", {
				code: result.code,
				bill: this.bill,
				block: this.block
			});

			this.svResetVariable();

			// 设置需要显示的输入框
			this.setShowRange();

			this.$store.commit('UPDATE_F3STATE', true)
			console.log("bill", this.bill);
		},

		// bill
		setBill(bill) {
			if (tools.isLousy(bill.pictures)) {
				this.toasted.warning("未获取到图片！", toastedOptions);
				bill.pictures = [];
			}

			bill.phoTotal = bill.pictures.length;
			bill.pageIndex = 0;

			this.bill = bill;
			this.billID = this.bill.ID
		},

		// fieldsList
		async setFieldsList(fieldsList) {
			if (tools.isYummy(fieldsList)) {
				this.fieldsList = tools.deepClone(fieldsList);
			} else {
				this.fieldsList = [];
				return;
			}

			// 设置 fieldsList 默认信息
			this.fieldsList.map((fields, fieldsIndex) => {
				fields.map((field, fieldIndex) => {
					field.codeValues = this.codeValues
					// resultInput(op1 | op2 | opq)：op[1|2|q]Input 的值不为 no 才需要赋值
					if (field[`${this.op}Input`] !== "no") {
						field.resultInput = field[`${this.op}Input`];
					}
					this.setFieldEffectKeyValue(fieldsIndex, field, fieldIndex);

					// 设置默认值
					this.op === "opq" && this.opqSetFieldEffectKeyValue(field);

					for (let key in field) {
						// 删除不需要的键值对
						if (delKeys.includes(key)) {
							delete field[key];
						}
					}
				});
			});
			await this.svInit({
				fieldsList: this.fieldsList,
			});
			await this.ocrPrompt({ fieldsList: this.fieldsList, })
			// 字段已生成
			this.svUpdateFields({
				fieldsList: this.fieldsList,
				codeValues: this.codeValues
			});

			// 生成字段校验规则
			for (let fields of this.fieldsList) {
				for (let field of fields) {
					const configField = tools.find(this.fieldsConfig, field.code);
					field.rules = this.setValidateRules({ field, configField });
					field.sessionRules = tools.deepClone(field.rules);
				}
			}

			this.findAndMarkFirstLastGoodField();
			// 若 isLoop === true，以 fieldsList 的第 0 个为基础创建模板
			if (this.isLoop) {
				this.tempFields = tools.deepClone(this.fieldsList[0]);

				this.tempFields.map(field => {
					delete field.ID;

					// 除特殊校验需要禁用的，默认均可录入
					{
						if (field.op1Input === "no_if") {
							field.disabled = false;
							field.op1Value = "";
							field.op1Input = "yes";
						}

						if (field.op2Input === "no_if") {
							field.disabled = false;
							field.op2Value = "";
							field.op2Input = "yes";
						}

						if (field.resultInput === "no_if") {
							field.disabled = false;
							field.resultValue = "";
							field.resultInput = "yes";
						}
					}
				});
			}
			const { firstGoodField } = this.findFieldsFirstLastGoodField(this.fieldsList[0]);

			this.focusFieldIndex = tools.isYummy(firstGoodField?.fieldIndex)
				? firstGoodField.fieldIndex
				: defaultFocusFieldIndex;

			// console.log(this.fieldsList);

			// 判断一个分块如果所有字段都处于屏蔽状态直接提交
			let showCount = 0
			let hidden = 0
			let disabledCount = 0
			for (let i = 0; i < this.fieldsList[0].length; i++) {
				if (this.fieldsList[0][i].show == false) hidden++
				if (this.fieldsList[0][i].show == true) showCount++
				if (this.fieldsList[0][i].disabled == true && this.fieldsList[0][i].show == true) disabledCount++
			}
			if (disabledCount > 0 && showCount == disabledCount) {
				this.fEightSubmit()
			}
			if (hidden == this.fieldsList[0].length) {
				this.fEightSubmit()
			}
			console.log(this.fieldsList);
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
			if (isTempField) return;

			// 根据`${ this.op }Input`的值设置 field 状态，yes || '' 显示 field，no 隐藏 field
			field.show = field[`${this.op}Input`] === "no" ? false : true;

			// 根据`${ this.op }Input`的值设置 field 状态，no_if 禁用 field
			field.disabled = field[`${this.op}Input`] === "no_if" ? true : false;

			// 底部提示
			field.prompt = tools.find(this.fieldsConfig, field.code)?.prompt;
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

				// 修复问题BUG：B0114 第一个[合法field]定位失败无法自动聚焦 原因： fc175字段 disabled == true 未生效
				if (this.bill.proCode == "B0114" && this.op != "op0") {
					fields.map(field => {
						if (field.code == "fc175") {
							field.disabled = true;
						}
					});
				}

				fields.map(field => {
					if (field.show !== false && field.disabled !== true) {
						// 设置 firstGoodField
						if (!firstGoodField) {
							firstGoodField = tools.deepClone(field);
							firstGoodField.isFirstGoodField = true;
							// 默认第一个[合法field]自动聚焦
							firstGoodField.autofocus = true;
						}

						// 设置 lastGoodField
						lastGoodField = tools.deepClone(field);
						lastGoodField.isLastGoodField = true;
					}
				});

				// 第一个及最后一个合法的field为同一个field
				if (firstGoodField?.fieldIndex === lastGoodField?.fieldIndex) {
					const sameGoodField = { ...firstGoodField, ...lastGoodField };

					// 不为{}
					if (tools.isYummy(sameGoodField)) {
						fields[lastGoodField.fieldIndex] = sameGoodField;
					}
				} else {
					if (tools.isYummy(firstGoodField)) {
						fields[firstGoodField.fieldIndex] = firstGoodField;
					}

					if (tools.isYummy(lastGoodField)) {
						fields[lastGoodField.fieldIndex] = lastGoodField;
					}
				}
			});
			this.fieldsList = [...this.fieldsList];
		},

		// 找到当前 fields 的第一个及最后一个[合法field]
		findFieldsFirstLastGoodField(fields) {
			let [firstGoodField, lastGoodField] = [null, null];

			fields.map(field => {
				if (field.show !== false && field.disabled !== true) {
					if (!firstGoodField) {
						firstGoodField = tools.deepClone(field);
					}

					lastGoodField = tools.deepClone(field);
				}
			});

			return {
				firstGoodField,
				lastGoodField
			};
		},

		// 聚焦
		async onFocusField({ event, field, fieldsIndex, fieldIndex }) {
			this.fields_Index = field.fieldIndex
			const { customValue: value } = event;

			field.force = false;

			this.focusFieldsIndex = fieldsIndex;
			this.focusFieldIndex = fieldIndex;

			// 默认都不聚焦
			{
				const flatFieldsList = tools.flatArray(this.fieldsList);

				flatFieldsList.map(field => {
					field.autofocus = false;
				});
			}
			this.fieldsList[fieldsIndex][field.fieldIndex].autofocus = true;
			// 聚焦即更新底部 prompt
			this.$store.commit("UPDATE_CHANNEL", { prompt: field.prompt });

			// 设置需要显示的输入框
			this.setShowRange();

			this.op === "opq" && this.opqMarkFieldFirstDiffIndex(field);
			if (this.proCode == 'B0108' && (this.op === "op1" || this.op === "op2")) {
				const fieldArr = ['fc080', 'fc143', 'fc145', 'fc147', 'fc149', 'fc151', 'fc153', 'fc155']
				if (fieldArr.includes(field.code) && field.resultValue.includes('?')) {
					this.opqMarkFieldFirstDiffIndex(field);
				} else {
					this.selectAll(field);
				}
			} else if (field.resultValue && !field.resultValue.includes('?') && this.op != 'opq') {
				this.selectAll(field);
			}

			if (sessionStorage.get('isApp').isApp === 'true') {
				if (value != '') {
					await this.requestDropFields({ value, field })
				}
			} else {
				this.svSearchConstants({ value, field });
			}

			this.fieldsList = [...this.fieldsList];

			this.scrollUpDn({ field });
		},

		// 回车
		onEnterField({ event, field, fieldsIndex }) {
			this.op === "opq" && this.opqNextDiffIndex(field);
			if (this.proCode == 'B0108' && (this.op === "op1" || this.op === "op2")) {
				const fieldArr = ['fc080', 'fc143', 'fc145', 'fc147', 'fc149', 'fc151', 'fc153', 'fc155']
				if (fieldArr.includes(field.code) && field.resultValue.includes('?')) {
					if (!event.ctrlKey) {
						this.opqNextDiffIndex(field);
						return
					}
				}
			}

			field.ctrlKey = event.ctrlKey;

			this.fuckAllowForce({ field });
			if (this.op == 'opq') field.allowForce = true
			// 是否允许强制通过
			if (field.allowForce !== false) {
				// ctrlKey为true则强制通过
				field.force = field.ctrlKey;

				// 强制通过需要清理当前field下的所有校验
				if (field.force) field.rules = [];
			}

			let isValid
			// 校验当前字段是否匹配所有校验规则
			if (field.resultValue == '?') {
				isValid = true
			} else {
				isValid = this.validateField({ field });
			}

			if (!isValid) return;
			// fieldsList
			this.fieldsListLength = this.fieldsList.length;

			const result = this.svEnterUpdateField({
				codeValues: this.codeValues,
				field,
				fieldsList: this.fieldsList,
				focusFieldsIndex: fieldsIndex
			});

			this.findAndMarkFirstLastGoodField();

			const sameField = this.fieldsList[fieldsIndex][field.fieldIndex];

			this.fieldsIndex_s = fieldsIndex
			this.fields_Index = field.fieldIndex

			this.ifIsLoopTrueAlwaysBeOneTempFields({ sameField, fieldsIndex, index: this.fieldsList.length - 1 });

			// 修复问题 缺少 enter回车，result== false 校验拦截
			if (result == false) return false

			this.autofocusToNextField({ sameField, fieldsIndex });

			console.log("field", field);
		},

		// 用户输入值
		async onInputField(value, field, fieldsIndex, fieldIndex) {
			field.rules = field.sessionRules;

			this.fieldsList[fieldsIndex][fieldIndex][`${this.op}Value`] = value;
			this.fieldsList[fieldsIndex][fieldIndex].resultValue = value;

			if (sessionStorage.get('isApp').isApp === 'true') {
				if (value != '') {
					await this.requestDropFields({ value, field })
				}
			} else {
				this.svSearchConstants({ value, field });
			}

			this.fieldsList = [...this.fieldsList];
		},

		// 若 isLoop === true，尾部永远有一个 tempFields
		async ifIsLoopTrueAlwaysBeOneTempFields({ sameField, fieldsIndex, index }) {
			if (this.isLoop) {
				console.log('sameField--------', sameField);
				// 修改问题：循环分块拿不到上一分块的模板-----暂时解决---后续bug未知

				// const fieldsListLastIndex = this.fieldsListLength - 1

				// fieldsIndex === fieldsListLastIndex &&
				// console.log(fieldsIndex);
				// console.log(index);
				// 114OCR循环分块bug  只在最后一页最后一个字段enter后克隆一个新的fields
				if (sameField.isLastGoodField && fieldsIndex >= index - 1) {
					const tempFields = tools.deepClone(this.tempFields);

					tempFields.map((tempField, tempFieldIndex) => {
						this.setFieldEffectKeyValue(
							fieldsIndex + 1,
							tempField,
							tempFieldIndex,
							true
						);
					});
					console.log('tempFields------------', tempFields);
					for (let field of tempFields) {
						// if (field.resultValue) {
						field.op0Value = ''
						field.op1Value = ''
						field.op2Value = ''
						field.opqValue = ''
						field.resultValue = ''
						// await this.svSearchConstants({ value: field.resultValue, field })
						// await this.hintFc({field})
						// }
					}

					this.fieldsList = [...this.fieldsList, tempFields];

					this.focusFieldsIndex = fieldsIndex + 1;


					// 字段已生成
					this.svUpdateFields({
						fieldsList: this.fieldsList,
						codeValues: this.codeValues
					});

					console.log({ fieldsList: this.fieldsList });
				}
			}
		},

		// 回车后聚焦到下一个 field
		autofocusToNextField({ sameField, fieldsIndex }) {
			this.fieldsListLength = this.fieldsList.length;
			const fieldsListLastIndex = this.fieldsListLength - 1;

			// 默认均为false
			{
				const flatFieldsList = tools.flatArray(this.fieldsList);

				flatFieldsList.map(field => {
					field.autofocus = false;
					field.uniqueKey = `enter_${field.uniqueId}_${Date.now()}`;
				});
			}
			const surplusFields = this.fieldsList[fieldsIndex].slice(sameField.fieldIndex + 1);
			this.nextGoodField = tools.find(surplusFields, { show: true, disabled: false });
			// 当前 fields 的下一个[合法field]
			if (this.nextGoodField) {
				this.nextGoodField.autofocus = true;
			} else {
				// 下一个分块
				const focusToNextFields = () => {
					this.focusFieldsIndex = fieldsIndex + 1;
					this.findAndMarkFirstLastGoodField()
					this.nextGoodField = tools.find(this.fieldsList[this.focusFieldsIndex], {
						isFirstGoodField: true
					});
					this.nextGoodField.autofocus = true;

					const { firstGoodField } = this.findFieldsFirstLastGoodField(
						this.fieldsList[this.focusFieldsIndex]
					);
					this.focusFieldIndex = firstGoodField.fieldIndex;
				};

				// 最后一个 fields
				if (fieldsIndex === fieldsListLastIndex) {
					// isLoop === true: 循环分块
					if (this.isLoop) {
						focusToNextFields();
					} else {
						this.nextGoodField = sameField;
						this.nextGoodField.autofocus = true;

						// 某些字段的校验不通过
						if (sameField.force) {
							this.submitTask();
						} else {
							this.fEightSubmit();
						}
					}
				} else {
					focusToNextFields();
				}
			}
		},

		// 按下向上键
		onDnKey({ event, field, fields, fieldsIndex, fieldIndex }) {
			if (!event.ctrlKey) {
				// fieldsList
				this.fieldsListLength = this.fieldsList.length;

				const sameField = this.fieldsList[fieldsIndex][field.fieldIndex];

				this.autofocusToPrevField({ sameField, fieldsIndex });
			}
		},

		// 按下向上键后聚焦到上一个 field
		autofocusToPrevField({ sameField, fieldsIndex }) {
			// 默认均为false
			{
				const flatFieldsList = tools.flatArray(this.fieldsList);

				flatFieldsList.map(field => {
					field.autofocus = false;
					field.uniqueKey = `up_${field.uniqueId}_${Date.now()}`;
				});
			}

			const reverseSurplusFields = this.fieldsList[fieldsIndex]
				.slice(0, sameField.fieldIndex)
				.reverse();
			this.prevGoodField = tools.find(reverseSurplusFields, { show: true, disabled: false });

			if (this.prevGoodField) {
				this.prevGoodField.autofocus = true;
			} else {
				if (fieldsIndex > 0) {
					this.focusFieldsIndex = fieldsIndex - 1;
					this.prevGoodField = tools.find(this.fieldsList[this.focusFieldsIndex], {
						isLastGoodField: true
					});
					this.prevGoodField.autofocus = true;

					const { lastGoodField } = this.findFieldsFirstLastGoodField(
						this.fieldsList[this.focusFieldsIndex]
					);
					this.focusFieldIndex = lastGoodField.fieldIndex;
				} else {
					this.focusFieldsIndex = 0;
					this.prevGoodField = tools.find(this.fieldsList[this.focusFieldsIndex], {
						isFirstGoodField: true
					});
					this.prevGoodField.autofocus = true;
				}
			}

			this.fieldsList = [...this.fieldsList];
		},

		// 提交当前分块
		async submitTask() {
			let data = {
				bill: this.bill,
				block: this.block,
				fields: nifty.deepClone(this.fieldsList),
				op: this.op
			};

			// 删除前端设置的key
			for (let fields of data.fields) {
				for (let field of fields) {
					field.resultValue = field.resultValue.trim()
					field.op0Value = field.op0Value.trim()
					field.op1Value = field.op1Value.trim()
					field.op2Value = field.op2Value.trim()
					field.opqValue = field.opqValue.trim()
					delete field.autofocus;
					delete field.desserts;
					delete field.disabled;
					delete field.effectValidations;
					delete field.items;
					delete field.rules;
					delete field.uniqueId;
					delete field.uniqueKey;
					delete field.table;
					delete field.codeValues;
				}
			}

			let surplus = JSON.parse(JSON.stringify(data.fields))
			for (let fields of surplus) {
				if (fields[0].hasOwnProperty('ID')) continue
				else {
					console.log(123456);
					let count = 0
					for (let field of fields) {
						if (field.resultValue == '') {
							count++
						}
					}
					if (count == fields.length) fields.length = 0
				}
			}
			surplus = surplus.filter(el => {
				return el.length != 0
			})
			data = { ...data, fields: surplus }


			//开发环境默认不提交
			if (process.env.NODE_ENV === 'development') {
				console.log("提交当前分块，开发环境默认不提交------", data)
				return
			}

			// let surplus = data.fields
			// surplus = surplus.filter(el => {
			// 	return el[0].hasOwnProperty('ID')
			// })

			// data = {...data, fields: surplus}

			// console.log("提交当前分块，开发环境默认不提交------", data);

			const result = await this.$store.dispatch("INPUT_SUBMIT_TASK", data);

			this.toasted.dynamic(result.msg, result.code, toastedOptions);

			if (result.code === 200) {
				this.prevNums = 0;
				this.getTask({ status: "new" });
				this.$nextTick(() => {
					this.$refs.col.scrollTop = 0
				})
			}
		},

		// 按F8提交
		async fEightSubmit() {
			this.clearNVs = true;
			this.$nextTick(async () => {
				const svValid = await this.svValidateFields();
				let valid;
				if (this.$refs[`${this.op}Form`].validate.length == 0) {
					valid = true;
				} else {
					valid = this.$refs[`${this.op}Form`].validate();
				}
				console.log(this.$refs[`${this.op}Form`]);
				console.log({ F8校验: svValid, valid });
				svValid && valid && this.limitSubmitTask();
				this.clearNVs = false;
			});
		},

		// 聚焦field或者提交
		fuckEvents(event) {
			event = event || window.event;

			switch (event.keyCode) {
				// focus 到下一个未被屏蔽的 field(F4)
				case 115:
					event.preventDefault();

					this.svDisableFields();

					// 修复一二码问题件F4后不能自动聚焦的问题
					this.fieldsList = [...this.fieldsList];
					const sameField = this.fieldsList[this.focusFieldsIndex][this.fields_Index];
					this.autofocusToNextField({ sameField, fieldsIndex: this.focusFieldsIndex });
					break;
			}
		},

		// 控制需要显示的field
		showRange(field, fieldIndex, fieldsIndex) {
			// (fieldIndex >= focusFieldIndex - upperShowRange && fieldIndex <= focusFieldIndex + lowerShowRange) && (focusFieldsIndex === fieldsIndex && field.show)
			let isYummyUpper = fieldIndex >= this.focusFieldIndex - this.upperShowRange;
			let isYummylower = fieldIndex <= this.focusFieldIndex + this.lowerShowRange;

			if (this.focusFieldIndex === -1) {
				isYummyUpper = true;
				isYummylower = true;
			}
			if (
				isYummyUpper &&
				isYummylower &&
				this.focusFieldsIndex === fieldsIndex &&
				field.show
			) {
				return true;
			}

			return false;
		},

		// 是否为OCR识别字段， 是开始计时
		starting() {
			this.start = false
			for (let el of this.fieldsList) {
				for (let item of el) {
					if (item.op0Input == 'ocr') {
						this.start = true
						break
					}
				}
				if (this.start) break
			}

			if (this.start) {
				let timer = setInterval(() => {
					this.times--
					if (this.times == 0) {
						clearInterval(timer)
					}
				}, 1000)
			}
		}

		// 图片路径
		// imgSrc(proCode, fileUrl, downloadPath, picture) {
		// 	if (proCode == "B0108" || proCode == "B0114" || proCode == "B0118") {
		// 		return `${fileUrl}${downloadPath}${picture}`;
		// 	} else {
		// 		return `${fileUrl.replace("files/", "")}${downloadPath}${picture}`;
		// 	}
		// },
	},

	computed: {
		...mapGetters(["task"]),

		computedPrevFieldsName() {
			return function (field) {
				if (/项目名称[0-9]*/.test(field.name)) {
					return `${field.name}：${field.resultValue}`;
				}
			};
		}
	},
};
