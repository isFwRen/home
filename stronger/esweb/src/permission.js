import router from "./router";
import { store } from "@/store/index";

import NProgress from "nprogress"; // progress bar
import "nprogress/nprogress.css"; // progress bar style
NProgress.configure({ showSpinner: false }); // NProgress Configuration

const whiteList = ["Login", "Register"]; //直接用to.meta.requireAuth
const project = "B115";

router.beforeEach(async (to, from, next) => {
  NProgress.start(); // start progress bar
  const token = store.getters["user/token"];
  const sysRoles = store.getters["user/roles"];

  // 有token
  if (token) {
      //有效token直接进去
    // if (to.path === "/") { 
    //   next({ path: "/HomePage" });
    //   NProgress.done();
    // }
    const expiresAt = store.getters["user/expiresAt"];
    const nowUnix = new Date().getTime();
    const hasExpires = expiresAt - nowUnix < 0;
    //有token
    if (hasExpires) {
      // token 过期
      await store.dispatch("user/LoginOut");
      next({ path: "/", query: { redirect: to.fullPath } });
    } else {
      // token 没过期

      if (to.meta.roles && to.meta.roles.length > 0 && sysRoles[project]) {
        var arr = new Set(sysRoles[project]);
        // 有交集
        let intersection = to.meta.roles.filter(item => arr.has(item));
        //   该路由需要权限
        if (intersection.length > 0 && intersection!=null) {
          next();
        } else {
          alert("没有权限");
        }
      } else {
        //   该路由不需要权限
        next();
      }
    }
  } else {
    //没有token
    if (!to.meta.requireAuth) {
      // 免登录
      next();
    } else {
      next({ path: "/", query: { redirect: to.fullPath } });
      NProgress.done();
    }
  }
});

router.afterEach(() => {
  NProgress.done(); // 关闭 progress bar
});
