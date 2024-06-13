// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'

// import Antd from 'ant-design-vue';
import './assets/font/iconfont.css'
// import 'ant-design-vue/dist/antd.css';


Vue.config.productionTip = false

// import VueI18n from 'vue-i18n'
// Vue.use(VueI18n) // 通过插件的形式挂载
// const i18n = new VueI18n({
//     locale: window.localStorage.getItem('localeLanguage') || 'en',    // 语言标识
//     //this.$i18n.locale // 通过切换locale的值来实现语言切换
//     messages: {
//       'zh': require('./lang/zh'),   // 中文语言包
//       'en': require('./lang/en')    // 英文语言包
//     }
// })

// import TEditor from '@/components/TEditor.vue'
// Vue.component('TEditor',TEditor)
// canvas背景插件
// import vueParticleLine from 'vue-particle-line'
// import 'vue-particle-line/dist/vue-particle-line.css'
// Vue.use(vueParticleLine)
import VueParticles from 'vue-particles'
Vue.use(VueParticles)
import VueContextMenu from 'vue-contextmenu'
Vue.use(VueContextMenu)
// import VueQuillEditor from 'vue-quill-editor';
// import 'quill/dist/quill.core.css';
// import 'quill/dist/quill.snow.css';
// import 'quill/dist/quill.bubble.css';
// Vue.use(VueQuillEditor)

import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
// 全局配置elementui的dialog不能通过点击遮罩层关闭
ElementUI.Dialog.props.closeOnClickModal.default = false
Vue.use(ElementUI);
//Vue.use(Antd);
import '@/permission'
import { store } from '@/store/index'
/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  // i18n,
  components: { App },
  template: '<App/>'
})
