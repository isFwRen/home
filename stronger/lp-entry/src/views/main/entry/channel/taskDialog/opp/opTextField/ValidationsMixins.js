import { tools } from 'vue-rocket'
import { validators } from './rules'

/**
 * N/n => normal
 * S/s => special
 * V/v => validation
 * H/h => hint
 */

export default {
  props: {
    bill: {
      type: Object,
      default: () => ({})
    },

    block: {
      type: Object,
      default: () => ({})
    },

    clearNVs: {
      type: Boolean,
      default: false
    },

    clearVs: {
      type: Boolean,
      default: false
    },

    includes: {
      type: Array,
      default: () => ([])
    },

    svHints: {
      type: Array,
      required: false
    },

    svValidations: {
      type: Array,
      required: false
    },

    validations: {
      type: Array,
      required: false
    }
  },

  data() {
    return {
      billObject: {},
      blockObject: {},

      error: false,
      errorMessages: [],

      // 普通校验 - 对象
      nValidations: [],

      // 特殊校验 - 对象
      sValidations: [],

      // 校验规则 - 函数
      rules: [],
      // 临时保存校验规则 - 函数
      sessionRules: [],

      // 跳过校验
      freeValueRule: null,

      hints: [],

      // 普通校验数量
      nvCount: 0
    }
  },

  watch: {
    // 忽略普通校验规则
    clearNVs: {
      handler(clearNVs) {
        if (clearNVs) {
          this.rules.splice(0, this.nvCount)
        }
        // console.log({ [this.label]: clearNVs })
      }
    },

    // 忽略校验规则
    clearVs: {
      handler(clearVs) {
        if (clearVs) {
          this.rules.splice(0)
        }
        // console.log({ [this.label]: clearVs })
      }
    }
  },

  created() {
    const { force } = this.field

    this.billObject = { ...this.bill }

    this.blockObject = { ...this.block }

    this.rules = []

    this.hints = []

    !force && this.setNVs()
    // !force && this.setSVs()
    this.setSVs()
    this.setSHs()

    // 禁用/隐藏不需要校验
    if (this.field.disabled || !this.field.show) {
      this.rules = []
    }

    this.sessionRules = tools.deepClone(this.rules)
  },

  methods: {
    // 设置普通校验
    setNVs() {
      this.nvCount = 0

      if (!tools.isYummy(this.validations)) {
        return
      }

      // this.removeShowLengthRule()

      // this.removeSpacedRule()

      this.freeValueRule = this.getFreeValueRule()

      const effectValidations = {
        freeValue: this.freeValueRule
      }

      for (let validation of this.validations) {
        this.nValidations.push({ ...validation, effectValidations })

        // 填入校验规则
        this.rules.push(() =>
          validators[validation['key']]({ ...this.nvArgs(validation), effectValidations })
        )

        ++this.nvCount
      }

      this.field.effectValidations = effectValidations

      // console.log(this.field.name, this.nValidations)
    },

    // 设置特殊校验
    setSVs() {
      if (!tools.isYummy(this.svValidations)) return

      const effectValidations = { freeValue: this.freeValueRule }

      for (let svValidation of this.svValidations) {
        // 过滤不需要校验的字段
        if (!svValidation['fields']?.includes(this.field.code)) continue

        // console.log(this.field.code, svValidation)

        for (let func in svValidation) {
          if (!tools.isFunction(svValidation[func])) continue

          this.sValidations.push({ func: svValidation[func], effectValidations })

          // 填入校验规则
          if (svValidation.index === 0) {
            this.rules.unshift(() =>
              svValidation[func]({ ...this.svArgs(), effectValidations })
            )
          }
          else {
            this.rules.push(() =>
              svValidation[func]({ ...this.svArgs(), effectValidations })
            )
          }
        }
      }
    },

    // 设置特殊提示文本
    setSHs() {
      if (!tools.isYummy(this.svHints)) return

      for (let svHint of this.svHints) {
        if (!svHint['fields']?.includes(this.field.code)) continue

        for (let funcName in svHint) {
          if (!tools.isFunction(svHint[funcName])) continue

          const message = svHint[funcName](this.shArgs())

          if (tools.isString(message)) {
            this.hints.push(message)
          }
        }
      }
    },

    // 回车/输入文本后设置 field 状态
    validateAfterEnterInput() {
      const { code } = this.field

      let svRules = []

      this.freeValueRule = this.getFreeValueRule()

      const effectValidations = {
        freeValue: this.freeValueRule
      }

      // 普通校验
      for (let i = 0; i < this.validations.length; i += 1) {
        const { rule, message } = this.validations[i]

        const result = this.rules[i]?.({ effectValidations, value: this.value, rule, message })

        if (result !== true) {
          this.error = true
          this.errorMessages = result
          return
        }

        this.error = false
        this.errorMessages = []

        svRules = this.rules.slice(i + 1)
      }

      // 特殊校验
      for (let svValidation of this.svValidations) {
        if (!svValidation['fields']?.includes(code)) continue

        for (let rule of svRules) {
          const result = rule(this.svArgs(), effectValidations)

          if (result !== true) {
            this.error = true
            this.errorMessages = result
            return
          }

          this.error = false
          this.errorMessages = []
        }
      }
    },

    // 普通校验所需参数
    nvArgs(validation) {
      return {
        field: this.field,
        fieldsList: this.fieldsList,
        label: this.label,
        message: validation['message'],
        op: this.op,
        rule: validation['rule'],
        value: this.value
      }
    },

    // 特殊校验所需参数
    svArgs() {
      return {
        bill: this.billObject,
        block: this.blockObject,
        field: this.field,
        fieldsIndex: this.fieldsIndex,
        fieldsList: this.fieldsList,
        fieldsObject: this.fieldsObject,
        includes: this.includes,
        items: this.items,
        sameFieldValue:this.sameFieldValue,
        label: this.label,
        op: this.op,
        thumbIndex: this.thumbIndex,
        value: this.value
      }
    },

    // 特殊提示文本所需参数
    shArgs() {
      return {
        bill: this.billObject,
        block: this.blockObject,
        field: this.field,
        fieldsIndex: this.fieldsIndex,
        fieldsList: this.fieldsList,
        label: this.label,
        op: this.op,
        value: this.value,
        items: this.items,
      }
    },

    // 是否需要显示 value 的长度
    removeShowLengthRule() {
      const index = tools.findIndex(this.validations, { key: 'show_length' })

      if (index !== -1) {
        this.showLength = true
        this.validations.splice(index, 1)
      }

      this.showLength = false
    },

    // 是否允许 value 包含空格
    removeSpacedRule() {
      const index = tools.findIndex(this.validations, { key: 'spaced' })

      if (index !== -1) {
        this.allowSpace = true
        this.validations.splice(index, 1)
      }

      this.allowSpace = false
    },

    // 跳过校验的字符
    getFreeValueRule() {
      const index = tools.findIndex(this.validations, { key: 'free_value' })

      if (index !== -1) {
        return this.validations.splice(index, 1)[0]
      }

      return null
    }
  }
}