import _ from "lodash";
import { mapGetters } from "vuex";
import nifty from "nifty-util";
import { tools, localStorage, sessionStorage } from "vue-rocket";
import { toastedOptions } from "../../cells";
import { MessageBox } from 'element-ui';
import moment from 'moment'

// 定义特殊字段的常量
const fields = [
	["图片页码", "PAGE_FIELD"],
	["模板类型字段", "TEMP_FIELD"],
	["显示范围", "RANG_FIELD"]
];
// F8开关锁
// let flagF8 = true

// 特殊字段的固定下标
const [PAGE_INDEX, TEMP_INDEX, RANG_INDEX] = [0, 1, 2];

// 特殊字段(special fields)
const mapSFields = new Map(fields);

// 默认不显示(default hidden fields)
const mapDHFields = new Map([fields[0], fields[2]]);

const debounce = (() => {
	let timer = null;

	return (fn, delay = 333) => {
		if (timer) {
			clearTimeout(timer);
		}

		timer = setTimeout(() => {
			fn();
		}, delay);
	};
})();

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

			// 排序后对应的图片index
			newfieldsObjectArray: {},
			imgArray: [],
			// 图片标红的
			redimgArray: [],
			// 首次默认选中第一张
			flag: true,
			// 缩略图分类名称
			imgName: [],

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

			// 记录用户录入值
			memoFields: {},

			// 更新操作记录(切图范围、旋转方向、宽、高)
			memoMarks: {},

			// focus
			focusFieldsIndex: -1,
			// 判断是否在新的[模板类型字段]
			prevFocusFieldsIndex: -1,
			focusFieldIndex: 0,

			fieldsListLength: 0,

			// 图片加载中
			imageLoading: false,
			// 当前编辑图片
			modifyImage: "",
			// 图片大小
			imageSize: 0,

			// 绘图后生成的图片名
			drewImageName: "",

			// 释放单
			release: 0,
			// 修改时间
			cacheTime: '',

			// F8锁
			downF8: true,

			clearNVs: false,

			// 计时器
			timer: null,

			// 记录field
			FileLists: [],

			// 记录每组field的校验方法
			fieldCheckMethod: {},

			// 记录Enter按键
			recordEnter: '',

			// 记录F4按键
			recordF4: '',

			// 记录F8按键
			recordF8: ''
		};
	},

	created() {
		this.setConfigs();
		this.getTask({ status: "new" });
	},

	methods: {
		// 重置数据
		resetData() {
			this.bill = {};
			this.block = {};

			this.thumbIndex = 0;
			this.newfieldsObjectArray = {};
			this.flag = true;
			this.fieldsObject = {};
			this.fieldsList = [];
			this.tempFields = [];
			this.sessionFieldsList = [];

			this.prevGoodField = null;
			this.nextGoodField = null;
			this.tempField = {};
			this.freezeTempField = {};
			this.pageField = {};
			this.rangField = {};

			this.memoFields = {};

			this.memoMarks = {};

			this.focusFieldsIndex = -1;

			this.prevFocusFieldsIndex = -1;
			this.focusFieldIndex = 0;

			this.fieldsListLength = 0;

			this.modifyImage = "";

			this.drewImageName = "";

			this.clearNVs = false;

			this.prevNums = 0;

			this.FileLists = [];

			this.recordEnter = '';

			this.recordF4 = '';

			this.recordF8 = ''
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
			this.$store.commit("RESET_LOG");

			this.bill = {};
			this.block = {};
			this.codeValues = {};

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
					this.cacheTime = result.data.cacheTime
					break;
			}

			// F3修改初审内容
			if (status == 'modify') {
				// this.release = prevNums
				console.log(this.cacheTime);
				this.toasted.dynamic(`请注意修改时间为${this.cacheTime}秒，超时自动提交`, { duration: 3000 })
				setTimeout(async () => {
					// if (this.release == prevNums) {

					if (sessionStorage.get('isApp')?.isApp === 'true') {
						await MessageBox.alert(`修改时间已超时${this.cacheTime}秒，此初审内容已被提交`, '请检查', {
							type: 'warning',
							confirmButtonText: '确定',
							showClose: false,
						})
					} else {
						alert(`修改时间已超时${this.cacheTime}秒，此初审内容已被提交`)
					}

					this.getTask({ status: "new" });
					// this.navigatePage('Op0')
					setTimeout(() => {
						this.$store.commit('UPDATE_F3STATE', true)
					}, 5000);
					// }
				}, this.cacheTime * 1000);
			}


			const { bill = {}, block = {}, fields = [] } = tools.isYummy(result.data) ? result.data : {};

			if (result.code === 200) {
				if (tools.isLousy(fields)) {
					this.toasted.warning("fields为空!", toastedOptions);
					return;
				}

				this.setBill(bill);
				this.block = block;
				this.setFieldsObject(fields);
				this.hasOp = true;
				console.log("bill", this.bill);
			} else {
				this.hasOp = false;
				this.toasted.warning(result.msg, toastedOptions);
				console.log(result.msg);
				// 清除底部 prompt
				this.$store.commit("UPDATE_CHANNEL", { prompt: "" });
				return
			}

			this.$emit("gotTaskResponse", { code: result.code, bill: this.bill, block: this.block });

			this.svResetVariable();

			await this.svInit();

			await this.svUpdateFields({
				fieldsList: this.fieldsObject[this.thumbIndex]?.fieldsList,
				focusFieldsIndex: this.focusFieldsIndex
			});
			// await this.ocrPrompt({ fieldsList: this.fieldsObject[this.thumbIndex]?.fieldsList, })
		},

		// 账单
		setBill(bill) {
			if (tools.isLousy(bill.pictures)) {
				this.toasted.warning("未获取到图片！", toastedOptions);
				bill.pictures = [];
				this.thumbIndex = -1;
			}

			bill.phoTotal = bill.pictures.length;
			bill.thumbIndex = this.thumbIndex;

			this.bill = bill;
		},

		// fieldsObject
		setFieldsObject(fieldsList) {
			if (!tools.isYummy(this.bill.pictures)) {
				this.toasted.warning("获取图片失败，请联系后端!", toastedOptions);
				return;
			}

			// 初始化 fieldsObject
			this.bill.pictures.map((path, index) => {
				this.fieldsObject[index] = {
					sessionStorage: false,
					path: path || this.bill.pictures[0],
					initFieldsList: [],
					fieldsList: []
				};
			});

			if (this.bill.imagesType && this.bill.imagesType[0]) {
				this.imagesSort();
			}

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
					this.fieldsObject[pageField.op0Value].fieldsList.push(tools.deepClone(fields));
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

					this.fieldsObject[key].fieldsList.push(tools.deepClone(tempFields));
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
					this.fieldsObject[key].initFieldsList = tools.deepClone(this.fieldsObject[key].fieldsList);
				});
			}

			this.getNodes();
			// console.log('img_pages_fieldsObject', { fieldsObject: this.fieldsObject })
		},

		// 图片分类排序
		// 申请书—>发票信息—>结算单—>第三方报销—>清单—>疾病诊断—>身份信息—>银行卡正面—>银行卡背面
		imagesSort() {
			let shenqingshu = [];
			let fapiao = [];
			let jiesuandan = [];
			let qingdan = [];
			let zhenduan = [];
			let imagesSortArray = [];
			let shenqingshuObj = [];
			let shenqingshuName = [];
			let fapiaoObj = [];
			let fapiaoName = [];
			let jiesuandanObj = [];
			let jiesuandanName = [];
			let qingdanObj = [];
			let qingdanName = [];
			let zhenduanObj = [];
			let zhenduanName = [];
			let shenfenzheng1 = [];
			let shenfenzheng1Name = [];
			let shenfenzheng1Obj = [];
			let yinhangka = [];
			let yinhangkaName = [];
			let yinhangkaObj = [];

			let other = [];
			let otherObj = [];
			let otherName = [];

			let baoxiao = []
			let baoxiaoName = []
			let baoxiaoObj = [];

			let final = []
			let finalName = []
			let finalObj = [];

			this.imgArray = [];
			if (this.bill.imagesType.length !== 0) {
				for (let i = 0, j = this.bill.pictures.length; i < j; i++) {
					if (this.bill.imagesType[i]?.includes("申请书")) {
						shenqingshu = [...shenqingshu, this.bill.pictures[i]];
						shenqingshuObj.push(i);
						shenqingshuName.push(this.bill.imagesType[i]);
					} else if (this.bill.imagesType[i].includes('发票')) {
						fapiao = [...fapiao, this.bill.pictures[i]];
						fapiaoObj.push(i);
						fapiaoName.push(this.bill.imagesType[i]);
					} else if (this.bill.imagesType[i].includes('结算')) {
						jiesuandan = [...jiesuandan, this.bill.pictures[i]];
						jiesuandanObj.push(i);
						jiesuandanName.push(this.bill.imagesType[i]);
					} else if (this.bill.imagesType[i].includes('报销')) {
						baoxiao = [...qingdan, this.bill.pictures[i]];
						baoxiaoObj.push(i);
						baoxiaoName.push(this.bill.imagesType[i]);
					} else if (this.bill.imagesType[i].includes('清单')) {
						qingdan = [...qingdan, this.bill.pictures[i]];
						qingdanObj.push(i);
						qingdanName.push(this.bill.imagesType[i]);
					} else if (this.bill.imagesType[i].includes('诊断')) {
						zhenduan = [...zhenduan, this.bill.pictures[i]];
						zhenduanObj.push(i);
						zhenduanName.push(this.bill.imagesType[i]);
					} else if (this.bill.imagesType[i].includes('身份证')) {
						shenfenzheng1 = [...shenfenzheng1, this.bill.pictures[i]];
						shenfenzheng1Obj.push(i);
						shenfenzheng1Name.push(this.bill.imagesType[i]);
					} else if (this.bill.imagesType[i].includes('银行卡')) {
						yinhangka = [...yinhangka, this.bill.pictures[i]];
						yinhangkaObj.push(i);
						yinhangkaName.push(this.bill.imagesType[i]);
					} else if (this.bill.imagesType[i] == '') {
						other = [...other, this.bill.pictures[i]];
						otherObj.push(i);
						otherName.push('其它')
					} else {
						final = [...final, this.bill.pictures[i]];
						finalObj.push(i);
						finalName.push(this.bill.imagesType[i]);
					}
				}
			}
			// 将打乱的图片排序
			imagesSortArray = [
				...shenqingshu,
				...fapiao,
				...jiesuandan,
				...baoxiao,
				...qingdan,
				...zhenduan,
				...shenfenzheng1,
				...yinhangka,
				...final,
				...other
			];
			// 排序后图像的页码数组
			this.imgArray = [
				...shenqingshuObj,
				...fapiaoObj,
				...jiesuandanObj,
				...baoxiaoObj,
				...qingdanObj,
				...zhenduanObj,
				...shenfenzheng1Obj,
				...yinhangkaObj,
				...finalObj,
				...otherObj
			];
			// 图像分类名
			this.imgName = [
				...shenqingshuName,
				...fapiaoName,
				...jiesuandanName,
				...baoxiaoName,
				...qingdanName,
				...zhenduanName,
				...shenfenzheng1Name,
				...yinhangkaName,
				...finalName,
				...otherName
			];
			// 图片标红的数组
			this.redimgArray = [
				...shenqingshuObj,
				...fapiaoObj,
				...jiesuandanObj,
				...baoxiaoObj,
				...qingdanObj,
				...zhenduanObj,
				...shenfenzheng1Obj,
				...yinhangkaObj,
				// ...finalObj,
			];
			// 替换掉原来展示顺序的图片数组
			this.bill.pictures = imagesSortArray;
			// 构造一个排序后的图片页码与之前页码一一映射的对象
			this.imgArray.forEach((item, index) => {
				this.newfieldsObjectArray[index] = item;
			});
		},

		// 前端自定义 keyValue
		setFieldEffectKeyValue({ thumbIndex, fieldsIndex, field, fieldIndex }) {
			// 累加的下标
			{
				const frontFieldsList = this.fieldsObject[thumbIndex].fieldsList.slice(0, fieldsIndex);

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

			// 获取字段配置的校验规则
			const configField = tools.find(this.fieldsConfig, field.code);

			// 校验规则
			field.rules = this.setValidateRules({ field, configField });
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
			if (field.resultValue && !field.resultValue.includes('?') && this.op != 'opq') {
				this.selectAll(field);
			}
			debounce(async () => {
				const { customValue: value } = event;

				this.focusFieldsIndex = fieldsIndex;
				this.focusFieldIndex = fieldIndex;

				const fieldsList = this.fieldsObject[this.thumbIndex]?.fieldsList;

				// 图片显示
				if (this.prevFocusFieldsIndex !== this.focusFieldsIndex) {
					this.prevFocusFieldsIndex = this.focusFieldsIndex;
					this.setModifyImage({ fieldsList });
					this.flag = false;
				}

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

				if (sessionStorage.get('isApp')?.isApp === 'true') {
					if (value != '') {
						await this.requestDropFields({ value: value || event.target.value, field })
					}
				} else {
					this.svSearchConstants({ value: value || event.target.value, field });
				}

				this.scrollUpDn({ field });

				this.getImageSize();

				// OpSpecificValidationsMixins.js  初审F4校验传值
				this.getFieldFirst(field, fieldsIndex, fieldIndex)
				this.recordF4 = '第' + this.$store.state['recordKey'].page + '页' + ',' + field.code + ':' + field.resultValue
				this.recordF8 = '第' + this.$store.state['recordKey'].page + '页' + ',' + field.code + ':' + field.resultValue
				// console.log(fieldsList[fieldsIndex]);
			});
		},

		// 回车
		onEnterField(event, field, fieldsIndex) {
			// this.recordKey.push('按键:Enter回车' + ',' + this.recordEnter)
			// console.log(this.recordKey);
			if (event.ctrlKey) {
				let operate = '强过:' + 'Ctrl + Enter' + ',' + this.recordEnter
				this.$store.commit('UPDATE_KEY', operate)

			} else this.$store.commit('UPDATE_KEY', '按键:Enter回车' + ',' + this.recordEnter)

			let { customValue: value } = event;

			this.svUpdateMemoFieldValues({
				code: field.code
			});

			field.ctrlKey = event.ctrlKey;

			// 是否允许强制通过
			if (field.allowForce !== false) {
				// ctrlKey为true则强制通过
				field.force = field.ctrlKey;

				// 强制通过需要清理当前field下的所有校验
				if (field.force) field.rules = [];
			}

			// 校验当前字段是否匹配所有校验规则
			const isValid = this.validateField({ field });

			if (!isValid) return;

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
							const { items: svMemoItems, value: svMemoValue } = this.svMemoFields[cloneTempField.code] || {
								value: "",
								items: []
							};

							const { items: memoItems, value: memoValue } = this.returnMemoField(cloneTempField.uniqueId);

							cloneTempField.items = memoItems || svMemoItems;
							cloneTempField.op0Value = memoValue || svMemoValue || tools.find(this.fieldsConfig, fField.fCode)?.defaultVal;

							// cloneTempField.op0Value = this.memoFields[cloneTempField.uniqueId]?.value || svMemoValue
							// cloneTempField.items = this.memoFields[cloneTempField.uniqueId]?.items || svMemoItems
						}

						cloneTempField.resultValue = cloneTempField.op0Value;

						// prompt
						const configField = tools.find(this.fieldsConfig, cloneTempField.code);
						cloneTempField.prompt = configField?.prompt;

						// 校验规则
						cloneTempField.rules = this.setValidateRules({ field: cloneTempField, configField });

						subFields.push(tools.deepClone(cloneTempField));
					});
				}

				fieldsList[fieldsIndex] = [...fieldsList[fieldsIndex].slice(0, 3), ...subFields];

				this.svUpdateDropdown();
			}

			// fieldsList
			this.fieldsListLength = fieldsList.length;

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

			this.svUpdateFields({
				fieldsList: this.fieldsObject[this.thumbIndex].fieldsList,
				focusFieldsIndex: this.focusFieldsIndex
			});

			this.findAndMarkLastGoodField(fieldsIndex);

			this.alwaysBeOneTempFields(fieldsList[fieldsIndex][field.fieldIndex], fieldsIndex);

			this.autofocusToNextField(fieldsList[fieldsIndex][field.fieldIndex], fieldsIndex);

			// this.fieldsObject = { ...this.fieldsObject }

			// console.log({ fieldsObject: this.fieldsObject })
			console.log("field", field);
		},

		// 模板类型字段输入内容自动生成跳转
		onAutoEnterField(value, field, fieldsIndex) {
			this.svUpdateMemoFieldValues({
				code: field.code
			});

			// 校验当前字段是否匹配所有校验规则
			const isValid = this.validateField({ field });

			if (!isValid) return;

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
							const { items: svMemoItems, value: svMemoValue } = this.svMemoFields[cloneTempField.code] || {
								value: "",
								items: []
							};

							const { items: memoItems, value: memoValue } = this.returnMemoField(cloneTempField.uniqueId);

							cloneTempField.items = memoItems || svMemoItems;
							// cloneTempField.op0Value = memoValue || svMemoValue || tools.find(this.fieldsConfig, fField.fCode)?.defaultVal;
							cloneTempField.op0Value = memoValue || svMemoValue || tools.find(this.fieldsConfig, fField.fCode)?.defaultVal;

							// cloneTempField.op0Value = this.memoFields[cloneTempField.uniqueId]?.value || svMemoValue
							// cloneTempField.items = this.memoFields[cloneTempField.uniqueId]?.items || svMemoItems
						}

						cloneTempField.resultValue = cloneTempField.op0Value;

						// prompt
						const configField = tools.find(this.fieldsConfig, cloneTempField.code);
						cloneTempField.prompt = configField?.prompt;

						// 校验规则
						cloneTempField.rules = this.setValidateRules({ field: cloneTempField, configField });

						subFields.push(tools.deepClone(cloneTempField));
					});
				}

				// console.log(subFields);
				fieldsList[fieldsIndex] = [...fieldsList[fieldsIndex].slice(0, 3), ...subFields];

				this.svUpdateDropdown();
			}

			// fieldsList
			this.fieldsListLength = fieldsList.length;
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

			this.svUpdateFields({
				fieldsList: this.fieldsObject[this.thumbIndex].fieldsList,
				focusFieldsIndex: this.focusFieldsIndex
			});

			this.findAndMarkLastGoodField(fieldsIndex);

			this.alwaysBeOneTempFields(fieldsList[fieldsIndex][field.fieldIndex], fieldsIndex);

			this.autofocusToNextField(fieldsList[fieldsIndex][field.fieldIndex], fieldsIndex);

			// this.fieldsObject = { ...this.fieldsObject }

			// console.log({ fieldsObject: this.fieldsObject })
			console.log("field", field);
		},

		// 用户输入值
		async onInputField(value, field, fieldsIndex, fieldIndex) {
			if (!/[\u4e00-\u9fa5]/.test(value)) {
				this.recordEnter = moment(new Date).format('YYYY/MM/DD HH:mm:ss') + ',' + '第' + this.$store.state['recordKey'].page + '页' + ',' + field.code + ':' + value
			}

			const sameField = this.fieldsObject[this.thumbIndex].fieldsList[fieldsIndex][fieldIndex];
			sameField.op0Value = value;
			sameField.resultValue = value;

			// console.log("用户输入值value", value);
			// && sessionStorage.get('isApp').isApp
			if (sessionStorage.get('isApp')?.isApp === 'true') {
				if (value != '') {
					await this.requestDropFields({ field, value })
				}
			} else {
				this.svSearchConstants({ field, value });
			}

			this.fieldsObject = { ...this.fieldsObject };

			_.set(this.memoFields, `${field.uniqueId}.value`, value);
			_.set(this.memoFields, `${field.uniqueId}.items`, field.items);
			// console.log(value);

			// const { uniqueId } = field
			const { nValidations, sValidations, nvArgs, svArgs } = this.$refs[field.uniqueId][0]
			if (nValidations && nvArgs) {
				field.checkMethod = { nValidations, sValidations, nvArgs, svArgs }
			}
			// console.log(sameField);
			// console.log(this.fieldsObject);
			// this.FileLists = []
			// let unIDArray = []
			// for (let page in this.fieldsObject) {
			// 	for (let fields of this.fieldsObject[page].fieldsList) {
			// 		for (let _field of fields) {
			// 			if (_field.name != '图片页码' && _field.name != '显示范围') {
			// 				unIDArray.push(_field.uniqueId)
			// 			}
			// 		}
			// 	}
			// }

			// let flag = false
			// for (let el of this.FileLists) {
			// 	if (el && el.uniqueId == uniqueId) {
			// 		flag = true
			// 		break
			// 	}
			// }
			// if (!flag) {
			// 	this.FileLists.push(field)
			// } else {
			// 	for (let index in this.FileLists) {
			// 		if (this.FileLists[index] && this.FileLists[index].uniqueId == uniqueId && unIDArray.includes(uniqueId)) {
			// 			this.FileLists[index] = field
			// 		}
			// 	}
			// }
			// this.FileLists = this.FileLists.filter(el => {
			// 	return el.resultValue != ''
			// })

			// console.log(unIDArray);
			// console.log(this.FileLists);
			if (field.name === '模板类型字段') {
				clearTimeout(this.timer)
				this.timer = setTimeout(() => {
					if (!/[\u4e00-\u9fa5]/.test(value)) {
						console.log(value);
						this.recordEnter = moment(new Date).format('YYYY/MM/DD HH:mm:ss') + ',' + '第' + this.$store.state['recordKey'].page + '页' + ',' + field.code + ':' + value
						this.$store.commit('UPDATE_KEY', '按键:Enter回车' + ',' + this.recordEnter)
					}
					this.onAutoEnterField(value, field, fieldsIndex)
				}, 200);
			}
		},

		// 创建字段模板
		createTempFields(tempFields) {
			const delKeys = ["CreatedAt", "ID", "UpdatedAt", "feedbackDate"];
			const emptyKeys = ["op0Value", "op1Input", "op1Value", "op2Input", "op2Value", "opqInput", "opqValue", "resultValue"];

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

			// next field
			if (this.nextGoodField) {
				this.nextGoodField.autofocus = true;
				this.focusFieldIndex = this.nextGoodField.fieldIndex;
			}
			// next fields
			else {
				this.focusFieldsIndex = fieldsIndex + 1;
				this.focusFieldIndex = TEMP_INDEX;
			}

			this.$set(fieldsList[this.focusFieldsIndex][this.focusFieldIndex], "autofocus", true);

			return true;
		},

		// 按下向上键
		onDnKey({ ctrlKey }, field, fieldsIndex) {
			if (!ctrlKey) {
				const fieldsList = this.fieldsObject[this.thumbIndex]?.fieldsList;

				this.autofocusToPrevField(fieldsList[fieldsIndex][field.fieldIndex], fieldsIndex);
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

			const reverseSurplusFields = fieldsList[fieldsIndex].slice(0, sameField.fieldIndex).reverse();
			this.prevGoodField = tools.find(reverseSurplusFields, { show: true, disabled: false });

			if (this.prevGoodField) {
				this.focusFieldIndex = this.prevGoodField.fieldIndex;
			} else {
				// previous fields
				if (fieldsIndex > 0) {
					this.focusFieldsIndex = fieldsIndex - 1;
					this.prevGoodField = tools.find(fieldsList[this.focusFieldsIndex], { isLastGoodField: true });
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

		// 自动聚焦到顶层 field
		autofocusToTopField() {
			const fieldsList = this.fieldsObject[this.thumbIndex]?.fieldsList;

			// 默认均为false
			{
				const flatFieldsList = tools.flatArray(fieldsList);

				flatFieldsList.map(field => {
					field.autofocus = false;
					field.uniqueKey = `up_${field.uniqueId}_${Date.now()}`;
				});
			}

			this.$set(this.fieldsObject[this.thumbIndex].fieldsList[0][1], "autofocus", true);
			this.fieldsObject = { ...this.fieldsObject };

			return true;
		},

		// 自动聚焦到底部 field
		autofocusToBottomField() {
			const fieldsList = this.fieldsObject[this.thumbIndex]?.fieldsList;

			// 默认均为false
			{
				const flatFieldsList = tools.flatArray(fieldsList);

				flatFieldsList.map(field => {
					field.autofocus = false;
					field.uniqueKey = `up_${field.uniqueId}_${Date.now()}`;
				});
			}

			let flag1 = this.fieldsObject[this.thumbIndex].fieldsList.length - 1
			let flags = this.fieldsObject[this.thumbIndex].fieldsList[flag1].length - 1

			if (this.fieldsObject[this.thumbIndex].fieldsList[flag1].length == 3) {
				this.$set(this.fieldsObject[this.thumbIndex].fieldsList[flag1][1], "autofocus", true);
			}

			if (this.fieldsObject[this.thumbIndex].fieldsList[flag1].length > 3) {
				console.log(this.fieldsObject[this.thumbIndex].fieldsList[flag1]);
				let flag2
				for (let index in this.fieldsObject[this.thumbIndex].fieldsList[flag1]) {
					if (index > 2 && this.fieldsObject[this.thumbIndex].fieldsList[flag1][index] == '') flag2 = index
				}
				if (this.fieldsObject[this.thumbIndex].fieldsList[flag1][flags] != '') flag2 = flags
				this.$set(this.fieldsObject[this.thumbIndex].fieldsList[flag1][flag2], "autofocus", true);
			}

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

		// 临时保存(按F4)自动聚焦到下一组[模板类型字段]
		autofocusToNextFields() {
			const fieldsList = this.fieldsObject[this.thumbIndex]?.fieldsList;
			const lastFieldsIndex = fieldsList.length - 1;

			// 若为最后一组，不执行
			if (this.focusFieldsIndex === lastFieldsIndex) return;

			// 默认均为false
			{
				const flatFieldsList = tools.flatArray(fieldsList);

				flatFieldsList.map(field => {
					field.autofocus = false;
					field.uniqueKey = `enter_${field.uniqueId}_${Date.now()}`;
				});
			}

			this.focusFieldsIndex += 1;
			this.focusFieldIndex = TEMP_INDEX;

			this.$set(fieldsList[this.focusFieldsIndex][this.focusFieldIndex], "autofocus", true);

			this.setModifyImage({ fieldsList });
		},

		// 记录需要临时保存的字段
		async setSessionSaveField() {
			this.clearValues = {}
			this.clearFieldValues = {}
			this.f4SetFieldValues()
			const { fieldsList } = this.fieldsObject[this.thumbIndex];
			// fieldsList.map(fields => {
			//   fields.map(field => {
			//     field.sessionStorage = true
			//   })
			// })
			console.log('记录需要临时保存的字段---fieldList----', fieldsList);
			for (let i = 0, len = fieldsList.length; i < len; i++) {
				for (let j = 0, lens = fieldsList[i].length; j < lens; j++) {
					if (i == 0) {
						fieldsList[i][j].sessionStorage = true;
						if (typeof fieldsList[i][j].op0Value === 'object') {
							fieldsList[i][j].op0Value = ''
							fieldsList[i][j].F4Value = fieldsList[i][j].op0Value;
						} else {
							fieldsList[i][j].F4Value = fieldsList[i][j].op0Value;
						}
					}
					if (i > 0) {
						if (fieldsList[i][TEMP_INDEX].resultValue == '') {
							fieldsList[i][TEMP_INDEX].F4Value = fieldsList[i][TEMP_INDEX].op0Value;
							break
						}
						if (typeof fieldsList[i][j].op0Value === 'object') {
							fieldsList[i][j].op0Value = ''
							fieldsList[i][j].F4Value = fieldsList[i][j].op0Value;
						} else {
							fieldsList[i][j].F4Value = fieldsList[i][j].op0Value;
						}
						fieldsList[i][j].sessionStorage = true;
					}
				}
			}

			for (let key of Object.keys(this.memoFields)) {
				// const memoField = this.memoFields[key]
				// memoField
				this.memoFields[key].sessionStorage = true;
			}
			// console.log(fieldsList);
			// console.log(this.memoFields);
			// console.log(this.fieldsObject[this.thumbIndex].fieldsList)
		},

		// 获取临时保存的字段
		getSessionSaveField() {
			let { fieldsList, initFieldsList, sessionStorage } = this.fieldsObject[this.thumbIndex];
			// console.log(fieldsList);
			// console.log(this.fieldsObject[this.thumbIndex]);
			let arr = fieldsList[fieldsList.length - 1]
			// 未按F4，还原为初始状态
			if (!sessionStorage) {
				this.fieldsObject[this.thumbIndex].fieldsList = tools.deepClone(initFieldsList);
				return;
			}

			// 默认均为false
			{
				const flatFieldsList = tools.flatArray(fieldsList);

				flatFieldsList.map(field => {
					field.autofocus = false;
					field.uniqueKey = `enter_${field.uniqueId}_${Date.now()}`;
				});
			}

			// 过滤出field的sessionStorage字段为true的数据
			{
				let fieldsListArr = fieldsList.filter(fields => {
					return fields[TEMP_INDEX + 1].hasOwnProperty('sessionStorage') && fields[TEMP_INDEX].sessionStorage == true
				})
				if (arr.length == 3 && arr[TEMP_INDEX].resultValue == '') {
					fieldsListArr.push(arr)
				}
				// console.log('获取临时字段的----fieldsList-------', fieldsList);
				for (let i = 0; i < fieldsListArr.length; i++) {
					for (let j = 0; j < fieldsListArr[i].length; j++) {
						let x = fieldsListArr[i][j].F4Value
						if (fieldsListArr[i][j].name == '图片页码' || fieldsListArr[i][j].name == '显示范围') continue
						fieldsListArr[i][j].op0Value = x
						// this.$set(fieldsListArr[i][j], "resultValue", x);
					}
				}
				fieldsList = fieldsListArr
				this.$set(this.fieldsObject[this.thumbIndex], "fieldsList", fieldsListArr);
			}

			// 检查是否为按F4后录入的值
			{
				const length = fieldsList.length;

				for (let fieldsIndex = 1; fieldsIndex < length; fieldsIndex++) {
					const fields = fieldsList[fieldsIndex];
					const tempField = fields[TEMP_INDEX];

					if (!tempField.sessionSave) {
						this.focusFieldsIndex = fieldsIndex - 1;
						this.focusFieldIndex = TEMP_INDEX;

						// this.$set(fieldsList[this.focusFieldsIndex][this.focusFieldIndex], "autofocus", true);
						//修复问题BUG临时保存光标
						// 最后一个字段无内容
						if (fieldsList[length - 1][1].resultValue == "") {
							// let flag = fieldsList[length - 2].length - 1
							this.$set(fieldsList[length - 1][1], "autofocus", true);
							return;
						} else {
							// 最后一个字段有内容跳转至最后一个字段
							if (fieldsList[length - 1].length == 3) {
								this.$set(fieldsList[length - 1][1], "autofocus", true);
							} else {
								for (let i = 3; i < fieldsList[length - 1].length; i++) {
									if (fieldsList[length - 1][i].resultValue == "") {
										this.$set(fieldsList[length - 1][i - 1], "autofocus", true);
										return;
									}
								}
							}
						}

						if (fieldsList[length - 1][fieldsList[length - 1].length - 1].resultValue != '') {
							this.$set(fieldsList[length - 1][fieldsList[length - 1].length - 1], "autofocus", true);
							return;
						}

						// fieldsList.splice(fieldsIndex, 1)
						return;
					}
				}
			}

			if (this.fieldsObject[this.thumbIndex].fieldsList.length == 1 && this.fieldsObject[this.thumbIndex].fieldsList[0].length == 3) {
				this.$set(this.fieldsObject[this.thumbIndex].fieldsList[0][1], "autofocus", true);
				return;
			}

			// console.log(this.fieldsObject[this.thumbIndex].fieldsList);
			// 如果只有一个模板字符段光标停留在最后一个输入框
			if (this.fieldsObject[this.thumbIndex].fieldsList.length == 1) {
				const length = this.fieldsObject[this.thumbIndex].fieldsList[0].length;
				for (let i = 3; i < length; i++) {
					if (this.fieldsObject[this.thumbIndex].fieldsList[0][i].resultValue != "") {
						// console.log(6543231);
						// console.log(this.fieldsObject[this.thumbIndex].fieldsList[0][i]);
						// this.$set(this.fieldsObject[this.thumbIndex].fieldsList[0][i], "autofocus", true);
						let flag = this.fieldsObject[this.thumbIndex].fieldsList[0].length - 1
						this.$set(this.fieldsObject[this.thumbIndex].fieldsList[0][flag], "autofocus", true);
						return;
					} else {
						let flag = this.fieldsObject[this.thumbIndex].fieldsList[0].length - 1
						this.$set(this.fieldsObject[this.thumbIndex].fieldsList[0][flag], "autofocus", true);
						return
					}
				}
			}

			// console.log(this.fieldsObject[this.thumbIndex].fieldsList[0]);
		},

		// 绘制后的图片
		async handleDrewSave({ file, modified }) {
			this.FileLists = []
			for (let key in this.fieldsObject) {
				for (let fields of this.fieldsObject[key].fieldsList) {
					for (let field of fields) {
						if (field.resultValue != '' && field.name != '图片页码' && field.name != '显示范围' && field.hasOwnProperty('checkMethod')) {
							this.FileLists.push(field)
						}
					}
				}
			}
			let result = {};
			// 用户操作过图片
			// console.log(this.bill.downloadPath, this.bill.proCode, `${this.bill.billName}.${this.thumbIndex}.${this.focusFieldsIndex}.png`);
			if (tools.isYummy(file.name) && modified) {
				let data;
				if (this.bill.proCode === "B0108") {
					data = {
						file,
						path: this.bill.downloadPath,
						name: `${this.bill.billName}.${this.thumbIndex}.${this.focusFieldsIndex}.png`,
						op: this.op
					};
				} else if (this.bill.proCode === "B0114") {
					data = {
						file,
						path: this.bill.downloadPath,
						name: `${this.bill.billName}.${this.thumbIndex}.${this.focusFieldsIndex}.png`,
						op: this.op
					};
				} else if (this.bill.proCode === "B0118") {
					data = {
						file,
						path: this.bill.downloadPath,
						name: `${this.bill.billName}.${this.thumbIndex}.${this.focusFieldsIndex}.png`,
						op: this.op
					};
				} else {
					data = {
						file,
						path: this.bill.downloadPath.replace("files/", ""),
						name: `${this.bill.billName}.${this.thumbIndex}.${this.focusFieldsIndex}.png`,
						op: this.op
					};
				}

				result = await this.$store.dispatch("INPUT_SUBMIT_MODIFIED_IMAGE", data);

				if (result.code === 200) {
					let pass = await this.validateFieldss(this.FileLists, 'F4保存')
					if (pass) {
						this.modifyImage = result.data;

						// 显示范围(存储编辑后的图片)
						this.saveModifiedImage({ sessionStorage: true });
					}
				}
			}

			// 优化：不等待返回结果直接聚焦到下一模板
			this.autofocusToNextFields();

			// 防止手速过快对应字段没有出现就保存或提交了
			const { fieldsList } = this.fieldsObject[this.thumbIndex];
			let mdStr = fieldsList[0][1].valChange
			let mdObj = localStorage.get('mb001')
			let change1 = mdStr.split(';')
			let mapArr = []
			change1.forEach(el => {
				let x = el.split('=').reverse()
				mapArr.push(x)
			})
			// mapArr.pop()
			const map = new Map(mapArr)
			let objMap = {}
			for (let name in mdObj) {
				objMap[map.get(name)] = mdObj[name]
			}
			// console.log(objMap);
			objMap[''] = 3
			let flag = true
			fieldsList.forEach(el => {
				let length = el.length
				if (length != objMap[el[1].resultValue]) {
					flag = false
					return
				}
			})
			if (flag == true) {
				let pass = await this.validateFields(this.FileLists, 'F4保存')
				if (!pass) return
				this.limitSessionSave(result);
			} else return
		},

		// 临时保存(按F4则临时保存,否则不管是切换图片还是更新[模板类型字段]的值,均恢复为初始化)
		sessionSave({ code, msg }) {
			this.saveModifiedImage({ sessionStorage: true });

			// this.autofocusToNextFields();

			this.setSessionSaveField();

			if (code) {
				this.toasted.success(`${msg}，字段临时保存成功！`, toastedOptions);
			} else {
				this.toasted.success("字段临时保存成功！", toastedOptions);
			}
		},

		// 限制临时保存次数
		limitSessionSave(result) {
			// this.recordKey.push('按键:F4保存' + ',' + moment(new Date).format('YYYY/MM/DD HH:mm:ss') + ',' + this.recordF4)
			// console.log(this.recordKey);
			this.$store.commit('UPDATE_KEY', '按键:F4保存' + ',' + moment(new Date).format('YYYY/MM/DD HH:mm:ss') + ',' + this.recordF4)

			if (!this.oldTime) {
				this.sessionSave(result);
				this.oldTime = Date.now();
			} else {
				if (Date.now() - this.oldTime > 500) {
					this.sessionSave(result);
					this.oldTime = Date.now();
				}
			}
		},

		// 提交(F8)
		async submitTask(check = true) {
			let logResult = await this.$store.dispatch("INPUT_SUBMIT_TASK_KEY_OPERATION", { log: this.$store.state['recordKey'].recordKey, billNum: this.bill.billNum });
			console.log(logResult);
			// this.recordKey.push('按键:F8提交' + ',' + moment(new Date).format('YYYY/MM/DD HH:mm:ss') + ',' + this.recordF8)
			// console.log(this.recordKey);
			this.$store.commit('UPDATE_KEY', '按键:F8提交' + ',' + moment(new Date).format('YYYY/MM/DD HH:mm:ss') + ',' + this.recordF8)
			// console.log(this.$store.state['recordKey'].recordKey);
			// 释放单release
			this.release = 0
			// 当前 thumbIndex 下用户未按 F4，也需要保存并提交
			await this.setSessionSaveField();
			if (!this.fieldsObject[this.thumbIndex]?.sessionStorage && this.fieldsObject[this.thumbIndex].fieldsList.length != 0) {
				this.saveModifiedImage({ sessionStorage: true });
			}

			const mergeFieldsList = [];
			// console.log('F8提交this.fieldsObject----------', this.fieldsObject);
			// 模板字段 均为空或均有值可提交 其中有空不可提交 屏蔽字段除外
			// 记录每组restValue为空的个数
			// let count = 0;
			// let disabled = 0;
			// for (let fields of this.fieldsObject[0].fieldsList) {
			// 	// console.log(fields);
			// 	for (let field of fields) {
			// 		if (field.resultValue == "") count++;
			// 		if (field.disabled == true) disabled++;
			// 	}
			// 	if (count + 1 != fields.length) {
			// 		// console.log(count, disabled, fields.length);
			// 		if (count > disabled + 1) {
			// 			count = 0;
			// 			disabled = 0;

			// 			if (sessionStorage.get('isApp')?.isApp === 'true') {
			// 				await MessageBox.confirm("存在模板未屏蔽字段值为空，不可提交, 请检查并按要求录入后提交！！！", '请检查', {
			// 					type: 'warning',
			// 					confirmButtonText: '确定',
			// 					showCancelButton: false,
			// 					showClose: false,
			// 				}).then(() => {
			// 					return false;
			// 				})
			// 			} else {
			// 				alert("存在模板未屏蔽字段值为空，不可提交, 请检查并按要求录入后提交！！！");
			// 				return false;
			// 			}

			// 		}
			// 	}
			// 	count = 0;
			// 	disabled = 0;
			// }

			// for (let el in this.fieldsObject) {
			// 	for (let fields of this.fieldsObject[el].fieldsList) {
			// 		for (let field of fields) {
			// 			if (field.resultValue == "") count++;
			// 			if (field.disabled == true) disabled++;
			// 		}
			// 		if (count + 1 != fields.length) {
			// 			// console.log(count, disabled, fields.length);
			// 			if (count > disabled + 1) {
			// 				count = 0;
			// 				disabled = 0;
			// 				alert("存在模板未屏蔽字段值为空，不可提交, 请检查并按要求录入后提交！！！");
			// 				return false;
			// 			}
			// 		}
			// 	}
			// 	count = 0;
			// 	disabled = 0;
			// }

			//
			for (let key in this.fieldsObject) {
				// 若有做临时保存,则提交到服务器
				if (this.fieldsObject[key].sessionStorage) {
					// 若当前 fields 的 [模板类型字段] 不为空，则添加 fields
					this.fieldsObject[key].fieldsList.map(fields => {
						if (tools.isYummy(fields[TEMP_INDEX].op0Value)) {
							mergeFieldsList.push(fields);
						}
					});
				}
			}
			// 默认ID
			if (tools.isYummy(mergeFieldsList)) {
				this.fieldsList[0].map((field, index) => {
					mergeFieldsList[0][index].ID = field.ID;
				});
			}

			const data = {
				bill: this.bill,
				block: this.block,
				fields: nifty.deepClone(mergeFieldsList),
				op: this.op
			};
			// console.log("F8提交Data---------", data);
			for (let key in data.fields) {
				for (let field of data.fields[key]) {
					if (field.name == "模板类型字段" && field.resultValue.trim() == "") {
						data.fields[key] = [];
					}
					if (field.hasOwnProperty('F4Value') && field.F4Value != '') {
						field.resultValue = field.F4Value
					}
					if (field.hasOwnProperty('F4Value') && field.F4Value == '' && field.op0Value != '' && field.name != '图片页码' && field.name != '显示范围') {
						field.resultValue = field[`${this.op}Value`]
					}
				}
			}
			let flag = [];
			for (let field of data.fields) {
				if (field.length != 0) {
					flag.push(field);
				}
			}
			data.fields = flag;
			console.log("F8提交清除模板类型字段为空的数据data", data);
			// 删除前端设置的key
			for (let i = 0; i < data.fields.length; i++) {
				for (let field of data.fields[i]) {
					field.resultValue = field.resultValue.trim()
					field.op0Value = field.op0Value.trim()
					field.op1Value = field.op1Value.trim()
					field.op2Value = field.op2Value.trim()
					field.opqValue = field.opqValue.trim()
					delete field.F4Value;
					delete field.autofocus;
					delete field.desserts;
					delete field.disabled;
					delete field.effectValidations;
					delete field.items;
					delete field.rules;
					delete field.uniqueId;
					delete field.uniqueKey;
					delete field.table;
					if (i > 0) {
						delete field.ID;
					}
				}
			}

			if (check) {
				const svResult = await this.svValidateFields({
					fieldsObject: tools.deepClone(this.fieldsObject),
					mergeFieldsList: tools.deepClone(mergeFieldsList)
				});

				if (svResult) {
					// 开发环境默认不提交
					if (process.env.NODE_ENV === "development") {
						console.log("F8提交Data---------", data);
						return;
					}

					console.log("F8提交Data---------", data);
					let result
					if (this.downF8) {
						this.downF8 = false
						result = await this.$store.dispatch("INPUT_SUBMIT_TASK", data);
						let logResult = await this.$store.dispatch("INPUT_SUBMIT_TASK_KEY_OPERATION", { log: this.$store.state['recordKey'].recordKey, billNum: this.bill.billNum });
					} else {
						this.toasted.warning('正在提交中， 请勿重复按F8提交')
					}

					this.toasted.dynamic(result.msg, result.code, toastedOptions);

					if (result.code === 200) {
						this.downF8 = true
						this.getTask({ status: "new" });
						setTimeout(() => {
							this.navigatePage('Op0')
						}, 300);
					}

					this.$emit("submittedTaskResponse", { code: result.code });
				}
			}
		},

		// 释放单
		navigatePage(name) {
			const proCode = sessionStorage.get('proCode')

			if (proCode) {
				// 用于点击当前工序重新领取任务
				this.$router.replace({ name: "Channel", query: { proCode, op: -1 } });

				this.$nextTick(() => {
					this.$router.push({ name, query: { proCode } });
				});
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
			this.clearNVs = true;

			this.$nextTick(async () => {
				// this.svMemoFieldValues = {};

				// const valid = this.$refs.op0Form.validate()
				// console.log('this.fieldsObject-----F8------', this.fieldsObject);

				this.FileLists = []
				for (let key in this.fieldsObject) {
					for (let fields of this.fieldsObject[key].fieldsList) {
						for (let field of fields) {
							if (field.resultValue != '' && field.name != '图片页码' && field.name != '显示范围' && field.hasOwnProperty('checkMethod')) {
								this.FileLists.push(field)
							}
						}
					}
				}
				// console.log('this.FileLists------F8-----', this.FileLists);
				let pass = await this.validateFields(this.FileLists, 'F8提交')
				if (!pass) return
				let valid;
				if (this.$refs[`${this.op}Form`].validate.length == 0) {
					valid = true;
				} else {
					valid = this.$refs[`${this.op}Form`].validate();
				}

				valid && this.limitSubmitTask();

				this.clearNVs = false;
			});
		},

		// 设置当前编辑图片
		setModifyImage({ fieldsList }) {
			const rangField = fieldsList[this.focusFieldsIndex][RANG_INDEX];
			if (rangField) {
				if (this.flag && this.bill.imagesType && this.bill.imagesType[0]) {
					const { path } = this.fieldsObject[this.newfieldsObjectArray[0]];
					this.modifyImage = path;
					this.thumbIndex = this.newfieldsObjectArray[0];
				} else {
					const { path } = this.fieldsObject[this.thumbIndex];
					if (rangField.op0Value || this.modifyImage !== path) {
						this.modifyImage = rangField.op0Value || path;

						this.$set(fieldsList[this.focusFieldsIndex][rangField.fieldIndex], "op0Value", rangField.op0Value);
					}
				}
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
				this.fieldsObject[this.thumbIndex].fieldsList = tools.deepClone(fieldsList);
			}
		},

		// 获取当前编辑图片的大小
		async getImageSize() {
			const path = `${this.bill.downloadPath}${this.modifyImage}`;
			const result = await this.$store.dispatch("INPUT_GET_IMAGE_SIZE", { op: this.op, path });

			this.imageSize = result.data;
		},

		// 恢复初始化图片
		handleInitImage() {
			const { path } = this.fieldsObject[this.thumbIndex];
			const rangField = this.fieldsObject[this.thumbIndex].fieldsList[this.focusFieldsIndex][RANG_INDEX];

			this.modifyImage = path;

			rangField.op0Value = "";
			rangField.resultValue = "";
		},

		// 根据模板字段生成编辑后生成图片名
		setDrewImageName() {
			this.drewImageName = `${this.bill.billName}.${this.thumbIndex}.${this.focusFieldsIndex}.jpeg`;
		}
	},

	computed: {
		...mapGetters(["task"])
	}
};
