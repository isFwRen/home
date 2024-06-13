<template>
  <div class="z-flex justify-center lp-desensitization">
    <div 
      v-if="update" 
      id="screenshot"
      class="screenshot" 
    >
      <img 
        id="canvasImg" 
        :src="imgAddress" 
        :width="imgWidth" 
        height="auto"
      >
      <div class="canvas">
        <canvas id="cutCanvas"></canvas>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    name: 'LPDesensitization',

    props: {
      defaultStartXY: {
        type: Object,
        required: false
      },

      defaultRectWH: {
        type: Object,
        required: false
      },

      imgAddress: {
        type: String,
        required: true
      },

      imgWidth: {
        type: [Number, String],
        default: 500
      },

      update: {
        type: Boolean,
        default: false
      }
    },

    data() {
      return {
        width: 0,
        height: 0,
        startXY: {
          x: null,
          y: null
        },
        endXY: {
          x: null,
          y: null
        },
        // rectX: null,
        // rectY: null,
        rectW: null,
        rectH: null,

        cutCanvas: undefined,
        cutCtx: undefined,

        originImg: {
          w: null,
          h: null
        },

        localImg: {
          w: null,
          h: null
        }
      }
    },

    mounted() {
      !this.update && this._setDefaultCanvas()
    },

    methods: {
      /**
       * 初始化
       */ 
      init() {
        if(this.defaultStartXY) {
          this.startXY = { ...this.defaultStartXY }
        }else {
          this.startXY = null
        }

        if(this.defaultRectWH) {
          const { w, h } = this.defaultRectWH
          this.rectW = w
          this.rectH = h
        }

        // if(this.startXY) {
        //   [this.rectX, this.rectY, this.startXY] = [this.startXY.x, this.startXY.y, null]
        //   this.startXY = null
        // }

        this._screenshot()
      },

      /**
       * 截图
       */ 
      takeScreenshot() {
        const vm = this

        // 按下
        this.cutCanvas.onmousedown = function(event) {
          vm.startXY = {
            x: event.offsetX,
            y: event.offsetY
          }
        }

        // 移动
        this.cutCanvas.onmousemove = function(event) {
          if(vm.startXY) {
            vm.rectW = event.offsetX - vm.startXY.x
            vm.rectH = event.offsetY - vm.startXY.y

            let [x, y, w, h] = [vm.startXY.x, vm.startXY.y, vm.rectW, vm.rectH]

            if(w < 0) {
              x = x + w
              w = Math.abs(w)
            }

            if(h < 0) {
              y = y + h
              h = Math.abs(h)
            }

            vm._screenshot()
          }
        }

        // 抬起
        this.cutCanvas.onmouseup = function(event) {
          vm.endXY = {
            x: event.offsetX,
            y: event.offsetY
          }

          vm._originImg()
          vm._localImg()

          vm.$emit('coordinate', {
            startX: vm.startXY.x,
            startY: vm.startXY.y,
            endX: vm.endXY.x,
            endY: vm.endXY.y,
            oImgW: vm.originImg.w,
            oImgH: vm.originImg.h,
            lImgW: vm.localImg.w,
            lImgH: vm.localImg.h
          })
        }

        document.addEventListener('mouseup', function() {
          if(vm.startXY) {
            // [vm.rectX, vm.rectY, vm.startXY] = [vm.startXY.x, vm.startXY.y, null]
            vm.startXY = null
          }
        })
      },

      /**
       * @description 定义如何截图
       */ 
      _screenshot() {
        if(!this.startXY) return

        this.cutCtx.fillStyle = 'rgba(0, 0, 0, .6)'
        this.cutCtx.strokeStyle = '#1976d2'

        this.cutCtx.clearRect(0, 0, this.width, this.height)

        this.cutCtx.beginPath()

        // 遮罩
        this.cutCtx.globalCompositeOperation = 'source-over'
        this.cutCtx.fillRect(0, 0, this.width, this.height)

        //画框
        this.cutCtx.globalCompositeOperation = 'destination-out'
        this.cutCtx.fillRect(this.startXY.x, this.startXY.y, this.rectW, this.rectH)

        //描边
        this.cutCtx.globalCompositeOperation = 'source-over'
        this.cutCtx.moveTo(this.startXY.x, this.startXY.y)
        this.cutCtx.lineWidth = 2
        this.cutCtx.lineTo(this.startXY.x + this.rectW, this.startXY.y)
        this.cutCtx.lineTo(this.startXY.x + this.rectW, this.startXY.y + this.rectH)
        this.cutCtx.lineTo(this.startXY.x, this.startXY.y + this.rectH)
        this.cutCtx.lineTo(this.startXY.x, this.startXY.y)
        this.cutCtx.stroke()
        this.cutCtx.closePath()
      },

      /**
       * @description 设置默认画布
       */ 
      _setDefaultCanvas() {
        this.$nextTick(() => {
          const img = document.getElementById('canvasImg')

          img.onload = (event) => {
            const { width, height } = event.target

            this.width = width
            this.height = height

            this.cutCanvas = document.getElementById('cutCanvas')
            this.cutCanvas.width = this.width
            this.cutCanvas.height = this.height
            this.cutCtx = this.cutCanvas.getContext('2d')

            this.init()

            this.takeScreenshot()

            this._originImg()
            this._localImg()

            this.$emit('coordinate', {
              startX: this.startXY ? this.startXY.x : 0,
              startY: this.startXY ? this.startXY.y : 0,
              endX: this.startXY ? this.startXY.x + this.rectW : 0,
              endY: this.startXY ? this.startXY.y + this.rectH : 0,
              oImgW: this.originImg.w,
              oImgH: this.originImg.h,
              lImgW: this.localImg.w,
              lImgH: this.localImg.h
            })

            this.startXY = null
          }
        })
      },

      /**
       * @description 源图片信息
       */ 
      _originImg() {
        const img = new Image()
        img.src = this.imgAddress
        this.originImg = {
          w: img.width,
          h: img.height
        }
      },

      /**
       * @description 本地图片信息
       */ 
      _localImg() {
        const img = document.getElementById('canvasImg')
        this.localImg = {
          w: img.offsetWidth,
          h: img.offsetHeight
        }
      }
    },

    watch: {
      imgAddress: {
        handler(addr) {
          console.log(addr)
        },
        immediate: true
      },

      update: {
        handler(update) {
          if(update) {
            this._setDefaultCanvas()
          }
        },
        immediate: true
      }
    }
  }
</script>

<style lang="scss">
  .screenshot {
    position: relative;
    display: inline-block;
    box-sizing: border-box;
    

    &>img {
      display: block;
    }

    .canvas {
      #canvas {
        position: absolute;
        left: 0;
        top: 0;
        z-index: 1;
      }

      #cutCanvas {
        position: absolute;
        left: 0;
        top: 0;
        z-index: 2;
      }
    }
  }
</style>