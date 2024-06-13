import { localStorage } from 'vue-rocket'

const [titleHeight, range] = [24, 15]

let [fieldsNode, tLoadingNode, bLoadingNode] = [void 0, void 0, void 0]

export default {
  data() {
    return {
      upperShowRange: range,
      lowerShowRange: range,
      fieldHeight: 94,
      deadLine: 0
    }
  },

  created() {
    const { height } = localStorage.get('viewport')

    if (this.op === 'opq') {
      this.fieldHeight = 138
    }

    this.deadLine = this.fieldHeight * 2 + titleHeight

    this.deadLine = height * .45
  },

  watch: {
    'fieldsList.length': {
      handler(length) {
        if (this.op === 'op0' || length < 1) return

        this.getNodes()
      },
      immediate: true
    }
  },

  methods: {
    setShowRange() {
      const fields = this.fieldsList[this.focusFieldsIndex]
      const fieldIndex = this.focusFieldIndex

      if (!fields) return

      const upperFields = fields.slice(0, fieldIndex)
      const lowerFields = fields.slice(fieldIndex)
      const [upperShowFields, lowerShowFields] = [[], []]

      // upper
      for (let _field of upperFields) {
        if (_field.show) {
          upperShowFields.push(_field)

          if (upperShowFields.length === range) {
            break
          }
        }
      }

      // lower
      for (let _field of lowerFields) {
        if (_field.show) {
          lowerShowFields.unshift(_field)

          if (lowerShowFields.length === range) {
            break
          }
        }
      }

      const upperField = upperShowFields[0]
      const lowerField = lowerShowFields[0]

      if (upperField) {
        this.upperShowRange = fieldIndex - upperField.fieldIndex
      }
      else {
        this.upperShowRange = 0
      }

      if (lowerField) {
        this.lowerShowRange = lowerField.fieldIndex
      }
      else {
        this.lowerShowRange = fields.length
      }

      // console.log(this.focusFieldIndex, { lowerShowRange: this.lowerShowRange, upperShowRange: this.upperShowRange })
    },

    observer() {
      const options = {
        root: fieldsNode,
        threshold: 0
      }

      const ob = new IntersectionObserver((entries) => {
        const [entry] = entries

        const nodeValue = entry.target.attributes['data-value'].nodeValue

        if (nodeValue === 'tLoading') {
          if (entry.isIntersecting) {
            this.upperShowRange += 10
            fieldsNode.scrollTop = 50
          }
        }
        else {
          if (entry.isIntersecting) {
            this.lowerShowRange += 10
          }
        }
      }, options)

      ob.observe(tLoadingNode)
      ob.observe(bLoadingNode)
    },

    // 向上向下滚动
    scrollUpDn({ field }) {
      this.$nextTick(() => {
        const timer = setTimeout(() => {
          const input = document.getElementById(field?.uniqueId)

          if (input?.offsetTop >= this.deadLine) {
            // 兼容客户端
            if (document.documentElement.clientHeight < 700) fieldsNode.scrollTop = input.offsetTop - 200
            else fieldsNode.scrollTop = input.offsetTop - 433.5
            // fieldsNode.scrollTop = input.offsetTop - this.deadLine
          }

          clearTimeout(timer)
        }, 25)
      })
    },

    // 返回顶部(暂未使用)
    scrollToTop() {
      const timer = setTimeout(() => {
        if (fieldsNode) {
          fieldsNode.scrollTop = 0
        }

        clearTimeout(timer)
      }, 25)
    },

    getNodes() {
      this.$nextTick(() => {
        fieldsNode = document.querySelector(`#${this.op}FieldsNode`)
        tLoadingNode = document.querySelector('.t-loading')
        bLoadingNode = document.querySelector('.b-loading')

        if (fieldsNode && tLoadingNode && bLoadingNode) {
          this.observer()
        }
      })
    }
  }
}