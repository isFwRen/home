process.env.DURIAN_MODE = "CLIENT";
const child_process = require("node:child_process");
const fs = require("node:fs/promises");
const path = require("node:path");
const qs = require("node:querystring");
const lodash = require("lodash");
const Async = require("async");
const crypt = require("./crypt");
const logger = require("./logger");
const app = require("./app");
const Db = require("./db");
const seed = require("./seed");
const zip = require("./zip");
const httpclient = require("./httpclient");
const sleep = require("./sleep");

const {
  DURIAN_API,
  DURIAN_API_PREFIX,
  DURIAN_DIR,
  DURIAN_MODE,
  DURIAN_ON_UPDATED,
  DURIAN_ON_DOWNLOAD_FAILED,
  DURIAN_DB_SECRET,
  DURIAN_HTTP_DOWNLOAD_TIMEOUT,
  DURIAN_PEER_DOWNLOAD_PROBABILITY,
} = require("./env");

// 主动查询服务器的数据库信息，和本机对比，若有不同则同步
async function _check_update() {
  if (!app._online) {
    app._socket.off("connect", check_update).once("connect", check_update);
    return;
  }
  await Db.ready();
  logger.debug("检查更新");
  // 获取服务器的数据库列表
  const r_databases = await list_remote_databases();
  // 获取本机的数据库列表
  const l_databases = await Db.listDatabases();
  // logger.debug("本机的数据库列表", l_databases);
  // 是否有删除的数据库
  const databases_del = l_databases.filter((o) => r_databases.indexOf(o) < 0);
  // 是否有新加的数据库
  const databases_add = r_databases.filter((o) => l_databases.indexOf(o) < 0);
  // 是否有更新的数据库
  const databases_update = [];
  for (const db_name of l_databases.filter(
    (o) => r_databases.indexOf(o) >= 0
  )) {
    const l_db = await Db(db_name);
    const l_tables = await l_db.listTables();
    const l_names = l_tables.map((o) => o.name);
    // const l_totals = {};
    // for (const coll_name of l_names) {
    //   const coll = l_db.table(coll_name);
    //   if (!coll) continue;
    //   l_totals[coll_name] = await coll.countEx({});
    // }
    const r_tables = await list_remote_tables(db_name);
    const r_releases = r_tables.filter((o) => o.hash);
    const r_names = r_tables.map((o) => o.name);
    // 删除表
    const names_del = l_names.filter((o) => r_names.indexOf(o) < 0);
    // 添加表
    const tables_add = r_releases
      .map((o) => o.name)
      .filter((o) => l_names.indexOf(o) < 0)
      .map((coll_name) => lodash.find(r_releases, (o) => o.name === coll_name));
    // 更新表
    const tables_update = l_names
      .map((coll_name) => {
        if (r_names.indexOf(coll_name) < 0) return;
        const r_info = lodash.find(r_tables, (o) => o.name === coll_name);
        if (r_info.released !== 1) return;
        const l_info = lodash.find(l_tables, (o) => o.name === coll_name);
        if (r_info.hash !== l_info.hash) return r_info;
        if (r_info.total !== l_info.total) return r_info;
        return;
      })
      .filter((o) => o);
    if (
      names_del.length === 0 &&
      tables_add.length === 0 &&
      tables_update.length === 0
    )
      continue;
    databases_update.push({
      db_name,
      names_del,
      tables_add,
      tables_update,
    });
  }
  if (
    databases_del.length === 0 &&
    databases_add.length === 0 &&
    databases_update.length === 0
  ) {
    logger.info("无需更新");
    app._io.emit("update", {
      action: "开始更新",
      current: 0,
      total: 0,
    });
    setTimeout(() => {
      app._io.emit("update", {
        action: "结束更新",
        current: 0,
        total: 0,
      });
    }, 40);
    setTimeout(check_update, CHECK_UPDATE_INTERVAL);
    start_share();
    return;
  }
  await update_queue.pushAsync({
    reason: "自动检查",
    databases_del,
    databases_add,
    databases_update,
  });
}
async function check_update() {
  try {
    _check_update();
  } catch (err) {
    logger.error("检查更新失败", err);
  }
}
setTimeout(check_update, 2000);
// 每4小时检查更新，不必太勤，因为服务器会主动推送更新消息
const CHECK_UPDATE_INTERVAL = 4 * 3600 * 1000;

