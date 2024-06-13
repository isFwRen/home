import moment from 'moment'
import BigNumber from 'bignumber.js'
import { tools, sessionStorage } from 'vue-rocket'
import { codesList46, codesList48, validate50CommonFn } from './cells'
import { ignore } from '../tools'
import { getNode, getNodeValue } from './tools'
import { ALLDESSERTS } from '../../mixins/DropCells'
import { MessageBox, Notification } from 'element-ui';

const B0114 = {
  op0: {
    // 记录最后一次存储的合法field
    memoFields: [],

    // 记录相同 code 的 field 的值
    memoFieldValues: [],

    // fields 的值从 targets 里的值选择
    dropdownFields: [
      // 3
      {
        targets: ['fc094'],
        fields: ['fc091', 'fc092', 'fc093']
      }
    ],

    // 校验规则
    rules: [
      // 2
      {
        fields: ['fc094'],
        validate2: function ({ field, fieldsObject, thumbIndex, value }) {
          const fc094Values = []

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage

            if (sessionStorage || thumbIndex === +key) {
              const _fieldsList = fieldsObject[key].fieldsList

              for (let _fields of _fieldsList) {
                for (let _field of _fields) {
                  if (_field.code === 'fc094' && _field.uniqueId !== field.uniqueId) {
                    fc094Values.push(_field.resultValue)
                  }
                }
              }
            }
          }

          if (fc094Values.includes(value)) {
            return '发票/报销单属性不能重复!'
          }

          return true
        }
      },

      // 4
      {
        fields: ['fc091', 'fc092'],
        validate4: function ({ includes, value }) {
          if (includes) {
            const result = includes.find(text => text === value)

            if (!result) {
              return '没有此发票/报销单，请核实!'
            }
          }

          return true
        }
      },

      // 12
      {
        fields: ['fc008'],
        validate12: function ({ effectValidations, value, items }) {
          if (ignore({ effectValidations, value })) return true

          const result = items.find(text => text === value)

          if (result) {
            return true
          }

          return '事故地所属区县录入错误，请根据下拉提示内容选录.'
        }
      },

      // 14
      {
        fields: ['fc003'],
        validate14: function ({ effectValidations, field, fieldsIndex, fieldsObject, items, thumbIndex, value }) {
          if (ignore({ effectValidations, value })) return true

          const fc008Field = fieldsObject[thumbIndex].fieldsList[fieldsIndex][3]
          const { targetDesserts } = fc008Field
          if (!targetDesserts) return true

          if (!ignore({ effectValidations: fc008Field.effectValidations, value: fc008Field.resultValue })) {
            const cities = []

            for (let dessert of targetDesserts) {
              cities.push(dessert[3])
            }

            if (!cities.includes(value)) {
              return '市与区县不匹配，如无匹配的内容请录入A'
            }
          }

          return true
        }
      },

      // 34
      {
        fields: ['fc106'],
        validate34: function ({ effectValidations, field, items, value }) {
          if (!value) return true
          if (ignore({ effectValidations, value })) return true

          field.allowForce = true
          const result = items.find(text => text === value)

          if (result) {
            return true
          }

          return '治疗医院录入错误，请按下拉提示内容进行选录.'
        }
      }
    ],

    // 提示文本
    hints: [],

    // 工序完成初始化
    init: {
      methods: {
        // 9(xml)
        validateOtherInfo9: function ({ bill }) {
          const otherInfo = bill.otherInfo
          const values = getNodeValue(otherInfo, 'causeCode')
          const str = values.toString()

          if (str === '03' || str === '04') {

            if (sessionStorage.get('isApp')?.isApp === 'true') {
              // MessageBox.alert('当前案件为身故案件，不需要录入发票、清单、报销单!', '请注意', {
              //   type: 'warning',
              //   confirmButtonText: '确定',
              //   showClose: false,
              // })
              return Notification({
                type: 'warning',
                title: '提醒(5s后自动关闭)',
                message: '当前案件为身故案件，不需要录入发票、清单、报销单!',
                duration: 5000,
                position: 'top-left'
              })
            } else {
              alert('当前案件为身故案件，不需要录入发票、清单、报销单!')
            }

          }
        },

        // 51(xml)
        validateOtherInfo51: function ({ bill }) {
          const otherInfo = bill.otherInfo
          const xmls = getNode(otherInfo, 'applyCauses')
          let values = []

          for (let xml of xmls) {
            const vals = getNodeValue(xml, 'causeCode')
            values = [...values, ...vals]
          }

          if (values.includes('09') || values.includes('10')) {

            if (sessionStorage.get('isApp')?.isApp === 'true') {
              // MessageBox.alert('理赔申请项目为重疾，需要切重疾分块.', '请注意', {
              //   type: 'warning',
              //   confirmButtonText: '确定',
              //   showClose: false,
              // })
              return Notification({
                type: 'warning',
                title: '提醒(5s后自动关闭)',
                message: '理赔申请项目为重疾，需要切重疾分块.',
                duration: 5000,
                position: 'top-left'
              })
            } else {
              alert('理赔申请项目为重疾，需要切重疾分块.')
            }

          }

          if (values.includes('17') || values.includes('18')) {

            if (sessionStorage.get('isApp')?.isApp === 'true') {
              // MessageBox.alert('理赔申请项目为轻症，需要切轻症分块.', '请注意', {
              //   type: 'warning',
              //   confirmButtonText: '确定',
              //   showClose: false,
              // })
              return Notification({
                type: 'warning',
                title: '提醒(5s后自动关闭)',
                message: '理赔申请项目为轻症，需要切轻症分块.',
                duration: 5000,
                position: 'top-left'
              })
            } else {
              alert('理赔申请项目为轻症，需要切轻症分块.')
            }

          }

          if (values.includes('27') || values.includes('28')) {

            if (sessionStorage.get('isApp')?.isApp === 'true') {
              // MessageBox.alert('理赔申请项目为中症，需要切中症分块.', '请注意', {
              //   type: 'warning',
              //   confirmButtonText: '确定',
              //   showClose: false,
              // })
              return Notification({
                type: 'warning',
                title: '提醒(5s后自动关闭)',
                message: '理赔申请项目为中症，需要切中症分块.',
                duration: 5000,
                position: 'top-left'
              })
            } else {
              alert('理赔申请项目为中症，需要切中症分块.')
            }

          }

          if (values.includes('05') || values.includes('06')) {

            if (sessionStorage.get('isApp')?.isApp === 'true') {
              // MessageBox.alert('理赔申请项目为全残，需要切全残分块.', '请注意', {
              //   type: 'warning',
              //   confirmButtonText: '确定',
              //   showClose: false,
              // })
              return Notification({
                type: 'warning',
                title: '提醒(5s后自动关闭)',
                message: '理赔申请项目为全残，需要切全残分块.',
                duration: 5000,
                position: 'top-left'
              })
            } else {
              alert('理赔申请项目为全残，需要切全残分块.')
            }

          }

          if (values.includes('07') || values.includes('08')) {

            if (sessionStorage.get('isApp')?.isApp === 'true') {
              // MessageBox.alert('理赔申请项目为残疾，需要切残疾分块.', '请注意', {
              //   type: 'warning',
              //   confirmButtonText: '确定',
              //   showClose: false,
              // })
              return Notification({
                type: 'warning',
                title: '提醒(5s后自动关闭)',
                message: '理赔申请项目为残疾，需要切残疾分块.',
                duration: 5000,
                position: 'top-left'
              })
            } else {
              alert('理赔申请项目为残疾，需要切残疾分块.')
            }

          }

          if (values.includes('25') || values.includes('26')) {

            if (sessionStorage.get('isApp')?.isApp === 'true') {
              // MessageBox.alert('理赔申请项目为特种病，需要切特种病分块.', '请注意', {
              //   type: 'warning',
              //   confirmButtonText: '确定',
              //   showClose: false,
              // })
              return Notification({
                type: 'warning',
                title: '提醒(5s后自动关闭)',
                message: '理赔申请项目为特种病，需要切特种病分块.',
                duration: 5000,
                position: 'top-left'
              })
            } else {
              alert('理赔申请项目为特种病，需要切特种病分块.')
            }

          }
        },

        // 54(xml)
        validateOtherInfo54: function ({ bill }) {
          const otherInfo = bill.otherInfo
          const xmls = getNode(otherInfo, 'applyCauses')
          let values = []

          for (let xml of xmls) {
            const vals = getNodeValue(xml, 'causeCode')
            values = [...values, ...vals]
          }

          if (values.includes('21') || values.includes('22')) {

            if (sessionStorage.get('isApp')?.isApp === 'true') {
              // MessageBox.alert('理赔申请项目为津贴，请注意切“紧急救护车使用津贴', '请注意', {
              //   type: 'warning',
              //   confirmButtonText: '确定',
              //   showClose: false,
              // })
              return Notification({
                type: 'warning',
                title: '提醒(5s后自动关闭)',
                message: '理赔申请项目为津贴，请注意切“紧急救护车使用津贴',
                duration: 5000,
                position: 'top-left'
              })
            } else {
              alert('理赔申请项目为津贴，请注意切“紧急救护车使用津贴”.')
            }

          }
        }
      }
    },

    // 字段已生成
    updateFields: {
      methods: {
        // 12
        setConstants12: function ({ flatFieldsList }) {
          const fields = ['fc008']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0114_华夏理赔_华夏理赔地址库',
                query: '区县'
              }
            }
          })
        },

        // 15
        setConstants15: function ({ flatFieldsList }) {
          const fields = ['fc003']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0114_华夏理赔_华夏理赔地址库',
                query: '市'
              }
            }
          })
        },

        // 16
        setConstants16: function ({ flatFieldsList }) {
          const fields = ['fc004']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0114_华夏理赔_华夏理赔地址库',
                query: '省'
              }
            }
          })
        },

        // 34
        setConstants34: function ({ flatFieldsList }) {
          const fields = ['fc106']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0114_华夏理赔_医院名称表',
                query: '医院名称'
              }
            }
          })
        }
      }
    },

    // 回车
    enter: {
      methods: {
        // 13
        disable13({ field, fieldsList, focusFieldsIndex }) {
          if (field.code === 'fc008') {
            // const { desserts, resultValue } = field
            const { resultValue } = field
            // 统计区/县是否唯一
            let count = 0
            const targetDesserts = []

            for (let dessert of ALLDESSERTS) {
              if (dessert[4] === resultValue) {
                ++count
                targetDesserts.push(dessert)
              }
            }

            field.targetDesserts = targetDesserts

            const fields = fieldsList[focusFieldsIndex]

            fields.map(field => {
              if (field.code === 'fc003' || field.code === 'fc004') {
                if (count === 1) {
                  field.op0Value = ''
                  field.resultValue = ''
                  field.disabled = true
                }
                else {
                  field.disabled = false
                }
              }
            })
          }
        },

        // 15
        disable15({ field, fieldsList, focusFieldsIndex }) {
          if (field.code === 'fc003') {
            if (field.items.includes(field.resultValue)) {
              const fields = fieldsList[focusFieldsIndex]

              fields.map(field => {
                if (field.code === 'fc004') {
                  field.op0Value = ''
                  field.resultValue = ''
                  field.disabled = true
                }
              })
            }

          }
        }
      }
    },

    // F8(提交前校验)
    beforeSubmit: {
      methods: {
        // 5
        validate5({ fieldsObject }) {
          const [fc092Values, fc094Values] = [[], []]

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  const { code, resultValue } = field

                  if (code === 'fc092') {
                    resultValue && fc092Values.push(resultValue)
                  }

                  if (code === 'fc094') {
                    resultValue && fc094Values.push(resultValue)
                  }
                }
              }
            }
          }

          for (let value of fc092Values) {
            if (!fc094Values.includes(value)) {
              return {
                popup: 'confirm',
                errorMessage: `报销单${value}没有匹配的发票!`
              }
            }
          }

          return true
        },

        // 6
        validate6({ fieldsObject }) {
          const [fc091Values, fc094Values] = [[], []]

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  const { code, resultValue } = field

                  if (code === 'fc091') {
                    resultValue && fc091Values.push(resultValue)
                  }

                  if (code === 'fc094') {
                    resultValue && fc094Values.push(resultValue)
                  }
                }
              }
            }
          }

          for (let value of fc094Values) {
            if (!fc091Values.includes(value)) {
              return {
                popup: 'confirm',
                errorMessage: `发票${value}没有匹配的清单!`
              }
            }
          }

          return true
        },

        // 7
        validate7({ fieldsObject }) {
          const nums = [0, 1, 2, 20]
          const fc179Values = []

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  const { code, resultValue } = field

                  if (code === 'fc179') {
                    resultValue && fc179Values.push(+resultValue)
                  }
                }
              }
            }
          }

          for (let num of fc179Values) {
            if (nums.includes(num)) {
              return true
            }
          }

          return {
            errorMessage: '案件必须有申请表!'
          }
        },

        // 10
        validate10({ fieldsObject }) {
          const fc179Values = []

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  const { code, resultValue } = field

                  if (code === 'fc179') {
                    resultValue && fc179Values.push(+resultValue)
                  }
                }
              }
            }
          }

          if (fc179Values.includes(6)) {
            return true
          }

          return {
            errorMessage: '无匹配的诊断书，请确认'
          }
        },

        // 50
        validate50a({ bill, fieldsObject }) {
          const { values, fc179Values } = validate50CommonFn({ bill, fieldsObject })

          if (values.includes('09') || values.includes('10')) {
            if (!fc179Values.includes(61)) {
              return {
                errorMessage: '理赔申请项目为重疾，需要切重疾分块'
              }
            }
          }

          return true
        },

        validate50b({ bill, fieldsObject }) {
          const { values, fc179Values } = validate50CommonFn({ bill, fieldsObject })

          if (values.includes('17') || values.includes('18')) {
            if (!fc179Values.includes(62)) {
              return {
                errorMessage: '理赔申请项目为轻症，需要切轻症分块'
              }
            }
          }

          return true
        },

        validate50c({ bill, fieldsObject }) {
          const { values, fc179Values } = validate50CommonFn({ bill, fieldsObject })

          if (values.includes('27') || values.includes('28')) {
            if (!fc179Values.includes(63)) {
              return {
                errorMessage: '理赔申请项目为中症，需要切中症分块'
              }
            }
          }

          return true
        },

        validate50d({ bill, fieldsObject }) {
          const { values, fc179Values } = validate50CommonFn({ bill, fieldsObject })

          if (values.includes('05') || values.includes('06')) {
            if (!fc179Values.includes(64)) {
              return {
                errorMessage: '理赔申请项目为全残，需要切全残分块'
              }
            }
          }

          return true
        },

        validate50e({ bill, fieldsObject }) {
          const { values, fc179Values } = validate50CommonFn({ bill, fieldsObject })

          if (values.includes('07') || values.includes('08')) {
            if (!fc179Values.includes(65)) {
              return {
                errorMessage: '理赔申请项目为残疾，需要切残疾分块'
              }
            }
          }

          return true
        },

        validate50f({ bill, fieldsObject }) {
          const { values, fc179Values } = validate50CommonFn({ bill, fieldsObject })

          if (values.includes('25') || values.includes('26')) {
            if (!fc179Values.includes(66)) {
              return {
                errorMessage: '理赔申请项目为特种病，需要切特种病分块'
              }
            }
          }

          return true
        },

        // 52(xml)
        validate52({ bill, fieldsObject }) {
          const otherInfo = bill.otherInfo
          const values = getNodeValue(otherInfo, 'isAllowance')
          const str = values.toString()

          const getFc179Values = () => {
            const fc179Values = []

            for (let key in fieldsObject) {
              const sessionStorage = fieldsObject[key].sessionStorage
              const fieldsList = fieldsObject[key].fieldsList

              if (sessionStorage) {
                for (let fields of fieldsList) {
                  for (let field of fields) {
                    const { code, resultValue } = field

                    if (code === 'fc179') {
                      resultValue && fc179Values.push(+resultValue)
                    }
                  }
                }
              }
            }

            return fc179Values
          }

          if (str === 'Y') {
            const fc179Values = getFc179Values()

            if (!fc179Values.includes(4) && !fc179Values.includes(41)) {
              alert('“icu/烧伤病房/抢救室”不为空')
            }
          }

          return true
        }
      }
    }
  },

  op1op2opq: {
    // 校验规则
    rules: [
      // 19
      {
        fields: ['fc014'],
        validate19: function ({ effectValidations, items, value }) {
          if (ignore({ effectValidations, value })) return true

          const result = items.find(text => text === value)

          if (result) {
            return true
          }

          return '意外原因错误，请根据下拉提示内容选录.'
        }
      },

      // 20
      {
        fields: ['fc015'],
        validate20: function ({ effectValidations, items, value }) {
          if (ignore({ effectValidations, value })) return true

          const result = items.find(text => text === value)

          if (result) {
            return true
          }

          return '损伤外部原因错误，请根据下拉提示内容选录.'
        }
      },

      // 22
      {
        fields: ['fc096', 'fc187', 'fc190'],
        validate22: function ({ effectValidations, items, value }) {
          if (!value) return true
          if (ignore({ effectValidations, value })) return true

          const result = items.find(text => text === value)

          if (result) {
            return true
          }

          return '手术术士编码错误，请根据下拉提示内容选录.'
        }
      },

      // 23
      {
        fields: ['fc101', 'fc181', 'fc182'],
        validate23: function ({ effectValidations, items, value }) {
          if (!value) return true
          if (ignore({ effectValidations, value })) return true

          const result = items.find(text => text === value)

          if (result) {
            return true
          }

          return '疾病代码错误，请根据下拉提示内容选录.'
        }
      },

      // 33
      {
        fields: ['fc267'],
        validate33: function ({ effectValidations, fieldsList, fieldsIndex, value }) {
          if (ignore({ effectValidations, value })) return true

          const fields = fieldsList[fieldsIndex]
          const fc266Field = tools.find(fields, { code: 'fc266' }) || {}

          if (!['2', '3'].includes(fc266Field.resultValue)) return true

          if (['1', '2'].includes(value)) return true

          return '医保类型错误.'
        }
      },

      // 34
      {
        fields: ['fc106'],
        validate34: function ({ effectValidations, field, items, value }) {
          if (!value) return true
          if (ignore({ effectValidations, value })) return true

          field.allowForce = true

          const result = items.find(text => text === value)

          if (result) {
            return true
          }

          return '治疗医院录入错误，请按下拉提示内容进行选录.'
        }
      },

      // 38
      {
        fields: ['fc016'],
        validate38: function ({ effectValidations, items, value }) {
          if (!value) return true
          if (ignore({ effectValidations, value })) return true

          const result = items.find(text => text === value)

          if (result) {
            return true
          }

          return '重大疾病编码错误，请根据下拉提示内容选录.'
        }
      },

      // 39
      {
        fields: ['fc019'],
        validate39: function ({ effectValidations, items, value }) {
          if (!value) return true
          if (ignore({ effectValidations, value })) return true

          const result = items.find(text => text === value)

          if (result) {
            return true
          }

          return '轻症重疾编码错误，请根据下拉提示内容选录.'
        }
      },

      // 40
      {
        fields: ['fc002', 'fc005', 'fc006', 'fc007', 'fc011', 'fc017', 'fc018', 'fc020', 'fc022', 'fc024', 'fc027', 'fc028', 'fc108', 'fc110', 'fc111', 'fc112'],
        validateDate: function ({ effectValidations, value }) {
          if (!value) return true
          if (ignore({ effectValidations, value })) return true

          // if(/[A, \?]/.test(value)) {
          //   return true
          // }

          if (value.length !== 6 || moment(`20${value}`).format('YYYYMMDD') === 'Invalid date') {
            return '日期格式错误! '
          }

          return true
        }
      },

      // 42
      {
        fields: ['fc021'],
        validate42: function ({ effectValidations, items, value }) {
          if (!value) return true
          if (ignore({ effectValidations, value })) return true

          const result = items.find(text => text === value)

          if (result) {
            return true
          }

          return '中症疾病编码错误，请根据下拉提示内容选录.'
        }
      },

      // 43
      {
        fields: ['fc023'],
        validate43: function ({ effectValidations, items, value }) {
          if (!value) return true
          if (ignore({ effectValidations, value })) return true

          const result = items.find(text => text === value)

          if (result) {
            return true
          }

          return '特种病编码错误，请根据下拉提示内容选录.'
        }
      },

      // 44
      {
        fields: ['fc025'],
        validate44: function ({ effectValidations, items, value }) {
          if (!value) return true
          if (ignore({ effectValidations, value })) return true

          const result = items.find(text => text === value)

          if (result) {
            return true
          }

          return '全残项目编码错误，请根据下拉提示内容选录.'
        }
      },

      // 45
      {
        fields: ['fc029'],
        validate44: function ({ effectValidations, items, value }) {
          if (!value) return true
          if (ignore({ effectValidations, value })) return true

          const result = items.find(text => text === value)

          if (result) {
            return true
          }

          return '伤残项目编码错误，请根据下拉提示内容选录.'
        }
      },

      // 55
      {
        index: 0,
        fields: ['fc188', 'fc191'],
        validate55: function ({ effectValidations, field, fieldsList, fieldsIndex, value }) {
          if (ignore({ effectValidations, value })) return true

          const fields = fieldsList[fieldsIndex]

          const fc187Field = tools.find(fields, { code: 'fc187' }) || {}
          const fc190Field = tools.find(fields, { code: 'fc190' }) || {}

          if (field.code === 'fc188') {
            if (fc187Field.resultValue) {
              if (!value) return '手术诊断类型2必录.'
            }
          }

          if (field.code === 'fc191') {
            if (fc190Field.resultValue) {
              if (!value) return '手术诊断类型3必录.'
            }
          }

          return true
        }
      },

      // CSB0114RC0133000
      {
        fields: ['fc138', 'fc142', 'fc146', 'fc150', 'fc154', 'fc158', 'fc162', 'fc166'],
        validate133: function ({ value }) {
          if (value == 'DR' || value == '？' || value == '?' || value == '') return true
          // 是否有中文字符
          let pattern = /[\u4E00-\u9FFF\u3400-\u4DFF\uF900-\uFAFF]/;
          // 有几个中文字符
          let patterns = /[\u4E00-\u9FFF\u3400-\u4DFF\uF900-\uFAFF]/g;
          // value元素个数
          let count = [...value].length
          // 中文字符个数
          let matches = value.match(patterns)?.length;

          if (!pattern.test(value)) {
            return '录入内容有误，请检查'
          } else if (!pattern.test(value) && count == 1) {
            return '录入内容有误，请检查'
          } else {
            return true
          }
          // else if (matches * 2 < count - matches) {
          //   return '录入内容有误，请检查'
          // } 
        },
      },
    ],

    // 提示文本
    hints: [
      // CSB0114RC0213000
      {
        fields: ['fc138', 'fc142', 'fc146', 'fc150', 'fc154', 'fc158', 'fc162', 'fc166'],
        hintFc152({ field }) {
          // const codes =  ['fc138', 'fc142', 'fc146', 'fc150', 'fc154', 'fc158', 'fc162', 'fc166']

          // if (codes.includes(field.code)) {
          field.allowForce = true
          // console.log(field.ocrBlur);
          if (field.items) {
            const result = field.items.find(text => text === field.resultValue)
            if (result) return true
          }
          if (field.ocrBlur && field.ocrBlur?.length != 0) {
            const result = field.ocrBlur.find(text => text === field.resultValue)
            // console.log(field.resultValue);
            console.log(result);
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
                  console.log('count', count);
                  console.log('flag1', flag1, 'flag2', flag2, 'flag3', flag3);
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
              console.log(showItems);
              if (!showItems[0]) {
                field.hint = ''
                return true
              }

              if (showItems[1]) {
                return `<p style="color: blue; fontSize: 14px; margin-top: -3px; margin-bottom: 0px"">${showItems[0]}，${showItems[1]}</p>`
                // return false
              } else if (!showItems[1]) {
                return `<p style="color: blue; fontSize: 14px; margin-top: -3px; margin-bottom: 0px"">${showItems[0]}</p>`
                // return false
              }
            }
          }
          return true
          // }
        },
      }
    ],

    // 字段已生成
    updateFields: {
      methods: {
        // 17
        validateOtherInfo17({ bill, flatFieldsList }) {
          const otherInfo = bill.otherInfo
          const values = getNodeValue(otherInfo, 'accidentDesc')
          const str = values.toString()

          if (str.length > 30) {
            flatFieldsList.map(_field => {
              if (_field.code === 'fc012') {
                _field.disabled = true
              }
            })
          }
        },

        // 18(xml)
        validateOtherInfo18({ bill, flatFieldsList }) {
          const otherInfo = bill.otherInfo
          const values = getNodeValue(otherInfo, 'causeCode')
          const nums = ['02', '04', '06', '08', '10', '12', '14', '16', '18', '20', '22', '24', '26', '28']

          const hasSome = values.some(val => nums.includes(val))

          if (!hasSome) {
            flatFieldsList.map(_field => {
              if (_field.code === 'fc014') {
                _field.disabled = true
              }

              if (_field.code === 'fc015') {
                _field.disabled = true
              }
            })
          }
        },

        // 16
        setConstants16({ flatFieldsList }) {
          const fields = ['fc015']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0114_华夏理赔_华夏理赔损伤外部原因表',
                query: '描述'
              }
            }
          })
        },

        // 19
        setConstants19({ flatFieldsList }) {
          const fields = ['fc014']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0114_华夏理赔_华夏理赔意外原因表',
                query: '意外编码描述'
              }
            }
          })
        },

        // 20
        setConstants20({ flatFieldsList }) {
          const fields = ['fc015']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0114_华夏理赔_华夏理赔损伤外部原因表',
                query: '描述'
              }
            }
          })
        },

        // 22
        setConstants22({ flatFieldsList }) {
          const fields = ['fc096', 'fc187', 'fc190']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0114_华夏理赔_手术术式编码表',
                query: '手术名称'
              }
            }
          })
        },

        // 23
        setConstants23({ flatFieldsList }) {
          const fields = ['fc101', 'fc181', 'fc182']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0114_华夏理赔_ICD疾病代码表',
                query: '疾病名称'
              }
            }
          })
        },

        // 34
        setConstants34: function ({ flatFieldsList }) {
          const fields = ['fc106']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0114_华夏理赔_医院名称表',
                query: '医院名称'
              }
            }
          })
        },

        // 38 39
        setConstants3839({ flatFieldsList }) {
          const fields = ['fc016', 'fc019']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0114_华夏理赔_重疾和轻症疾病名称表',
                query: '疾病名称'
              }
            }
          })
        },

        // 42 43
        setConstants4243({ flatFieldsList }) {
          const fields = ['fc021', 'fc023']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0114_华夏理赔_重疾和轻症疾病名称表',
                query: '疾病名称'
              }
            }
          })
        },

        // 44
        setConstants44({ flatFieldsList }) {
          const fields = ['fc025']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0114_华夏理赔_全残信息表',
                query: '伤残名称'
              }
            }
          })
        },

        // 45
        setConstants45({ flatFieldsList }) {
          const fields = ['fc029']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0114_华夏理赔_伤残项目表',
                query: '伤残名称'
              }
            }
          })
        },

        // 24
        disable24({ fieldsList }) {
          const codes = [
            'fc099', 'fc183', 'fc184',
            'fc102', 'fc185', 'fc186'
          ]

          fieldsList?.map(fields => {
            fields?.map(_field => {
              if (codes.includes(_field.code)) {
                _field.disabled = true
              }
            })
          })
        },

        // 37(xml)
        disable37({ bill, flatFieldsList }) {
          const otherInfo = bill.otherInfo
          const xmls = getNode(otherInfo, 'applyCauses')
          let values = []

          for (let xml of xmls) {
            const vals = getNodeValue(xml, 'causeCode')
            values = [...values, ...vals]
          }

          if (!values.includes('21') && !values.includes('22')) {
            flatFieldsList.map(_field => {
              if (_field.code === 'fc048') {
                _field.disabled = true
              }
            })
          }
        },

        // 41 
        disable41({ op, fieldsList, focusFieldsIndex }) {
          if (op === 'op0') {
            return
          }

          const codesList = [
            ['fc098', 'fc189', 'fc192'],
            ['fc256', 'fc116', 'fc117', 'fc118', 'fc119', 'fc120', 'fc121', 'fc122', 'fc123', 'fc124', 'fc125', 'fc129', 'fc114'],
            ['fc172', 'fc173', 'fc176'],
            ['fc133', 'fc134', 'fc135', 'fc136', 'fc137'],
            ['fc026'],
            ['fc031'],
            ['fc170', 'fc171', 'fc174', 'fc175'],
            ['fc045'],
            ['fc274', 'fc275', 'fc276', 'fc277', 'fc278', 'fc279', 'fc282', 'fc280', 'fc281'],
            ['fc103', 'fc090', 'fc115', 'fc128']
          ]
          const flatCodesList = []

          codesList.map(codes => {
            flatCodesList.push(...codes)
          })

          const fields = fieldsList[focusFieldsIndex]

          fields?.map(_field => {
            if (flatCodesList.includes(_field.code)) {
              _field.disabled = true
            }
          })
        },

        // 46
        disable46({ block, fieldsList, focusFieldsIndex }) {
          const codes = []

          for (let _codes of codesList46) {
            codes.push(..._codes.slice(1))
          }

          const fields = fieldsList[focusFieldsIndex]

          fields.map(_field => {
            if (codes.includes(_field.code)) {
              _field.disabled = true
            }
          })
        },

        // 48
        disable48({ fieldsList, focusFieldsIndex }) {
          const codes = []

          for (let _codes of codesList48) {
            codes.push(..._codes.slice(2))
          }

          const fields = fieldsList[focusFieldsIndex]

          fields.map(_field => {
            if (codes.includes(_field.code)) {
              _field.disabled = true
            }
          })
        },

        // 49
        disable49({ bill, block, flatFieldsList, op }) {
          if (op !== 'op2') return
          if (block.temp !== 'MB002') return
          if (block.code !== 'bc010') return

          const codesList = [
            ['fc138', 'fc139', 'fc140', 'fc141', 'fc231', 'fc247', 'fc207', 'fc215', 'fc257'],
            ['fc142', 'fc143', 'fc144', 'fc145', 'fc232', 'fc248', 'fc208', 'fc216', 'fc258'],
            ['fc146', 'fc147', 'fc148', 'fc149', 'fc233', 'fc249', 'fc209', 'fc217', 'fc259'],
            ['fc150', 'fc151', 'fc152', 'fc153', 'fc234', 'fc250', 'fc210', 'fc218', 'fc260'],
            ['fc154', 'fc155', 'fc156', 'fc157', 'fc235', 'fc251', 'fc211', 'fc219', 'fc261'],
            ['fc158', 'fc159', 'fc160', 'fc161', 'fc236', 'fc252', 'fc212', 'fc220', 'fc262'],
            ['fc162', 'fc163', 'fc164', 'fc165', 'fc237', 'fc253', 'fc213', 'fc221', 'fc263'],
            ['fc166', 'fc167', 'fc168', 'fc169', 'fc238', 'fc254', 'fc214', 'fc222', 'fc264']
          ]

          // 1. 左边八列字段二码录入时进行屏蔽，字段内容显示一码录入内容
          {
            const firstEightCodes = []

            codesList.map(codes => {
              firstEightCodes.push(...codes.slice(0, 8))
            })

            flatFieldsList.map(_field => {
              if (firstEightCodes.includes(_field.code)) {
                _field[`${op}Value`] = _field.op1Value
                _field.resultValue = _field.op1Value
                _field.disabled = true
              }
            })
          }

          // 4. 二码录入时，当左边两列的字段录入值均为空时，屏蔽对应的第九列字段
          {
            const firstTwoCodes = []
            const ninthOneCodes = []

            codesList.map(codes => {
              firstTwoCodes.push(codes.slice(0, 2))
              ninthOneCodes.push(codes[8])
            })

            firstTwoCodes.map((codes, codesIndex) => {
              const [code1, code2] = codes
              field1 = tools.find(flatFieldsList, { code: code1 }) || {}
              field2 = tools.find(flatFieldsList, { code: code2 }) || {}

              if (!field1.resultValue && !field2.resultValue) {
                for (let _field of flatFieldsList) {
                  if (_field.code === ninthOneCodes[codesIndex]) {
                    _field.disabled = true
                  }
                }
              }
            })
          }
        },

        // 53
        disable53({ bill, fieldsList, focusFieldsIndex }) {
          const otherInfo = bill.otherInfo
          const values = getNodeValue(otherInfo, 'isAllowance')
          const str = values.toString()
          const fields = fieldsList[focusFieldsIndex]

          const codes2 = ['fc207', 'fc208', 'fc209', 'fc210', 'fc211', 'fc212', 'fc213', 'fc214']

          if (str === 'N') {
            fields.map(_field => {
              if (codes2.includes(_field.code)) {
                _field.disabled = true
              }
            })
          }
        },

        // 48
        set48Items({ bill, codeValues = {}, fieldsList, focusFieldsIndex }) {
          const { fc008, fc003, fc004 } = codeValues
          if (!fc008 && !fc003 && !fc004) return
          const db = window['constantsDB']['B0114']
          if (!db) return
          const collections1 = db['B0114_华夏理赔_华夏理赔地址库']
          const collections2 = db['B0114_华夏理赔_机构编码对应表']


          // fc008常量名列表、fc003常量名列表、fc004常量名列表、上级机构常量名列表
          const [fc008Items, fc003Items, fc004Items, agencyItems] = [[], [], [], []]

          // 所属医疗目录常量名列表
          const [medicalItems1, medicalItems2] = [[], []]

          // 所属医疗目录常量名
          let [medicalValue1, medicalValue2] = [void 0, void 0]
          // 集合
          let [medicalHeaders, medicalDesserts1, medicalDesserts2] = [[], [], []]

          for (let dessert of collections1.desserts) {
            medicalItems1.push(dessert[0])
            fc008Items.push(dessert[4])
            fc003Items.push(dessert[3])
            fc004Items.push(dessert[2])
          }

          const indexFc008 = fc008Items.indexOf(fc008)
          const lastIndexFc008 = fc008Items.lastIndexOf(fc008)

          const indexFc003 = fc008Items.indexOf(fc003)
          const lastIndexFc003 = fc008Items.lastIndexOf(fc003)

          const indexFc004 = fc008Items.indexOf(fc004)
          const lastIndexFc004 = fc008Items.lastIndexOf(fc004)

          // 若fc008唯一
          if (indexFc008 === lastIndexFc008) {
            medicalValue1 = medicalItems1[indexFc008]
          }
          // 若fc003唯一
          else if (indexFc003 === lastIndexFc003) {
            medicalValue1 = medicalItems1[indexFc003]
          }
          // 若fc004唯一
          else if (indexFc004 === lastIndexFc004) {
            medicalValue1 = medicalItems1[indexFc004]
          }

          const otherInfo = bill.otherInfo
          const values = getNodeValue(otherInfo, 'policyBranchCode')
          if (values.length) {
            for (let dessert of collections2.desserts) {
              medicalItems2.push(dessert[0])
              agencyItems.push(dessert[4])
            }

            const index = agencyItems.indexOf(values.toString())

            medicalValue2 = medicalItems2[index]
          }

          if (medicalValue1) {
            const medicalCollections1 = db[`B0114_华夏理赔_省份-${medicalValue1}`]

            if (medicalCollections1) {
              medicalHeaders = medicalCollections1.headers
              medicalDesserts1 = medicalCollections1.desserts
            }
          }

          if (medicalValue2) {
            const medicalCollections2 = db[`B0114_华夏理赔_省份-${medicalValue2}`]

            if (medicalCollections2) {
              medicalHeaders = medicalCollections2.headers
              medicalDesserts2 = medicalCollections2.desserts
            }
          }

          const assignDesserts = [...medicalDesserts1, ...medicalDesserts2]

          const codes = []

          for (let _codes of codesList48) {
            codes.push(_codes[0])
          }

          const fields = fieldsList[focusFieldsIndex]

          for (let _field of fields) {
            if (codes.includes(_field.code)) {
              const index = tools.findIndex(codes, _field.code)
              const field = tools.find(fields, { code: codes[index] }) || {}

              let key
              if (medicalValue1) {
                key = `B0114_华夏理赔_省份-${medicalValue1}`
              }
              if (medicalValue2) {
                key = `B0114_华夏理赔_省份-${medicalValue2}`
              }

              if (medicalValue1 && medicalValue2) {
                key = `B0114_华夏理赔_省份-${medicalValue1}-${medicalValue2}`
              }
              // const key = `B0114_华夏理赔_省份-${medicalValue1}-${medicalValue2}`

              window['constantsDB']['B0114'][key] = {
                headers: medicalHeaders,
                desserts: assignDesserts
              }

              field.table = {
                name: key,
                query: '中文名称'
              }

              // field.desserts = assignDesserts
            }
          }
        }
      }
    },

    // 回车
    enter: {
      methods: {
        // 25
        disable25({ field, fieldsList, focusFieldsIndex, op }) {
          if (field.code !== 'fc265') return

          const codes = ['fc266']
          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue === '2') {
            fields.map(_field => {
              if (codes.includes(_field.code)) {
                _field[`${op}Value`] = '1'
                _field.resultValue = '1'
                _field.disabled = true
              }
            })
          }
          else {
            fields.map(_field => {
              if (codes.includes(_field.code)) {
                _field[`${op}Value`] = ''
                _field.resultValue = ''
                _field.disabled = false
              }
            })
          }
        },

        // 26
        disable26({ field, fieldsList, focusFieldsIndex, op }) {
          if (field.code !== 'fc265') return

          const codes = ['fc271']
          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue === '2') {
            fields.map(_field => {
              if (codes.includes(_field.code)) {
                _field[`${op}Value`] = ''
                _field.resultValue = ''
                _field.disabled = true
              }
            })
          }
          else {
            fields.map(_field => {
              if (codes.includes(_field.code)) {
                _field.disabled = false
              }
            })
          }
        },

        // 27
        disable27({ field, fieldsList, focusFieldsIndex, op }) {
          if (field.code !== 'fc265' && field.code !== 'fc266') return

          const fields = fieldsList[focusFieldsIndex]
          const fc265Field = tools.find(fields, { code: 'fc265' }) || {}
          const fc266Field = tools.find(fields, { code: 'fc266' }) || {}
          const fc266EnableValues = ['2', '3', '4', '6']
          const codes = ['fc268', 'fc269']

          const yes = fc265Field.resultValue === '1' && fc266EnableValues.includes(fc266Field.resultValue)

          for (let _field of fields) {
            if (codes.includes(_field.code)) {
              if (yes) {
                _field[`${op}Value`] = ''
                _field.resultValue = ''
                _field.disabled = true
              }
              else {
                _field.disabled = false
              }
            }
          }

        },

        // 28
        disable28({ field, fieldsList, focusFieldsIndex, op }) {
          if (field.code !== 'fc266') return

          const fields = fieldsList[focusFieldsIndex]
          const fc266Field = tools.find(fields, { code: 'fc266' }) || {}
          const fc266Value = fc266Field.resultValue

          const codes1 = ['fc267']
          const codes2 = ['fc271', 'fc272', 'fc273']

          for (let _field of fields) {
            if (fc266Value === '5' || fc266Value === '6') {
              if (codes1.includes(_field.code)) {
                _field[`${op}Value`] = ''
                _field.resultValue = ''
                _field.disabled = true
              }
            }
            else {
              if (codes1.includes(_field.code)) {
                _field.disabled = false
              }
            }

            if (fc266Value === '6') {
              if (codes2.includes(_field.code)) {
                _field[`${op}Value`] = ''
                _field.resultValue = ''
                _field.disabled = true
              }
            }
            else {
              if (codes2.includes(_field.code)) {
                _field.disabled = false
              }
            }
          }
        },

        // 29
        disable29({ field, fieldsList, focusFieldsIndex, op }) {
          if (field.code !== 'fc105') return

          const fields = fieldsList[focusFieldsIndex]
          const fc105Field = tools.find(fields, { code: 'fc105' }) || {}
          const codes1 = ['fc111', 'fc112']
          const codes2 = ['fc110']

          for (let _field of fields) {
            if (fc105Field.resultValue === '1') {
              if (codes1.includes(_field.code)) {
                _field[`${op}Value`] = ''
                _field.resultValue = ''
                _field.disabled = true
              }
            }
            else {
              if (codes1.includes(_field.code)) {
                _field.disabled = false
              }
            }

            if (fc105Field.resultValue === '2') {
              if (codes2.includes(_field.code)) {
                _field[`${op}Value`] = ''
                _field.resultValue = ''
                _field.disabled = true
              }
            }
            else {
              if (codes2.includes(_field.code)) {
                _field.disabled = false
              }
            }
          }
        },

        // 30
        disable30({ field, fieldsList, focusFieldsIndex, op }) {
          const FC109 = 'fc109'

          if (field.code !== FC109) return

          const fields = fieldsList[focusFieldsIndex]
          const fc109Field = tools.find(fields, { code: FC109 }) || {}

          for (let _field of fields) {
            if (_field.code !== 'fc126') continue

            if (['A', '3'].includes(fc109Field.resultValue)) {
              _field.disabled = true
              _field[`${op}Value`] = ''
              _field.resultValue = ''
            }
            else {
              _field.disabled = false
            }
          }
        },

        // 32
        disable32({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc267') return

          const fields = fieldsList[focusFieldsIndex]
          const fc267Field = tools.find(fields, { code: 'fc267' }) || {}
          const codes1 = ['fc269']
          const codes2 = ['fc268']

          for (let _field of fields) {
            if (['1', '2'].includes(fc267Field.resultValue)) {
              if (codes1.includes(_field.code)) {
                _field.disabled = true
              }
            }
            else {
              if (codes1.includes(_field.code)) {
                _field.disabled = false
              }
            }

            if (fc267Field.resultValue === '3') {
              if (codes2.includes(_field.code)) {
                _field.disabled = true
              }
            }
            else {
              if (codes2.includes(_field.code)) {
                _field.disabled = false
              }
            }
          }
        },

        // 53
        disable53({ bill, fieldsList, focusFieldsIndex }) {
          const otherInfo = bill.otherInfo
          const values = getNodeValue(otherInfo, 'isAllowance')
          const str = values.toString()

          if (str !== 'Y') return

          const codes1 = ['fc138', 'fc142', 'fc146', 'fc150', 'fc154', 'fc158', 'fc162', 'fc166']
          const codes2 = ['fc207', 'fc208', 'fc209', 'fc210', 'fc211', 'fc212', 'fc213', 'fc214']

          const texts = ['抢救室', 'ICU', '烧伤病']

          const fields = fieldsList[focusFieldsIndex]

          fields.map(_field => {
            const code1Index = codes1.indexOf(_field.code)

            if (code1Index > -1) {
              const code2 = codes2[code1Index]
              const field2 = tools.find(fields, { code: code2 })

              if (/抢救室|ICU|烧伤病/.test(_field.resultValue)) {
                field2.disabled = false
              }
              else {
                field2.disabled = true
              }
            }
          })
        },

        // 46 
        set46DefaultValue({ codeValues = {}, fieldsList, focusFieldsIndex, op }) {
          const [firCodes, fifCodes, sixCodes] = [[], [], []]
          const { fc113 } = codeValues

          for (let _codes of codesList46) {
            firCodes.push(_codes[0])
            fifCodes.push(_codes[4])
            sixCodes.push(_codes[5])
          }

          const fields = fieldsList[focusFieldsIndex]

          for (let _field of fields) {
            if (firCodes.includes(_field.code)) {
              const index = tools.findIndex(firCodes, _field.code)
              const fifField = tools.find(fields, { code: fifCodes[index] }) || {}
              const sixField = tools.find(fields, { code: sixCodes[index] }) || {}

              if (_field.resultValue) {
                fifField[`${op}Value`] = '1'
                fifField.resultValue = '1'

                sixField[`${op}Value`] = fc113
                sixField.resultValue = fc113
              }
              else {
                fifField[`${op}Value`] = ''
                fifField.resultValue = ''

                sixField[`${op}Value`] = ''
                sixField.resultValue = ''
              }
            }
          }
        },

        // CSB0114RC02110000
        // disable211({ field, op }) {
        //   const codes = ['fc138', 'fc142', 'fc146', 'fc150', 'fc154', 'fc158', 'fc162', 'fc166']
        //   if (codes.includes(field.code)) {
        //     if (field.resultValue.includes('（') || field.resultValue.includes('）')) {
        //       field.resultValue = field.resultValue.replace('（', '(')
        //       field.resultValue = field.resultValue.replace('）', ')')
        //       field[`${op}Value`] = field.resultValue
        //     }
        //   }
        // },

        // CSB0114RC0213000
        // hintFc152({ field }) {
        //   const codes =  ['fc138', 'fc142', 'fc146', 'fc150', 'fc154', 'fc158', 'fc162', 'fc166']

        //   if (codes.includes(field.code)) {
        //     field.allowForce = true
        //     console.log(field.ocrBlur);
        //     if (field.items) {
        //       const result = field.items.find(text => text === field.resultValue)
        //       if(result) return true
        //     }
        //     if (field.ocrBlur && field.ocrBlur?.length != 0) {
        //       const result = field.ocrBlur.find(text => text === field.resultValue)
        //       // console.log(field.resultValue);
        //       console.log(result);
        //       if (!result) {
        //         // let dropItems = field.ocrBlur.slice(0, 2)
        //         let showItems = []
        //         let showItems1 = []

        //         for (let el of field.ocrBlur) {
        //           let valueArr = [...field.resultValue]
        //           let elArr = [...el]
        //           // 记录第一个不相同下标
        //           let flag1 = -1
        //           // 记录第二个不相同下标
        //           let flag2 = -1
        //           // 记录第三个不相同下标
        //           let flag3 = -1
        //           // 记录不相同个数
        //           let count = 0
        //           for (let i = 0; i < elArr.length; i++) {
        //             if (elArr[i] != valueArr[i]) {
        //               if (count == 0) flag1 = i
        //               if (count == 1) flag2 = i
        //               if (count == 2) flag3 = i
        //               count++
        //             }
        //           }

        //           if (count <= 3) {
        //             console.log('count',count);
        //             console.log('flag1',flag1, 'flag2',flag2, 'flag3',flag3);
        //             if (flag1 != -1 && flag2 == -1 && flag3 == -1) {
        //               let front1 = el.slice(0,flag1)
        //               let behind1 = el.slice(flag1 + 1)
        //               let str = `${front1}<span style="color: red;">${elArr[flag1]}</span>${behind1}`
        //               showItems.push(str)
        //               showItems1.push(str)
        //             }
        //             if (flag1 != -1 && flag2 != -1 && flag3 == -1) {
        //               let front1 = el.slice(0,flag1)
        //               let behind1 = el.slice(flag1 + 1, flag2)
        //               let behind2 = el.slice(flag2 + 1)
        //               let str = `${front1}<span style="color: red;">${elArr[flag1]}</span>${behind1}<span style="color: red;">${elArr[flag2]}</span>${behind2}`
        //               showItems.push(str)
        //             }
        //             if (flag1 != -1 && flag2 != -1 && flag3 != -1) {
        //               let front1 = el.slice(0,flag1)
        //               let behind1 = el.slice(flag1 + 1, flag2)
        //               let behind2 = el.slice(flag2 + 1, flag3)
        //               let behind3 = el.slice(flag3 + 1)
        //               let str = `${front1}<span style="color: red;">${elArr[flag1]}</span>${behind1}<span style="color: red;">${elArr[flag2]}</span>${behind2}<span style="color: red;">${elArr[flag3]}</span>${behind3}`
        //               showItems.push(str)
        //             }
        //           }
        //         }
        //         if (showItems1.length >= 2) showItems = showItems1
        //         console.log(showItems);
        //         if (!showItems[0]) {
        //           field.hint = ''
        //           return true
        //         }

        //         if (showItems[1]) {
        //           field.hint = `<p style="color: blue; fontSize: 14px; margin-top: -3px; margin-bottom: 0px"">${showItems[0]}，${showItems[1]}</p>`
        //           return false
        //         } else if(!showItems[1]) {
        //           field.hint = `<p style="color: blue; fontSize: 14px; margin-top: -3px; margin-bottom: 0px"">${showItems[0]}</p>`
        //           return false
        //         }
        //       }
        //     }

        //   }
        // },
      }
    },

    // 临时保存
    sessionSave: {
      methods: {
        // 36
        disable36({ fieldsList, focusFieldsIndex }) {
          const codesList = [
            ['fc130', 'fc131'],
            ['fc062', 'fc063'],
            ['fc064', 'fc065'],
            ['fc066', 'fc067'],
            ['fc068', 'fc069'],
            ['fc070', 'fc071'],
            ['fc072', 'fc073'],
            ['fc074', 'fc075'],
            ['fc076', 'fc077'],
            ['fc078', 'fc079'],
            ['fc080', 'fc081'],
            ['fc082', 'fc083'],
            ['fc084', 'fc085'],
            ['fc086', 'fc087'],
            ['fc088', 'fc089']
          ]

          const col2Codes = []

          codesList.map(codes => {
            col2Codes.push(codes[1])
          })

          const fields = fieldsList[focusFieldsIndex]

          const focusField = fields.find(field => field.autofocus)
          const codeIndex = col2Codes.indexOf(focusField.code)

          if (codeIndex > -1) {
            const restCodes = []
            let sliceIndex = -1

            for (let codesIndex in codesList) {
              if (codesList[codesIndex].includes(focusField.code)) {
                sliceIndex = +codesIndex + 1
                break
              }
            }

            const restCodesList = codesList.slice(sliceIndex)

            restCodesList.map(codes => {
              restCodes.push(...codes)
            })

            const restFields = fields.slice(focusField.fieldIndex + 1)

            restFields?.map(restField => {
              if (restCodes.includes(restField.code)) {
                restField.disabled = true
              }
            })
          }
        }
      }
    },

    // F8(提交前校验)
    beforeSubmit: {
      methods: {
        // 21
        validate21({ block, fieldsList }) {
          const { code, temp } = block
          if (temp !== 'MB002' || code !== 'bc008') return true

          const fields = fieldsList[0]
          const fc187Field = fields.find(field => field.code === 'fc187') || {}
          const fc188Field = fields.find(field => field.code === 'fc188') || {}
          const fc190Field = fields.find(field => field.code === 'fc190') || {}
          const fc191Field = fields.find(field => field.code === 'fc191') || {}

          if (fc187Field.resultValue) {
            if (!fc188Field.resultValue) return false
          }

          if (fc190Field.resultValue) {
            if (!fc191Field.resultValue) return false
          }

          return true
        },

        // 35
        validate35({ fieldsList }) {
          const fields = fieldsList[0]
          const fc132Field = tools.find(fields, { code: 'fc132' })
          const codes = ['fc131', 'fc063', 'fc065', 'fc067', 'fc069', 'fc071', 'fc073', 'fc075', 'fc077', 'fc079', 'fc081', 'fc083', 'fc085', 'fc087', 'fc089']

          if (!fc132Field) {
            return true
          }

          let fc132Value = fc132Field.resultValue

          if (fc132Value === 'A' || fc132Value.includes('?')) {
            fc132Value = 0
          }

          for (let field of fields) {
            if (codes.includes(field.code)) {
              if (field.resultValue === '?') {
                return true
              }
            }
          }

          let count = new BigNumber(0)

          for (let field of fields) {
            if (codes.includes(field.code)) {
              let fieldValue = field.resultValue

              if (!fieldValue || fieldValue.includes('?')) {
                fieldValue = 0
              }

              const resultValue = +field.resultValue
              count = count.plus(resultValue)
            }
          }

          const fc062Value = new BigNumber(+fc132Value)

          const diff = fc062Value.minus(count).toString()

          if (diff != 0) {
            return {
              popup: 'confirm',
              errorMessage: `明细金额与发票总金额${fc132Value}不一致，差额为${diff}，请确认并修改!`
            }
          }

          return true
        },

        // 47
        validate47({ block, fieldsList }) {
          const { code, temp } = block
          if (temp !== 'MB002' || code !== 'bc010') return true

          const fields = fieldsList[0]
          const fc138Field = fields.find(field => field.code === 'fc138') || {}

          if (!fc138Field.resultValue) {
            return false
          }

          return true
        },

        // CSB0114RC0210000
        validate210({ block, fieldsList }) {
          const { code, temp } = block
          if (temp !== 'MB002' || code !== 'bc010') return true

          const codes = ['fc138', 'fc142', 'fc146', 'fc150', 'fc154', 'fc158', 'fc162', 'fc166']
          // 是否有中文字符
          let pattern = /[\u4E00-\u9FFF\u3400-\u4DFF\uF900-\uFAFF]/;
          // 有几个中文字符
          let patterns = /[\u4E00-\u9FFF\u3400-\u4DFF\uF900-\uFAFF]/g;

          let count = 0
          for (let el of fieldsList) {
            for (let item of el) {
              if (item.hasOwnProperty('ID')) {
                count++
                break
              }
            }
          }
          // console.log(count);
          let fieldsArr
          if (fieldsList.length > count) {
            fieldsArr = fieldsList[fieldsList.length - 1]
            fieldsList.splice(fieldsList.length - 1)
          }

          console.log(fieldsList);
          for (let fields of fieldsList) {
            for (let field of fields) {
              if (codes.includes(field.code)) {
                if (field.resultValue == '') continue
                // value元素个数
                let count = [...field.resultValue]?.length
                // 中文字符个数
                let matches = field.resultValue.match(patterns)?.length;

                if (field.resultValue == 'DR' || field.resultValue == '?' || field.resultValue == '' || field.ctrlKey) continue
                else if (!pattern.test(field.resultValue)) {
                  fieldsList.push(fieldsArr)
                  return {
                    errorMessage: `${field.code}---${field.resultValue}录入内容有误，请检查`
                  }
                } else if (!pattern.test(field.resultValue) && count == 1) {
                  fieldsList.push(fieldsArr)
                  return {
                    errorMessage: `${field.code}---${field.resultValue}录入内容有误，请检查`
                  }
                }
                // else if (matches * 2 < count - matches) {
                //   fieldsList.push(fieldsArr)
                //   return {
                //     errorMessage: `${field.code}---${field.resultValue}录入内容有误，请检查`
                //   }
                // }
              }
            }
          }
          return true
        },
      }
    }
  }
}

export default B0114
