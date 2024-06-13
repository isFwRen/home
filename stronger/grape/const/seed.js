const fs = require("node:fs/promises");
const crypto = require("node:crypto");
const qs = require("node:querystring");
const path = require("node:path");
const util = require("node:util");
const Async = require("async");
const WebTorrent = require("webtorrent");
const zip = require("./zip");
const httpclient = require("./httpclient");
const logger = require("./logger");
const {
  DURIAN_HOST,
  DURIAN_DIR,
  DURIAN_TEMP,
  DURIAN_TRACKER_PORT,
  DURIAN_TORRENT_PORT,
  DURIAN_DHT_BOOTSTRAP,
  DURIAN_DHT_PORT,
  DURIAN_DOWNLOAD_LIMIT,
  DURIAN_UPLOAD_LIMIT,
  DURIAN_TRACKERS,
  DURIAN_MODE,
  DURIAN_PEER_DOWNLOAD_PROBABILITY,
  DURIAN_PEER_DOWNLOAD_TIMEOUT,
  DURIAN_HTTP_DOWNLOAD_TIMEOUT,
  DURIAN_HTTP_DOWNLOAD_MAX,
  DURIAN_API,
} = require("./env");

const ALL_SEEDS = {};
exports.hash = (...args) => {
  const [hash, v] = args;
  if (args.length === 0) return ALL_SEEDS;
  if (args.length === 1) return ALL_SEEDS[hash];
  if (v === undefined || v === null) {
    const r = ALL_SEEDS[hash];
    delete ALL_SEEDS[hash];
    return r;
  }
  ALL_SEEDS[hash] = v;
  return v;
};
exports.get_info = (database, name) =>
  Object.values(ALL_SEEDS).find(
    (o) => o.database === database && o.name === name
  );

// 发布更新
exports.release = async function release(coll) {
  const ident = `${coll._db._name}.${coll._name}`;
  logger.info(`发布更新 ${ident}`);
  // 整理数据表
  try {
    await coll.checkIndex(true);
  } catch (error) {
    logger.error(error);
  }
  await coll.compactCollectionEx();
  // 停止旧版数据表的BT分享
  await exports.stopShare(coll._db, coll._name);
  // 删除旧文件
  const pathname_zip = `${DURIAN_DIR}/seed/${coll._db._name}/${coll._name}.zip`;
  for (const f of [pathname_zip, `${pathname_zip}.torrent`]) {
    try {
      await fs.rm(f, { force: true });
    } catch (error) {}
  }
  // 生成新的.zip和.torrent
  try {
    await fs.mkdir(`${DURIAN_DIR}/seed/${coll._db._name}`, { recursive: true });
  } catch (error) {}
  await zip.compress([`${coll._db._dir}/${coll._name}`], pathname_zip);
  // 开始分享新的数据表文件
  const torrent = await startShare(pathname_zip, {
    // 注释
    comment: JSON.stringify({
      database: coll._db._name,
    }),
  });
  // 更新 meta
  const coll_meta = coll._db.table("META");
  const $set = {
    hash: torrent.infoHash,
    magnet: torrent.magnetURI,
    total: await coll.countEx({}),
    releasedAt: Date.now(),
    released: 1,
  };
  await coll_meta.updateEx({ name: coll._name }, { $set });
  await coll_meta.compactCollectionEx();
  // 推送更新消息
  return torrent;
};

// 停止旧版数据表的BT分享
exports.stopShare = async function (db, coll_name) {
  const info = await db.getTableMeta(coll_name);
  if (!info || !info.hash) return;
  await exports.stopHash(info.hash, { keep: true });
};

