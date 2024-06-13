import moment from 'moment'
export default {
  mounted() {
    this.$refs.drawImage?.eventCut()

    window.addEventListener('keydown', this.canvasFuckEvents)
  },

  beforeDestroy() {
    window.removeEventListener('keydown', this.canvasFuckEvents)
  },

  methods: {
    // 通过按键选择绘图方式或者提交
    canvasFuckEvents(event) {
      event = event || window.event
      const { ctrlKey, shiftKey, keyCode } = event

      if (keyCode === 46) {
        // 清除当前激活对象(Delete)
        this.$refs.drawImage.eventClearActivatedCtx()
      }
      else if (keyCode === 81) {
        // 切图(Q)
        if (ctrlKey) {
          event.preventDefault()
          this.$refs.drawImage.eventCut()
        }
      }
      else if (keyCode === 88) {
        // 框图(x)
        if (ctrlKey) {
          event.preventDefault()
          this.$refs.drawImage.eventRect()
          this.$store.commit('UPDATE_KEY', '按键:框图' + ',' + moment(new Date).format('YYYY/MM/DD HH:mm:ss') + ',' + '第' + this.$store.state['recordKey'].page + '页')
        }
      }
      else if (keyCode === 69) {
        // 文字(E)
        if (ctrlKey) {
          event.preventDefault()
          this.$refs.drawImage.eventText()
        }
      }
      else if (keyCode === 76) {
        // 左微旋转(L)
        if (shiftKey && ctrlKey) {
          event.preventDefault()
          this.$refs.drawImage.eventRotateMinLeft()
          this.$store.commit('UPDATE_KEY', '按键:左微旋转' + ',' + moment(new Date).format('YYYY/MM/DD HH:mm:ss') + ',' + '第' + this.$store.state['recordKey'].page + '页')
          return
        }

        // 左旋转(L)
        if (ctrlKey) {
          event.preventDefault()
          this.$refs.drawImage.eventRotateLeft()
          this.$store.commit('UPDATE_KEY', '按键:左旋转' + ',' + moment(new Date).format('YYYY/MM/DD HH:mm:ss') + ',' + '第' + this.$store.state['recordKey'].page + '页')
        }
      }
      else if (keyCode === 82) {
        // 右微旋转(R)
        if (shiftKey && ctrlKey) {
          event.preventDefault()
          this.$refs.drawImage.eventRotateMinRight()
          this.$store.commit('UPDATE_KEY', '按键:右微旋转' + ',' + moment(new Date).format('YYYY/MM/DD HH:mm:ss') + ',' + '第' + this.$store.state['recordKey'].page + '页')
          return
        }

        // 右旋转(R)
        if (ctrlKey) {
          event.preventDefault()
          this.$refs.drawImage.eventRotateRight()
          this.$store.commit('UPDATE_KEY', '按键:右旋转' + ',' + moment(new Date).format('YYYY/MM/DD HH:mm:ss') + ',' + '第' + this.$store.state['recordKey'].page + '页')
        }
      }
      else if (keyCode === 90) {
        // 清除(Z)
        if (ctrlKey) {
          event.preventDefault()
          this.$refs.drawImage.eventClear()
        }
      }
      else if (keyCode === 117) {
        // 还原(F6)
        event.preventDefault()
        this.handleInitImage()
      }
    }
  }
}