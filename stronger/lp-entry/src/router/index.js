import Vue from "vue";
import Router from "vue-router";
import { localStorage, sessionStorage } from "vue-rocket";

import routes from "./routes";

// Fix vue-router NavigationDuplicated
const originalPush = Router.prototype.push;
Router.prototype.push = function push(location) {
	return originalPush.call(this, location).catch(err => err);
};

const originalReplace = Router.prototype.replace;
Router.prototype.replace = function replace(location) {
	return originalReplace.call(this, location).catch(err => err);
};

Vue.use(Router);

const router = new Router({
	mode: "history",
	routes
});

router.afterEach((to, from) => {
	if (from.path == "/login") {
		// 修改网页标题
		let title = sessionStorage.get("proCode");
		console.log(title, to, "routerrouter-to");
		if (title == "null") document.title = "理赔2.0-录入系统";
		else document.title = `${title}-理赔2.0`;
	}
});
// 在路由导航之前修改网站标题
router.beforeEach((to, from, next) => {
	// if (to.path == '/login') {
	// const maps = new Map([
	//   // 正式测试 内网
	//   ['9000', 'B0103'],
	//   ['6700', 'B0106'],
	//   ['6900', 'B0108'],
	//   ['8800', 'B0108'],
	//   ['6600', 'B0113'],
	//   ['6500', 'B0114'],
	//   ['6200', 'B0118'],
	//   ['7800', 'B0121'],
	//   // 测试 外网
	//   ['39000', 'B0103'],
	//   ['36700', 'B0106'],
	//   ['38800', 'B0108'],
	//   ['36900', 'B0110'],
	//   ['36600', 'B0113'],
	//   ['36500', 'B0114'],
	//   ['36200', 'B0118'],
	//   ['37800', 'B0121'],
	//   // 正式 外网
	//   ['2100', 'B0108'],
	//   ['2300', 'B0113'],
	//   ['2200', 'B0114'],
	//   ['13001', 'B0118'],
	//   ['2400', 'B0121'],
	// ])
	// if (window.location.href.indexOf('m') != -1) {
	//   // 外网
	//   let flag = window.location.href.indexOf('m') + 2
	//   let port1 = window.location.href.slice(flag, flag + 4)
	//   let port2 = window.location.href.slice(flag, flag + 5)
	//   let project1 = maps.get(port1)
	//   let project2 = maps.get(port2)
	//   if(project1) document.title = `${project1}-理赔2.0`
	//   else if (project2) document.title = `${project2}-理赔2.0`
	//   else document.title = `理赔2.0-录入系统`
	// } else if (window.location.href.includes('202.18:')) {
	//   // 测试内网
	//   let flag = window.location.href.indexOf('8') + 9
	//   let port1 = window.location.href.slice(flag, flag + 4)
	//   let port2 = window.location.href.slice(flag, flag + 5)
	//   let project1 = maps.get(port1)
	//   let project2 = maps.get(port2)
	//   if (project1) document.title = `${project1}-理赔2.0`
	//   else if (project2) document.title = `${project2}-理赔2.0`
	//   else document.title = `理赔2.0-录入系统`
	// } else {
	//   // 正式内网
	//   let flag = window.location.href.indexOf('0') + 4
	//   let port1 = window.location.href.slice(flag, flag + 4)
	//   let port2 = window.location.href.slice(flag, flag + 5)
	//   let project1 = maps.get(port1)
	//   let project2 = maps.get(port2)
	//   if (project1) document.title = `${project1}-理赔2.0`
	//   else if (project2) document.title = `${project2}-理赔2.0`
	//   else document.title = `理赔2.0-录入系统`
	// }
	// }

	if (to.path != "/login") {
		// console.log(sessionStorage.get('proCode'));

		let title = sessionStorage.get("proCode");
		if (title == "null") document.title = "理赔2.0-录入系统";
		else document.title = `${title}-理赔2.0`;
	}
	if (to.path === "/download") {
		document.title = `${to.meta.title}`;
	}

	next();
});

export { router };
export default router;
