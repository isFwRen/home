import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

import { forms } from "vue-rocket";

// modules


const files = require.context("./modules", true, /\.js$/);
const modules = {};

files.keys().map(key => {
	modules[key.replace(/(\.\/|\.js)/g, "")] = files(key).default;
});
// auth
import auth from "./auth/auth";

// constants
import constants from "./constants/constants";

import userInfo from './userInfo/userInfo.js'

export default new Vuex.Store({
	...forms,
	modules: {
		...modules,
		auth,
		constants,
		userInfo
	}
});