const run_each = (arr, func, limit = 8) => {
  return Async.eachLimit(arr, limit, (item, next) => {
    func(item).finally(() => setTimeout(next, 1));
  });
};

async function do_update(param) {
  // @TODO: 更新加锁
  const { databases_del, databases_add, databases_update } = param;
  let total = (databases_del || []).length + (databases_add || []).length;
  for (const { names_del, tables_add, tables_update } of databases_update ||
    []) {
    total +=
      (names_del || []).length +
      (tables_add || []).length +
      (tables_update || []).length;
  }
  let current = 0;
  logger.info("开始更新 来自", param.reason);
  // logger.debug("更新细节", param);
  app._io.emit("update", {
    action: "开始更新",
    param,
    current,
    total,
  });
  await run_each(databases_del, async (db_name) => {
    logger.info("删除数据库:", db_name);
    current += 1;
    app._io.emit("update", {
      action: "删除数据库",
      param: db_name,
      current,
      total,
    });
    await Db.dropEx(db_name);
  });
  await run_each(
    databases_add,
    async (db_name) => {
      const remote_tables = await list_remote_tables(db_name);
      if (remote_tables.filter((o) => o.hash).length === 0) {
        logger.info("忽略无发布的数据库", db_name);
        return;
      }
      logger.info("新增数据库:", db_name);
      current += 1;
      app._io.emit("update", {
        action: "新增数据库",
        param: db_name,
        current,
        total,
      });
      await download_database(db_name).catch((err) => {
        logger.error("新增数据库失败", err);
      });
    },
    2
  );
  await run_each(
    databases_update,
    async ({ db_name, names_del, tables_add, tables_update }) => {
      const l_db = await Db(db_name, { create: true });
      await run_each(names_del, async (coll_name) => {
        logger.info("删除数据表:", `${db_name}.${coll_name}`);
        current += 1;
        app._io.emit("update", {
          action: "删除数据表",
          param: `${db_name}.${coll_name}`,
          current,
          total,
        });
        await l_db.dropCollectionEx(coll_name);
      });
      await run_each(
        tables_add,
        async (info) => {
          // logger.info("新增数据表:", `${db_name}.${info.name}`);
          current += 1;
          app._io.emit("update", {
            action: "新增数据表",
            param: `${db_name}.${info.name}`,
            current,
            total,
          });
          await download_table(db_name, info);
        },
        4
      );
      await run_each(
        tables_update,
        async (info) => {
          // logger.info("更新数据表:", `${db_name}.${info.name}`);
          current += 1;
          app._io.emit("update", {
            action: "更新数据表",
            param: `${db_name}.${info.name}`,
            current,
            total,
          });
          await download_table(db_name, info);
        },
        4
      );
    },
    2
  );
  decompress_failed = 0;
  app._io.emit("update", {
    action: "结束更新",
    current: total,
    total: total,
  });
  setTimeout(check_update, CHECK_UPDATE_INTERVAL);
  start_share();
}
const update_queue = Async.queue(do_update, 1);
update_queue.drain(() => {
  logger.info("结束更新");
  // if (DURIAN_ON_UPDATED === "quit") return quit("更新完毕，结束进程");
  // if (DURIAN_ON_UPDATED === "restart") return restart("更新完毕，重启进程");
});

let shared = false;
async function start_share() {
  if (
    shared ||
    (DURIAN_MODE !== "BACKEND" && DURIAN_PEER_DOWNLOAD_PROBABILITY < 1)
  )
    return;
  shared = true;
  for (const dbname of await Db.listDatabases()) {
    try {
      const db = await Db(dbname);
      if (db) seed.load_db_shares(db);
    } catch (error) {}
  }
}

// 结束程序
function quit(msg = "结束进程") {
  logger.info(msg);
  try {
    app._server.closeAllConnections();
    app._server.close();
  } catch (error) {}
  process.exit(0);
}

// 重新运行
function restart(msg = "重启进程") {
  logger.info(msg);
  try {
    app._server.closeAllConnections();
    app._server.close();
  } catch (error) {}
  // 运行自身
  const cmd = process.argv[0];
  const prog = path.basename(cmd, ".exe").toLowerCase();
  const args = prog === "node" ? process.argv.slice(1) : process.argv.slice(2);
  child_process.spawn(cmd, args, {
    detached: true,
  });
  process.exit(0);
}

