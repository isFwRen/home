// 跳转到项目中转页面默认项目
let defaultProjects = new Proxy(['B0118', 'B0103', 'B0116', 'B0102', 'B0121', 'B0113', 'B0110', 'B0106', 'B0114', 'B0122'], {
  get: function (target, key, receiver) {
    return target[key]
  },
  set: function (target, key, value, receiver) {
    console.warn('不允许修改项目配置')
    return false
  },
});

export {
  defaultProjects
}
