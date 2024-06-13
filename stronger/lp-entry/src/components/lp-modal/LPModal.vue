<template>
  <div 
    class="modal"
    v-if="visible" 
    @click="off"
  >
    <div class="tip-dialog animate__animated" :class="[visible ? 'animate__headShake' : '']">
      <v-card>
        <v-card-title class="z-flex justify-center text-h6">
          <span v-html="title"></span>
        </v-card-title>

        <v-card-text class="text-center">
          <span v-html="content"></span>
        </v-card-text>

        <v-card-actions class="z-flex justify-center">
          <z-btn
            color="normal"
            class="mr-3"
            depressed
            small
            @click="onCancel"
          >
            取消
          </z-btn>

          <z-btn
            color="primary"
            depressed
            small
            @click="onConfirm"
          >
            确认
          </z-btn>
        </v-card-actions>
      </v-card>
    </div>
  </div>
</template>

<script>
  export default {
    name: 'Modal',

    data() {
      return {
        visible: false,
        title: '',
        content: ''
      }
    },

    methods: {
      onCancel() {
        this.$modal.cancel()
        this.onClose()
      },

      onConfirm() {
        this.$modal.confirm()
        this.onClose()
      },

      onClose() {
        this.visible = false
      },

      onOpen() {
        this.visible = true
      },

      off(event) {
        if(event.target === event.currentTarget) {
          this.visible = false
        }
      }
    }
  }
</script>

<style lang="scss">
  .modal {
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(33, 33, 33, .46);
    z-index: 202;

    .tip-dialog {
      position: relative;
      left: 50%;
      top: 50%;
      margin-left: -100px;
      margin-top: -145px;
      width: 290px;
      height: 200px;
      z-index: 203;

      button.primary {
        background-color: #1976d2 !important;
      }
    }
  }
</style>