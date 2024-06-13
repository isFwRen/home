import { mapGetters } from "vuex";
import nifty from "nifty-util";
import { tools, localStorage, sessionStorage } from "vue-rocket";
import { toastedOptions } from "../../cells";
import { del } from "vue";

const defaultBill = {
	pictures: [],
	phoTotal: 0,
	pageIndex: 0
};

const defaultFocusFieldIndex = -1;

export default {
	data() {
		return {
			// block
			isLoop: false,

			// fieldsList
			fieldsListLength: 0,

			// fields
			tempFields: [],
			focusFieldsIndex: 0,
			focusFieldIndex: defaultFocusFieldIndex,

			// field
			prevGoodField: null,
			nextGoodField: null,

			// 考核题目id
			assessmentProblem: '',
			// n页数据
			pagesFields: [],
			// 总共多少页
			pages: 1,
			// 第几页
			page: 0,
			// 当前页对象
			fieldsListObj: [],
			// 当前页数据
			fieldsList: [],
			// 当前页图片
			fieldImg: '',

			// 题目
			answerList: [],

			// 保存要提交的数据
			submitData: {},

			// 考核结果
			examResult: {},

			// 错误提示
			errorIndex: null
		};
	},

	created() {
		this.getTask();
	},

	mounted() {
		window.addEventListener("keydown", this.fuckEvents);
	},

	beforeDestroy() {
		window.removeEventListener("keydown", this.fuckEvents);
	},

	methods: {
		// 重置数据
		resetData() {
			this.block = {};

			this.fieldsListLength = 0;

			this.fieldsList = [];
			this.tempFields = [];
			this.focusFieldsIndex = 0;
			this.focusFieldIndex = defaultFocusFieldIndex;

			this.prevGoodField = null;
			this.nextGoodField = null;


			// 考核题目id
			this.assessmentProblem = null;
			// n页数据
			this.pagesFields = [];
			// 总共多少页
			this.pages = 1;
			// 第几页
			this.page = 0;
			// 当前页对象
			this.fieldsListObj = [];
			// 当前页数据
			this.fieldsList = [];
			// 当前页图片
			this.fieldImg = null;
			// 小题
			this.answerList = []
			// 保存要提交的数据
			this.submitData = {}
			// 考核结果
			this.examResult = {}

			this.errorIndex = null
		},

		// 领任务
		async getTask() {
			this.resetData();
			let result = await this.$store.dispatch("GET_EXAM_TASK", this.$route.query.proCode);
			console.log('result', result);

			if (result.code === 200) {
				this.assessmentProblem = result.data.assessmentProblem
				this.submitData.assessmentProblemId = this.assessmentProblem
				this.pagesFields = result.data.assessmentBlockList
				this.pages = this.pagesFields.length
				this.fieldsListObj = this.pagesFields[this.page]
				this.fieldImg = this.fieldsListObj.blockImgUrl
				// this.fieldsList = [this.fieldsListObj.assessmentSingleList]
				this.setFieldsList([this.fieldsListObj.assessmentSingleList]);
			} else {
				this.toasted.warning(result.msg, toastedOptions);
			}

			this.$emit("examContent", {
				code: 200,
			});
			this.progress = (((this.page + 1) / this.pages) * 100).toFixed(2)
			console.log("this.fieldsList", this.fieldsList);
		},

		// fieldsList
		setFieldsList(fieldsList) {
			if (tools.isYummy(fieldsList)) {
				this.fieldsList = tools.deepClone(fieldsList);
			} else {
				this.fieldsList = [];
				return;
			}
			this.findAndMarkFirstLastGoodField();
			const { firstGoodField } = this.findFieldsFirstLastGoodField(this.fieldsList[0]);

			this.focusFieldIndex = tools.isYummy(firstGoodField?.fieldIndex)
				? firstGoodField.fieldIndex
				: defaultFocusFieldIndex;
		},

		// 找到第一个及最后一个[合法field]，并标记
		findAndMarkFirstLastGoodField() {
			this.fieldsList.map(fields => {
				let [firstGoodField, lastGoodField] = [null, null];
				// 默认不设置 firstGoodField lastGoodField
				fields.map((field, index) => {
					if (field.hasOwnProperty("isFirstGoodField")) delete field.isFirstGoodField;
					if (field.hasOwnProperty("isLastGoodField")) delete field.isLastGoodField;
					field.fieldIndex = index
				});

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
			})
			this.fieldsList = [...this.fieldsList];
		},

		// 找到当前 fields 的第一个及最后一个[合法field]
		findFieldsFirstLastGoodField(fields) {
			console.log(fields);
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
			// console.log('聚焦', event, field, fieldsIndex, fieldIndex);

			this.focusFieldsIndex = fieldsIndex;
			this.focusFieldIndex = fieldIndex;

			// 默认都不聚焦
			{
				const flatFieldsList = tools.flatArray(this.fieldsList);

				flatFieldsList.map(field => {
					field.autofocus = false;
				});
			}

			this.fieldsList[fieldsIndex][field.fieldIndex].autofocus = true

			this.fieldsList = [...this.fieldsList];
		},

		// 回车
		onEnterField({ event, field, fieldsIndex }) {
			// if (this.page < this.pages) this.nextPage()
			// else {
			// 	// fieldsList
			// 	this.fieldsListLength = this.fieldsList.length;

			// 	this.findAndMarkFirstLastGoodField();

			// 	const sameField = this.fieldsList[fieldsIndex][field.fieldIndex];

			// 	this.autofocusToNextField({ sameField, fieldsIndex });
			// }

			// console.log("field", field);
		},

		// 用户输入值
		async onInputField(value, field, fieldsIndex, fieldIndex) {

			this.fieldsList[fieldsIndex][fieldIndex].optionList[0] = value;

			this.fieldsList = [...this.fieldsList];
			// console.log(this.fieldsList);
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
		onDnKey({ event, field, fieldsIndex }) {
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

		// radio
		radioChange(value, fieldsIndex, fieldIndex) {
			// console.log('radioChange', value, fieldsIndex, fieldIndex);
			this.fieldsList[fieldsIndex][fieldIndex].optionLists = value;

			this.fieldsList = [...this.fieldsList];
			// console.log(this.fieldsList);
		},

		// checkBox
		boxChange(value, fieldsIndex, fieldIndex) {
			// console.log('boxChange', value, fieldsIndex, fieldIndex);
			this.fieldsList[fieldsIndex][fieldIndex].optionLists = value;

			this.fieldsList = [...this.fieldsList];
			// console.log(this.fieldsList);
		},

		// 下一页
		nextPage() {
			this.errorIndex = null
			for (let el of this.fieldsList[0]) {
				if (el.problemType != 1 && el.isRequire == 1) {
					if (!el.hasOwnProperty('optionLists') || el.optionLists.length == 0) {
						this.errorIndex = el.fieldIndex
						return alert(`第${el.fieldIndex + 1}题，${el.question}未填写，该项为必填，请返回填写后提交`)
					}
				} else {
					if (el.optionList.length == 0) {
						this.errorIndex = el.fieldIndex
						return alert(`第${el.fieldIndex + 1}题，${el.question}未填写，该项为必填，请返回填写后提交`)
					}
				}
			}
			for (let el of this.fieldsList[0]) {
				delete el.autofocus
				delete el.fieldIndex
				delete el.isFirstGoodField
				delete el.isLastGoodField
				delete el.orderNum
				delete el.question
				if (el.hasOwnProperty('optionLists')) {
					el.optionList = el.optionLists
				}
				delete el.optionLists
				this.answerList = [...this.answerList, el]
			}

			this.page++
			this.fieldsListObj = this.pagesFields[this.page]
			this.fieldImg = this.fieldsListObj.blockImgUrl
			// this.fieldsList = [this.fieldsListObj.assessmentSingleList]
			this.setFieldsList([this.fieldsListObj.assessmentSingleList]);
			this.submitData.answerList = this.answerList
			this.progress = (((this.page + 1) / this.pages) * 100).toFixed(2)
			console.log("this.fieldsList", this.fieldsList);
			// console.log('this.assessmentProblem', this.assessmentProblem);
			// console.log('this.page', this.page);
			// console.log('this.submitData', this.submitData);
		},

		// 提交当前分块
		async submitTask() {
			this.errorIndex = null
			for (let el of this.fieldsList[0]) {
				if (el.isRequire && el.optionList.length == 0) {
					this.errorIndex = el.fieldIndex
					return alert(`第${el.fieldIndex + 1}题，${el.question}未填写，该项为必填，请返回填写后提交`)
				}
			}
			for (let el of this.fieldsList[0]) {
				delete el.autofocus
				// delete el.fieldIndex
				delete el.isFirstGoodField
				delete el.isLastGoodField
				delete el.orderNum
				// delete el.question
				if (el.hasOwnProperty('optionLists')) {
					el.optionList = el.optionLists
				}
				delete el.optionLists
				this.answerList = [...this.answerList, el]
			}

			this.submitData.answerList = this.answerList
			console.log('考核最终提交的数据------', this.submitData);
			const result = await this.$store.dispatch("SUBMIT_EXAM_TASK", this.submitData)
			if (result.code == 200) {
				this.examResult = result.data
				this.dialog = true;
			}
		},

		// 按F8提交
		async fEightSubmit() {
			let valid = this.$refs[`${this.op}Form`].validate();
			valid && this.limitSubmitTask();
		},
	},
};
