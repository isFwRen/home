import { mapGetters } from "vuex";
import { tools } from "vue-rocket";
import _ from "lodash";

// 定义特殊字段的常量
const fields = [
	["图片页码", "PAGE_FIELD"],
	["模板类型字段", "TEMP_FIELD"],
	["显示范围", "RANG_FIELD"]
];

// 特殊字段的固定下标
const [PAGE_INDEX, TEMP_INDEX, RANG_INDEX] = [0, 1, 2];

// 特殊字段(special fields)
const mapSFields = new Map(fields);

// 默认不显示(default hidden fields)
const mapDHFields = new Map([fields[0], fields[2]]);

export default {
	data() {
		return {
			hasOp: false,
			valid: true,

			// config
			mb001Config: [],
			fieldsConfig: [],

			// bill
			bill: {},

			// block
			block: {},

			thumbIndex: 0,

			// fields
			fieldsObject: {},
			fieldsList: [],
			tempFields: [],
			sessionFieldsList: [],

			// field
			prevGoodField: null,
			nextGoodField: null,
			tempField: {},
			freezeTempField: {},
			pageField: {},
			rangField: {},

			// 记录用户的操作
			memoFields: {},

			// focus
			focusFieldsIndex: 0,
			focusFieldIndex: 0,

			fieldsListLength: 0,

			// 当前编辑图片
			modifyImage: "",

			// 画布宽
			canvasWidth: 0,

			// 绘图后生成的图片名
			drewImageName: "",

			oldTime: null
		};
	},

	created() {
		this.setConfigs();
		this.getTask({ status: "new" });
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
				bill = {},
				block = {},
				fields = []
			} = tools.isYummy(result.data) ? result.data : {};

			if (result.code === 200) {
				if (tools.isLousy(fields)) {
					this.toasted.warning("fields为空!", { position: "top" });
					return;
				}

				this.setBill(bill);
				this.block = block;
				this.setFieldsObject(fields);
				this.hasOp = true;
			} else {
				this.hasOp = false;
				this.toasted.dynamic(result.msg, result.code, { position: "top" });

				// 清除底部 prompt
				this.$store.commit("UPDATE_CHANNEL", { prompt: "" });
			}

			this.$emit("gotTaskResponse", {
				code: result.code,
				bill: this.bill,
				block: this.block
			});

			this.svResetVariable();
			this.scrollToTop();
		},

		// 账单
		setBill(bill) {
			if (tools.isLousy(bill.pictures)) {
				this.toasted.warning("未获取到图片！", { position: "top" });
				bill.pictures = [];
				this.thumbIndex = -1;
			}

			bill.phoTotal = bill.pictures.length;
			bill.thumbIndex = this.thumbIndex;

			this.bill = bill;
		},

		// fieldsObject
		setFieldsObject(fieldsList) {
			// 默认取第一个 fields 下的图片下标/图片路径
			if (tools.isYummy(fieldsList)) {
				this.focusFieldsIndex = 0;
				this.fieldsList = tools.deepClone(fieldsList);

				// 图片
				this.setModifyImage({ fieldsList });
			} else {
				this.fieldsList = [];
				return;
			}

			// 初始化 fieldsObject
			const length = this.bill.pictures.length;

			for (let i = 0; i < length; i++) {
				this.fieldsObject[i] = {
					sessionStorage: false,
					initFieldsList: [],
					fieldsList: []
				};
			}

			// 设置 fieldsList 默认信息
			{
				const flatFieldsList = tools.flatArray(this.fieldsList);

				flatFieldsList.map(field => {
					const result = tools.find(this.fieldsConfig, field.code);
					const { name, prompt, valChange } = result;

					field.name = name;
					field.prompt = prompt;
					field.valChange = valChange;

					// 特殊字段及 op0Input 的值不为 no 的字段才需要赋值
					if (mapSFields.get(field.name) || field.op0Input !== "no") {
						field.resultInput = field.op0Input;
					}

					// field 类型为 模板类型字段，field 显示且不禁用
					if (mapSFields.get(field.name) === "TEMP_FIELD") {
						field.show = true;
						field.disabled = false;
					}

					// field 类型为 图片页码 || 显示范围，field 不显示
					if (mapDHFields.get(field.name)) {
						field.show = false;
					}
				});
			}

			// 创建字段模板
			this.createTempFields(this.fieldsList[0].slice(0, 3));

			// 根据后端返回的fieldsList 设置 fieldsObject 的 fieldsList
			this.fieldsList.map(fields => {
				const pageField = fields[PAGE_INDEX];
				const existed = this.fieldsObject.hasOwnProperty(pageField.op0Value);

				if (existed) {
					this.fieldsObject[pageField.op0Value].fieldsList.push(fields);
				}
			});

			// fieldsObject 的 fieldsList 默认有一个[模板类型字段]field
			for (let key in this.fieldsObject) {
				if (tools.isLousy(this.fieldsObject[key].fieldsList)) {
					const tempFields = tools.deepClone(this.tempFields);

					// 生成图片页码
					const pageField = tempFields[PAGE_INDEX];
					pageField.op0Value = key;
					pageField.resultValue = pageField.op0Value;

					this.fieldsObject[key].fieldsList.push(tempFields);
				}
			}

			// 默认第一个 field 获取焦点
			for (let key in this.fieldsObject) {
				this.fieldsObject[key]?.fieldsList[0]?.map(field => {
					if (mapSFields.get(field.name) === "TEMP_FIELD") {
						field.autofocus = true;
					}
				});
			}

			// 设置前端自定义 keyValue
			for (let key in this.fieldsObject) {
				this.fieldsObject[key]?.fieldsList.map((fields, fieldsIndex) => {
					fields.map((field, fieldIndex) => {
						this.setFieldEffectKeyValue({
							thumbIndex: key,
							fieldsIndex,
							field,
							fieldIndex
						});
					});

					// 复制 fieldsList
					this.fieldsObject[key].initFieldsList = tools.deepClone(
						this.fieldsObject[key].fieldsList
					);
				});
			}

			// console.log({ fieldsObject: this.fieldsObject })
		},

		// 前端自定义 keyValue
		setFieldEffectKeyValue({ thumbIndex, fieldsIndex, field, fieldIndex }) {
			// 累加的下标
			{
				const frontFieldsList = this.fieldsObject[thumbIndex].fieldsList.slice(
					0,
					fieldsIndex
				);

				let frontFieldsLength = 0;

				frontFieldsList.map(frontFields => {
					frontFieldsLength += frontFields.length;
				});

				field.waterfallIndex = frontFieldsLength + fieldIndex;
			}

			// 唯一id
			field.uniqueId = `${thumbIndex}_${fieldsIndex}_${fieldIndex}_${field.code}`;

			// 唯一key
			field.uniqueKey = field.uniqueId;

			// 下标
			field.fieldIndex = fieldIndex;

			// 模板类型字段永远为[合法field]
			if (mapSFields.get(field.name) === "TEMP_FIELD") {
				field.show = true;
				field.disabled = false;
			}

			// 根据 op0Input 的值设置 field 状态，yes || '' 显示 field，no 隐藏 field，no_if 禁用 field
			if (!mapSFields.get(field.name)) {
				if (field.op0Input === "no") {
					field.show = false;
				} else {
					field.show = true;
				}

				if (field.op0Input === "no_if") {
					field.disabled = true;
				} else {
					field.disabled = false;
				}
			}

			// 校验规则
			field.rules = this.setValidateRules(field);
		},

		// 找到最后一个[合法field]，并标记
		findAndMarkLastGoodField(fieldsIndex) {
			const fields = this.fieldsObject[this.thumbIndex]?.fieldsList[fieldsIndex];
			let lastGoodField = null;

			// 默认不设置 lastGoodField
			fields.map(field => {
				if (field.hasOwnProperty("isLastGoodField")) {
					delete field.isLastGoodField;
				}
			});

			fields.map(field => {
				if (field.show !== false && field.disabled !== true) {
					lastGoodField = tools.deepClone(field);
					lastGoodField.isLastGoodField = true;
				}
			});

			fields[lastGoodField.fieldIndex] = lastGoodField;
		},

		// 聚焦
		onFocusField(event, field, fieldsIndex, fieldIndex) {
			const { customValue: value } = event;

			this.focusFieldsIndex = fieldsIndex;
			this.focusFieldIndex = fieldIndex;

			const fieldsList = this.fieldsObject[this.thumbIndex]?.fieldsList;

			// 图片显示
			this.setModifyImage({ fieldsList });

			// 默认都不聚焦
			{
				const flatFieldsList = tools.flatArray(fieldsList);

				flatFieldsList.map(field => {
					field.autofocus = false;
				});
			}

			this.$set(fieldsList[this.focusFieldsIndex][this.focusFieldIndex], "autofocus", true);

			// 聚焦即更新底部 prompt
			this.$store.commit("UPDATE_CHANNEL", { prompt: field.prompt });

			this.setDrewImageName();

			this.svSearchConstants({ value, field });
		},

		// 回车
		onEnterField(event, field, fieldsIndex, fieldIndex) {
			let { customValue: value } = event;

			// 校验当前字段是否匹配所有校验规则
			const isValidField = this.validateField({
				event,
				field,
				fieldsList: this.fieldsObject[this.thumbIndex].fieldsList,
				fieldIndex
			});

			if (!isValidField) {
				return;
			}

			const fieldsList = this.fieldsObject[this.thumbIndex].fieldsList;

			// 根据[模板类型字段]的值更新 fieldsList
			if (value && mapSFields.get(field.name) === "TEMP_FIELD") {
				const code = this.convertValChange(field.valChange, value);

				// 找到模板配置对应的 fields
				const result = code ? tools.find(this.mb001Config, code) : { fields: [] };

				const subFields = [];

				if (result) {
					result.fields.map((fField, fFieldIndex) => {
						const cloneTempField = tools.deepClone(this.freezeTempField);

						cloneTempField.name = fField.fName;
						cloneTempField.code = fField.fCode;

						this.setFieldEffectKeyValue({
							thumbIndex: this.thumbIndex,
							fieldsIndex,
							field: cloneTempField,
							fieldIndex: fFieldIndex + 3
						});

						// 默认值
						{
							const { value: svMemoValue, items: svMemoItems } = this.svMemoFields[
								cloneTempField.code
							] || {
								value: "",
								items: []
							};

							cloneTempField.op0Value =
								this.memoFields[cloneTempField.uniqueId]?.value || svMemoValue;
							cloneTempField.items =
								this.memoFields[cloneTempField.uniqueId]?.items || svMemoItems;
						}

						cloneTempField.resultValue = cloneTempField.op0Value;

						// 校验规则
						cloneTempField.rules = this.setValidateRules(cloneTempField);

						// prompt
						const result = tools.find(this.fieldsConfig, cloneTempField.code);
						cloneTempField.prompt = result.prompt;

						subFields.push(tools.deepClone(cloneTempField));
					});
				}

				fieldsList[fieldsIndex] = [...fieldsList[fieldsIndex].slice(0, 3), ...subFields];

				this.svUpdateDropdown();
			}

			this.findAndMarkLastGoodField(fieldsIndex);

			// this.op0FindAndMarkScrollUpDnField()

			// fieldsList
			this.fieldsListLength = fieldsList.length;

			this.alwaysBeOneTempFields(fieldsList[fieldsIndex][field.fieldIndex], fieldsIndex);

			this.autofocusToNextField(fieldsList[fieldsIndex][field.fieldIndex], fieldsIndex);

			this.fieldsObject = { ...this.fieldsObject };

			this.svUpdateMemoFields({
				value,
				code: field.code,
				items: field.items
			});

			this.svEnterUpdateField({
				field,
				fieldsList,
				focusFieldsIndex: fieldsIndex,
				memoFields: this.memoFields
			});

			this.svInit({
				fieldsList: this.fieldsObject[this.thumbIndex].fieldsList,
				focusFieldsIndex: this.focusFieldsIndex
			});

			this.scrollUpDn({ field: this.nextGoodField });

			// console.log({ fieldsObject: this.fieldsObject })
		},

		// 用户输入值
		onInputField(value, field, fieldsIndex, fieldIndex) {
			const sameField =
				this.fieldsObject[this.thumbIndex].fieldsList[fieldsIndex][fieldIndex];

			sameField.op0Value = value;
			sameField.resultValue = value;

			this.svSearchConstants({ value, field });

			this.fieldsObject = { ...this.fieldsObject };

			_.set(this.memoFields, `${field.uniqueId}.value`, value);
			_.set(this.memoFields, `${field.uniqueId}.items`, field.items);
		},

		// 创建字段模板
		createTempFields(tempFields) {
			const delKeys = ["CreatedAt", "ID", "UpdatedAt", "feedbackDate"];
			const emptyKeys = [
				"op0Value",
				"op1Input",
				"op1Value",
				"op2Input",
				"op2Value",
				"opqInput",
				"opqValue",
				"resultValue"
			];

			this.tempFields = tools.deepClone(tempFields);

			this.tempFields.map(field => {
				// 设置 tempFields(前三个) 的默认值
				for (let key in field) {
					if (delKeys.includes(key)) {
						delete field[key];
					}

					if (emptyKeys.includes(key)) {
						field[key] = "";
					}
				}

				// 每个字段的初始状态
				if (mapSFields.get(field.name) === "TEMP_FIELD") {
					const cleanValues = {
						code: "",
						disabled: false,
						fieldIndex: -1,
						name: "",
						op0Value: "",
						op0Input: "",
						prompt: "",
						resultInput: "",
						resultValue: "",
						show: true
					};

					this.freezeTempField = Object.assign(tools.deepClone(field), cleanValues);
				}
			});
		},

		// 保证当前 thumbIndex 末尾永远有一个[模板类型字段]
		alwaysBeOneTempFields(sameField, fieldsIndex) {
			const lastFieldsIndex = this.fieldsListLength - 1;

			if (fieldsIndex === lastFieldsIndex && sameField.isLastGoodField) {
				const tempFields = tools.deepClone(this.tempFields);

				// 图片页码 field
				const pageField = tempFields[PAGE_INDEX];

				tempFields.map((tempField, tempFieldIndex) => {
					this.setFieldEffectKeyValue({
						thumbIndex: this.thumbIndex,
						fieldsIndex: fieldsIndex + 1,
						field: tempField,
						fieldIndex: tempFieldIndex
					});

					// 生成页码
					if (tempField.name === pageField.name) {
						tempField.op0Value = String(this.thumbIndex);
						tempField.resultValue = tempField.op0Value;
					}
				});

				this.fieldsObject[this.thumbIndex]?.fieldsList.push(tempFields);
			}
		},

		// 自动聚焦到下一个 field
		autofocusToNextField(sameField, fieldsIndex) {
			const fieldsList = this.fieldsObject[this.thumbIndex]?.fieldsList;

			// 默认均为false
			{
				const flatFieldsList = tools.flatArray(fieldsList);

				flatFieldsList.map(field => {
					field.autofocus = false;
					field.uniqueKey = `enter_${field.uniqueId}_${Date.now()}`;
				});
			}

			const surplusFields = fieldsList[fieldsIndex].slice(sameField.fieldIndex + 1);
			this.nextGoodField = tools.find(surplusFields, { show: true, disabled: false });

			if (this.nextGoodField) {
				this.nextGoodField.autofocus = true;
				this.focusFieldIndex = this.nextGoodField.fieldIndex;
			} else {
				this.focusFieldsIndex = fieldsIndex + 1;
				this.focusFieldIndex = 1;
			}

			// 图片显示
			this.setModifyImage({ fieldsList });

			this.$set(fieldsList[this.focusFieldsIndex][this.focusFieldIndex], "autofocus", true);

			return true;
		},

		// 按下向上键
		onDnKey({ ctrlKey }, field, fieldsIndex) {
			if (!ctrlKey) {
				const fieldsList = this.fieldsObject[this.thumbIndex]?.fieldsList;

				this.autofocusToPrevField(fieldsList[fieldsIndex][field.fieldIndex], fieldsIndex);

				this.scrollUpDn({ field: this.prevGoodField });
			}
		},

		// 自动聚焦到上一个 field
		autofocusToPrevField(sameField, fieldsIndex) {
			const fieldsList = this.fieldsObject[this.thumbIndex]?.fieldsList;

			// 默认均为false
			{
				const flatFieldsList = tools.flatArray(fieldsList);

				flatFieldsList.map(field => {
					field.autofocus = false;
					field.uniqueKey = `up_${field.uniqueId}_${Date.now()}`;
				});
			}

			const reverseSurplusFields = fieldsList[fieldsIndex]
				.slice(0, sameField.fieldIndex)
				.reverse();
			this.prevGoodField = tools.find(reverseSurplusFields, { show: true, disabled: false });

			if (this.prevGoodField) {
				this.focusFieldIndex = this.prevGoodField.fieldIndex;
			} else {
				if (fieldsIndex > 0) {
					this.focusFieldsIndex = fieldsIndex - 1;
					this.prevGoodField = tools.find(fieldsList[this.focusFieldsIndex], {
						isLastGoodField: true
					});
					this.focusFieldIndex = this.prevGoodField.fieldIndex;
				} else {
					this.focusFieldsIndex = 0;
					this.focusFieldIndex = 1;
				}
			}

			this.$set(fieldsList[this.focusFieldsIndex][this.focusFieldIndex], "autofocus", true);
			this.fieldsObject = { ...this.fieldsObject };

			return true;
		},

		// 字符串转换为键值对集合
		convertValChange(oStr, value) {
			if (!oStr) return oStr;

			const str = oStr.replace(/\s/g, "");
			const arr = str.split(";");
			const valChangeMap = new Map();

			for (let val of arr) {
				if (val) {
					const item = val.split("=");
					valChangeMap.set(item[0], item[1]);
				}
			}

			return valChangeMap.get(value);
		},

		// 绘制后的图片
		async drewImage({ file }) {
			let result = {};

			// 用户操作过图片
			if (tools.isYummy(file.name)) {
				const data = {
					file,
					path: this.bill.downloadPath,
					name: `${this.bill.billName}.${this.thumbIndex}.${this.focusFieldsIndex}.png`,
					op: this.op
				};

				// console.log(`${ this.fileUrl }${ this.bill.downloadPath }${ this.modifyImage }`)

				result = await this.$store.dispatch("UPDATE_LP_TASK_FIELD_IMAGE", data);

				if (result.code === 200) {
					this.modifyImage = result.data;

					// 显示范围(存储编辑后的图片)
					this.saveModifiedImage({ sessionStorage: true });
				}
			}

			this.limitSessionSave(result);
		},

		// 临时保存(按F4则临时保存,否则不管是切换图片还是更新[模板类型字段]的值,均恢复为初始化)
		sessionSave({ code, msg }) {
			this.saveModifiedImage({ sessionStorage: true });

			if (code) {
				this.toasted.success(`${msg}，字段临时保存成功！`, { position: "top" });
			} else {
				this.toasted.success("字段临时保存成功！", { position: "top" });
			}
		},

		// 限制临时保存次数
		limitSessionSave(result) {
			if (!this.oldTime) {
				this.sessionSave(result);
				this.oldTime = Date.now();
			} else {
				if (Date.now() - this.oldTime > 1000) {
					this.sessionSave(result);
					this.oldTime = Date.now();
				}
			}
		},

		// 提交(F8)
		async submitTask() {
			// 当前 thumbIndex 下用户未按 F4，也需要保存并提交
			if (!this.fieldsObject[this.thumbIndex]?.sessionStorage) {
				this.saveModifiedImage({ sessionStorage: true });
			}

			const fieldsList = [];

			// 若有做临时保存,则提交到服务器
			for (let key in this.fieldsObject) {
				if (this.fieldsObject[key].sessionStorage) {
					// 若当前 fields 的 [模板类型字段] 不为空，则添加 fields
					this.fieldsObject[key].fieldsList.map(fields => {
						if (tools.isYummy(fields[TEMP_INDEX].op0Value)) {
							fieldsList.push(fields);
						}
					});
				}
			}

			// 默认ID
			if (tools.isYummy(fieldsList)) {
				this.fieldsList[0].map((field, index) => {
					fieldsList[0][index].ID = field.ID;
				});
			}

			const data = {
				bill: this.bill,
				block: this.block,
				fields: fieldsList,
				op: this.op
			};

			console.log(data);

			const svResult = this.svValidateFields({ fieldsList });

			// return

			if (!svResult) {
				return;
			}

			const result = await this.$store.dispatch("UPDATE_LP_TASK", data);

			this.toasted.dynamic(result.msg, result.code, { position: "top" });

			if (result.code === 200) {
				this.prevNums = 0;
				this.memoFields = {};

				this.getTask({ status: "new" });
			}

			this.$emit("submittedTaskResponse", { code: result.code });
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
		fEightSubmit() {
			{
				const fieldsList = this.fieldsObject[this.thumbIndex]?.fieldsList;
				const lastIndex = fieldsList.length - 1;
				const tempField = fieldsList[lastIndex][TEMP_INDEX];

				// 其它页面有保存，当前页面只有一个模板类型字段，且为空
				if (lastIndex === 0) {
					for (let key in this.fieldsObject) {
						if (this.fieldsObject[key].sessionStorage) {
							if (!tempField.op0Value) {
								fieldsList.length = lastIndex;
								this.focusFieldsIndex = lastIndex;
							}
							break;
						}
					}

					this.fieldsObject = { ...this.fieldsObject };
				}
				// 最后一个字段为模板类型字段，且为空
				else {
					if (!tempField.op0Value) {
						fieldsList.length = lastIndex;
						this.focusFieldsIndex = lastIndex - 1;
					}

					this.fieldsObject = { ...this.fieldsObject };
				}
			}

			this.$nextTick(() => {
				const valid = this.$refs.op0Form.validate();

				valid && this.limitSubmitTask();
			});
		},

		// 设置当前编辑图片
		setModifyImage({ fieldsList }) {
			const rangField = fieldsList[this.focusFieldsIndex][RANG_INDEX];

			if (rangField) {
				if (!rangField.op0Value) {
					rangField.op0Value =
						this.bill.pictures[this.thumbIndex] || this.bill.pictures[0];
				}

				this.modifyImage = rangField.op0Value;

				this.$set(
					fieldsList[this.focusFieldsIndex][rangField.fieldIndex],
					"op0Value",
					rangField.op0Value
				);
			}
		},

		// 存储编辑后的图片
		saveModifiedImage({ sessionStorage }) {
			this.fieldsObject[this.thumbIndex].sessionStorage = sessionStorage;

			if (sessionStorage) {
				const { fieldsList } = this.fieldsObject[this.thumbIndex];

				fieldsList[this.focusFieldsIndex]?.map(field => {
					if (mapSFields.get(field.name) === "RANG_FIELD") {
						field.op0Value = this.modifyImage;
						field.resultValue = field.op0Value;
					}
				});

				const initFieldsList = tools.deepClone(fieldsList);

				this.fieldsObject[this.thumbIndex].initFieldsList = initFieldsList;
				this.fieldsObject[this.thumbIndex].fieldsList = fieldsList;
			}
		},

		// // 恢复初始化图片
		// setInitImage() {
		//   this.modifyImage = this.bill.pictures[this.thumbIndex] || this.bill.pictures[0]
		//   this.saveModifiedImage({ sessionStorage: true })
		// },

		// 根据模板字段生成编辑后生成图片名
		setDrewImageName() {
			this.drewImageName = `${this.bill.billName}.${this.thumbIndex}.${this.focusFieldsIndex}.png`;
		}
	},

	computed: {
		...mapGetters(["task"])
	}
};
