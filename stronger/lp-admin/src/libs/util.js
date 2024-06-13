import Vue from "vue";
import storage from "./util.storage";
import tools from "./util.tools";

const bus = new Vue({});

const util = {
	storage,
	tools,
	bus
};

export { storage, tools, bus };

export default util;
