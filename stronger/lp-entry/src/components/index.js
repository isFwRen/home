import Vue from 'vue'

import LPCalendar from './lp-calendar'
import LPDesensitization from './lp-desensitization'
import LPDialog from './lp-dialog'
import LPDropdown from './lp-dropdown'
import LPTabs from './lp-tabs'
import LPTooltipBtn from './lp-tooltip-btn'

const components = () => {
	Vue.component('lp-calendar', LPCalendar)
	Vue.component('lp-desensitization', LPDesensitization)
	Vue.component('lp-dialog', LPDialog)
	Vue.component('lp-dropdown', LPDropdown)
	Vue.component('lp-tabs', LPTabs)
	Vue.component('lp-tooltip-btn', LPTooltipBtn)
}

Vue.use(components)