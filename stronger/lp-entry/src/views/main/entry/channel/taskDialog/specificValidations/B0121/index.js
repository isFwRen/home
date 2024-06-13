import BigNumber from "bignumber.js";
import { tools, sessionStorage } from 'vue-rocket'
import { ignoreFreeValue } from "../tools";
import moment from "moment";
import { MessageBox, Notification } from 'element-ui';

const B0121 = {
  op0: {
    // 记录最后一次存储的合法field
    memoFields: [],

    // 记录相同 code 的 field 的值
    memoFieldValues: ['fc002'],

    // fields 的值从 targets 里的值选择
    dropdownFields: [
      {
        targets: ['fc108'],
        fields: ['fc109', 'fc118', 'fc361', 'fc362', 'fc363']
      }
    ],

    // 校验规则
    rules: [
      {
        fields: ['fc108'],
        validate2: function ({ field, fieldsObject, thumbIndex, value }) {
          const fc108Values = []

          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage

            if (sessionStorage || thumbIndex === +key) {
              const _fieldsList = fieldsObject[key].fieldsList

              for (let _fields of _fieldsList) {
                for (let _field of _fields) {
                  if (_field.code === 'fc108' && _field.uniqueId !== field.uniqueId) {
                    fc108Values.push(_field.resultValue)
                  }
                }
              }
            }
          }

          if (fc108Values.includes(value)) {
            return '该发票属性已被使用，请重新定义'
          }

          return true
        }
      },
      {
        fields: ['fc109', 'fc118', 'fc361', 'fc362', 'fc363'],
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
        fields: ['fc065'],
        validatefc065: function ({ effectValidations, field, items, value }) {
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
            return '录入内容与代码表不一致，请按单录入后强过'
          }

          field.allowForce = true

          return true
        }
      },

      // CSB0121RC0125000
      {
        fields: ['fc002'],
        validate01: function ({ fieldsObject }) {
          let arr = []
          let flag = 0
          for (let key in fieldsObject) {
            const fieldsList = fieldsObject[key].fieldsList

            for (let fields of fieldsList) {
              for (let field of fields) {
                const { code, resultValue } = field

                if (code == 'fc002') {

                  arr.push(resultValue)
                }
              }
            }
          }
          arr.forEach(el => {
            if (el == '1') flag++
          })

          if (flag > 1) return '申请书已录入'
          return true
        }
      },

      // CSB0121RC0126000
      {
        fields: ['fc002'],
        validate49: function ({ bill, op, fieldsObject }) {
          if (op != 'op0') return true
          if (bill.agency == '8611') {
            let arr = []
            for (let key in fieldsObject) {
              const fieldsList = fieldsObject[key].fieldsList

              for (let fields of fieldsList) {
                for (let field of fields) {
                  const { code, resultValue } = field

                  if (code == 'fc002') {

                    arr.push(resultValue)
                  }
                }
              }
            }

            if (arr.includes('2')) {
              return '8611机构案件, 发票切14'
            }
          }
          return true
        },
      },

      // CSB0121RC0127000
      {
        fields: ['fc002'],
        validate50: function ({ bill, op, fieldsObject }) {
          if (op != 'op0') return true

          if (bill.agency != '8611') {
            let arr = []
            for (let key in fieldsObject) {
              const fieldsList = fieldsObject[key].fieldsList

              for (let fields of fieldsList) {
                for (let field of fields) {
                  const { code, resultValue } = field

                  if (code == 'fc002') {

                    arr.push(resultValue)
                  }
                }
              }
            }

            if (arr.includes('14')) return '非8611机构案件, 发票切2'

          }
          return true
        },
      },

      // CSB0121RC0128000
      {
        fields: ['fc002'],
        validate51: function ({ bill, op, fieldsObject }) {
          if (op != 'op0') return true

          if (bill.agency != '8611') {
            let arr = []
            for (let key in fieldsObject) {
              const fieldsList = fieldsObject[key].fieldsList

              for (let fields of fieldsList) {
                for (let field of fields) {
                  const { code, resultValue } = field

                  if (code == 'fc002') {

                    arr.push(resultValue)
                  }
                }
              }
            }

            if (arr.includes('15')) return '非8611机构案件, 不需要切诊断'
          }
          return true
        },
      }
    ],

    // 提示文本
    hints: [],

    // 工序完成初始化
    init: {
      methods: {}
    },

    // 字段已生成
    updateFields: {
      methods: {
        setConstants40({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc065']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0121_百年人寿团险理赔_百年理赔医院代码表',
                query: '医院名称'
              }
            }
          })
        },
        disable01({ fieldsList, focusFieldsIndex }) {
          const fields = fieldsList[focusFieldsIndex];
          const fc161Field = fields.find(field => field.code === "fc161");
          if (fc161Field) fc161Field.disabled = true;
        }
      }
    },

    // 回车
    enter: {
      methods: {
        // 同分块
        disable01({ op, field, fieldsList, focusFieldsIndex }) {
          if (field.code !== 'fc364') return
          const fields = fieldsList[focusFieldsIndex]

          if (field.resultValue == '1') {
            fields?.map(_field => {
              if (_field.code == 'fc365') {
                _field.disabled = true
              }
            })
          } else {
            fields?.map(_field => {
              if (_field.code == 'fc365') {
                _field[`${op}Value`] = ''
                _field.resultValue = ''
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
        validate5({ fieldsObject }) {
          const [fc108Values, fc109Values, fc118Values] = [[], [], []]
          for (let key in fieldsObject) {
            const sessionStorage = fieldsObject[key].sessionStorage
            const fieldsList = fieldsObject[key].fieldsList

            if (sessionStorage) {
              for (let fields of fieldsList) {
                for (let field of fields) {
                  const { code, resultValue } = field

                  if (code === 'fc108') {
                    resultValue && fc108Values.push(resultValue)
                  }

                  if (code === 'fc109') {
                    resultValue && fc109Values.push(resultValue)
                  }

                  if (code === 'fc118') {
                    resultValue && fc118Values.push(resultValue)
                  }
                }
              }
            }
          }
          let value
          let uniArr = [...new Set([...fc109Values, ...fc118Values])]
          const flag = fc108Values.every(item => {
            value = item
            return uniArr.includes(item)
          })
          if (flag) {
            return true
          } else {
            return {
              popup: 'confirm',
              errorMessage: `发票${value}没有匹配的清单或报销单,请检查`
            }
          }
        },
        validate45: function ({ mergeFieldsList, op }) {
          const flatFieldsList = tools.flatArray(mergeFieldsList)
          const fc002Field = flatFieldsList.find(field => field.code === 'fc002')
          const fc321Field = flatFieldsList.find(field => field.code === 'fc321')

          if (fc002Field && fc321Field) {
            const [fc002Values, fc321Values] = [[], []]

            flatFieldsList.map(field => {
              if (field.code === 'fc002') {
                if (!fc002Values.includes(field[`${op}Value`])) {
                  fc002Values.push(field[`${op}Value`])
                }
              }

              if (field.code === 'fc321') {
                if (!fc321Values.includes(field[`${op}Value`])) {
                  fc321Values.push(field[`${op}Value`])
                }
              }
            })
            if (fc321Values.includes('1') && !fc002Values.includes('16')) {
              return {
                errorMessage: '漏切手术'
              }
            }
          }
          return true
        },
        validate46: function ({ bill, mergeFieldsList, op }) {
          if (bill.agency == '8611') {
            const flatFieldsList = tools.flatArray(mergeFieldsList)

            const fc002Field = flatFieldsList.find(field => field.code === 'fc002')

            if (fc002Field) {
              const fc002Values = []

              flatFieldsList.map(field => {
                if (field.code === 'fc002') {
                  if (!fc002Values.includes(field[`${op}Value`])) {
                    fc002Values.push(field[`${op}Value`])
                  }
                }
              })

              if (fc002Values.includes('2')) {
                return {
                  errorMessage: '机构为8611, 发票全部切14, 不切2'
                }
              }
            }
          }
          return true

        },
        validate47: function ({ bill, mergeFieldsList, op }) {
          if (bill.agency == '8611') {
            const flatFieldsList = tools.flatArray(mergeFieldsList)

            const fc002Field = flatFieldsList.find(field => field.code === 'fc002')

            if (fc002Field) {
              const fc002Values = []

              flatFieldsList.map(field => {
                if (field.code === 'fc002') {
                  if (!fc002Values.includes(field[`${op}Value`])) {
                    fc002Values.push(field[`${op}Value`])
                  }
                }
              })

              if (!fc002Values.includes('15')) {
                return {
                  errorMessage: '机构为8611, 发票全部切14, 并切一个诊断15, 不切2'
                }
              }
            }
          }
          return true
        },

        // CSB0121RC0124000
        validate48({ fieldsObject }) {
          let arr = []
          for (let key in fieldsObject) {
            const fieldsList = fieldsObject[key].fieldsList

            for (let fields of fieldsList) {
              for (let field of fields) {
                const { code, resultValue } = field

                if (code == 'fc002') {

                  arr.push(resultValue)
                }
              }
            }
          }
          if (!arr.includes('1')) {
            return {
              errorMessage: `缺少申请书，请确认`
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
      {
        fields: ['fc214', 'fc215', 'fc216', 'fc217', 'fc218', 'fc219', 'fc220', 'fc221', 'fc222', 'fc223', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246'],
        validatefc214: function ({ effectValidations, field, items, value }) {
          if (ignoreFreeValue({ effectValidations, value })) return true

          const result = items.find((text) => text === value)

          if (value.includes('?')) {
            field.allowForce = true
            return true
          }
          else if (value.includes('？')) {
            field.allowForce = true
            return true
          }

          if (!result) {
            field.allowForce = false
            return '录入内容与代码表不一致，请按相似内容进行选录，不可以进行强过'
          }

          field.allowForce = true

          return true
        }
      },
      {
        fields: ['fc082', 'fc083', 'fc084', 'fc085', 'fc086', 'fc087', 'fc088', 'fc089'],
        validatefc082: function ({ effectValidations, field, items, value }) {
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
            return '录入内容与医疗目录不一致，请按相似内容进行选录。'
          }

          field.allowForce = true

          return true
        }
      },
      {
        fields: ['fc066', 'fc068', 'fc069'],
        validateDate: function ({ value }) {
          if (!value) return true

          if (/[A,B, \?]/.test(value)) {
            return true
          }

          if (value.length !== 6 || moment(`20${value}`).format('YYYYMMDD') === 'Invalid date') {
            return '日期不合法， 请按单录入后强过。'
          }

          return true
        }
      },
      {
        fields: ['fc020'],
        validateDate: function ({ value }) {
          if (!value) return true

          if (/[A,B,1, \?]/.test(value)) {
            return true
          }

          if (value.length !== 6 || moment(`20${value}`).format('YYYYMMDD') === 'Invalid date') {
            return '日期不合法， 请按单录入后强过。'
          }

          return true
        }
      },
      {
        fields: ['fc019'],
        validatefc019: function ({ fieldsList, value, op }) {
          if (fieldsList[0][0]?.[`${op}Value`] == 0 || fieldsList[0][0]?.[`${op}Value`] == 5) {
            var reg = /^(\d{18}|\d{15}|\d{17}x)$/;
            if (!reg.test(fieldsList[0][1]?.[`${op}Value`])) {
              return '证件类型为身份证或户口本， 证件号码不为身份证， 请确认后按单录入'
            }
          }
          return true
        }
      }
    ],

    // 提示文本
    hints: [],

    // 字段已生成
    updateFields: {
      methods: {
        disableFields123: function ({ op, fieldsList, focusFieldsIndex }) {
          if (op === 'op0') {
            return
          }
          const codesList = [
            ['fc360'],
            ['fc320', 'fc270', 'fc272', 'fc274', 'fc276', 'fc278', 'fc280'],
            ['fc350'],
            ['fc347', 'fc351']
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
        setConstants39: function ({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc214', 'fc215', 'fc216', 'fc217', 'fc218', 'fc219', 'fc220', 'fc221', 'fc222', 'fc223', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0121_百年人寿团险理赔_百年理赔费用项目代码表',
                query: '费用项目名称'
              }
            }
          })
        },

        setConstants57: function ({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc082', 'fc083', 'fc084', 'fc085', 'fc086', 'fc087', 'fc088', 'fc089']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0121_百年人寿团险理赔_全国',
                query: '项目名称'
              }
            }
          })
        },

        setConstants33: function ({ flatFieldsList }) {
          const fields = ['fc024']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0121_百年人寿团险理赔_百年理赔地址库',
                query: '区',
                targets: ['省', '市', '区']
              }
            }
          })
        },

        setConstants34: function ({ flatFieldsList }) {
          const fields = ['fc023']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0121_百年人寿团险理赔_百年理赔地址库',
                query: '市',
                targets: ['省', '市']
              }
            }
          })
        },

        // 34
        setConstants35: function ({ flatFieldsList }) {
          const fields = ['fc022']

          flatFieldsList.map(_field => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0121_百年人寿团险理赔_百年理赔地址库',
                query: '省'
              }
            }
          })
        },
        disableFields: function ({ op, fieldsList, focusFieldsIndex }) {
          if (op === 'op0') {
            return
          }

          const codesList = [
            ['fc327'],
            ['fc257'],
            ['fc204', 'fc205', 'fc206', 'fc207', 'fc208', 'fc209', 'fc210', 'fc211', 'fc212', 'fc213', 'fc247', 'fc248', 'fc249', 'fc250', 'fc251']
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
        setConstants41({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc348']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0121_百年人寿团险理赔_手术代码',
                query: '手术名称'
              }
            }
          })
        },
        setConstants42({ fieldsList }) {
          const flatFieldsList = tools.flatArray(fieldsList)
          const fields = ['fc349']

          flatFieldsList.map((_field) => {
            if (fields.includes(_field.code)) {
              _field.table = {
                name: 'B0121_百年人寿团险理赔_疾病诊断',
                query: '疾病名称'
              }
            }
          })
        },
        disable05({ fieldsList, focusFieldsIndex }) {
          const fields = fieldsList[focusFieldsIndex];
          const fc327Field = fields.find(field => field.code === "fc327");
          if (fc327Field) fc327Field.disabled = true;
        },
        // 不同分块
        disable01({ fieldsList, focusFieldsIndex, codeValues = {} }) {
          const { fc364 } = codeValues
          const codesList = ['fc366', 'fc368', 'fc369', 'fc370']
          const fields = fieldsList[focusFieldsIndex]
          fields?.map(_field => {
            if (fc364 == '1')
              if (codesList.includes(_field.code)) {
                _field.disabled = true
              }
          })
        },

        disable367({ fieldsList, focusFieldsIndex }) {
          const fields = fieldsList[focusFieldsIndex];
          const fc367Field = fields.find(field => field.code === "fc367");
          if (fc367Field) fc367Field.disabled = true;
        },

        disable370({ fieldsList, focusFieldsIndex, codeValues = {} }) {
          const { fc365 } = codeValues
          const fields = fieldsList[focusFieldsIndex]

          if (fc365 == '1') {
            fields?.map(_field => {
              if (_field.code == 'fc370') {
                _field.disabled = true
              }
            })
          }

          if (fc365 == '2') {
            fields?.map(_field => {
              if (_field.code == 'fc369') {
                _field.disabled = true
              }
            })
          }
        },
      }
    },

    // 回车
    enter: {
      methods: {
        validate33And34: function ({ op, field, fieldsList, focusFieldsIndex, memoFields }) {
          if (field.code === "fc024") {
            const fields = fieldsList[focusFieldsIndex];

            const fc024Value = field.resultValue;
            const fc023Field = fields.find(field => field.code === "fc023");
            const fc022Field = fields.find(field => field.code === "fc022");

            if (fc024Value.includes("-")) {
              const values = fc024Value.split("-");

              field[`${op}Value`] = values[2];
              field.resultValue = values[2];
              _.set(memoFields, `${field.uniqueId}.value`, values[2]);

              fc023Field[`${op}Value`] = values[1];
              fc023Field.resultValue = values[1];
              _.set(memoFields, `${fc023Field.uniqueId}.value`, values[1]);

              fc022Field[`${op}Value`] = values[0];
              fc022Field.resultValue = values[0];
              _.set(memoFields, `${fc022Field.uniqueId}.value`, values[0]);
            }
          }
        },
        disable33({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== "fc319") return;

          const codes = ["fc335"];
          const fields = fieldsList[focusFieldsIndex];

          if (field.resultValue != 5) {
            fields.map(_field => {
              if (codes.includes(_field.code)) {
                _field.disabled = true;
              }
            });
          } else {
            fields.map(_field => {
              if (codes.includes(_field.code)) {
                _field.disabled = false;
              }
            });
          }
        },
        disable04({ op, field, fieldsList, focusFieldsIndex }) {
          if (field.code != "fc063") return;
          const fields = fieldsList[focusFieldsIndex];
          const fc367Field = fields.find(field => field.code === "fc367");
          if (fc367Field) {
            fc367Field[`${op}Value`] = field.resultValue;
            fc367Field.resultValue = field.resultValue;
            fc367Field.disabled = true;
          }
        },
        // fc105
        disable64({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== "fc105") return;

          const codes = [
            "fc016",
            "fc017",
            "fc018",
            "fc019",
            "fc020",
            "fc021",
            "fc022",
            "fc023",
            "fc024",
            "fc025",
            "fc026",
            "fc352"
          ];
          const validValues = ["A"];
          const fields = fieldsList[focusFieldsIndex];

          if (validValues.includes(field.resultValue)) {
            fields.map(_field => {
              if (codes.includes(_field.code)) {
                _field.disabled = true;
              }
            });
          } else {
            fields.map(_field => {
              if (codes.includes(_field.code)) {
                _field.disabled = false;
              }
            });
          }
        },
        disable65({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== "fc268") return;
          const validValues = ["A"];
          const fields = fieldsList[focusFieldsIndex];

          if (field.resultValue == "A") {
            fields.map(_field => {
              if (_field.code == "fc259") {
                _field.disabled = true;
              }
            });
          } else {
            fields.map(_field => {
              if (_field.code == "fc259") {
                _field.disabled = false;
              }
            });
          }
        },
        validate36({ field, fieldsList, focusFieldsIndex }) {
          const mapCodesList = new Map([
            ["fc260", ["fc098", "fc308"]],
            ["fc261", ["fc099", "fc309"]],
            ["fc262", ["fc100", "fc310"]],
            ["fc263", ["fc101", "fc311"]],
            ["fc264", ["fc102", "fc312"]],
            ["fc265", ["fc103", "fc313"]],
            ["fc266", ["fc104", "fc314"]],
            ["fc267", ["fc105", "fc315"]]
          ]);

          const codes = mapCodesList.get(field.code);

          if (codes) {
            const fields = fieldsList[focusFieldsIndex];
            const fiveField = fields.find(field => field.code === codes[0]);
            const sixField = fields.find(field => field.code === codes[1]);

            fiveField.disabled = false;
            sixField.disabled = false;

            if (field.resultValue === "1") {
              fiveField.disabled = true;
            } else if (field.resultValue === "4") {
              sixField.disabled = true;
            } else {
              fiveField.disabled = true;
              sixField.disabled = true;
            }
          }
        },
        validate37({ field, fieldsList, focusFieldsIndex }) {
          const mapCodesList = new Map([
            ["fc260", ["fc098", "fc308"]],
            ["fc261", ["fc099", "fc309"]],
            ["fc262", ["fc100", "fc310"]],
            ["fc263", ["fc101", "fc311"]],
            ["fc264", ["fc102", "fc312"]],
            ["fc265", ["fc103", "fc313"]],
            ["fc266", ["fc104", "fc314"]],
            ["fc267", ["fc105", "fc315"]]
          ]);

          const codes = mapCodesList.get(field.code);

          if (codes) {
            const fields = fieldsList[focusFieldsIndex];
            const fiveField = fields.find(field => field.code === codes[0]);
            const sixField = fields.find(field => field.code === codes[1]);

            fiveField.disabled = false;
            sixField.disabled = false;

            if (field.resultValue === "1") {
              fiveField.disabled = true;
            } else if (field.resultValue === "4") {
              sixField.disabled = true;
            } else {
              fiveField.disabled = true;
              sixField.disabled = true;
            }
          }
        },
        disable38({ field, fieldsList, focusFieldsIndex }) {
          if (field.code !== "fc082") return;

          const codes = [
            "fc074",
            "fc075",
            "fc076",
            "fc077",
            "fc078",
            "fc079",
            "fc080",
            "fc081",
            "fc110",
            "fc111",
            "fc112",
            "fc113",
            "fc114",
            "fc115",
            "fc116",
            "fc117"
          ];
          const fields = fieldsList[focusFieldsIndex];

          if (field.resultValue != "") {
            fields.map(_field => {
              if (codes.includes(_field.code)) {
                _field.disabled = true;
              }
            });
          }
        },
        validate21({ op, fieldsList, focusFieldsIndex, block }) {
          const fields = fieldsList[focusFieldsIndex];
          const fc062Field = fields.find(field => field.code == "fc062");
          // const fc066Field = fields.find(field => field.code == 'fc066')
          const fc068Field = fields.find(field => field.code == "fc068");
          const fc069Field = fields.find(field => field.code == "fc069");

          if (fc062Field?.[`${op}Value`] == 1) {
            fc069Field[`${op}Value`] = fc068Field[`${op}Value`];
            if (fc068Field[`${op}Value`] != "") fc069Field.resultValue = fc068Field[`${op}Value`];
            fc069Field.disabled = true;
          }
          if (fc062Field?.[`${op}Value`] == 2) {
            fc069Field.disabled = false;
          }
        },
        disable03({ op, field, fieldsList, focusFieldsIndex, codeValues = {} }) {
          if (field.code != "fc072") return;
          const { fc365 } = codeValues;
          const fields = fieldsList[focusFieldsIndex];
          const fc369Field = fields.find(field => field.code == "fc369");
          if (fc365 == "2") {
            fc369Field[`${op}Value`] = field.resultValue;
            fc369Field.resultValue = field.resultValue;
          }
        },
        disable02({ op, field, fieldsList, focusFieldsIndex }) {
          if (field.code != 'fc063') return
          const fields = fieldsList[focusFieldsIndex]
          fields?.map(_field => {
            if (_field.code == 'fc367') {
              _field[`${op}Value`] = field.resultValue
              _field.resultValue = field.resultValue
              _field.disabled = true
            }
          })
        },
      }
    },

    // 临时保存
    sessionSave: {
      methods: {
        disable33({ fieldsList, focusFieldsIndex }) {
          const codesList = [
            ["fc214", "fc224"],
            ["fc215", "fc225"],
            ["fc216", "fc226"],
            ["fc217", "fc227"],
            ["fc218", "fc228"],
            ["fc219", "fc229"],
            ["fc220", "fc230"],
            ["fc221", "fc231"],
            ["fc222", "fc232"],
            ["fc223", "fc233"],
            ["fc242", "fc252"],
            ["fc243", "fc253"],
            ["fc244", "fc254"],
            ["fc245", "fc255"],
            ["fc246", "fc256"]
          ];

          const col2Codes = [];

          codesList.map(codes => {
            col2Codes.push(codes[1]);
          });

          const fields = fieldsList[focusFieldsIndex];

          const focusField = fields.find(field => field.autofocus);
          const codeIndex = col2Codes.indexOf(focusField.code);

          if (codeIndex > -1) {
            const restCodes = [];
            let sliceIndex = -1;

            for (let codesIndex in codesList) {
              if (codesList[codesIndex].includes(focusField.code)) {
                sliceIndex = +codesIndex + 1;
                break;
              }
            }

            const restCodesList = codesList.slice(sliceIndex);

            restCodesList.map(codes => {
              restCodes.push(...codes);
            });

            const restFields = fields.slice(focusField.fieldIndex + 1);

            restFields?.map(restField => {
              if (restCodes.includes(restField.code)) {
                restField.disabled = true;
              }
            });
          }
        }
      }
    },

    // 提交前
    beforeSubmit: {
      methods: {
        validate39({ block, fieldsList }) {
          if (block.code !== "bc006") {
            return true;
          }

          for (let fields of fieldsList) {
            for (let field of fields) {
              if (field.resultValue) {
                return true;
              }
            }
          }

          return {
            errorMessage: "清单不能空白提交，请检查，如清单内容无法录入则录入一组数据后按F8提交."
          };
        },
        validate44: function ({ block, fieldsList, op }) {
          if (block.code !== "bc005") return true;
          if (op !== "op1") return true;

          const fields = fieldsList[0];
          const fc072Field = tools.find(fields, { code: "fc072" });
          const codes = [
            "fc224",
            "fc225",
            "fc226",
            "fc227",
            "fc228",
            "fc229",
            "fc230",
            "fc231",
            "fc232",
            "fc233",
            "fc252",
            "fc253",
            "fc254",
            "fc255",
            "fc256"
          ];
          if (!fc072Field || fc072Field.resultValue === "?") {
            return true;
          }

          for (let field of fields) {
            if (codes.includes(field.code)) {
              if (field.resultValue === "?") {
                return true;
              }
            }
          }

          // let count = 0
          let count = new BigNumber(0);

          for (let field of fields) {
            if (codes.includes(field.code)) {
              const resultValue = +field.resultValue || 0;
              // count += resultValue
              count = count.plus(resultValue);
            }
          }

          // const fc072Value = +fc072Field.resultValue

          // const diff = fc072Value - count

          const fc072Value = new BigNumber(+fc072Field.resultValue);

          const diff = fc072Value.minus(count).toString();

          if (diff != 0) {
            return {
              popup: 'confirm',
              errorMessage: `明细金额与总金额不一致，差额为${diff}，请确认并修改!`
            };
          }
          return true;
        },
        validate53: function ({ fieldsList }) {
          const fields = fieldsList[0];
          const codesListMap = new Map([
            ["fc082", "fc090", "fc234", "fc260"],
            ["fc083", "fc091", "fc235", "fc261"],
            ["fc084", "fc092", "fc236", "fc262"],
            ["fc085", "fc093", "fc237", "fc263"],
            ["fc086", "fc094", "fc238", "fc264"],
            ["fc087", "fc095", "fc239", "fc265"],
            ["fc088", "fc096", "fc240", "fc266"],
            ["fc089", "fc097", "fc241", "fc267"]
          ]);

          for (let field of fields) {
            const { code, resultValue } = field;
            const targetCode = codesListMap.get(code);

            if (targetCode && resultValue) {
              const targetField = tools.find(fields, { code: targetCode });

              if (!targetField?.resultValue) {
                return {
                  errorMessage: "清单内容录入遗漏，请检查!"
                };
              }
            }
          }
          return true;
        },

        // CSB0121RC0123000
        validate33({ fieldsList }) {
          const fields = fieldsList[0]

          let codeMaps = new Map([
            ['fc214', 'fc224'],
            ['fc215', 'fc225'],
            ['fc216', 'fc226'],
            ['fc217', 'fc227'],
            ['fc218', 'fc228'],
            ['fc219', 'fc229'],
            ['fc220', 'fc230'],
            ['fc221', 'fc231'],
            ['fc222', 'fc232'],
            ['fc223', 'fc233'],
            ['fc242', 'fc252'],
            ['fc243', 'fc253'],
            ['fc244', 'fc254'],
            ['fc245', 'fc255'],
            ['fc246', 'fc256'],
          ])

          let arr = ['fc214', 'fc215', 'fc216', 'fc217', 'fc218', 'fc219', 'fc220', 'fc221', 'fc222', 'fc223', 'fc242', 'fc243', 'fc244', 'fc245', 'fc246']

          for (let field of fields) {

            if (arr.includes(field.code) && field.resultValue == '') {

              for (let _field of fields) {
                if (_field.code == codeMaps.get(field.code) && _field.resultValue == '') break
                else if (_field.code == codeMaps.get(field.code) && _field.resultValue != '') {
                  if (_field.resultValue == '?' || field.resultValue == '？') break
                  return {
                    errorMessage: `${field.code}_${field.name}不能为空， 请核查`
                  }
                }
              }
            }

            if (arr.includes(field.code) && field.resultValue != '') {
              if (field.resultValue == '?' || field.resultValue == '？') break

              for (let _field of fields) {
                if (_field.code == codeMaps.get(field.code) && _field.resultValue == '') {
                  return {
                    errorMessage: `${_field.code}_${_field.name}不能为空， 请核查`
                  }
                }
              }
            }
          }

          return true
        }
      }
    }
  }
};

export default B0121;
