# duriand 常量数据库服务

## 关于测试

可在同不同的主机上分开运行 duriand 和 duriand。如果在同一台设备上运行，需要防止端口冲突。以下是一个测试的例子：
```sh
# 运行服务器节点 duriand
# 在 192.168.1.61 上启动 duriand
# 访问 http://192.168.1.61:30080/?api_token=... 可以管理常量表
# 密钥是必须的，可放在查询参数 ?api_token=... 或 http头 Authorization: Bearer ...
# 限速20k是为了降低数据同步的速度，方便测试和观察
DURIAN_HOST=192.168.1.61 \
DURIAN_UPLOAD_LIMIT=20480 \ DURIAN_API_TOKEN=Q2xpZW50cyBTSE9VTEQgbWFrZSBhdXRoZW50aWNhdGVkIHJlcXVlc3RzIHdpdGggYSBiZWFyZXIgdG9rZW4K \
node backend.js
# 或者使用编译的 duriand
DURIAN_HOST=192.168.1.61 DURIAN_UPLOAD_LIMIT=20480 DURIAN_API_TOKEN=... duriand
```

```sh
# 运行客户端节点 durianc
# 可以在同一主机另开终端运行，也可以在能访问到 http://192.168.1.61:30080 的其他主机上运行
# 注意，DURIAN_HOST 必须是服务器节点的 IP 或者域名
# 客户端的页面只能本地访问
# 访问 http://127.0.0.1:30080/ 可以查看客户端的页面，但管理功能不可用
DURIAN_HOST=192.168.1.61 node client.js
# 或者使用编译的 durianc
DURIAN_HOST=192.168.1.61 durianc
```

数据同步是全自动的，不需要任何干预。已经考虑了超时和降级的处理，代码需要更多的测试。

## 接口简介

### 导入常量 POST /sys-const/import
- 导入excel文件（xlsx,xls,csv）的第一个sheet
- 文件名去掉扩展名的部分作为表名


### 列出某项目的常量清单
- GET /sys-const/info-list/:proCode
- POST /sys-const/info-list
	JSON表单: { "proCode": "B0114" }


### 常量分页查询 POST /sys-const/page
表单 
```js
{
	proCode: "项目编码",
	name: "常量表名",
	queryNames: {查询条件}
	respNames: [返回字段]
	sort: 排序条件
	pageIndex: 从1开始的页码
	pageSize: 页尺寸
}
```

返回
```js
{
  status: 200,
  msg: "操作成功"
  total: 总数,
  pageIndex: 从1开始的页码,
  pageSize: 页尺寸,
  list:[{查到的内容}]
}
```

举例
```js
(async function(){
  const response = await fetch("/sys-const/page", {
    body: JSON.stringify({
      proCode: '项目编码',
      name: 'B0113_百年理赔_百年理赔医院代码表',
      queryNames: {
        "医院名称": {
          $regex: '/^北京/'
        }
      },
      sort:[{"序号": 1}]
      respNames: ["医院编码", "医院名称"]
    }),
    method: "POST",
    headers: { "Content-Type": "application/json",},
  });
  const docs = await response.json();
  console.debug(docs);
})()
```