// 停止正在下载或分享的资源
exports.stopHash = async function (hash, opts = {}) {
  const info = ALL_SEEDS[hash] || { ident: hash };
  delete ALL_SEEDS[hash];
  if (seed_peer && seed_peer.get(hash)) {
    if (!opts.silence) logger.debug("BT取消", info.ident);
    try {
      await seed_peer.removeEx(hash);
    } catch (err) {}
  }
  if (!info.filename) return;
  const files = {
    zip: `${info.path}/${info.filename}`,
    torrent: `${info.path}/${info.filename}.torrent`,
  };
  if (opts.keep === false) {
    for (const f of Object.values(files)) {
      try {
        await fs.rm(f, { force: true });
      } catch (err) {}
    }
    return;
  }
  if (DURIAN_MODE !== "BACKEND" && DURIAN_PEER_DOWNLOAD_PROBABILITY < 1) return;
  // 旧文件移动到tmp继续分享，防止其他节点无法完成更新
  const tmp_dir = `${DURIAN_TEMP}/${info.hash}`;
  const tmp_files = {
    zip: `${tmp_dir}/${info.filename}`,
    torrent: `${tmp_dir}/${info.filename}.torrent`,
  };
  try {
    await fs.rm(tmp_dir, { recursive: true, force: true });
  } catch (error) {}
  await fs.mkdir(tmp_dir, { recursive: true });
  for (const k of Object.keys(files)) {
    try {
      await fs.rename(files[k], tmp_files[k]);
    } catch (err) {
      if (err.code === "EXDEV") {
        try {
          await fs.copyFile(files[k], tmp_files[k]);
        } catch (error) {
          var failed = true;
          // logger.warn("分享旧版文件失败", files[k]);
        } finally {
          await fs.rm(files[k], { force: true });
        }
        if (failed) return;
      }
    }
  }
  const db = await require("./db")(info.database);
  if (db)
    exports.share_table({
      hash: info.hash,
      db,
      coll_name: info.name,
      path: tmp_dir,
      autofix: false,
    });
  return;
};

// BT下载
function _bt_download(hash, opts, callback) {
  if (!seed_peer) return _http_download(hash, opts, callback);
  const timeout_check = opts.downgrade
    ? setTimeout(async () => {
        const ident = get_ident(torrent);
        logger.debug("BT超时", ident);
        _http_download(hash, opts, callback);
      }, DURIAN_PEER_DOWNLOAD_TIMEOUT * 1000)
    : null;
  const torrent = seed_peer.add(hash, {
    announce: DURIAN_TRACKERS,
    ...opts,
  });
  torrent.once("done", async () => {
    // @TODO 验证下载结果
    clearTimeout(timeout_check);
    await fs.writeFile(
      `${torrent.path}/${torrent.name}.torrent`,
      torrent.torrentFile
    );
    callback(null, torrent);
  });
  watch_torrent(torrent);
}
const bt_download = util.promisify(_bt_download);

const run_each = (arr, func, limit = 8) => {
  return Async.eachLimit(arr, limit, (item, next) => {
    func(item).finally(() => setTimeout(next, 1));
  });
};

// http下载模式
async function http_download(hash, opts) {
  const info = ALL_SEEDS[hash];
  if (!info) throw new Error(`未知资源: ${hash}`);
  await exports.stopHash(hash, { silence: true });
  await run_each(
    [
      {
        url: `${DURIAN_API}/release/${qs.escape(info.database)}/${qs.escape(
          info.filename
        )}.torrent`,
        pathname: `${info.path}/${info.filename}.torrent`,
      },
      {
        url: `${DURIAN_API}/release/${qs.escape(info.database)}/${qs.escape(
          info.filename
        )}`,
        pathname: `${info.path}/${info.filename}`,
      },
    ],
    async (param) => {
      try {
        await fs.rm(param.pathname, { force: true });
      } catch (error) {}
      while (true) {
        try {
          await http_download_file(param);
          break;
        } catch (err) {
          // 1      If an unknown error occurred.
          // 2      If time out occurred.
          // 5      If a download aborted because download speed was too slow.  See --lowest-speed-limit option.
          // 6      If network problem occurred.
          // 7      If  there were unfinished downloads. This error is only reported if all finished downloads were successful and there were unfinished downloads in a queue when aria2 exited by pressing Ctrl-C by an user or sending TERM or INT signal.
          if ([1, 2, 5, 6, 7].indexOf(err.code) >= 0) {
            logger.warn("HTTP重试", param);
            await fs.rm(param.pathname, { force: true });
            continue;
          }
          logger.warn("HTTP下载失败", param);
          throw err;
        }
      }
    }
  );
  // 降级下载并分享
  if (DURIAN_PEER_DOWNLOAD_PROBABILITY > 0) {
    _bt_download(
      hash,
      {
        ...opts,
        downgrade: false,
      },
      () => {
        logger.debug("分享", info.ident);
      }
    );
  }
  return {
    infoHash: hash,
    path: info.path,
    name: `${info.name}.zip`,
    torrentFile: await fs.readFile(`${info.path}/${info.filename}.torrent`),
  };
}
const _http_download = util.callbackify(http_download);

// 按照BT/HTTP的比例，随机选择下载方式
exports.download = async (...args) => {
  if (crypto.randomInt(100) < DURIAN_PEER_DOWNLOAD_PROBABILITY) {
    if (args[1].ident) logger.debug("BT获取", args[1].ident);
    return bt_download(...args);
  }
  // if (args[1].ident) logger.debug("HTTP获取", args[1].ident);
  return http_download(...args);
};

