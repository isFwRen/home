import BigNumber from 'bignumber.js'
import { tools, sessionStorage } from 'vue-rocket'
import { ignoreFreeValue } from '../tools'
import moment from 'moment'
import { MessageBox, Notification } from 'element-ui';

const B0108 = {
  op0: {
    // 记录最后一次存储的合法field
    memoFields: [],

    // 记录相同 code 的 field 的值
    memoFieldValues: ['fc096', 'fc097', 'fc084', 'fc101', 'fc213', 'fc089', 'fc091', 'fc101', 'fc213', 'fc225'],

    // fields 的值从 targets 里的值选择
    dropdownFields: [
      {
        targets: ['fc096', 'fc097'],
        fields: ['fc089', 'fc091']
      },

      {
        targets: ['fc213'],
        fields: ['fc097']
      },

      {
        targets: ['fc225'],
        fields: ['fc097']
      }
    ],

    // 校验规则
    rules: [
      // 6
      {
        fields: ['fc096', 'fc097'],
        validate1: function ({ field, fieldsObject, thumbIndex, value }) {
          const [fc096Values, fc097Values] = [[], []]

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage

            if (sessionStorage || thumbIndex === +key) {
              const _fieldsList = fieldsObject[key].fieldsList

              for (let _fields of _fieldsList) {
                for (let _field of _fields) {
                  if (_field.code === 'fc096' && _field.uniqueId !== field.uniqueId) {
                    fc096Values.push(_field.resultValue)
                  }

                  if (_field.code === 'fc097' && _field.uniqueId !== field.uniqueId) {
                    fc097Values.push(_field.resultValue)
                  }
                }
              }
            }
          }

          if (fc096Values.includes(value) || fc097Values.includes(value)) {
            return '发票属性不能重复!'
          }

          return true
        }
      },

      // 7
      {
        fields: ['fc089', 'fc091'],
        validate2: function ({ includes, value }) {
          if (includes) {
            const result = includes.find(text => text === value)

            if (!result) {
              return '没有此发票，请核实!'
            }
          }

          return true
        }
      },

      // 11/13
      {
        fields: ['fc213', 'fc225'],
        validate3: function ({ fieldsObject, thumbIndex, value }) {
          const fc097Values = []

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage

            if (sessionStorage || thumbIndex === +key) {
              const _fieldsList = fieldsObject[key].fieldsList

              for (let _fields of _fieldsList) {
                for (let _field of _fields) {
                  if (_field.code === 'fc097') {
                    _field.resultValue && fc097Values.push(_field.resultValue)
                  }
                }
              }
            }
          }

          if (!fc097Values.includes(value)) {
            return '没有此发票住院发票，请核实!'
          }

          return true
        }
      },

      // 22/23/24
      {
        fields: ['fc084'],
        validate4: function ({ bill, field }) {
          if (bill?.saleChannel === '秒赔') {
            const value = field.op0Value

            switch (+value) {
              case 1:
              case 2:
                return '该案件不需要切“申请表”!';

              case 11:
                return '该案件不需要切“证件信息”!';

              case 12:
                return '该案件不需要“存折/银行卡”!';
            }
          }

          return true
        }
      },

      // 25
      {
        fields: ['fc084'],
        validate5: function ({ bill, field }) {
          // const agencies = ['83010', '83012', '83310', '83311']
          const agencies = ['00083000', '00083002', '00083300', '00083301']

          if (!agencies.includes(bill.agency)) {
            const value = field.op0Value

            if (value === '7')
              return '该案件发票不需要切清单!'
          }

          return true
        }
      },

      {
        fields: ['fc084'],
        validate6: function ({ bill, field }) {
          const agencies = ['00183000', '00183010', '00183002', '00183012', '00183300', '00183310', '00183301', '00183311']

          if (agencies.includes(bill.agency)) {
            const value = field.op0Value

            if (value === '8')
              return '该案件发票不需要切报销单!'
          }

          return true
        }
      },

      // CSB0108RC0014000
      {
        fields: ['fc084'],
        validate7: function ({ value, sameFieldValue }) {
          let fc084Values = sameFieldValue.fc084?.values
          if (!fc084Values) return
          let codesMap = new Map([
            ['1', '申请书重复录入'],
            ['2', '申请书重复录入'],
            // ['9', '诊断重复录入'],
            ['10', '死亡证明重复录入'],
            ['11', '证件信息重复录入'],
            ['12', '银行卡重复录入'],
            ['15', '病人姓名重复录入'],
          ])

          let duplicates = [];

          fc084Values = fc084Values.map(el => {
            if (el == '1' || el == '2') {
              return el = '1'
            } else {
              return el
            }
          })

          for (let i = 0; i < fc084Values.length; i++) {
            for (let j = i + 1; j < fc084Values.length; j++) {
              if (fc084Values[i] === fc084Values[j] && !duplicates.includes(fc084Values[i])) {
                duplicates.push(fc084Values[i]);
                break;
              }
            }
          }

          duplicates = [...new Set([...duplicates])]

          if (duplicates.includes(value) && codesMap.get(value)) {
            return codesMap.get(value)
          }
          if (sameFieldValue.fc084?.values.includes('1') && sameFieldValue.fc084?.values.includes('2')) {
            return codesMap.get('1')
          }
          return true
        }
      },

      // CSB0108RC0015000
      {
        fields: ['fc101'],
        validate8: function ({ value, sameFieldValue }) {
          let fc101Values = sameFieldValue.fc101?.values

          let count = 0
          for (let el of fc101Values) {
            if (el == value) count++
            if (count >= 2) return '手术属性不能重复'
          }
          return true
        }
      },

      // CSB0108RC0261000
      {
        fields: ['fc213'],
        validate9: function ({ value, sameFieldValue }) {
          let fc213Values = sameFieldValue.fc213?.values

          let count = 0
          for (let el of fc213Values) {
            if (el == value) count++
            if (count >= 2) return '住院日期属性不能重复'
          }
          return true
        }
      },

      // CSB0108RC0309000
      {
        fields: ['fc101'],
        validate10: function ({ field, fieldsObject, thumbIndex, value }) {
          const Values = []

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage

            if (sessionStorage || thumbIndex === +key) {
              const _fieldsList = fieldsObject[key].fieldsList

              for (let _fields of _fieldsList) {
                for (let _field of _fields) {
                  if ((_field.code === 'fc096' || _field.code === 'fc097') && _field.uniqueId !== field.uniqueId) {
                    Values.push(_field.resultValue)
                  }
                }
              }
            }
          }

          if (!Values.includes(value)) {
            return '没有此发票，请核实'
          }

          return true
        }
      },
    ],

    // 提示文本
    hints: [
      // CSB0108RC0262000
      {
        fields: ['fc084'],
        validate20: function ({ field, bill }) {
          const agency = bill?.agency
          const agencies = ['00083000', '00083002', '00083300', '00083301']

          if (field.resultValue == '5' || field.resultValue == '6') {
            if (agencies.includes(agency)) {
              return '<p style="color: red; margin-top: -3px; margin-bottom: 0px; font-weight:1000; font-size:16px">该案件需要切清单</p>'
            }
            if (!agencies.includes(agency)) {
              return '<p style="color: red; margin-top: -3px; margin-bottom: 0px; font-weight:1000; font-size:16px">该案件不需要切清单</p>'
            }
          }
        }
      },
    ],

    // 工序完成初始化
    init: {
      methods: {
        validate18({ bill }) {
          const { agency } = bill
          // const agencies1 = ['83000', '83002', '83300', '83301']
          // const agencies2 = ['83010', '83012', '83310', '83311']

          const agencies1 = ['00083000', '00083002', '00083300', '00083301']
          const agencies2 = ['00083000', '00083002', '00083300', '00083301']

          if (agencies1.includes(agency)) {
            if (sessionStorage.get('isApp')?.isApp === 'true') {
              // MessageBox.alert('该案件发票需要切清单!', '请注意', {
              //   type: 'warning',
              //   confirmButtonText: '确定',
              //   showClose: false,
              // })
              return Notification({
                type: 'warning',
                title: '提醒(5s后自动关闭)',
                message: '该案件发票需要切清单!',
                duration: 5000,
                position: 'top-left'
              })
            } else {
              alert('该案件发票需要切清单!')
            }
          }

          if (!agencies2.includes(agency)) {
            if (sessionStorage.get('isApp')?.isApp === 'true') {
              // MessageBox.alert('该案件发票不需要切清单!', '请注意', {
              //   type: 'warning',
              //   confirmButtonText: '确定',
              //   showClose: false,
              // })
              return Notification({
                type: 'warning',
                title: '提醒(5s后自动关闭)',
                message: '该案件发票不需要切清单!',
                duration: 5000,
                position: 'top-left'
              })
            } else {
              alert('该案件发票不需要切清单!')
            }
          }
        }
      }
    },

    // 字段已生成
    updateFields: {
    },

    // 回车
    enter: {
      methods: {
        // 3/4/5
        disableFields({ fieldsList, focusFieldsIndex }) {
          const codesList = [
            ['fc189'],
            ['fc003', 'fc004'],
            ['fc217']
          ]

          const flatCodesList = tools.flatArray(codesList)
          const fields = fieldsList[focusFieldsIndex]

          fields?.map(_field => {
            if (flatCodesList.includes(_field.code)) {
              _field.disabled = true
            }
          })
        },

        disableFields01({ fieldsList, focusFieldsIndex }) {
          const codesList = [
            ['fc259', 'fc260']
          ]

          const flatCodesList = tools.flatArray(codesList)
          const fields = fieldsList[focusFieldsIndex]

          fields?.map(_field => {
            if (flatCodesList.includes(_field.code)) {
              _field.disabled = true
            }
          })
        },

        disableFields02({ fieldsList, focusFieldsIndex }) {
          const codesList = [
            ['fc261', 'fc262', 'fc264', 'fc265', 'fc266'],
          ]

          const flatCodesList = tools.flatArray(codesList)
          const fields = fieldsList[focusFieldsIndex]

          fields?.map(_field => {
            if (flatCodesList.includes(_field.code)) {
              _field.disabled = true
            }
          })
        },
      }
    },

    // 临时保存
    sessionSave: {
      methods: {
        validate14: function ({ sameFieldValue }) {
          let fc084Values = sameFieldValue.fc084?.values
          if (!fc084Values) return
          let codesMap = new Map([
            ['1', '申请书重复录入,保存失败'],
            ['2', '申请书重复录入,保存失败'],
            // ['9', '诊断重复录入,保存失败'],
            ['10', '死亡证明重复录入,保存失败'],
            ['11', '证件信息重复录入,保存失败'],
            ['12', '银行卡重复录入,保存失败'],
            ['15', '病人姓名重复录入,保存失败'],
          ])

          let duplicates = [];

          fc084Values = fc084Values.map(el => {
            if (el == '1' || el == '2') {
              return el = '1'
            } else {
              return el
            }
          })

          for (let i = 0; i < fc084Values.length; i++) {
            for (let j = i + 1; j < fc084Values.length; j++) {
              if (fc084Values[i] === fc084Values[j] && !duplicates.includes(fc084Values[i])) {
                duplicates.push(fc084Values[i]);
                break;
              }
            }
          }

          duplicates = [...new Set([...duplicates])]

          if (duplicates == []) return true
          else {
            for (let el of duplicates) {
              if (codesMap.get(el)) {
                return {
                  errorMessage: codesMap.get(el)
                }
              }
            }
          }
          return true
        },

        // CSB0108RC0263000
        validate263: function ({ fieldsList }) {
          let fc101 = ''
          let fc096 = ''
          let fc140 = ''
          for (let fields of fieldsList) {
            for (let field of fields) {
              if (field.code == 'fc101') {
                fc101 = field.resultValue
              }
              if (field.code == 'fc096') {
                fc096 = field.resultValue
              }
              if (field.code == "fc140" && field.resultValue == '2') {
                fc140 = '2'
              }
            }
          }

          if (fc101 && fc096) {
            if (fc140 == '2' && fc101 == fc096) {
              return {
                errorMessage: '保存失败，该发票无手术， 无需切手术'
              }
            }
          }

          return true
        },

        // CSB0108RC0264000
        validate264: function ({ fieldsList }) {
          let fc142 = ''
          let fc101 = ''
          let fc097 = ''
          for (let fields of fieldsList) {
            for (let field of fields) {
              if (field.code == 'fc101') fc101 = field.resultValue
              if (field.code == 'fc097') fc097 = field.resultValue
              if (field.code == "fc142" && field.resultValue == '2') fc142 = '2'
            }
          }
          if (fc101 && fc097) {
            if (fc142 == '2' && fc101 == fc097) {
              return {
                errorMessage: '保存失败，该发票无手术， 无需切手术'
              }
            }
          }
          return true
        },

        // CSB108RC0265000
        validate265: function ({ bill, sameFieldValue }) {
          let fc084Values = sameFieldValue.fc084?.values
          if (!fc084Values) return
          if (bill.saleChannel != '秒赔') return
          let codesMap = new Map([
            ['1', '申请书'],
            ['2', '申请书'],
            ['10', '死亡证明'],
            ['11', '证件信息'],
            ['12', '银行卡'],
          ])
          let codes = ['1', '2', '10', '11', '12']

          for (let el of fc084Values) {
            if (codes.includes(el)) {
              return {
                errorMessage: `秒赔案件不需要切${codesMap.get(el)}`
              }
            }
          }
          return true
        },

        validate293({ fieldsList }) {
          const codes = ['fc096', 'fc097', 'fc089', 'fc091', 'fc101', 'fc213', 'fc225']
          for (let fields of fieldsList) {
            for (let field of fields) {
              if (codes.includes(field.code) && field.resultValue == '') {
                return {
                  errorMessage: `${field.name}属性不能为空，保存失败，请检查!`
                }
              }
            }
          }

          return true
        },
      }
    },

    // F8(提交前校验)
    beforeSubmit: {
      methods: {
        // 8
        validate8({ fieldsObject }) {
          const [fc089Values, fc096Values, fc097Values] = [[], [], []]

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  const { code, resultValue } = field

                  if (code === 'fc089') {
                    resultValue && fc089Values.push(resultValue)
                  }

                  if (code === 'fc096') {
                    resultValue && fc096Values.push(resultValue)
                  }

                  if (code === 'fc097') {
                    resultValue && fc097Values.push(resultValue)
                  }
                }
              }
            }
          }

          const mergeValues = [...fc096Values, ...fc097Values]

          for (let value of fc089Values) {
            if (!mergeValues.includes(value)) {
              return {
                errorMessage: `清单明细${value}没有匹配的发票，请检查!`
              }
            }
          }

          return true
        },

        // 9
        validate9({ bill, fieldsObject }) {
          // const secondLast = bill?.agency?.split('').reverse()[1]
          const agency = bill?.agency

          if (agency == '00083000' || agency == '00083002' || agency == '00083300' || agency == '00083301') {
            const [fc089Values, fc096Values, fc097Values] = [[], [], []]

            for (let key in fieldsObject) {
              const sessionStorage = fieldsObject[key].sessionStorage
              const fieldsList = fieldsObject[key].fieldsList

              if (sessionStorage) {
                for (let fields of fieldsList) {
                  for (let field of fields) {
                    const { code, resultValue } = field

                    if (code === 'fc089') {
                      resultValue && fc089Values.push(resultValue)
                    }

                    if (code === 'fc096') {
                      resultValue && fc096Values.push(resultValue)
                    }

                    if (code === 'fc097') {
                      resultValue && fc097Values.push(resultValue)
                    }
                  }
                }
              }
            }

            const mergeValues = [...fc096Values, ...fc097Values]

            for (let value of mergeValues) {
              if (!fc089Values.includes(value)) {
                return {
                  errorMessage: `发票属性为${value}的发票漏切清单!`
                }
              }
            }
          }

          return true
        },

        // 10
        validate10({ fieldsObject }) {
          const [fc091Values, fc096Values, fc097Values] = [[], [], []]

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  if (field.code === 'fc091') {
                    fc091Values.push(field.resultValue)
                  }

                  if (field.code === 'fc096') {
                    fc096Values.push(field.resultValue)
                  }

                  if (field.code === 'fc097') {
                    fc097Values.push(field.resultValue)
                  }
                }
              }
            }
          }

          for (let value of fc091Values) {
            if (![...fc096Values, ...fc097Values].includes(value)) {
              return {
                errorMessage: `报销单${value}没有匹配的发票!`
              }
            }
          }

          // console.log({ fc091Values, fc096Values, fc097Values })

          return true
        },

        // 12
        validate12({ fieldsObject }) {
          const [fc097Values, fc213Values] = [[], []]

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  const { code, resultValue } = field

                  if (code === 'fc097') {
                    resultValue && fc097Values.push(resultValue)
                  }

                  if (code === 'fc213') {
                    resultValue && fc213Values.push(resultValue)
                  }
                }
              }
            }
          }

          for (let value of fc097Values) {
            if (!fc213Values.includes(value)) {
              return {
                errorMessage: `住院发票${value}没有匹配的住院日期,请检查!`
              }
            }
          }

          return true
        },

        // 11/13
        validate11_13({ fieldsObject }) {
          const [fc213Values, fc225Values, fc097Values] = [[], [], []]

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  if (field.code === 'fc213') {
                    field.resultValue && fc213Values.push(field.resultValue)
                  }

                  if (field.code === 'fc225') {
                    field.resultValue && fc225Values.push(field.resultValue)
                  }

                  if (field.code === 'fc097') {
                    field.resultValue && fc097Values.push(field.resultValue)
                  }
                }
              }
            }
          }

          for (let value of fc213Values) {
            if (!fc097Values.includes(value)) {
              return {
                errorMessage: '没有此住院发票，请核实!'
              }
            }
          }

          for (let value of fc225Values) {
            if (!fc097Values.includes(value)) {
              return {
                errorMessage: '没有此住院发票，请核实!'
              }
            }
          }

          return true
        },

        // 14
        validate14({ fieldsObject }) {
          const [fc097Values, fc225Values] = [[], []]

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  const { code, resultValue } = field

                  if (code === 'fc097') {
                    resultValue && fc097Values.push(resultValue)
                  }

                  if (code === 'fc225') {
                    resultValue && fc225Values.push(resultValue)
                  }
                }
              }
            }
          }

          for (let value of fc097Values) {
            if (!fc225Values.includes(value)) {
              return {
                errorMessage: `住院发票${value}没有匹配的诊查费天数,请检查!`
              }
            }
          }

          return true
        },

        // 15
        validate15({ fieldsObject, bill }) {
          const agency = bill.agency
          const agencies = ['00183000', '00183010', '00183002', '00183012', '00183300', '00183310', '00183301', '00183311']
          if (agencies.includes(agency)) return
          const fc084Values = []

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  const { code, resultValue } = field

                  if (code === 'fc084') {
                    resultValue && fc084Values.push(resultValue)
                  }
                }
              }
            }
          }

          const has5 = fc084Values.includes('5')
          const has6 = fc084Values.includes('6')

          if (has5 || has6) {
            const has7 = fc084Values.includes('7')
            const has8 = fc084Values.includes('8')

            if (!has7 && !has8) {
              return {
                popup: 'confirm',
                errorMessage: '发票缺少对应清单或报销单,请确认!'
              }
            }
          }

          return true
        },

        // 16
        validate16({ block, fieldsObject }) {
          // if(block?.code.toLowerCase() === 'bc002') {
          const fc101Values = []

          for (let key in fieldsObject) {
            const fieldsList = fieldsObject[key].fieldsList

            for (let fields of fieldsList) {
              for (let field of fields) {
                // fc101
                if (field?.code === 'fc101') {
                  field.resultValue && fc101Values.push(field.resultValue)
                }
              }
            }
          }

          for (let key in fieldsObject) {
            const fieldsList = fieldsObject[key].fieldsList

            for (let fields of fieldsList) {
              const fc140Field = tools.find(fields, { code: 'fc140' })

              if (fc140Field?.resultValue === '1') {
                const fc096Field = tools.find(fields, { code: 'fc096' })
                const fc096ResultValue = fc096Field.resultValue

                if (!fc101Values.includes(fc096ResultValue)) {
                  return {
                    errorMessage: `门诊发票${fc096ResultValue}存在手术费用，漏切手术，请确认!`
                  }
                }
              }
            }
          }

          // }

          return true
        },

        // 17
        validate17({ block, fieldsObject }) {
          // if(block?.code.toLowerCase() === 'bc003') {
          const fc101Values = []

          for (let key in fieldsObject) {
            const fieldsList = fieldsObject[key].fieldsList

            for (let fields of fieldsList) {
              // fc101
              const fc101Field = tools.find(fields, { code: 'fc101' })
              fc101Field?.resultValue && fc101Values.push(fc101Field.resultValue)
            }
          }

          for (let key in fieldsObject) {
            const fieldsList = fieldsObject[key].fieldsList

            for (let fields of fieldsList) {
              const fc142Field = tools.find(fields, { code: 'fc142' })

              if (fc142Field?.resultValue === '1') {
                const fc097Field = tools.find(fields, { code: 'fc097' })
                const fc097ResultValue = fc097Field.resultValue

                if (!fc101Values.includes(fc097ResultValue)) {
                  return {
                    errorMessage: `住院发票${fc097ResultValue}存在手术费用，漏切手术，请确认!`
                  }
                }
              }
            }
          }

          // }

          return true
        },

        // 20
        validate20({ fieldsObject }) {
          const [values1, values2, values9, values10, values11, values12] = [[], [], [], [], [], []]

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  if (field.code === 'fc084') {
                    const resultValue = field.resultValue

                    switch (resultValue) {
                      case '1':
                        values1.push('1')
                        break;
                      case '2':
                        values2.push('2')
                        break;
                      case '9':
                        values9.push('9')
                        break;
                      case '10':
                        values10.push('10')
                        break;
                      case '11':
                        values11.push('11')
                        break;
                      case '12':
                        values12.push('12')
                        break;
                    }
                  }
                }
              }
            }
          }

          if (values1.length > 1) {
            return {
              errorMessage: '重复录入申请书，请核查!'
            }
          }

          if (values2.length > 1) {
            return {
              errorMessage: '重复录入申请书，请核查!'
            }
          }

          if (values1.length != 0 && values2.length != 0) {
            return {
              errorMessage: '重复录入申请书，请核查!'
            }
          }

          // if (values9.length > 1) {
          //   return {
          //     errorMessage: '重复录入诊断书，请核查!'
          //   }
          // }

          if (values10.length > 1) {
            return {
              errorMessage: '重复录入死亡证明书，请核查!'
            }
          }

          if (values11.length > 1) {
            return {
              errorMessage: '重复录入证件信息，请核查!'
            }
          }

          if (values12.length > 1) {
            return {
              errorMessage: '重复录入存折/银行卡，请核查!'
            }
          }

          return true
        },

        // 21
        validate21({ bill, fieldsObject }) {
          const firstTwo = bill.billNum?.slice(0, 2)

          const values = []

          if (firstTwo === '00') {
            for (let key in fieldsObject) {
              const sessionStorage = fieldsObject[key].sessionStorage
              const fieldsList = fieldsObject[key].fieldsList

              if (sessionStorage) {
                for (let fields of fieldsList) {
                  for (let field of fields) {
                    const { code, resultValue } = field

                    if (code === 'fc084') {
                      resultValue && values.push(resultValue)
                    }
                  }
                }
              }
            }

            if (!values.includes('1') && !values.includes('2')) {
              return {
                errorMessage: '缺少申请书，请确认!'
              }
            }

            if (!values.includes('11')) {
              return {
                errorMessage: '缺少证件信息，请确认!'
              }
            }

            if (!values.includes('12')) {
              return {
                errorMessage: '缺少存折/银行卡，请确认!'
              }
            }
          }

          return true
        },

        // 26
        validate26({ bill, fieldsObject }) {
          const { agency } = bill
          // const agencies = ['83000', '83002', '83300', '83301']
          const agencies = ['00083000', '00083002', '00083300', '00083301']

          if (agencies.includes(agency)) {
            const fc084Values = []

            for (let key in fieldsObject) {
              const sessionStorage = fieldsObject[key].sessionStorage
              const fieldsList = fieldsObject[key].fieldsList

              if (sessionStorage) {
                for (let fields of fieldsList) {
                  const fc084Field = tools.find(fields, { code: 'fc084' })
                  fc084Field?.resultValue && fc084Values.push(fc084Field.resultValue)
                }
              }
            }

            if (!fc084Values.includes('5') && !fc084Values.includes('6')) {
              return true
            }

            if (!fc084Values.includes('7')) {
              return {
                errorMessage: '该案件没有录入清单，请确认!'
              }
            }
          }

          return true
        },

        // 53
        validate53({ flatFieldsList }) {
          let boolean = true;
          const rules = [
            'fc080',
            'fc081',
            'fc143',
            'fc144',
            'fc145',
            'fc146',
            'fc147',
            'fc148',
            'fc149',
            'fc150',
            'fc151',
            'fc152',
            'fc153',
            'fc154',
            'fc155',
            'fc156',
          ];
          let codeObjList = {};
          flatFieldsList.forEach((e) => {
            codeObjList.code = e;
          });
          for (let key = 0; key < rules.length; key += 2) {
            if (
              (codeObjList[rules[key]] && codeObjList[rules[key + 1]]) ||
              !codeObjList[rules[key]]
            ) {
              continue;
            } else {
              boolean = false;
            }
          }
          return (
            boolean || {
              errorMessage: `清单内容录入遗漏，请检查`,
            }
          );
        },

        validate19({ mergeFieldsList, op }) {
          const flatFieldsList = tools.flatArray(mergeFieldsList)
          const fc084Field = tools.find(flatFieldsList, { code: 'fc084' })

          if (fc084Field) {
            const fc084Values = []

            for (let field of flatFieldsList) {
              if (field.code === 'fc084') {
                fc084Values.push(field[`${op}Value`])
              }
            }

            if (fc084Values.includes('9')) {
              return true
            }
            else {
              return {
                errorMessage: '缺少诊断书，请确认!'
              }
            }
          }

          return {
            errorMessage: '缺少诊断书，请确认!'
          }
        },

        validate191({ mergeFieldsList, op }) {
          const flatFieldsList = tools.flatArray(mergeFieldsList)
          const fc084Field = tools.find(flatFieldsList, { code: 'fc084' })

          if (fc084Field) {
            const fc084Values = []

            for (let field of flatFieldsList) {
              if (field.code === 'fc084') {
                fc084Values.push(field[`${op}Value`])
              }
            }

            if (fc084Values.includes('15')) {
              return true
            }
            else {
              return {
                errorMessage: '缺少病人姓名，请确认!'
              }
            }
          }

          return true
        },

        validate265: function ({ bill, sameFieldValue }) {
          let fc084Values = sameFieldValue.fc084?.values
          if (!fc084Values) return
          if (bill.saleChannel != '秒赔') return
          let codesMap = new Map([
            ['1', '申请书'],
            ['2', '申请书'],
            ['10', '死亡证明'],
            ['11', '证件信息'],
            ['12', '银行卡'],
          ])
          let codes = ['1', '2', '10', '11', '12']

          for (let el of fc084Values) {
            if (codes.includes(el)) {
              return {
                errorMessage: `秒赔案件不需要切${codesMap.get(el)}`
              }
            }
          }
          return true
        },

        validate292({ fieldsObject }) {
          const [fc096Values, fc097Values] = [[], []]

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  const { code, resultValue } = field

                  if (code === 'fc096') {
                    resultValue && fc096Values.push(resultValue)
                  }

                  if (code === 'fc097') {
                    resultValue && fc097Values.push(resultValue)
                  }
                }
              }
            }
          }

          let fc096 = new Set(fc096Values);
          let fc097 = new Set(fc097Values);
          let newSet = new Set(
            [...fc096].filter(item => fc097.has(item))
          );

          let arr = [...newSet]
          let duplicatesFc096 = fc096Values.filter((item, index) => fc096Values.indexOf(item) !== index);
          let duplicatesFc097 = fc097Values.filter((item, index) => fc097Values.indexOf(item) !== index);


          if (arr[0]) {
            console.log(arr[0]);
            return {
              errorMessage: `发票属性${arr[0]}重复，不能提交，请修改!`
            }
          }
          if (duplicatesFc096[0]) {
            console.log(duplicatesFc096[0]);
            return {
              errorMessage: `发票属性${duplicatesFc096[0]}重复，不能提交，请修改!`
            }
          }

          if (duplicatesFc097[0]) {
            console.log(duplicatesFc097[0]);
            return {
              errorMessage: `发票属性${duplicatesFc097[0]}重复，不能提交，请修改!`
            }
          }

          return true
        },

        validate293({ fieldsObject }) {
          const codes = ['fc096', 'fc097', 'fc089', 'fc091', 'fc101', 'fc213', 'fc225']

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  if (codes.includes(field.code) && field.resultValue == '') {
                    return {
                      errorMessage: `${field.name}属性不能为空，不能提交，请检查!`
                    }
                  }
                }
              }
            }
          }

          return true
        },

        validate299({ bill, sameFieldValue }) {
          // const len = bill?.agency?.length - 2
          // const num = bill?.agency?.charAt(len)
          const agency = bill?.agency
          const agencies = ['00083000', '00083002', '00083300', '00083301']

          if (!agencies.includes(agency)) {
            let fc084Values = sameFieldValue.fc084?.values

            if (fc084Values.includes('7')) {
              return {
                errorMessage: '该案件发票不需要切清单！请注意字段拦截提示！！！'
              }
            }
          }

          return true
        },

        validate300({ bill, sameFieldValue }) {
          const agency = bill?.agency
          const agencies = ['00183000', '00183010', '00183002', '00183012', '00183300', '00183310', '00183301', '00183311']
          console.log(agency);
          if (agencies.includes(agency)) {
            let fc084Values = sameFieldValue.fc084?.values

            if (fc084Values.includes('8')) {
              return {
                errorMessage: '该案件发票不需要切报销单！请注意字段拦截提示！！！'
              }
            }
          }

          return true
        },
      }
    }
  },

  op1op2opq: {
    // 校验规则
    rules: [
      // 43
      {
        fields: ['fc165', 'fc196', 'fc197', 'fc198', 'fc199'],
        validate43: function ({ value, items }) {
          if (value === '?') {
            return true
          }

          const result = items.find((text) => text === value)

          if (!result) {
            return '录入内容不在数据库中，请确认!'
          }

          return true
        }
      },

      // 47
      {
        fields: ['fc113', 'fc114', 'fc115', 'fc116', 'fc117', 'fc118', 'fc119', 'fc120', 'fc121', 'fc122', 'fc123', 'fc124', 'fc125', 'fc126', 'fc127', 'fc128', 'fc129', 'fc130', 'fc131', 'fc132'],
        validate47: function ({ field, value, items }) {
          if (!value) return true

          const result = items.find((text) => text === value)

          if (value.includes('?')) {
            field.allowForce = true
            return true
          }
          else {
            field.allowForce = false
          }

          if (!result) {
            field.allowForce = false
            return '录入内容不在数据库中，请确认!'
          }

          field.allowForce = true

          return true
        }
      },

      // 67
      {
        fields: ['fc030', 'fc190'],
        validate67: function ({ effectValidations, field, value, items }) {
          if (ignoreFreeValue({ effectValidations, value })) return true
          const result = items.find(text => text === value)

          if (value.includes('?')) {
            field.allowForce = true
            return true
          }
          else if (value == 'B') {
            return true
          }
          else {
            field.allowForce = false
          }

          if (!result) {
            field.allowForce = false
            return '录入内容不在数据库中，请确认!'
          }

          field.allowForce = true

          return true
        }
      },

      // 35
      {
        fields: ['fc034', 'fc055'],
        validate35: function ({ field, value, items }) {
          if (!value) return true

          const result = items.find(text => text === value)

          if (value.includes('?')) {
            field.allowForce = true
            return true
          }
          else {
            field.allowForce = false
          }

          if (!result) {
            return '录入内容不在数据库中，请确认!'
          }

          field.allowForce = true

          return true
        }
      },

      {
        fields: ['fc018', 'fc112', 'fc036', 'fc104', 'fc058', 'fc161', 'fc162', 'fc205', 'fc206', 'fc207', 'fc208', 'fc209', 'fc239', 'fc240', 'fc241', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246', 'fc247', 'fc248'],
        validatefc018: function ({ effectValidations, field, items, value }) {
          if (ignoreFreeValue({ effectValidations, value })) return true

          const result = items.find((text) => text === value)

          if (value.includes('?')) {
            field.allowForce = true
            return true
          }
          else {
            field.allowForce = false
          }

          if (!result) {
            return '录入内容不在数据库中，请确认!'
          }

          field.allowForce = true

          return true
        }
      },

      {
        fields: ['fc204'],
        validate41: function ({ value, fieldsList }) {
          const rules = {
            0: 'fc164、fc165、fc196、fc200、fc197、fc201、fc198、fc202、fc199、fc203'.split(
              `、`
            ),
            1: 'fc196、fc200、fc197、fc201、fc198、fc202、fc199、fc203'.split(
              `、`
            ),
            2: 'fc197、fc201、fc198、fc202、fc199、fc203'.split(`、`),
            3: 'fc198、fc202、fc199、fc203'.split(`、`),
            4: 'fc199、fc203'.split(`、`),
          };
          //对象化
          let allObj = {};
          for (let key in rules) {

            rules[key].forEach((e) => {
              allObj[e] = 1;
            });
          }
          var flatFieldsList = tools.flatArray(fieldsList)
          flatFieldsList && flatFieldsList.forEach((e) => {
            if (allObj[e.code]) {
              e.disabled = false;
            }
            if (rules[value] && rules[value].includes(e.code)) {

              e.disabled = true;
            }
          });
          return true;
        },
      },

      // CSB0108RC0237000
      {
        fields: ['fc019'],
        validate26: function ({ field, value }) {

          if (value === 'A') {
            field.allowForce = false
            return '如无勾选，则录入问号，不可录入A'
          }

          if (!value.includes('/')) {
            field.allowForce = false
            let number = Number(value)
            if (number < 10 || number > 20) {
              return '录入值不在选项中'
            }
          }

          if (value.includes('/')) {
            field.allowForce = false
            let arr = value.split('/')

            for (let el of arr) {
              let number = Number(el)
              if (number < 10 || number > 20) {
                return '录入值不在选项中'
              }
            }
          }

          field.allowForce = true
          return true
        }
      },

      // CSB0108RC0071000
      {
        fields: ['fc017', 'fc035', 'fc056', 'fc057', 'fc099', 'fc158', 'fc164', 'fc200', 'fc201', 'fc202', 'fc203', 'fc214', 'fc215', 'fc217'],
        validateDate02: function ({ value }) {
          if (!value) return true

          if (value == 'A') return true

          if (/[A, \?]/.test(value)) {
            return true
          }

          if (value.length !== 6 || moment(`20${value}`).format('YYMMDD') === 'Invalid date') {
            return '日期格式录入错误'
          }

          return true
        }
      },
    ],

    // 提示文本
    hints: [
      // CSB0108RC0243000
      {
        fields: ['fc226', 'fc227'],
        hintFc01: function () {
          return '<p style="color: red; margin-top: -3px; margin-bottom: 0px">优先录入盖章医院全称，如盖章模糊可录入印刷体医院名称，注意不要录入错别字</p>'
        }
      },
      {
        fields: ['fc056'],
        hintFc02: function () {
          return '<p style="color: red; margin-top: -3px; margin-bottom: 0px">录入发票中最早入院日期</p>'
        }
      },
      {
        fields: ['fc041', 'fc064'],
        hintFc03: function () {
          return '<p style="color: red; margin-top: -3px; margin-bottom: 0px">录入”手写标注“其他扣除或者第三方保险公司的理赔金额</p>'
        }
      },
      {
        fields: ['fc035'],
        hintFc04: function () {
          return '<p style="color: red; margin-top: -3px; margin-bottom: 0px">录入发票中最早就诊日期</p>'
        }
      },
      {
        fields: ['fc001', 'fc002', 'fc013'],
        hintFc05: function () {
          return '<p style="color: red; margin-top: -3px; margin-bottom: 0px">只录入"非上海地区"，”手写标注“的自费金额</p>'
        }
      },
      {
        fields: ['fc014', 'fc015', 'fc016'],
        hintFc06: function () {
          return '<p style="color: red; margin-top: -3px; margin-bottom: 0px">只录入上海发票自费金额</p>'
        }
      },
      {
        fields: ['fc045', 'fc082', 'fc083'],
        hintFc07: function () {
          return '<p style="color: red; margin-top: -3px; margin-bottom: 0px">只录入固定模板：四川、陕西、重庆、宁波自费金额</p>'
        }
      },
      {
        fields: ['fc090', 'fc105', 'fc106'],
        hintFc08: function () {
          return '<p style="color: red; margin-top: -3px; margin-bottom: 0px">录入影像中的统筹金额、报销金额、基金支付、暂缓统筹、医保保险支付、挂账统筹金额</p>'
        }
      },
      {
        fields: ['fc111', 'fc138', 'fc139'],
        hintFc09: function () {
          return '<p style="color: red; margin-top: -3px; margin-bottom: 0px">录入各种补助类金额，如其他支付、大病支付、大额支付、公务员补助、新农合、超限垫支、医院承担、医院负担等金额</p>'
        }
      },
      {
        fields: ['fc036', 'fc104', 'fc058', 'fc161', 'fc162', 'fc205', 'fc206', 'fc207', 'fc208', 'fc209', 'fc239', 'fc240', 'fc241', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246', 'fc247', 'fc248'],
        hintFc10: function () {
          return '<p style="color: red; margin-top: -3px; margin-bottom: 0px;font-weight:bolder">注意：肉瘤、癌、恶性是恶性肿瘤；肿瘤/腺瘤/原位癌/瘤变/上皮病变/占位性病变是良性肿瘤</p>'
        }
      },
      {
        fields: ['fc035'],
        hintFc11: function () {
          return '<p style="color: red; margin-top: -3px; margin-bottom: 0px;font-weight:bolder">打印时间不是就诊日期,不要录入</p>'
        }
      },
      {
        fields: ['fc057'],
        hintFc12: function () {
          return '<p style="color: red; margin-top: -3px; margin-bottom: 0px;font-weight:bolder">开票日期不是出院日期，也不是入院日期，不要录入；</p>'
        }
      },
      {
        fields: ['fc062'],
        hintFc13: function () {
          return '<p style="color: red; margin-top: -3px; margin-bottom: 0px;font-weight:bolder">增值税发票：金额需要录入“金额和税额”的合计</p>'
        }
      },
      {
        fields: ['fc183'],
        hintFc14: function () {
          return '<p style="color: red; margin-top: -3px; margin-bottom: 0px;font-weight:bolder">未填写/模糊不需要录入</p>'
        }
      },
      {
        fields: ['fc177'],
        hintFc15: function () {
          return '<p style="color: red; margin-top: -3px; margin-bottom: 0px;font-weight:bolder">未填写/模糊不需要录入</p>'
        }
      },
    ],

    // 字段已生成
    updateFields: {
      methods: {
        // 28
        disable28({ bill, fieldsList }) {
          if (bill?.saleChannel !== '秒赔') return

          const codes = ['fc009', 'fc019', 'fc099', 'fc010', 'fc011', 'fc095', 'fc012', 'fc027', 'fc028', 'fc029', 'fc030', 'fc031', 'fc158', 'fc112', 'fc169', 'fc170', 'fc171', 'fc172', 'fc173', 'fc110', 'fc190', 'fc191', 'fc174', 'fc175']

          fieldsList?.map(fields => {
            fields?.map(_field => {
              if (codes.includes(_field.code)) {
                _field.disabled = true
              }
            })
          })
        },

        // 48
        disable48({ fieldsList }) {
          const codes = ['fc044', 'fc046', 'fc047', 'fc048', 'fc049', 'fc050', 'fc051', 'fc052', 'fc053', 'fc100', 'fc102']

          fieldsList?.map(fields => {
            fields?.map(_field => {
              if (codes.includes(_field.code)) {
                _field.disabled = true
              }
            })
          })
        },

        // 51
        disable51({ block, fieldsList }) {
          if (block?.code.toLowerCase() === 'bc004') {

            fieldsList?.map(fields => {
              fields?.map(_field => {
                if (_field.code === 'fc188') {
                  _field.disabled = true
                }
              })
            })
          }
        },

        // 63
        disable63({ fieldsList }) {
          const codes = ['fc066', 'fc067', 'fc068', 'fc069', 'fc070', 'fc071', 'fc072', 'fc073', 'fc074', 'fc075', 'fc076', 'fc077', 'fc078', 'fc079']

          fieldsList.map(fields => {
            fields?.map(_field => {
              if (codes.includes(_field.code)) {
                _field.disabled = true
              }
            })
          })
        },

        // 35
        notAllowForce35({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const codes = ['fc034', 'fc055']

          flatFieldsList.map(_field => {
            if (codes.includes(_field.code)) {
              _field.allowForce = false
            }
          })
        },

        // 39
        notAllowForce39({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const codes = ['fc018', 'fc112', 'fc036', 'fc104', 'fc058', 'fc161', 'fc162', 'fc205', 'fc206', 'fc207', 'fc208', 'fc209']

          flatFieldsList.map(_field => {
            if (codes.includes(_field.code)) {
              _field.allowForce = false
            }
          })
        },

        // 43
        notAllowForce43({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const codes = ['fc106', 'fc165', 'fc196', 'fc197', 'fc198', 'fc199']

          flatFieldsList.map(_field => {
            if (codes.includes(_field.code)) {
              _field.allowForce = false
            }
          })
        },

        // 47
        notAllowForce47({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const codes = ['fc113', 'fc114', 'fc115', 'fc116', 'fc117', 'fc118', 'fc119', 'fc120', 'fc121', 'fc122', 'fc123', 'fc124', 'fc125', 'fc126', 'fc127', 'fc128', 'fc129', 'fc130', 'fc131', 'fc132']

          flatFieldsList.map(_field => {
            if (codes.includes(_field.code)) {
              _field.allowForce = false
            }
          })
        },

        // 67
        notAllowForce67({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const codes = ['fc030', 'fc190']

          flatFieldsList.map(_field => {
            if (codes.includes(_field.code)) {
              _field.allowForce = false
            }
          })
        },

        // 56
        hint56({ fieldsList, focusFieldsIndex, op }) {
          const codesObj = {
            fc081: 'fc080',
            fc144: 'fc143',
            fc146: 'fc145',
            fc148: 'fc147',
            fc150: 'fc149',
            fc152: 'fc151',
            fc154: 'fc153',
            fc156: 'fc155'
          }

          const fields = fieldsList[focusFieldsIndex]

          for (let field of fields) {

            if (codesObj[field.code]) {
              const _field = tools.find(fields, { code: codesObj[field.code] })
              field.hint = `<p  style="color: blue; fontSize: 14px; margin-top: -3px; margin-bottom: 0px">${_field[`${op}Value`]}</p>`
            }
          }
        },

        // 28
        setConstants28({ bill, fieldsList }) {
          if (bill?.saleChannel !== '秒赔') {
            const flatFieldsList = tools.flatArray(fieldsList)
            const rules = ['fc009', 'fc019', 'fc099', 'fc010', 'fc011', 'fc095', 'fc012', 'fc027', 'fc028', 'fc029', 'fc030', 'fc031', 'fc158', 'fc112', 'fc169', 'fc170', 'fc171', 'fc172', 'fc173', 'fc110', 'fc190', 'fc191', 'fc174', 'fc175']

            for (let _field of flatFieldsList) {
              if (rules.includes[_field.code]) {
                _field.disabled = true
              }
            }
          }
        },

        setConstants57({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc080', 'fc143', 'fc145', 'fc147', 'fc149', 'fc151', 'fc153', 'fc155']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0108_太平理赔_全国',
                query: '项目名称'
              }
            }
          })
        },

        // 66
        setConstants66({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc030', 'fc190']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0108_太平理赔_银行代码表',
                query: '银行名称'
              }
            }
          })
        },

        setConstants35({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc034', 'fc055']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0108_太平理赔_医院代码表',
                query: '医院名称'
              }
            }
          })
        },

        setConstants39({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc018', 'fc112', 'fc036', 'fc104', 'fc058', 'fc161', 'fc162', 'fc205', 'fc206', 'fc207', 'fc208', 'fc209', , 'fc239', 'fc240', 'fc241', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246', 'fc247', 'fc248']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0108_太平理赔_诊断代码表',
                query: 'ACCIDENT_NAME'
              }
            }
          })
        },

        setConstants42({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc165', 'fc196', 'fc197', 'fc198', 'fc199']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0108_太平理赔_手术代码表',
                query: '手术名称'
              }
            }
          })
        },

        // CSB0108RC0240000
        disable240({ fieldsList, focusFieldsIndex, codeValues = {} }) {
          const { fc110 } = codeValues

          const fields = fieldsList[focusFieldsIndex]
          fields?.map(_field => {
            if (fc110 == '1')
              if (_field.code == 'fc238') {
                _field.disabled = true
              }
          })
        },

        // CSB0108RC0208001
        setConstants45({ fieldsList, block }) {

          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc113', 'fc114', 'fc115', 'fc116', 'fc117', 'fc118', 'fc119', 'fc120', 'fc121', 'fc122', 'fc123', 'fc124', 'fc125', 'fc126', 'fc127', 'fc128', 'fc129', 'fc130', 'fc131', 'fc132']

          if (block.code === 'bc002') {
            flatFieldsList.map((_field) => {
              if (fields.includes(_field.code)) {
                _field.table = {
                  name: 'B0108_太平理赔_门诊发票大项代码表',
                  query: '大项名称'
                }
              }
            })
          }

          if (block.code === 'bc003') {
            flatFieldsList.map((_field) => {
              if (fields.includes(_field.code)) {
                _field.table = {
                  name: 'B0108_太平理赔_住院发票大项代码表',
                  query: '大项名称'
                }
              }
            })
          }
        },

        disable241({ fieldsList, focusFieldsIndex, bill }) {
          const agency = bill.agency
          const agencies = ['00183000', '00183010', '00183002', '00183012', '00183300', '00183310', '00183301', '00183311']
          const arr = [
            'fc039', 'fc042', 'fc040', 'fc043', 'fc090', 'fc111', 'fc001', 'fc014', 'fc045',
            'fc062', 'fc103', 'fc063', 'fc065', 'fc105', 'fc138', 'fc002', 'fc015', 'fc082',
            'fc219', 'fc220', 'fc221', 'fc222', 'fc223', 'fc113', 'fc114', 'fc115', 'fc116',
            'fc117', 'fc118', 'fc119', 'fc120', 'fc121', 'fc122', 'fc123', 'fc124', 'fc125',
            'fc126', 'fc127', 'fc128', 'fc129', 'fc130', 'fc131', 'fc132', 'fc133', 'fc134',
            'fc135', 'fc136', 'fc137', 'fc005', 'fc006', 'fc007', 'fc008', 'fc037', 'fc038',
            'fc060', 'fc061', 'fc086', 'fc087', 'fc088', 'fc092', 'fc093', 'fc098', 'fc109'
          ]
          const fields = fieldsList[focusFieldsIndex]
          if (agencies.includes(agency)) {
            fields?.map(_field => {
              if (arr.includes(_field.code))
                _field.disabled = true
            })
          }
        },
      }
    },

    // 回车
    enter: {
      methods: {
        // 29
        disable29({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc026') return

          const codes = ['fc024', 'fc025']
          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue === '2') {
            fields.map(_field => {
              if (codes.includes(_field.code)) {
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

        // 30
        disable30({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc020') return

          const codes = ['fc021']
          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue === '0') {
            fields.map(_field => {
              if (codes.includes(_field.code)) {
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

        // 31
        disable31({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc022') return

          const codes = ['fc023']
          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue === '0') {
            fields.map(_field => {
              if (codes.includes(_field.code)) {
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

        // 32
        disable32({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc110') return

          const codes = ['fc009']
          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue === '1') {
            fields.map(_field => {
              if (codes.includes(_field.code)) {
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

        // 33
        disable33({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc157') return

          const codes = ['fc018']
          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue === '2') {
            fields.map(_field => {
              if (codes.includes(_field.code)) {
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

        // 36
        disable36({ field, fieldsList, focusFieldsIndex }) {
          const fields = fieldsList[focusFieldsIndex]

          if (field.code === 'fc034') {
            if (field.resultValue === 'B') {
              fields.map(_field => {
                if (_field.code === 'fc226') {
                  _field.disabled = false
                }
              })
            }
            else {
              fields.map(_field => {
                if (_field.code === 'fc226') {
                  _field.disabled = true
                }
              })
            }
          }
        },

        // 37(禁用字段需要清空值)
        disable37({ field, fieldsList, focusFieldsIndex, op }) {
          const codesListMap = new Map([
            ['0', ['fc036', 'fc104', 'fc058', 'fc161', 'fc162', 'fc205', 'fc206', 'fc207', 'fc208', 'fc209', 'fc228', 'fc229', 'fc230', 'fc231', 'fc232', 'fc233', 'fc234', 'fc235', 'fc236', 'fc237', 'fc239', 'fc240', 'fc241', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246', 'fc247', 'fc248', 'fc249', 'fc250', 'fc251', 'fc252', 'fc253', 'fc254', 'fc255', 'fc256', 'fc257', 'fc258']],
            ['1', ['fc104', 'fc058', 'fc161', 'fc162', 'fc205', 'fc206', 'fc207', 'fc208', 'fc209', 'fc229', 'fc230', 'fc231', 'fc232', 'fc233', 'fc234', 'fc235', 'fc236', 'fc237', 'fc239', 'fc240', 'fc241', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246', 'fc247', 'fc248', 'fc249', 'fc250', 'fc251', 'fc252', 'fc253', 'fc254', 'fc255', 'fc256', 'fc257', 'fc258']],
            ['2', ['fc058', 'fc161', 'fc162', 'fc205', 'fc206', 'fc207', 'fc208', 'fc209', 'fc230', 'fc231', 'fc232', 'fc233', 'fc234', 'fc235', 'fc236', 'fc237', 'fc239', 'fc240', 'fc241', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246', 'fc247', 'fc248', 'fc249', 'fc250', 'fc251', 'fc252', 'fc253', 'fc254', 'fc255', 'fc256', 'fc257', 'fc258']],
            ['3', ['fc161', 'fc162', 'fc205', 'fc206', 'fc207', 'fc208', 'fc209', 'fc231', 'fc232', 'fc233', 'fc234', 'fc235', 'fc236', 'fc237', 'fc239', 'fc240', 'fc241', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246', 'fc247', 'fc248', 'fc249', 'fc250', 'fc251', 'fc252', 'fc253', 'fc254', 'fc255', 'fc256', 'fc257', 'fc258']],
            ['4', ['fc162', 'fc205', 'fc206', 'fc207', 'fc208', 'fc209', 'fc232', 'fc233', 'fc234', 'fc235', 'fc236', 'fc237', 'fc239', 'fc240', 'fc241', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246', 'fc247', 'fc248', 'fc249', 'fc250', 'fc251', 'fc252', 'fc253', 'fc254', 'fc255', 'fc256', 'fc257', 'fc258']],
            ['5', ['fc205', 'fc206', 'fc207', 'fc208', 'fc209', 'fc233', 'fc234', 'fc235', 'fc236', 'fc237', 'fc239', 'fc240', 'fc241', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246', 'fc247', 'fc248', 'fc249', 'fc250', 'fc251', 'fc252', 'fc253', 'fc254', 'fc255', 'fc256', 'fc257', 'fc258']],
            ['6', ['fc206', 'fc207', 'fc208', 'fc209', 'fc234', 'fc235', 'fc236', 'fc237', 'fc239', 'fc240', 'fc241', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246', 'fc247', 'fc248', 'fc249', 'fc250', 'fc251', 'fc252', 'fc253', 'fc254', 'fc255', 'fc256', 'fc257', 'fc258']],
            ['7', ['fc207', 'fc208', 'fc209', 'fc235', 'fc236', 'fc237', 'fc239', 'fc240', 'fc241', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246', 'fc247', 'fc248', 'fc249', 'fc250', 'fc251', 'fc252', 'fc253', 'fc254', 'fc255', 'fc256', 'fc257', 'fc258']],
            ['8', ['fc208', 'fc209', 'fc236', 'fc237', 'fc239', 'fc240', 'fc241', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246', 'fc247', 'fc248', 'fc249', 'fc250', 'fc251', 'fc252', 'fc253', 'fc254', 'fc255', 'fc256', 'fc257', 'fc258']],
            ['9', ['fc209', 'fc237', 'fc239', 'fc240', 'fc241', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246', 'fc247', 'fc248', 'fc249', 'fc250', 'fc251', 'fc252', 'fc253', 'fc254', 'fc255', 'fc256', 'fc257', 'fc258']],
            ['10', ['fc239', 'fc240', 'fc241', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246', 'fc247', 'fc248', 'fc249', 'fc250', 'fc251', 'fc252', 'fc253', 'fc254', 'fc255', 'fc256', 'fc257', 'fc258']],
            ['11', ['fc240', 'fc241', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246', 'fc247', 'fc248', 'fc250', 'fc251', 'fc252', 'fc253', 'fc254', 'fc255', 'fc256', 'fc257', 'fc258']],
            ['12', ['fc241', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246', 'fc247', 'fc248', 'fc251', 'fc252', 'fc253', 'fc254', 'fc255', 'fc256', 'fc257', 'fc258']],
            ['13', ['fc242', 'fc243', 'fc244', 'fc245', 'fc246', 'fc247', 'fc248', 'fc252', 'fc253', 'fc254', 'fc255', 'fc256', 'fc257', 'fc258']],
            ['14', ['fc243', 'fc244', 'fc245', 'fc246', 'fc247', 'fc248', 'fc253', 'fc254', 'fc255', 'fc256', 'fc257', 'fc258']],
            ['15', ['fc244', 'fc245', 'fc246', 'fc247', 'fc248', 'fc254', 'fc255', 'fc256', 'fc257', 'fc258']],
            ['16', ['fc245', 'fc246', 'fc247', 'fc248', 'fc255', 'fc256', 'fc257', 'fc258']],
            ['17', ['fc246', 'fc247', 'fc248', 'fc256', 'fc257', 'fc258']],
            ['18', ['fc247', 'fc248', 'fc257', 'fc258']],
            ['19', ['fc248', 'fc258']],
          ])

          const fields = fieldsList[focusFieldsIndex]

          if (field.code === 'fc160') {
            const codes = codesListMap.get(field.resultValue)
            const keys = Array.from(codesListMap.keys())

            fields.map(_field => {
              for (let key of keys) {
                const values = codesListMap.get(key)

                if (values.includes(_field.code)) {
                  _field.disabled = false
                }
              }

              // 屏蔽
              if (codes?.includes(_field.code)) {
                _field.disabled = true
                _field[`${op}Value`] = ''
                _field.resultValue = ''
              }
            })
          }
        },

        // 40
        disable40({ field, fieldsList, focusFieldsIndex }) {
          const codesListMap = new Map([
            ['fc036', 'fc228'],
            ['fc104', 'fc229'],
            ['fc058', 'fc230'],
            ['fc161', 'fc231'],
            ['fc162', 'fc232'],
            ['fc205', 'fc233'],
            ['fc206', 'fc234'],
            ['fc207', 'fc235'],
            ['fc208', 'fc236'],
            ['fc209', 'fc237'],
            ['fc239', 'fc249'],
            ['fc240', 'fc250'],
            ['fc241', 'fc251'],
            ['fc242', 'fc252'],
            ['fc243', 'fc253'],
            ['fc244', 'fc254'],
            ['fc245', 'fc255'],
            ['fc246', 'fc256'],
            ['fc247', 'fc257'],
            ['fc248', 'fc258'],
          ])

          const fields = fieldsList[focusFieldsIndex]
          const code = codesListMap.get(field.code)

          if (code) {
            const resultValue = field.resultValue

            fields.map(_field => {
              if (_field.code === code) {
                _field.disabled = false

                if (resultValue !== 'B') {
                  _field.disabled = true
                }
              }
            })
          }
        },

        // 45
        disable45({ field, fieldsList, focusFieldsIndex }) {
          const { code, resultValue } = field

          if (code === 'fc042') {
            const fields = fieldsList[focusFieldsIndex]
            const codes = ['fc040', 'fc041', 'fc090', 'fc111']

            fields.map(_field => {
              if (codes.includes(_field.code)) {
                _field.disabled = false

                if (resultValue === '2') {
                  _field.disabled = true
                }
              }
            })
          }
        },

        // 58
        disable58({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc055') return

          const codes = ['fc227']
          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue !== 'B') {
            fields.map(_field => {
              if (codes.includes(_field.code)) {
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

        // 59(禁用字段需要清空值)
        disable59({ field, fieldsList, focusFieldsIndex, op }) {
          const { code, resultValue } = field
          const codesListObj = {
            0: ['fc177', 'fc178', 'fc179', 'fc180', 'fc181'],
            1: ['fc178', 'fc179', 'fc180', 'fc181'],
            2: ['fc179', 'fc180', 'fc181'],
            3: ['fc180', 'fc181'],
            4: ['fc181']
          }

          const codeValues = []

          for (let key in codesListObj) {
            codeValues.push(...codesListObj[key])
          }

          if (code === 'fc176') {
            const fields = fieldsList[focusFieldsIndex]

            fields.map(_field => {
              if (codeValues.includes(_field.code)) {
                _field.disabled = false
              }

              const targetCodeValues = codesListObj[resultValue]

              if (targetCodeValues?.includes(_field.code)) {
                _field.disabled = true
                _field[`${op}Value`] = ''
                _field.resultValue = ''
              }
            })
          }
        },

        // 60
        disable60({ field, fieldsList, focusFieldsIndex }) {
          const { code, resultValue } = field
          const codesListObj = {
            0: ['fc183', 'fc184', 'fc185', 'fc186', 'fc187'],
            1: ['fc184', 'fc185', 'fc186', 'fc187'],
            2: ['fc185', 'fc186', 'fc187'],
            3: ['fc186', 'fc187'],
            4: ['fc187']
          }

          const codeValues = []

          for (let key in codesListObj) {
            codeValues.push(...codesListObj[key])
          }

          if (code === 'fc182') {
            const fields = fieldsList[focusFieldsIndex]

            fields.map(_field => {
              if (codeValues.includes(_field.code)) {
                _field.disabled = false
              }

              const targetCodeValues = codesListObj[resultValue]

              if (targetCodeValues?.includes(_field.code)) {
                _field.disabled = true
              }
            })
          }
        },

        // 62
        disable62({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc103') return

          const codes = ['fc063', 'fc064', 'fc105', 'fc138']
          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue === '2') {
            fields.map(_field => {
              if (codes.includes(_field.code)) {
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

        // 64
        disable64({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc027') return

          const codes = ['fc028', 'fc030']
          const validValues = ['1', 'A']
          const fields = fieldsList[focusFieldsIndex]

          if (validValues.includes(field.resultValue)) {
            fields.map(_field => {
              if (codes.includes(_field.code)) {
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

        // 65(Waiting backend)
        // disable65({ field, fieldsList, focusFieldsIndex }) {

        // },

        // 68
        disable68({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc173') return

          const codes = ['fc169', 'fc170', 'fc171', 'fc172']
          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue === '2') {
            fields.map(_field => {
              if (codes.includes(_field.code)) {
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

        // 56
        hint56({ field, fieldsList, focusFieldsIndex }) {
          const codesMap = new Map([
            ['fc080', 'fc081'],
            ['fc143', 'fc144'],
            ['fc145', 'fc146'],
            ['fc147', 'fc148'],
            ['fc149', 'fc150'],
            ['fc151', 'fc152'],
            ['fc153', 'fc154'],
            ['fc155', 'fc156']
          ])

          const code = codesMap.get(field.code)

          if (code) {
            const fields = fieldsList[focusFieldsIndex]

            for (let _field of fields) {
              if (_field.code === code) {
                _field.hint = `<p  style="color: blue; fontSize: 14px; margin-top: -3px; margin-bottom: 0px">${field.resultValue}</p>`
              }
            }
          }
        },

        disable69({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc177' && field.code !== 'fc178' && field.code !== 'fc179' && field.code !== 'fc180') return

          const fields = fieldsList[focusFieldsIndex]
          if (field.code == 'fc177') {
            if (field.resultValue === 'A') {
              const codes = ['fc178', 'fc179', 'fc180', 'fc181']
              fields.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = true
                }
              })
            } else {
              const codes = ['fc178', 'fc179', 'fc180', 'fc181']
              fields.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = false
                }
              })
            }
          }

          if (field.code == 'fc178') {
            if (field.resultValue === 'A') {
              const codes = ['fc179', 'fc180', 'fc181']
              fields.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = true
                }
              })
            } else {
              const codes = ['fc179', 'fc180', 'fc181']
              fields.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = false
                }
              })
            }
          }

          if (field.code == 'fc179') {
            if (field.resultValue === 'A') {
              const codes = ['fc180', 'fc181']
              fields.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = true
                }
              })
            } else {
              const codes = ['fc180', 'fc181']
              fields.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = false
                }
              })
            }
          }

          if (field.code == 'fc180') {
            if (field.resultValue === 'A') {
              const codes = ['fc181']
              fields.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = true
                }
              })
            } else {
              const codes = ['fc181']
              fields.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = false
                }
              })
            }
          }
        },

        disable70({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc183' && field.code !== 'fc184' && field.code !== 'fc185' && field.code !== 'fc186') return

          const fields = fieldsList[focusFieldsIndex]
          if (field.code == 'fc183') {
            if (field.resultValue === 'A') {
              const codes = ['fc184', 'fc185', 'fc186', 'fc187']
              fields.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = true
                }
              })
            } else {
              const codes = ['fc184', 'fc185', 'fc186', 'fc187']
              fields.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = false
                }
              })
            }
          }

          if (field.code == 'fc184') {
            if (field.resultValue === 'A') {
              const codes = ['fc185', 'fc186', 'fc187']
              fields.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = true
                }
              })
            } else {
              const codes = ['fc185', 'fc186', 'fc187']
              fields.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = false
                }
              })
            }
          }

          if (field.code == 'fc185') {
            if (field.resultValue === 'A') {
              const codes = ['fc186', 'fc187']
              fields.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = true
                }
              })
            } else {
              const codes = ['fc186', 'fc187']
              fields.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = false
                }
              })
            }
          }

          if (field.code == 'fc186') {
            if (field.resultValue === 'A') {
              const codes = ['fc187']
              fields.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = true
                }
              })
            } else {
              const codes = ['fc187']
              fields.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = false
                }
              })
            }
          }
        },
      }
    },

    // 临时保存
    sessionSave: {
      methods: {
        // 50
        disabled50({ field, fieldsList, focusFieldsIndex }) {
          const fields = fieldsList[focusFieldsIndex]
          const codesList = [
            ['fc113', 'fc133'],
            ['fc114', 'fc134'],
            ['fc115', 'fc135'],
            ['fc116', 'fc136'],
            ['fc117', 'fc137'],
            ['fc118', 'fc005'],
            ['fc119', 'fc006'],
            ['fc120', 'fc007'],
            ['fc121', 'fc008'],
            ['fc122', 'fc037'],
            ['fc123', 'fc038'],
            ['fc124', 'fc060'],
            ['fc125', 'fc061'],
            ['fc126', 'fc086'],
            ['fc127', 'fc087'],
            ['fc128', 'fc088'],
            ['fc129', 'fc092'],
            ['fc130', 'fc093'],
            ['fc131', 'fc098'],
            ['fc132', 'fc109']
          ]

          let restCodes = []

          for (let i = 0; i < codesList.length; i++) {
            const codes = codesList[i]

            if (codes.includes(field.code)) {
              restCodes = codesList.slice(i + 1)
              break
            }
          }

          const flatRestCodes = tools.flatArray(restCodes)

          fields.map(_field => {
            if (flatRestCodes.includes(_field.code)) {
              _field.disabled = true
            }
          })
        },

        // 55
        disabled55({ field, fieldsList, focusFieldsIndex, flatFieldsList }) {
          const codes = ['fc080', 'fc143', 'fc145', 'fc147', 'fc149', 'fc151', 'fc153', 'fc155']
          if (!codes.includes(field.code)) return
          const fields = fieldsList[focusFieldsIndex]
          const codesList = [
            ['fc080', 'fc081'],
            ['fc143', 'fc144'],
            ['fc145', 'fc146'],
            ['fc147', 'fc148'],
            ['fc149', 'fc150'],
            ['fc151', 'fc152'],
            ['fc153', 'fc154'],
            ['fc155', 'fc156']
          ]

          let restCodes = []

          for (let i = 0; i < codesList.length; i++) {
            const codes = codesList[i]

            if (field.resultValue != '' && codes.includes(field.code)) {
              restCodes = codesList.slice(i + 1)
              break
            }
            if (field.resultValue == '' && codes.includes(field.code)) {
              restCodes = codesList.slice(i)
              break
            }
          }

          const flatRestCodes = tools.flatArray(restCodes)

          fields.map(_field => {
            if (flatRestCodes.includes(_field.code)) {
              _field.disabled = true
            }
          })
          for (let field of fields) {
            if (!field.disabled && field.resultValue == '') {
              flatFieldsList.map(_field => {
                _field.autofocus = false;
                _field.uniqueKey = `enter_${field.uniqueId}_${Date.now()}`;
              });
              field.autofocus = false
              break
            }
          }
        }
      }
    },

    // 提交前
    beforeSubmit: {
      methods: {
        // 44
        validate44: function ({ block, fieldsList, op }) {
          if (block.code !== 'bc002') return true
          if (op !== 'op1') return true

          const fields = fieldsList[0]
          const fc039Field = tools.find(fields, { code: 'fc039' })
          const codes = ['fc133', 'fc134', 'fc135', 'fc136', 'fc137', 'fc005', 'fc006', 'fc007', 'fc008', 'fc037', 'fc038', 'fc060', 'fc061', 'fc086', 'fc087', 'fc088', 'fc092', 'fc093', 'fc098', 'fc109']

          if (!fc039Field || fc039Field.resultValue === '?') {
            return true
          }

          for (let field of fields) {
            if (codes.includes(field.code)) {
              if (field.resultValue === '?') {
                return true
              }
            }
          }

          // let count = 0
          let count = new BigNumber(0)

          for (let field of fields) {
            if (codes.includes(field.code)) {
              const resultValue = +field.resultValue || 0
              // count += resultValue
              count = count.plus(resultValue)
            }
          }

          // const fc039Value = +fc039Field.resultValue

          // const diff = fc039Value - count

          const fc039Value = new BigNumber(+fc039Field.resultValue)

          const diff = fc039Value.minus(count).toString()

          if (diff != 0) {
            return {
              popup: 'confirm',
              errorMessage: `门诊发票明细金额与总金额不一致，差额为${diff}，请确认并修改!`
            }
          }

          return true
        },

        // 49
        validate49: function ({ fieldsList }) {
          const fields = fieldsList[0]
          const codesListMap = new Map([
            ['fc113', 'fc133'],
            ['fc114', 'fc134'],
            ['fc115', 'fc135'],
            ['fc116', 'fc136'],
            ['fc117', 'fc137'],
            ['fc118', 'fc005'],
            ['fc119', 'fc006'],
            ['fc120', 'fc007'],
            ['fc121', 'fc008'],
            ['fc122', 'fc037'],
            ['fc123', 'fc038'],
            ['fc124', 'fc060'],
            ['fc125', 'fc061'],
            ['fc126', 'fc086'],
            ['fc127', 'fc087'],
            ['fc128', 'fc088'],
            ['fc129', 'fc092'],
            ['fc130', 'fc093'],
            ['fc131', 'fc098'],
            ['fc132', 'fc109']
          ])

          for (let field of fields) {
            const { code, resultValue } = field
            const targetCode = codesListMap.get(code)

            if (targetCode && resultValue) {
              const targetField = tools.find(fields, { code: targetCode })

              if (!targetField?.resultValue) {
                return {
                  errorMessage: '发票大项内容录入遗漏，请检查!'
                }
              }
            }
          }

          return true
        },

        // 52
        validate52({ block, fieldsList }) {
          if (block.code !== 'bc004') return true

          const flatFieldsList = tools.flatArray(fieldsList)
          const values = []

          flatFieldsList?.map(_field => {
            _field.resultValue && values.push(_field.resultValue)
          })

          if (!values?.length) {
            return {
              errorMessage: '清单不能空白提交，请检查!'
            }
          }

          return true
        },

        // 53
        validate1153: function ({ fieldsList }) {
          const fields = fieldsList[0]
          const codesListMap = new Map([
            ['fc080', 'fc081'],
            ['fc143', 'fc144'],
            ['fc145', 'fc146'],
            ['fc147', 'fc148'],
            ['fc149', 'fc150'],
            ['fc151', 'fc152'],
            ['fc153', 'fc154'],
            ['fc155', 'fc156']
          ])

          for (let field of fields) {
            const { code, resultValue } = field
            const targetCode = codesListMap.get(code)

            if (targetCode && resultValue) {
              const targetField = tools.find(fields, { code: targetCode })

              if (!targetField?.resultValue) {
                return {
                  errorMessage: '清单内容录入遗漏，请检查!'
                }
              }
            }
          }

          return true
        },

        // 54
        validate54: function ({ block, fieldsList }) {
          if (!block.isLoop) return true

          const fields = fieldsList[0]
          const restFields = fieldsList.slice(1)
          const flatRestFields = tools.flatArray(restFields)
          const codes = [
            'fc080', 'fc081',
            'fc143', 'fc144',
            'fc145', 'fc146',
            'fc147', 'fc148',
            'fc149', 'fc150',
            'fc151', 'fc152',
            'fc153', 'fc154',
            'fc155', 'fc156'
          ]

          // 若 >=2 的循环分块有值，不需要校验
          for (let field of flatRestFields) {
            const { code, resultValue } = field

            if (codes.includes(code)) {
              if (resultValue) {
                return true
              }
            }
          }

          // 若第一个循环分块有字段没有值，不需要校验
          for (let field of fields) {
            if (codes.includes(field.code)) {
              if (!field.resultValue) {
                return true
              }
            }
          }

          return {
            popup: 'confirm',
            // errorMessage: '请确认图片所有项目名称是否录入完整，未录入完整请点击“不完整”按回车后继续录入!'
            errorMessage: '请确认图片所有项目名称是否录入完整，未录入完整请点击“取消”按回车后继续录入!'
          }
        },

        // 61
        validate61({ block, fieldsList, op }) {
          if (block.code !== 'bc003') return true
          if (op !== 'op1') return true

          const fields = fieldsList[0]
          const fc062Field = tools.find(fields, { code: 'fc062' })
          const codes = ['fc133', 'fc134', 'fc135', 'fc136', 'fc137', 'fc005', 'fc006', 'fc007', 'fc008', 'fc037', 'fc038', 'fc060', 'fc061', 'fc086', 'fc087', 'fc088', 'fc092', 'fc093', 'fc098', 'fc109']

          if (!fc062Field || fc062Field.resultValue === '?') {
            return true
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
              const resultValue = +field.resultValue || 0
              count = count.plus(resultValue)
            }
          }

          const fc062Value = new BigNumber(+fc062Field.resultValue)

          const diff = fc062Value.minus(count).toString()

          if (diff != 0) {
            return {
              popup: 'confirm',
              errorMessage: `住院发票明细金额与总金额不一致，差额为${diff}，请确认并修改!`
            }
          }

          return true
        },

        // CSB0108RC0242000
        validate53: function ({ fieldsList, op }) {
          const fields = fieldsList[0]

          for (let field of fields) {
            let result = /^[?]+$/.test(field.resultValue)
            if (field.disabled && result) {
              field[`${op}Value`] = ''
              field.resultValue = ''
            }
          }

          return true
        },
      }
    }
  }
}

export default B0108
