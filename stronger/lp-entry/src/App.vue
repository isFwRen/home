<template>
  <v-app>
    <div class="lp-app">
      <router-view></router-view>
    </div>
  </v-app>
</template>

<script>
  export default {
    name: 'App'
  }
</script>

<style lang="scss">
  html, body {
    color: #333;
    font-size: 14px;
    overflow: hidden !important;
  }

  .lp-app {
    width: 100%;
    height: 100%;
  }
</style>
