import { tools, sessionStorage } from 'vue-rocket'
import BigNumber from 'bignumber.js'
import moment from 'moment'
import { MessageBox, Notification } from 'element-ui';

const B0122 = {
  op0: {
    // 记录最后一次存储的合法field
    memoFields: ['fc205', 'fc204'],

    // 记录相同 code 的 field 的值
    memoFieldValues: ['fc153', 'fc477', 'fc156', 'fc157', 'fc158', 'fc159'],

    // fields 的值从 targets 里的值选择
    dropdownFields: [
      {
        targets: ['fc156', 'fc157'],
        fields: ['fc158', 'fc159', 'fc160']
      },
    ],

    // 校验规则
    rules: [
      // 2
      // CSB0122RC0003000
      {
        fields: ['fc156', 'fc157'],
        validate02: function ({ field, fieldsObject, thumbIndex, value }) {
          let fc156Values = []
          let fc157Values = []
          let allValues = []

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage

            if (sessionStorage || thumbIndex === +key) {
              const _fieldsList = fieldsObject[key].fieldsList

              for (let _fields of _fieldsList) {
                for (let _field of _fields) {
                  if (_field.code === 'fc156' && _field.uniqueId !== field.uniqueId) {
                    fc156Values.push(_field.resultValue)
                  }
                }
              }

              for (let _fields of _fieldsList) {
                for (let _field of _fields) {
                  if (_field.code === 'fc157' && _field.uniqueId !== field.uniqueId) {
                    fc157Values.push(_field.resultValue)
                  }
                }
              }
            }
          }

          allValues = [...fc156Values, ...fc157Values]

          if (fc156Values.includes(value)) {
            return '发票属性不能重复!'
          }

          if (fc157Values.includes(value)) {
            return '发票属性不能重复!'
          }

          if (allValues.includes(value)) {
            return '发票属性不能重复!'
          }

          return true
        }
      },

      // 3
      // CSB0122RC0004000
      {
        fields: ['fc158', 'fc159', 'fc160'],
        validate03: function ({ includes, value }) {
          if (includes) {
            const result = includes.find(text => text === value)

            if (!result) {
              return '没有此发票，请核实!'
            }
          }

          return true
        }
      },

      // 9
      // CSB0122RC0010000
      // {
      //   fields: ['fc037', 'fc040', 'fc041'],
      //   validate09: function ({ field, value, items }) {
      //     field.allowForce = false
      //     if (value === '?') {
      //       return true
      //     }

      //     const result = items.find((text) => text === value)

      //     if (!result) {
      //       return '录入内容不在数据库中，不可强过!'
      //     }

      //     return true
      //   }
      // },

      // 10
      // CSB0122RC0011000
      {
        fields: ['fc161'],
        validate10: function ({ field, value, items }) {
          field.allowForce = false
          if (value === '?') {
            return true
          }

          const result = items.find((text) => text === value)

          if (!result) {
            return '录入内容不在数据库中，不可强过!'
          }

          return true
        }
      },

      // 38 || 39 || 40
      // CSB0122RC0039000
      // CSB0122RC0040000
      // CSB0122RC0041000
      {
        fields: ['fc153'],
        validate38: function ({ bill, value }) {
          if (bill?.saleChannel === '微理赔') {
            if (value == '2' || value == '3' || value == '4') {
              return '该案件申请书切1'
            }
          } else if (bill?.saleChannel === '移动理赔') {
            if (value == '1' || value == '3' || value == '4') {
              return '该案件申请书切2'
            }
          } else if (bill?.saleChannel === '纸质') {
            if (value == '1' || value == '2') {
              return '该案件申请书切3或4， 请注意影像右上角标记'
            }
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
        // 37
        // CSB0122RC0038000
        validate37: function ({ bill }) {
          if (bill?.saleChannel) {
            const value = bill?.saleChannel
            switch (value) {
              case '微理赔':

                if (sessionStorage.get('isApp')?.isApp === 'true') {
                  // return MessageBox.alert('申请书切1-微理赔', '请注意', {
                  //   type: 'warning',
                  //   confirmButtonText: '确定',
                  //   showClose: false,
                  // })
                  return Notification({
                    type: 'warning',
                    title: '提醒(5s后自动关闭)',
                    message: '申请书切1-微理赔',
                    duration: 5000,
                    position: 'top-left'
                  })
                } else {
                  return alert('申请书切1-微理赔')
                }

              case '移动理赔':

                if (sessionStorage.get('isApp')?.isApp === 'true') {
                  // return MessageBox.alert('申请书切2-移动理赔', '请注意', {
                  //   type: 'warning',
                  //   confirmButtonText: '确定',
                  //   showClose: false,
                  // })
                  return Notification({
                    type: 'warning',
                    title: '提醒(5s后自动关闭)',
                    message: '申请书切2-移动理赔',
                    duration: 5000,
                    position: 'top-left'
                  })
                } else {
                  return alert('申请书切2-移动理赔')
                }

              case '纸质':

                if (sessionStorage.get('isApp')?.isApp === 'true') {
                  // return MessageBox.alert('申请书请仔细确认是小额理赔还是理赔', '请注意', {
                  //   type: 'warning',
                  //   confirmButtonText: '确定',
                  //   showClose: false,
                  // })
                  return Notification({
                    type: 'warning',
                    title: '提醒(5s后自动关闭)',
                    message: '申请书请仔细确认是小额理赔还是理赔',
                    duration: 5000,
                    position: 'top-left'
                  })
                } else {
                  return alert('申请书请仔细确认是小额理赔还是理赔')
                }

            }
          }

          return true
        }
      }
    },

    // 字段已生成
    updateFields: {
      methods: {
        // 默认屏蔽
        // 8
        // CSB0122RC0009000
        disable08({ fieldsList, focusFieldsIndex }) {
          const codesList = ['fc262', 'fc038', 'fc039', 'fc485']

          const fields = fieldsList[focusFieldsIndex]

          fields?.map(_field => {
            if (codesList.includes(_field.code)) {
              _field.disabled = true
            }
          })
        },

        // 9
        // CSB0122RC0010000
        setConstants09: function ({ flatFieldsList }) {
          const fields = ['fc037', 'fc040', 'fc041', 'fc478']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0122_信诚理赔_中保信医院代码表',
                query: '医院名称'
              }

              _field.tables = {
                name: 'B0122_信诚理赔_内部医院代码表',
                query: '医院名称'
              }
            }
          })
        },

        // 10
        // CSB0122RC0011000
        set10Items({ bill, fieldsList, focusFieldsIndex }) {
          const dropArr = ['fc161']
          // // 数据库
          // const db = window['constantsDB']['B0122']
          // if (!db) return
          // const collections = db['B0122_信诚理赔_机构号代码表']
          // // 案件的机构号
          // const agency = bill.agency
          // // 机构代码: agencyCode 地址: address
          // const [address, agencyCode] = [[], []]
          // // 本级机构代码对应的地址
          // let addressValue = ''
          // // 指定的地址目录序号
          // let index = ''
          // if (agency) {
          //   for (let dessert of collections.desserts) {
          //     address.push(dessert[0])
          //     agencyCode.push(dessert[1])
          //   }

          //   index = agencyCode.indexOf(agency)
          //   addressValue = address[index]
          // }

          const fields = fieldsList[focusFieldsIndex]

          if (bill.agency == '北京') {
            fields.map(_field => {
              if (dropArr.includes(_field.code)) {
                _field.table = {
                  name: `B0122_信诚理赔_北京发票大项明细表`,
                  query: '名称'
                }
              }
            })
          } else {
            fields.map(_field => {
              if (dropArr.includes(_field.code)) {
                _field.table = {
                  name: `B0122_信诚理赔_发票大项明细表`,
                  query: '名称'
                }
              }
            })
          }
        },

        // CSB0122RC0051000
        disable09({ fieldsList, focusFieldsIndex }) {
          const codesList = ['fc496']

          const fields = fieldsList[focusFieldsIndex]

          fields?.map(_field => {
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
        // CSB0122RC0170000
        disable170({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc071') return

          const fields = fieldsList[focusFieldsIndex]
          const fc072Field = fields.find(field => field.code === 'fc072')

          if (field.resultValue != '1') {
            fc072Field.disabled = true
          } else {
            fc072Field.disabled = false
          }
        },

        // CSB0122RC0171000
        disable171({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc072') return

          const fields = fieldsList[focusFieldsIndex]
          const fc037Field = fields.find(field => field.code === 'fc037')

          if (field.resultValue == '2') {
            fc037Field.hint = `<p  style="color: red; fontSize: 14px; margin-top: -3px; margin-bottom: 0px">就诊医院需要查找病历资料进行补充</p>`
          } else {
            fc037Field.hint = ''
          }
        },

        disabled190({ field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc484') return
          const codeArr = ['fc055', 'fc056']

          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue != '1') {
            for (let _field of fields) {
              if (codeArr.includes(_field.code)) {
                _field.disabled = true
              }
            }
          } else {
            for (let _field of fields) {
              if (codeArr.includes(_field.code)) {
                _field.disabled = false
              }
            }
          }
        },

        disabled191({ op, field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc135') return

          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue == '2') {
            for (let _field of fields) {
              if (_field.code == 'fc486') {
                _field[`${op}Value`] = '';
                _field.resultValue = '';
                _field.disabled = true
              }
            }
          } else {
            for (let _field of fields) {
              if (_field.code == 'fc486') {
                _field.disabled = false
              }
            }
          }
        },
      }
    },

    // F8(提交前校验)
    beforeSubmit: {
      methods: {
        // 4
        // CSB0122RC0005000
        // validate04({ fieldsObject }) {
        //   const [fc158Values, fc157Values, fc156Values] = [[], [], []]
        //   for (let key in fieldsObject) {
        //     const sessionStorage = fieldsObject[key].sessionStorage
        //     const fieldsList = fieldsObject[key].fieldsList
        //     if (sessionStorage) {
        //       for (let fields of fieldsList) {
        //         for (let field of fields) {
        //           const { code, resultValue } = field
        //           if (code === 'fc158') {
        //             resultValue && fc158Values.push(resultValue)
        //           }
        //           if (code === 'fc156') {
        //             resultValue && fc156Values.push(resultValue)
        //           }
        //           if (code === 'fc157') {
        //             resultValue && fc157Values.push(resultValue)
        //           }
        //         }
        //       }
        //     }
        //   }
        //   const mergeValues = [...fc157Values, ...fc156Values]
        //   for (let value of mergeValues) {
        //     if (!fc158Values.includes(value)) {
        //       return {
        //         errorMessage: `发票属性为${value}的发票漏切清单!`
        //       }
        //     }
        //   }

        //   return true
        // },

        // 5
        // CSB0122RC0006000
        validate05({ fieldsObject }) {
          let fc153Values = []

          for (let key in fieldsObject) {
            const fieldsList = fieldsObject[key].fieldsList

            for (let fields of fieldsList) {
              const fc153Field = tools.find(fields, { code: 'fc153' })
              if (fc153Field) {
                fc153Values.push(fc153Field?.resultValue)
              }
            }
          }

          if (!fc153Values.includes('10')) {
            return {
              errorMessage: `漏切诊断书!`
            }
          }

          return true
        },

        // 6
        // CSB0122RC0007000
        validate06({ fieldsObject }) {
          let fc153Values = []

          for (let key in fieldsObject) {
            const fieldsList = fieldsObject[key].fieldsList

            for (let fields of fieldsList) {
              const fc153Field = tools.find(fields, { code: 'fc153' })
              if (fc153Field) {
                fc153Values.push(fc153Field?.resultValue)
              }
            }
          }

          if (fc153Values.includes('24')) {
            if (!fc153Values.includes('26')) {
              return {
                errorMessage: `身故案件必须切受益人及其关系声明文件!`
              }
            }
          }
          return true
        },

        // 7
        // CSB0122RC0008000
        validate07({ sameFieldValue }) {
          let fc153Values = sameFieldValue.fc153?.values
          if (!fc153Values) return
          let codesMap = new Map([
            ['1', '重复录入: 申请书， 请核查'],
            ['2', '重复录入: 申请书， 请核查'],
            ['3', '重复录入: 申请书， 请核查'],
            ['4', '重复录入: 申请书， 请核查'],
            // ['10', '重复录入: 诊断书， 请核查'],
            // ['11', '重复录入: 诊断书， 请核查'],
            // ['12', '重复录入: 12-手术， 请核查'],
            ['13', '重复录入: 13-受款人银行卡， 请核查'],
            ['14', '重复录入: 14-事故人身份证， 请核查'],
            ['15', '重复录入: 15-申请人身份证， 请核查'],
            ['16', '重复录入: 16-报警回执、不予立案通知， 请核查'],
            ['17', '重复录入: 17-立案决定书， 请核查'],
            ['18', '重复录入: 18-立案报告书， 请核查'],
            ['19', '重复录入: 19-交通事故责任认定书， 请核查'],
            ['20', '重复录入: 20-驾驶证， 请核查'],
            ['21', '重复录入: 21-行驶证， 请核查'],
            ['22', '重复录入: 22-鉴定报告， 请核查'],
            ['23', '重复录入: 23-病理报告， 请核查'],
            ['24', '重复录入: 24-身故（殡葬证明、户籍注销证明、死亡证明）， 请核查'],
            ['25', '重复录入: 25-投保人与被投保人关系证明， 请核查'],
            ['26', '重复录入: 26-受益人与被保人关系证明， 请核查'],

            ['11', '重复录入: 11-受款人身份证'],
            ['31', '重复录入: 31-受益人身份证1'],
            ['32', '重复录入: 32-受益人银行卡1'],
            ['33', '重复录入: 33-受益人身份证2'],
            ['34', '重复录入: 34-受益人银行卡2'],
            ['35', '重复录入: 35-受益人身份证3'],
            ['36', '重复录入: 36-受益人银行卡3'],
            ['37', '重复录入: 37-受益人身份证4'],
            ['38', '重复录入: 38-受益人银行卡4'],
            ['39', '重复录入: 39-受益人身份证5'],
            ['40', '重复录入: 40-受益人银行卡5'],
          ])

          let duplicates = [];

          fc153Values = fc153Values.map(el => {
            if (el == '1' || el == '2' || el == '3' || el == '4') {
              return el = '1'
            } else {
              return el
            }
            // if (el == '10' || el == '11') {
            //   return '10'
            // }
          })

          for (let i = 0; i < fc153Values.length; i++) {
            for (let j = i + 1; j < fc153Values.length; j++) {
              if (fc153Values[i] === fc153Values[j] && !duplicates.includes(fc153Values[i])) {
                duplicates.push(fc153Values[i]);
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

        // CSB0122RC0173000
        validate173({ fieldsObject }) {

          for (let key in fieldsObject) {
            const fieldsList = fieldsObject[key].fieldsList
            for (let fields of fieldsList) {
              for (let field of fields) {
                const fc156 = fields.find(el => el.code == 'fc156')
                if (field.code == 'fc037' && field.resultValue == '') {
                  return {
                    errorMessage: `发票属性${fc156.resultValue}就诊医院为空， 不能提交`
                  }
                }
              }
            }
          }

          return true
        },


        validate192({ sameFieldValue }) {
          let fc153Values = sameFieldValue.fc153?.values
          if (!fc153Values) return
          let codesMap = new Map([
            ['1', '漏切: 申请书， 请核查'],
            ['2', '漏切: 申请书， 请核查'],
            ['3', '漏切: 申请书， 请核查'],
            ['4', '漏切: 申请书， 请核查'],
            ['11', '漏切: 11-受款人身份证， 请核查'],
            ['13', '漏切: 13-受款人银行卡， 请核查'],
            ['14', '漏切: 14-事故人身份证， 请核查'],
            ['25', '漏切: 25-意外'],
          ])

          let codeArr = ['1', '11', '13', '14', '25']
          fc153Values = fc153Values.map(el => {
            if (el == '1' || el == '2' || el == '3' || el == '4') {
              return el = '1'
            } else {
              return el
            }
          })

          for (let el of codeArr) {
            if (!fc153Values.includes(el)) {
              return {
                errorMessage: codesMap.get(el)
              }
            }
          }
          return true
        },

        validate193({ sameFieldValue }) {
          let fc153Values = sameFieldValue.fc153?.values
          if (!fc153Values) return
          if (fc153Values.includes('26')) {
            if (!fc153Values.includes('31')) {
              return {
                errorMessage: '身故案件必须切受益人身份证及银行卡'
              }
            }
          }

          if (fc153Values.includes('31') || fc153Values.includes('32')) {
            if (!(fc153Values.includes('31') && fc153Values.includes('32'))) {
              return {
                errorMessage: '受益人身份证1及银行卡1需同时切'
              }
            }
          }


          if (fc153Values.includes('33') || fc153Values.includes('34')) {
            if (!(fc153Values.includes('33') && fc153Values.includes('34'))) {
              return {
                errorMessage: '受益人身份证2及银行卡2需同时切'
              }
            }
          }

          if (fc153Values.includes('35') || fc153Values.includes('36')) {
            if (!(fc153Values.includes('35') && fc153Values.includes('36'))) {
              return {
                errorMessage: '受益人身份证3及银行卡3需同时切'
              }
            }
          }


          if (fc153Values.includes('37') || fc153Values.includes('38')) {
            if (!(fc153Values.includes('37') && fc153Values.includes('38'))) {
              return {
                errorMessage: '受益人身份证4及银行卡4需同时切'
              }
            }
          }

          if (fc153Values.includes('39') || fc153Values.includes('40')) {
            if (!(fc153Values.includes('39') && fc153Values.includes('40'))) {
              return {
                errorMessage: '受益人身份证5及银行卡5需同时切'
              }
            }
          }
        },

        validate194({ sameFieldValue, bill }) {
          let fc153Values = sameFieldValue.fc153?.values
          if (!fc153Values) return
          if (bill.insuranceType == '身故') {
            if (!fc153Values.includes('24')) {
              return {
                errorMessage: `该案件为身故案件， 需切24-死亡证明`
              }
            }
          }

          if (bill.insuranceType == '全残') {
            if (!fc153Values.includes('22')) {
              return {
                errorMessage: `该案件为全残案件， 需切22-鉴定报告`
              }
            }
          }

          if (bill.insuranceType == '部分残疾') {
            if (!fc153Values.includes('22')) {
              return {
                errorMessage: `该案件为部分残疾案件， 需切22-鉴定报告`
              }
            }
          }

          if (bill.insuranceType == '重大疾病') {
            if (!fc153Values.includes('23')) {
              return {
                errorMessage: `该案件为重疾案件， 需切23-病理报告`
              }
            }
          }

          return true
        },

        // validate195({ sameFieldValue }) {
        //   let fc153Values = sameFieldValue.fc153?.values
        //   let fc477Values = sameFieldValue.fc477?.values
        //   if (!fc153Values || !fc477Values) return
        //   if (fc477Values.includes('1')) {
        //     if (!fc153Values.includes('25')) {
        //       return {
        //         errorMessage: `该案件为意外案件， 必需切25-意外`
        //       }
        //     }
        //   }
        // },

        validate196({ sameFieldValue }) {
          let fc156Values = sameFieldValue.fc156?.values
          let fc158Values = sameFieldValue.fc158?.values || []
          let fc159Values = sameFieldValue.fc159?.values || []

          let arr = [...fc158Values, ...fc159Values]
          if (fc156Values) {
            for (let el of fc156Values) {
              if (!arr.includes(el)) {
                return {
                  errorMessage: `发票属性${el}缺少对应清单或报销单`
                }
              }
            }
          }
        },

        validate197({ sameFieldValue }) {
          let fc157Values = sameFieldValue.fc157?.values
          let fc158Values = sameFieldValue.fc158?.values || []
          let fc159Values = sameFieldValue.fc159?.values || []

          let arr = [...fc158Values, ...fc159Values]
          if (fc157Values) {
            for (let el of fc157Values) {
              if (!arr.includes(el)) {
                return {
                  errorMessage: `非医疗发票属性${el}缺少对应清单或报销单`
                }
              }
            }
          }
        }
      }
    }
  },

  op1op2opq: {
    // 校验规则
    rules: [
      // 11 
      // CSB0122RC0012000
      {
        fields: [
          'fc003', 'fc004', 'fc013', 'fc014', 'fc027', 'fc028', 'fc054', 'fc055', 'fc056',
          'fc081', 'fc082', 'fc189', 'fc190', 'fc191', 'fc192', 'fc193', 'fc197', 'fc198',
          'fc210', 'fc212', 'fc213', 'fc217', 'fc225', 'fc226', 'fc235', 'fc236', 'fc245',
          'fc246', 'fc255', 'fc256', 'fc278', 'fc202', 'fc203', 'fc061'
        ],
        validateDate11: function ({ value, field }) {
          if (!value) return true

          field.allowForce = false

          if (value == 'A') return true

          if (/[A, \?]/.test(value)) {
            return true
          }

          if (moment(`20${value}`).format('YYMMDD') === 'Invalid date') {
            return '日期格式录入错误'
          }

          return true
        }
      },

      // 12
      // CSB0122RC0013000
      // {
      //   fields: ['fc202', 'fc203'],
      //   validateDate12: function ({ value, field }) {
      //     if (!value) return true

      //     field.allowForce = true

      //     if(value == 'A') return true

      //     if (/[A, \?]/.test(value)) {
      //       return true
      //     }

      //     if (value.length !== 6 || moment(`20${value}`).format('YYYYMM') === 'Invalid date') {
      //       return '日期格式录入错误'
      //     }

      //     return true
      //   }
      // },

      // 13
      // CSB0122RC0014000
      {
        fields: ['fc019', 'fc029', 'fc227', 'fc237', 'fc247', 'fc257'],
        validateDate13: function ({ value, field }) {
          if (!value) return true

          field.allowForce = true

          if (value == 'A') return true

          if (/[A, \?]/.test(value)) {
            return true
          }

          if (value.length !== 8 || moment(`${value}`).format('YYYYMMDD') === 'Invalid date') {
            return '日期格式录入错误'
          }

          return true
        }
      },

      // 25
      // CSB0122RC0026000
      {
        fields: ['fc164', 'fc165', 'fc166', 'fc167', 'fc168', 'fc169', 'fc170', 'fc171', 'fc172', 'fc173', 'fc497', 'fc498', 'fc499', 'fc500', 'fc501', 'fc502', 'fc503', 'fc504', 'fc505', 'fc506'],
        validate25: function ({ field, value, items }) {
          if (value === '?') {
            return true
          }
          if (value === 'B') {
            return true
          }
          field.allowForce = false
          const result = items.find((text) => text === value)

          if (!result) {
            return '录入内容不在数据库中，不可强过!'
          }

          return true
        }
      },

      // 35
      {
        fields: ['fc205'],
        validate35: function ({ field, value, items }) {

          field.allowForce = false
          // 录入内容为F时不做校验
          if (value === 'F') {
            return true
          }

          const result = items.find(text => text == value)

          if (result) {
            return true
          }

          return '录入内容不在数据库中，不可强过！'
        }
      },

      {
        fields: ['fc204'],
        validate351: function ({ field, value, items }) {

          field.allowForce = false
          // 录入内容为F时不做校验
          if (value === 'F') {
            return true
          }

          const result = items.find(text => text == value)

          if (result) {
            return true
          }

          return '录入内容不在数据库中，不可强过！'
        }
      },

      {
        fields: ['fc474'],
        validate352: function ({ field, value, items }) {

          field.allowForce = false
          // 录入内容为F时不做校验
          if (value === 'F') {
            return true
          }

          const result = items.find(text => text == value)

          if (result) {
            return true
          }

          return '录入内容不在数据库中，不可强过！'
        }
      },

      {
        fields: ['fc475'],
        validate353: function ({ field, value, items }) {

          field.allowForce = false
          // 录入内容为F时不做校验
          if (value === 'F') {
            return true
          }

          const result = items.find(text => text == value)

          if (result) {
            return true
          }

          return '录入内容不在数据库中，不可强过！'
        }
      },

      {
        fields: ['fc476'],
        validate354: function ({ field, value, items }) {

          field.allowForce = false
          // 录入内容为F时不做校验
          if (value === 'F') {
            return true
          }

          const result = items.find(text => text == value)

          if (result) {
            return true
          }

          return '录入内容不在数据库中，不可强过！'
        }
      },

      // 26
      {
        fields: ['fc184', 'fc185', 'fc186', 'fc187', 'fc188'],
        validate26: function ({ field, value, items }) {
          field.allowForce = false
          // 录入内容为F时不做校验
          if (value === 'F') {
            return true
          }

          const result = items.find(text => text == value)

          if (result) {
            return true
          }

          return '录入内容不在数据库中，不可强过！'
        }
      },

      // 27
      {
        fields: ['fc022', 'fc032', 'fc230', 'fc240', 'fc250', 'fc260', 'fc263'],
        validate27: function ({ field, value, items }) {
          field.allowForce = false
          // 录入内容为F时不做校验
          if (value === 'F') {
            return true
          }

          const result = items.find(text => text == value)

          if (result) {
            return true
          }

          return '录入内容不在数据库中，不可强过！'
        }
      },

      // 28
      {
        fields: ['fc023', 'fc033', 'fc231', 'fc241', 'fc251', 'fc261'],
        validate28: function ({ field, value, items }) {
          field.allowForce = false
          // 录入内容为F时不做校验
          if (value === 'F') {
            return true
          }

          const result = items.find(text => text.includes(value))

          if (result) {
            return true
          }

          return '录入内容不在数据库中，不可强过！'
        }
      },

      // 29
      {
        fields: ['fc214'],
        validate29: function ({ field, value, items }) {
          field.allowForce = false
          // 录入内容为F时不做校验
          if (value === 'F') {
            return true
          }

          const result = items.find(text => text.includes(value))

          if (result) {
            return true
          }

          return '录入内容不在数据库中，不可强过！'
        }
      },

      // 30
      {
        fields: ['fc216'],
        validate30: function ({ field, value, items }) {
          field.allowForce = false
          // 录入内容为F时不做校验
          if (value === 'F') {
            return true
          }

          const result = items.find(text => text.includes(value))

          if (result) {
            return true
          }

          return '录入内容不在数据库中，不可强过！'
        }
      },

      // 31
      {
        fields: ['fc211'],
        validate31: function ({ field, value, items }) {
          field.allowForce = false
          // 录入内容为F时不做校验
          if (value === 'F') {
            return true
          }

          const result = items.find(text => text == value)

          if (result) {
            return true
          }

          return '录入内容不在数据库中，不可强过！'
        }
      },

      // 32
      {
        fields: ['fc218'],
        validate32: function ({ field, value, items }) {
          field.allowForce = false
          // 录入内容为F时不做校验
          if (value === 'F' || value === '?') {
            return true
          }

          const result = items.find(text => text.includes(value))

          if (result) {
            return true
          }

          return '录入内容不在数据库中，不可强过！'
        }
      },

      // 16
      {
        fields: [
          'fc402', 'fc403', 'fc404', 'fc405', 'fc406', 'fc407', 'fc408', 'fc409', 'fc410', 'fc411',
          'fc412', 'fc413', 'fc414', 'fc415', 'fc416', 'fc417', 'fc418', 'fc419', 'fc420', 'fc421',
          'fc422', 'fc423', 'fc424', 'fc425', 'fc426', 'fc427', 'fc428', 'fc429', 'fc430', 'fc431',
        ],
        validate16: function ({ field, value, items }) {
          field.allowForce = false
          if (value === '?') {
            return true
          }
          const result = items.find(text => text == value)

          if (result) {
            return true
          }

          return '录入内容不在数据库中，不可强过！'
        }
      }
    ],

    // 提示文本
    hints: [
      {
        fields: ['fc280'],
        hintFc001: function () {
          return `<p style="color: red;">注意：不要漏掉“其它支付、大病支付、大额支付、医疗救助、公务员补助、各种补助报销金额”</p>`
        }
      },
      {
        fields: ['fc130'],
        hintFc002: function () {
          return `<p style="color: red;">注意：不要漏掉“其它支付、大病支付、大额支付、医疗救助、公务员补助、各种补助报销金额”</p>`
        }
      },
      {
        fields: ['fc338'],
        hintFc003: function () {
          return `<p style="color: red;">注意：发票上的个人自费=总金额不录入；自费+统筹金额大于或等于总金额不录入；部分自付、自付二、乙类、超限价不录入；</p>`
        }
      },
    ],

    // 字段已生成
    updateFields: {
      methods: {
        // 14
        // CSB0122RC0015000
        disable14({ op, block, fieldsList, focusFieldsIndex }) {
          if (block.code != 'bc004') return
          const fields = fieldsList[focusFieldsIndex]
          const codes = ['fc016', 'fc021']
          const fc015Field = fields.find(field => field.code == 'fc015')

          fc015Field[`${op}Value`] = '2'
          fc015Field.resultValue = '2'
          fc015Field.disabled = true

          fields?.map(_field => {
            if (codes.includes(_field.code)) {
              _field.disabled = true
            }
          })

        },

        // 16
        // CSB0122RC0017000
        set16Items({ bill, fieldsList, focusFieldsIndex }) {
          const dropArr = [
            'fc402', 'fc403', 'fc404', 'fc405', 'fc406', 'fc407', 'fc408', 'fc409', 'fc410', 'fc411',
            'fc412', 'fc413', 'fc414', 'fc415', 'fc416', 'fc417', 'fc418', 'fc419', 'fc420', 'fc421',
            'fc422', 'fc423', 'fc424', 'fc425', 'fc426', 'fc427', 'fc428', 'fc429', 'fc430', 'fc431',
          ]
          // 数据库
          // const db = window['constantsDB']['B0122']
          // if (!db) return
          // const collections = db['B0122_信诚理赔_机构号代码表']
          // // 案件的机构号
          // const agency = bill.agency
          // // 机构代码: agencyCode 地址: address
          // const [address, agencyCode] = [[], []]
          // // 本级机构代码对应的地址
          // let addressValue = ''
          // // 指定的地址目录序号
          // let index = ''
          // if (agency) {
          //   for (let dessert of collections.desserts) {
          //     address.push(dessert[0])
          //     agencyCode.push(dessert[1])
          //   }

          //   index = agencyCode.indexOf(agency)
          //   addressValue = address[index]
          // }

          const fields = fieldsList[focusFieldsIndex]

          if (bill.agency == '北京') {
            fields.map(_field => {
              if (dropArr.includes(_field.code)) {
                _field.table = {
                  name: `B0122_信诚理赔_北京发票大项明细表`,
                  query: '名称'
                }
              }
            })
          } else {
            fields.map(_field => {
              if (dropArr.includes(_field.code)) {
                _field.table = {
                  name: `B0122_信诚理赔_发票大项明细表`,
                  query: '名称'
                }
              }
            })
          }
        },

        // 19
        // CSB0122RC0020000
        disable19({ block, fieldsList, focusFieldsIndex }) {
          if (block.code !== 'bc019') return true
          const fields = fieldsList[focusFieldsIndex]
          const codes = [
            'fc138', 'fc296', 'fc297', 'fc298', 'fc299', 'fc300', 'fc301', 'fc302',
            'fc140', 'fc310', 'fc311', 'fc312', 'fc313', 'fc314', 'fc315', 'fc316',
            'fc149', 'fc374', 'fc375', 'fc376', 'fc377', 'fc378', 'fc379', 'fc380',
            'fc150', 'fc381', 'fc382', 'fc383', 'fc384', 'fc385', 'fc386', 'fc387',
            'fc151', 'fc388', 'fc389', 'fc390', 'fc391', 'fc392', 'fc393', 'fc394',
            'fc152', 'fc395', 'fc396', 'fc397', 'fc398', 'fc399', 'fc400', 'fc401',
          ]

          fields?.map(_field => {
            if (codes.includes(_field.code)) {
              _field.disabled = true
            }
          })
        },

        // 21
        // CSB0122RC0022000
        setConstants21({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc139', 'fc303', 'fc304', 'fc305', 'fc306', 'fc307', 'fc308', 'fc309']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0122_信诚理赔_全国',
                query: '项目名称'
              }
            }
          })
        },

        // 25
        // CSB0122RC0026000
        setConstants25({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc164', 'fc165', 'fc166', 'fc167', 'fc168', 'fc169', 'fc170', 'fc171', 'fc172', 'fc173', , 'fc497', 'fc498', 'fc499', 'fc500', 'fc501', 'fc502', 'fc503', 'fc504', 'fc505', 'fc506']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0122_信诚理赔_伤病代码表',
                query: '名称'
              }
            }
          })
        },

        // 26
        setConstants26({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc184', 'fc185', 'fc186', 'fc187', 'fc188']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0122_信诚理赔_手术代码表',
                query: '名称'
              }
            }
          })
        },

        // 27
        setConstants27({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc022', 'fc032', 'fc230', 'fc240', 'fc250', 'fc260', 'fc263']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0122_信诚理赔_银行代码表',
                query: '全称'
              }
            }
          })
        },

        // 28
        setConstants28({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc023', 'fc033', 'fc231', 'fc241', 'fc251', 'fc261']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0122_信诚理赔_开户城市表',
                query: '开户城市'
              }
            }
          })
        },

        // 29
        setConstants29({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc214']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0122_信诚理赔_意外原因代码表',
                query: '名称'
              }
            }
          })
        },

        // 30
        setConstants30({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc216']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0122_信诚理赔_意外发生地',
                query: '发生地'
              }
            }
          })
        },

        // 31
        setConstants31({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc211']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0122_信诚理赔_重疾原因代码表',
                query: '名称'
              }
            }
          })
        },

        // 32
        setConstants32({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc218']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0122_信诚理赔_身故原因代码表',
                query: '名称'
              }
            }
          })
        },

        // 35
        setConstants33: function ({ flatFieldsList }) {
          const fields = ['fc205']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0122_信诚理赔_行政区划代码',
                query: '市',
                // targets: ['省', '市']
              }
            }
          })
        },

        // 35
        setConstants34: function ({ flatFieldsList }) {
          const fields = ['fc204']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0122_信诚理赔_行政区划代码',
                query: '省'
              }
            }
          })
        },

        setConstants331: function ({ flatFieldsList }) {
          const fields = ['fc474']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0122_信诚理赔_行政区划代码',
                query: '区',
                // targets: ['省', '市', '区']
              }
            }
          })
        },

        setConstants332: function ({ flatFieldsList }) {
          const fields = ['fc475']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0122_信诚理赔_行政区划代码',
                query: '市',
                // targets: ['省', '市']
              }
            }
          })
        },

        setConstants333: function ({ flatFieldsList }) {
          const fields = ['fc476']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0122_信诚理赔_行政区划代码',
                query: '省',
                // targets: ['省']
              }
            }
          })
        },

        // 42
        // CSB0122RC0043000
        disable43({ block, fieldsList, focusFieldsIndex, codeValues = {} }) {
          if (block.code != 'bc019') return
          const { fc010 } = codeValues
          const codesList = ['fc143', 'fc331', 'fc332', 'fc333', 'fc334', 'fc335', 'fc336', 'fc337']
          const fields = fieldsList[focusFieldsIndex]
          fields?.map(_field => {
            if (codesList.includes(_field.code) && fc010 == '2') {
              _field.disabled = true
            }
          })
        },

        disabled44({ fieldsList, focusFieldsIndex, codeValues = {} }) {
          const { fc265 } = codeValues
          const codeMaps = new Map([
            ['1', ['fc274', 'fc275', 'fc276', 'fc277', 'fc266', 'fc267', 'fc268', 'fc269']],
            ['2', ['fc275', 'fc276', 'fc277', 'fc267', 'fc268', 'fc269']],
            ['3', ['fc276', 'fc277', 'fc268', 'fc269']],
            ['4', ['fc277', 'fc269']],
          ])

          const fields = fieldsList[focusFieldsIndex]

          if (codeMaps.get(fc265)) {
            let arr = codeMaps.get(fc265)
            for (let _field of fields) {
              if (arr.includes(_field.code)) {
                _field.disabled = true
              }
            }
          }
        },

        disable196({ fieldsList, focusFieldsIndex }) {

          const fields = fieldsList[focusFieldsIndex]
          const fc196Field = fields.find(field => field.code == 'fc196')

          if (fc196Field) {
            fc196Field.disabled = true
          }
        },

        disable197({ fieldsList, focusFieldsIndex }) {

          const fields = fieldsList[focusFieldsIndex]
          const fc490Field = fields.find(field => field.code == 'fc490')
          const fc491Field = fields.find(field => field.code == 'fc491')

          if (fc490Field && fc491Field) {
            fc490Field.disabled = true
            fc491Field.disabled = true
          }
        },

        disable198({ fieldsList, focusFieldsIndex, codeValues = {} }) {
          const { fc135 } = codeValues
          const codesList = ['fc487', 'fc488', 'fc489']
          const fields = fieldsList[focusFieldsIndex]
          fields?.map(_field => {
            if (codesList.includes(_field.code) && fc135 == '2') {
              _field.disabled = true
            }
          })
        },

        disable199({ fieldsList, focusFieldsIndex, codeValues = {} }) {
          const { fc486 } = codeValues

          const fields = fieldsList[focusFieldsIndex]
          const codeList = ['fc487', 'fc488', 'fc489']
          if (fc486 == '1') {
            for (let _field of fields) {
              if (codeList.includes(_field.code)) {
                _field.disabled = true
              }
            }
          }
        },

        disabled200({ fieldsList, focusFieldsIndex, codeValues = {} }) {

          const { fc486 } = codeValues
          const fields = fieldsList[focusFieldsIndex]
          const codeList = ['fc487', 'fc488', 'fc489', 'fc490']
          const fc491Field = fields.find(field => field.code == 'fc491')

          if (fc491Field && fc491Field.resultValue == '' && fc486 == '1') {
            for (let _field of fields) {
              if (codeList.includes(_field.code)) {
                _field.disabled = false
              }
            }
          }
        },

        disable201({ op, block, fieldsList, focusFieldsIndex, codeValues = {} }) {
          if (this.op != 'op0') return
          const { fc484 } = codeValues
          console.log(fc484);
          const codesList = ['fc061', 'fc062', 'fc063', 'fc064', 'fc065', 'fc066', 'fc067', 'fc068', 'fc069', 'fc070', 'fc077', 'fc076', 'fc075', 'fc058', 'fc059']
          const fields = fieldsList[focusFieldsIndex]
          fields?.map(_field => {
            if (codesList.includes(_field.code) && fc484 != '1') {
              _field.disabled = true
            }
          })
        },
      }
    },

    // 回车
    enter: {
      methods: {
        // 14
        // CSB0122RC0015000
        disabled14({ block, field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc015') return
          const fields = fieldsList[focusFieldsIndex]

          const fc016Field = fields.find(field => field.code == 'fc016')
          const codes = ['fc022', 'fc023', 'fc017', 'fc021']

          if (block.code == 'bc008' && field.resultValue == '2') {
            fc016Field.disabled = true
          } else {
            fc016Field.disabled = false
          }

          if (block.code == 'bc012' || block.code == 'bc016') {
            if (field.resultValue == '1') {
              fields?.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = true
                }
              })
            } else {
              fields?.map(_field => {
                if (codes.includes(_field.code)) {
                  _field.disabled = false
                }
              })
            }

            if (field.resultValue == '2') {
              fc016Field.disabled = true
            } else {
              fc016Field.disabled = false
            }
          }
        },

        // 15
        // CSB0122RC0016000
        disabled15({ field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc279') return
          const codeMaps = new Map([
            ['1', ['fc081', 'fc082']],
            ['2', ['fc278']],
            ['3', ['fc081', 'fc082']],
          ])

          const codeValue = ['fc081', 'fc082', 'fc278', 'fc081', 'fc082']
          const fields = fieldsList[focusFieldsIndex]

          if (codeMaps.get(field.resultValue)) {
            for (let _field of fields) {
              if (codeValue.includes(_field.code)) {
                _field.disabled = false
              }
            }
            let arr = codeMaps.get(field.resultValue)

            for (let _field of fields) {
              if (arr.includes(_field.code)) {
                _field.disabled = true
              }
            }
          }
        },

        // 23
        // CSB0122RC002400
        disabled23({ op, field, flatFieldsList }) {
          const codesList = new Map([
            ['fc145', ['fc144', 'fc462']],
            ['fc346', ['fc339', 'fc463']],
            ['fc347', ['fc340', 'fc464']],
            ['fc348', ['fc341', 'fc465']],
            ['fc349', ['fc342', 'fc466']],
            ['fc350', ['fc343', 'fc467']],
            ['fc351', ['fc344', 'fc468']],
            ['fc352', ['fc345', 'fc469']],
          ])
          const codes = ['fc145', 'fc346', 'fc347', 'fc348', 'fc349', 'fc350', 'fc351', 'fc352']

          if (codes.includes(field.code)) {
            if ((field.resultValue == '1' || field.resultValue == '2' || field.resultValue == '3' || field.resultValue == '5')) {
              let disableArr = codesList.get(field.code)
              for (let field of flatFieldsList) {
                if (disableArr.includes(field.code)) {
                  field.disabled = true
                }
              }
            } else if (field.resultValue == '4') {
              let disableArr = codesList.get(field.code)[1]
              let disableArrs = codesList.get(field.code)
              for (let field of flatFieldsList) {
                if (disableArrs.includes(field.code)) {
                  field.disabled = false
                }
              }
              for (let field of flatFieldsList) {
                if (disableArr == field.code) {
                  // field[`${op}Value`] = ''
                  // field.resultValue = ''
                  field.disabled = true
                }
              }
            } else if (field.resultValue == '6') {
              let disableArr = codesList.get(field.code)[0]
              let disableArrs = codesList.get(field.code)
              for (let field of flatFieldsList) {
                if (disableArrs.includes(field.code)) {
                  field.disabled = false
                }
              }
              for (let field of flatFieldsList) {
                if (disableArr == field.code) {
                  field.disabled = true
                }
              }
            } else {
              for (let field of flatFieldsList) {
                field.disabled = false
              }
            }
          }
        },

        // 24
        // CSB0122RC0025000
        disabled24({ field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc194') return
          const codeMaps = new Map([
            ['1', ['fc165', 'fc166', 'fc167', 'fc168', 'fc169', 'fc170', 'fc171', 'fc172', 'fc173', 'fc175', 'fc176', 'fc177', 'fc178', 'fc179', 'fc180', 'fc181', 'fc182', 'fc183', 'fc497', 'fc498', 'fc499', 'fc500', 'fc501', 'fc502', 'fc503', 'fc504', 'fc505', 'fc506', 'fc507', 'fc508', 'fc509', 'fc510', 'fc511', 'fc512', 'fc513', 'fc514', 'fc515']],
            ['2', ['fc166', 'fc167', 'fc168', 'fc169', 'fc170', 'fc171', 'fc172', 'fc173', 'fc176', 'fc177', 'fc178', 'fc179', 'fc180', 'fc181', 'fc182', 'fc183', 'fc497', 'fc498', 'fc499', 'fc500', 'fc501', 'fc502', 'fc503', 'fc504', 'fc505', 'fc506', 'fc507', 'fc508', 'fc509', 'fc510', 'fc511', 'fc512', 'fc513', 'fc514', 'fc515', 'fc516']],
            ['3', ['fc167', 'fc168', 'fc169', 'fc170', 'fc171', 'fc172', 'fc173', 'fc177', 'fc178', 'fc179', 'fc180', 'fc181', 'fc182', 'fc183']],
            ['4', ['fc168', 'fc169', 'fc170', 'fc171', 'fc172', 'fc173', 'fc178', 'fc179', 'fc180', 'fc181', 'fc182', 'fc183']],
            ['5', ['fc169', 'fc170', 'fc171', 'fc172', 'fc173', 'fc179', 'fc180', 'fc181', 'fc182', 'fc183']],
            ['6', ['fc170', 'fc171', 'fc172', 'fc173', 'fc180', 'fc181', 'fc182', 'fc183']],
            ['7', ['fc171', 'fc172', 'fc173', 'fc181', 'fc182', 'fc183']],
            ['8', ['fc172', 'fc173', 'fc182', 'fc183']],
            ['9', ['fc173', 'fc183']]
          ])

          const fields = fieldsList[focusFieldsIndex]
          const codeValue = ['fc165', 'fc166', 'fc167', 'fc168', 'fc169', 'fc170', 'fc171', 'fc172', 'fc173', 'fc175', 'fc176', 'fc177', 'fc178', 'fc179', 'fc180', 'fc181', 'fc182', 'fc183']

          if (codeMaps.get(field.resultValue)) {
            for (let _field of fields) {
              if (codeValue.includes(_field.code)) {
                _field.disabled = false
              }
            }
            let arr = codeMaps.get(field.resultValue)
            for (let _field of fields) {
              if (arr.includes(_field.code)) {
                _field.disabled = true
              }
            }
          }
        },

        // 26
        // CSB0122RC0027000
        // disabled26({ field, fieldsList, focusFieldsIndex }) {
        //   const fields = fieldsList[focusFieldsIndex]
        //   const codeMap = new Map([
        //     ['fc164', 'fc174'],
        //     ['fc165', 'fc175'],
        //     ['fc166', 'fc176'],
        //     ['fc167', 'fc177'],
        //     ['fc168', 'fc178'],
        //     ['fc169', 'fc179'],
        //     ['fc170', 'fc180'],
        //     ['fc171', 'fc181'],
        //     ['fc172', 'fc182'],
        //     ['fc173', 'fc183'],
        //   ])
        //   const codes = ['fc164', 'fc165', 'fc166', 'fc167', 'fc168', 'fc169', 'fc170', 'fc171', 'fc172', 'fc173']


        //   if (codes.includes(field.code) && field.resultValue == 'B') {
        //     let code = codeMap.get(field.code)
        //     console.log(code);
        //     fields?.map(_field => {
        //       if (_field.code == code) {
        //         _field.disabled = false
        //       }
        //     })
        //   } else if (codes.includes(field.code) && field.resultValue != 'B') {
        //     let code = codeMap.get(field.code)
        //     fields?.map(_field => {
        //       if (_field.code == code) {
        //         _field.disabled = true
        //       }
        //     })
        //   }
        // },

        // 27
        // CSB0122RC0028000
        disabled27({ field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc195') return
          const codeMaps = new Map([
            ['1', ['fc185', 'fc186', 'fc187', 'fc188', 'fc190', 'fc191', 'fc192', 'fc193']],
            ['2', ['fc186', 'fc187', 'fc188', 'fc191', 'fc192', 'fc193']],
            ['3', ['fc187', 'fc188', 'fc192', 'fc193']],
            ['4', ['fc188', 'fc193']],
          ])

          const fields = fieldsList[focusFieldsIndex]
          const codeValue = ['fc185', 'fc186', 'fc187', 'fc188', 'fc190', 'fc191', 'fc192', 'fc193']

          if (codeMaps.get(field.resultValue)) {
            for (let _field of fields) {
              if (codeValue.includes(_field.code)) {
                _field.disabled = false
              }
            }
            let arr = codeMaps.get(field.resultValue)
            for (let _field of fields) {
              if (arr.includes(_field.code)) {
                _field.disabled = true
              }
            }
          }
        },

        // 36
        // CSB0122RC0037000
        // disabled36({ field, fieldsList, focusFieldsIndex }) {
        //   if (field.code != 'fc205') return
        //   const fields = fieldsList[focusFieldsIndex]
        //   const fc204Field = fields.find(field => field.code == 'fc204')

        //   if (field.resultValue == '北京市' || field.resultValue == '上海市' || field.resultValue == '天津市' || field.resultValue == '重庆市') {
        //     fc204Field.disabled = true
        //   } else {
        //     fc204Field.disabled = false
        //   }
        // },

        // 35
        // CSB0122RC0036000
        validate33And35: function ({ field, fieldsList, focusFieldsIndex, memoFields }) {
          if (field.resultValue == 'F') return
          if (field.code === 'fc205') {
            const fields = fieldsList[focusFieldsIndex]

            const fc204Field = fields.find(field => field.code === 'fc204')

            // if (fc205Value.includes('-')) {
            //   fc204Field.disabled = false
            //   const values = fc205Value.split('-')
            //   if (values[0] == values[1] || values[1] == '') {
            //     fc204Field.disabled = true
            //   } else {
            //     fc204Field.disabled = false
            //   }
            //   if (values[1] != '') {
            //     field[`${op}Value`] = values[1]
            //     field.resultValue = values[1]
            //     _.set(memoFields, `${field.uniqueId}.value`, values[1])
            //   } else {
            //     field[`${op}Value`] = values[0]
            //     field.resultValue = values[0]
            //     _.set(memoFields, `${field.uniqueId}.value`, values[0])
            //   }
            // }
            let count = 0
            for (let el of field.items) {
              if (el == field.resultValue) count++
            }

            if (count >= 2) {
              fc204Field.disabled = false
            } else {
              fc204Field.disabled = true
            }
          }
        },

        validate331And351: function ({ field, fieldsList, focusFieldsIndex }) {
          if (field.resultValue == 'F') return
          if (field.code === 'fc474') {
            const fields = fieldsList[focusFieldsIndex]

            const fc475Field = fields.find(field => field.code === 'fc475')
            const fc476Field = fields.find(field => field.code === 'fc476')

            let count = 0
            for (let el of field.items) {
              if (el == field.resultValue) count++
            }

            if (count >= 2) {
              fc475Field.disabled = false
              fc476Field.disabled = false
            } else {
              fc475Field.disabled = true
              fc476Field.disabled = true
            }
          }
        },

        validate332And352: function ({ op, field, fieldsList, focusFieldsIndex }) {
          if (field.resultValue == 'F') return
          if (field.code === 'fc475') {
            const fields = fieldsList[focusFieldsIndex]

            const fc476Field = fields.find(field => field.code === 'fc476')

            let count = 0
            for (let el of field.items) {
              if (el == field.resultValue) count++
            }

            if (count >= 2) {
              fc476Field.disabled = false
            } else {
              fc476Field.disabled = true
            }

          }
        },

        // 41
        // CSB0122RC0042000
        disabled42({ field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc265') return
          const codeMaps = new Map([
            ['0', ['fc033', 'fc035', 'fc231', 'fc270', 'fc241', 'fc271', 'fc251', 'fc272', 'fc261', 'fc273']],
            ['1', ['fc231', 'fc270', 'fc241', 'fc271', 'fc251', 'fc272', 'fc261', 'fc273']],
            ['2', ['fc241', 'fc271', 'fc251', 'fc272', 'fc261', 'fc273']],
            ['3', ['fc251', 'fc272', 'fc261', 'fc273']],
            ['4', ['fc261', 'fc273']],
          ])

          const fields = fieldsList[focusFieldsIndex]

          const codeValue = ['fc033', 'fc035', 'fc231', 'fc270', 'fc241', 'fc271', 'fc251', 'fc272', 'fc261', 'fc273']

          if (codeMaps.get(field.resultValue)) {
            for (let _field of fields) {
              if (codeValue.includes(_field.code)) {
                _field.disabled = false
              }
            }
            let arr = codeMaps.get(field.resultValue)
            for (let _field of fields) {
              if (arr.includes(_field.code)) {
                _field.disabled = true
              }
            }
          }
        },

        disabled185({ field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc063') return

          let arr = ['fc064', 'fc065', 'fc066']
          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue == 'A' || field.resultValue == '2') {
            for (let _field of fields) {
              if (arr.includes(_field.code)) {
                _field.disabled = true
              }
            }
          } else {
            for (let _field of fields) {
              if (arr.includes(_field.code)) {
                _field.disabled = false
              }
            }
          }
        },

        disabled186({ field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc067') return

          let arr = ['fc068', 'fc069', 'fc070']
          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue == 'A') {
            for (let _field of fields) {
              if (arr.includes(_field.code)) {
                _field.disabled = true
              }
            }
          } else {
            for (let _field of fields) {
              if (arr.includes(_field.code)) {
                _field.disabled = false
              }
            }
          }
        },

        disabled187({ field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc051') return

          let arr = ['fc052', 'fc053']
          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue == 'A') {
            for (let _field of fields) {
              if (arr.includes(_field.code)) {
                _field.disabled = true
              }
            }
          } else {
            for (let _field of fields) {
              if (arr.includes(_field.code)) {
                _field.disabled = false
              }
            }
          }
        },

        disabled188({ field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc045') return

          let arr = ['fc046', 'fc047']
          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue == 'A') {
            for (let _field of fields) {
              if (arr.includes(_field.code)) {
                _field.disabled = true
              }
            }
          } else {
            for (let _field of fields) {
              if (arr.includes(_field.code)) {
                _field.disabled = false
              }
            }
          }
        },

        disabled189({ field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc007') return

          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue == '4') {
            for (let _field of fields) {
              if (_field.code == 'fc130') {
                _field.disabled = true
              }
            }
          } else {
            for (let _field of fields) {
              if (_field.code == 'fc130') {
                _field.disabled = false
              }
            }
          }
        },

        disabled190({ field, fieldsList, focusFieldsIndex }) {
          const codeArr = ['fc018', 'fc026', 'fc224', 'fc234', 'fc244', 'fc254']
          if (!codeArr.includes(field.code)) return
          const fieldMap = new Map([
            ['fc018', ['fc019', 'fc020']],
            ['fc026', ['fc029', 'fc030']],
            ['fc224', ['fc227', 'fc228']],
            ['fc234', ['fc237', 'fc238']],
            ['fc244', ['fc247', 'fc248']],
            ['fc254', ['fc257', 'fc258']],
          ])

          // 18位身份证号码的正则表达式
          const pattern1 = /^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[1-2][0-9]|3[0-1])\d{3}(\d|X|x)$/;
          // 15位身份证号码的正则表达式
          const pattern2 = /^[1-9]\d{5}\d{2}(0[1-9]|1[0-2])(0[1-9]|[1-2][0-9]|3[0-1])\d{2}\d{1}$/;

          const fields = fieldsList[focusFieldsIndex]

          if (pattern1.test(field.resultValue) || pattern2.test(field.resultValue)) {
            for (let _field of fields) {
              if (fieldMap.get(field.code).includes(_field.code)) {
                _field.disabled = true
              }
            }
          } else {
            for (let _field of fields) {
              if (fieldMap.get(field.code).includes(_field.code)) {
                _field.disabled = false
              }
            }
          }
        },

        disabled191({ field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc214') return

          const fields = fieldsList[focusFieldsIndex]

          const codeArr = ['fc215', 'fc216']
          if (field.resultValue == 'A') {
            for (let _field of fields) {
              if (codeArr.includes(_field.code)) {
                _field.disabled = true
              }
            }
          } else {
            for (let _field of fields) {
              if (codeArr.includes(_field.code)) {
                _field.disabled = false
              }
            }
          }
        },

        disabled192({ field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc338') return

          const fields = fieldsList[focusFieldsIndex]

          const valueArr = ['fc492', 'fc493', 'fc494', 'fc495']
          if (field.resultValue == 'A') {
            for (let _field of fields) {
              if (valueArr.includes(_field.code)) {
                _field.disabled = true
              }
            }
          } else {
            for (let _field of fields) {
              if (valueArr.includes(_field.code)) {
                _field.disabled = false
              }
            }
          }
        },

        // disabled190({ field, fieldsList, focusFieldsIndex }) {
        //   if (field.code != 'fc048') return

        //   const fields = fieldsList[focusFieldsIndex]

        //   if (!(field.resultValue).includes('特需')) {
        //     for (let _field of fields) {
        //       if (_field.code == 'fc044') {
        //         _field.disabled = true
        //       }
        //     }
        //   } else {
        //     for (let _field of fields) {
        //       if (_field.code == 'fc044') {
        //         _field.disabled = false
        //       }
        //     }
        //   }
        // },
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
        // 17
        // CSB0122RC0018000
        validate17: function ({ fieldsList }) {
          const fields = fieldsList[0]
          const codeMap = new Map([
            ['fc402', 'fc432'],
            ['fc403', 'fc433'],
            ['fc404', 'fc434'],
            ['fc405', 'fc435'],
            ['fc406', 'fc436'],
            ['fc407', 'fc437'],
            ['fc408', 'fc438'],
            ['fc409', 'fc439'],
            ['fc410', 'fc440'],
            ['fc411', 'fc441'],
            ['fc412', 'fc442'],
            ['fc413', 'fc443'],
            ['fc414', 'fc444'],
            ['fc415', 'fc445'],
            ['fc416', 'fc446'],
            ['fc417', 'fc447'],
            ['fc418', 'fc448'],
            ['fc419', 'fc449'],
            ['fc420', 'fc450'],
            ['fc421', 'fc451'],
            ['fc422', 'fc452'],
            ['fc423', 'fc453'],
            ['fc424', 'fc454'],
            ['fc425', 'fc455'],
            ['fc426', 'fc456'],
            ['fc427', 'fc457'],
            ['fc428', 'fc458'],
            ['fc429', 'fc459'],
            ['fc430', 'fc460'],
            ['fc431', 'fc461'],
          ])
          const codes = [
            'fc402', 'fc403', 'fc404', 'fc405', 'fc406', 'fc407', 'fc408', 'fc409', 'fc410', 'fc411',
            'fc412', 'fc413', 'fc414', 'fc415', 'fc416', 'fc417', 'fc418', 'fc419', 'fc420', 'fc421',
            'fc422', 'fc423', 'fc424', 'fc425', 'fc426', 'fc427', 'fc428', 'fc429', 'fc430', 'fc431',
          ]

          for (let field of fields) {
            if (codes.includes(field.code) && field.resultValue != '') {
              let code = codeMap.get(field.code)
              for (let _field of fields) {
                if (_field.code == code && _field.resultValue == '') {
                  _field.autofocus = true
                  return {
                    errorMessage: `内容录入遗漏， 请检查!`
                  }
                }
              }
            }
          }

          return true
        },

        // 18
        // CSB0122RC0019000
        validate18: function ({ block, fieldsList }) {
          if (block.code !== 'bc017') return true
          const fields = fieldsList[0]
          const fc080Field = tools.find(fields, { code: 'fc080' })
          const fc432Field = tools.find(fields, { code: 'fc432' })
          if (fc432Field.show == false) return true
          const codes = [
            'fc432', 'fc433', 'fc434', 'fc435', 'fc436', 'fc437', 'fc438', 'fc439', 'fc440', 'fc441',
            'fc442', 'fc443', 'fc444', 'fc445', 'fc446', 'fc447', 'fc448', 'fc449', 'fc450', 'fc451',
            'fc451', 'fc453', 'fc454', 'fc455', 'fc456', 'fc457', 'fc458', 'fc459', 'fc460', 'fc461',
          ]

          if (!fc080Field || fc080Field.resultValue === '?') {
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

          // const fc080Value = +fc080Field.resultValue

          // const diff = fc080Value - count

          const fc080Value = new BigNumber(+fc080Field.resultValue)

          const diff = fc080Value.minus(count).toString()

          if (diff != 0) {
            return {
              popup: 'confirm',
              errorMessage: `发票明细金额与总金额不一致，差额为${diff}，请确认并修改!`
            }
          }

          return true
        },

        // 20
        // CSB0122RC0021000
        validate20({ block, fieldsList }) {
          if (block.code !== 'bc019') return true

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

        // 22
        // CSB0122RC0023000
        // 清单内容遗漏请检查， 光标自动定位
        validate22: function ({ fieldsList, flatFieldsList }) {
          const fields = fieldsList[0]
          const codeMap = new Map([
            ['fc139', ['fc145', 'fc141', 'fc142', 'fc143', 'fc144', 'fc462']],
            ['fc303', ['fc346', 'fc317', 'fc324', 'fc331', 'fc339', 'fc463']],
            ['fc304', ['fc347', 'fc318', 'fc325', 'fc332', 'fc340', 'fc464']],
            ['fc305', ['fc348', 'fc319', 'fc326', 'fc333', 'fc341', 'fc465']],
            ['fc306', ['fc349', 'fc320', 'fc327', 'fc334', 'fc342', 'fc466']],
            ['fc307', ['fc350', 'fc321', 'fc328', 'fc335', 'fc343', 'fc467']],
            ['fc308', ['fc351', 'fc322', 'fc329', 'fc336', 'fc344', 'fc468']],
            ['fc309', ['fc352', 'fc323', 'fc330', 'fc337', 'fc345', 'fc469']],
          ])

          const codes = ['fc139', 'fc303', 'fc304', 'fc305', 'fc306', 'fc307', 'fc308', 'fc309']

          for (let field of fields) {
            let codeArr = codeMap.get(field.code)
            if (codes.includes(field.code) && field.resultValue != '' && field.isPractice == false) {
              for (let _field of fields) {
                if (_field.disabled == true) continue
                if (codeArr?.includes(_field.code) && _field.resultValue == '') {
                  flatFieldsList.map(field => {
                    field.autofocus = false;
                    field.uniqueKey = `enter_${field.uniqueId}_${Date.now()}`;
                  });
                  _field.autofocus = true
                  return {
                    errorMessage: `清单内容录入遗漏， 请检查!`
                  }
                }
              }
            } else if (codes.includes(field.code) && field.op1Value != '') {
              for (let _field of fields) {
                if (_field.disabled == true) continue
                if (codeArr?.includes(_field.code) && _field.op1Value == '') {
                  flatFieldsList.map(field => {
                    field.autofocus = false;
                    field.uniqueKey = `enter_${field.uniqueId}_${Date.now()}`;
                  });
                  _field.autofocus = true
                  return {
                    errorMessage: `清单内容录入遗漏， 请检查!`
                  }
                }
              }
            }
          }
          return true
        },

        // CSB0122RC0147000
        validate147: function ({ fieldsList, op }) {
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

export default B0122
