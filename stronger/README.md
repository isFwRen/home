## 总体目标

1. 更高的运行效率。操作流程、操作功能、操作细节进行优化。
2. 优化大案件录入完毕后存在无法导出的情况。
3. 录入流程优化，已达到缩短处理时间。
4. 更友好易用的管理界面。
5. 系统卡顿的优化,增加通知公告等的功能。


### 1. 相关的技术
服务器
  * nginx
  * postgres 
  * redis
  * Jenkins 自动化部署

后端
  * go 开发语言 https://go-zh.org/
  * gin 后端框架 https://github.com/gin-gonic/gin
  * casbin  权限控制 https://casbin.org/docs/en/overview
  * jwt  用于单点登录
  * zap  日志库 https://github.com/uber-go/zap
  * grom   orm框架 https://gorm.io/zh_CN/docs/dbresolver.html
  * go-redis   用于缓存 https://github.com/go-redis/redis
  * swaggo 接口文档 https://github.com/swaggo/swag
  * websocket 通知 https://pkg.go.dev/github.com/gorilla/websocket#pkg-functions
  * sqlite3 内存数据库 https://github.com/mattn/go-sqlite3
  * GF 多功能库 https://www.bookstack.cn/read/gf-v1.8/index.md
  * go-memoize 函数缓存 https://github.com/kofalt/go-memoize



前端
  * vue  https://cn.vuejs.org/
  * vuex 状态管理 https://vuex.vuejs.org/zh/guide/
  * vue-router 路由 https://router.vuejs.org/zh/
  * axios  请求 http://www.axios-js.com/
  * ant  ui库 https://www.antdv.com/docs/vue/getting-started-cn/
  * webpack  打包 https://www.webpackjs.com/
  * sass  CSS的预编译 https://sass.bootcss.com/documentation
  * element-ui  ui库 https://element.eleme.cn/
  * PouchDB 大常量管理 https://pouchdb.com/adapters.html
  * nodejs  npm  yarn

开发工具
  * golang 编辑器
    * 注释的模板在 Preferences | Editor | Live Templates 可以编辑
  * vscode 编辑器

代码版本控制
  * svn


### 2. 项目的运行模式

  每个项目分为：下载节点、录入节点、二次加载、导出和导出、任务节点
  每个节点打包入口分别为
  * 任务节点 main.go
  * 下载节点 download.go
  * 录入节点 main.go 
  * 二次加载 type_load.go
  * 导出节点 export.go
  * 回传节点 upload.go

  录入节点以端口号运行
  前端nginx端口来区分项目，请求后端api(前端相关配置文件 config/index.js && src/utils/request.js)

### 3. 任务分配
后序再补充


### 4 开发相关的约束
##### 后端

* 一定要写注释  可以配模板
* 时间访问
* 一律从后端获取系统当前时间，前端向后端发送获取时间请求
* 时间格式采用 date
* 不同的请求，响应和数据库结构体放到对应文件夹  modal下的request，resp，根目录
* 校验用工具类
* 开发新功能或修改现有功能时，遵循以下顺序
    1. 充分沟通，达成共识
    2. 编写或修订文档，确定功能定义、后端接口定义
    3. 完成后端功能
    5. 完成前端功能
* 接口封装返回
* 文件目录基本结构
```

├── main.go             //程序入口
├── config_XXX.yaml     //全局配置文件
├── go.mod              //go的包管理
├── README.md           
├── build
│   ├── myexec          //打包后可执行文件
├── config              //配置的结构体映射
│   ├── config.gogo get -u github.com/swaggo/swag/cmd/swag
├── core                //核心配置，开启服务，log配置等
│   ├── config.go
├── db                  //数据库打包备份一下
├── global              //全局变量声明
│   ├── global.go
│   └── response        //api返回的封装
├── initialize          //初始化调用
├── log                 //日志记录
├── middleware          //中间件，jwt、操作记录、casbin权限等
├── resource            //静态资源配置
├── module             //功能模块
    ├── model          //请求，返回等结构体
    ├── router         //接口路由
    │   └── ...
    ├── api             //具体api的方法
    ├── service             //操作数据库方法，dao层
└── utils               //工具类
```

* 日志配置

| 配置名         | 配置的类型 | 说明                                                         |
| -------------- | ---------- | ------------------------------------------------------------ |
| level          | string     | level的模式的详细说明,请看[zap官方文档](https://pkg.go.dev/go.uber.org/zap?tab=doc#pkg-constants) <br />info: info模式,无错误的堆栈信息,只输出信息<br />debug:debug模式,有错误的堆栈详细信息<br />warn:warn模式<br />error: error模式,有错误的堆栈详细信息<br />dpanic: dpanic模式<br />panic: panic模式<br />fatal: fatal模式<br /> |
| format         | string     | console: 控制台形式输出日志<br />json: json格式输出日志      |
| prefix         | string     | 日志的前缀                                                   |
| director       | string     | 存放日志的文件夹,修改即可,不需要手动创建                     |
| link_name      | string     | 在server目录下会生成一个link_name的[软连接文件](https://baike.baidu.com/item/%E8%BD%AF%E9%93%BE%E6%8E%A5),链接的是director配置项的最新日志文件 |
| show_line      | bool       | 显示行号, 默认为true,不建议修改                              |
| encode_level   | string     | LowercaseLevelEncoder:小写<br /> LowercaseColorLevelEncoder:小写带颜色<br />CapitalLevelEncoder: 大写<br />CapitalColorLevelEncoder: 大写带颜色 |
| stacktrace_key | string     | 堆栈的名称,即在json格式输出日志时的josn的key                 |
| log_in_console | bool       | 是否输出到控制台,默认为true                                  |






##### 前端  
* 一定要写注释
* 统一格式化代码 比如vscode prettier格式化
* 按功能分文件夹目录 尽量结构为：要不然功能多了找不到对应的js和css了

或者

    ├── login  登录文件夹
    │   ├── index.vue
    │   ├── index.js
    │   └── index.sass 

或者

    ├── login  登录文件夹
    
    │   ├── index.vue
    
    │   └── index.sass 

或者

    ├── login  登录文件夹   
    
    │   └── index.vue



### 5. 自动化脚本
配置Jenkins管理脚本
各种节点的启动、结束脚本。确保数据的保存和临时数据的清理。
数据库、节点的开机自动启动，异常终止后自动重启动。
源代码更新脚本。
各种异常报警脚本。
文件自动化清理脚本（长时间的运单，只保留原始文件，需要时重新生成切图）。



### 6. 计划任务
* 定时统计
* 定时备份