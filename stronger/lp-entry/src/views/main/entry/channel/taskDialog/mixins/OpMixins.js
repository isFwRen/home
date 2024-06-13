import localForage from 'localforage'
import { tools, sessionStorage } from 'vue-rocket'
import { LP2 } from '@/api/syncPouchDB'
import { tools as lpTools } from '@/libs/util'
import { mapRules, validators } from '../components/opTextField/rules'
import { toastedOptions } from '../../cells'
import { MessageBox } from 'element-ui';
import { reqConst } from '@/api/reqConst'
import moment from 'moment'
import keyMap from './keycode'

const { baseURLApi } = lpTools.baseURL()

// // 与字段配置的"录入工序"对应
// const inputProcessMap = new Map([
//   [2, 'op1'],
//   [3, 'op2']
// ])
sessionStorage.set('select', [])
export default {
  data() {
    return {
      fileUrl: `${baseURLApi}files/`,

      // 返回重录(1：上一单，2：前一单，以此类推...)
      prevNums: 0,

      oldTime: void 0,

      fieldName: void 0,
    }
  },

  computed: {
    // 字段名
    computedLabel() {
      const env = process.env.NODE_ENV

      if (env === 'development' || env === 'test') {
        return function ({ field, fieldIndex, fieldsIndex }) {
          const { code, name, allowForce } = field

          return `(${code}_${fieldsIndex}-${fieldIndex})${name}${allowForce === false ? ' - 不允许强过' : ''}`
        }
      }

      return function ({ field }) {
        return field.name
      }
    },

    // 校验规则
    computedRules() {
      return function (rules) {
        return tools.isYummy(rules) ? tools.deepClone(rules) : []
      }
    },

    // 字段提示(特殊校验)
    computedSVHints() {
      if (this.op === 'op0') {
        return this.specificProject?.op0?.hints
      }

      return this.specificProject?.op1op2opq?.hints
    },

    // 字段校验(特殊校验)
    computedSVValidations() {
      if (this.op === 'op0') {
        return this.specificProject?.op0?.rules
      }

      return this.specificProject?.op1op2opq?.rules
    }
  },

  // created() {
  //   this.reloadConstants()
  // },

  mounted() {
    window.addEventListener('keydown', this.fuckOpShortcut)
  },

  beforeDestroy() {
    window.removeEventListener('keydown', this.fuckOpShortcut)
  },

  methods: {
    // 重新获取本地常量表
    async reloadConstants() {
      if (!window['constantsDB']) {
        const forage = await localForage.getItem(LP2)
        window['constantsDB'] = forage || {}
      }

      console.log(window['constantsDB'])
    },

    // 获取客户端常量表
    async reqConstants() {
      let result = await reqConst({
        url: `/sys-const/info-list/${this.proCode}`,
        method: "GET",
      })
      let proList = result.list
      window.sessionStorage.setItem('proList', JSON.stringify(proList))
    },

    // 生成校验规则
    setValidateRules({ field, configField }) {
      // 只有[合法field]才需生成校验规则
      if (!field.show || field.disabled) return

      // console.log(configField.name, configField)

      const rules = []

      // Normal validations
      // if(inputProcessMap.get(configField.inputProcess) === this.op) {
      // 从 field 的属性中获取校验规则
      for (let key in configField) {
        if (!mapRules[key] || !tools.isYummy(configField[key])) continue

        if (key === 'specChar') {
          if (configField.specChar) {
            const values = configField.specChar.split(';').filter(f => f)

            rules.push({
              key: mapRules['specChar'],
              rule: values,
              message: `只有${configField.specChar}为可通过字符.`
            })
          }
        }
        else if (key === 'fixValue') {
          if (configField.fixValue) {
            const values = configField.fixValue.split(';').filter(f => f)

            rules.push({
              key: mapRules['fixValue'],
              rule: values,
              message: `必须为${configField.fixValue}中的一项.`
            })
          }
        }
        else {
          rules.push({ key: mapRules[key], rule: configField[key] })
        }
      }

      // 从 field 的 validations 属性中获取校验规则
      configField?.validations?.map(value => {
        if (mapRules[value]) {
          rules.push({ key: mapRules[value], rule: 'NO' })
        }
      })
      // }

      // Specific validations
      if (tools.isYummy(field.includes)) {
        rules.push({
          key: mapRules['includes'],
          rule: field.includes.items,
          message: field.includes.message
        })
      }

      // 问题件
      if (this.op === 'opq') {
        rules.push(
          {
            key: 'loose_excluded',
            rule: ['?'],
            message: '请检查录入数据，确认没有问题，请强制通过!'
          },

          {
            key: 'included',
            rule: [field.op1Value, field.op2Value],
            message: '请检查录入数据，问题件与1、2码录入数据不一致!'
          }
        )
      }

      // console.log(field.name, rules)

      return rules
    },

    // 校验当前字段
    validateField({ field }) {
      if (field.force) return true

      const index = tools.findIndex(field?.rules, { key: 'free_value' })
      if (field?.rules && field?.rules[index]?.rule.includes(field.resultValue)) return true
      const { nValidations, sValidations, nvArgs, svArgs } = this.$refs[field.uniqueId][0]
      let common
      let special
      // 普通校验
      const nVerified = nValidations.every(validation => {
        const result = validators[validation['key']]({ ...nvArgs(validation), effectValidations: validation.effectValidations })
        common = result
        return result === true
      })

      // 特殊校验
      const sVerified = sValidations.every(validation => {
        const result = validation['func']({ ...svArgs(), effectValidations: validation.effectValidations })
        return result === true
      })

      if (!nVerified || !sVerified) {
        console.log({
          '普通校验': nVerified ? '通过' : '不通过',
          '业务校验': sVerified ? '通过' : '不通过',
          '普通信息': common,
          '特殊信息': special
        })
      }

      return nVerified && sVerified
    },

    // F8 F4校验所有字段
    async validateFields(fields, key) {
      console.log("所有字段", fields);
      let result = true
      for (let field of fields) {
        // const { nValidations, nvArgs } = this.$refs[field.uniqueId][0]
        const { nValidations, sValidations, nvArgs, svArgs } = field.checkMethod
        let common1
        let common2
        // 普通校验
        const nVerified = nValidations.every(validation => {
          const result = validators[validation['key']]({ ...nvArgs(validation), effectValidations: validation.effectValidations })
          common1 = result
          return result === true
        })

        // 特殊校验
        const sVerified = sValidations.every(validation => {
          const result = validation['func']({ ...svArgs(), effectValidations: validation.effectValidations })
          common2 = result
          return result === true
        })

        if (field.name.includes('医院') || field.name.includes('医疗机构')) {
          common1 = true
          common2 = true
          result = true
        }
        if ((typeof common1) == 'string') {
          let page = Number(field.uniqueId.slice(0, 1)) + 1
          let fieldName = field.uniqueId.slice(2, 5)

          if (sessionStorage.get('isApp')?.isApp === 'true') {
            await MessageBox.confirm(`禁止${key}，请修改或删除第${page}张图的${field.code}_${fieldName}${field.name}的内容, 禁止原因：内容${field.resultValue}不满足条件${common1}`, '请检查', {
              type: 'warning',
              confirmButtonText: '确定',
              showCancelButton: false,
              showClose: false,
            }).then(() => {
              result = false
            })
          } else {
            result = alert(`禁止${key}，请修改或删除第${page}张图的${field.code}_${fieldName}${field.name}的内容, 禁止原因：内容${field.resultValue}不满足条件${common1},`)
          }
        }

        if ((typeof common2) == 'string') {
          let page = Number(field.uniqueId.slice(0, 1)) + 1
          let fieldName = field.uniqueId.slice(2, 5)

          if (sessionStorage.get('isApp')?.isApp === 'true') {
            await MessageBox.confirm(`禁止${key}，请修改或删除第${page}张图的${field.code}_${fieldName}${field.name}的内容, 禁止原因：内容${field.resultValue}不满足条件${common2}`, '请检查', {
              type: 'warning',
              confirmButtonText: '确定',
              showCancelButton: false,
              showClose: false,
            }).then(() => {
              result = false
            })
          } else {
            result = alert(`禁止${key}，请修改或删除第${page}张图的${field.code}_${fieldName}${field.name}的内容, 禁止原因：内容${field.resultValue}不满足条件${common2},`)
          }
        }

      }
      return result
    },

    async validateFieldss(fields, key) {
      console.log("所有字段", fields);
      let result = true
      for (let field of fields) {
        // const { nValidations, nvArgs } = this.$refs[field.uniqueId][0]
        const { nValidations, sValidations, nvArgs, svArgs } = field.checkMethod
        let common1
        let common2
        // 普通校验
        const nVerified = nValidations.every(validation => {
          const result = validators[validation['key']]({ ...nvArgs(validation), effectValidations: validation.effectValidations })
          common1 = result
          return result === true
        })

        // 特殊校验
        const sVerified = sValidations.every(validation => {
          const result = validation['func']({ ...svArgs(), effectValidations: validation.effectValidations })
          common2 = result
          return result === true
        })

        if (field.name.includes('医院') || field.name.includes('医疗机构')) {
          common1 = true
          common2 = true
          result = true
        }

        if ((typeof common1) == 'string' || (typeof common2) == 'string') result = false
      }
      return result
    },

    // 问题件录入的值在常量表，即使与一码、二码不匹配也可以强过
    fuckAllowForce({ field }) {
      if (this.op !== 'opq') return
      if (!tools.isYummy(field.items)) return

      if (field.items.includes(field.resultValue)) {
        field.allowForce = true
      }
    },

    // 快捷键
    async fuckOpShortcut(event) {
      const { keyCode } = event || window.event

      switch (keyCode) {
        // 清空当前输入框(ESC)
        case 27:
          if (this.op === 'op0') {
            const field = this.fieldsObject[this.thumbIndex].fieldsList[this.focusFieldsIndex][this.focusFieldIndex]

            field.op0Value = ''
            field.resultValue = ''
          }
          else {
            const field = this.fieldsList[this.focusFieldsIndex][this.focusFieldIndex]

            field[`${this.op}Value`] = ''
            field.resultValue = ''

            this.fieldsList = [...this.fieldsList]
          }
          break;

        // 返回修改(F3)
        case 114:
          event.preventDefault()

          // const prevNums = ++this.prevNums
          // this.getTask({
          //   status: 'modify',
          //   prevNums
          // })
          if (!this.$store.state['entry/task'].f3state) {
            this.toasted.warning('没有要修改的初审内容，请不要重复操作F3', { duration: 3000 })
          }
          if (this.$store.state['entry/task'].f3state) {
            const prevNums = ++this.prevNums
            this.getTask({
              status: 'modify',
              prevNums
            })
            this.$store.commit('UPDATE_F3STATE', false)
          }
          break;

        // 提交(F8)
        case 119:
          event.preventDefault()
          // 开发环境默认可以通过F8提交
          if (process.env.NODE_ENV === 'development') {
            console.log(process.env.NODE_ENV);
            // if (this.times != 0 && this.block.name == '清单' && this.start && this.times != 0) {

            //   if (sessionStorage.get('isApp')?.isApp === 'true') {
            //     await MessageBox.alert(`请认真核对清单内容是否与影像内容一致---${this.times}秒后可以提交`, '请检查', {
            //       type: 'warning',
            //       confirmButtonText: '确定',
            //       showClose: false,
            //     })
            //   } else {
            //     alert(`请认真核对清单内容是否与影像内容一致---${this.times}秒后可以提交`)
            //   }

            // } else {
            if (this.op == 'op0') {
              setTimeout(() => {
                this.fEightSubmit()
              }, 300)
            } else this.fEightSubmit()

            // }
          }
          else {
            if (this.block.fEight) {
              // console.log(this.block.fEight);
              // console.log('不是development');
              // console.log(this.bill);
              // if (this.times != 0 && this.block.name == '清单' && this.start && this.times != 0) {

              // if (sessionStorage.get('isApp')?.isApp === 'true') {
              //   await MessageBox.alert(`请认真核对清单内容是否与影像内容一致---${this.times}秒后可以提交`, '请检查', {
              //     type: 'warning',
              //     confirmButtonText: '确定',
              //     showClose: false,
              //   })
              // } else {
              //   alert(`请认真核对清单内容是否与影像内容一致---${this.times}秒后可以提交`)
              // }

              // } else {
              if (this.op == 'op0') {
                setTimeout(() => {
                  this.fEightSubmit()
                }, 300)
              } else this.fEightSubmit()
              // }
            }
            // else {
            //   this.toasted.warning('当前任务不可通过F8提交!', toastedOptions)
            // }
          }

          // console.warn({ fEight: this.block.fEight })
          break;
      }

      if (keyCode != 13 && keyCode != 115 && keyCode != 119 && keyMap.get(keyCode)) {
        let operate = '按键:' + keyMap.get(keyCode) + ',' + moment(new Date).format('YYYY/MM/DD HH:mm:ss') + ',' + '第' + this.$store.state['recordKey'].page + '页'
        // this.recordKey.push(operate)
        this.$store.commit('UPDATE_KEY', operate)
      }
    },

    // 限制提交次数
    limitSubmitTask() {
      if (!this.oldTime) {
        this.submitTask()
        this.oldTime = Date.now()
      } else {
        if (Date.now() - this.oldTime > 1000) {
          this.submitTask()
          this.oldTime = Date.now()
        }
      }
    },

    // 打开字段规则弹框
    handleFieldRulesDialog({ name }) {
      this.fieldName = name
      this.$refs.fieldRules.onOpen()
    },

    // 阻止表单默认事件
    preventForm(event) {
      event.preventDefault()
    }
  }
}