# Grape

### 运行

```bash
npm install
npm run start
npm run package
```

### 路径

==mac==

```bash

日志
~/Library/Logs/grape/main.log

文件
/Users/mjl/Library/Application\ Support/grape/
```


==win==

```bash

日志
%USERPROFILE%\AppData\Roaming\grape\logs

文件

```

### 杀掉常量服务进程

```bash
kill `lsof -i:30080 | grep 30080 | awk '{print $2}'`
```


### iss打包成安装向导（innosetup-6.0.5）
