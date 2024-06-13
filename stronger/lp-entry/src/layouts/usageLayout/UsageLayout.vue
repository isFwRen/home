<template>
  <div class="usage-layout">
    <app-bar @drawer="onDrawer"></app-bar>
    
    <lp-drawer 
      ref="drawer"
      :menus="menus"
      app
      @select="onSelect"
    >
      <div class="drawer-top" slot="top"></div>
    </lp-drawer>

    <v-main>
      <div class="pa-4 z-usage">
        <slot></slot>
      </div>
    </v-main>
  </div>
</template>

<script>
  import { localStorage } from 'vue-rocket'
  import { menus } from './cells'

  export default {
    name: 'UsageLayout',
    
    data() {
      return {
        menus
      }
    },

    methods: {
      onDrawer() {
        this.$refs.drawer.onToggle()
      },

      onSelect({ id }) {
        if(id === 'home') {
          const token = localStorage.get('token')
          
          if(!!token) {
            this.$router.push({ path: '/main/entry/channel' })
          }
          else {
            location.replace(`${ location.origin }/login`)
          }
        }
      }
    },

    components: {
      'app-bar': () => import('../appBar'),
      'lp-drawer': () => import('@/components/lp-drawer')
    }
  }
</script>

<style scoped>
  .drawer-top {
    height: 64px;
  }

  .z-usage {
    height: calc(100vh - 64px);
    overflow-y: auto;
  }
</style>