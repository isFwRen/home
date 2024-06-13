const child_process = require("child_process");
const { app, BrowserWindow, Notification, autoUpdater, dialog, BrowserView, utilityProcess } = require("electron");
const path = require("path");
const log = require("electron-log");
const process = require("process");
const kill = require("./util/kill-port");
const { packageName, customConfigs, constKey, apiToken, durianDir, constNodeOptions } = require("./config.js");

// 开启日志
log.transports.console.level = "silly";
// 内网ip访问需要忽略ssl
app.commandLine.appendSwitch("ignore-certificate-errors");


const { trackers, isCommon, isIntranet, ip, port, constServer } = customConfigs[packageName];

log.info("管理:", isCommon, " packageName:", packageName, " ip:", ip, " port:", port, " constServer:", constServer);

function createWindow() {
  const mainWindow = new BrowserWindow({
    show: false,
  });
  mainWindow.maximize()
  mainWindow.show()

  if (process.env.CONST) {
    log.info(`打开常量调试页面`);
    mainWindow.loadFile("./const/public/index.html", { search: `api_token=${apiToken}` })
  } else {
    log.info(`打开网页：  https://${ip}:${port}?isApp=true`);
    mainWindow.loadURL(`https://${ip}:${port}?isApp=true`);
  }


  mainWindow.on('close', function (event) {
    event.preventDefault()
    dialog.showMessageBox({
      type: "info",
      title: "提示",
      message: "确认关闭?",
      buttons: ["确定", "取消"],
      cancelId: 1,
    }).then(res => {
      console.log(res);
      if (res.response == 0) {
        app.exit()
      }
    })
  })

  // Open the DevTools.
  // mainWindow.webContents.openDevTools()
}

function startConstServerEXE() {
  log.info("start grape ", app.isPackaged, process.platform);

  let url = "assets/lib/durianc.exe"
  if (app.isPackaged) {
    url = path.join(path.dirname(__dirname), url);
  }
  log.info("cmd", url);
  log.info("ip", ip);
  log.info("isIntranet", isIntranet);
  log.info("path", path.join(path.dirname(__dirname), ".durian"));

  const startConstServerProcess = child_process.spawn(url, {
    env: {
      DURIAN_HOST: ip,
      DURIAN_DB_SECRET: constKey,
      DURIAN_API_TOKEN: apiToken,
      DURIAN_TRACKERS: trackers,
      // DURIAN_PEER_DOWNLOAD_PROBABILITY: isIntranet ? 0 : 100, //常量更新 0:内网走http 100：外网走BT
      DURIAN_PEER_DOWNLOAD_PROBABILITY: 0, //常量更新 0:内网走http
      DURIAN_DIR: path.join(path.dirname(__dirname), durianDir),
      NODE_OPTIONS: constNodeOptions,
    },
  });
  startConstServerProcess.stdout.on("data", (data) => {
    log.info(`startConstServerProcess stdout: ${data}`);
  });
  startConstServerProcess.stderr.on("data", (data) => {
    log.error(`startConstServerProcess stderr: ${data}`);
  });
  startConstServerProcess.on("close", (code) => {
    startConstServer()
    log.info(`startConstServerProcess child process exited with code ${code}`);
  });
}

function startConstServer() {
  process.env.DURIAN_HOST = ip
  process.env.DURIAN_API = constServer
  process.env.DURIAN_DB_SECRET = constKey
  process.env.DURIAN_API_TOKEN = apiToken
  process.env.DURIAN_TRACKERS = trackers
  // process.env.DURIAN_PEER_DOWNLOAD_PROBABILITY= isIntranet ? 0 : 100, //常量更新 0:内网走http 100：外网走BT
  process.env.DURIAN_PEER_DOWNLOAD_PROBABILITY = 0 //常量更新 0:内网走http
  process.env.DURIAN_DIR = path.join(path.dirname(__dirname), durianDir)
  log.info("path", path.join(__dirname));
  // process.env.NODE_OPTIONS = constNodeOptions
  const child = utilityProcess.fork(path.join(__dirname, "./const/client.js"), { stdio: "pipe" })
  child.stdout.on('data', (data) => {
    log.info(`startConstServerProcess stdout: ${data}`)
  })
  child.stderr.on('data', (data) => {
    log.error(`startConstServerProcess stderr: ${data}`)
  })
  // require("./const/client.js");
  // const startConstServerProcess = child_process.fork(path.join(path.dirname(__dirname), "./const/client.js"));
  // startConstServerProcess.on("error", (data) => {
  //   log.info(`startConstServerProcess error: ${data}`);
  // });
}

