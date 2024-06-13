import { tools, sessionStorage } from 'vue-rocket'

export default {
  props: {
    items: {
      type: Array,
      required: false
    },
    sameFieldValue: {
      type: Object | Array,
      required: false
    }
  },

  data() {
    return {
      typing: false,
      dropdownItems: []
    }
  },

  methods: {
    handleHoverListItem({ item, index }) {
      if (this.typing) {
        this.$emit('dropdown', {
          index: index,
          value: item
        })

        this.$refs[`input_${this.id}`].focus()
      }
    },

    handleSelectListItem({ item, index }) {
      if (this.typing) {
        this.$emit('dropdown', {
          index: index,
          value: item
        })
        let selectArr = sessionStorage.get('select')
        if (!selectArr.includes(item)) {
          selectArr.push(item)
        }
        sessionStorage.set('select', selectArr)
        this.$refs[`input_${this.id}`].focus()
        this.typing = false
      }
    }
  },

  watch: {
    items: {
      handler(items) {
        if (items.length) {
          this.dropdownItems = tools.deepClone(items)
        }
        else {
          this.dropdownItems = []
        }
      },
      immediate: true
    }
  }
}