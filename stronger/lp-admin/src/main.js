import Vue from "vue";
import App from "./App.vue";

import "./assets/styles/reset.css";
import "animate.css";
import EventBus from "./plugins/EventBus";
import TooltipPlugin from './plugins/tooltip.js';

Vue.use(TooltipPlugin);
import store from "./store";
import router from "./router";
import vuetify from "./plugins/vuetify";

import "./plugins/highlight";
import "./plugins/vue-rocket";
import "./plugins/vxe-table";
import "./plugins/toasted";
import "./plugins/vue-video-player";
import "./plugins/vue-toasted";
import "./plugins/modal";
import "./filters";

import "./layouts";
import "./components";

import { Toasted, tools, R } from "vue-rocket";

const customTypes = new Map([
	[200, "success"],
	[400, "warning"],
	[true, "success"],
	[false, "warning"]
]);

const toasted = new Toasted(customTypes, { duration: 3000 });

Vue.prototype.$EventBus = EventBus;
Vue.prototype.toasted = toasted;
Vue.prototype.tools = tools;
Vue.prototype.R = R;

import { request } from "./api/service";
Vue.prototype.request = request;

import { storage, tools as _tools, bus } from "./libs/util";
Vue.prototype.storage = storage;
Vue.prototype._tools = _tools;
Vue.prototype.bus = bus;

Vue.config.productionTip = false;

new Vue({
	store,
	router,
	vuetify,
	render: h => h(App)
}).$mount("#app");