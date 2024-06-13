let thumbsNode = null

export default {
  data() {
    return {
      thumbsDeadLine: 300
    }
  },

  watch: {
    'bill.pictures': {
      handler(pictures) {
        if(!pictures || !pictures.length) return

        this.$nextTick(() => {
          thumbsNode = document.querySelector('#op0ThumbsNode')
        })
      },
      immediate: true
    }
  },

  methods: {
    // 向上向下滚动
    scrollThumbsUpDn({ id, thumbIndex }) {
      this.$nextTick(() => {
        const thumb = document.getElementById(id)

        if(thumb?.offsetTop >= 300) {
          thumbsNode.scrollTop = thumb.offsetTop - this.thumbsDeadLine
        }
        else {
          thumbsNode.scrollTop = 0
        }
      })
    } 
  }
}