// 订阅更新消息，自动更新数据库
app._socket.only("notify", async (m) => {
  const u = {
    db_name: m.proCode,
    names_del: [],
    tables_add: [],
    tables_update: [],
  };
  if (m.action === "drop") u.names_del.push(m.name);
  else if (m.info.hash) u.tables_update.push(m.info);
  await update_queue.pushAsync({
    reason: "服务器推送",
    databases_update: [u],
  });
});

// // 退出
// app.post(`${DURIAN_API_PREFIX}/quit`, (req, res) => {
//   res.status(200).json({ status: 200, msg: "操作成功" });
//   quit();
// });
// // 重启
// app.post(`${DURIAN_API_PREFIX}/restart`, (req, res) => {
//   res.status(200).json({ status: 200, msg: "操作成功" });
//   restart();
// });

async function download_database(db_name) {
  const r_tables = (await list_remote_tables(db_name)).filter((o) => o.hash);
  if (r_tables.length === 0) {
    logger.debug("空数据库", db_name);
    return;
  }
  // 下载解压全部文件
  let i = 0;
  let meta_time = 0;
  await run_each(
    r_tables,
    async (r_info) => {
      i += 1;
      logger.debug(
        "正在获取",
        `${i}/${r_tables.length}`,
        `${db_name}.${r_info.name}`
      );
      await download_table_file({
        database: db_name,
        name: r_info.name,
        hash: r_info.hash,
      });
      const coll_file = `${DURIAN_DIR}/db/${db_name}/${r_info.name}`;
      // 切换解压的新文件
      await sleep(100);
      await fs.rm(coll_file, { force: true });
      await fs.rename(`${coll_file}.new`, coll_file);
      if (r_info.releasedAt) {
        if (meta_time < r_info.releasedAt) meta_time = r_info.releasedAt;
        const t = new Date(r_info.releasedAt);
        await fs.utimes(coll_file, t, t);
      }
    },
    4
  );
  logger.info("数据库已下载", db_name);
  // logger.debug("创建META");
  const meta_file = `${DURIAN_DIR}/db/${db_name}/META`;
  const bak_file = `${meta_file}.${Date.now()}.old`;
  try {
    await fs.rm(bak_file, { force: true });
  } catch (error) {}
  try {
    await fs.rename(meta_file, bak_file);
  } catch (error) {}
  let meta_data =
    r_tables.map((o) => JSON.stringify(o)).join("\n") +
    `
{"$$indexCreated":{"fieldName":"name","unique":true,"sparse":false}}
{"$$indexCreated":{"fieldName":"hash","unique":false,"sparse":false}}
{"$$indexCreated":{"fieldName":"released","unique":false,"sparse":false}}
`;
  if (DURIAN_DB_SECRET && DURIAN_DB_SECRET !== "") {
    meta_data =
      meta_data
        .trim()
        .split(/[\n\r]+/)
        .map((line) => crypt.encrypt(line))
        .join("\n") + "\n";
  }
  await fs.writeFile(meta_file, meta_data);
  meta_time = new Date(meta_time);
  await fs.utimes(meta_file, meta_time, meta_time);
  const l_db = await Db(db_name);
  // 要不要再检查一次数据个数
  // logger.debug("数据库已同步", db_name);
}

async function download_table(db_name, info) {
  const ident = `${db_name}.${info.name}`;
  if (!info.hash) return logger.warn("数据表尚未发布", ident);
  // 下载解压数据表文件
  await download_table_file({
    database: db_name,
    name: info.name,
    hash: info.hash,
  });
  // 存档现有数据表
  const l_db = await Db(db_name, { create: true });
  try {
    await l_db.table(info.name).archive();
  } catch (error) {}
  const coll_file = `${DURIAN_DIR}/db/${db_name}/${info.name}`;
  // 切换解压的新文件
  await sleep(100);
  try {
    await fs.rm(coll_file, { force: true });
  } catch (error) {}
  await fs.rename(`${coll_file}.new`, coll_file);
  logger.debug("更新META", `${db_name}.${info.name}`, info.header);
  const coll_meta = l_db.table("META");
  coll_meta.allowEdit = true;
  await coll_meta.removeEx(
    { name: info.name },
    {
      multi: true,
    }
  );
  await coll_meta.insertEx([info]);
  await coll_meta.compactCollectionEx();
  coll_meta.allowEdit = false;
  if (info.releasedAt) {
    const t = new Date(info.releasedAt);
    await fs.utimes(coll_file, t, t);
    await fs.utimes(`${DURIAN_DIR}/db/${db_name}/META`, t, t);
  }
  // logger.debug("数据表已同步", ident, await l_db.table(info.name).countEx({}));
  logger.debug("数据表已同步", ident);
}

