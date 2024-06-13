import moment from 'moment'
import BigNumber from 'bignumber.js'
import { tools, sessionStorage } from 'vue-rocket'
import { MessageBox, Notification } from 'element-ui';


const B0118 = {
  op0: {
    // 记录最后一次存储的合法field
    memoFields: ['fc054', 'fc180', 'fc181'],

    // 记录相同 code 的 field 的值
    memoFieldValues: ['fc059'],

    // fields 的值从 targets 里的值选择
    dropdownFields: [
      {
        targets: ['fc059'],
        fields: ['fc060', 'fc061', 'fc201']
      }
    ],

    // 校验规则
    rules: [
      // 2
      {
        fields: ['fc059'],
        validate2: function ({ block, field, fieldsObject, thumbIndex, value }) {
          if (block?.code.toLowerCase() === 'bc001') {
            const fc082Values = []

            for (let key in fieldsObject) {
              const sessionStorage = fieldsObject[key].sessionStorage

              if (sessionStorage || thumbIndex === +key) {
                const _fieldsList = fieldsObject[key].fieldsList

                for (let _fields of _fieldsList) {
                  for (let _field of _fields) {
                    if (_field.code === 'fc082') {
                      _field.resultValue && fc082Values.push(_field.resultValue)
                    }
                  }
                }
              }
            }

            const allValueIs1 = !!fc082Values.length && fc082Values.every(v => v == 1)

            if (allValueIs1) {
              return true
            }

            const fc059Values = []

            for (let key in fieldsObject) {
              const sessionStorage = fieldsObject[key].sessionStorage

              if (sessionStorage || thumbIndex === +key) {
                const _fieldsList = fieldsObject[key].fieldsList

                for (let _fields of _fieldsList) {
                  for (let _field of _fields) {
                    if (_field.code === 'fc059' && _field.uniqueId !== field.uniqueId) {
                      fc059Values.push(_field.resultValue)
                    }
                  }
                }
              }
            }

            if (fc059Values.includes(value)) {
              return '发票属性不能重复！'
            }
          }

          return true
        }
      },

      // 4
      {
        fields: ['fc060', 'fc061', 'fc201'],
        validate6: function ({ value, items }) {
          const result = items.find(text => text === value)

          if (result) {
            return true
          }

          return '没有此发票，请核实！'
        }
      },

      // 23
      {
        fields: ['fc054'],
        validate23: function ({ value, items }) {
          const result = items.find(text => text === value)

          if (result) {
            return true
          }

          return '录入内容不在代码表中，请录入“本代码表中不存在的其他医院”.'
        }
      },

      // 33
      {
        fields: ['fc180'],
        validate33: function ({ value, items }) {
          // 录入内容为F时不做校验
          if (value === 'F') {
            return true
          }

          const result = items.find(text => text.includes(value))

          if (result) {
            return true
          }

          return '区不在常量表中时，录入F.'
        }
      },

      // 43
      {
        fields: ['fc152'],
        validate27: function ({ field, value, items }) {
          if (value === '?') return true

          const result = items.find(text => text === value)

          if (result) {
            return true
          }
          else {
            field.allowForce = false
            return '请选择常量表名称录入'
          }
        }
      },
    ],

    // 提示文本
    hints: [
      // {
      //   fields: ['fc152'],
      //   hintFc152: function () {
      //     return '<p style="color: blue;">冻干人用狂犬病疫苗、破伤风疫苗、破伤风类毒素需录入西药费；工本费、病历费、卡费、复印费、陪护费、陪人费此类大项名称按单录入强过；照相费选择放射费、处置费选择治疗费、材料费选择卫材费!</p>'
      //   }
      // },
      {
        fields: ['fc056'],
        hintFc162: function ({ bill, field }) {
          if (bill.saleChannel == '1') {
            if (field.resultValue == '3') {
              return '<p style="color: blue;">特约件乙类不用录</p>'
            }
          }
        }
      },
    ],

    // 工序完成初始化
    init: {
      methods: {
        // 14
        validateBillNum14: function ({ bill }) {
          const firstTwo = bill.billNum?.slice(0, 2)

          if (firstTwo === '30') {

            if (sessionStorage.get('isApp')?.isApp === 'true') {
              // MessageBox.alert('收据上有手术费的，需要切手术！', '请注意', {
              //   type: 'warning',
              //   confirmButtonText: '确定',
              //   showClose: false,
              // })
              return Notification({
                type: 'warning',
                title: '提醒(5s后自动关闭)',
                message: '收据上有手术费的，需要切手术！',
                duration: 5000,
                position: 'top-left'
              })
            } else {
              alert('收据上有手术费的，需要切手术！')
            }

          }
        }
      }
    },

    // 字段已生成
    updateFields: {
      methods: {
        // 23
        setConstants23: async function ({ flatFieldsList }) {
          const fields = ['fc054']
          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0118_中意理赔_医院代码表',
                query: '名称'
              }
            }
          })
        },

        // 30
        setConstants30: function ({ flatFieldsList }) {
          const fields = ['fc152']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0118_中意理赔_发票大项类型',
                query: '费用名称'
              }
            }
          })
        },

        // 33
        setConstants33: function ({ flatFieldsList }) {
          const fields = ['fc180']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0118_中意理赔_地址库',
                query: '市中文名',
                targets: ['省中文名', '市中文名']
              }
            }
          })
        },

        // 34
        setConstants34: function ({ flatFieldsList }) {
          const fields = ['fc181']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0118_中意理赔_地址库',
                query: '省中文名'
              }
            }
          })
        }
      }
    },

    // 回车
    enter: {
      methods: {
        // 33 34
        validate33And34: function ({ op, field, fieldsList, focusFieldsIndex, memoFields }) {
          if (field.code === 'fc180') {
            const fields = fieldsList[focusFieldsIndex]

            const fc180Value = field.resultValue
            const fc181Field = fields.find(field => field.code === 'fc181')

            if (fc180Value.includes('-')) {
              const values = fc180Value.split('-')

              field[`${op}Value`] = values[1]
              field.resultValue = values[1]
              _.set(memoFields, `${field.uniqueId}.value`, values[1])

              fc181Field[`${op}Value`] = values[0]
              fc181Field.resultValue = values[0]
              _.set(memoFields, `${fc181Field.uniqueId}.value`, values[0])
            }
          }
        },

        validate401({ field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc277') return
          const fieldArr = ['fc278']
          const fields = fieldsList[focusFieldsIndex]
          if (field.resultValue == '2') {
            fields?.map(_field => {
              if (fieldArr.includes(_field.code)) {
                _field.disabled = true
              }
            })
          } else {
            fields?.map(_field => {
              if (fieldArr.includes(_field.code)) {
                _field.disabled = false
              }
            })
          }
        },

        // 44
        // validate44({ field, fieldsList, focusFieldsIndex }) {
        //   if (field.code != 'fc061') return
        //   const fieldsLists = fieldsList
        //   const invoice = []

        //   for (let el of fieldsLists) {
        //     const fc059Field = el.find(field => field.code === 'fc059')
        //     const fc180Field = el.find(field => field.code === 'fc180')
        //     if (fc180Field && fc180Field.resultValue == '南京') invoice.push(fc059Field.resultValue)
        //     else if(fc180Field && fc180Field.resultValue != '南京') invoice.filter((el)=> el != fc059Field.resultValue)
        //   }

        //   const fields = fieldsList[focusFieldsIndex]
        //   const fc056Field = fields.find(field => field.code === 'fc056')
        //   if (invoice.includes(field.op0Value)) fc056Field.hint = `<p style="color: blue; fontSize: 14px; margin-top: -3px; margin-bottom: 0px">该案件地区为南京地区， 请核实是否正确</p>`
        //   else fc056Field.hint = ''         
        // }
      }
    },

    // F8(提交前校验)
    beforeSubmit: {
      methods: {
        // 5
        validate5({ mergeFieldsList, op }) {
          const flatFieldsList = tools.flatArray(mergeFieldsList)
          const fc060Field = flatFieldsList.find(field => field.code === 'fc060')

          if (fc060Field) {
            const [fc059Values, fc060Values] = [[], []]

            flatFieldsList.map(field => {
              if (field.code === 'fc059') {
                if (!fc059Values.includes(field[`${op}Value`])) {
                  fc059Values.push(field[`${op}Value`])
                }
              }

              if (field.code === 'fc060') {
                if (!fc060Values.includes(field[`${op}Value`])) {
                  fc060Values.push(field[`${op}Value`])
                }
              }
            })

            for (let value of fc060Values) {
              if (!fc059Values.includes(value)) {
                return {
                  errorMessage: `清单明细${value}没有匹配的发票，请检查！`
                }
              }
            }
          }

          return true
        },

        // 6
        validate6({ mergeFieldsList }) {
          const flatFieldsList = tools.flatArray(mergeFieldsList)
          const fc061Field = flatFieldsList.find(field => field.code === 'fc061')

          if (fc061Field) {
            const [fc059Values, fc061Values] = [[], []]

            flatFieldsList.map(field => {
              if (field.code === 'fc059') {
                if (!fc059Values.includes(field.resultValue)) {
                  fc059Values.push(field.resultValue)
                }
              }

              if (field.code === 'fc061') {
                if (!fc061Values.includes(field.resultValue)) {
                  fc061Values.push(field.resultValue)
                }
              }
            })

            for (let value of fc061Values) {
              if (!fc059Values.includes(value)) {
                return {
                  errorMessage: `报销单${value}没有匹配的发票，请检查！`
                }
              }
            }
          }

          return true
        },

        // 7
        validate7({ fieldsObject }) {
          const fc056Values = []

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  const { code, resultValue } = field

                  if (code === 'fc056') {
                    resultValue && fc056Values.push(resultValue)
                  }
                }
              }
            }
          }

          if (!fc056Values.includes('2')) {
            return {
              errorMessage: '没有录入发票，请确认!'
            }
          }

          return true
        },


        validate9({ mergeFieldsList }) {
          const flatFieldsList = tools.flatArray(mergeFieldsList)
          const fc059Field = flatFieldsList.find(field => field.code === 'fc059')
          const fc060Field = flatFieldsList.find(field => field.code === 'fc060')

          if (fc059Field) {
            const [fc059Values, fc201Values] = [[], []]

            flatFieldsList.map(field => {
              const op0Value = field.op0Value
              if (field.code === 'fc059') {
                !fc059Values.includes(op0Value) && fc059Values.push(op0Value)
              }

              if (field.code === 'fc201') {
                !fc201Values.includes(op0Value) && fc201Values.push(op0Value)
              }
            })

            if (fc201Values.length == 0) {
              return {
                errorMessage: `漏切发票大项！`
              }
            }
            if (fc059Values.length != 0) {
              for (let i = 0; i < fc059Values.length; i++) {
                if (!fc201Values.includes(fc059Values[i])) {
                  return {
                    errorMessage: `发票${fc059Values[i]}属性没有对应发票大项内容，请修改！`
                  }
                }
              }
            }
            return true
          }

          return true
        },


        // 10
        validate10({ mergeFieldsList }) {
          const flatFieldsList = tools.flatArray(mergeFieldsList)
          const fc201Field = flatFieldsList.find(field => field.code === 'fc201')

          if (fc201Field) {
            const [fc059Values, fc201Values] = [[], []]

            flatFieldsList.map(field => {
              if (field.code === 'fc059') {
                if (!fc059Values.includes(field.resultValue)) {
                  fc059Values.push(field.resultValue)
                }
              }

              if (field.code === 'fc201') {
                if (!fc201Values.includes(field.resultValue)) {
                  fc201Values.push(field.resultValue)
                }
              }
            })

            for (let value of fc201Values) {
              if (!fc059Values.includes(value)) {
                return {
                  errorMessage: `发票大项${value}没有匹配的发票，请检查！`
                }
              }
            }

            return true
          }

          return true
        },

        // 11
        validate11({ mergeFieldsList }) {
          const flatFieldsList = tools.flatArray(mergeFieldsList)
          const fc056Field = flatFieldsList.find(field => field.code === 'fc056')

          if (fc056Field) {
            const fc056Values = []

            for (let field of flatFieldsList) {
              if (field.code === 'fc056') {
                fc056Values.push(field.op0Value)
              }
            }

            if (fc056Values.includes('5')) {
              return true
            }
            else {
              return {
                errorMessage: '缺少诊断书，请确认!'
              }
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
      {
        fields: ['fc191', 'fc192', 'fc193', 'fc194', 'fc195', 'fc196', 'fc197', 'fc198', 'fc199', 'fc200'],
        validateDot: function ({ value }) {
          if (!value) return true

          if (!/^[A-Z]/.test(value)) {
            return '疾病诊断第一个必须为大写字母，请检查!'
          }

          if (/\./.test(value)) {
            const index = value.indexOf('.')
            const restValue = value.slice(index + 1)

            if (restValue.length > 1) {
              return '否则提示录入内容只能录入到“.”符号后的一个字符，请检查!'
            }
          }

          return true
        }
      },

      {
        fields: ['fc005', 'fc006', 'fc007', 'fc280'],
        validateDate: function ({ value }) {
          if (!value) return true

          if (/[A, \?]/.test(value)) {
            return true
          }

          if (value.length !== 6 || moment(`20${value}`).format('YYYYMMDD') === 'Invalid date') {
            return '日期格式错误! '
          }

          return true
        }
      },

      // 16
      {
        fields: ['fc053', 'fc182', 'fc183', 'fc184', 'fc185', 'fc186', 'fc187', 'fc188', 'fc189', 'fc190'],
        validate16: function ({ value, items }) {
          // 录入内容为A或包含?时不做校验
          if (value === 'A' || !value || value.includes('?')) {
            return true
          }

          const result = items.find(text => text === value)

          if (result) {
            return true
          }
          else {
            return '疾病诊断录入错误，请根据下拉提示内容选录.'
          }
        }
      },

      // 26
      {
        fields: ['fc055'],
        validate26: function ({ value, items }) {
          // 录入内容为A、为空或包含?时不做校验
          if (value === 'A' || !value || value.includes('?')) {
            return true
          }

          const result = items.find(text => text === value)

          if (result) {
            return true
          }
          else {
            return '手术录入错误，请根据下拉提示内容选录.'
          }
        }
      },

      // 27
      {
        fields: ['fc009', 'fc011', 'fc013', 'fc015', 'fc017', 'fc019', 'fc021', 'fc023', 'fc025', 'fc027', 'fc029', 'fc031', 'fc033', 'fc035', 'fc037', 'fc039', 'fc041', 'fc043', 'fc045', 'fc047'],
        validate27: function ({ field, value, items }) {
          if (value === '?') return true

          const result = items.find(text => text === value)

          if (result) {
            return true
          }
          else {
            return '录入内容与代码表不一致，请选择相近的内容录入，如无相近则录入其他费.'
          }
        }
      },

      // 30
      {
        fields: ['fc152'],
        validate30: function () {
          return true
        }
      },

      // 37
      {
        fields: ['fc172', 'fc173', 'fc174', 'fc175', 'fc176', 'fc177', 'fc178', 'fc179'],
        validate37: function ({ value, field, fieldsIndex, fieldsList }) {
          if (!value) return true

          const mapCodesList = new Map([
            ['fc172', ['fc154', 'fc092']],
            ['fc173', ['fc155', 'fc093']],
            ['fc174', ['fc156', 'fc094']],
            ['fc175', ['fc157', 'fc095']],
            ['fc176', ['fc158', 'fc096']],
            ['fc177', ['fc159', 'fc097']],
            ['fc178', ['fc160', 'fc098']],
            ['fc179', ['fc161', 'fc099']]
          ])

          const codes = mapCodesList.get(field.code)

          if (codes) {
            const fields = fieldsList[fieldsIndex]
            const col1Field = fields.find(_field => _field.code === codes[0])
            const col2Field = fields.find(_field => _field.code === codes[1])

            if (+col1Field.resultValue === 4) {
              const col2FieldValue = col2Field.resultValue || 0

              if (+value > col2FieldValue) {
                return '自付金额不能大于总金额.'
              }
            }
          }

          return true
        }
      },

      // 38
      {
        fields: ['fc084', 'fc085', 'fc086', 'fc087', 'fc088', 'fc089', 'fc090', 'fc091'],
        validate38: function ({ value, field, items }) {
          if (!value || field.force) {
            return true
          }

          const result = items.find(text => text === value)

          if (result) {
            return true
          }
          else {
            return '录入内容不在常量表中，请选择相近内容进行录入，如完全没有相近的则按单录入强过.'
          }
        }
      },

      // 47
      {
        fields: ['fc007'],
        validate47: function ({ fieldsIndex, fieldsList }) {
          const fields = fieldsList[fieldsIndex]

          const fc006Field = tools.find(fields, { code: 'fc006' })
          const fc007Field = tools.find(fields, { code: 'fc007' })

          let fc006Value = fc006Field?.resultValue
          let fc007Value = fc007Field?.resultValue

          if (fc007Value.includes('?') || fc007Value.includes('？')) {
            return true
          }

          fc006Value = +fc006Value
          fc007Value = +fc007Value

          if (fc007Value < fc006Value) {
            return '出院日期不得早于住院日期!'
          }

          return true
        }
      },
    ],

    // 提示文本
    hints: [
      // {
      //   fields: ['fc084', 'fc067'],
      //   hintFc067: function ({ bill }) {
      //     const { billNum } = bill

      //     if (/^64/.test(billNum)) {
      //       return '江苏徐州市的结算单模板，社保自费请按单录入“先行支付”+“自费”!'
      //     }
      //     return true
      //   }
      // },

      {
        fields: ['fc152'],
        hintFc152: function () {
          return '<p style="color: blue;">冻干人用狂犬病疫苗、破伤风疫苗、破伤风类毒素需录入西药费；工本费、病历费、卡费、复印费、陪护费、陪人费此类大项名称按单录入强过；照相费选择放射费、处置费选择治疗费、材料费选择卫材费!</p>'
        }
      },

      // CSB0118RC0320001
      {
        fields: ['fc067'],
        hintFc172: function ({ bill }) {
          if (bill.saleChannel == '1') {
            return '<p style="color: red; margin-top: -3px; margin-bottom: 0px">特约件，乙类自付无需录入， 仅需录入自费金额、超限价自费</p>'
          }
        }
      },
    ],

    // 字段已生成
    updateFields: {
      methods: {
        disableFields: function ({ op, fieldsList, focusFieldsIndex }) {
          if (op === 'op0') {
            return
          }

          const codesList = [
            ['fc057', 'fc058'],
            ['fc144', 'fc145', 'fc146'],
            ['fc202', 'fc203', 'fc204', 'fc211', 'fc212', 'fc209', 'fc210'],
            ['fc103', 'fc150', 'fc151', 'fc170', 'fc171'],
            ['fc225', 'fc226', 'fc227'],
            ['fc228', 'fc229', 'fc230'],
            ['fc231', 'fc232', 'fc233'],
            ['fc104', 'fc105', 'fc106', 'fc107', 'fc108', 'fc109', 'fc110', 'fc111', 'fc112', 'fc113', 'fc114', 'fc115', 'fc116', 'fc117', 'fc118', 'fc119', 'fc120', 'fc121', 'fc122', 'fc123', 'fc124', 'fc125', 'fc126', 'fc127', 'fc128', 'fc129', 'fc130', 'fc131', 'fc132', 'fc133', 'fc134', 'fc135', 'fc136', 'fc137', 'fc138', 'fc139', 'fc140', 'fc141', 'fc142', 'fc143']
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

        // 16
        setConstants16: function ({ flatFieldsList }) {
          const fields = ['fc053', 'fc182', 'fc183', 'fc184', 'fc185', 'fc186', 'fc187', 'fc188', 'fc189', 'fc190']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0118_中意理赔_疾病诊断',
                query: '名称'
              }
            }
          })
        },

        // 26
        setConstants26: function ({ flatFieldsList }) {
          const fields = ['fc055']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0118_中意理赔_手术编码',
                query: '名称'
              }
            }
          })
        },

        // 27
        setConstants27: function ({ flatFieldsList }) {
          const fields = ['fc009', 'fc011', 'fc013', 'fc015', 'fc017', 'fc019', 'fc021', 'fc023', 'fc025', 'fc027', 'fc029', 'fc031', 'fc033', 'fc035', 'fc037', 'fc039', 'fc041', 'fc043', 'fc045', 'fc047']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0118_中意理赔_发票大项类型',
                query: '费用名称'
              }
            }
          })
        },

        // 38
        setConstants28: async function ({ flatFieldsList, codeValues }) {
          const fields = ['fc084', 'fc085', 'fc086', 'fc087', 'fc088', 'fc089', 'fc090', 'fc091']
          const constants = window.constantsDB['B0118'] || {}
          const { fc180, fc181 } = codeValues || {}
          const prefix = 'B0118_中意理赔_省份-'


          let proList = JSON.parse(window.sessionStorage.getItem('proList'))
          let constList = proList && proList.map(el => el.name)

          if (proList) {
            flatFieldsList.map(_field => {
              if (fields.includes(_field.code)) {
                if (constList.includes(`${prefix}${fc180}`)) {
                  _field.table = {
                    name: `${prefix}${fc180}`,
                    query: '中文名称'
                  }
                }
                else if (constList.includes(`${prefix}${fc181}`)) {
                  _field.table = {
                    name: `${prefix}${fc181}`,
                    query: '中文名称'
                  }
                }
                else {
                  _field.table = {
                    name: `${prefix}全国`,
                    query: '中文名称'
                  }
                }
              }
            })
          } else {
            flatFieldsList.map(_field => {
              if (fields.includes(_field.code)) {
                if (constants[`${prefix}${fc180}`]) {
                  _field.table = {
                    name: `${prefix}${fc180}`,
                    query: '中文名称'
                  }
                }
                else if (constants[`${prefix}${fc181}`]) {
                  _field.table = {
                    name: `${prefix}${fc181}`,
                    query: '中文名称'
                  }
                }
                else {
                  _field.table = {
                    name: `${prefix}全国`,
                    query: '中文名称'
                  }
                }
              }
            })
          }
        },

        disableFields29: function ({ op, fieldsList, focusFieldsIndex }) {
          if (op === 'op0') {
            return
          }
          const codesList = [
            [
              'fc235', 'fc255',
              'fc236', 'fc256',
              'fc237', 'fc257',
              'fc238', 'fc258',
              'fc239', 'fc259',
              'fc240', 'fc260',
              'fc241', 'fc261',
              'fc242', 'fc262',
              'fc243', 'fc263',
              'fc244', 'fc264',
              'fc245', 'fc265',
              'fc246', 'fc266',
              'fc247', 'fc267',
              'fc248', 'fc268',
              'fc249', 'fc269',
              'fc250', 'fc270',
              'fc251', 'fc271',
              'fc252', 'fc272',
              'fc253', 'fc273',
              'fc254', 'fc274'
            ]
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

        disableFields30({ fieldsList, focusFieldsIndex, codeValues = {} }) {
          const { fc180 } = codeValues
          const fields = fieldsList[focusFieldsIndex]
          const fc275Field = fields.find(field => field.code === 'fc275')
          console.log(fc180);
          if (fc180 != '徐州' && fc275Field) {
            fc275Field.disabled = true
          }
        },

        // disableFields31({ fieldsList, focusFieldsIndex, codeValues = {} }) {
        //   const { fc277 } = codeValues
        //   console.log(fc277);
        //   if (fc277 == '2') {
        //     const fieldArr = ['fc279', 'fc280', 'fc281', 'fc282']
        //     const fields = fieldsList[focusFieldsIndex]
        //     fields?.map(_field => {
        //       if (fieldArr.includes(_field.code)) {
        //         _field.disabled = true
        //       }
        //     })
        //   }
        // },

        disable32({ fieldsList }) {
          const codes = ['fc282']
          fieldsList?.map(fields => {
            fields?.map(_field => {
              if (codes.includes(_field.code)) {
                _field.disabled = true
              }
            })
          })
        },

        // disableFields33: function ({ op, fieldsList, focusFieldsIndex }) {
        //   if (op != 'op1') return 

        //   const codesList = [
        //     'fc009', 'fc010',
        //     'fc011', 'fc012',
        //     'fc013', 'fc014',
        //     'fc015', 'fc016',
        //     'fc017', 'fc018',
        //     'fc019', 'fc020',
        //     'fc021', 'fc022',
        //     'fc023', 'fc024',
        //     'fc025', 'fc026',
        //     'fc027', 'fc028',
        //     'fc029', 'fc030',
        //     'fc031', 'fc032',
        //     'fc033', 'fc034',
        //     'fc035', 'fc036',
        //     'fc037', 'fc038',
        //     'fc039', 'fc040',
        //     'fc041', 'fc042',
        //     'fc043', 'fc044',
        //     'fc045', 'fc046',
        //     'fc047', 'fc048'
        //   ]

        //   const fields = fieldsList[focusFieldsIndex]

        //   fields?.map(_field => {
        //     if (codesList.includes(_field.code)) {
        //       _field.disabled = true
        //     }
        //   })
        // },

        // 49
        // disableFields49({ op, fieldsList, focusFieldsIndex, codeValues }) {
        //   if (op == 'op0') return

        //   const { fc056 } = codeValues || {}
        //   if( fc056 < 2 ) return

        //   const fields = fieldsList[focusFieldsIndex]

        //   fields?.map(_field => {
        //     if (_field.code == 'fc062') {
        //       _field.disabled = true
        //     }
        //   })
        // },

        // 53
        // setConstants53: function ({ fieldsList, focusFieldsIndex, flatFieldsList }) {
        //   const fieldsArr = ['fc009', 'fc011', 'fc013', 'fc015', 'fc017', 'fc019', 'fc021', 'fc023', 'fc025', 'fc027', 'fc029', 'fc031', 'fc033', 'fc035', 'fc037', 'fc039', 'fc041', 'fc043', 'fc045', 'fc047']
        //   const fields = fieldsList[focusFieldsIndex]
        //   const fc003field = fields.find(field => field.code === 'fc003')

        //   if (fc003field && fc003field.resultValue == '1') {
        //     flatFieldsList.map(_field => {
        //       if (fieldsArr.includes(_field.code)) {
        //         _field.table = {
        //           name: 'B0118_中意理赔_门诊费用类型',
        //           query: '费用名称'
        //         }               
        //       }
        //     })
        //   } else if (fc003field && fc003field.resultValue == '2') {
        //     flatFieldsList.map(_field => {
        //       if (fieldsArr.includes(_field.code)) {          
        //         _field.table = {
        //           name: 'B0118_中意理赔_住院费用类型',
        //           query: '费用名称'
        //         }
        //       }
        //     })
        //   }
        // },

        // setConstants056: function ({ flatFieldsList, codeValues }) {

        //   const { fc056 } = codeValues || {}

        //   flatFieldsList.map(_field => {
        //     if (_field.code == 'fc084' && fc056 == '3') {
        //       _field.hint = `<p style="color: blue; fontSize: 14px; margin-top: -3px; margin-bottom: 0px">有框图， 注意框外内容不能录入</p>`
        //     }
        //   })
        // },
      }
    },

    // 回车
    enter: {
      methods: {

        // 17
        // CSB0118RC0186001
        validate17({ op, field, fieldsList, focusFieldsIndex }) {
          const fields = fieldsList[focusFieldsIndex]

          if (field.code === 'fc053') {
            const ifADisableCodes = ['fc053', 'fc182', 'fc183', 'fc184', 'fc185', 'fc186', 'fc187', 'fc188', 'fc189', 'fc190']
            const ifNotADisableCodes = ['fc191', 'fc192', 'fc193', 'fc194', 'fc195', 'fc196', 'fc197', 'fc198', 'fc199', 'fc200']

            if (field.resultValue == 'A') {
              fields.map(_field => {
                if (ifADisableCodes.includes(_field.code)) {
                  _field.disabled = true
                  _field[`${op}Value`] = ''
                  _field.resultValue = ''
                }

                if (ifNotADisableCodes.includes(_field.code)) {
                  _field.disabled = false
                }
              })
            } else if (field.resultValue != 'A') {
              console.log(123456);
              fields.map(_field => {
                if (ifADisableCodes.includes(_field.code)) {
                  _field.disabled = false
                }

                if (ifNotADisableCodes.includes(_field.code)) {
                  _field.disabled = true
                  _field[`${op}Value`] = ''
                  _field.resultValue = ''
                }
              })
            }
          }
        },

        // 20
        validate20({ op, fieldsList, focusFieldsIndex }) {
          const fields = fieldsList[focusFieldsIndex]

          const fc003Field = fields.find(field => field.code === 'fc003')
          const fc005Field = fields.find(field => field.code === 'fc005')
          const fc006Field = fields.find(field => field.code === 'fc006')
          const fc007Field = fields.find(field => field.code === 'fc007')

          if (fc003Field?.[`${op}Value`] == 1) {
            fc005Field.disabled = false
            fc006Field.disabled = true
            fc007Field.disabled = true

            const fc005Value = fc005Field[`${op}Value`]

            if (fc005Value) {
              fc006Field[`${op}Value`] = fc005Value
              fc006Field.resultValue = fc005Value

              fc007Field[`${op}Value`] = fc005Value
              fc007Field.resultValue = fc005Value
            }
            else {
              fc006Field[`${op}Value`] = ''
              fc006Field.resultValue = ''

              fc007Field[`${op}Value`] = ''
              fc007Field.resultValue = ''
            }
          }
          else if (fc003Field?.resultValue == 2) {
            fc005Field.disabled = true
            fc006Field.disabled = false
            fc007Field.disabled = false

            fc005Field[`${op}Value`] = ''
            fc005Field.resultValue = ''
          }
        },

        // 22
        validate22({ field, fieldsList, focusFieldsIndex }) {
          if (field.code === 'fc004') {
            const fields = fieldsList[focusFieldsIndex]
            const codes = ['fc063', 'fc064', 'fc065', 'fc066', 'fc067', 'fc068', 'fc069', 'fc070', 'fc071', 'fc100', 'fc072', 'fc074', 'fc075', 'fc076', 'fc077', 'fc050']

            fields.map(_field => {
              if (codes.includes(_field.code)) {
                _field.disabled = false
              }
            })

            if (field.resultValue === '3' || field.resultValue === '4') {
              fields.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = true
                }
              })
            }
          }
        },

        // 24 47
        validate24({ field, fieldsList, focusFieldsIndex }) {
          if (field.code === 'fc101') {
            const fields = fieldsList[focusFieldsIndex]
            const fc079Field = fields.find(field => field.code === 'fc079')
            const fc051Field = fields.find(field => field.code === 'fc051')
            const fc102Field = fields.find(field => field.code === 'fc102')
            const fc052Field = fields.find(field => field.code === 'fc052')
            const fc067Field = fields.find(field => field.code === 'fc067')
            const fc275Field = fields.find(field => field.code === 'fc275')

            fc079Field.disabled = false
            fc051Field.disabled = false
            fc102Field.disabled = false
            fc052Field.disabled = false
            fc067Field.disabled = false
            // fc275Field.disabled = false

            if (field.resultValue === '1' || field.resultValue === '2') {
              fc051Field.disabled = true
              fc102Field.disabled = true
              fc079Field.disabled = true
            }
            else {
              fc052Field.disabled = true
              fc067Field.disabled = true
              fc275Field.disabled = true
            }
          }
        },

        // 36
        validate36({ field, fieldsList, focusFieldsIndex }) {
          const mapCodesList = new Map([
            ['fc154', ['fc162', 'fc172']],
            ['fc155', ['fc163', 'fc173']],
            ['fc156', ['fc164', 'fc174']],
            ['fc157', ['fc165', 'fc175']],
            ['fc158', ['fc166', 'fc176']],
            ['fc159', ['fc167', 'fc177']],
            ['fc160', ['fc168', 'fc178']],
            ['fc161', ['fc169', 'fc179']]
          ])

          const codes = mapCodesList.get(field.code)

          if (codes) {
            const fields = fieldsList[focusFieldsIndex]
            const rightFirField = fields.find(field => field.code === codes[1])
            const rightSecField = fields.find(field => field.code === codes[0])

            rightFirField.disabled = false
            rightSecField.disabled = false


            if (field.resultValue === '1') {
              rightFirField.disabled = true
            }
            else if (field.resultValue === '4') {
              rightSecField.disabled = true
            }
            else {
              rightFirField.disabled = true
              rightSecField.disabled = true
            }
          }
        },

        // 45
        validate45({ op, field, fieldsList, focusFieldsIndex }) {
          if (field.code === 'fc205') {
            const fields = fieldsList[focusFieldsIndex]
            const fc206Field = fields.find(field => field.code === 'fc206')
            const fc207Field = fields.find(field => field.code === 'fc207')
            const fc208Field = fields.find(field => field.code === 'fc208')

            fc206Field.disabled = false
            fc207Field.disabled = false
            fc208Field.disabled = false

            if (field.resultValue === 'A') {
              fc206Field[`${op}Value`] = ''
              fc206Field.restValue = ''
              fc206Field.disabled = true
              fc207Field.disabled = true
              fc208Field.disabled = true
            }
            else if (field.resultValue == 1 || field.resultValue == 2) {
              fc207Field.disabled = true
              fc208Field.disabled = true
            }
            else {
              fc206Field.disabled = true
            }
          }
        },

        validate53({ op, field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc004') return

          const fields = fieldsList[focusFieldsIndex]
          const fc003field = fields.find(field => field.code === 'fc003')
          const fc205field = fields.find(field => field.code === 'fc205')

          if (fc003field.resultValue == '1' && (field.resultValue == '2' || field.resultValue == '4')) {
            field.hint = `<p style="color: blue; fontSize: 14px; margin-top: -3px; margin-bottom: 0px">账单类型录入错误， 请修改</p>`
          } else if (fc003field.resultValue == '2' && (field.resultValue == '1' || field.resultValue == '3')) {
            field.hint = `<p style="color: blue; fontSize: 14px; margin-top: -3px; margin-bottom: 0px">账单类型录入错误， 请修改</p>`
          } else {
            field.hint = ''
          }

          if (field.resultValue == '3' || field.resultValue == '4') {
            fc205field[`${op}Value`] = 'A'
            fc205field.resultValue = 'A'
          } else {
            fc205field[`${op}Value`] = ''
            fc205field.resultValue = ''
          }
        },

        // 47
        validate47({ op, field }) {
          const fields = [
            'fc010', 'fc012', 'fc014', 'fc016', 'fc018',
            'fc020', 'fc022', 'fc024', 'fc026', 'fc028',
            'fc030', 'fc032', 'fc034', 'fc036', 'fc038',
            'fc040', 'fc042', 'fc044', 'fc046', 'fc048'
          ]
          if (fields.includes(field.code) && field.resultValue.includes('.')) {
            let flag = field.resultValue.indexOf('.')
            let str = field.resultValue.slice(flag + 1).split('')

            if (str[0] == 0 && [...new Set(str)].length == 1) {
              field[`${op}Value`] = field.resultValue.slice(0, flag)
              field.resultValue = field.resultValue.slice(0, flag)
            }
          }

          if (fields.includes(field.code) && field.resultValue.includes('.') && field.resultValue.charAt(field.resultValue.length - 1) == '0') {
            field[`${op}Value`] = field.resultValue.slice(0, -1)
            field.resultValue = field.resultValue.slice(0, -1)
          }
        },

        validate339({ op, field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc275') return

          const fields = fieldsList[focusFieldsIndex]
          const fc067Field = fields.find(field => field.code === 'fc067')

          if (field.resultValue == '1') {
            fc067Field[`${op}Value`] = '0'
            fc067Field.resultValue = '0'
            fc067Field.disabled = true
          } else {
            fc067Field[`${op}Value`] = ''
            fc067Field.resultValue = ''
            fc067Field.disabled = false
          }
        },
      },
    },

    // 临时保存
    sessionSave: {
      methods: {
        // 18
        // CSB0118RC0187001
        disable18({ op, fieldsList, focusFieldsIndex }) {
          const codesList = [
            ['fc182', 'fc183', 'fc184', 'fc185', 'fc186', 'fc187', 'fc188', 'fc189', 'fc190'],
            ['fc191', 'fc192', 'fc193', 'fc194', 'fc195', 'fc196', 'fc197', 'fc198', 'fc199', 'fc200']
          ]

          const fields = fieldsList[focusFieldsIndex]
          const focusField = fields.find(field => field.autofocus)

          let codes = []

          for (let index in codesList) {
            if (codesList[index].includes(focusField.code)) {
              codes = codesList[index]
              break
            }
          }

          const codeIndex = codes.indexOf(focusField.code)

          if (codeIndex > -1) {
            const restCodes = codes.slice(codeIndex + 1)
            const restFields = fields.slice(focusField.fieldIndex + 1)

            restFields?.map(restField => {
              if (restCodes.includes(restField.code)) {
                restField.disabled = true
                restField[`${op}Value`] = ''
                restField.resultValue = ''
              }
            })
          }
        },

        // 31
        disable31({ fieldsList, focusFieldsIndex }) {
          const codesList = [
            ['fc009', 'fc010'],
            ['fc011', 'fc012'],
            ['fc013', 'fc014'],
            ['fc015', 'fc016'],
            ['fc017', 'fc018'],
            ['fc019', 'fc020'],
            ['fc021', 'fc022'],
            ['fc023', 'fc024'],
            ['fc025', 'fc026'],
            ['fc027', 'fc028'],
            ['fc029', 'fc030'],
            ['fc031', 'fc032'],
            ['fc033', 'fc034'],
            ['fc035', 'fc036'],
            ['fc037', 'fc038'],
            ['fc039', 'fc040'],
            ['fc041', 'fc042'],
            ['fc043', 'fc044'],
            ['fc045', 'fc046'],
            ['fc047', 'fc048']
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
        },

        // 41
        disable41({ fieldsList, focusFieldsIndex }) {
          const codesList = [
            ['fc084', 'fc154', 'fc092', 'fc162', 'fc172'],
            ['fc085', 'fc155', 'fc093', 'fc163', 'fc173'],
            ['fc086', 'fc156', 'fc094', 'fc164', 'fc174'],
            ['fc087', 'fc157', 'fc095', 'fc165', 'fc175'],
            ['fc088', 'fc158', 'fc096', 'fc166', 'fc176'],
            ['fc089', 'fc159', 'fc097', 'fc167', 'fc177'],
            ['fc090', 'fc160', 'fc098', 'fc168', 'fc178'],
            ['fc091', 'fc161', 'fc099', 'fc169', 'fc179']
          ]
          const fields = fieldsList[focusFieldsIndex]
          const focusField = fields.find(field => field.autofocus)
          let sliceIndex = -1

          for (let codesIndex in codesList) {
            if (codesList[codesIndex].includes(focusField.code)) {
              sliceIndex = +codesIndex + 1
              break
            }
          }

          if (sliceIndex > -1) {
            const restCodesList = codesList.slice(sliceIndex)
            const restCodes = []

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

    // 提交前
    beforeSubmit: {
      methods: {
        // 32
        validate32({ block, op, fieldsList }) {
          const fields = fieldsList[0]
          const codes = ['fc010', 'fc012', 'fc014', 'fc016', 'fc018', 'fc020', 'fc022', 'fc024', 'fc026', 'fc028', 'fc030', 'fc032', 'fc034', 'fc036', 'fc038', 'fc040', 'fc042', 'fc044', 'fc046', 'fc048']
          const fc008Field = fields.find(field => field.code === 'fc008')
          let total = 0

          if (block.code !== 'bc007' || op === 'opq') {
            return true
          }
          if (op === 'opq') {
            return true
          }

          if (!fc008Field) {
            return true
          }

          for (let field of fields) {
            let value = field[`${op}Value`]

            if (/\?/.test(value)) {
              return true
            }

            if (codes.includes(field.code)) {
              if (!value) {
                value = 0
              }

              total += Number(value)
            }
          }

          const diff = +fc008Field[`${op}Value`] - total.toFixed(2)

          if (diff === 0) {
            return true
          } else {
            return {
              popup: 'confirm',
              errorMessage: `明细金额与总金额不一致，差额为${diff}，请确认并修改!`
            }
          }
        },

        // 39
        validate39({ block, fieldsList }) {
          if (block.code !== 'bc002') {
            return true
          }

          for (let fields of fieldsList) {
            for (let field of fields) {
              if (field.resultValue) {
                return true
              }
            }
          }

          return {
            errorMessage: '清单不能空白提交，请检查，如清单内容无法录入则录入一组数据后按F8提交.'
          }
        },

        // 40
        validate40({ fieldsList }) {
          const error = {
            errorMessage: '清单内容录入遗漏，请检查!'
          }

          for (let fields of fieldsList) {
            const codesList = [
              { fc084: '', fc154: '', fc092: '', fc162: '', fc172: '' },
              { fc085: '', fc155: '', fc093: '', fc163: '', fc173: '' },
              { fc086: '', fc156: '', fc094: '', fc164: '', fc174: '' },
              { fc087: '', fc157: '', fc095: '', fc165: '', fc175: '' },
              { fc087: '', fc158: '', fc096: '', fc166: '', fc176: '' },
              { fc088: '', fc159: '', fc097: '', fc167: '', fc177: '' },
              { fc090: '', fc160: '', fc098: '', fc168: '', fc178: '' },
              { fc091: '', fc161: '', fc099: '', fc169: '', fc179: '' }
            ]

            for (let field of fields) {
              for (let codeItem of codesList) {
                if (codeItem.hasOwnProperty(field.code)) {
                  codeItem[field.code] = field.resultValue
                }
              }
            }

            for (let codeItem of codesList) {
              const keys = Object.keys(codeItem)
              const firKey = keys[0]
              const secKey = keys[1]
              const thiKey = keys[2]
              const fouKey = keys[3]
              const fifKey = keys[4]

              if (codeItem[secKey] == 1) {
                if ((!codeItem[firKey] || !codeItem[thiKey] || !codeItem[fouKey]) || codeItem[fifKey]) {
                  return error
                }
              }

              if (codeItem[secKey] == 2 || codeItem[secKey] == 3) {
                if ((!codeItem[firKey] || !codeItem[thiKey]) || (codeItem[fouKey] || codeItem[fifKey])) {
                  return error
                }
              }

              if (codeItem[secKey] == 4) {
                if ((!codeItem[firKey] || !codeItem[thiKey] || !codeItem[fifKey]) || codeItem[fouKey]) {
                  return error
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

export default B0118