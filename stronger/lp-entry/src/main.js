import Vue from 'vue'
import App from './App.vue'

import './assets/styles/reset.css'
import 'animate.css'

import store from './store'
import router from './router'
import vuetify from './plugins/vuetify'

import EventBus from "./plugins/EventBus";
import './plugins/highlight'
import './plugins/vue-rocket'
import './plugins/vxe-table'
import './plugins/vue-toasted'
import './plugins/modal'
import './filters'

import './layouts'
import './components'


import 'element-ui/lib/theme-chalk/index.css';

import { Toasted, tools } from 'vue-rocket'

const customTypes = new Map([
  [200, 'success'],
  [400, 'warning']
])

const toasted = new Toasted(customTypes, { duration: 3000 })


Vue.prototype.$EventBus = EventBus;

Vue.prototype.toasted = toasted
Vue.prototype.tools = tools

import { request } from './api/service'
Vue.prototype.request = request

import { tools as _tools, bus } from './libs/util'
Vue.prototype._tools = _tools
Vue.prototype.bus = bus

Vue.config.productionTip = false

new Vue({
  store,
  router,
  vuetify,
  render: h => h(App),
}).$mount('#app')
