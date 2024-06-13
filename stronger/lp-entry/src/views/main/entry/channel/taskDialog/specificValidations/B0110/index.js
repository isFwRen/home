import BigNumber from 'bignumber.js'
import { tools, sessionStorage } from 'vue-rocket'
import { getNode } from './tools'
import moment from 'moment'
import { MessageBox, Notification } from 'element-ui';
import { reqConst } from '@/api/reqConst'

// 两次F8通过
let flags = true
// 记录fieldObject
let fieldsObjects = []
// 记录fc019对应的发票下的清单页码
let pages = []


const B0110 = {
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
      // 第一道工序:fc018录入内容不能重复，否则出错误提示:该发票属性已被使用，请重新定义(不允许强过)
      {
        fields: ['fc018'],
        validate1: function ({ field, fieldsObject, thumbIndex, value }) {
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
      // 第一道工序:fc019和fc020录入内容必须为fc018的录入内容，否则出错误提示：没有此发票，请核实
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
      // 第一道工序:fc019录入内容相同时，fc008录入内容需按从小到大的顺序录入且不能重复，否则出错误提示：该清单页码录入错误，请重新定义(不允许强过)
      // {
      //   fields: ['fc008'],
      //   validate2: function ({ field, fieldsObject, thumbIndex, value }) {
      //     if (field.code != 'fc008') return true
      //     const fc019Values = []
      //     const fc008Values = []
      //     for (let key in fieldsObject) {
      //       const sessionStorage = fieldsObject[key].sessionStorage

      //       if (sessionStorage || thumbIndex === +key) {
      //         const _fieldsList = fieldsObject[key].fieldsList

      //         for (let _fields of _fieldsList) {
      //           for (let _field of _fields) {
      //             if (_field.code === 'fc019' && _field.uniqueId !== field.uniqueId) {
      //               fc019Values.push(_field.resultValue)
      //             }
      //             if (_field.code === 'fc008' && _field.uniqueId !== field.uniqueId) {
      //               fc008Values.push(_field.resultValue)
      //             }
      //           }
      //         }
      //       }
      //     }
      //     if (pages.includes(value)) {
      //       return '该清单页码录入错误，请重新定义(不允许强过)'
      //     }

      //     // 判断页码是否按照从小到大的顺序排列
      //     let sort = JSON.parse(JSON.stringify(pages))
      //     if (value) sort.push(value)
      //     for (let i = 0; i < sort.length; i++) {
      //       for (let j = i + 1; j < sort.length; j++) {
      //         if (sort[i] > sort[j]) {
      //           return '该清单页码录入错误，请重新定义(不允许强过)'
      //         }
      //       }
      //     }
      //     return true
      //   }
      // },
      {
        fields: ['fc006'],
        validate05: function ({ value, items }) {
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
        validateOtherInfo01: async function ({ bill }) {
          let otherInfo = bill.otherInfo
          let bpoSendRemark = JSON.parse(otherInfo).bpoSendRemark
          if (bpoSendRemark && !bpoSendRemark.includes('【')) {
            if (sessionStorage.get('isApp')?.isApp === 'true') {
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
        setConstants01({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc006']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0110_新疆国寿理赔_医疗机构65',
                query: '医院名称'
              }
            }
          })
        },

        setConstants02({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc074']

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
        disable02({ op, fieldsList, focusFieldsIndex }) {
          const fields = fieldsList[focusFieldsIndex]
          const fc066Field = fields.find(field => field.code == 'fc066')
          const fc067Field = fields.find(field => field.code == 'fc067')
          if (fc066Field && fc066Field.resultValue == '2') {
            fc067Field[`${op}Value`] = ''
            fc067Field.resultValue = ''
            fc067Field.disabled = true
          }
          if (fc066Field && fc066Field.resultValue != '2') {
            fc067Field.disabled = false
          }
        },
        // 对应rules第一道工序:fc019录入内容相同时
        disable03({ field, fieldsList }) {
          if (field.code !== 'fc019') return
          fieldsObjects.push(...fieldsList)
          for (let fields of fieldsObjects) {
            for (let _field of fields) {
              if (_field.code == 'fc019' && _field.resultValue == field.resultValue) {
                for (let _field of fields) {
                  if (_field.code == 'fc008') {
                    if (_field.resultValue != '' && !pages.includes(_field.resultValue)) {
                      pages.push(_field.resultValue)
                    }
                  }
                }
              }
            }
          }
        },
        // 对应rules第一道工序:fc019录入内容相同时 清空pages
        disable04({ field }) {
          if (field.code !== 'fc008') return
          pages = []
        },
      }
    },

    // F8(提交前校验)
    beforeSubmit: {
      methods: {
        // 对应rules第一道工序:fc019录入内容相同时 清空fieldsObjects
        // validate01: function ({ bill }) {
        //   fieldsObjects = []
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

        // CSB0110RC093000
        validate02({ sameFieldValue }) {
          if (!sameFieldValue.fc003?.values.includes('4')) {
            return {
              errorMessage: `缺少诊断书`
            }
          }
        },

        validate03({ fieldsObject }) {
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

        validate04({ sameFieldValue }) {
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
        fields: ['fc022', 'fc023'],
        validateDate03: function ({ value }) {
          if (!value) return true

          if (/[A, \?]/.test(value)) {
            return true
          }

          if (value.length !== 6 || moment(`20${value}`).format('YYYYMMDD') === 'Invalid date') {
            return '日期格式错误, 确认按单录入后， 强过处理'
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
        fields: ['fc007'],
        validate05: function ({ value, items }) {
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
        fields: ['fc005'],
        validate07: function ({ field, value, items }) {

          if (value.includes('?')) {
            return true
          }

          let { fc067 } = field.codeValues

          if (fc067 == '3' && value.length != 20) {
            return '票据号不够20位数，请检查！'
          }
          return true
        }
      },
    ],

    // 提示文本
    hints: [
      {
        fields: ['fc047'],
        hintFc001: function ({ bill }) {
          let otherInfo = bill.otherInfo
          let bpoSendRemark = JSON.parse(otherInfo).bpoSendRemark
          return `<p style="color: blue;">${bpoSendRemark}</p>`
        }
      },
    ],

    // 字段已生成
    updateFields: {
      methods: {
        setConstants01({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc006']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0110_新疆国寿理赔_医疗机构65',
                query: '医院名称'
              }
            }
          })
        },
        setConstants02({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc007']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0110_新疆国寿理赔_ICD10疾病编码',
                query: '疾病名称'
              }
            }
          })
        },
        setConstants03({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc061']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0110_新疆国寿理赔_第三方出具单位',
                query: '第三方出具单位名称'
              }
            }
          })
        },
        async disable04({ bill, fieldsList, focusFieldsIndex }) {
          // 数据库
          const db = window['constantsDB']['B0110']
          if (!db) return
          const collections = db['B0110_新疆国寿理赔_数据库编码对应表']
          // 案件的机构号
          const agency = bill.agency
          // 本级机构代码: agencyCode heighAgencyCode上级机构代码
          const [agencyCode, heighAgencyCode] = [[], []]

          let index
          let proList = JSON.parse(window.sessionStorage.getItem('proList'))
          let constant


          if (agency && proList) {
            console.log('客户端-----------');
            const data = {
              proCode: 'B0110',
              name: 'B0110_新疆国寿理赔_数据库编码对应表',
              queryNames: {
                ['本级机构代码']: agency
              },
              respNames: [
                '本级机构代码', '上级机构代码', '医疗目录编码'
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

            for (let dessert of constant) {
              agencyCode.push(dessert[0]['本机机构代码'])
              heighAgencyCode.push(dessert[0]['上级机构代码'])
            }
            index = agencyCode.indexOf(agency)
          } else {
            for (let dessert of collections.desserts) {
              agencyCode.push(dessert[1])
              heighAgencyCode.push(dessert[2])
            }
            index = agencyCode.indexOf(agency)
          }

          const fields = fieldsList[focusFieldsIndex]

          if (heighAgencyCode[index] != '650100' && heighAgencyCode[index] != '652200' && heighAgencyCode[index] != '652300') {
            fields.map(_field => {
              if (_field.code == 'fc017') {
                _field.disabled = true
              }
            })
          }
          if (heighAgencyCode[index] == '652200') {
            fields.map(_field => {
              if (_field.code == 'fc017') {
                _field.hint = `<p  style="color: blue; fontSize: 14px; margin-top: -3px; margin-bottom: 0px">机构为哈密,第三方报销只录入公补金额</p>`
              }
            })
          }
        },
        // 默认屏蔽
        disable06({ fieldsList, focusFieldsIndex }) {
          const codesList = ['fc069', 'fc072', 'fc073']

          const fields = fieldsList[focusFieldsIndex]

          fields?.map(_field => {
            if (codesList.includes(_field.code)) {
              _field.disabled = true
            }
          })
        },
        // 默认屏蔽
        disable05({ fieldsList, focusFieldsIndex }) {
          const codesList = [
            ['fc048', 'fc049', 'fc050', 'fc051', 'fc052', 'fc016'],
            ['fc053', 'fc054', 'fc055', 'fc056', 'fc057', 'fc058', 'fc059', 'fc060'],
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
        // 默认屏蔽 不同分块
        disable07({ op, fieldsList, focusFieldsIndex, codeValues = {} }) {
          const { fc066, fc067 } = codeValues
          const codesList = ['fc068', 'fc070', 'fc071']
          const fields = fieldsList[focusFieldsIndex]
          fields?.map(_field => {
            // if (fc066 == '2' || fc067 == '1')
            if (fc066 == '2')
              if (codesList.includes(_field.code)) {
                _field[`${op}Value`] = ''
                _field.resultValue = ''
                _field.disabled = true
              }
          })
        },
        async set48Items({ op, bill, fieldsList, focusFieldsIndex }) {
          const dropArr = ['fc012', 'fc025', 'fc026', 'fc027', 'fc028', 'fc029', 'fc030', 'fc031']
          // 案件的机构号
          const agency = bill.agency
          // 本级机构代码: agencyCode heighAgencyCode上级机构代码 医疗目录编码: medicalCode
          const [agencyCode, heighAgencyCode, medicalCode] = [[], [], []]
          // 本级机构代码对应的医疗目录编码
          let medicalValue = ''
          // 指定的医疗目录编码
          let index = ''

          const prefix = 'B0110_新疆国寿理赔_医疗目录'
          let proList = JSON.parse(window.sessionStorage.getItem('proList'))
          let constant

          if (agency && proList) {
            const data = {
              proCode: 'B0110',
              name: 'B0110_新疆国寿理赔_数据库编码对应表',
              queryNames: {
                ['本级机构代码']: agency
              },
              respNames: [
                '本级机构代码', '上级机构代码', '医疗目录编码'
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

            medicalValue = constant[0]['医疗目录编码']
          } else {
            // 数据库
            const db = window['constantsDB']['B0110']
            if (!db) return
            const collections = db['B0110_新疆国寿理赔_数据库编码对应表']
            for (let dessert of collections.desserts) {
              agencyCode.push(dessert[1])
              heighAgencyCode.push(dessert[2])
              medicalCode.push(dessert[3])
            }

            index = agencyCode.indexOf(agency)
            medicalValue = medicalCode[index]
          }

          const fields = fieldsList[focusFieldsIndex]

          if (heighAgencyCode[index] == '650100') {
            fields.map(_field => {
              if (dropArr.includes(_field.code)) {
                _field.table = {
                  name: `${prefix}6501000002`,
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

        disable08({ fieldsList, focusFieldsIndex, codeValues = {} }) {
          const { fc067 } = codeValues
          const codesList = ['fc068', 'fc071',]
          const fields = fieldsList[focusFieldsIndex]
          fields?.map(_field => {
            if (fc067 == '3')
              if (codesList.includes(_field.code)) {
                _field.disabled = true
              }
          })
        },

        // CSB0110RC00900000
        // disable08({ fieldsList, focusFieldsIndex, codeValues = {} }) {
        //   const { fc067 } = codeValues
        //   const codesList = ['fc068', 'fc070', 'fc071']
        //   const fields = fieldsList[focusFieldsIndex]
        //   const fc073Field = fields.find(field => field.code == 'fc073')
        //   fields?.map(_field => {
        //     if (codesList.includes(_field.code) && fc067 == '1' && fc073Field.resultValue == '') {
        //       _field.disabled = false
        //     }
        //   })
        // },
      }
    },

    // 回车
    enter: {
      methods: {
        disable01({ op, field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc021') return

          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue == '2') {
            fields.map(_field => {
              if (_field.code == 'fc010') {
                _field[`${op}Value`] = ''
                _field.resultValue = ''
                _field.disabled = true
              }
            })
          } else {
            fields.map(_field => {
              if (_field.code == 'fc010') {
                _field.disabled = false
              }
            })
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
                _field[`${op}Value`] = ''
                _field.resultValue = ''
              }
            })
          }
        },

        hint01111({ op, field, fieldsList, focusFieldsIndex }) {
          const codesMap = new Map([
            ['fc012', ['fc014', 'fc013']],
            ['fc025', ['fc039', 'fc032']],
            ['fc026', ['fc040', 'fc033']],
            ['fc027', ['fc041', 'fc034']],
            ['fc028', ['fc042', 'fc035']],
            ['fc029', ['fc043', 'fc036']],
            ['fc030', ['fc044', 'fc037']],
            ['fc031', ['fc045', 'fc038']]
          ])
          const codeList = ['fc012', 'fc025', 'fc026', 'fc027', 'fc028', 'fc029', 'fc030', 'fc031',]
          let content = ['床位', '人间', '病房', '陪护床', '走廊', '护理']
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
            const codes = fields.find(field => field.code == code[1])
            codes[`${op}Value`] = ''
            codes.resultValue = ''
            for (let field of fields) {
              if (field.code == code[0]) {
                field.hint = `<p  style="color: blue; fontSize: 14px; margin-top: -3px; margin-bottom: 0px">项目名称属于[床位费类],请按单录入对应单价</p>`
              }
              if (field.code == code[1]) {
                field.disabled = false
                field.hint = `<p  style="color: blue; fontSize: 14px; margin-top: -3px; margin-bottom: 0px">项目名称属于[床位费类],请按单录入对应数量</p>`
              }
            }
          } else {
            const code = codesMap.get(field.code)
            if (!code) return
            for (let field of fields) {
              if (field.code == code[0]) {
                field.hint = ''
              }
              if (field.code == code[1]) {
                field.disabled = true
                field.hint = ''
              }
            }
          }
        },

        // fc009
        validate05({ op, field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc009') return
          const fields = fieldsList[focusFieldsIndex]

          const fc023Field = fields.find(field => field.code == 'fc023')

          if (field.resultValue == '1') {
            fc023Field.disabled = true
          } else {
            fc023Field[`${op}Value`] = ''
            fc023Field.resultValue = ''
            fc023Field.disabled = false
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
          if (block.code !== 'bc001') return true
          const fields = fieldsList[0]
          for (let field of fields) {
            if (field.code == 'fc011' && field.resultValue == '') {
              return {
                errorMessage: `没有录入内容，请检查`
              }
            }
          }
          return true
        },

        validate02({ block, fieldsList }) {
          if (block.code !== 'bc002') return true

          const flatFieldsList = tools.flatArray(fieldsList)
          const values = []

          flatFieldsList?.map(_field => {
            _field.resultValue && values.push(_field.resultValue)
          })

          if (!values?.length) {
            return {
              errorMessage: '清单不能空白提交，请检查，如清单内容无法录入则录入一组数据后提交!'
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

export default B0110