let decompress_failed = 0;
// 下载文件，然后解压
async function _download_table_file(hash, opts) {
  const torrent = await seed.download(hash, opts);
  logger.debug("已取得", seed.get_ident(torrent));
  const db_name = path.basename(torrent.path);
  const db_dir = `${DURIAN_DIR}/db/${db_name}`;
  const coll_name = torrent.name.replace(/\.zip$/i, "");
  const coll_file = `${db_dir}/${coll_name}.new`;
  await sleep(100); // 等待一会，避免与病毒防护软件冲突
  try {
    fs.rm(coll_file, { force: true });
  } catch (err) {}
  try {
    await fs.mkdir(db_dir, { recursive: true });
  } catch (err) {}
  // 解压缩
  const zip_pathname = `${torrent.path}/${torrent.name}`;
  try {
    await zip.decompress(zip_pathname, coll_name, coll_file);
  } catch (err) {
    err.files = [zip_pathname, coll_file];
    decompress_failed += 1;
    logger.error(`解压缩失败${decompress_failed}次`, zip_pathname);
    if (decompress_failed >= 6) {
      logger.error(`多次解压缩失败，请立即重启进程`, zip_pathname);
      // if (DURIAN_ON_DOWNLOAD_FAILED === "quit") {
      //   quit("多次解压缩失败，结束进程");
      // } else {
      //   restart("多次解压缩失败，重启进程");
      // }
      // await sleep(500);
    }
    throw err;
  }
  return coll_file;
}

async function download_table_file(opts) {
  const { database, name, hash } = opts;
  const pathname = `${DURIAN_DIR}/seed/${database}`;
  const ident = `${database}.${name}`;
  // 停止分享此数据表的旧资源
  const old_info = seed.get_info(database, name);
  if (old_info && old_info.hash !== hash)
    await seed.stopHash(old_info.hash, { keep: true });
  // 停止分享此资源，然后准备重新下载
  await seed.stopHash(hash);
  seed.hash(hash, {
    database,
    name,
    hash,
    filename: `${name}.zip`,
    path: pathname,
    ident,
  });
  try {
    return await _download_table_file(hash, {
      ident,
      path: pathname,
      downgrade: true,
    });
  } catch (err) {
    logger.error("获取失败", ident, err.message);
    if (err.files) {
      for (const f of err.files)
        try {
          await fs.rm(f, { force: true });
        } catch (error) {}
    }
    // if (DURIAN_ON_DOWNLOAD_FAILED === "retry") {
    //   logger.debug("重试", ident);
    //   return await download_table_file(opts);
    // } else if (DURIAN_ON_DOWNLOAD_FAILED === "restart") {
    //   return restart("下载失败，重启进程");
    // } else if (DURIAN_ON_DOWNLOAD_FAILED === "quit") {
    //   return quit("下载失败，结束进程");
    // }
    logger.debug("重试", ident);
    return await download_table_file(opts);
  }
}

async function list_remote_databases(api_url = `${DURIAN_API}/info-list`) {
  const response = await httpclient.fetch({
    url: api_url,
    timeout: DURIAN_HTTP_DOWNLOAD_TIMEOUT * 1000,
    retry: 5,
  });
  const text = response.body;
  const r = JSON.parse(text);
  return r.list || [];
}

async function list_remote_tables(
  db_name,
  api_url = `${DURIAN_API}/info-list`
) {
  const response = await httpclient.fetch({
    url: `${api_url}/${qs.escape(db_name)}`,
    timeout: DURIAN_HTTP_DOWNLOAD_TIMEOUT * 1000,
    retry: 5,
  });
  const text = response.body;
  const r = JSON.parse(text);
  return r.list || [];
}