async function do_http_download({ url, pathname }) {
  await fs.mkdir(path.dirname(pathname), { recursive: true });
  try {
    await httpclient.download({
      url,
      pathname,
      timeout: DURIAN_HTTP_DOWNLOAD_TIMEOUT * 1000,
    });
  } catch (err) {
    logger.error("HTTP", err);
    throw err;
  }
}
const http_download_queue = Async.queue(
  do_http_download,
  DURIAN_HTTP_DOWNLOAD_MAX
);
http_download_queue.drain(() => {
  logger.debug("HTTP下载队列空闲");
});

async function http_download_file(param) {
  await http_download_queue.pushAsync(param);
}

// 分享文件
exports.share_table = async ({
  hash,
  db,
  coll_name,
  path: pathname,
  autofix,
}) => {
  if (!db) {
    logger.error("share_table", {
      hash,
      db,
      coll_name,
      path: pathname,
      autofix,
    });
    return;
  }
  const coll = db.table(coll_name);
  if (coll._shared) return;
  const ident = `${db._name}.${coll_name}`;
  const filename = `${coll_name}.zip`;
  if (!pathname) pathname = `${DURIAN_DIR}/seed/${db._name}`;
  const zip_pathname = `${pathname}/${filename}`;
  try {
    var buf = await fs.readFile(`${zip_pathname}.torrent`);
    // logger.debug("分享资源", ident);
    if (!ALL_SEEDS[hash])
      ALL_SEEDS[hash] = {
        database: db._name,
        name: coll_name,
        hash,
        filename,
        path: pathname,
        ident,
      };
  } catch (err) {
    if (DURIAN_MODE !== "BACKEND") {
      // logger.debug("未能分享", hash, ident);
      coll._shared = false;
    } else {
      logger.warn("未能分享", ident, err.message);
      if (autofix) {
        logger.warn("自动重新发布", ident);
        await exports.release(coll);
        coll._shared = true;
      }
    }
    return;
  }
  if (!seed_peer) return;
  const torrent = seed_peer.add(buf || hash, {
    path: pathname,
    announce: DURIAN_TRACKERS,
  });
  coll._shared = true;
  watch_torrent(torrent);
  torrent.rescanFiles((err) => {
    if (err) logger.error("BT扫描", ident, err);
    else torrent.resume();
  });
};

// 加载已有的数据库文件分享
exports.load_db_shares = async (db) => {
  if (DURIAN_MODE !== "BACKEND" && DURIAN_PEER_DOWNLOAD_PROBABILITY < 1) return;
  const infos = await db.listTables();
  await Promise.all(
    infos
      .filter((info) => info.hash)
      .map((info) => {
        return exports.share_table({
          hash: info.hash,
          db,
          coll_name: info.name,
          autofix: true,
        });
      })
  );
};

function init_tracker() {
  const TrackerServer = require("bittorrent-tracker").Server;
  const tracker = new TrackerServer({
    trustProxy: true,
    filter: (hash, params, cb) => {
      // logger.debug("tracker check", hash);
      setTimeout(() => {
        if (ALL_SEEDS[hash]) return cb(null);
        const err = new Error(`disallowed torrent ${hash}`);
        err.code = "EDISALLOWED";
        cb(err);
      }, 1000);
    },
  });
  // buggy client, bad data...
  tracker.on("warning", (err) => {
    // if (["EDISALLOWED"].indexOf(err.code) >= 0) return;
    logger.warn("节点", err.message || JSON.stringify(err));
  });
  // fatal error
  tracker.on("error", (err) =>
    logger.error("节点", err.message || JSON.stringify(err))
  );
  tracker.listen(
    DURIAN_TRACKER_PORT,
    {
      http: DURIAN_HOST,
      udp: DURIAN_HOST,
      ws: false,
    },
    () => {
      const httpAddr = tracker.http.address();
      const httpHost =
        httpAddr.address !== "::" ? httpAddr.address : "localhost";
      const httpPort = httpAddr.port;
      // const wsAddr = tracker.ws.address();
      // const wsHost = wsAddr.address !== "::" ? wsAddr.address : "localhost";
      // const wsPort = wsAddr.port;
      const udpAddr = tracker.udp.address();
      const udpHost = udpAddr.address;
      const udpPort = udpAddr.port;
      logger.info("tracker listening on", [
        `udp://${udpHost}:${udpPort}/announce`,
        `http://${httpHost}:${httpPort}/announce`,
      ]);
    }
  );
  // require("bittorrent-tracker/lib/common").EVENT_NAMES
  const tracker_evt = {
    // update: "更新",
    complete: "完成",
    // start: "开始",
    stop: "结束",
    pause: "暂停",
  };
  // action: 1,
  // transactionId: 1265694717,
  // type: 'udp',
  // info_hash: '8258647a64dfa09f6eddb86481fbe9a0ae85a704',
  // left: 0,
  // ip: '192.168.66.115',
  // key: 0,
  // numwant: 50,
  // port: 30083,
  // addr: '192.168.66.115:30083',
  // compact: 1
  Object.keys(tracker_evt).map((k) => {
    tracker.on(k, (addr, { info_hash, downloaded, uploaded }) => {
      addr = addr.replace(/^::ffff:/, "");
      if (addr === BITTORRENT_SERVER_ADDR) return;
      const ident = ALL_SEEDS[info_hash]
        ? ALL_SEEDS[info_hash].ident
        : info_hash;
      logger.debug(
        "节点",
        addr,
        tracker_evt[k],
        ident,
        "↓",
        downloaded,
        "B",
        "↑",
        uploaded,
        "B"
      );
    });
  });
}
const BITTORRENT_SERVER_ADDR =
  DURIAN_MODE === "BACKEND" ? `${DURIAN_HOST}:${DURIAN_TORRENT_PORT}` : "-";
