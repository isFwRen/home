import { sessionStorage } from 'vue-rocket'
import { tools as lpTools } from '@/libs/util'

const { baseURLApi } = lpTools.baseURL()

sessionStorage.set('select', [])
export default {
  data() {
    return {
      fileUrl: `https://www.i-confluence.com:12345/api/`,

      oldTime: void 0,

      fieldName: void 0,
    }
  },

  mounted() {
    window.addEventListener('keydown', this.fuckOpShortcut)
  },

  beforeDestroy() {
    window.removeEventListener('keydown', this.fuckOpShortcut)
  },

  methods: {
    // 快捷键
    async fuckOpShortcut(event) {
      const { keyCode } = event || window.event

      switch (keyCode) {
        // 提交(F8)
        case 119:
          event.preventDefault()
          // 开发环境默认可以通过F8提交
          if (process.env.NODE_ENV === 'development') {
            this.fEightSubmit()
          }
          else {
            this.fEightSubmit()
          }
          break;
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

    // 阻止表单默认事件
    preventForm(event) {
      event.preventDefault()
    }
  }
}