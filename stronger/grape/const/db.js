const fs = require("node:fs/promises");
const util = require("node:util");
const Async = require("async");
const lodash = require("lodash");
const logger = require("./logger");
const engine_nedb = require("./db-nedb");
const seed = require("./seed");
const sleep = require("./sleep");
const { DURIAN_DIR, DURIAN_MODE, DURIAN_DB_LOAD_MAX } = require("./env");
const databases = {};

// 加载数据库比较花费时间，用队列防止在加载某个库的中途再次加载
async function loadDatabase(dbname, options = {}) {
  dbname = dbname.toUpperCase();
  if (!databases[dbname])
    databases[dbname] = {
      queue: Async.queue(async (task) => {
        return await _loadDatabase(task.dbname, task.options);
      }),
    };
  if (databases[dbname].db) return databases[dbname].db;
  // logger.debug("加载数据库", dbname);
  const db = await databases[dbname].queue.pushAsync({
    dbname,
    options,
  });
  databases[dbname].db = db;
  // seed.load_db_shares(db);
  // logger.debug("database loaded:", dbname);
  return db;
}

// 最近加载的数据库
const recent_dbs = [];

async function _loadDatabase(dbname, options = {}) {
  if (!dbname) return;
  if ([".", ".."].indexOf(dbname) >= 0) return;
  if (/[\/\\:\*\?<>|]/.test(dbname)) return;
  if (databases[dbname].db) return databases[dbname].db;

  // 处理最近加载的数据库
  lodash.remove(recent_dbs, (o) => o === dbname);
  recent_dbs.push(dbname);
  // 只保留最近的两个数据库
  while (recent_dbs.length > DURIAN_DB_LOAD_MAX) {
    const dbn = recent_dbs.shift();
    if (databases[dbn].db) {
      // logger.debug("释放", dbn);
      databases[dbn].db._colls = {};
      delete databases[dbn].db;
    }
  }

  const proj_dir = `${DURIAN_DIR}/db/${dbname}`;
  const mk = options.create ? proj_dir : `${DURIAN_DIR}/db`;
  try {
    await fs.mkdir(mk, { recursive: true });
  } catch (error) {}
  try {
    const stats = await fs.stat(proj_dir);
    if (!stats.isDirectory()) return;
  } catch (error) {
    return;
  }
  const db = await engine_nedb(dbname, options);
  databases[dbname].db = db;
  databases[dbname].load_at = Date.now();
  return db;
}

let loaded = false;
function _ready(callback) {
  if (loaded) return callback();
  setTimeout(() => {
    if (loaded) {
      return callback();
    }
    _ready(callback);
  }, 1000);
}

loadDatabase.ready = util.promisify(_ready);
loadDatabase.waitReady = (req, res, next) => _ready(next);
loadDatabase.listDatabases = async () => {
  const db_dir = `${DURIAN_DIR}/db`;
  await fs.mkdir(db_dir, { recursive: true });
  return await fs.readdir(db_dir);
};
loadDatabase.dropEx = async (db_name) => {
  await loadDatabase.ready();
  if (!databases[db_name])
    return await fs.rm(`${DURIAN_DIR}/db/${db_name}`, {
      force: true,
      recursive: true,
      maxRetries: 100,
    });
  const db = await loadDatabase(db_name);
  return await db.dropEx();
};

async function init_dbs() {
  let dbs = await loadDatabase.listDatabases();
  for (const dbname of dbs) {
    const proj_dir = `${DURIAN_DIR}/db/${dbname}`;
    let has_meta = false;
    try {
      const stats = await fs.stat(`${proj_dir}/META`);
      has_meta = stats.isFile() && stats.size > 128;
    } catch (err) {}
    if (!has_meta) {
      logger.warn("删除破损的数据库目录", proj_dir);
      await fs.rm(proj_dir, { recursive: true, force: true });
      continue;
    }
    if (DURIAN_MODE === "BACKEND") {
      const db = await loadDatabase(dbname);
      await sleep(1);
      seed.load_db_shares(db);
    }
  }
  loaded = true;
  logger.info("数据库服务就绪");
}
init_dbs();

module.exports = loadDatabase;
