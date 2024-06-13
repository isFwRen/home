import _ from 'lodash'

export default {
  methods: {
    // 下拉选中回车
    onDropdownField({ value, index }, field) {
      if (index > -1) {
        let val = value
        console.log(val);
        field[`op1Value`] = val
        // field.resultValue = val
        field.autofocus = true

        if (this.op === 'op0') {
          _.set(this.memoFields, `${field.uniqueId}.value`, val)
          this.fieldsObject = { ...this.fieldsObject }
        }
        else {
          this.fieldsList = [...this.fieldsList]
          console.log(this.fieldsList);
        }
      }
    },

    // 下拉按上键
    onDropdownUpField() {

    }
  }
}
