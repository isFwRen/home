<template>
  <div class="field-rules-dialog">
    <lp-dialog
      ref="dialog"
      :title="title"
      width="700"
      @dialog="handleDialog"
    >
      <div class="pt-6" slot="main">
        <img width="100%" :src="`${ baseURLApi }${ path }`" />
      </div>

      <div class="z-flex" slot="actions">
        <z-btn
          class="mr-3"
          color="primary"
          outlined
          @click="onClose"
        >关闭</z-btn>
      </div>
    </lp-dialog>
  </div>
</template>

<script>
  import { tools as lpTools } from '@/libs/util'
  import DialogMixins from '@/mixins/DialogMixins'

  const { baseURLApi } = lpTools.baseURL()

  export default {
    name: 'FieldRulesDialog',
    mixins: [DialogMixins],

    props: {
      fieldName: {
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
      async getRules() {
        const params = {
          proCode: this.proCode,
          fieldsName: this.fieldName
        }

        const result = await this.$store.dispatch('GET_PM_TEACHING_FIELD_RULES_LIST', params)

        if(result.code === 200) {
          this.path = result.data.list[0]?.rulePicture[0]
        }
      }
    },

    watch: {
      dialog: {
        handler(dialog) {
          if(dialog) {
            this.getRules()
          }
        }
      }
    }
  }
</script>