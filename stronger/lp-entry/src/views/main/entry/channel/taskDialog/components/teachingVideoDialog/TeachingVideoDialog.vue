<template>
  <div class="view-dialog">
    <lp-dialog
      ref="dialog"
      cardTextClass="pa-0"
      :cardTextStyle="{
        height: '100vh',
        overflow: 'hidden'
      }"
      fullscreen
      :title="title"
      toolbarColor="#272727"
      @dialog="handleDialog"
    >
      <div class="z-flex main" slot="main">
        <div class="flex-grow-1 wrap">
          <video ref="video" width="100%" height="100%" controls>
            <source type="video/mp4">
            您的浏览器不支持 HTML5 video 标签。
          </video>
        </div>
      </div>
    </lp-dialog>
  </div>
</template>

<script>
  import { tools as lpTools } from '@/libs/util'
  import DialogMixins from '@/mixins/DialogMixins'

  const { baseURLApi } = lpTools.baseURL()

  export default {
    name: 'TeachingFieldRulesViewDialog',
    mixins: [DialogMixins],

    props: {
      blockCode: {
        type: String,
        required: false
      },

      proCode: {
        type: String,
        required: false
      }
    },

    data() {
      return {
        baseURLApi,
        path: ''
      }
    },

    methods: {
      async getMenuList() {
        const body = {
          blockName: this.blockCode,
          pageSize: 20,
          pageIndex: 1,
          proCode: this.proCode,
          rule: '有'
        }

        const result = await this.$store.dispatch('RULE_VIDEO_GET_LIST', body)

        if(result.code === 200) {
          this.$refs.video.src = `${ baseURLApi }${ result.data?.list?.[0].video[0].path }`
        }
      }
    },

    watch: {
      dialog: {
        handler(dialog) {
          if(dialog) {
            this.getMenuList()
          }
        },
        immediate: true
      }
    }
  }
</script>

<style scoped lang="scss">
  .main {
    padding-top: 64px;
    height: 100vh;
    box-sizing: border-box;
    border-bottom: 1px solid rgba(0, 0, 0, .2);
    background-color: #121212;
  }
</style>