if (DURIAN_MODE === "BACKEND") init_tracker();

const seed_peer =
  DURIAN_MODE !== "BACKEND" && DURIAN_PEER_DOWNLOAD_PROBABILITY < 1
    ? null
    : new WebTorrent({
        // blocklist: argv.blocklist,
        torrentPort: DURIAN_MODE === "BACKEND" ? DURIAN_TORRENT_PORT : 0,
        dhtPort: DURIAN_MODE === "BACKEND" ? DURIAN_DHT_PORT : 0,
        downloadLimit: DURIAN_DOWNLOAD_LIMIT,
        uploadLimit: DURIAN_UPLOAD_LIMIT,
        dht: {
          bootstrap: DURIAN_DHT_BOOTSTRAP,
        },
        // tracker: { announce: [] }
      });
if (seed_peer) {
  seed_peer.on("error", (err) => logger.error("seed", err));
  seed_peer.on("infoHash", (hash) => logger.debug("parsed", hash));
  seed_peer.removeEx = util.promisify((hash, callback) =>
    seed_peer.remove(hash, callback)
  );
}

let last_bt_server_status;
function bt_server_status() {
  const torrent_count = seed_peer.torrents.length;
  const downloadSpeed = Math.round(seed_peer.downloadSpeed);
  const uploadSpeed = Math.round(seed_peer.uploadSpeed);

  // 阻止重复打印
  const current_status = `${torrent_count};${downloadSpeed};${uploadSpeed}`;
  if (last_bt_server_status === current_status) return;
  last_bt_server_status = current_status;

  logger.debug(
    "BT状态 资源数",
    torrent_count,
    "↓",
    downloadSpeed,
    "B/S",
    "↑",
    uploadSpeed,
    "B/S"
  );
  setTimeout(bt_server_status, 60000);
}
let last_bt_client_status;
function bt_client_status() {
  const seed_count = seed_peer.torrents.filter((o) => o.done).length;
  const download_count = seed_peer.torrents.filter((o) => !o.done).length;
  const downloadSpeed = Math.round(seed_peer.downloadSpeed);
  const uploadSpeed = Math.round(seed_peer.uploadSpeed);
  let peer_count = 0;
  seed_peer.torrents.map((o) => (peer_count += o.numPeers));

  // 阻止重复打印
  const current_status = `${download_count};${seed_count};${peer_count};${downloadSpeed};${uploadSpeed}`;
  if (last_bt_client_status === current_status) return;
  last_bt_client_status = current_status;

  logger.debug(
    "BT状态 共",
    download_count,
    "分享",
    seed_count,
    "相关节点",
    peer_count,
    "↓",
    downloadSpeed,
    "B/S",
    "↑",
    uploadSpeed,
    "B/S"
  );
  // if (download_count > 0)
  //   seed_peer.torrents.map((o) => {
  //     if (o.done) {
  //       const ratio = Math.floor(o.ratio * 1000) / 10;
  //       logger.debug(
  //         "做种",
  //         // o.infoHash,
  //         exports.get_ident(o),
  //         "贡献",
  //         ratio,
  //         "%",
  //         "累计",
  //         o.uploaded,
  //         "B",
  //         "↑",
  //         Math.round(o.uploadSpeed),
  //         "B/S"
  //       );
  //     } else {
  //       const progress = Math.floor(o.progress * 1000) / 10;
  //       logger.debug(
  //         "BT",
  //         // o.infoHash,
  //         exports.get_ident(o),
  //         "进度",
  //         progress,
  //         "%",
  //         o.downloaded,
  //         "/",
  //         o.length,
  //         "速度",
  //         Math.round(o.downloadSpeed),
  //         "B/S",
  //         "剩余时间",
  //         Math.ceil(o.timeRemaining / 1000),
  //         "S"
  //       );
  //     }
  //   });
  setTimeout(bt_client_status, 60000);
}
const bt_status =
  DURIAN_MODE === "BACKEND" ? bt_server_status : bt_client_status;
