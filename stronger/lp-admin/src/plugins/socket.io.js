import Vue from "vue";
import VueSocketIO from "vue-socket.io";
import SocketIO from "socket.io-client";
import { tools as lpTools } from "@/libs/util";
const { baseURL, baseURLApi } = lpTools.baseURL();

Vue.use(
	new VueSocketIO({
		debug: false,
		connection: SocketIO(`${baseURL}constProject`)
	})
);
