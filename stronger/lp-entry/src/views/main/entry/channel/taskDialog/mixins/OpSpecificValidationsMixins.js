import { tools, sessionStorage } from 'vue-rocket'
import _ from 'lodash'
import { change } from './DropCells'
import { reqConst } from '@/api/reqConst'
import commonFunctionOp12q from './commonFun'
import { MessageBox } from 'element-ui';

export default {
  data() {
    return {
      proCode: void 0,

      specificProject: null,

      // 记录最后一次存储的合法field
      svMemoFields: {},

      // 记录当前字段的所有值
      svMemoFieldValues: {},

      // 以某些字段的值为下拉
      svDropdownFields: {},

      // 记录当前字段录入的所有值(包括重复的值)
      sameFieldValue: {},

      fieldFirst: '',
      firstObject: {},
      result: false,

      clearValues: {},
      clearFieldValues: {},
    }
  },

  created() {
    this.proCode = this.$route.query?.proCode
    const files = require.context('../specificValidations', true, /index\.js$/)

    files.keys().map(key => {
      if (key.includes(this.proCode)) {
        this.specificProject = { ...files(key).default }
      }
    })
  },

  methods: {
    // 重置变量
    svResetVariable() {
      this.svMemoFields = {}
      this.svMemoFieldValues = {}
      this.svDropdownFields = {}
      // 记录字段录入的所有值
      this.sameFieldValue = {}
      this.clearValues = {}
      this.clearFieldValues = {}
    },

    // 记录最后一次存储的合法field 
    svUpdateMemoFields({ value, code, items }) {
      const { op0, op1op2opq } = this.specificProject
      const memoFields = []

      if (this.op === 'op0') {
        memoFields.push(...op0?.memoFields)
      }
      else {
        memoFields.push(...op1op2opq?.memoFields)
      }

      // 记录
      if (memoFields.includes(code)) {
        _.set(this.svMemoFields, `${code}.value`, value)
        _.set(this.svMemoFields, `${code}.items`, items)
      }
    },

    /**
     * @description 记录相同 code 的 field 的值
     * @param { String } code 字段编码
     */
    svUpdateMemoFieldValues({ code }) {
      const { op0, op1op2opq } = this.specificProject
      const memoFieldValues = []

      if (this.op === 'op0') {
        memoFieldValues.push(...op0?.memoFieldValues)
      }
      else {
        memoFieldValues.push(...op1op2opq?.memoFieldValues)
      }

      // 记录
      if (memoFieldValues.includes(code)) {
        const values = []
        const fieldValues = []

        const clearValues = []
        const clearFieldValues = []

        for (let key in this.fieldsObject) {
          const fieldsList = this.fieldsObject[key].fieldsList
          const sessionStorage = this.fieldsObject[key].sessionStorage
          // 临时存储或为当前页面
          if (sessionStorage || this.thumbIndex === +key) {
            for (let fields of fieldsList) {
              for (let field of fields) {
                if (code === field.code) {
                  const value = field[`${this.op}Value`]
                  !values.includes(value) && values.push(value)
                }
              }
            }
          }

          if (sessionStorage || this.thumbIndex === +key) {
            for (let fields of fieldsList) {
              for (let field of fields) {
                if (code === field.code) {
                  const value = field[`${this.op}Value`]
                  fieldValues.push(value)
                }
              }
            }
          }

          // 清除未F4保存的
          if (!sessionStorage && this.thumbIndex === +key) {
            for (let fields of fieldsList) {
              for (let field of fields) {
                if (code === field.code) {
                  const value = field[`${this.op}Value`]
                  !clearValues.includes(value) && clearValues.push(value)
                }
              }
            }
          }

          if (!sessionStorage && this.thumbIndex === +key) {
            for (let fields of fieldsList) {
              for (let field of fields) {
                if (code === field.code) {
                  const value = field[`${this.op}Value`]
                  clearFieldValues.push(value)
                }
              }
            }
          }
        }

        _.set(this.svMemoFieldValues, `${code}.values`, values)
        _.set(this.sameFieldValue, `${code}.values`, fieldValues)

        _.set(this.clearValues, `${code}.values`, clearValues)
        _.set(this.clearFieldValues, `${code}.values`, clearFieldValues)

        // console.log('this.svMemoFieldValues', this.svMemoFieldValues);
        // console.log('this.sameFieldValue', this.sameFieldValue);
      }
    },

    f4SetFieldValues() {
      const { op0, op1op2opq } = this.specificProject
      const memoFieldValues = []

      if (this.op === 'op0') {
        memoFieldValues.push(...op0?.memoFieldValues)
      }
      else {
        memoFieldValues.push(...op1op2opq?.memoFieldValues)
      }

      // 记录
      this.svMemoFieldValues = {}
      this.sameFieldValue = {}

      for (let el of memoFieldValues) {
        _.set(this.svMemoFieldValues, `${el}.values`, [])
        _.set(this.sameFieldValue, `${el}.values`, [])
      }

      for (let key in this.fieldsObject) {
        const fieldsList = this.fieldsObject[key].fieldsList
        // 临时存储或为当前页面
        for (let fields of fieldsList) {
          for (let field of fields) {
            if (memoFieldValues.includes(field.code)) {
              const value = field[`${this.op}Value`]
              !this.svMemoFieldValues[field.code].values.includes(value) && _.set(this.svMemoFieldValues, `${field.code}.values`, [...this.svMemoFieldValues[field.code].values, value])
              _.set(this.sameFieldValue, `${field.code}.values`, [...this.sameFieldValue[field.code].values, value])
            }
          }
        }
      }

      // console.log('this.svMemoFieldValues---f4', this.svMemoFieldValues);
      // console.log('this.sameFieldValue---f4', this.sameFieldValue);
    },

    // 以某个 field 的所有值为下拉
    svUpdateDropdown() {
      const { op0, op1op2opq } = this.specificProject
      const dropdownFields = []

      if (this.op === 'op0') {
        dropdownFields.push(...op0?.dropdownFields)
      }
      else {
        dropdownFields.push(...op1op2opq?.dropdownFields)
      }

      for (let record of dropdownFields) {
        const [fields, items, flatFieldsList, flatInitFieldsList] = [[], [], [], []]

        for (let key in this.fieldsObject) {
          const sessionStorage = this.fieldsObject[key].sessionStorage
          flatFieldsList.push(...tools.flatArray(this.fieldsObject[key].fieldsList))
          flatInitFieldsList.push(...tools.flatArray(this.fieldsObject[key].initFieldsList))

          // 临时存储或为当前页面
          if (sessionStorage || this.thumbIndex === +key) {
            flatFieldsList.map(field => {
              if (record.targets.includes(field.code)) {
                fields.push(tools.deepClone(field))
              }
            })
          }
          else {
            flatInitFieldsList.map(field => {
              if (record.targets.includes(field.code)) {
                fields.push(tools.deepClone(field))
              }
            })
          }
        }

        // 找到需要将值设置为 dropdown 的 field
        fields.map(field => {
          if (record.targets.includes(field.code)) {
            if (!items.includes(field[`${this.op}Value`]) && field[`${this.op}Value`]) {
              items.push(field[`${this.op}Value`])
            }
          }
        })

        // 给符合条件的 field 设置 dropdown
        record.fields.map(fieldCode => {
          _.set(this.svDropdownFields, `${fieldCode}.desserts`, items)
        })
      }
    },

    // 工序完成初始化 
    async svInit() {
      const { op0, op1op2opq } = this.specificProject
      let methods = {}

      if (this.op === 'op0') {
        methods = { ...methods, ...op0?.init?.methods }
      }
      else {
        methods = { ...methods, ...op1op2opq?.init?.methods }
      }

      // 执行函数
      if (tools.isYummy(methods)) {
        for (let funcName in methods) {
          await methods[funcName]({
            op: this.op,
            bill: this.bill,
            block: this.block
          })
        }
      }
    },

    // 工序初始化ocr首次加载提示
    async ocrPrompt({ fieldsList }) {
      for (let fields of fieldsList) {
        for (let field of fields) {
          if (field.resultValue) {
            await this.svSearchConstants({ value: field.resultValue, field })
            await this.hintFc({ field })
          }
        }
      }
    },
    // 字段已生成
    async svUpdateFields({ fieldsList }) {
      const { op0, op1op2opq } = this.specificProject
      let methods = {}

      if (this.op === 'op0') {
        methods = { ...methods, ...op0?.updateFields?.methods }
      }
      else {
        methods = { ...methods, ...op1op2opq?.updateFields.methods }
      }

      // 执行函数
      if (tools.isYummy(methods)) {
        for (let funcName in methods) {
          methods[funcName]({
            op: this.op,
            bill: this.bill,
            block: this.block,
            codeValues: this.codeValues,
            focusFieldsIndex: this.focusFieldsIndex,
            fieldsList,
            flatFieldsList: tools.flatArray(fieldsList),
          })
        }
      }
      if (this.proCode == 'B0114') this.ocrPrompt({ fieldsList })
    },

    // 回车后操作字段
    svEnterUpdateField({ field, fieldsList, focusFieldsIndex, memoFields }) {
      const { op0, op1op2opq } = this.specificProject
      let methods = {}
      if (this.op === 'op0') {
        methods = { ...methods, ...op0.enter?.methods }
      }
      else {
        methods = { ...methods, ...op1op2opq.enter?.methods }
      }
      this.fieldFirst = fieldsList
      let flag
      // 执行函数
      if (tools.isYummy(methods)) {
        for (let funcName in methods) {
          let result = methods[funcName]({
            op: this.op,
            bill: this.bill,
            block: this.block,
            codeValues: this.codeValues,
            field,
            fieldsList,
            focusFieldsIndex,
            memoFields,
            sameCodeValues: this.svMemoFieldValues,
            sameFieldValue: this.sameFieldValue,
            flatFieldsList: tools.flatArray(fieldsList)
          })
          if (result == false) flag = false
        }
      }
      return flag
    },

    svDisableFields() {
      const { op0, op1op2opq } = this.specificProject
      let methods = {}
      if (this.op === 'op0') {
        methods = { ...methods, ...op0.sessionSave?.methods }
      }
      else {
        methods = { ...methods, ...op1op2opq.sessionSave?.methods }
      }

      // 执行禁用函数
      if (tools.isYummy(methods)) {
        for (let funcName in methods) {
          methods[funcName]({
            op: this.op,
            bill: this.bill,
            block: this.block,
            field: this.fieldsList[this.focusFieldsIndex][this.focusFieldIndex],
            fieldsList: this.fieldsList,
            focusFieldIndex: this.focusFieldIndex,
            focusFieldsIndex: this.focusFieldsIndex,
            flatFieldsList: tools.flatArray(this.fieldsList)
          })
        }
      }
    },

    // 初审通过聚焦拿到field，fieldsIndex, fieldIndex的值
    getFieldFirst(field, fieldsIndex, fieldIndex) {
      this.firstObject = { field, fieldsIndex, fieldIndex }
    },
    // 初审按F4 108_F4校验弹框需求
    async svDisableFieldsFirst() {
      const { op0 } = this.specificProject
      let methods = {}
      const errors = []
      if (this.op === 'op0') {
        methods = { ...methods, ...op0.sessionSave?.methods }
      }

      // 执行禁用函数
      if (tools.isYummy(methods)) {
        for (let funcName in methods) {
          const result = methods[funcName]({
            bill: this.bill,
            field: this.firstObject.field,
            fieldsIndex: this.firstObject.fieldsIndex,
            fieldIndex: this.firstObject.fieldIndex,
            fieldsList: this.fieldFirst,
            sameFieldValue: this.sameFieldValue,
          })

          if (result !== true) {
            errors.push(result)
          }
        }
      }

      if (!errors.length) {
        return true
      }
      const [normalErrors, confirmErrors] = [[], []]

      for (let error of errors) {
        if (error === false) {
          return false
        }
        if (!error) continue
        else if (!error?.popup) {
          normalErrors.push(error)
        }
        else if (error?.popup === 'confirm') {
          confirmErrors.push(error)
        }
      }

      const assignErrors = [...normalErrors, ...confirmErrors]

      console.log(assignErrors)
      if (!assignErrors[0]) return true

      for (let error of assignErrors) {
        const { errorMessage, popup } = error

        if (popup === 'confirm') {

          // if (!result) {
          //   // this.submitTask(false)
          //   return true
          // }
          if (sessionStorage.get('isApp')?.isApp === 'true') {
            await MessageBox.confirm(errorMessage, '请检查(ESC提交 回车返回修改)', {
              confirmButtonText: '修改',
              cancelButtonText: '提交',
              type: 'warning',
              // closeOnPressEscape: false,
              showClose: false,
              closeOnClickModal: false,
            }).then(() => {
              this.toasted.warning("请认真核对后修改", 3000);
            }).catch(() => {
              this.result = true
            });
            if (this.result) {
              // this.submitTask(false)
              this.result = false
              return true
            }
          } else {
            const result = confirm(errorMessage + '，-----确定按钮返回修改，取消按钮提交内容')

            if (!result) {
              // this.submitTask(false)
              // result = false
              return true
            }
          }
        }
        else {
          if (sessionStorage.get('isApp')?.isApp === 'true') {
            await MessageBox.alert(errorMessage, '请检查', {
              type: 'warning',
              confirmButtonText: '确定',
              showClose: false,
            })
          } else {
            alert(errorMessage)
          }

        }
      }
      this.cursorPosition()
      return false
    },

    /**
     * @description 提交前校验(F8)
     * @param fieldsObject 来自初审
     * @param mergeFieldsList 来自初审
     */
    async svValidateFields({ fieldsObject, mergeFieldsList } = {}) {
      const { op0, op1op2opq } = this.specificProject
      let methods = {}
      const errors = []
      if (this.op === 'op0') {
        methods = { ...methods, ...op0.beforeSubmit?.methods }
      }
      else {
        methods = { ...methods, ...op1op2opq.beforeSubmit?.methods, ...commonFunctionOp12q.methods }
      }

      // 执行校验函数
      if (tools.isYummy(methods)) {
        for (let funcName in methods) {
          const result = methods[funcName]({
            op: this.op,
            bill: this.bill,
            block: this.block,
            fieldsObject,
            mergeFieldsList,
            sameCodeValues: this.svMemoFieldValues,
            sameFieldValue: this.sameFieldValue,
            fieldsList: this.fieldsList,
            flatFieldsList: tools.flatArray(this.fieldsList)
          })

          if (result !== true) {
            errors.push(result)
          }
        }
      }

      if (!errors.length) {
        return true
      }
      const [normalErrors, confirmErrors] = [[], []]

      for (let error of errors) {
        if (error === false) {
          return false
        }
        if (!error) continue
        else if (!error?.popup) {
          normalErrors.push(error)
        }
        else if (error?.popup === 'confirm') {
          confirmErrors.push(error)
        }
      }

      const assignErrors = [...normalErrors, ...confirmErrors]

      if (!assignErrors[0]) return true

      for (let error of assignErrors) {
        const { errorMessage, popup } = error

        if (popup === 'confirm') {

          if (sessionStorage.get('isApp')?.isApp === 'true') {
            await MessageBox.confirm(errorMessage, '请检查(ESC提交 回车返回修改)', {
              confirmButtonText: '修改',
              cancelButtonText: '提交',
              type: 'warning',
              // closeOnPressEscape: false,
              showClose: false,
              closeOnClickModal: false,
            }).then(() => {
              this.toasted.warning("请认真核对后修改", 3000);
            }).catch(() => {
              this.result = true
            });
            if (this.result) {
              // this.submitTask(false)
              this.result = false
              return true
            }
          } else {
            const result = confirm(errorMessage + '，-----确定按钮返回修改，取消按钮提交内容')
            if (!result) {
              // this.submitTask(false)
              // result = false
              return true
            }
          }
        }
        else {
          if (sessionStorage.get('isApp')?.isApp === 'true') {
            await MessageBox.alert(errorMessage, '请检查', {
              type: 'warning',
              confirmButtonText: '确定',
              showClose: false,
            })
          } else {
            alert(errorMessage)
          }
        }
      }
      this.cursorPosition()

      return false
    },

    // 对话框光标自动定位
    cursorPosition() {
      // 对话框提示光标定位最后一个字段
      if (this.op == 'op0') {
        let { fieldsList } = this.fieldsObject[this.thumbIndex];
        // 默认均为false
        {
          const flatFieldsList = tools.flatArray(fieldsList);

          flatFieldsList.map(field => {
            field.autofocus = false;
            field.uniqueKey = `enter_${field.uniqueId}_${Date.now()}`;
          });
        }
        const length = fieldsList.length;

        if (fieldsList[length - 1].length == 3) {
          this.$set(fieldsList[length - 1][fieldsList[length - 1].length - 2], "autofocus", true);
        } else {
          this.$set(fieldsList[length - 1][fieldsList[length - 1].length - 1], "autofocus", true);
        }
      }
    },

    // 从常量库查询匹配字段
    async svSearchConstants({ value, field }) {
      // 存放单个常量表
      let items = []
      // 存放两个常量表
      let itemsTwo = []
      // 存放需要查询常量的Ocr
      let ocrBlur = []
      let ocrBlur1 = []
      // 需要查询常量
      if (field.table) {
        const { name: tableName, query, targets } = field.table
        const tableInfo = window.constantsDB[this.proCode]?.[tableName]
        // 当前常量存在
        if (tableInfo) {
          const targetIndexes = []

          // 修改问题：优化初审录入数据卡顿，数据多浏览器爆掉-----暂时解决---后续bug未知(B0114出现无法回车--已解决)
          // field.desserts = tableInfo.desserts
          change(tableInfo.desserts)

          if (targets) {
            tableInfo.headers.map((header, headerIndex) => {
              if (targets.includes(header)) {
                targetIndexes.push(headerIndex)
              }
            })
          }

          // 获取过滤名称下标
          const queryIndex = tableInfo.headers.indexOf(query)
          if (value) {
            if (!targets) {
              for (let dessert of tableInfo.desserts) {
                const text = dessert?.[queryIndex] || ''
                if (text.includes(value)) {
                  if (this.proCode == 'B0122') {
                    // if (!items.includes(text)) {
                    items.push(text?.trim())
                    itemsTwo.push(text?.trim())
                    // }
                  } else {
                    if (!items.includes(text)) {
                      items.push(text?.trim())
                      itemsTwo.push(text?.trim())
                    }
                  }

                }
              }
            }
            else {
              for (let dessert of tableInfo.desserts) {
                let texts = ''
                // dessert.map((text, textIndex) => {
                //   if (targetIndexes.includes(textIndex)) {
                //     const lastTargetIndex = targetIndexes[targetIndexes.length - 1]

                //     texts += textIndex !== lastTargetIndex ? `${text}-` : text
                //   }
                // })

                const text = dessert?.[queryIndex] || ''
                if (text.includes(value)) {
                  dessert.map((text, textIndex) => {
                    if (targetIndexes.includes(textIndex)) {
                      const lastTargetIndex = targetIndexes[targetIndexes.length - 1]

                      texts += textIndex !== lastTargetIndex ? `${text}-` : text
                    }
                  })
                }

                if (texts.includes(value)) {
                  if (!items.includes(texts)) {
                    items.push(texts?.trim())
                    itemsTwo.push(texts?.trim())
                  }
                }
              }
            }
          }
          else {
            items = []
          }

          // 存放需要查询常量的Ocr
          let str1 = value.slice(0, 2)
          let str2 = value.slice(2, 4)
          // 分两次模糊查找 -----可优化 B0114 代码编码：CSB0114RC0213000
          if (str1) {
            if (!targets) {
              for (let dessert of tableInfo.desserts) {
                const text = dessert?.[queryIndex] || ''
                if (text.includes(str1)) {
                  if (!ocrBlur.includes(text)) {
                    ocrBlur.push(text?.trim())
                  }
                }
              }
            }
            else {
              console.log(tableInfo.desserts);
              for (let dessert of tableInfo.desserts) {
                let texts = ''
                dessert.map((text, textIndex) => {
                  if (targetIndexes.includes(textIndex)) {
                    const lastTargetIndex = targetIndexes[targetIndexes.length - 1]

                    texts += textIndex !== lastTargetIndex ? `${text}-` : text
                  }
                })

                if (texts.includes(str1) && texts.length == value.length) {
                  if (!ocrBlur.includes(texts)) {
                    ocrBlur.push(texts?.trim())
                  }
                }
              }
            }
          }
          else {
            ocrBlur = []
          }

          if (str2) {
            if (!targets) {
              for (let dessert of tableInfo.desserts) {
                const text = dessert?.[queryIndex] || ''
                if (text.includes(str2)) {
                  if (!ocrBlur1.includes(text)) {
                    ocrBlur1.push(text?.trim())
                  }
                }
              }
            }
            else {
              for (let dessert of tableInfo.desserts) {
                let texts = ''
                dessert.map((text, textIndex) => {
                  if (targetIndexes.includes(textIndex)) {
                    const lastTargetIndex = targetIndexes[targetIndexes.length - 1]

                    texts += textIndex !== lastTargetIndex ? `${text}-` : text
                  }
                })

                if (texts.includes(str2) && texts.length == value.length) {
                  if (!ocrBlur1.includes(texts)) {
                    ocrBlur1.push(texts?.trim())
                  }
                }
              }
              console.log('ocrBlur1---------', ocrBlur1);
            }
          }
          else {
            ocrBlur1 = []
          }

          items.sort((a, b) => a.length - b.length);

          let selectArr = sessionStorage.get('select')
          if (selectArr.length != 0) {
            for (let el of selectArr) {
              if (el != null && el.includes(value) && items.includes(el)) {
                let flag = items.indexOf(el)
                items.splice(flag, 1)
                items?.unshift(el)
                ocrBlur.splice(flag, 1)
                ocrBlur1.splice(flag, 1)
                ocrBlur?.unshift(el)
                ocrBlur1?.unshift(el)
              }
            }
          }
          field.items = items?.slice(0, 100)
          if (ocrBlur.length != 0) field.ocrBlur = ocrBlur?.slice(0, 100)
          else field.ocrBlur = ocrBlur1?.slice(0, 100)
        }
      }
      // 查询多个常量 ---代码可优化
      if (field.tables) {
        const { name: tableName, query, targets } = field.tables
        const tableInfo = window.constantsDB[this.proCode]?.[tableName]
        // 当前常量存在
        if (tableInfo) {
          const targetIndexes = []

          // 修改问题：优化初审录入数据卡顿，数据多浏览器爆掉-----暂时解决---后续bug未知(B0114出现无法回车--已解决)
          // field.desserts = tableInfo.desserts
          change(tableInfo.desserts)

          if (targets) {
            tableInfo.headers.map((header, headerIndex) => {
              if (targets.includes(header)) {
                targetIndexes.push(headerIndex)
              }
            })
          }

          // 获取过滤名称下标
          const queryIndex = tableInfo.headers.indexOf(query)
          if (value) {
            if (!targets) {
              for (let dessert of tableInfo.desserts) {
                const text = dessert?.[queryIndex] || ''
                if (text.includes(value)) {
                  if (!items.includes(text)) {
                    itemsTwo.push(text?.trim())
                  }
                }
              }
            }
            else {
              for (let dessert of tableInfo.desserts) {
                let texts = ''
                dessert.map((text, textIndex) => {
                  if (targetIndexes.includes(textIndex)) {
                    const lastTargetIndex = targetIndexes[targetIndexes.length - 1]

                    texts += textIndex !== lastTargetIndex ? `${text}-` : text
                  }
                })

                if (texts.includes(value)) {
                  if (!items.includes(texts)) {
                    itemsTwo.push(texts?.trim())
                  }
                }
              }
            }
          }
          else {
            items = []
          }

          itemsTwo.sort((a, b) => a.length - b.length);

          let selectArr = sessionStorage.get('select')
          if (selectArr.length != 0) {
            for (let el of selectArr) {
              if (el != null && el.includes(value) && itemsTwo.includes(el)) {
                let flag = itemsTwo.indexOf(el)
                itemsTwo.splice(flag, 1)
                itemsTwo?.unshift(el)
              }
            }
          }

          field.items = itemsTwo?.slice(0, 50)
        }
      }
      // 当前字段的值从某个字段的所有值中选择
      if (this.svDropdownFields?.[field.code]) {
        const DD = this.svDropdownFields?.[field.code]

        if (value) {
          for (let text of DD.desserts) {
            if (text.includes(value) && !items.includes(text)) {
              items.push(text?.trim())
            }
          }
        }
        else {
          items = []
        }
        // console.log(value);
        // console.log(field.resultValue);

        let selectArr = sessionStorage.get('select')
        if (selectArr.length != 0) {
          for (let el of selectArr) {
            if (el != null && el.includes(value) && items.includes(el)) {
              let flag = items.indexOf(el)
              items.splice(flag, 1)
              items?.unshift(el)
            }
          }
        }
        field.items = items.slice(0, 50)
      }

      field.sameFieldValue = this.sameFieldValue
    },

    // OCR首次加载提示 
    async hintFc({ field }) {
      const codes = ['fc138', 'fc142', 'fc146', 'fc150', 'fc154', 'fc158', 'fc162', 'fc166']
      if (codes.includes(field.code)) {
        field.allowForce = true
        if (field.items) {
          const result = field.items.find(text => text === field.resultValue)
          if (result) return true
        }
        if (field.ocrBlur && field.ocrBlur?.length != 0) {
          const result = field.ocrBlur.find(text => text === field.resultValue)
          console.log('ocr识别的---------', field.resultValue);
          console.log('常量模糊0-100条---------', field.ocrBlur);
          if (!result) {
            // let dropItems = field.ocrBlur.slice(0, 2)
            let showItems = []
            let showItems1 = []

            for (let el of field.ocrBlur) {
              let valueArr = [...field.resultValue]
              let elArr = [...el]
              // 记录第一个不相同下标
              let flag1 = -1
              // 记录第二个不相同下标
              let flag2 = -1
              // 记录第三个不相同下标
              let flag3 = -1
              // 记录不相同个数
              let count = 0
              for (let i = 0; i < elArr.length; i++) {
                if (elArr[i] != valueArr[i]) {
                  if (count == 0) flag1 = i
                  if (count == 1) flag2 = i
                  if (count == 2) flag3 = i
                  count++
                }
              }

              if (count <= 3) {
                // console.log('count',count);
                // console.log('flag1',flag1, 'flag2',flag2, 'flag3',flag3);
                if (flag1 != -1 && flag2 == -1 && flag3 == -1) {
                  let front1 = el.slice(0, flag1)
                  let behind1 = el.slice(flag1 + 1)
                  let str = `${front1}<span style="color: red;">${elArr[flag1]}</span>${behind1}`
                  showItems.push(str)
                  showItems1.push(str)
                }
                if (flag1 != -1 && flag2 != -1 && flag3 == -1) {
                  let front1 = el.slice(0, flag1)
                  let behind1 = el.slice(flag1 + 1, flag2)
                  let behind2 = el.slice(flag2 + 1)
                  let str = `${front1}<span style="color: red;">${elArr[flag1]}</span>${behind1}<span style="color: red;">${elArr[flag2]}</span>${behind2}`
                  showItems.push(str)
                }
                if (flag1 != -1 && flag2 != -1 && flag3 != -1) {
                  let front1 = el.slice(0, flag1)
                  let behind1 = el.slice(flag1 + 1, flag2)
                  let behind2 = el.slice(flag2 + 1, flag3)
                  let behind3 = el.slice(flag3 + 1)
                  let str = `${front1}<span style="color: red;">${elArr[flag1]}</span>${behind1}<span style="color: red;">${elArr[flag2]}</span>${behind2}<span style="color: red;">${elArr[flag3]}</span>${behind3}`
                  showItems.push(str)
                }
              }
            }
            if (showItems1.length >= 2) showItems = showItems1
            if (!showItems[0]) {
              field.hint = ''
              return true
            }

            if (showItems[1]) {
              field.hint = `<p style="color: blue; fontSize: 14px; margin-top: -3px; margin-bottom: 0px"">${showItems[0]}，${showItems[1]}</p>`
            } else if (!showItems[1]) {
              field.hint = `<p style="color: blue; fontSize: 14px; margin-top: -3px; margin-bottom: 0px"">${showItems[0]}</p>`
            }
          }
        }

      }
    },


    // 请求客户端常量库
    async requestDropFields({ value, field }) {
      let items = []
      if (field.table) {
        const { name, query, targets } = field.table
        var regexStr = "/" + value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&') + "/"; // 转换为正则表达式字符串
        const data = {
          proCode: this.proCode,
          name,
          queryNames: {
            [field.table.query]: {
              $regex: regexStr
            }
          },
          respNames: [
            query
          ],
          pageSize: 50,
          pageIndex: 1
        }

        if (targets) {
          data.respNames = targets
        }

        const result = await reqConst({
          url: '/sys-const/page',
          method: "POST",
          data,
        })

        let fieldContent = result.list
        if (targets) {
          let fieldItems
          if (targets.length == 2) {
            fieldItems = fieldContent.map(item => {
              return item[targets[0]] + '-' + item[targets[1]]
            })
          }
          if (targets.length == 3) {
            fieldItems = fieldContent.map(item => {
              return item[targets[0]] + '-' + item[targets[1]] + '-' + item[targets[2]]
            })
          }
          field.items = fieldItems
        } else {
          let fieldItems = fieldContent.map(item => {
            return item[query]
          })
          fieldItems.sort((a, b) => a.length - b.length);
          field.items = fieldItems
        }
      }

      // 当前字段的值从某个字段的所有值中选择
      // console.log(this.svDropdownFields);
      if (this.svDropdownFields?.[field.code]) {
        const DD = this.svDropdownFields?.[field.code]

        if (value) {
          for (let text of DD.desserts) {
            if (text.includes(value) && !items.includes(text)) {
              items.push(text?.trim())
            }
          }
        }
        else {
          items = []
        }
        // console.log(value);
        // console.log(field.resultValue);
        field.items = items.slice(0, 50)
      }

      field.sameFieldValue = this.sameFieldValue

      // if (field.table) {
      //   const data = {
      //     proCode: 'B0113',
      //     name: 'B0113_百年理赔_百年理赔地址库',
      //     queryNames: {
      //       ['市']: {
      //         $regex: `/${value}/`
      //       }
      //     },
      //     respNames: [
      //       '市', '省'
      //     ],
      //     pageSize: 10,
      //     pageIndex: 1
      //   }

      //   const result = await reqConst({
      //     url: '/sys-const/page',
      //     method: "POST",
      //     data,
      //   })
      //   console.log(result);
      //   let fieldContent = result.list
      //   let fieldItems = fieldContent.map(item => {
      //     return item['市']
      //   })
      //   field.items = fieldItems
      // }
    }
  }
}