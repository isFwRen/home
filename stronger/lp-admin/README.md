# 理赔2.0 管理系统

### 框架
 * vuetify https://vuetifyjs.com
 * vue-rocket https://zhtgithub.github.io/vue-rocket.github.io/index.html#/home
 * vxe-table https://xuliangzhan_admin.gitee.io/vxe-table/#/table/base/basic
 
### 开发工具
  * Visual Studio Code

### 代码版本控制
  * gitlab

### 目录结构（编写中...）
  ├── node_modules  存放用包管理工具下载安装的包  
  ├── public 静态资源文件夹，不会被webpack处理  
  ├── src 业务代码  
  │   ├── api 请求封装  
  │   │   ├── pocketdb.js pouchdb封装  
  │   │   └── service.js axios封装  
  │   ├── assets 静态资源文件夹  
  │   │   ├── images 图片  
  │   │   ├── styles 样式  
  │   │   └── logo.png 项目LOGO  
  │   ├── components 业务组件  
  │   │   ├── lp-calendar 日历组件  
  │   │   ├── lp-desensitization 脱敏配置组件  
  │   │   ├── lp-dialog 通过组件引入的弹窗  
  │   │   ├── lp-drawer 侧边栏  
  │   │   ├── lp-dropdown 按钮菜单  
  │   │   ├── lp-modal 直接使用js调用的弹窗组件  
  │   │   ├── lp-spinners 全屏加载  
  │   │   ├── lp-tabs 选项卡  
  │   │   └── lp-tooltip-btn 提示按钮  
  │   ├── filters 过滤器  
  │   │   └── index.js  
  │   ├── layouts 不同模块下的页面布局  
  │   │   ├── appBar 应用栏  
  │   │   ├── mainLayout 用户登录后大部分模块所使用到的布局  
  │   │   ├── normalLayout 小部分模块需要使用的布局  
  │   │   └── usageLayout 使用说明模块用到的布局  
  │   ├── libs 常用工具文件夹  
  │   │   ├── util.message.js 用户操作结果提示  
  │   │   ├── util.storage.js 本地存储(以后会抛弃，改为使用vue-rocket封装好的localStorage)  
  │   │   └── util.tools.js 一些业务常用功能的封装  
  │   ├── mixins 混入  
  │   │   ├── ButtonMixins.js 按钮复用功能  
  │   │   ├── DialogMixins.js 弹窗复用功能  
  │   │   ├── PanelsMixins.js  
  │   │   ├── SocketsConstMixins.js
  │   │   ├── SocketsMixins.js  
  │   │   └── TableMixins.js 表格复用功能  
  │   ├── plugins 第三方依赖包  
  │   │   ├── highlight.js 代码高亮插件  
  │   │   ├── modal.js 直接使用js调用的弹窗组件(@/components/lp-modal)  
  │   │   ├── socket.io.js 即时通讯  
  │   │   ├── toasted.js Toasted  
  │   │   ├── vue-rocket.js Vue-rocket  
  │   │   ├── vuetify.js Vuetify  
  │   │   ├── vxe-table.js 表格插件  
  │   │   └── watermark.js 水印  
  │   ├── router 路由配置  
  │   │   ├── entry 管理系统不需要(以后可能会移除)  
  │   │   ├── login 登录  
  │   │   ├── main 用户登录后  
  │   │   │   ├── complaint 客户投诉  
  │   │   │   ├── entry 录入通道  
  │   │   │   ├── error 错误查询  
  │   │   │   ├── home 首页  
  │   │   │   ├── rule 项目规则  
  │   │   │   ├── salary 我的工资  
  │   │   │   └── yield 产量查询  
  │   │   ├── normal 小部分模块的路由配置  
  │   │   └── usage 使用规则  
  │   ├── store 状态管理  
  │   │   ├── auth 权限  
  │   │   ├── constants 常量  
  │   │   ├── modules 登录  
  │   │   │   ├── entry  
  │   │   │   ├── case.js  
  │   │   │   ├── errorDetail.js  
  │   │   │   ├── login.js  
  │   │   │   ├── pt.js  
  │   │   │   └── staff.js  
  │   ├── views 页面  
  │   │   ├── login 登录  
  │   │   ├── main 登录后  
  │   │   │   ├── complaint 客户投诉  
  │   │   │   ├── entry 录入通道  
  │   │   │   │   └── channel 录入  
  │   │   │   ├── error 错误查询  
  │   │   │   ├── home 首页  
  │   │   │   ├── rule 项目规则  
  │   │   │   │   ├── business 业务规则  
  │   │   │   │   ├── template 报销单模板  
  │   │   │   │   └── video 视频教学  
  │   │   │   ├── salary 我的工资  
  │   │   │   └── yield 产量查询  
  │   │   ├── normal 小部分模块  
  │   │   │   └── viewImages 查看图片  
  │   │   └── usage 使用说明  
  │   │   │   ├── code 无法获取验证码  
  │   │   │   ├── forgotJobNumber 忘记工号  
  │   │   │   ├── forgotPassword 忘记密码  
  │   │   │   ├── loginGuide 如何登录  
  │   │   │   ├── registrationGuide 如何注册  
  │   │   │   └── restoreJobNumber 恢复工号  
  │   ├── App.vue 页面资源的首加载项  
  │   ├── main.js 项目入口文件  
  ├── .env.development 开发环境下的配置文件  
  ├── .env.production 生产环境下的配置文件  
  ├── .gitignore 用以设置Git的忽略规则  
  ├── babel.config.js  
  ├── package-lock.json  
  ├── package.json 项目的描述文件  
  ├── README.md 项目手册  
  └── vue.config.js 可选配置文件  