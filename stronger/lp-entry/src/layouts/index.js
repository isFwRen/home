import Vue from 'vue'

import mainLayout from './mainLayout'
import usageLayout from './usageLayout'
import normalLayout from './normalLayout'

const layouts = () => {
	Vue.component('main-layout', mainLayout)
	Vue.component('usage-layout', usageLayout)
	Vue.component('normal-layout', normalLayout)
}

Vue.use(layouts)