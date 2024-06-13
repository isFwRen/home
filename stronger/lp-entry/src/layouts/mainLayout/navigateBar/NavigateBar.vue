<template>
  <v-app-bar 
    app 
    color="#fcb69f"
    dark
    clipped-left 
    elevate-on-scroll
    style="z-index:7"
  >
    <v-app-bar-nav-icon @click="onDrawer"></v-app-bar-nav-icon>

    <v-avatar class="mr-4" size="52">
      <img src="@/assets/logo.png" alt="汇流" />
    </v-avatar>

    <template v-slot:img="{ props }">
      <v-img
        v-bind="props"
        gradient="to top right, rgba(19,84,122,.5), rgba(128,208,199,.8)"
      ></v-img>
    </template>

    <v-toolbar-title>珠海汇流理赔数据处理平台2.0</v-toolbar-title>

    <v-spacer></v-spacer>

    <div class="z-flex align-center mr-3 user-info">
      <v-avatar
        v-if="user.headerImg"
        color="indigo"
        size="42"
      >
        <img :src="user.headerImg" >
      </v-avatar>

      <v-icon 
        v-else
        size="26"
      >mdi-account-circle</v-icon>

      <span class="pl-1">{{ user.nickName }}</span>

      <v-icon class="pl-4">mdi-card-account-details</v-icon>

      <span class="pl-1">{{ user.authorityId }}</span>
    </div>

    <z-btn 
      icon 
      @click="onSignOut"
    >
      <v-icon>mdi-export</v-icon>
    </z-btn>
  </v-app-bar>
</template>

<script>
  import { localStorage } from 'vue-rocket'

  export default {
    name: 'NavigateBar',

    data() {
      return {
        user: {}
      }
    },

    mounted() {
      const user = localStorage.get('user')
      if(user) {
        this.user = user
      }
    },

    methods: {
      onDrawer() {
        this.$emit('drawer')
      },

      onSignOut() {
        const vm = this
        this.$modal({
          visible: true,
          title: '退出提示',
          content: `请确认是否要退出？`,
          confirm() {
            vm.clearInfo()
            location.replace(`${ location.origin }/login`)
          }
        })
      },

      clearInfo() {
        localStorage.delete(['token', 'user','secret'])
      }
    }
  }
</script>