if (seed_peer) setTimeout(bt_status, 5000);

const TORRENT_EVENTS = {
  // ready: "就绪",
  blockedPeer: "黑名单",
  hotswap: "热交换",
  close: "关闭",
  // done: "完成",
  // trackerAnnounce,
  // dhtAnnounce,
  // idle,
  // interested,
  // uninterested,
  // metadata,
  // peer,
  // invalidPeer,
  // download,
  // upload,
};

function get_ident(torrent) {
  const info = ALL_SEEDS[torrent.infoHash];
  if (!info) return torrent.name || torrent.infoHash;
  return info.ident || torrent.name;
}
exports.get_ident = get_ident;

function watch_torrent(torrent) {
  torrent.on("_infoHash", (hash) => {
    if (ALL_SEEDS[hash]) return;
    if (!torrent.name) return;
    const database = path.basename(torrent.path);
    const name = torrent.name.replace(/\.[^\.]*$/, "");
    ident = `${database}.${name}`;
    ALL_SEEDS[hash] = {
      database,
      name,
      hash: torrent.infoHash,
      filename: torrent.name,
      path: torrent.path,
      ident,
    };
  });
  // Object.keys(TORRENT_EVENTS).map((k) => {
  //   const event = TORRENT_EVENTS[k];
  //   torrent.on(k, (...args) =>
  //     logger.debug(
  //       `BT${event}`,
  //       get_ident(torrent),
  //       ...args
  //     )
  //   );
  // });
  const torrent_warn = (err) => {
    if (["ENOTFOUND"].indexOf(err.code) >= 0) return;
    if (["Request timed out", "No nodes to query"].indexOf(err.message) >= 0)
      return;
    if (/^tracker request timed out/.test(err.message)) return;
    if (/^Cannot add duplicate/i.test(err.message)) return;
    if (/^Invalid torrent identifier/i.test(err.message)) return;
    logger.debug(
      "BT状态",
      get_ident(torrent),
      err.message || JSON.stringify(err)
    );
  };
  torrent.on("warning", torrent_warn);
  torrent.on("error", torrent_warn);
  /*
  torrent._web = torrent.createServer();
  torrent._web.listen(0).on("error", (err) => {
    if (err.code === "EADDRINUSE" || err.code === "EACCES") {
      torrent._web.close();
      torrent._web.listen(0);
      return torrent._web;
    } else
      return logger.error(
        "torrent",
        torrent.infoHash,
        torrent.name,
        err.message || JSON.stringify(err)
      );
  });
  torrent._web.once("connection", () => (serving = true));
  torrent._web.once("listening", () => {
    if (torrent.ready) onReady();
    else torrent.once("ready", onReady);
  });
  async function onReady() {
    logger.debug(
      "local url",
      `http://127.0.0.1:${torrent._web.address().port}/0/${torrent.name}`
    );
  }
*/
}

function _startShare(pathname, opts, callback) {
  const db_name = path.basename(path.dirname(pathname));
  const coll_name = path.basename(pathname, ".zip");
  logger.debug("BT创建", `${db_name}.${coll_name}`);
  const torrent = seed_peer.seed(pathname, {
    announce: DURIAN_TRACKERS,
    ...opts,
  });
  watch_torrent(torrent);
  torrent.once("done", async () => {
    if (torrent.torrentFile) {
      const torrent_pathname = `${pathname}.torrent`;
      await fs.writeFile(torrent_pathname, torrent.torrentFile);
    }
    callback(null, torrent);
  });
}
const startShare = util.promisify(_startShare);
