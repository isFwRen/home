const commonFunctionOP12q = {
  methods: {
    // 问题件将屏蔽字段的问号清空
    validate99999: function ({ fieldsList, op }) {
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

export default commonFunctionOP12q