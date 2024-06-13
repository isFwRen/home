import { tools, sessionStorage } from 'vue-rocket'
import moment from 'moment'
import { ignoreFreeValue } from '../tools'
import { MessageBox, Notification } from 'element-ui';

const B0116 = {
  op0: {
    // 记录最后一次存储的合法field
    memoFields: ['fc021', 'fc196', 'fc212', 'fc022', 'fc195', 'fc211', 'fc023', 'fc194', 'fc210'],

    // 记录相同 code 的 field 的值
    memoFieldValues: ['fc088', 'fc089', 'fc090', 'fc112'],

    // fields 的值从 targets 里的值选择
    dropdownFields: [
      {
        targets: ['fc088'],
        fields: ['fc089', 'fc090']
      }
    ],


    // 校验规则
    rules: [
      // CSB0116RC0001000
      {
        fields: ['fc088'],
        validate01: function ({ field, fieldsObject, thumbIndex, value }) {
          const fc088Values = []

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage

            if (sessionStorage || thumbIndex === +key) {
              const _fieldsList = fieldsObject[key].fieldsList

              for (let _fields of _fieldsList) {
                for (let _field of _fields) {
                  if (_field.code === 'fc088' && _field.uniqueId !== field.uniqueId) {
                    fc088Values.push(_field.resultValue)
                  }
                }
              }
            }
          }

          if (fc088Values.includes(value)) {
            return '发票属性不能重复'
          }

          return true
        }
      },

      // CSB0116RC0002000
      {
        fields: ['fc089', 'fc090'],
        validate02: function ({ includes, value }) {
          if (includes) {
            const result = includes.find(text => text === value)

            if (!result) {
              return '没有此发票，请核实'
            }
          }

          return true
        }
      },

      // CSB0116RC0026000
      {
        fields: ['fc111'],
        validate26: function ({ value, fieldsIndex, fieldsList }) {
          let arr = fieldsList[fieldsIndex]
          let fc112 = ''
          for (let field of arr) {
            if (field.code == 'fc112') fc112 = field.resultValue
          }
          if (fc112 == '1') {
            if (value == '1' || value == '2') return '北京账单只能选录对应的账单信息， 若无法判断， 则根据影像默认选录3或4'
          }

          if (fc112 == '2') {
            if (value !== '1' && value !== '2') return '非北京账单只能选录1或2'
          }
          return true
        }
      },

      {
        fields: ['fc035'],
        validate25: function ({ value, items }) {
          if (value.includes('?') || value.includes('？') || !value || value == 'A' || value == 'F') {
            return true
          }

          const result = items.find((text) => text === value)

          if (!result) {
            return '医院录入错误，没有下拉提示时需要按单录入强过'
          }

          return true
        }
      },
    ],

    // 提示文本
    hints: [],

    // 工序完成初始化
    init: {
      methods: {

      }
    },

    // 字段已生成
    updateFields: {
      methods: {
        // CSB0116RC0020000
        setConstants20({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc035']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0116_华夏人寿团险理赔_医院代码',
                query: '中文名'
              }
            }
          })
        },
      }
    },

    // 回车
    enter: {
      methods: {

      }
    },

    // F8(提交前校验)
    beforeSubmit: {
      methods: {
        // CSB0116RC0003000
        validate03({ fieldsObject }) {
          const [fc090Values, fc088Values] = [[], []]

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  if (field.code === 'fc090') {
                    fc090Values.push(field.resultValue)
                  }

                  if (field.code === 'fc088') {
                    fc088Values.push(field.resultValue)
                  }
                }
              }
            }
          }

          for (let value of fc090Values) {
            if (![...fc088Values].includes(value)) {
              return {
                errorMessage: `报销单${value}没有匹配的发票!`
              }
            }
          }

          return true
        },

        // CSB0116RC0004000
        validate04({ fieldsObject }) {
          const [fc089Values, fc088Values] = [[], []]

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  if (field.code === 'fc089') {
                    fc089Values.push(field.resultValue)
                  }

                  if (field.code === 'fc088') {
                    fc088Values.push(field.resultValue)
                  }
                }
              }
            }
          }

          for (let value of fc089Values) {
            if (![...fc088Values].includes(value)) {
              return {
                errorMessage: `清单${value}没有匹配的发票，不可提交!`
              }
            }
          }

          return true
        },

        // CSB0116RC0005000
        validate05({ fieldsObject }) {
          const fc086Values = []

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  if (field.code === 'fc086') {
                    fc086Values.push(field.resultValue)
                  }
                }
              }
            }
          }

          if (!fc086Values.includes('1') && !fc086Values.includes('2')) {
            return {
              popup: 'confirm',
              errorMessage: `案件必须要有申请表，确认提交吗`
            }
          }

          return true
        },

        // CSB0116RC0006000
        // validate06({ fieldsObject }) {
        //   const [fc088Values, fc089Values, fc090Values] = [[], [], []]

        //   for (let key in fieldsObject) {
        //     const sessionStorage = fieldsObject[key].sessionStorage
        //     const fieldsList = fieldsObject[key].fieldsList

        //     if (sessionStorage) {
        //       for (let fields of fieldsList) {
        //         for (let field of fields) {
        //           const { code, resultValue } = field

        //           if (code === 'fc088') {
        //             resultValue && fc088Values.push(resultValue)
        //           }

        //           if (code === 'fc089') {
        //             resultValue && fc089Values.push(resultValue)
        //           }

        //           if (code === 'fc090') {
        //             resultValue && fc090Values.push(resultValue)
        //           }
        //         }
        //       }
        //     }
        //   }

        //   const mergeValues = [...fc089Values, ...fc090Values]

        //   for (let value of mergeValues) {
        //     if (!fc088Values.includes(value)) {
        //       return {
        //         errorMessage: `发票${value}没有匹配的清单/报销单，不可提交!`
        //       }
        //     }
        //   }

        //   return true
        // },

        // CSB0116RC0007000
        validate07({ fieldsObject }) {
          const fc086Values = []

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  if (field.code === 'fc086') {
                    fc086Values.push(field.resultValue)
                  }
                }
              }
            }
          }

          if (!fc086Values.includes('6')) {
            return {
              errorMessage: `缺少诊断书，不可提交`
            }
          }

          return true
        },

        // CSB0116RC0008000
        validate08({ fieldsObject }) {
          const fc086Values = []

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  if (field.code === 'fc086') {
                    fc086Values.push(field.resultValue)
                  }
                }
              }
            }
          }
          let count = 0
          for (let value of fc086Values) {
            if (value == '11') count++
          }
          if (count >= 2) {
            return {
              errorMessage: `同一案件仅能有一个重疾/轻症，不可提交`
            }
          }
          return true
        },

        // CSB0116RC0010000
        validate10({ fieldsObject }) {
          const fc086Values = []

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  if (field.code === 'fc086') {
                    fc086Values.push(field.resultValue)
                  }
                }
              }
            }
          }

          if (!fc086Values.includes('8')) {
            return {
              errorMessage: `缺少受益人身份证，不可提交`
            }
          }

          return true
        },
      },
    }
  },

  op1op2opq: {
    // 校验规则
    rules: [
      // CSB0116RC0001200 
      {
        fields: ['fc005', 'fc008', 'fc017', 'fc018', 'fc019', 'fc192', 'fc054', 'fc055', 'fc032', 'fc036', 'fc216', 'fc217', 'fc220', 'fc221'],
        validateDate12: function ({ value }) {
          if (!value) return true

          if (/[A,1 \?]/.test(value)) {
            return true
          }

          if (value.length !== 8 || moment(value).format('YYYYMMDD') === 'Invalid date') {
            return '日期格式错误'
          }

          return true
        }
      },

      // CSB0116RC0014000
      // {
      //   fields: ['fc009'],
      //   validateDate14: function ({ value, fieldsList }) {
      //     if (!value) return true
      //     const fields = fieldsList[0]
      //     const fc008Field = fields.find(field => field.code == 'fc008')

      //     if (fc008Field.resultValue == '' && value == '') {
      //       return '证件有效期和证件是否长期有效至少录一个'
      //     }

      //     if ((fc008Field.resultValue != '' || fc008Field.resultValue == 'A') && (value != '' || value != 'A')) {
      //       return '证件有效期和证件是否长期有效不能同时录入'
      //     }

      //     return true
      //   }
      // },

      // CSB0116RC0018000
      {
        fields: ['fc015'],
        validate18: function ({ value, items }) {
          if (value.includes('?') || value.includes('？') || !value || value == 'F' || value == 'A') {
            return true
          }

          const result = items.find((text) => text === value)

          if (!result) {
            return '意外出险细节错误， 请根据下拉提示内容选录'
          }

          return true
        }
      },

      // CSB0116RC0019000
      {
        fields: ['fc016'],
        validate19: function ({ value, items }) {
          if (value.includes('?') || value.includes('？') || !value || value == 'F' || value == 'A') {
            return true
          }

          const result = items.find((text) => text === value)

          if (!result) {
            return '损伤外部原因错误， 请根据下拉提示内容选录'
          }

          return true
        }
      },

      // CSB0116RC0021000
      {
        fields: ['fc023', 'fc194', 'fc210'],
        validate21: function ({ value, items }) {
          if (value == 'A' || value == 'F' || value == '?' || value == '？') return true
          const result = items.find((text) => text === value)

          if (!result) {
            return '录入有误， 请根据下拉提示内容选录'
          }

          return true
        }
      },

      // CSB0116RC0021000
      {
        fields: ['fc022', 'fc195', 'fc211'],
        validate211: function ({ value, items }) {
          if (value == 'A' || value == 'F' || value == '?' || value == '？') return true
          const result = items.find((text) => text === value)
          if (items[0].includes(value)) return true
          if (!result) {
            return '录入有误， 请根据下拉提示内容选录'
          }

          return true
        }
      },

      // CSB0116RC0021000
      {
        fields: ['fc021', 'fc196', 'fc212'],
        validate212: function ({ value, items }) {
          if (value == 'A' || value == 'F' || value == '?' || value == '？') return true
          const result = items.find((text) => text === value)
          if (items[0].includes(value)) return true
          if (!result) {
            return '录入有误， 请根据下拉提示内容选录'
          }

          return true
        }
      },

      // CSB0116RC0032000
      {
        fields: ['fc027'],
        validate32: function ({ value, items }) {
          if (value.includes('?') || value.includes('？') || !value || value == 'F' || value == 'A') {
            return true
          }

          const result = items.find((text) => text === value)

          if (!result) {
            return '残疾名称错误， 请根据下拉提示内容选录'
          }

          return true
        }
      },

      // CSB0116RC0033000
      {
        fields: ['fc108'],
        validate32: function ({ value, items }) {
          if (value.includes('?') || value.includes('？') || !value || value == 'F' || value == 'A') {
            return true
          }

          const result = items.find((text) => text === value)

          if (!result) {
            return '重疾名称错误， 请根据下拉提示内容选录'
          }

          return true
        }
      },

      // CSB0116RC0024000
      {
        fields: ['fc044', 'fc096', 'fc099'],
        validate24: function ({ value, items }) {
          if (value.includes('?') || value.includes('？') || !value || value == 'A' || value == 'F') {
            return true
          }

          const result = items.find((text) => text === value)

          if (!result) {
            return '伤病诊断错误， 请根据下拉提示内容选录'
          }

          return true
        }
      },

      // CSB0116RC0025000
      {
        fields: ['fc041', 'fc102', 'fc105'],
        validate25: function ({ value, items }) {
          if (value.includes('?') || value.includes('？') || !value || value == 'A' || value == 'F') {
            return true
          }

          const result = items.find((text) => text === value)

          if (!result) {
            return '手术术式错误， 请根据下拉提示内容选录'
          }

          return true
        }
      },
    ],

    // 提示文本
    hints: [

    ],

    // 字段已生成
    updateFields: {
      methods: {
        // CSB0116RC0015000
        disable15({ fieldsList }) {
          const codes = ['fc010', 'fc011', 'fc043', 'fc097', 'fc100', 'fc005', 'fc009', 'fc188', 'fc193', 'fc147', 'fc148', 'fc149', 'fc150', 'fc151', 'fc152', 'fc153', 'fc154', 'fc051', 'fc155', 'fc156', 'fc157', 'fc158', 'fc159', 'fc160', 'fc161', 'fc162', 'fc163', 'fc164', 'fc165', 'fc166', 'fc167', 'fc168', 'fc169', 'fc170', 'fc171', 'fc172', 'fc173', 'fc174', 'fc026', 'fc028', 'fc030', 'fc033']

          fieldsList?.map(fields => {
            fields?.map(_field => {
              if (codes.includes(_field.code)) {
                _field.disabled = true
              }
            })
          })
        },

        // CSB0116RC0016000
        disable16({ fieldsList, focusFieldsIndex, codeValues = {} }) {
          const { fc013 } = codeValues

          const fields = fieldsList[focusFieldsIndex]
          const fc019Field = fields.find(field => field.code == 'fc019')

          if (fc013 != '9' && fc019Field) {
            fc019Field.disabled = true
          }
        },

        // CSB0116RC0018000
        setConstants18({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc015']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0116_华夏人寿团险理赔_意外出险细节',
                query: '中文名'
              }
            }
          })
        },

        // CSB0116RC0019000
        setConstants19({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc016']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0116_华夏人寿团险理赔_损伤外部原因',
                query: '中文名'
              }
            }
          })
        },

        // CSB0116RC0021000
        setConstants21: function ({ flatFieldsList }) {
          const fields = ['fc023', 'fc194', 'fc210']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0116_华夏人寿团险理赔_地址库',
                query: '区、县中文名',
                targets: ['省中文名', '市中文名', '区、县中文名']
              }
            }
          })
        },

        // CSB0116RC0021000
        setConstants211: function ({ flatFieldsList }) {
          const fields = ['fc022', 'fc195', 'fc211']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0116_华夏人寿团险理赔_地址库',
                query: '市中文名',
                targets: ['省中文名', '市中文名']
              }
            }
          })
        },

        // CSB0116RC0021000
        setConstants212: function ({ flatFieldsList }) {
          const fields = ['fc021', 'fc196', 'fc212']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0116_华夏人寿团险理赔_地址库',
                query: '省中文名'
              }
            }
          })
        },

        // CSB0116RC0024000
        setConstants24: function ({ flatFieldsList }) {
          const fields = ['fc044', 'fc096', 'fc099']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0116_华夏人寿团险理赔_疾病编码',
                query: '中文名'
              }
            }
          })
        },

        // CSB0116RC0025000
        setConstants25: function ({ flatFieldsList }) {
          const fields = ['fc041', 'fc102', 'fc105']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0116_华夏人寿团险理赔_手术术士编码',
                query: '中文名'
              }
            }
          })
        },

        // CSB0116RC0032000
        setConstants32: function ({ flatFieldsList }) {
          const fields = ['fc027']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0116_华夏人寿团险理赔_残疾',
                query: '残疾中文名'
              }
            }
          })
        },

        // CSB0116RC0033000
        setConstants33: function ({ flatFieldsList }) {
          const fields = ['fc108']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0116_华夏人寿团险理赔_重疾',
                query: '中文名'
              }
            }
          })
        },

        // CSB0116RC0035000
        setConstants35: function ({ flatFieldsList }) {
          const fields = ['fc203', 'fc206']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0116_华夏人寿团险理赔_职业编码',
                query: '职业名称'
              }
            }
          })
        },

        // CSB0116RC0030000
        disable30({ fieldsList, focusFieldsIndex, codeValues = {} }) {
          const { fc112 } = codeValues
          console.log(fc112);
          const fields = fieldsList[focusFieldsIndex]
          const fc202Field = fields.find(field => field.code == 'fc202')

          if (fc202Field && fc112 == '1') {
            fc202Field.disabled = true
          }
        },

        // CSB0116RC0027000
        disable27({ fieldsList }) {
          const codes = ['fc054', 'fc055', 'fc056', 'fc058', 'fc059', 'fc060', 'fc061', 'fc062', 'fc077', 'fc082', 'fc080', 'fc083', 'fc063', 'fc064', 'fc065', 'fc066', 'fc067', 'fc068', 'fc069', 'fc070', 'fc071', 'fc072', 'fc073', 'fc074', 'fc075', 'fc078', 'fc079', 'fc076', 'fc081', 'fc084']

          fieldsList?.map(fields => {
            fields?.map(_field => {
              if (codes.includes(_field.code)) {
                _field.disabled = true
              }
            })
          })
        },

        // CSB0116RC0027000
        disable271({ fieldsList, focusFieldsIndex, codeValues = {} }) {
          const { fc111 } = codeValues
          console.log(fc111);
          const codeMaps = new Map([
            ['1', ['fc054']],
            ['2', ['fc055', 'fc056']],
            ['3', ['fc054', 'fc058', 'fc059', 'fc060', 'fc061', 'fc062', 'fc063', 'fc064', 'fc065', 'fc066', 'fc067', 'fc068', 'fc069', 'fc070', 'fc071', 'fc072', 'fc073', 'fc074', 'fc075', 'fc076']],
            ['4', ['fc055', 'fc056', 'fc058', 'fc059', 'fc060', 'fc061', 'fc062', 'fc077', 'fc066', 'fc067', 'fc068', 'fc069', 'fc071', 'fc073', 'fc075', 'fc078', 'fc079']],
            ['5', ['fc055', 'fc056']],
            ['6', ['fc054', 'fc084']],
            ['7', ['fc055', 'fc056', 'fc084']],
            ['8', ['fc054', 'fc058', 'fc059', 'fc060', 'fc061', 'fc062', 'fc077', 'fc080', 'fc066', 'fc067', 'fc068', 'fc072', 'fc073', 'fc074', 'fc079', 'fc076']],
            ['9', ['fc055', 'fc056', 'fc058', 'fc059', 'fc060', 'fc061', 'fc062', 'fc077', 'fc080', 'fc066', 'fc067', 'fc068', 'fc069', 'fc070', 'fc072', 'fc073', 'fc074', 'fc079', 'fc076', 'fc081']],
            ['10', ['fc054', 'fc060', 'fc061', 'fc062', 'fc077', 'fc082', 'fc080', 'fc083']],
            ['11', ['fc055', 'fc056', 'fc060', 'fc061', 'fc062', 'fc077', 'fc082', 'fc080', 'fc083']],
          ])

          const fields = fieldsList[focusFieldsIndex]

          if (codeMaps.get(fc111)) {
            let arr = codeMaps.get(fc111)
            for (let _field of fields) {
              if (arr.includes(_field.code)) {
                _field.disabled = false
              }
            }
          }
        },

        // CSB0116RC0023000
        disable23({ op, fieldsList, focusFieldsIndex, codeValues = {} }) {
          const { fc204 } = codeValues

          const fields = fieldsList[focusFieldsIndex]
          const fieldCode = ['fc003', 'fc004', 'fc006', 'fc007', 'fc008', 'fc009', 'fc012', 'fc183', 'fc194', 'fc195', 'fc196', 'fc197', 'fc205', 'fc206', 'fc207']

          if (fc204 == '1') {
            fields.map(_field => {
              if (fieldCode.includes(_field.code)) {
                _field[`${op}Value`] = ''
                _field.resultValue = ''
                _field.disabled = true
              }
            })
          }
        },

        // CSB0116RC0034000
        disable34({ bill, fieldsList, focusFieldsIndex }) {
          if (bill.insuranceType == '1') {
            const fields = fieldsList[focusFieldsIndex]
            const fieldCode = ['fc004', 'fc005', 'fc006', 'fc007', 'fc008', 'fc009', 'fc010', 'fc011', 'fc012', 'fc017', 'fc019', 'fc022', 'fc023', 'fc024', 'fc025', 'fc183', 'fc184', 'fc185', 'fc186', 'fc187', 'fc188', 'fc189', 'fc190', 'fc191', 'fc192', 'fc193', 'fc194', 'fc195', 'fc196', 'fc197', 'fc198', 'fc199', 'fc200', 'fc203', 'fc204', 'fc205', 'fc206', 'fc207', 'fc208', 'fc209', 'fc210', 'fc211', 'fc212', 'fc213', 'fc214']

            fields.map(_field => {
              if (fieldCode.includes(_field.code)) {
                _field.disabled = true
              }
            })
          }
        },

        // CSB0116RC0028000
        disable28({ fieldsList, focusFieldsIndex, codeValues = {} }) {
          const { fc111 } = codeValues
          const fields = fieldsList[focusFieldsIndex]

          const fc215Field = fields.find(field => field.code == 'fc215')

          if (fc215Field && (fc111 == '1' || fc111 == '2')) {
            fc215Field.disabled = false
          } else {
            if (fc215Field) fc215Field.disabled = true
          }
        },
      }
    },

    // 回车
    enter: {
      methods: {
        // CSB0116RC0013000
        // disable13({ op, fieldsList, focusFieldsIndex }) {
        //   const fields = fieldsList[focusFieldsIndex]
        //   const fc008Field = fields.find(field => field.code == 'fc008')
        //   const fc009Field = fields.find(field => field.code == 'fc009')
        //   if (fc008Field && fc008Field.resultValue != 'A' || fc008Field && fc008Field.resultValue == '') {
        //     fc009Field[`${op}Value`] = ''
        //     fc009Field.resultValue = ''
        //     fc009Field.disabled = true
        //   } else {
        //     if (fc008Field) fc009Field.disabled = false
        //   }
        // },

        // CSB0116RC0017000
        disable17({ op, field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc014') return
          const fields = fieldsList[focusFieldsIndex]

          const fc015Field = fields.find(field => field.code == 'fc015')
          const fc016Field = fields.find(field => field.code == 'fc016')
          if (field.resultValue != '2') {
            fc015Field[`${op}Value`] = ''
            fc015Field.resultValue = ''
            fc016Field[`${op}Value`] = ''
            fc016Field.resultValue = ''
            fc015Field.disabled = true
            fc016Field.disabled = true
          } else {
            fc015Field.disabled = false
            fc016Field.disabled = false
          }
        },

        // CSB0116RC0022000
        validate21And211and212and1: function ({ op, field, fieldsList, focusFieldsIndex, memoFields }) {
          if (field.code === 'fc023') {
            const fields = fieldsList[focusFieldsIndex]

            const fc023Value = field.resultValue
            const fc022Field = fields.find(field => field.code === 'fc022')
            const fc021Field = fields.find(field => field.code === 'fc021')

            if (fc023Value.includes('-')) {
              const values = fc023Value.split('-')

              field[`${op}Value`] = values[2]
              field.resultValue = values[2]
              _.set(memoFields, `${field.uniqueId}.value`, values[2])

              fc022Field[`${op}Value`] = values[1]
              fc022Field.resultValue = values[1]
              _.set(memoFields, `${fc022Field.uniqueId}.value`, values[1])

              fc021Field[`${op}Value`] = values[0]
              fc021Field.resultValue = values[0]
              _.set(memoFields, `${fc021Field.uniqueId}.value`, values[0])
            }
          }
        },

        validate21And211and212and12: function ({ op, field, fieldsList, focusFieldsIndex, memoFields }) {
          if (field.code === 'fc022') {
            const fields = fieldsList[focusFieldsIndex]

            const fc022Value = field.resultValue
            const fc021Field = fields.find(field => field.code === 'fc021')

            if (fc022Value.includes('-')) {
              const values = fc022Value.split('-')

              field[`${op}Value`] = values[1]
              field.resultValue = values[1]
              _.set(memoFields, `${field.uniqueId}.value`, values[1])

              fc021Field[`${op}Value`] = values[0]
              fc021Field.resultValue = values[0]
              _.set(memoFields, `${fc021Field.uniqueId}.value`, values[0])
            }
          }
        },

        validate21And211and212and2: function ({ op, field, fieldsList, focusFieldsIndex, memoFields }) {
          if (field.code === 'fc194') {
            const fields = fieldsList[focusFieldsIndex]

            const fc194Value = field.resultValue
            const fc195Field = fields.find(field => field.code === 'fc195')
            const fc196Field = fields.find(field => field.code === 'fc196')

            if (fc194Value.includes('-')) {
              const values = fc194Value.split('-')

              field[`${op}Value`] = values[2]
              field.resultValue = values[2]
              _.set(memoFields, `${field.uniqueId}.value`, values[2])

              fc195Field[`${op}Value`] = values[1]
              fc195Field.resultValue = values[1]
              _.set(memoFields, `${fc195Field.uniqueId}.value`, values[1])

              fc196Field[`${op}Value`] = values[0]
              fc196Field.resultValue = values[0]
              _.set(memoFields, `${fc196Field.uniqueId}.value`, values[0])
            }
          }
        },

        validate21And211and212and21: function ({ op, field, fieldsList, focusFieldsIndex, memoFields }) {
          if (field.code === 'fc195') {
            const fields = fieldsList[focusFieldsIndex]

            const fc195Value = field.resultValue
            const fc196Field = fields.find(field => field.code === 'fc196')

            if (fc195Value.includes('-')) {
              const values = fc195Value.split('-')

              field[`${op}Value`] = values[1]
              field.resultValue = values[1]
              _.set(memoFields, `${field.uniqueId}.value`, values[1])

              fc196Field[`${op}Value`] = values[0]
              fc196Field.resultValue = values[0]
              _.set(memoFields, `${fc196Field.uniqueId}.value`, values[0])
            }
          }
        },

        validate21And211and212and3: function ({ op, field, fieldsList, focusFieldsIndex, memoFields }) {
          if (field.code === 'fc210') {
            const fields = fieldsList[focusFieldsIndex]

            const fc210Value = field.resultValue
            const fc211Field = fields.find(field => field.code === 'fc211')
            const fc212Field = fields.find(field => field.code === 'fc212')

            if (fc210Value.includes('-')) {
              const values = fc210Value.split('-')

              field[`${op}Value`] = values[2]
              field.resultValue = values[2]
              _.set(memoFields, `${field.uniqueId}.value`, values[2])

              fc211Field[`${op}Value`] = values[1]
              fc211Field.resultValue = values[1]
              _.set(memoFields, `${fc211Field.uniqueId}.value`, values[1])

              fc212Field[`${op}Value`] = values[0]
              fc212Field.resultValue = values[0]
              _.set(memoFields, `${fc212Field.uniqueId}.value`, values[0])
            }
          }
        },

        validate21And211and212and22: function ({ op, field, fieldsList, focusFieldsIndex, memoFields }) {
          if (field.code === 'fc211') {
            const fields = fieldsList[focusFieldsIndex]

            const fc211Value = field.resultValue
            const fc212Field = fields.find(field => field.code === 'fc212')

            if (fc211Value.includes('-')) {
              const values = fc211Value.split('-')

              field[`${op}Value`] = values[1]
              field.resultValue = values[1]
              _.set(memoFields, `${field.uniqueId}.value`, values[1])

              fc212Field[`${op}Value`] = values[0]
              fc212Field.resultValue = values[0]
              _.set(memoFields, `${fc212Field.uniqueId}.value`, values[0])
            }
          }
        },

        // CSB0116RC0031000
        disable31({ op, field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc052') return
          const fields = fieldsList[focusFieldsIndex]

          const fc053Field = fields.find(field => field.code == 'fc053')
          if (field.resultValue == '2') {
            fc053Field[`${op}Value`] = ''
            fc053Field.resultValue = ''
            fc053Field.disabled = true
          } else {
            fc053Field.disabled = false
          }
        },

        // CSB0116RC0029000
        disable29({ field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc021') return
          const fields = fieldsList[focusFieldsIndex]
          const fieldCode = ['fc131', 'fc132', 'fc133', 'fc134', 'fc135', 'fc136', 'fc137', 'fc138']

          if (field.resultValue != '' && !field.resultValue == 'A' && !field.resultValue.includes('?') && !field.resultValue.includes('？')) {
            fields.map(_field => {
              if (fieldCode.includes(_field.code)) {
                _field.table = {
                  name: `'B0116_华夏人寿团险理赔_${field.resultValue}'`,
                  query: '项目名称'
                }
              }
            })
          }
        },
      }
    },

    // 临时保存
    sessionSave: {
      methods: {

      }
    },

    // 提交前
    beforeSubmit: {
      methods: {

      }
    }
  }
}

export default B0116
