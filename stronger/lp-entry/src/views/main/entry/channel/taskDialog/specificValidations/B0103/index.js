import { tools, sessionStorage } from 'vue-rocket'
import moment from 'moment'
import { ignoreFreeValue } from '../tools'
import { MessageBox, Notification } from 'element-ui';
import { reqConst } from '@/api/reqConst'

// 两次F8通过
// let flags = true
// 记录数据库查询编码
let medicalValues = ''
let codeMaps = ''

const B0103 = {
  op0: {
    // 记录最后一次存储的合法field
    memoFields: [],

    // 记录相同 code 的 field 的值
    memoFieldValues: ['fc018', 'fc003'],

    // fields 的值从 targets 里的值选择
    dropdownFields: [
      {
        targets: ['fc018'],
        fields: ['fc019', 'fc020', 'fc085']
      }
    ],


    // 校验规则
    rules: [
      // 第一道工序:fc018录入内容不能重复，否则出错误提示:该发票属性已被使用，请重新定义(不允许强过)
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
      // 第一道工序:fc019,fc020和fc085录入内容必须为fc018的录入内容，否则出错误提示：没有此发票，请核实
      {
        fields: ['fc019', 'fc020', 'fc085'],
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
        setConstants01({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc006']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0103_广西贵州国寿理赔_医疗机构52',
                query: '医院名称'
              }
            }
          })
        },
      }
    },

    // 回车
    enter: {
      methods: {
        disable01({ op, field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc090') return

          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue == '2') {
            fields.map(_field => {
              if (_field.code == 'fc091') {
                // _field[`${op}Value`] = ''
                // _field.resultValue = ''
                _field.disabled = true
              }
            })
          } else {
            fields.map(_field => {
              if (_field.code == 'fc091') {
                _field.disabled = false
              }
            })
          }
        },
        disable02({ fieldsList, focusFieldsIndex }) {
          const fields = fieldsList[focusFieldsIndex]
          fields.map(_field => {
            if (_field.code == 'fc098') {
              _field.disabled = true
            }
          })
        }
      },
    },

    // F8(提交前校验)
    beforeSubmit: {
      methods: {
        // validate01: function ({ bill }) {
        //   // fieldsObjects = []
        //   let otherInfo = bill.otherInfo
        //   let bpoSendRemark = JSON.parse(otherInfo).bpoSendRemark
        //   if (bpoSendRemark && !bpoSendRemark.includes('【') && flags) {
        //     flags = false
        //     return alert(bpoSendRemark)
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
        // 发票xxx没有匹配的清单或报销单
        validate02({ fieldsObject }) {
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
        // fc003录入多少个4就必须存在多少个fc085， 诊断书无对应属性，请核实
        validate03({ fieldsObject }) {
          const fc003Values = []
          const fc085Values = []

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  const { code, resultValue } = field

                  if (code === 'fc003' && resultValue == '4') {
                    fc003Values.push(+resultValue)
                  }
                  if (code === 'fc085' && resultValue) {
                    fc085Values.push(+resultValue)
                  }
                }
              }
            }
          }

          if (fc003Values.length == fc085Values.length) {
            return true
          }

          return {
            errorMessage: '诊断书无对应属性，请核实'
          }
        },
        // 同一发票录入报销单与清单
        validate04({ fieldsObject }) {
          const [fc018Values, fc019Values, fc020Values, fc089Values, fc089Values1, fc018Values1] = [[], [], [], [], [], []]
          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  const { code, resultValue } = field

                  if (code === 'fc018' && !fc018Values.includes(resultValue)) {
                    resultValue && fc018Values.push(resultValue)
                  }

                  if (code === 'fc019' && !fc019Values.includes(resultValue)) {
                    resultValue && fc019Values.push(resultValue)
                  }

                  if (code === 'fc020' && !fc020Values.includes(resultValue)) {
                    resultValue && fc020Values.push(resultValue)
                  }
                  // 若fc089的录入值为1，标记fc089Values对应的分组方便拿到 fc018的值
                  if (code === 'fc089') {
                    resultValue && fc089Values.push(resultValue)
                    if (resultValue == '1') {
                      fc089Values1.push(fields)
                    }
                  }
                }
              }
            }
          }
          // 判断fc018 fc019 fc020 录入的值是否相同
          if ((fc019Values.length + fc020Values.length) > fc018Values.length) {
            return {
              popup: 'confirm',
              errorMessage: `同一发票录入报销单与清单`
            }
          }
          // 通过遍历fc089值为1的数组找出对应的fc018的值存放在fc018Values1数组中
          for (let fields of fc089Values1) {
            for (let field of fields) {
              if (field.code == 'fc018') {
                fc018Values1.push(field.resultValue)
              }
            }
          }
          // 判断fc089录入内容为1时， fc018的值是否与fc019的值相同
          if (fc018Values1.length != 0 && fc019Values.length != 0) {
            let fc018 = new Set(fc018Values1);
            let fc019 = new Set(fc019Values);
            let newSet = new Set(
              [...fc018].filter(item => fc019.has(item))
            )
            if (newSet.length != 0) {
              return {
                popup: 'confirm',
                errorMessage: `同一发票录入报销单与清单`
              }
            }
          }
          return true
        },

        // CSB0103RC0097000
        validate05({ sameFieldValue }) {
          if (!sameFieldValue.fc003?.values.includes('4')) {
            return {
              errorMessage: `漏切诊断书`
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
        validateDate02: function ({ value }) {
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
        validate03: function ({ fieldsIndex, fieldsList }) {
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
        fields: ['fc006'],
        validate006: function ({ effectValidations, field, items, value }) {
          if (ignoreFreeValue({ effectValidations, value })) return true

          const result = items.find((text) => text === value)

          if (value.includes('?')) {
            field.allowForce = true
            return true
          }
          else {
            // field.allowForce = false
          }

          if (!result) {
            return '医院名称错误，请确认'
          }

          field.allowForce = true

          return true
        }
      },
      {
        fields: ['fc007', 'fc074', 'fc075', 'fc076', 'fc077', 'fc078', 'fc079', 'fc080', 'fc081', 'fc082'],
        validate04: function ({ value, items }) {
          if (value === '?' || value === 'A') {
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
        fields: ['fc012', 'fc025', 'fc026', 'fc027', 'fc028', 'fc029', 'fc030', 'fc031'],
        validate05: function ({ value, items }) {
          if (value === '?' || value === '？') {
            return true
          }

          const result = items.find((text) => text === value)

          if (!result) {
            return '项目名称必须优先录入数据库内容，否则容易引起投诉， 请再次确认并修改'
          }

          return true
        }
      },

      {
        fields: ['fc014', 'fc039', 'fc040', 'fc041', 'fc042', 'fc043', 'fc044', 'fc045'],
        validate03: async function ({ block, fieldsList, value, field }) {
          if (block.code !== 'bc002') return true

          const fields = fieldsList[0]
          // 数据库
          const db = window['constantsDB']['B0103']
          if (!db) return
          const collections = db[`B0103_广西贵州国寿理赔_医疗目录${medicalValues}`]
          const maps = []

          const codesMap = new Map([
            ['fc014', 'fc012'],
            ['fc039', 'fc025'],
            ['fc040', 'fc026'],
            ['fc041', 'fc027'],
            ['fc042', 'fc028'],
            ['fc043', 'fc029'],
            ['fc044', 'fc030'],
            ['fc045', 'fc031'],
          ])

          let proList = JSON.parse(window.sessionStorage.getItem('proList'))
          let constant

          if (proList && codeMaps == '') {
            const data = {
              proCode: 'B0103',
              name: `B0103_广西贵州国寿理赔_医疗目录${medicalValues}`,
              queryNames: {
                ['项目名称']: {
                  $regex: `/${value}/`
                }
              },
              respNames: [
                '最高限价', '项目名称'
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
            // 将项目名称和最高限价推入map数组maps中
            for (let dessert of constant) {
              maps.push([dessert['项目名称'], dessert['最高限价']])
            }
            codeMaps = new Map([...maps])
          } else {
            // 将项目名称和最高限价推入map数组maps中
            for (let dessert of collections.desserts) {
              maps.push([dessert[2], dessert[0]])
            }
            codeMaps = new Map([...maps])
          }


          if (codesMap.get(field.code)) {
            for (let _field of fields) {
              if (_field.code == codesMap.get(field.code)) {
                let price = codeMaps.get(_field.resultValue)
                if (price == '') return true
                if (value > price) return `${_field.resultValue}价格超限额${price}, 请确认`
              }
            }
          }

          return true
        },



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
        fields: ['fc046'],
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
                name: 'B0103_广西贵州国寿理赔_医疗机构52',
                query: '医院名称'
              }
            }
          })
        },
        setConstants02({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc007', 'fc074', 'fc075', 'fc076', 'fc077', 'fc078', 'fc079', 'fc080', 'fc081', 'fc082']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0103_广西贵州国寿理赔_ICD10疾病编码',
                query: '疾病名称'
              }
            }
          })
        },
        setConstants03({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc062']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0103_广西贵州国寿理赔_第三方出具单位',
                query: '第三方出具单位名称'
              }
            }
          })
        },
        // 默认屏蔽
        disable05({ fieldsList, focusFieldsIndex }) {
          const codesList = [
            ['fc008', 'fc055', 'fc056', 'fc057', 'fc058', 'fc059', 'fc060', 'fc061'],
            ['fc047', 'fc048', 'fc049', 'fc050', 'fc051', 'fc052', 'fc053', 'fc054'],
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
        // 机构不包含5221 屏蔽fc015
        disable04({ bill, fieldsList, focusFieldsIndex }) {
          if (!bill.agency.includes('5221')) {
            const fields = fieldsList[focusFieldsIndex]
            fields?.map(_field => {
              if (_field.code == 'fc015') {
                _field.disabled = true
              }
            })
          }
        },
        async set48Items({ bill, fieldsList, focusFieldsIndex }) {
          const dropArr = ['fc012', 'fc025', 'fc026', 'fc027', 'fc028', 'fc029', 'fc030', 'fc031']
          // 案件的机构号
          const agency = bill.agency
          // 本级机构代码: agencyCode heighAgencyCode上级机构代码 医疗目录编码: medicalCode
          const [agencyCode, heighAgencyCode, medicalCode] = [[], [], []]
          // 本级机构代码对应的医疗目录编码
          let medicalValue = ''

          const prefix = 'B0103_广西贵州国寿理赔_医疗目录'
          let proList = JSON.parse(window.sessionStorage.getItem('proList'))
          let constant

          if (agency && proList) {
            const data = {
              proCode: 'B0103',
              name: 'B0103_广西贵州国寿理赔_数据库编码对应表',
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
            medicalValues = medicalValue
          } else {
            // 数据库
            const db = window['constantsDB']['B0103']
            if (!db) return
            const collections = db['B0103_广西贵州国寿理赔_数据库编码对应表']
            for (let dessert of collections.desserts) {
              agencyCode.push(dessert[1])
              heighAgencyCode.push(dessert[2])
              medicalCode.push(dessert[3])
            }

            const index = agencyCode.indexOf(agency)

            medicalValue = medicalCode[index]
            medicalValues = medicalValue
          }
          const fields = fieldsList[focusFieldsIndex]

          fields.map(_field => {
            if (dropArr.includes(_field.code)) {
              _field.table = {
                name: `${prefix}${medicalValue}`,
                query: '项目名称'
              }
            }
          })
        },

        disabled100({ block, fieldsList, focusFieldsIndex }) {
          if (block.code !== 'bc001') return true
          const fields = fieldsList[focusFieldsIndex]
          fields?.map(_field => {
            if (_field.code == 'fc017') {
              _field.disabled = true
            }
          })
        },

        disable01({ fieldsList, focusFieldsIndex, codeValues = {} }) {
          const { fc090 } = codeValues
          const codesList = ['fc092', 'fc094', 'fc095', 'fc096']
          const fields = fieldsList[focusFieldsIndex]
          fields?.map(_field => {
            if (fc090 == '2')
              if (codesList.includes(_field.code)) {
                _field.disabled = true
              }
          })
        },

        disable03({ fieldsList, focusFieldsIndex }) {
          const codesList = ['fc097']
          const fields = fieldsList[focusFieldsIndex]
          fields?.map(_field => {
            if (codesList.includes(_field.code)) {
              _field.disabled = true
            }
          })
        },

        disable06({ fieldsList, focusFieldsIndex, codeValues = {} }) {
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
      }
    },

    // 回车
    enter: {
      methods: {
        // fc021
        disable01({ op, field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc021') return

          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue == '2') {
            fields.map(_field => {
              if (_field.code == 'fc010') {
                // _field[`${op}Value`] = ''
                // _field.resultValue = ''
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
        // fc009
        validate02({ op, field, fieldsList, focusFieldsIndex }) {
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
        // bc004录入内容为A
        disable03({ field, block, fieldsList, focusFieldsIndex }) {
          if (block.code !== 'bc004') return
          const codesMap = new Map([
            ['fc074', ['fc075', 'fc076', 'fc077', 'fc078', 'fc079', 'fc080', 'fc081', 'fc082']],
            ['fc075', ['fc076', 'fc077', 'fc078', 'fc079', 'fc080', 'fc081', 'fc082']],
            ['fc076', ['fc077', 'fc078', 'fc079', 'fc080', 'fc081', 'fc082']],
            ['fc077', ['fc078', 'fc079', 'fc080', 'fc081', 'fc082']],
            ['fc078', ['fc079', 'fc080', 'fc081', 'fc082']],
            ['fc079', ['fc080', 'fc081', 'fc082']],
            ['fc080', ['fc081', 'fc082']],
            ['fc081', ['fc082']],
          ])
          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue === 'A') {
            const getCode = codesMap.get(field.code)
            fields.map(_field => {
              if (getCode && getCode.includes(_field.code)) {
                _field.disabled = true
              }
            })
          }
          else {
            const getCode = codesMap.get(field.code)
            fields.map(_field => {
              if (getCode && getCode.includes(_field.code)) {
                _field.disabled = false
              }
            })
          }
        },
        // 左边字段有录入值， 屏蔽右边字段， 同时赋值1
        disable04({ op, field, fieldsList, focusFieldsIndex }) {
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
        // bc003录入内容为A
        disable05({ field, block, fieldsList, focusFieldsIndex }) {
          if (block.code !== 'bc003') return
          const codesMap = new Map([
            ['fc063', ['fc062', 'fc064', 'fc065', 'fc066', 'fc067', 'fc068', 'fc069', 'fc070', 'fc071', 'fc072', 'fc073']],
            ['fc065', ['fc064', 'fc066', 'fc067', 'fc068', 'fc069', 'fc070', 'fc071', 'fc072', 'fc073']],
            ['fc067', ['fc066', 'fc068', 'fc069', 'fc070', 'fc071', 'fc072', 'fc073']],
            ['fc069', ['fc068', 'fc070', 'fc071', 'fc072', 'fc073']],
            ['fc071', ['fc070', 'fc072', 'fc073']],
            ['fc073', ['fc072']],
          ])
          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue === 'A') {
            const getCode = codesMap.get(field.code)
            fields.map(_field => {
              if (getCode && getCode.includes(_field.code)) {
                _field.disabled = true
              }
            })
          }
          else {
            const getCode = codesMap.get(field.code)
            fields.map(_field => {
              if (getCode && getCode.includes(_field.code)) {
                _field.disabled = false
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
        validate01: function ({ block, fieldsList }) {
          if (block.code !== 'bc002') return true

          const fields = fieldsList[0]
          const codes = ['fc012', 'fc025', 'fc026', 'fc027', 'fc028', 'fc029', 'fc030', 'fc031']
          for (let field of fields) {
            if (codes.includes(field.code) && (field.resultValue.includes('?') || field.resultValue.includes('？'))) {
              return {
                errorMessage: `清单名称不可包含?`
              }
            }
          }
          return true
        },

        // CSB0103RC0094000
        validate02: function ({ block, fieldsList }) {
          if (block.code !== 'bc001') return true

          const fields = fieldsList[0]
          const codes = ['fc005', 'fc009', 'fc011']
          for (let field of fields) {
            if (codes.includes(field.code) && field.resultValue == '') {
              return {
                errorMessage: `票据号、诊疗方式、发票总金额不能为空，请检查`
              }
            }
          }
          return true
        },

        // CSB0103RC0093000
        validate03({ block, fieldsList }) {
          if (block.code !== 'bc002') return true

          const flatFieldsList = tools.flatArray(fieldsList)
          const values = []

          flatFieldsList?.map(_field => {
            _field.resultValue && values.push(_field.resultValue)
          })

          if (!values?.length) {
            return {
              errorMessage: '清单不能空白提交，请检查，如清单内容无法录入则录入一组数据后按F8提交!'
            }
          }

          return true
        },

        // CSB0103RC0098000
        validate04: function ({ block, fieldsList }) {
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

        // CSB0103RC0124000
        validate05: function ({ block, fieldsList }) {
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

export default B0103