## 配置
环境变量与默认值
```sh
# 服务地址与端口，请注意，同时运行服务器端和客户端，需要防止端口冲突
DURIAN_HOST=127.0.0.1
DURIAN_PORT=30080
# 存储路径名，服务器版本是 .duriand，客户端版本是 .durian
# 数据存储路径，默认在用户home目录
DURIAN_DIR=~/.durian
# 用于处理上传文件的临时路径，默认在系统临时目录
DURIAN_TEMP=$TEMP/.durian
# API前缀
DURIAN_API_PREFIX=/sys-const
# 完整的API路径，正式版本应该用 https
DURIAN_API=http://127.0.0.1:30080/sys-const
# 分页设置
DURIAN_PAGE_SIZE=10
DURIAN_PAGE_SIZE_MAX=100
# 导入文件时批处理数据个数
DURIAN_BULK=5120
# 屏幕日志时间格式
DURIAN_LOG_TIME="MM-DD HH:mm:ss"
# 数据库文件的 AES256 密钥，32字符长，默认不加密
DURIAN_DB_SECRET=
# BT设置
# 自部署的 tracker 服务端口，用于管理分布式数据，包括TCP和UDP，不能和api接口共用
DURIAN_TRACKER_PORT=30081
# 用于管理客户端集群的 dht 服务端口
DURIAN_DHT_PORT=30082
# 数据传输端口
DURIAN_TORRENT_PORT=30083
# 下载及上传限速, -1 表示不限速
DURIAN_DOWNLOAD_LIMIT=-1
DURIAN_UPLOAD_LIMIT=-1
# 用于发布种子的 tracker 列表, 用逗号隔开，可参考 https://github.com/ngosang/trackerslist
DURIAN_TRACKERS=http://bt.endpot.com:80/announce,http://tracker.openbittorrent.com:80/announce,https://1337.abcvg.info:443/announce,https://t.zerg.pw/announce,https://t1.hloli.org:443/announce,https://tr.ready4.icu:443/announce,https://tracker.cloudit.top:443/announce,https://tracker.gbitt.info:443/announce,https://tracker.imgoingto.icu:443/announce,https://tracker.ipfsscan.io:443/announce,https://tracker.kuroy.me:443/announce,https://tracker.lilithraws.org:443/announce,https://tracker.loligirl.cn:443/announce,https://tracker.tamersunion.org:443/announce,https://tracker1.520.jp:443/announce,https://tracker2.ctix.cn:443/announce,udp://exodus.desync.com:6969/announce,udp://explodie.org:6969/announce,udp://open.stealth.si:80/announce,udp://opentracker.i2p.rocks:6969/announce,udp://private.anonseed.com:6969/announce,udp://retracker01-msk-virt.corbina.net:80/announce,udp://sanincode.com:6969/announce,udp://tracker.4.babico.name.tr:3131/announce,udp://tracker.auctor.tv:6969/announce,udp://tracker.openbittorrent.com:6969/announce,udp://tracker.opentrackr.org:1337/announce,udp://tracker.theoks.net:6969/announce,udp://tracker.tiny-vps.com:6969/announce,udp://tracker.torrent.eu.org:451/announce,udp://tracker1.bt.moack.co.kr:80/announce,udp://uploads.gamecoast.net:6969/announce
```

## 打包
需要提前在各系统下编译nexe所需的node v18.12.1，然后把编译结果放在.nexe目录
```sh
# 编译 nexe 所需的文件，只需要一次，编译很耗时间，建议根据核心数添加-j并行编译参数
yarn gen-duriand-linux --build "--make=-j`nproc`"
```
之后的打包
```sh
# 生成 windows app 所需的常量数据库服务(x86)
yarn gen-durianc-win
# 生成 windows 下开发调试用的常量数据库服务(x86)
yarn gen-duriand-win
# 生成 linux 后端需要的常量数据库服务(x64)
yarn gen-duriand-linux
# 生成 macos 下开发调试用的常量数据库服务(x64)
yarn gen-duriand-mac
```

### 剥离源码
原始的编译结果是 `.nexe/18.12.1/out/Release/node` (在windows下是 `node.exe` ), 建议把此文件复制到外层:
```sh
cp .nexe/18.12.1/out/Release/node .nexe/mac-x64-18.12.1
```
这样，就不再需要源码目录 `.nexe/18.12.1` 了

### 减小文件尺寸
建议使用 upx 压缩编译结果，不建议使用 strip 命令删除可执行文件里的符号段（减不了太多，还有副作用）。
upx 需要更新到最近版本 v4.0.2，旧版使用最高压缩比时存在问题。
```sh
upx --best -k .nexe/mac-x64-18.12.1
```
