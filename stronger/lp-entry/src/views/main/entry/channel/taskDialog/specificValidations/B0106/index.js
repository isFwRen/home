import { tools, sessionStorage } from 'vue-rocket'
import moment from 'moment'
import { MessageBox, Notification } from 'element-ui';
import { reqConst } from '@/api/reqConst'

// 两次F8通过
let flags = true
// 记录医疗编码
let medical = ''

const B0106 = {
  op0: {
    // 记录最后一次存储的合法field
    memoFields: [],

    // 记录相同 code 的 field 的值
    memoFieldValues: ['fc018', 'fc003'],

    // fields 的值从 targets 里的值选择
    dropdownFields: [
      {
        targets: ['fc018'],
        fields: ['fc019', 'fc020']
      }
    ],

    // 校验规则
    rules: [
      {
        fields: ['fc018'],
        validate2: function ({ field, fieldsObject, thumbIndex, value }) {
          const fc018Values = []

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage

            if (sessionStorage || thumbIndex === +key) {
              const _fieldsList = fieldsObject[key].fieldsList

              for (let _fields of _fieldsList) {
                for (let _field of _fields) {
                  if (_field.code === 'fc018' && _field.uniqueId !== field.uniqueId) {
                    fc018Values.push(_field.resultValue)
                  }
                }
              }
            }
          }

          if (fc018Values.includes(value)) {
            return '该发票属性已被使用，请重新定义(不允许强过)'
          }

          return true
        }
      },
      {
        fields: ['fc019', 'fc020'],
        validate4: function ({ includes, value }) {
          if (includes) {
            const result = includes.find(text => text === value)

            if (!result) {
              return '没有此发票，请核实。'
            }
          }

          return true
        }
      },
      {
        fields: ['fc006', 'fc016'],
        validate05: function ({ value, items }) {
          if (value === '?') {
            return true
          }

          const result = items.find((text) => text === value)

          if (value === '蒲城县中医医院' || value === '蒲城县中医院') {
            return '应该录入:蒲城县中医医院（陕西发票专用）'
          }

          if (!result) {
            return '录入内容不在数据库中，请确认!'
          }

          return true
        }
      },
      {
        fields: ['fc006'],
        validate06: function ({ value }) {
          if (value === '?') {
            return true
          }

          if (value === '蒲城县中医医院' || value === '蒲城县中医院') {
            return '应该录入:蒲城县中医医院（陕西发票专用）'
          }
          return true
        }
      },
      {
        fields: ['fc003'],
        validate9: function ({ value, sameFieldValue }) {
          let fc003Values = sameFieldValue.fc003?.values

          let count = 0
          if (value == '4') {
            for (let el of fc003Values) {
              if (el == value) count++
              if (count >= 2) return '诊断书重复录入'
            }
          }
          return true
        }
      },
    ],

    // 提示文本
    hints: [
      {
        fields: ['fc003'],
        hintFc001: function ({ bill }) {
          let otherInfo = bill.otherInfo
          let bpoSendRemark = JSON.parse(otherInfo).bpoSendRemark
          return `<p style="color: red;">${bpoSendRemark}</p>`
        }
      },
    ],

    // 工序完成初始化
    init: {
      methods: {
        validateOtherInfo01: function ({ bill }) {
          let otherInfo = bill.otherInfo
          let bpoSendRemark = JSON.parse(otherInfo).bpoSendRemark
          if (bpoSendRemark && !bpoSendRemark.includes('【')) {
            if (sessionStorage.get('isApp')?.isApp === 'true') {
              // MessageBox.alert(bpoSendRemark, '请注意', {
              //   type: 'warning',
              //   confirmButtonText: '确定',
              //   showClose: false,
              // })
              return Notification({
                type: 'warning',
                title: '提醒(5s后自动关闭)',
                message: bpoSendRemark,
                duration: 5000,
                position: 'top-left'
              })
            } else {
              return alert(bpoSendRemark)
            }
          } else {
            // 【】外的内容
            let flag1 = bpoSendRemark.indexOf('【')
            let flag2 = bpoSendRemark.indexOf('】')
            let inner = bpoSendRemark.slice(flag1 + 1, flag2)
            let left = bpoSendRemark.slice(0, flag1)
            let right = bpoSendRemark.slice(flag2 + 1)
            let out = left + right
            if (out) {
              if (sessionStorage.get('isApp')?.isApp === 'true') {
                // MessageBox.alert(out, '请注意', {
                //   type: 'warning',
                //   confirmButtonText: '确定',
                //   showClose: false,
                // })
                return Notification({
                  type: 'warning',
                  title: '提醒(5s后自动关闭)',
                  message: out,
                  duration: 5000,
                  position: 'top-left'
                })
              } else {
                return alert(out)
              }
            }
          }
        }
      }
    },

    // 字段已生成
    updateFields: {
      methods: {
        setConstants40({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc006', 'fc016']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0106_陕西国寿理赔_医疗机构61',
                query: '医院名称'
              }
            }
          })
        },

        setConstants41({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc098']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.disabled = true
            }
          })
        },
      }
    },

    // 回车
    enter: {
      methods: {
        disable1({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc090') return

          const codes = ['fc091']
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
      }
    },

    // F8(提交前校验)
    beforeSubmit: {
      methods: {
        validate01({ fieldsObject }) {
          const [fc018Values, fc019Values, fc020Values] = [[], [], []]
          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  const { code, resultValue } = field

                  if (code === 'fc018') {
                    resultValue && fc018Values.push(resultValue)
                  }

                  if (code === 'fc019') {
                    resultValue && fc019Values.push(resultValue)
                  }

                  if (code === 'fc020') {
                    resultValue && fc020Values.push(resultValue)
                  }
                }
              }
            }
          }
          let value
          let uniArr = [...new Set([...fc019Values, ...fc020Values])]
          const flag = fc018Values.every(item => {
            value = item
            return uniArr.includes(item)
          })
          if (flag) {
            return true
          } else {
            return {
              errorMessage: `发票${value}没有匹配的清单或报销单,请检查`
            }
          }
        },
        // validate02: function ({ bill }) {
        //   // fieldsObjects = []
        //   let otherInfo = bill.otherInfo
        //   let bpoSendRemark = JSON.parse(otherInfo).bpoSendRemark
        //   if (bpoSendRemark && !bpoSendRemark.includes('【') && flags) {
        //     flags = false
        //     if (sessionStorage.get('isApp')?.isApp === 'true') {
        //       MessageBox.alert(bpoSendRemark, '请注意', {
        //         type: 'warning',
        //         confirmButtonText: '确定',
        //         showClose: false,
        //       })
        //     } else {
        //       return alert(bpoSendRemark)
        //     }
        //   }
        //   let flag1 = bpoSendRemark.indexOf('【')
        //   let flag2 = bpoSendRemark.indexOf('】')
        //   let inner = bpoSendRemark.slice(flag1 + 1, flag2)
        //   let left = bpoSendRemark.slice(0, flag1)
        //   let right = bpoSendRemark.slice(flag2 + 1)
        //   let out = left + right
        //   if (out && flags) {
        //     flags = false
        //     return {
        //       errorMessage: `${out}`
        //     }
        //   }
        //   return true
        // },
        validate05({ sameFieldValue }) {
          if (!sameFieldValue.fc003?.values.includes('4')) {
            return {
              errorMessage: `缺少诊断书`
            }
          }
        },
        validate06({ sameFieldValue }) {
          let fc003Values = sameFieldValue.fc003?.values
          let count = 0
          for (let el of fc003Values) {
            if (el == '4') count++
            if (count >= 2) return {
              errorMessage: '诊断书重复录入'
            }
          }
        }
      }
    }
  },

  op1op2opq: {
    // 校验规则
    rules: [
      {
        fields: ['fc006', 'fc016'],
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
      {
        fields: ['fc022', 'fc023'],
        validateDate03: function ({ value }) {
          if (!value) return true

          if (/[A, \?]/.test(value)) {
            return true
          }
          let today = moment(new Date).format('YYYYMMDD')

          if (value.length !== 6 || moment(`20${value}`).format('YYYYMMDD') === 'Invalid date') {
            return '日期格式错误, 确认按单录入后， 强过处理'
          }
          if (Number(value) > Number(today.slice(2))) {
            return '录入内容不能晚于当天日期'
          }

          return true
        }
      },
      {
        fields: ['fc023'],
        validate04: function ({ fieldsIndex, fieldsList }) {
          const fields = fieldsList[fieldsIndex]

          const fc022Field = tools.find(fields, { code: 'fc022' })
          const fc023Field = tools.find(fields, { code: 'fc023' })

          let fc022Value = fc022Field?.resultValue
          let fc023Value = fc023Field?.resultValue

          if (fc023Value.includes('?') || fc023Value.includes('？')) {
            return true
          }

          fc022Value = +fc022Value
          fc023Value = +fc023Value

          if (fc023Value < fc022Value) {
            return '出院日期早于入院日期'
          }

          return true
        }
      },
      {
        fields: ['fc007', 'fc008'],
        validate06: function ({ value, items }) {
          if (value === 'A') {
            return true
          }

          if (value.includes('?')) {
            return true
          }

          const result = items.find((text) => text === value)

          if (!result) {
            return '录入内容不在数据库中，请确认!'
          }
          return true
        }
      },

      {
        fields: ['fc005'],
        validate07: function ({ field, value, items }) {

          if (value.includes('?')) {
            return true
          }

          let { fc091 } = field.codeValues

          if (fc091 == '3' && value.length != 20) {
            return '票据号不够20位数，请检查！'
          }
          return true
        }
      },
    ],

    // 提示文本
    hints: [
      {
        fields: ['fc022'],
        hintFc022: function () {
          return '<p style="color: blue; margin-bottom: 0px">请按发票中最早的日期录入， 注意日期位置不固定</p>'
        }
      },
    ],

    // 字段已生成
    updateFields: {
      methods: {
        setConstants01({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc006', 'fc016']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0106_陕西国寿理赔_医疗机构61',
                query: '医院名称'
              }
            }
          })
        },
        setConstants02({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc007', 'fc008']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0106_陕西国寿理赔_ICD10疾病编码',
                query: '疾病名称'
              }
            }
          })
        },
        setConstants03({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc087']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0106_陕西国寿理赔_第三方出具单位',
                query: '第三方出具单位名称'
              }
            }
          })
        },
        disable04({ op, fieldsList }) {
          const codesList = [
            'fc050', 'fc052', 'fc053', 'fc054',
            'fc046', 'fc047', 'fc048',
            'fc055', 'fc059', 'fc063', 'fc067', 'fc071', 'fc075', 'fc079', 'fc083',
          ]

          fieldsList[0]?.map(_field => {
            if (codesList.includes(_field.code)) {
              _field[`${op}Value`] = ''
              _field.resultValue = ''
              _field.disabled = true
            }
          })
        },
        disable05({ block, fieldsList }) {
          if (block?.code.toLowerCase() === 'bc002') {

            fieldsList?.map(fields => {
              fields?.map(_field => {
                if (_field.code === 'fc015') {
                  _field.disabled = true
                }
              })
            })
          }
        },
        async set48Items({ bill, fieldsList, focusFieldsIndex }) {
          const dropArr = ['fc012', 'fc025', 'fc026', 'fc027', 'fc028', 'fc029', 'fc030', 'fc031']
          let otherInfo = bill.otherInfo
          let bpoSendRemark = JSON.parse(otherInfo).bpoSendRemark
          // 案件的机构号
          const agency = bill.agency
          // 本级机构代码: agencyCode heighAgencyCode上级机构代码 医疗目录编码: medicalCode
          const [agencyCode, heighAgencyCode, medicalCode] = [[], [], []]
          // 本级机构代码对应的医疗目录编码
          let medicalValue = ''
          // 指定的医疗目录编码
          let index = ''

          const prefix = 'B0106_陕西国寿理赔_医疗目录'
          let proList = JSON.parse(window.sessionStorage.getItem('proList'))
          let constant
          if (agency && proList) {
            const data = {
              proCode: 'B0106',
              name: 'B0106_陕西国寿理赔_数据库编码对应表',
              queryNames: {
                ['本级机构代码']: agency
              },
              respNames: [
                '本级机构代码', '医疗目录编码'
              ],
              pageSize: 10,
              pageIndex: 1
            }

            const result = await reqConst({
              url: '/sys-const/page',
              method: "POST",
              data,
            })

            constant = result.list
            console.log(result);
            medicalValue = constant[0]['医疗目录编码']
          } else {
            // 数据库
            const db = window['constantsDB']['B0106']
            if (!db) return
            const collections = db['B0106_陕西国寿理赔_数据库编码对应表']
            for (let dessert of collections.desserts) {
              agencyCode.push(dessert[1])
              heighAgencyCode.push(dessert[2])
              medicalCode.push(dessert[3])
            }

            index = agencyCode.indexOf(agency)

            medicalValue = medicalCode[index]
            medical = medicalValue
          }

          const fields = fieldsList[focusFieldsIndex]

          if (bpoSendRemark.includes('生育')) {
            fields.map(_field => {
              if (dropArr.includes(_field.code)) {
                _field.table = {
                  // name: `B0106_陕西国寿理赔_医疗目录M6101002022001`,
                  name: `${prefix}M6101002022001`,
                  query: '项目名称'
                }
              }
            })
          } else {
            fields.map(_field => {
              if (dropArr.includes(_field.code)) {
                _field.table = {
                  name: `${prefix}${medicalValue}`,
                  query: '项目名称'
                }
              }
            })
          }
        },
        disable06({ fieldsList }) {
          const codesList = ['fc093', 'fc096', 'fc097']

          fieldsList[0]?.map(_field => {
            if (codesList.includes(_field.code)) {
              _field.disabled = true
            }
          })
        },

        disable07({ fieldsList, focusFieldsIndex, codeValues = {} }) {
          const { fc090 } = codeValues
          const codesList = ['fc092', 'fc094', 'fc095']
          const fields = fieldsList[focusFieldsIndex]
          fields?.map(_field => {
            if (fc090 == '2') {
              if (codesList.includes(_field.code)) {
                _field.disabled = true
              }
            }
          })
        },

        disable08({ fieldsList, focusFieldsIndex, codeValues = {} }) {
          const { fc091 } = codeValues
          const codesList = ['fc092', 'fc095',]
          const fields = fieldsList[focusFieldsIndex]
          fields?.map(_field => {
            if (fc091 == '3')
              if (codesList.includes(_field.code)) {
                _field.disabled = true
              }
          })
        },

        // disable08({ fieldsList, focusFieldsIndex, codeValues = {} }) {
        //   const { fc091 } = codeValues
        //   const codesList = ['fc092', 'fc094', 'fc095']
        //   const fields = fieldsList[focusFieldsIndex]
        //   fields?.map(_field => {
        //     if (fc091 == '1') {
        //       if (codesList.includes(_field.code)) {
        //         _field.disabled = true
        //       }
        //     }
        //   })
        // },

        // disable09({ fieldsList, focusFieldsIndex, codeValues = {} }) {
        //   const { fc091 } = codeValues
        //   const codes = ['fc092', 'fc094', 'fc095', 'fc096']
        //   const fields = fieldsList[focusFieldsIndex]
        //   const fc097field = fields.find(field => field.code === 'fc097')
        //   if (fc091 == '1' && fc097field && fc097field.resultValue == '') {
        //     fields.map(_field => {
        //       if (codes.includes(_field.code)) {
        //         _field.disabled = false
        //       }
        //     })
        //   }
        // },
      }
    },

    // 回车
    enter: {
      methods: {
        // fc009
        validate05({ op, field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc009') return
          const fields = fieldsList[focusFieldsIndex]

          const fc023Field = fields.find(field => field.code == 'fc023')
          const fc017Field = fields.find(field => field.code == 'fc017')

          if (field.resultValue == '1') {
            fc023Field.disabled = true
          } else {
            fc023Field[`${op}Value`] = ''
            fc023Field.resultValue = ''
            fc023Field.disabled = false
          }
          if (field.resultValue == '2') {
            fc017Field[`${op}Value`] = ''
            fc017Field.resultValue = ''
            fc017Field.disabled = true
          } else {
            fc017Field.disabled = false
          }
        },
        validate06({ op, field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc022') return
          const fields = fieldsList[focusFieldsIndex]

          const fc009Field = fields.find(field => field.code == 'fc009')
          const fc022Field = fields.find(field => field.code == 'fc022')
          const fc023Field = fields.find(field => field.code == 'fc023')

          if (fc009Field && fc009Field.resultValue == '1') {
            fc023Field[`${op}Value`] = field[`${op}Value`]
            fc023Field.resultValue = field[`${op}Value`]
          }
        },
        disable02({ op, field, fieldsList, focusFieldsIndex }) {
          const codesListMap = new Map([
            ['fc012', 'fc013'],
            ['fc025', 'fc032'],
            ['fc026', 'fc033'],
            ['fc027', 'fc034'],
            ['fc028', 'fc035'],
            ['fc029', 'fc036'],
            ['fc030', 'fc037'],
            ['fc031', 'fc038'],
          ])

          const fields = fieldsList[focusFieldsIndex]
          const code = codesListMap.get(field.code)

          if (code && field.code != "") {
            fields.map(_field => {
              if (_field.code === code) {
                _field.disabled = true
                _field[`${op}Value`] = '1'
                _field.resultValue = '1'
              }
            })
          } else if (code && field.code == "") {
            fields.map(_field => {
              if (_field.code === code) {
                _field.disabled = false
                _field[`${op} Value`] = ''
                _field.resultValue = ''
              }
            })
          }
        },
        hint01({ op, field, fieldsList, focusFieldsIndex }) {
          const codesMap = new Map([
            ['fc012', 'fc013'],
            ['fc025', 'fc032'],
            ['fc026', 'fc033'],
            ['fc027', 'fc034'],
            ['fc028', 'fc035'],
            ['fc029', 'fc036'],
            ['fc030', 'fc037'],
            ['fc031', 'fc038']
          ])
          const codeList = ['fc012', 'fc025', 'fc026', 'fc027', 'fc028', 'fc029', 'fc030', 'fc031',]
          let content = ['床位', '人间', '病房', '陪护床', '走廊']
          let noContent = ['取暖', '空调']
          const fields = fieldsList[focusFieldsIndex]
          let flag1 = '1'
          let flag2 = '1'
          for (let i of content) {
            if (field.resultValue.includes(i)) {
              flag1 = '0'
              break
            }
          }
          for (let i of noContent) {
            if (field.resultValue.includes(i)) {
              flag2 = '0'
              break
            }
          }

          if (codeList.includes(field.code) && flag1 == '0' && flag2 == '1') {
            const code = codesMap.get(field.code)
            const codes = fields.find(field => field.code == code)
            codes[`${op}Value`] = ''
            codes.resultValue = ''
            for (let field of fields) {
              if (field.code == code) {
                field.disabled = false
                field.hint = `<p  style="color: blue; fontSize: 14px; margin-top: -3px; margin-bottom: 0px">项目名称属于[床位费类],请按单录入对应数量</p>`
              }
            }
          } else {
            const code = codesMap.get(field.code)
            if (!code) return
            for (let field of fields) {
              if (field.code == code) {
                field.disabled = true
                field.hint = ''
              }
            }
          }
        },
        validate10({ block, fieldsList }) {
          if (block.code !== 'bc002') return true

          const fields = fieldsList[0]
          // 数据库
          const db = window['constantsDB']['B0106']
          if (!db) return
          const collections = db[`B0106_陕西国寿理赔_医疗目录${medical}`]
          const maps = []

          const codesMap = new Map([
            ['fc012', 'fc014'],
            ['fc025', 'fc039'],
            ['fc026', 'fc040'],
            ['fc027', 'fc041'],
            ['fc028', 'fc042'],
            ['fc029', 'fc043'],
            ['fc030', 'fc044'],
            ['fc031', 'fc045'],
          ])
          const codes = ['fc012', 'fc025', 'fc026', 'fc027', 'fc028', 'fc029', 'fc030', 'fc031']
          // 将项目名称和最高限价推入map数组maps中
          for (let dessert of collections.desserts) {
            maps.push([dessert[2], dessert[0]])
          }
          const codeMaps = new Map([...maps])

          for (let field of fields) {
            if (codes.includes(field.code) && codesMap.get(field.code)) {
              if (codeMaps.get(field.resultValue)) {
                for (let _field of fields) {
                  if (_field.code == codesMap.get(field.code) && Number(_field.resultValue) > Number(codeMaps.get(field.resultValue))) {
                    _field.hint = `<p style="color: red; fontSize: 14px; margin-top: -3px; margin-bottom: 0px">${field.resultValue}价格超限额${codeMaps.get(field.resultValue)}, 请确认</p>`
                    return false
                  } else if (_field.code == codesMap.get(field.code) && Number(_field.resultValue) <= Number(codeMaps.get(field.resultValue))) {
                    _field.hint = ''
                  }
                }
              }
            }
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
        validate01({ block, fieldsList }) {
          if (block.code !== 'bc002') return true
          for (let fields of fieldsList) {
            for (let field of fields) {
              if (field.resultValue) return true
            }
          }
          return {
            errorMessage: '清单不能空白提交, 请检查, 如清单内容无法录入则录入一组数据后按F8提交.'
          }
        },

        // CSB0103RC0087000
        validate02: function ({ block, fieldsList }) {
          if (block.code !== 'bc001') return true

          const fields = fieldsList[0]
          // const codes = ['fc011']
          for (let field of fields) {
            if (field.code == 'fc011' && field.resultValue == '') {
              return {
                errorMessage: `没有录入内容，请检查`
              }
            }
          }
          return true
        },

        validate03: function ({ block, fieldsList }) {
          if (block.code !== 'bc002') return true

          const fields = fieldsList[0]
          const codes = ['fc012', 'fc025', 'fc026', 'fc027', 'fc028', 'fc029', 'fc030', 'fc031']
          const codeMap = new Map([
            ['fc012', 'fc014'],
            ['fc025', 'fc039'],
            ['fc026', 'fc040'],
            ['fc027', 'fc041'],
            ['fc028', 'fc042'],
            ['fc029', 'fc043'],
            ['fc030', 'fc044'],
            ['fc031', 'fc045'],
          ])
          for (let field of fields) {
            if (codes.includes(field.code) && field.resultValue != '') {
              for (let _field of fields) {
                if (_field.code == codeMap.get(field.code) && _field.resultValue == '') {
                  return {
                    errorMessage: `金额字段不能为空`
                  }
                }
              }
            }
          }
          return true
        },
      }
    }
  }
}

export default B0106
