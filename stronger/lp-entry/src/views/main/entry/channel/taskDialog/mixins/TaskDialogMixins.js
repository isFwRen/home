import { mapGetters } from 'vuex'
import { sessionStorage, localStorage } from 'vue-rocket'
import { tools as lpTools } from '@/libs/util'
import io from 'socket.io-client'
import { toastedOptions } from '../cells'

const isIntranet = lpTools.isIntranet()

export default {
  data() {
    return {
      socket: null
    }
  },

  mounted() {
    window.addEventListener('keydown', this.commonFuckEvents)
  },

  beforeDestroy() {
    window.removeEventListener('keydown', this.commonFuckEvents)
  },

  methods: {
    // 顶部信息
    async getTaskTopInfo() {
      this.user = localStorage.get('user') || {}

      const data = {
        code: this.user.code
      }

      const result = await this.$store.dispatch('INPUT_GET_OPS_INFO', data)
      if (!result) return
      if (result.code === 200) {
        this.$store.commit('UPDATE_CHANNEL', { topInfo: result.data })
      }
    },

    // 删单
    async handleDelConfirm({ }, form) {
      const project = localStorage.get('project')
      console.log(this.op);
      let name
      if (this.op == 'op0') {
        name = '初审'
      } else if (this.op == 'op1') {
        name = '一码'
      } else if (this.op == 'op2') {
        name = '二码'
      } else {
        name = '问题件'
      }
      const body = {
        proCode: project.code,
        id: this.bill.ID,
        delRemarks: form.delRemarks + '-' + name
      }
      console.log(body.delRemarks);
      if (form.password !== 'huiliu2022') {
        this.toasted.warning('删单密码不正确', toastedOptions)
        return
      }

      const result = await this.$store.dispatch('DELETE_CASE_ITEM', body)

      if (result.code === 200) {

        const aBody = {
          id: this.bill.ID
        }

        const aResult = await this.$store.dispatch('TASK_ALLOCATION_TASK', aBody)

        if (aResult.code == 200) {
          this.getTaskTopInfo()
          this.$refs.delDynamic.close()
          this.$refs.opRouter.getTask({ status: 'new' })
        }

        return
      }

      if (result.code !== 200) {
        this.toasted.warning(result.msg, toastedOptions)
      }
    },

    commonFuckEvents(event) {
      const { keyCode } = event || window.event

      switch (keyCode) {
        // 删单(F10)
        case 121:
          event.preventDefault()
          this.$refs.delDynamic.open({ title: '删单' })
          break;
      }
    },

    /**
     * 释放分块/工序
     * @description 若 blockId 不为空，释放分块
     * @description 若 op 不为空，释放工序
     */
    subscribeRelease({ addr, dialog }) {
      this.socket.on('connect', () => {
        console.log('连接成功!')
      })

      this.socket.on('connect_error', (error) => {
        console.log('连接错误!')
      })

      this.socket.on('connect_timeout', () => {
        console.log('连接超时!')
      })

      this.socket.on('reconnect', () => {
        console.log('重连成功!')
      })

      this.socket.on('disconnect', (error) => {
        console.log('断开连接!')
      })

      if (dialog) {
        this.socket.on('release', result => {
          const { billId, blockId, op } = result.data

          if (result.code === 200) {
            if (this.bill?.ID === billId) {
              this.redirectTask('任务已删除!')
              return
            }

            if (this.block?.ID === blockId && this.op === op) {
              this.redirectTask('任务已释放!')
              return
            }

            if (this.block?.ID === blockId && !op) {
              this.redirectTask('任务已释放!')
              return
            }

            if (!blockId && this.op === op) {
              this.redirectTask('任务已释放!')
            }
          }
        })
      }
      else {
        this.socket.disconnect()
      }
    },

    // 释放任务
    async releaseTask() {
      const body = {
        id: this.block.ID,
        op: this.op,
        code: ''
      }

      const result = await this.$store.dispatch('TASK_ALLOCATION_TASK', body)

      if (result.code === 200) {
        this.redirectTask('超时，释放分块!')
      }
    },

    // 退出当前工序
    redirectTask(message) {
      const { proCode } = this.$route.query
      this.toasted.warning(message, toastedOptions)
      this.setOpTabs()
      // this.$router.replace({ path: '/main/entry/channel', query: { op: -1, proCode } })
      this.$router.replace({ path: '/main/entry/channel' })
    }
  },

  computed: {
    ...mapGetters(['task'])
  },

  watch: {
    dialog: {
      handler(dialog) {
        if (dialog) {
          if (JSON.stringify(sessionStorage.get('task').rowInfo) == '{}') return
          const { innerIp, inAppPort, outIp, outAppPort } = sessionStorage.get('task').rowInfo
          const code = localStorage.get('user').code
          const addr = isIntranet ? `${innerIp}:${inAppPort}` : `${outIp}:${outAppPort}`
          this.socket = io(`https://${addr}/global-release?userCode=${code}`, { timeout: 5 * 60 * 1000 })
          this.subscribeRelease({ addr, dialog })
        }
      },
      immediate: true
    }
  }
}