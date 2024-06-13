import { tools as lpTools } from '@/libs/util'

export default {
  props: {
    firstDiffIndex: {
      type: Number,
      default: -1
    },
    lastDiffIndex: {
      type: Number,
      default: -1
    }
  },

  data() {
    return {
      focusStrIndex: -1
    }
  },

  watch: {
    // 找到问题件第一个问号
    firstDiffIndex: {
      handler(index) {
        this.$nextTick(() => {
          this.focusStrIndex = index
          if (this.autofocus && this.id && this.focusStrIndex > -1) {
            const el = document.querySelector(`#input_${this.id}`)
            lpTools.setCursorPosition(el, this.firstDiffIndex, this.lastDiffIndex)
          }
        })
      },
      immediate: true
    }
  }
}