app.whenReady().then(() => {
  // 开启常量服务
  if (!isCommon) {
    kill(30080, (err, stdout, stderr) => {
      if (err || stderr) {
        log.error(
          "kill 30080 err: ",
          err,
          " \n stderr:",
          stderr,
          " \n stdout:",
          stdout
        );
        app.quit();
      } else {
        startConstServer();
        log.info("kill 30080 success stdout：", stdout);
      }
    });
  }

  createWindow();

  app.on("activate", function () {
    if (BrowserWindow.getAllWindows().length === 0) createWindow();
  });
});

app.on("window-all-closed", function () {
  if (process.platform !== "darwin") app.quit();
});


// 自动更新配置
// const updateServer = 'http://192.168.202.18:8080'
// // 检查更新频率
// const updateDur = 60000
// 程序更新推送
// log.info("app info ", process.platform, app.getVersion())
// const url = `${updateServer}/update/win64/${app.getVersion()}/`
// log.info("update url ", url)
// autoUpdater.setFeedURL({ url })
// setInterval(() => {
//   autoUpdater.checkForUpdates()
// }, updateDur)
// autoUpdater.on('update-downloaded', (event, releaseNotes, releaseName) => {
//   const dialogOpts = {
//     type: 'info',
//     buttons: ['Restart', 'Later'],
//     title: 'Application Update',
//     message: process.platform === 'win32' ? releaseNotes : releaseName,
//     detail:
//       'A new version has been downloaded. Restart the application to apply the updates.'
//   }

//   dialog.showMessageBox(dialogOpts).then((returnValue) => {
//     if (returnValue.response === 0) autoUpdater.quitAndInstall()
//   })
// })
// autoUpdater.on('error', (message) => {
//   log.error('There was a problem updating the application')
//   log.error(message)
// })

//每次启动自动更新检查 更新版本 --可以根据自己方式更新，定时或者什么
// autoUpdater.setFeedURL("https://192.168.202.18:2222");
// setInterval(() => {
//   autoUpdater.checkForUpdates()
// }, updateDur)

// autoUpdater.autoDownload = false;//这个必须写成false，写成true时，我这会报没权限更新，也没清楚什么原因
// autoUpdater.on('error', (error) => {
//   log.error(["检查更新失败", error])
// })
// //当有可用更新的时候触发。 更新将自动下载。
// autoUpdater.on('update-available', (info) => {
//   log.info('检查到有更新，开始下载新版本')
//   autoUpdater.downloadUpdate()
// })
// //当没有可用更新的时候触发。
// autoUpdater.on('update-not-available', () => {
//   log.info('没有可用更新')
// })
// //在更新下载完成的时候触发。
// autoUpdater.on('update-downloaded', (res) => {
//   log.info('下载完毕！提示安装更新')
//   //dialog 想要使用，必须在BrowserWindow创建之后
//   dialog.showMessageBox({
//     title: '升级提示！',
//     message: '已为您下载最新应用，点击确定马上替换为最新版本！'
//   }).then((index) => {
//     log.info('退出应用，安装开始！')
//     //重启应用并在下载后安装更新。 它只应在发出 update-downloaded 后方可被调用。
//     autoUpdater.quitAndInstall()
//   });
// })