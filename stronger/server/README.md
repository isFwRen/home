```bash
go build  -o ./build/myexec ./main.go 
./build/myexec 11111 (端口)

#或者
go run ./main.go 11111

指定配置路径  core/config.go

go mod tidy
```

接口文档 /router/api

```bash
go get -u github.com/swaggo/swag/cmd/swag
swag init

swag自动生成接口文档报错ParseComment error的解决方法之一： 
检查swag版本，swag -v，如果大于1.6.7，go get -u github.com/swaggo/swag/cmd/swag@v1.6.7，获取1.6.7版本

swag init --parseDependency --parseInternal -o ../docs

以上废弃

go install github.com/swaggo/swag/cmd/swag@latest
swag init --parseDependency -g ./entry/main.go -o ./docs
```

```go
@Summary 是对该接口的一个描述，格式：案件列表--获取案件列表（根模块--具体接口功能说明）
@Id 是一个全局标识符，所有的接口文档中 Id 不能标注
@Tags 是对接口的标注，同一个 tag 为一组，这样方便我们整理接口 格式：案件列表--获取案件列表（根模块--子模块）
@Version 表明该接口的版本
@Accept 表示该该请求的请求类型
@Param 表示参数 分别有以下参数 参数名词 参数类型 数据类型 是否必须 注释 属性(可选参数), 参数之间用空格隔开。
@Success 表示请求成功后返回，它有以下参数 请求返回状态码，参数类型，数据类型，注释
@Failure 请求失败后返回，参数同上
@Router 该函数定义了请求路由并且包含路由的请求方式。

// func name                                     //函数名称
// @Tags 用户信息                                 //swagger API分类标签, 同一个tag为一组 格式：案件列表（根模块名称）
// @Summary 案件列表--获取案件列表                  //接口概要说明 与系统sys_api挂钩的，格式：案件列表--获取案件列表（根模块--具体接口功能说明）
// @Description                                 //接口详细描述信息
// @accept json                                 //浏览器可处理数据类型，浏览器默认发 Accept: */*
// @Produce  json                               //设置返回数据的类型和编码
// @Param id path int true "ID"                 //url参数：（name；参数类型[query(?id=),path(/123)]；数据类型；required；参数描述）
// @Param data body model.ForgetPWD true "注释"  //放body的参数
// @Success 200 {object} response.Response      //成功返回的数据结构， 最后是示例
// @Failure 400 {object} response.Response
// @Router /test/{id} [get]                     //路由信息，一定要写上

```

文件目录基本结构

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

```Excel文件处理库官方文档地址```
``https://xuri.me/excelize/zh-hans/``

```CouchDB api库kivik的Githup地址```
``https://github.com/go-kivik/kivik``

```kivik的接口文档```
``https://godoc.org/github.com/vladimiroff/kivik``

1. 上传文件统一接口，上传返回路径
2. 导出excel统一接口（统一工具类）

# 接口风格

【GET】 模块/实体类/列表 /pro-config/sys-project/list

【GET】 模块/实体类/分页列表 /pro-config/sys-project/page?pageIndex=1&pageSize=10&keyword=...

【GET】 模块/实体类/详情 /pro-config/sys-project/info/{id}

【POST】 模块/实体类/换位置 /pro-config/sys-project/change

【DELETE】 模块/实体类/删除 /pro-config/sys-project/delete/ {ids} 放body

【POST】 模块/实体类/新增 /pro-config/sys-project/add

【POST】 模块/实体类/修改 /pro-config/sys-project/edit 路径全小写 用-分割 参数驼峰 首字母小写

# **编码规范**

## 命名准则

- 当变量名称在定义和最后依次使用之间的距离很短时，简短的名称看起来会更好。
- 变量命名应尽量描述其内容，而不是类型
- 常量命名应尽量描述其值，而不是如何使用这个值
- 在遇到for，if等循环或分支时，推荐单个字母命名来标识参数和返回值
- method、interface、type、package推荐使用单词命名
- package名称也是命名的一部分，请尽量将其利用起来
- 使用一致的命名风格

## 文件命名规范

- 全部小写加下划线 or 不要下划线
- 文件名称不宜过长

## 变量命名规范参考

- 首字母小写
- 驼峰命名
- 见名知义，避免拼音替代英文
- 不建议包含下划线(_)
- 不建议包含数字

**适用范围**

- 局部变量
- 函数出参、入参

## 函数、常量命名规范

- 驼峰式命名
- 可exported的必须首字母大写
- 不可exported的必须首字母小写
- 避免全部大写与下划线(_)组合

# 字段要求

    //module/sys_base/model/sys_model.go  BasePageInfo BaseTimeRange(时间范围搜索)
    分页用 pageSize  pageIndex 排序orderBy：JSON.stringify([["CreatedAt","desc"]])

    //module/sys_base/model/request/sys_common.go  ReqIds
    删除用ids []string  

    //module/sys_base/model/sys_model.go  Model
	CreatedAt time.Time
	UpdatedAt time.Time
    id string 用utils.GWorker.NextId()获取 新建时候已经有做配置
     
    //module/sys_base/model/response/common.go BasePageResult
    分页返回用这个结构体封装

    //module/sys_base/model/my_time.go MyTime  废弃  需要也可以用
    时间需要格式化可以用这个  2006-01-02 15:04:05

# 时间用 timestamptz

# 使用要求

使用中间件捕获异常，会统一处理返回。 Use(middleware.GinRecovery(false))该中间件会捕获panic的异常
request入口处有全局的panic异常处理(GinRecovery)
，只可以可以保证主程序不会因panic宕机，每一个并发程序中都需要加上该异常处理,即开协程记得加上捕获意想不到的异常

# 日志使用 
```
不要使用fmt.Println()，这个日志文件不会记录的
global.GLog.Info("")
global.GLog.Error("")
```